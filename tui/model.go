package tui

import (
	"fmt"
	"gh-actions-tui/github"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	workflowsView viewState = iota
	runsView
	jobsView
	logsView
)

type model struct {
	client           *github.Client
	repo             string
	list             list.Model
	spinner          spinner.Model
	viewport         viewport.Model
	help             help.Model
	keys             keyMap
	logsKeys         logsKeyMap
	error            error
	view             viewState
	loading          bool
	selectedWorkflow github.Workflow
	selectedRun      github.WorkflowRun
	selectedJob      github.Job
	// Search functionality
	searchMode    bool
	searchQuery   string
	searchResults []int  // line numbers of search matches
	currentMatch  int    // current match index
	logContent    string // store original log content for searching
	lastKey       string // for handling 'gg' sequence
}

type workflowsLoadedMsg []github.Workflow
type runsLoadedMsg []github.WorkflowRun
type jobsLoadedMsg []github.Job
type logsLoadedMsg string

func InitialModel(repo string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Initialize viewport with default dimensions
	vp := viewport.New(80, 24)

	logsKeys := logsKeyMap{
		Up:       keys.Up,
		Down:     keys.Down,
		PageUp:   keys.PageUp,
		PageDown: keys.PageDown,
		GoTop: key.NewBinding(
			key.WithKeys("gg"),
			key.WithHelp("gg", "go to top"),
		),
		GoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "go to bottom"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		NextMatch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next match"),
		),
		PrevMatch: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "prev match"),
		),
		Back: keys.Back,
		Quit: keys.Quit,
	}

	m := model{
		client:   github.NewClient(),
		repo:     repo,
		list:     list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		spinner:  s,
		viewport: vp,
		help:     help.New(),
		keys:     keys,
		logsKeys: logsKeys,
		view:     workflowsView,
		loading:  true,
	}
	m.list.Title = "GitHub Workflows"
	m.updateKeybindings()
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		workflows, err := m.client.ListWorkflows(m.repo)
		if err != nil {
			return err
		}
		return workflowsLoadedMsg(workflows)
	})
}

func (m *model) updateKeybindings() {
	m.keys.Enter.SetEnabled(m.view != logsView)
	m.keys.Back.SetEnabled(m.view > workflowsView)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle navigation keys with explicit key checking to prevent child components from consuming them
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.view == logsView {
				// Navigate back from logs to jobs
				m.view = jobsView
				m.list.Title = "Jobs for Run #" + fmt.Sprintf("%d", m.selectedRun.Number)
				m.viewport.SetContent("") // Clear viewport content
				m.loading = true
				m.updateKeybindings()
				return m, func() tea.Msg {
					jobs, err := m.client.ListJobs(m.repo, m.selectedRun.ID)
					if err != nil {
						return err
					}
					return jobsLoadedMsg(jobs)
				}
			}
			if m.view > workflowsView {
				m.loading = true
				m.view--
				m.updateKeybindings()
				switch m.view {
				case workflowsView:
					m.list.Title = "GitHub Workflows"
					return m, func() tea.Msg {
						workflows, err := m.client.ListWorkflows(m.repo)
						if err != nil {
							return err
						}
						return workflowsLoadedMsg(workflows)
					}
				case runsView:
					m.list.Title = "Workflow Runs for " + m.selectedWorkflow.Name
					return m, func() tea.Msg {
						runs, err := m.client.ListWorkflowRuns(m.repo, m.selectedWorkflow.ID)
						if err != nil {
							return err
						}
						return runsLoadedMsg(runs)
					}
				case jobsView:
					m.list.Title = "Jobs for Run #" + fmt.Sprintf("%d", m.selectedRun.Number)
					return m, func() tea.Msg {
						jobs, err := m.client.ListJobs(m.repo, m.selectedRun.ID)
						if err != nil {
							return err
						}
						return jobsLoadedMsg(jobs)
					}
				}
			}
			return m, nil
		case "enter":
			if m.loading || m.view == logsView {
				return m, nil
			}
			m.loading = true
			switch m.view {
			case workflowsView:
				m.selectedWorkflow = m.list.SelectedItem().(github.Workflow)
				m.view = runsView
				m.list.Title = "Workflow Runs for " + m.selectedWorkflow.Name
				m.updateKeybindings()
				return m, func() tea.Msg {
					runs, err := m.client.ListWorkflowRuns(m.repo, m.selectedWorkflow.ID)
					if err != nil {
						return err
					}
					return runsLoadedMsg(runs)
				}
			case runsView:
				m.selectedRun = m.list.SelectedItem().(github.WorkflowRun)
				m.view = jobsView
				m.list.Title = "Jobs for Run #" + fmt.Sprintf("%d", m.selectedRun.Number)
				m.updateKeybindings()
				return m, func() tea.Msg {
					jobs, err := m.client.ListJobs(m.repo, m.selectedRun.ID)
					if err != nil {
						return err
					}
					return jobsLoadedMsg(jobs)
				}
			case jobsView:
				m.selectedJob = m.list.SelectedItem().(github.Job)
				if m.selectedJob.Conclusion == "skipped" {
					m.error = fmt.Errorf("cannot fetch logs for skipped jobs")
					return m, nil
				}
				m.view = logsView
				m.list.Title = "Logs for " + m.selectedJob.Name
				m.updateKeybindings()
				return m, func() tea.Msg {
					logs, err := m.client.GetJobLogs(m.repo, m.selectedJob.ID)
					if err != nil {
						return err
					}
					return logsLoadedMsg(logs)
				}
			}
		}

		// Handle viewport navigation keys when in logs view
		if m.view == logsView && !m.loading {
			if m.searchMode {
				// Handle search input
				switch msg.String() {
				case "enter":
					m.searchMode = false
					m.performSearch()
					return m, nil
				case "esc":
					m.searchMode = false
					m.searchQuery = ""
					return m, nil
				case "backspace":
					if len(m.searchQuery) > 0 {
						m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					}
					return m, nil
				default:
					// Add character to search query
					if len(msg.String()) == 1 {
						m.searchQuery += msg.String()
					}
					return m, nil
				}
			} else {
				// Handle navigation keys
				switch msg.String() {
				case "up", "k":
					m.viewport.ScrollUp(1)
					m.lastKey = ""
					return m, nil
				case "down", "j":
					m.viewport.ScrollDown(1)
					m.lastKey = ""
					return m, nil
				case "pgup", "b", "u":
					m.viewport.PageUp()
					m.lastKey = ""
					return m, nil
				case "pgdown", "f", "d", " ":
					m.viewport.PageDown()
					m.lastKey = ""
					return m, nil
				case "home":
					m.viewport.GotoTop()
					m.lastKey = ""
					return m, nil
				case "end":
					m.viewport.GotoBottom()
					m.lastKey = ""
					return m, nil
				case "g":
					if m.lastKey == "g" {
						// 'gg' - go to top
						m.viewport.GotoTop()
						m.lastKey = ""
					} else {
						m.lastKey = "g"
					}
					return m, nil
				case "G":
					// Shift+G - go to bottom
					m.viewport.GotoBottom()
					m.lastKey = ""
					return m, nil
				case "/":
					// Start search mode
					m.searchMode = true
					m.searchQuery = ""
					m.lastKey = ""
					return m, nil
				case "n":
					// Next search result
					m.nextSearchResult()
					m.lastKey = ""
					return m, nil
				case "N":
					// Previous search result
					m.prevSearchResult()
					m.lastKey = ""
					return m, nil
				default:
					// Reset lastKey for any other key
					m.lastKey = ""
				}
			}
		}

		// If we reach here, the key wasn't handled by global handlers
		// Pass it to the appropriate child component
		var cmd tea.Cmd
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
		} else {
			switch m.view {
			case logsView:
				m.viewport, cmd = m.viewport.Update(msg)
			default:
				m.list, cmd = m.list.Update(msg)
			}
		}
		return m, cmd

	case tea.WindowSizeMsg:
		headerHeight := 3 // Title + spacing
		footerHeight := 2 // Help text
		verticalMarginHeight := headerHeight + footerHeight

		m.list.SetSize(msg.Width, msg.Height-verticalMarginHeight)

		// Properly size viewport for logs view
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	case workflowsLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg))
		for i, w := range msg {
			items[i] = w
		}
		m.list.SetItems(items)
		m.list.Select(0)
	case runsLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg))
		for i, r := range msg {
			items[i] = r
		}
		m.list.SetItems(items)
		m.list.Select(0)
	case jobsLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg))
		for i, j := range msg {
			items[i] = j
		}
		m.list.SetItems(items)
		m.list.Select(0)
	case logsLoadedMsg:
		m.loading = false
		m.logContent = string(msg) // Store original content for searching
		m.viewport.SetContent(m.logContent)
	case error:
		m.loading = false
		m.error = msg
	}

	m.updateKeybindings()
	return m, nil
}

// Search functionality methods
func (m *model) performSearch() {
	if m.searchQuery == "" {
		return
	}

	m.searchResults = []int{}
	m.currentMatch = 0

	lines := strings.Split(m.logContent, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(m.searchQuery)) {
			m.searchResults = append(m.searchResults, i)
		}
	}

	if len(m.searchResults) > 0 {
		m.goToSearchResult(0)
	}
}

func (m *model) nextSearchResult() {
	if len(m.searchResults) == 0 {
		return
	}

	m.currentMatch = (m.currentMatch + 1) % len(m.searchResults)
	m.goToSearchResult(m.currentMatch)
}

func (m *model) prevSearchResult() {
	if len(m.searchResults) == 0 {
		return
	}

	m.currentMatch = (m.currentMatch - 1 + len(m.searchResults)) % len(m.searchResults)
	m.goToSearchResult(m.currentMatch)
}

func (m *model) goToSearchResult(index int) {
	if index >= len(m.searchResults) {
		return
	}

	lineNum := m.searchResults[index]
	// Calculate the percentage position in the content
	lines := strings.Split(m.logContent, "\n")
	if len(lines) > 0 {
		percentage := float64(lineNum) / float64(len(lines))
		m.viewport.SetYOffset(int(percentage * float64(m.viewport.TotalLineCount())))
	}
}
