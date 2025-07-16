package models

import "github.com/charmbracelet/bubbles/key"

// KeyMap contains all the key bindings for the dashboard
type KeyMap struct {
	// Navigation
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding

	// Actions
	Enter   key.Binding
	Create  key.Binding
	Kill    key.Binding
	Refresh key.Binding
	Help    key.Binding

	// Search
	Search key.Binding

	// System
	Quit   key.Binding
	Escape key.Binding

	// Modal specific
	Submit key.Binding
	Cancel key.Binding

	// Detail panel
	ScrollUp   key.Binding
	ScrollDown key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Home       key.Binding
	End        key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next panel"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous panel"),
		),

		// Actions
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/attach"),
		),
		Create: key.NewBinding(
			key.WithKeys("n", "c"),
			key.WithHelp("n/c", "create session"),
		),
		Kill: key.NewBinding(
			key.WithKeys("d", "x"),
			key.WithHelp("d/x", "kill session"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),

		// Search
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),

		// System
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel/close"),
		),

		// Modal specific
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),

		// Detail panel
		ScrollUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+b"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+f"),
			key.WithHelp("pgdn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to bottom"),
		),
	}
}

// MainKeys returns the key bindings for the main dashboard
func (k KeyMap) MainKeys() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Tab, k.Enter, k.Create, k.Kill, k.Refresh, k.Search, k.Help, k.Quit,
	}
}

// ModalKeys returns the key bindings for modal dialogs
func (k KeyMap) ModalKeys() []key.Binding {
	return []key.Binding{
		k.Tab, k.ShiftTab, k.Submit, k.Cancel,
	}
}

// DetailKeys returns the key bindings for the detail panel
func (k KeyMap) DetailKeys() []key.Binding {
	return []key.Binding{
		k.ScrollUp, k.ScrollDown, k.PageUp, k.PageDown, k.Home, k.End, k.Escape,
	}
}

// Implement help.KeyMap interface
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Tab, k.Enter, k.Create, k.Kill, k.Refresh, k.Search, k.Help, k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab, k.Enter},
		{k.Create, k.Kill, k.Refresh, k.Help},
		{k.Quit, k.Escape},
	}
}

// TableKeys returns the key bindings for table navigation
func (k KeyMap) TableKeys() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Enter, k.Create, k.Kill, k.Refresh, k.Help, k.Quit,
	}
}
