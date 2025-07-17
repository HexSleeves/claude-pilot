package main

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/interfaces"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// ViewState represents the current view state of the TUI
type ViewState int

const (
	TableView ViewState = iota
	CreatePrompt
	Loading
	Error
	Help
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

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int

	// Help visibility
	showHelp bool
}

const (
	// Table column keys for evertras/bubble-table
	columnKeyID         = "id"
	columnKeyName       = "name"
	columnKeyStatus     = "status"
	columnKeyBackend    = "backend"
	columnKeyCreated    = "created"
	columnKeyLastActive = "last_active"
	columnKeyMessages   = "messages"
	columnKeyProject    = "project"

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
	// Create table model with predefined columns
	// t := table.New(
	// 	table.WithColumns(components.GetBubblesTableColumns()),
	// 	table.WithFocused(true),
	// 	table.WithHeight(15),
	// )

	// // Apply Claude theme styling to the table
	// t.SetStyles(styles.GetBubblesTableStyles())

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

		table: table.New([]table.Column{
			table.NewFlexColumn(columnKeyID, "ID", 1),
			table.NewFlexColumn(columnKeyName, "Name", 1),
			table.NewFlexColumn(columnKeyStatus, "Status", 1),
			table.NewFlexColumn(columnKeyBackend, "Backend", 1),
			table.NewFlexColumn(columnKeyCreated, "Created", 1),
			table.NewFlexColumn(columnKeyLastActive, "Last Active", 1),
			table.NewFlexColumn(columnKeyProject, "Project", 1),
			table.NewFlexColumn(columnKeyMessages, "Messages", 1),
		}).WithRows([]table.Row{}),
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
			}
		}

		// Handle view-specific keys
		switch m.currentView {
		case TableView:
			cmd = m.handleTableViewKeys(msg)
		case CreatePrompt:
			cmd = m.handleCreatePromptKeys(msg)
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
		} else {
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
	default:
		return renderTableView(m)
	}
}

func (m *Model) recalculateTable() {
	m.table = m.table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight())
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
					m.isLoading = true
					m.currentView = Loading
					return killSessionCmd(m.client, session.ID)
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

// updateTableData updates the table with current session data
// and manages memory efficiently to prevent leaks
func (m *Model) updateTableData() {
	if len(m.sessions) == 0 {
		// Clear table data to free memory
		m.table = m.table.WithRows([]table.Row{})
		return
	}

	// Pre-allocate slice with exact capacity to avoid reallocations
	rows := make([]table.Row, 0, len(m.sessions))

	// Convert sessions to table rows with optimized processing
	for _, session := range m.sessions {
		if session == nil {
			continue // Skip nil sessions
		}

		// Create row data for evertras/bubble-table
		// Use defensive copying to prevent reference retention
		rowData := table.RowData{
			columnKeyID:         session.ID,
			columnKeyName:       session.Name,
			columnKeyStatus:     string(session.Status),
			columnKeyBackend:    "claude", // Default backend for now
			columnKeyCreated:    session.CreatedAt.Format("2006-01-02 15:04"),
			columnKeyLastActive: session.LastActive.Format("2006-01-02 15:04"),
			columnKeyMessages:   fmt.Sprintf("%d", len(session.Messages)),
			columnKeyProject:    session.ProjectPath,
		}

		rows = append(rows, table.NewRow(rowData))
	}

	if len(rows) > 0 {

		// Update table with new rows, replacing old data completely
		m.table = m.table.WithRows(rows)
	}
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
