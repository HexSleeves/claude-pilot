package tui

import (
	"claude-pilot/shared/interfaces"
)

// Message types for the bubbletea message-passing system

// sessionsLoadedMsg contains the result of loading sessions from the API.
// This message is sent when the loadSessionsCmd completes, either successfully
// with a list of sessions or with an error.
type sessionsLoadedMsg struct {
	sessions []*interfaces.Session
	err      error
}

// sessionCreatedMsg contains the result of creating a new session.
// This message is sent when the createSessionCmd completes, containing
// either the newly created session or an error if creation failed.
type sessionCreatedMsg struct {
	session *interfaces.Session
	err     error
}

// sessionKilledMsg contains the result of terminating a session.
// This message is sent when the killSessionCmd completes, indicating
// success or failure of the session termination operation.
type sessionKilledMsg struct {
	sessionID string
	err       error
}

// errorMsg contains error information for display in the error view.
// This message is used to transition the TUI to an error state with
// user-friendly error display and retry options.
type errorMsg struct {
	error error
}

// statusMsg contains status text for temporary display in the main view.
// This message is used to show brief status updates to the user, such as
// "Session created successfully" or "Refreshing sessions...".
type statusMsg struct {
	message string
}

// Internal state messages

// viewStateMsg triggers a change between different TUI views.
// This message is used internally to transition between TableView,
// CreatePrompt, Loading, Error, and Help states.
type viewStateMsg struct {
	state ViewState
}

// Table action messages for enhanced table operations

// TableSortMsg requests sorting of table data by a specific column and direction.
// This message is sent when the user triggers a sort operation through keyboard shortcuts
// or other UI interactions. The Direction should be "asc" or "desc".
type TableSortMsg struct {
	Column    string
	Direction string
}

// TableFilterMsg requests filtering of table data based on the provided text.
// This message is sent when the user applies a filter to the table data.
// The FilterText is used to match against various session fields.
type TableFilterMsg struct {
	FilterText string
}

// TablePageMsg requests a pagination action on the table data.
// This message is sent when the user navigates between pages or changes page size.
// Action can be "next", "prev", "first", "last", or "size_change".
type TablePageMsg struct {
	Action string
	Page   int
}

// TableSelectMsg requests a selection action on table rows.
// This message is sent when the user performs row selection operations.
// Action can be "select", "deselect", "toggle", "all", "none", or "invert".
// RowIndex is used for single-row operations (-1 for multi-row operations).
type TableSelectMsg struct {
	Action   string
	RowIndex int
}

// TableRefreshMsg requests a refresh of table data from the data source.
// This message is sent when the user explicitly requests a data refresh
// or when automatic refresh is triggered.
type TableRefreshMsg struct{}

// TableConfigMsg requests changes to table configuration settings.
// This message is sent when the user modifies table display options,
// pagination settings, or other configuration parameters.
type TableConfigMsg struct {
	Config map[string]interface{}
}

// Table state messages for operation completion notifications

// TableSortedMsg indicates that a table sort operation has completed.
// This message is sent after the table data has been successfully sorted
// and the UI should reflect the new sort state.
type TableSortedMsg struct {
	Column    string
	Direction string
}

// TableFilteredMsg indicates that a table filter operation has completed.
// This message is sent after the table data has been successfully filtered
// and the UI should reflect the new filter state.
type TableFilteredMsg struct {
	FilterText string
	RowCount   int
}

// TablePageChangedMsg indicates that a table pagination operation has completed.
// This message is sent after the table page has been successfully changed
// and the UI should reflect the new page state.
type TablePageChangedMsg struct {
	CurrentPage int
	TotalPages  int
	PageSize    int
}

// TableSelectionChangedMsg indicates that a table selection operation has completed.
// This message is sent after the table selection state has been successfully updated
// and the UI should reflect the new selection state.
type TableSelectionChangedMsg struct {
	SelectedRows []int
	Action       string
}

// TableErrorMsg indicates that a table operation has failed.
// This message is sent when any table operation encounters an error
// and provides error details for user feedback.
type TableErrorMsg struct {
	Error error
}

// Table view messages for UI state changes

// TableViewModeMsg requests a change in table display mode.
// This message is sent when the user switches between different table view modes.
// Mode can be "compact", "expanded", "detailed", or other supported modes.
type TableViewModeMsg struct {
	Mode string
}

// TableHelpToggleMsg requests toggling of table-specific help display.
// This message is sent when the user wants to show or hide table-specific
// help information and keyboard shortcuts.
type TableHelpToggleMsg struct{}

// TableResizeMsg requests responsive updates to table dimensions.
// This message is sent when the terminal window is resized or when
// the table needs to adjust its layout for optimal display.
type TableResizeMsg struct {
	Width  int
	Height int
}

// Helper constructor functions for creating table messages

// NewTableSortMsg creates a new TableSortMsg with the specified column and direction.
// Direction should be "asc" for ascending or "desc" for descending sort order.
func NewTableSortMsg(column, direction string) TableSortMsg {
	return TableSortMsg{
		Column:    column,
		Direction: direction,
	}
}

// NewTableFilterMsg creates a new TableFilterMsg with the specified filter text.
// The filter text will be used to match against session fields for filtering.
func NewTableFilterMsg(text string) TableFilterMsg {
	return TableFilterMsg{
		FilterText: text,
	}
}

// NewTablePageMsg creates a new TablePageMsg with the specified action and page number.
// Action should be one of: "next", "prev", "first", "last", or "size_change".
func NewTablePageMsg(action string, page int) TablePageMsg {
	return TablePageMsg{
		Action: action,
		Page:   page,
	}
}

// NewTableSelectMsg creates a new TableSelectMsg with the specified action and row index.
// Action should be one of: "select", "deselect", "toggle", "all", "none", or "invert".
// Use -1 for rowIndex when the action applies to multiple rows.
func NewTableSelectMsg(action string, rowIndex int) TableSelectMsg {
	return TableSelectMsg{
		Action:   action,
		RowIndex: rowIndex,
	}
}
