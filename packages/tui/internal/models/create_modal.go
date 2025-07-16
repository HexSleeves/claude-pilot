package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/styles"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message types
type SessionCreatedMsg struct {
	Session *api.Session
	Error   error
}

// ErrorType represents different types of validation errors
type ErrorType int

const (
	ErrorTypeNone ErrorType = iota
	ErrorTypeValidation
	ErrorTypeWarning
	ErrorTypeAPI
)

// CreateModalModel represents the session creation modal
type CreateModalModel struct {
	client *api.Client
	width  int
	height int

	// Form state - using Bubbles textinput components
	inputs  []textinput.Model
	labels  []string
	focused int

	// Modal state
	completed bool
	err       error
	errType   ErrorType
	keys      KeyMap

	// State management
	isCreating bool // Track if session creation is in progress
}

// Form field indices
// Form field indices
const (
	inputName        = iota // Index for session name input field
	inputDescription        // Index for session description input field
	inputProjectPath        // Index for project path input field
)

// NewCreateModalModel creates a new create modal model
func NewCreateModalModel(client *api.Client) *CreateModalModel {
	// Create Bubbles textinput models
	inputs := make([]textinput.Model, 3)
	labels := []string{"Session Name *", "Description", "Project Path"}

	// Session name input (required)
	inputs[inputName] = textinput.New()
	inputs[inputName].Placeholder = "my-session"
	inputs[inputName].Focus()
	inputs[inputName].Width = 60
	inputs[inputName] = styles.ConfigureBubblesTextInput(inputs[inputName])

	// Description input (optional)
	inputs[inputDescription] = textinput.New()
	inputs[inputDescription].Placeholder = "Optional session description"
	inputs[inputDescription].Width = 60
	inputs[inputDescription] = styles.ConfigureBubblesTextInput(inputs[inputDescription])

	// Project path input (optional)
	inputs[inputProjectPath] = textinput.New()
	inputs[inputProjectPath].Placeholder = "/path/to/project (optional)"
	inputs[inputProjectPath].Width = 60
	inputs[inputProjectPath] = styles.ConfigureBubblesTextInput(inputs[inputProjectPath])

	return &CreateModalModel{
		client:  client,
		inputs:  inputs,
		labels:  labels,
		focused: 0,
		keys:    DefaultKeyMap(),
	}
}

// Init implements tea.Model
func (m *CreateModalModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m *CreateModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update input widths - dashboard handles modal sizing
		for i := range m.inputs {
			m.inputs[i].Width = 60 // Larger width for better layout in wider modal
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			m.Reset()
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			return m.handleEnter()

		case key.Matches(msg, m.keys.Tab):
			m.cycleFocus(false)
			return m, nil

		case key.Matches(msg, m.keys.ShiftTab):
			m.cycleFocus(true)
			return m, nil

		case key.Matches(msg, m.keys.Up):
			m.moveFocus(-1)
			return m, nil

		case key.Matches(msg, m.keys.Down):
			m.moveFocus(1)
			return m, nil

		default:
			// Handle input for the focused field
			if m.focused >= 0 && m.focused < len(m.inputs) {
				m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
				cmds = append(cmds, cmd)
				
				// Clear errors when user starts typing (for better UX)
				if m.err != nil && m.errType == ErrorTypeValidation {
					m.clearError()
				}
			}
		}

	case SessionCreatedMsg:
		// Reset creation state
		m.isCreating = false
		
		if msg.Error != nil {
			// Handle API error
			m.setError(msg.Error, ErrorTypeAPI)
		} else {
			// Success - mark as completed
			m.completed = true
			m.clearError()
		}
		return m, nil
	}

	// Update all inputs for cursor blinking (skip focused input to avoid double update)
	for i := range m.inputs {
		if i != m.focused {
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *CreateModalModel) View() string {
	return m.renderModal()
}

// renderModal renders the modal content (positioning handled by dashboard)
func (m *CreateModalModel) renderModal() string {
	// Build modal content
	var content strings.Builder

	// Title
	title := styles.TitleStyle.Render("Create New Session")
	content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(70).Render(title))
	content.WriteString("\n\n")

	// Form fields
	for i := range m.inputs {
		content.WriteString(m.renderInput(i))
		content.WriteString("\n")
	}

	// Error display with type-based styling
	if m.err != nil {
		content.WriteString("\n")
		content.WriteString(m.renderError())
		content.WriteString("\n")
	}

	// Loading indicator during creation
	if m.isCreating {
		content.WriteString("\n")
		content.WriteString(styles.InfoStyle.Render("Creating session..."))
		content.WriteString("\n")
	}

	// Instructions
	content.WriteString("\n")
	content.WriteString(m.renderInstructions())

	// Return just the content - dashboard will handle positioning and styling
	return content.String()
}

// renderInput renders a single input field using Bubbles textinput
func (m *CreateModalModel) renderInput(index int) string {
	var field strings.Builder

	// Enhanced label with required field indicator and focus state
	labelText := m.labels[index]
	if index == inputName && strings.TrimSpace(m.inputs[inputName].Value()) == "" {
		labelText += " *" // Required field indicator
	}

	labelStyle := styles.LabelStyle
	if index == m.focused {
		labelStyle = lipgloss.NewStyle().
			Foreground(styles.ClaudePrimary).
			Bold(true)
	}

	// Add visual focus indicator
	if index == m.focused {
		labelText = "▶ " + labelText
	} else {
		labelText = "  " + labelText
	}

	// Center the label
	centeredLabel := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(70).
		Render(labelStyle.Render(labelText))

	field.WriteString(centeredLabel)
	field.WriteString("\n")

	// Enhanced input field styling with better visual states
	inputView := m.inputs[index].View()
	
	var borderStyle lipgloss.Style
	if index == m.focused {
		// Focused state with enhanced visual feedback
		borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ClaudePrimary).
			Background(styles.BackgroundSecondary).
			Padding(0, 1)
	} else if m.inputs[index].Value() != "" {
		// Filled state
		borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.SuccessColor).
			Padding(0, 1)
	} else {
		// Empty/default state
		borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.TextMuted).
			Padding(0, 1)
	}

	styledInput := borderStyle.Render(inputView)

	// Center the input field
	centeredInput := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(70).
		Render(styledInput)

	field.WriteString(centeredInput)

	// Add field-specific help text
	if index == m.focused {
		helpText := m.getFieldHelpText(index)
		if helpText != "" {
			field.WriteString("\n")
			helpStyle := lipgloss.NewStyle().
				Foreground(styles.TextDim).
				Align(lipgloss.Center).
				Width(70)
			field.WriteString(helpStyle.Render(helpText))
		}
	}

	return field.String()
}

// getFieldHelpText returns contextual help text for each field
func (m *CreateModalModel) getFieldHelpText(index int) string {
	switch index {
	case inputName:
		return "Alphanumeric characters, hyphens, and underscores only"
	case inputDescription:
		return "Optional description for this session"
	case inputProjectPath:
		return "Optional working directory path"
	default:
		return ""
	}
}

// renderError renders the error message with appropriate styling based on type
func (m *CreateModalModel) renderError() string {
	if m.err == nil {
		return ""
	}

	var errorText string
	switch m.errType {
	case ErrorTypeValidation:
		errorText = styles.ErrorStyle.Render(fmt.Sprintf("✗ %v", m.err))
	case ErrorTypeWarning:
		errorText = styles.WarningStyle.Render(fmt.Sprintf("⚠ %v", m.err))
	case ErrorTypeAPI:
		errorText = styles.ErrorStyle.Render(fmt.Sprintf("✗ API Error: %v", m.err))
	default:
		errorText = styles.ErrorStyle.Render(fmt.Sprintf("✗ %v", m.err))
	}

	// Center the error message
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(70).
		Render(errorText)
}

// renderInstructions renders the help instructions
func (m *CreateModalModel) renderInstructions() string {
	var instructions []string
	
	if m.isCreating {
		instructions = []string{
			"Creating session, please wait...",
		}
	} else {
		instructions = []string{
			"Tab/Shift+Tab: Navigate fields",
			"↑↓: Move between fields",
			"Enter: Create session",
			"Esc: Cancel",
		}
	}

	// Style instructions with divider
	instructionText := strings.Join(instructions, " • ")
	
	// Different styling based on state
	style := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(styles.TextMuted).
		Padding(1, 0, 0, 0)
	
	if m.isCreating {
		style = style.Foreground(styles.InfoColor)
	}
	
	return style.Render(instructionText)
}

// handleEnter handles the enter key press
func (m *CreateModalModel) handleEnter() (tea.Model, tea.Cmd) {
	// Prevent multiple creation attempts
	if m.isCreating {
		return m, nil
	}

	// Validate required fields
	if err := m.validateForm(); err != nil {
		m.setError(err, ErrorTypeValidation)
		return m, nil
	}

	// Clear any previous errors and start creation
	m.clearError()
	m.isCreating = true

	// Create session
	return m, m.createSession()
}

// setError sets an error with the appropriate type for styling
func (m *CreateModalModel) setError(err error, errType ErrorType) {
	m.err = err
	m.errType = errType
}

// clearError clears the current error state
func (m *CreateModalModel) clearError() {
	m.err = nil
	m.errType = ErrorTypeNone
}

// validateForm validates the form inputs with comprehensive checks
func (m *CreateModalModel) validateForm() error {
	name := strings.TrimSpace(m.inputs[inputName].Value())
	projectPath := strings.TrimSpace(m.inputs[inputProjectPath].Value())

	// Check if session name is provided (required field)
	if name == "" {
		return fmt.Errorf("Session Name is required")
	}

	// Validate session name format (alphanumeric, hyphens, underscores)
	if !isValidSessionName(name) {
		return fmt.Errorf("Session name must contain only letters, numbers, hyphens, and underscores")
	}

	// Check for duplicate session names
	if err := m.checkDuplicateName(name); err != nil {
		return err
	}

	// Validate project path if provided
	if projectPath != "" {
		if err := m.validateProjectPath(projectPath); err != nil {
			return err
		}
	}

	return nil
}

// isValidSessionName checks if session name contains only valid characters
func isValidSessionName(name string) bool {
	// Allow alphanumeric characters, hyphens, and underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return len(name) > 0
}

// checkDuplicateName verifies that session name is unique
func (m *CreateModalModel) checkDuplicateName(name string) error {
	sessions, err := m.client.ListSessions()
	if err != nil {
		// If we can't check for duplicates, it's a warning but don't block creation
		// We'll let the API handle the duplicate check on the backend
		return nil
	}

	for _, session := range sessions {
		if session.Name == name {
			return fmt.Errorf("Session name '%s' already exists", name)
		}
	}
	return nil
}

// validateProjectPath checks if project path exists and is accessible
func (m *CreateModalModel) validateProjectPath(path string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Project path does not exist: %s", path)
		}
		return fmt.Errorf("Cannot access project path: %s", path)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("Project path must be a directory: %s", path)
	}

	return nil
}

// createSession creates a new session
func (m *CreateModalModel) createSession() tea.Cmd {
	return func() tea.Msg {
		name := strings.TrimSpace(m.inputs[inputName].Value())
		description := strings.TrimSpace(m.inputs[inputDescription].Value())
		projectPath := strings.TrimSpace(m.inputs[inputProjectPath].Value())

		session, err := m.client.CreateSession(api.CreateSessionRequest{
			Name:        name,
			Description: description,
			ProjectPath: projectPath,
		})

		if err != nil {
			return SessionCreatedMsg{Session: nil, Error: err}
		}

		return SessionCreatedMsg{Session: session, Error: nil}
	}
}

// Focus navigation

func (m *CreateModalModel) cycleFocus(reverse bool) {
	if reverse {
		m.moveFocus(-1)
	} else {
		m.moveFocus(1)
	}
}

func (m *CreateModalModel) moveFocus(delta int) {
	// Blur current input
	if m.focused >= 0 && m.focused < len(m.inputs) {
		m.inputs[m.focused].Blur()
	}

	// Update focus index
	m.focused += delta
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	} else if m.focused >= len(m.inputs) {
		m.focused = 0
	}

	// Focus new input
	if m.focused >= 0 && m.focused < len(m.inputs) {
		m.inputs[m.focused].Focus()
	}
}

// Text input methods are now handled by Bubbles textinput components

// State management

func (m *CreateModalModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Adjust input widths based on modal size for better responsiveness
	inputWidth := min(60, width-10) // Ensure inputs fit within modal
	for i := range m.inputs {
		m.inputs[i].Width = inputWidth
	}
}

func (m *CreateModalModel) IsCompleted() bool {
	return m.completed
}

func (m *CreateModalModel) Reset() {
	m.completed = false
	m.err = nil
	m.errType = ErrorTypeNone
	m.isCreating = false
	m.focused = 0

	// Clear all input values and reset focus
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}

	// Focus the first input
	if len(m.inputs) > 0 {
		m.inputs[0].Focus()
	}
}
