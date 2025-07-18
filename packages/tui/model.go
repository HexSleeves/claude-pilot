package tui

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/interfaces"
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

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int

	// Help visibility
	showHelp bool

	// Advanced table interaction state
	tableFilter        textinput.Model
	filterMode         bool
	tableSortColumn    string
	tableSortDirection string
	tablePageSize      int
	tableCurrentPage   int
	tableSelectedRows  []int
	showTableHelp      bool
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

	// Create table filter input
	tableFilter := textinput.New()
	tableFilter.Placeholder = "Filter sessions..."
	tableFilter.CharLimit = 100
	tableFilter.Width = 40

	return Model{
		client:           client,
		currentView:      TableView,
		keymap:           DefaultKeyMap(),
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		pathInput:        pathInput,
		activeInput:      nameInputIndex,
		sessions:         []*interfaces.Session{},
		isLoading:        false,
		showHelp:         false,

		// Initialize advanced table interaction state
		tableFilter:        tableFilter,
		filterMode:         false,
		tableSortColumn:    "",
		tableSortDirection: "asc",
		tablePageSize:      10,
		tableCurrentPage:   1,
		tableSelectedRows:  []int{},
		showTableHelp:      false,

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

	// Update filter input when in filter mode
	if m.filterMode && m.currentView == TableView {
		m.tableFilter, cmd = m.tableFilter.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle filter input changes
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				// Apply filter and exit filter mode
				m.filterMode = false
				m.tableFilter.Blur()
				cmd = m.handleTableFilter(m.tableFilter.Value())
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			case "esc":
				// Cancel filter mode without applying
				m.filterMode = false
				m.tableFilter.Blur()
			}
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
	default:
		return renderTableView(m)
	}
}

// Helper methods for table state management

// handleTableSort updates sort state and refreshes table
func (m *Model) handleTableSort(column, direction string) tea.Cmd {
	m.tableSortColumn = column
	m.tableSortDirection = direction
	m.refreshTableWithCurrentState()
	return nil
}

// toggleSortDirection toggles sort direction for a column
func (m *Model) toggleSortDirection(column string) string {
	if m.tableSortColumn == column && m.tableSortDirection == "asc" {
		return "desc"
	}
	return "asc"
}

// handleTableFilter updates filter and refreshes table
func (m *Model) handleTableFilter(filterText string) tea.Cmd {
	m.tableFilter.SetValue(filterText)
	m.refreshTableWithCurrentState()
	return nil
}

// handleTablePagination handles page navigation
func (m *Model) handleTablePagination(action string) tea.Cmd {
	totalPages := m.calculateTotalPages()

	switch action {
	case "next":
		if m.tableCurrentPage < totalPages {
			m.tableCurrentPage++
		}
	case "prev":
		if m.tableCurrentPage > 1 {
			m.tableCurrentPage--
		}
	case "first":
		m.tableCurrentPage = 1
	case "last":
		m.tableCurrentPage = totalPages
	}

	m.refreshTableWithCurrentState()
	return nil
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

// refreshTableWithCurrentState updates table with current sort, filter, and pagination settings
func (m *Model) refreshTableWithCurrentState() {
	m.updateTableData()
	m.syncTableState()
}

// updateTableDimensions handles responsive table sizing
func (m *Model) updateTableDimensions() {
	m.recalculateTable()
}

// syncTableState keeps model state in sync with table component state
func (m *Model) syncTableState() {
	// Apply current state to the table component
	// This method ensures the evertras table reflects our model state
	m.table = m.table.WithTargetWidth(m.calculateWidth()).WithMinimumHeight(m.calculateHeight())
}

// calculateTotalPages calculates total pages based on current data and page size
func (m *Model) calculateTotalPages() int {
	if m.tablePageSize <= 0 {
		return 1
	}

	totalRows := len(m.sessions)
	return (totalRows + m.tablePageSize - 1) / m.tablePageSize
}

// sortSessionData sorts session data by the specified column and direction
func (m *Model) sortSessionData(sessions []components.SessionData, column, direction string) []components.SessionData {
	if len(sessions) == 0 {
		return sessions
	}

	sortedSessions := make([]components.SessionData, len(sessions))
	copy(sortedSessions, sessions)

	sort.Slice(sortedSessions, func(i, j int) bool {
		switch column {
		case "id":
			if direction == "desc" {
				return sortedSessions[i].ID > sortedSessions[j].ID
			}
			return sortedSessions[i].ID < sortedSessions[j].ID
		case "name":
			if direction == "desc" {
				return sortedSessions[i].Name > sortedSessions[j].Name
			}
			return sortedSessions[i].Name < sortedSessions[j].Name
		case "status":
			if direction == "desc" {
				return sortedSessions[i].Status > sortedSessions[j].Status
			}
			return sortedSessions[i].Status < sortedSessions[j].Status
		case "backend":
			if direction == "desc" {
				return sortedSessions[i].Backend > sortedSessions[j].Backend
			}
			return sortedSessions[i].Backend < sortedSessions[j].Backend
		case "created":
			if direction == "desc" {
				return sortedSessions[i].Created.After(sortedSessions[j].Created)
			}
			return sortedSessions[i].Created.Before(sortedSessions[j].Created)
		case "last_active":
			if direction == "desc" {
				return sortedSessions[i].LastActive.After(sortedSessions[j].LastActive)
			}
			return sortedSessions[i].LastActive.Before(sortedSessions[j].LastActive)
		case "project":
			if direction == "desc" {
				return sortedSessions[i].ProjectPath > sortedSessions[j].ProjectPath
			}
			return sortedSessions[i].ProjectPath < sortedSessions[j].ProjectPath
		case "messages":
			if direction == "desc" {
				return sortedSessions[i].Messages > sortedSessions[j].Messages
			}
			return sortedSessions[i].Messages < sortedSessions[j].Messages
		default:
			return false
		}
	})

	return sortedSessions
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

// getCurrentPageInfo returns information about the current page
func (m *Model) getCurrentPageInfo() (currentPage, totalPages, startRow, endRow, totalRows int) {
	totalRows = len(m.sessions)

	if m.tablePageSize <= 0 {
		return 1, 1, 1, totalRows, totalRows
	}

	currentPage = m.tableCurrentPage
	totalPages = m.calculateTotalPages()

	startRow = (currentPage-1)*m.tablePageSize + 1
	endRow = currentPage * m.tablePageSize
	if endRow > totalRows {
		endRow = totalRows
	}

	return currentPage, totalPages, startRow, endRow, totalRows
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

	// Table sorting keys
	case key.Matches(msg, m.keymap.SortByName):
		return m.handleTableSort("name", m.toggleSortDirection("name"))
	case key.Matches(msg, m.keymap.SortByStatus):
		return m.handleTableSort("status", m.toggleSortDirection("status"))
	case key.Matches(msg, m.keymap.SortByCreated):
		return m.handleTableSort("created", m.toggleSortDirection("created"))
	case key.Matches(msg, m.keymap.SortByLastActive):
		return m.handleTableSort("last_active", m.toggleSortDirection("last_active"))
	case key.Matches(msg, m.keymap.ToggleSortDirection):
		if m.tableSortColumn != "" {
			direction := "asc"
			if m.tableSortDirection == "asc" {
				direction = "desc"
			}
			return m.handleTableSort(m.tableSortColumn, direction)
		}
	case key.Matches(msg, m.keymap.ClearSort):
		m.tableSortColumn = ""
		m.tableSortDirection = "asc"
		m.refreshTableWithCurrentState()

	// Table pagination keys
	case key.Matches(msg, m.keymap.NextPage):
		return m.handleTablePagination("next")
	case key.Matches(msg, m.keymap.PrevPage):
		return m.handleTablePagination("prev")
	case key.Matches(msg, m.keymap.FirstPage):
		return m.handleTablePagination("first")
	case key.Matches(msg, m.keymap.LastPage):
		return m.handleTablePagination("last")
	case key.Matches(msg, m.keymap.PageSizeIncrease):
		if m.tablePageSize < 50 {
			m.tablePageSize += 5
			m.refreshTableWithCurrentState()
		}
	case key.Matches(msg, m.keymap.PageSizeDecrease):
		if m.tablePageSize > 5 {
			m.tablePageSize -= 5
			m.refreshTableWithCurrentState()
		}

	// Table filtering keys
	case key.Matches(msg, m.keymap.ToggleFilter):
		m.filterMode = !m.filterMode
		if m.filterMode {
			m.tableFilter.Focus()
		} else {
			m.tableFilter.Blur()
		}
	case key.Matches(msg, m.keymap.ClearFilter):
		m.tableFilter.SetValue("")
		m.filterMode = false
		m.tableFilter.Blur()
		return m.handleTableFilter("")
	case key.Matches(msg, m.keymap.FocusFilter):
		m.filterMode = true
		m.tableFilter.Focus()

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
			Backend:     "claude", // Default backend for now
			Created:     session.CreatedAt,
			LastActive:  session.LastActive,
			Messages:    len(session.Messages),
			ProjectPath: session.ProjectPath,
		})
	}

	// Apply filtering if active
	if m.tableFilter.Value() != "" {
		filteredData := []components.SessionData{}
		filterText := strings.ToLower(m.tableFilter.Value())
		for _, session := range sessionData {
			if strings.Contains(strings.ToLower(session.Name), filterText) ||
				strings.Contains(strings.ToLower(session.Status), filterText) ||
				strings.Contains(strings.ToLower(session.ID), filterText) ||
				strings.Contains(strings.ToLower(session.ProjectPath), filterText) {
				filteredData = append(filteredData, session)
			}
		}
		sessionData = filteredData
	}

	// Apply sorting if active
	if m.tableSortColumn != "" {
		sessionData = m.sortSessionData(sessionData, m.tableSortColumn, m.tableSortDirection)
	}

	// Apply pagination if enabled
	if m.tablePageSize > 0 {
		sessionData = m.paginateSessionData(sessionData, m.tableCurrentPage, m.tablePageSize)
	}

	// Use shared component's method to convert to evertras rows
	rows := components.ToEvertrasSessionRows(sessionData)

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
