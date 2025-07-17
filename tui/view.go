package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
)

func (m model) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.error)
	}

	var s string
	if m.loading {
		s = fmt.Sprintf("\n   %s %s...\n\n", m.spinner.View(), m.list.Title)
	} else {
		switch m.view {
		case logsView:
			s = m.viewport.View()
		default:
			s = m.list.View()
		}
	}

	// Show different help based on current view
	var helpView string
	if m.view == logsView {
		if m.searchMode {
			// Show search input
			searchPrompt := fmt.Sprintf("Search: %s", m.searchQuery)
			helpView = searchPrompt
		} else {
			// Show search status if there are results
			if len(m.searchResults) > 0 {
				searchStatus := fmt.Sprintf("(%d/%d matches) ", m.currentMatch+1, len(m.searchResults))
				helpView = searchStatus + m.help.View(m.logsKeys)
			} else if m.searchQuery != "" {
				helpView = "No matches found • " + m.help.View(m.logsKeys)
			} else {
				helpView = m.help.View(m.logsKeys)
			}
		}
	} else {
		helpView = m.help.View(m.keys)
	}

	return headerStyle.Render(m.list.Title) + "\n" + s + "\n" + helpView
}