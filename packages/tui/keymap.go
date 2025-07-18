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

	// Table sorting
	SortByName          key.Binding
	SortByStatus        key.Binding
	SortByCreated       key.Binding
	SortByLastActive    key.Binding
	ToggleSortDirection key.Binding
	ClearSort           key.Binding

	// Table pagination
	NextPage         key.Binding
	PrevPage         key.Binding
	FirstPage        key.Binding
	LastPage         key.Binding
	PageSizeIncrease key.Binding
	PageSizeDecrease key.Binding

	// Table filtering
	ToggleFilter key.Binding
	ClearFilter  key.Binding
	FocusFilter  key.Binding

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

		// Table sorting
		SortByName: key.NewBinding(
			key.WithKeys("s", "n"),
			key.WithHelp("s/n", "sort by name"),
		),
		SortByStatus: key.NewBinding(
			key.WithKeys("s", "s"),
			key.WithHelp("s+s", "sort by status"),
		),
		SortByCreated: key.NewBinding(
			key.WithKeys("s", "c"),
			key.WithHelp("s+c", "sort by created"),
		),
		SortByLastActive: key.NewBinding(
			key.WithKeys("s", "a"),
			key.WithHelp("s+a", "sort by last active"),
		),
		ToggleSortDirection: key.NewBinding(
			key.WithKeys("s", "d"),
			key.WithHelp("s+d", "toggle sort direction"),
		),
		ClearSort: key.NewBinding(
			key.WithKeys("s", "x"),
			key.WithHelp("s+x", "clear sort"),
		),

		// Table pagination
		NextPage: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "next page"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "previous page"),
		),
		FirstPage: key.NewBinding(
			key.WithKeys("g", "g"),
			key.WithHelp("g+g", "first page"),
		),
		LastPage: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "last page"),
		),
		PageSizeIncrease: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "increase page size"),
		),
		PageSizeDecrease: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "decrease page size"),
		),

		// Table filtering
		ToggleFilter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "toggle filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear filter"),
		),
		FocusFilter: key.NewBinding(
			key.WithKeys("ctrl+f"),
			key.WithHelp("ctrl+f", "focus filter"),
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
		k.ToggleFilter,
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
		// Table Sorting
		{k.SortByName, k.SortByStatus, k.SortByCreated, k.SortByLastActive},
		{k.ToggleSortDirection, k.ClearSort},
		// Table Pagination
		{k.NextPage, k.PrevPage, k.FirstPage, k.LastPage},
		{k.PageSizeIncrease, k.PageSizeDecrease},
		// Table Filtering
		{k.ToggleFilter, k.ClearFilter, k.FocusFilter},
		// Table Selection
		{k.SelectAll, k.DeselectAll, k.ToggleRowSelection, k.InvertSelection},
		// Table View Options
		{k.ToggleRowNumbers, k.ToggleCompactView, k.RefreshTable},
		// Help and control
		{k.Help, k.Back, k.Quit},
	}
}
