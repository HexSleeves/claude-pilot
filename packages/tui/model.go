package tui

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/interfaces"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"

	"claude-pilot/shared/styles"
)

// ViewState represents the current view state of the TUI
type ViewState int

const (
	TableView ViewState = iota
	CreatePrompt
	Loading
	Error
	Help
	KillConfirmation
	FilterView
	ExportView
)

// Model represents the main TUI model implementing bubbletea.Model interface
type Model struct {
	// Core dependencies
	client *api.Client

	// State management
	currentView   ViewState
	errorMessage  string
	statusMessage string
	lastRefresh   time.Time

	// Session data
	sessions  []*interfaces.Session
	isLoading bool

	// UI components
	table  table.Model
	keymap KeyMap

	// Create session form inputs
	nameInput        textinput.Model
	descriptionInput textinput.Model
	pathInput        textinput.Model
	activeInput      int // 0=name, 1=description, 2=path

	// Kill confirmation state
	sessionToKill *interfaces.Session

	// Filter input
	filterInput textinput.Model

	// Export state
	exportFormat      string // "csv" or "json"
	exportFilename    textinput.Model
	exportActiveInput int // 0=format, 1=filename

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int

	// Help visibility
	showHelp bool

	// Advanced table interaction state
	tablePageSize     int
	tableCurrentPage  int
	tableSelectedRows []int
	showTableHelp     bool

	// Sorting state
	sortColumn    string
	sortDirection string // "asc" or "desc"

	// Filter state
	filterActive bool
	filterQuery  string
	// filteredSessions []*interfaces.Session
}

const (
	// Layout constants
	fixedVerticalMargin = 4 // Fixed margin for description & instructions

	// Input field indices for form navigation
	nameInputIndex        = 0
	descriptionInputIndex = 1
	pathInputIndex        = 2
	maxInputIndex         = 2
)

// NewModel creates a new TUI model with the provided API client.
// It initializes all UI components, sets up default styling, and returns
// a ready-to-use Model for the bubbletea application.
func NewModel(client *api.Client) Model {
	if client == nil {
		panic("NewModel: client cannot be nil")
	}
	// Create text inputs for session creation
	nameInput := textinput.New()
	nameInput.Placeholder = "Session name"
	nameInput.Focus()
	nameInput.CharLimit = 50
	nameInput.Width = 30

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Session description (optional)"
	descriptionInput.CharLimit = 100
	descriptionInput.Width = 50

	pathInput := textinput.New()
	pathInput.Placeholder = "Project path (optional)"
	pathInput.CharLimit = 200
	pathInput.Width = 60

	filterInput := textinput.New()
	filterInput.Placeholder = "Filter sessions (name, status, description...)"
	filterInput.CharLimit = 100
	filterInput.Width = 50

	exportFilename := textinput.New()
	exportFilename.Placeholder = "sessions-export"
	exportFilename.CharLimit = 200
	exportFilename.Width = 40

	return Model{
		client:           client,
		currentView:      Loading,
		keymap:           DefaultKeyMap(),
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		pathInput:        pathInput,
		filterInput:      filterInput,
		activeInput:      nameInputIndex,
		sessions:         []*interfaces.Session{},
		isLoading:        false,
		showHelp:         false,

		// Initialize advanced table interaction state
		tablePageSize:     10,
		tableCurrentPage:  1,
		tableSelectedRows: []int{},
		showTableHelp:     false,

		// Initialize filter state
		filterActive: false,
		filterQuery:  "",

		// Initialize export state
		exportFormat:      "csv",
		exportFilename:    exportFilename,
		exportActiveInput: 0,

		// Initialize sort state
		sortColumn:    "",
		sortDirection: "asc",

		// Initialize evertras table with shared component columns and Claude theme styling
		table: styles.ConfigureEvertrasTable(
			table.New(components.GetEvertrasTableColumns()),
		),
	}
}

// Init initializes the model and returns a command to load initial session data
func (m Model) Init() tea.Cmd {
	return loadSessionsCmd(m.client)
}

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()

	case tea.KeyMsg:
		// Handle global keys first
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Help):
			m.showHelp = !m.showHelp
			if m.currentView != Help && m.showHelp {
				m.currentView = Help
			} else if m.currentView == Help && !m.showHelp {
				m.currentView = TableView
			}

		case key.Matches(msg, m.keymap.Back):
			if m.currentView == CreatePrompt || m.currentView == Help {
				m.currentView = TableView
				m.resetCreateForm()
			} else if m.currentView == KillConfirmation {
				m.currentView = TableView
				m.sessionToKill = nil
			} else if m.currentView == FilterView {
				m.currentView = TableView
				m.filterInput.Blur()
				// Keep filter active - don't clear it on escape
			} else if m.currentView == ExportView {
				m.currentView = TableView
				m.exportFilename.Blur()
			}
		}

		// Handle view-specific keys
		switch m.currentView {
		case TableView:
			cmd = m.handleTableViewKeys(msg)
		case CreatePrompt:
			cmd = m.handleCreatePromptKeys(msg)
		case KillConfirmation:
			cmd = m.handleKillConfirmationKeys(msg)
		case FilterView:
			cmd = m.handleFilterViewKeys(msg)
		case ExportView:
			cmd = m.handleExportViewKeys(msg)
		case Error:
			if key.Matches(msg, m.keymap.Refresh) {
				m.currentView = TableView
				m.errorMessage = ""
				cmd = loadSessionsCmd(m.client)
			}
		}

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case sessionsLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.currentView = Error
			m.errorMessage = msg.err.Error()
		} else {
			m.sessions = msg.sessions
			m.updateTableData()
			m.lastRefresh = time.Now()
			if m.currentView == Loading {
				m.currentView = TableView
			}
		}

	case sessionCreatedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.currentView = Error
			m.errorMessage = msg.err.Error()
		} else {
			m.currentView = TableView
			m.resetCreateForm()
			m.statusMessage = "Session created successfully"
			cmds = append(cmds, loadSessionsCmd(m.client))
		}

	case sessionKilledMsg:
		m.isLoading = false
		if msg.err != nil {
			m.currentView = Error
			m.errorMessage = msg.err.Error()
			m.sessionToKill = nil
		} else {
			m.currentView = TableView
			m.sessionToKill = nil
			m.statusMessage = "Session killed successfully"
			cmds = append(cmds, loadSessionsCmd(m.client))
		}

	case errorMsg:
		m.isLoading = false
		m.currentView = Error
		m.errorMessage = msg.error.Error()

	case statusMsg:
		m.statusMessage = msg.message

	case viewStateMsg:
		m.currentView = msg.state

	case TableSortedMsg:
		// Update status message to show current sort
		direction := "ascending"
		if msg.Direction == "desc" {
			direction = "descending"
		}
		m.statusMessage = fmt.Sprintf("Sorted by %s (%s)", msg.Column, direction)
	}

	// Update table model
	m.table, cmd = m.table.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Update input models when in create view
	if m.currentView == CreatePrompt {
		switch m.activeInput {
		case nameInputIndex:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case descriptionInputIndex:
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		case pathInputIndex:
			m.pathInput, cmd = m.pathInput.Update(msg)
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update filter input when in filter view
	if m.currentView == FilterView {
		m.filterInput, cmd = m.filterInput.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		// Apply filter in real-time as user types
		m.filterQuery = m.filterInput.Value()
		m.table = m.table.WithFilterInput(m.filterInput)
	}

	// Update export input when in export view
	if m.currentView == ExportView {
		m.exportFilename, cmd = m.exportFilename.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the current view based on state
func (m Model) View() string {
	switch m.currentView {
	case TableView:
		return renderTableView(m)
	case CreatePrompt:
		return renderCreateView(m)
	case Loading:
		return renderLoadingView(m)
	case Error:
		return renderErrorView(m)
	case Help:
		return renderHelpView(m)
	case KillConfirmation:
		return renderKillConfirmationView(m)
	case FilterView:
		return renderFilterView(m)
	case ExportView:
		return renderExportView(m)
	default:
		return renderTableView(m)
	}
}

// handleTableSelection manages row selection
func (m *Model) handleTableSelection(action string, rowIndex int) tea.Cmd {
	switch action {
	case "select_all":
		m.tableSelectedRows = make([]int, len(m.sessions))
		for i := range m.sessions {
			m.tableSelectedRows[i] = i
		}
	case "deselect_all":
		m.tableSelectedRows = []int{}
	case "toggle":
		if rowIndex >= 0 && rowIndex < len(m.sessions) {
			// Check if row is already selected
			found := false
			for i, selected := range m.tableSelectedRows {
				if selected == rowIndex {
					// Remove from selection
					m.tableSelectedRows = append(m.tableSelectedRows[:i], m.tableSelectedRows[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				// Add to selection
				m.tableSelectedRows = append(m.tableSelectedRows, rowIndex)
			}
		}
	case "invert":
		newSelection := []int{}
		for i := range m.sessions {
			selected := false
			for _, selectedRow := range m.tableSelectedRows {
				if selectedRow == i {
					selected = true
					break
				}
			}
			if !selected {
				newSelection = append(newSelection, i)
			}
		}
		m.tableSelectedRows = newSelection
	}

	m.refreshTableWithCurrentState()
	return nil
}

// refreshTableWithCurrentState updates table with current sort, and pagination settings
func (m *Model) refreshTableWithCurrentState() {
	m.updateTableData()
	m.syncTableState()
}

// syncTableState keeps model state in sync with table component state
func (m *Model) syncTableState() {
	// Apply current state to the table component
	// This method ensures the evertras table reflects our model state
	m.table = m.table.WithTargetWidth(m.calculateWidth()).WithMinimumHeight(m.calculateHeight())
}

// paginateSessionData returns a specific page of session data
func (m *Model) paginateSessionData(sessions []components.SessionData, page, pageSize int) []components.SessionData {
	if pageSize <= 0 || page <= 0 {
		return sessions
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(sessions) {
		return []components.SessionData{} // Page beyond available data
	}

	if endIndex > len(sessions) {
		endIndex = len(sessions)
	}

	return sessions[startIndex:endIndex]
}

func (m *Model) recalculateTable() {
	width := m.calculateWidth()
	// height := m.calculateHeight()

	// Update evertras table with responsive columns and dimensions
	m.table = m.table.
		WithColumns(components.GetEvertrasTableColumns()).
		WithTargetWidth(width)
	// WithMinimumHeight(height)
}

func (m Model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m Model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - fixedVerticalMargin
}

// handleTableViewKeys handles keyboard input in table view
func (m *Model) handleTableViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keymap.Create):
		m.currentView = CreatePrompt
		m.nameInput.Focus()
		m.activeInput = 0

	case key.Matches(msg, m.keymap.Kill):
		if len(m.sessions) > 0 {
			highlightedRow := m.table.GetHighlightedRowIndex()
			if highlightedRow >= 0 && highlightedRow < len(m.sessions) {
				session := m.sessions[highlightedRow]
				if session != nil && session.ID != "" {
					m.sessionToKill = session
					m.currentView = KillConfirmation
				}
			}
		}

	case key.Matches(msg, m.keymap.Attach):
		if len(m.sessions) > 0 {
			highlightedRow := m.table.GetHighlightedRowIndex()
			if highlightedRow >= 0 && highlightedRow < len(m.sessions) {
				session := m.sessions[highlightedRow]
				if session != nil && session.ID != "" {
					return attachSessionCmd(m.client, session.ID)
				}
			}
		}

	case key.Matches(msg, m.keymap.Refresh):
		if m.client != nil {
			m.isLoading = true
			m.currentView = Loading
			return loadSessionsCmd(m.client)
		}

	// Table selection keys
	case key.Matches(msg, m.keymap.SelectAll):
		return m.handleTableSelection("select_all", -1)
	case key.Matches(msg, m.keymap.DeselectAll):
		return m.handleTableSelection("deselect_all", -1)
	case key.Matches(msg, m.keymap.ToggleRowSelection):
		highlightedRow := m.table.GetHighlightedRowIndex()
		return m.handleTableSelection("toggle", highlightedRow)
	case key.Matches(msg, m.keymap.InvertSelection):
		return m.handleTableSelection("invert", -1)

	// Table view options
	case key.Matches(msg, m.keymap.ToggleRowNumbers):
		// Toggle row numbers display (implementation depends on table component)
		m.refreshTableWithCurrentState()
	case key.Matches(msg, m.keymap.ToggleCompactView):
		// Toggle compact view (implementation depends on table component)
		m.refreshTableWithCurrentState()
	case key.Matches(msg, m.keymap.RefreshTable):
		m.refreshTableWithCurrentState()
		if m.client != nil {
			return loadSessionsCmd(m.client)
		}

	// Toggle table help
	case key.Matches(msg, m.keymap.Help):
		m.showTableHelp = !m.showTableHelp

	// Sorting
	case key.Matches(msg, m.keymap.SortByName):
		return m.toggleSort("name")
	case key.Matches(msg, m.keymap.SortByStatus):
		return m.toggleSort("status")
	case key.Matches(msg, m.keymap.SortByCreated):
		return m.toggleSort("created")
	case key.Matches(msg, m.keymap.SortByLastActive):
		return m.toggleSort("last_active")

	// Filter
	case key.Matches(msg, m.keymap.Filter):
		m.currentView = FilterView
		m.filterInput.Focus()
		m.filterInput.SetValue(m.filterQuery) // Restore previous filter

	// Export
	case key.Matches(msg, m.keymap.Export):
		m.currentView = ExportView
		m.exportActiveInput = 0 // Start with format selection
		m.exportFilename.Focus()
	}

	return nil
}

// handleCreatePromptKeys handles keyboard input in create prompt view
func (m *Model) handleCreatePromptKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keymap.Submit):
		// Optimize by trimming once and reusing
		name := strings.TrimSpace(m.nameInput.Value())
		if name != "" && m.client != nil {
			m.isLoading = true
			m.currentView = Loading

			// Pre-trim all values to avoid repeated operations
			description := strings.TrimSpace(m.descriptionInput.Value())
			projectPath := strings.TrimSpace(m.pathInput.Value())

			return createSessionCmd(m.client, name, description, projectPath)
		} else if name == "" {
			// Show error message for empty name
			m.statusMessage = "Session name is required"
		}

	case key.Matches(msg, m.keymap.NextInput):
		m.switchToNextInput()

	case key.Matches(msg, m.keymap.PrevInput):
		m.switchToPrevInput()
	}

	return nil
}

// handleKillConfirmationKeys handles keyboard input in kill confirmation view
func (m *Model) handleKillConfirmationKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keymap.Yes):
		if m.sessionToKill != nil && m.sessionToKill.ID != "" {
			m.isLoading = true
			m.currentView = Loading
			return killSessionCmd(m.client, m.sessionToKill.ID)
		}
	case key.Matches(msg, m.keymap.No):
		m.currentView = TableView
		m.sessionToKill = nil
	}
	return nil
}

// handleFilterViewKeys handles keyboard input in filter view
func (m *Model) handleFilterViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keymap.Submit):
		// Apply filter and return to table view
		m.currentView = TableView
		m.filterInput.Blur()
		if m.filterQuery == "" {
			m.filterActive = false
			m.statusMessage = "Filter cleared"
		} else {
			m.statusMessage = fmt.Sprintf("Filtering: %s", m.filterQuery)
		}
		m.updateTableData()
		return nil
	}
	return nil
}

// handleExportViewKeys handles keyboard input in export view
func (m *Model) handleExportViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keymap.Submit):
		// Validate and execute export
		filename := strings.TrimSpace(m.exportFilename.Value())
		if filename == "" {
			filename = "sessions-export"
		}

		// Add file extension if not present
		if !strings.HasSuffix(filename, "."+m.exportFormat) {
			filename += "." + m.exportFormat
		}

		// Convert to SessionData format for export
		sessionData := make([]components.SessionData, 0, len(m.sessions))
		for _, session := range m.sessions {
			if session == nil {
				continue
			}
			sessionData = append(sessionData, components.SessionData{
				ID:          session.ID,
				Name:        session.Name,
				Status:      string(session.Status),
				Backend:     session.Backend,
				Created:     session.CreatedAt,
				LastActive:  session.LastActive,
				Messages:    len(session.Messages),
				ProjectPath: session.ProjectPath,
			})
		}

		// Return to table view and execute export
		m.currentView = TableView
		m.exportFilename.Blur()
		return ExportTableDataCmd(m.exportFormat, filename, sessionData)

	case key.Matches(msg, m.keymap.NextInput):
		// Toggle between format and filename (currently only filename input)
		return nil

	case key.Matches(msg, m.keymap.PrevInput):
		// Toggle between format and filename (currently only filename input)
		return nil
	}
	return nil
}

// updateTableData updates the table with current session data
// using the shared component's conversion methods with enhanced features
func (m *Model) updateTableData() {
	if len(m.sessions) == 0 {
		// Clear table data to free memory
		m.table = m.table.WithRows([]table.Row{})
		return
	}

	// Convert interfaces.Session to components.SessionData for shared component utility
	sessionData := make([]components.SessionData, 0, len(m.sessions))
	for _, session := range m.sessions {
		if session == nil {
			continue // Skip nil sessions
		}

		sessionData = append(sessionData, components.SessionData{
			ID:          session.ID,
			Name:        session.Name,
			Status:      string(session.Status),
			Backend:     session.Backend,
			Created:     session.CreatedAt,
			LastActive:  session.LastActive,
			Messages:    len(session.Messages),
			ProjectPath: session.ProjectPath,
		})
	}

	// Apply pagination if enabled
	if m.tablePageSize > 0 {
		sessionData = m.paginateSessionData(sessionData, m.tableCurrentPage, m.tablePageSize)
	}

	// Use shared component's method to convert to evertras rows
	rows := components.GetEvertrasSessionRows(sessionData, m.calculateWidth())

	// Update evertras table with new rows
	m.table = m.table.WithRows(rows)
}

// resetCreateForm resets the create session form to its initial state
// and clears any sensitive data from memory
func (m *Model) resetCreateForm() {
	// Clear input values to prevent memory retention
	m.nameInput.SetValue("")
	m.descriptionInput.SetValue("")
	m.pathInput.SetValue("")

	// Reset focus state
	m.activeInput = nameInputIndex
	m.nameInput.Focus()
	m.descriptionInput.Blur()
	m.pathInput.Blur()

	// Clear any status messages to prevent memory buildup
	m.statusMessage = ""
}

// switchToNextInput switches focus to the next input field
func (m *Model) switchToNextInput() {
	if m == nil {
		return // Defensive check
	}
	m.blurAllInputs()
	m.activeInput = (m.activeInput + 1) % (maxInputIndex + 1)
	m.focusActiveInput()
}

// switchToPrevInput switches focus to the previous input field
func (m *Model) switchToPrevInput() {
	if m == nil {
		return // Defensive check
	}
	m.blurAllInputs()
	m.activeInput = (m.activeInput - 1 + (maxInputIndex + 1)) % (maxInputIndex + 1)
	m.focusActiveInput()
}

// blurAllInputs removes focus from all input fields
func (m *Model) blurAllInputs() {
	if m == nil {
		return // Defensive check
	}
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.pathInput.Blur()
}

// focusActiveInput sets focus on the currently active input field
func (m *Model) focusActiveInput() {
	if m == nil {
		return // Defensive check
	}
	switch m.activeInput {
	case nameInputIndex:
		m.nameInput.Focus()
	case descriptionInputIndex:
		m.descriptionInput.Focus()
	case pathInputIndex:
		m.pathInput.Focus()
	}
}

// toggleSort toggles sorting by the specified column
func (m *Model) toggleSort(column string) tea.Cmd {
	// If clicking the same column, toggle direction
	if m.sortColumn == column {
		if m.sortDirection == "asc" {
			m.sortDirection = "desc"
		} else {
			m.sortDirection = "asc"
		}
	} else {
		// New column, default to ascending
		m.sortColumn = column
		m.sortDirection = "asc"
	}

	// Apply the sort immediately
	m.applySorting()

	// Return a command to indicate sorting is complete
	return func() tea.Msg {
		return TableSortedMsg{
			Column:    m.sortColumn,
			Direction: m.sortDirection,
		}
	}
}

// applySorting sorts the sessions based on current sort settings
func (m *Model) applySorting() {
	if m.sortColumn == "" || len(m.sessions) == 0 {
		return
	}

	// Sort the sessions slice
	sort.Slice(m.sessions, func(i, j int) bool {
		var result bool

		switch m.sortColumn {
		case "name":
			result = m.sessions[i].Name < m.sessions[j].Name
		case "status":
			result = m.sessions[i].Status < m.sessions[j].Status
		case "created":
			result = m.sessions[i].CreatedAt.Before(m.sessions[j].CreatedAt)
		case "last_active":
			// Handle zero LastActive (which means never active)
			iZero := m.sessions[i].LastActive.IsZero()
			jZero := m.sessions[j].LastActive.IsZero()

			if iZero && jZero {
				result = false
			} else if iZero {
				result = false
			} else if jZero {
				result = true
			} else {
				result = m.sessions[i].LastActive.Before(m.sessions[j].LastActive)
			}
		default:
			return false
		}

		// Reverse for descending order
		if m.sortDirection == "desc" {
			result = !result
		}

		return result
	})

	// Update the table with sorted data
	m.updateTableData()
}

// applyFilter filters sessions based on the current filter query
// func (m *Model) applyFilter() {
// 	if m.filterQuery == "" {
// 		m.filterActive = false
// 		m.filteredSessions = m.sessions
// 	} else {
// 		m.filterActive = true
// 		m.filteredSessions = []*interfaces.Session{}
// 		query := strings.ToLower(m.filterQuery)

// 		for _, session := range m.sessions {
// 			// Search in multiple fields
// 			if strings.Contains(strings.ToLower(session.Name), query) ||
// 				strings.Contains(strings.ToLower(session.Description), query) ||
// 				strings.Contains(strings.ToLower(string(session.Status)), query) ||
// 				strings.Contains(strings.ToLower(session.ProjectPath), query) ||
// 				strings.Contains(strings.ToLower(session.ID), query) {
// 				m.filteredSessions = append(m.filteredSessions, session)
// 			}
// 		}
// 	}
// }
