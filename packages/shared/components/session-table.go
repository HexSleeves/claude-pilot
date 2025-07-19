package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

// Column key constants for evertras table
const (
	columnKeyID         = "id"
	columnKeyName       = "name"
	columnKeyStatus     = "status"
	columnKeyBackend    = "backend"
	columnKeyCreated    = "created"
	columnKeyLastActive = "last_active"
	columnKeyProject    = "project"
	columnKeyMessages   = "messages"
)

// TableConfig holds configuration for table rendering
type TableConfig struct {
	Width       int
	ShowHeaders bool
	Interactive bool
	SelectedRow int
	MaxRows     int

	// Sorting configuration
	SortEnabled   bool
	SortColumn    string
	SortDirection string // "asc" or "desc"

}

// SessionData represents session information for table display
type SessionData struct {
	ID          string
	Name        string
	Status      string
	Backend     string
	Created     time.Time
	LastActive  time.Time
	Messages    int
	ProjectPath string
}

// Table provides a unified table component wrapping evertras/bubble-table
// This is the primary table interface for both CLI and TUI components
type SessionTable struct {
	width int

	config TableConfig
	data   []SessionData

	// Evertras table model - this is the core wrapped component
	evertrasModel table.Model

	// Internal state
	initialized bool
	keyMap      TableKeyMap
}

// TableKeyMap defines key bindings for table interaction
type TableKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Enter    key.Binding
	Space    key.Binding
	Filter   key.Binding
	Sort     key.Binding
}

// DefaultTableKeyMap returns the default key bindings for table interaction
func DefaultTableKeyMap() TableKeyMap {
	return TableKeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		PageUp:   key.NewBinding(key.WithKeys("pgup", "b"), key.WithHelp("pgup/b", "page up")),
		PageDown: key.NewBinding(key.WithKeys("pgdown", "f"), key.WithHelp("pgdn/f", "page down")),
		Home:     key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("home/g", "top")),
		End:      key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("end/G", "bottom")),
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Space:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
		Filter:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Sort:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "sort")),
	}
}

// NewTable creates a new table instance with evertras wrapper
func NewSessionTable(config TableConfig) *SessionTable {

	t := &SessionTable{
		config: config,
		keyMap: DefaultTableKeyMap(),
	}

	// Initialize the evertras model
	t.initializeEvertrasModel()

	return t
}

// initializeEvertrasModel sets up the underlying evertras table model
func (t *SessionTable) initializeEvertrasModel() {
	// Create base model with default configuration
	t.evertrasModel = table.New([]table.Column{}).WithBaseStyle(
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.ClaudePrimary).
			Align(lipgloss.Left),
	)

	maxWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		maxWidth = 120
	}

	// Apply initial configuration
	if t.config.Width > 0 {
		t.width = t.config.Width
		t.evertrasModel = t.evertrasModel.WithTargetWidth(t.width)
	} else {
		// Calculate terminal width
		t.width = maxWidth
		t.evertrasModel = t.evertrasModel.WithTargetWidth(maxWidth)
	}

	if t.config.MaxRows > 0 {
		t.evertrasModel = t.evertrasModel.WithMinimumHeight(t.config.MaxRows)
	}

	// Configure for interactive use if enabled
	if t.config.Interactive {
		t.evertrasModel = t.evertrasModel.Focused(true)
	}

	t.initialized = true
}

// refreshEvertrasModel updates the evertras model with current data and config
func (t *SessionTable) refreshEvertrasModel() {
	// Update columns
	columns := GetEvertrasTableColumns()
	if len(columns) > 0 {
		t.evertrasModel = t.evertrasModel.WithColumns(columns)
	}

	// Update rows with processed data
	rows := GetEvertrasSessionRows(t.data, t.width)
	if len(rows) > 0 {
		t.evertrasModel = t.evertrasModel.WithRows(rows)
	}

	// Apply current configuration
	t.applyConfigToEvertrasModel()
}

// applyConfigToEvertrasModel applies current config to the evertras model
func (t *SessionTable) applyConfigToEvertrasModel() {
	if t.config.Width > 0 {
		t.evertrasModel = t.evertrasModel.WithTargetWidth(t.config.Width)
	}

	if t.config.MaxRows > 0 {
		t.evertrasModel = t.evertrasModel.WithMinimumHeight(t.config.MaxRows)
	}

	// Apply interactive state
	if t.config.Interactive {
		t.evertrasModel = t.evertrasModel.Focused(true)
	} else {
		t.evertrasModel = t.evertrasModel.Focused(false)
	}
}

// SetSessionData converts session data to table format
func (t *SessionTable) SetSessionData(sessions []SessionData) {
	t.data = sessions
}

// RenderCLI renders the table for CLI output (static)
func (t *SessionTable) RenderCLI() string {
	if len(t.data) == 0 {
		return styles.Dim("No data to display.")
	}

	t.refreshEvertrasModel()

	// Set the table to non-focused state for static CLI rendering
	t.evertrasModel = t.evertrasModel.Focused(false)

	// Configure table for CLI with appropriate styling
	t.evertrasModel = styles.ConfigureEvertrasTable(t.evertrasModel)

	// Apply any configured limits
	if t.config.MaxRows > 0 {
		t.evertrasModel = t.evertrasModel.WithMaxTotalWidth(t.config.Width)
	}

	// Render the table using evertras View() method
	return t.evertrasModel.View()
}

// Sorting Methods

// SetSort sets the sort column and direction
func (t *SessionTable) SetSort(column, direction string) error {
	if !t.validateSortColumn(column) {
		return fmt.Errorf("invalid sort column: %s", column)
	}
	if direction != "asc" && direction != "desc" {
		return fmt.Errorf("invalid sort direction: %s (must be 'asc' or 'desc')", direction)
	}

	t.config.SortColumn = column
	t.config.SortDirection = direction

	// Apply to evertras model if initialized
	if t.initialized {
		t.SortByColumn(column, direction == "asc")
	}

	return nil
}

// SortByColumn sorts the table by the specified column
func (t *SessionTable) SortByColumn(columnKey string, ascending bool) {
	if t.initialized {
		if ascending {
			t.evertrasModel = t.evertrasModel.SortByAsc(columnKey)
		} else {
			t.evertrasModel = t.evertrasModel.SortByDesc(columnKey)
		}
	}

	// Update internal config
	t.config.SortColumn = columnKey
	if ascending {
		t.config.SortDirection = "asc"
	} else {
		t.config.SortDirection = "desc"
	}
}

// validateSortColumn validates if the given column is valid for sorting
func (t *SessionTable) validateSortColumn(column string) bool {
	validColumns := []string{"id", "name", "status", "backend", "created", "last_active", "project", "messages"}
	return slices.Contains(validColumns, column)
}

// GetEvertrasTableColumns returns predefined column definitions for session table
func GetEvertrasTableColumns() []table.Column {
	columnStyles := styles.GetEvertrasColumnStyles()

	return []table.Column{
		table.NewFlexColumn(columnKeyID, "ID", 2).WithStyle(columnStyles.ID),
		table.NewFlexColumn(columnKeyName, "Name", 2).WithStyle(columnStyles.Name),
		table.NewFlexColumn(columnKeyStatus, "Status", 1).WithStyle(columnStyles.Status),
		table.NewFlexColumn(columnKeyBackend, "Backend", 1).WithStyle(columnStyles.Backend),
		table.NewFlexColumn(columnKeyCreated, "Created", 1).WithStyle(columnStyles.Timestamp),
		table.NewFlexColumn(columnKeyLastActive, "Last Active", 1).WithStyle(columnStyles.Timestamp),
		table.NewFlexColumn(columnKeyProject, "Project", 3).WithStyle(columnStyles.Project),
		table.NewFlexColumn(columnKeyMessages, "Messages", 1).WithStyle(columnStyles.Messages),
	}
}

// ToEvertrasSessionRows converts session data directly to table.Row format
func GetEvertrasSessionRows(sessions []SessionData, width int) []table.Row {
	if len(sessions) == 0 {
		return nil
	}

	rows := make([]table.Row, len(sessions))
	for i, session := range sessions {

		id := session.ID
		name := session.Name

		// Get the width of the project path column
		projectPathWidth := min(width-10, 50)
		projectPath := styles.FormatProjectPath(session.ProjectPath, projectPathWidth)

		timeAgo := styles.FormatTimeAgo(session.LastActive)
		messages := fmt.Sprintf("%d", session.Messages)
		created := styles.FormatTime(session.Created)

		statusStyle := styles.GetContextualColor(styles.ContextStatus, session.Status)
		status := lipgloss.NewStyle().Foreground(statusStyle).Render(styles.FormatStatus(session.Status))

		backendStyle := styles.GetContextualColor(styles.ContextBackend, session.Backend)
		backend := lipgloss.NewStyle().Foreground(backendStyle).Render(session.Backend)

		rows[i] = table.NewRow(table.RowData{
			columnKeyID:         id,
			columnKeyName:       name,
			columnKeyStatus:     status,
			columnKeyBackend:    backend,
			columnKeyCreated:    created,
			columnKeyLastActive: timeAgo,
			columnKeyMessages:   messages,
			columnKeyProject:    projectPath,
		})
	}

	return rows
}
