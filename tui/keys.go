package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Up    key.Binding
	Down  key.Binding
	PageUp key.Binding
	PageDown key.Binding
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "b", "u"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "f", "d"),
		key.WithHelp("pgdn", "page down"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

type logsKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	GoTop    key.Binding
	GoBottom key.Binding
	Search   key.Binding
	NextMatch key.Binding
	PrevMatch key.Binding
	Back     key.Binding
	Quit     key.Binding
}

func (k logsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.GoTop, k.GoBottom, k.Search, k.Back, k.Quit}
}

func (k logsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.GoTop, k.GoBottom, k.Search, k.NextMatch, k.PrevMatch},
		{k.Back, k.Quit},
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Quit},
	}
}
