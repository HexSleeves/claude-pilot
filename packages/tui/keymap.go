package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines keyboard shortcuts for the TUI
type KeyMap struct {
	// Navigation
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding

	// Actions
	Attach  key.Binding
	Create  key.Binding
	Kill    key.Binding
	Refresh key.Binding
	Quit    key.Binding

	// Confirmation
	Yes key.Binding
	No  key.Binding

	// Form navigation
	Submit    key.Binding
	NextInput key.Binding
	PrevInput key.Binding

	// Help and navigation
	Help key.Binding
	Back key.Binding

	// Table selection
	SelectAll          key.Binding
	DeselectAll        key.Binding
	ToggleRowSelection key.Binding
	InvertSelection    key.Binding

	// Table view options
	ToggleRowNumbers  key.Binding
	ToggleCompactView key.Binding
	RefreshTable      key.Binding
}

// DefaultKeyMap returns the default key mappings for the TUI.
// It configures all keyboard shortcuts with appropriate help text
// for navigation, actions, and form interactions.
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
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup/b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f"),
			key.WithHelp("pgdn/f", "page down"),
		),

		// Actions
		Attach: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "attach to session"),
		),
		Create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create session"),
		),
		Kill: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "kill session"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),

		// Confirmation
		Yes: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "yes"),
		),
		No: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "no"),
		),

		// Form navigation
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
		NextInput: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		PrevInput: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous field"),
		),

		// Help and navigation
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),

		// Table selection
		SelectAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
		),
		DeselectAll: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "deselect all"),
		),
		ToggleRowSelection: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle row selection"),
		),
		InvertSelection: key.NewBinding(
			key.WithKeys("ctrl+i"),
			key.WithHelp("ctrl+i", "invert selection"),
		),

		// Table view options
		ToggleRowNumbers: key.NewBinding(
			key.WithKeys("v", "n"),
			key.WithHelp("v+n", "toggle row numbers"),
		),
		ToggleCompactView: key.NewBinding(
			key.WithKeys("v", "c"),
			key.WithHelp("v+c", "toggle compact view"),
		),
		RefreshTable: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh table"),
		),
	}
}

// ShortHelp returns brief key hints for footer display
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Attach,
		k.Create,
		k.Kill,
		k.ToggleRowSelection,
		k.Help,
		k.Quit,
	}
}

// FullHelp returns detailed help text for help modal
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// Navigation
		{k.Up, k.Down, k.PageUp, k.PageDown},
		// Actions
		{k.Attach, k.Create, k.Kill, k.Refresh},
		// Confirmation
		{k.Yes, k.No},
		// Table Selection
		{k.SelectAll, k.DeselectAll, k.ToggleRowSelection, k.InvertSelection},
		// Table View Options
		{k.ToggleRowNumbers, k.ToggleCompactView, k.RefreshTable},
		// Help and control
		{k.Help, k.Back, k.Quit},
	}
}
