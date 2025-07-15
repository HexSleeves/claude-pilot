package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/styles"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CreateModalModel represents the session creation modal
type CreateModalModel struct {
	client *api.Client
	width  int
	height int

	// Form state
	inputs  []textInput
	focused int

	// Modal state
	completed bool
	err       error
}

// textInput represents a form input field
type textInput struct {
	label         string
	placeholder   string
	value         string
	required      bool
	cursor        int
	cursorVisible bool
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
	inputs := []textInput{
		{
			label:       "Session Name",
			placeholder: "my-session",
			required:    true,
		},
		{
			label:       "Description",
			placeholder: "Optional session description",
			required:    false,
		},
		{
			label:       "Project Path",
			placeholder: "/path/to/project (optional)",
			required:    false,
		},
	}

	return &CreateModalModel{
		client:  client,
		inputs:  inputs,
		focused: 0,
	}
}

// Init implements tea.Model
func (m *CreateModalModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return CursorBlinkMsg{}
	})
}

// Update implements tea.Model
func (m *CreateModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case CursorBlinkMsg:
		// Toggle cursor visibility
		if m.focused >= 0 && m.focused < len(m.inputs) {
			m.inputs[m.focused].cursorVisible = !m.inputs[m.focused].cursorVisible
		}
		return m, tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
			return CursorBlinkMsg{}
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Reset()
			return m, nil

		case "enter":
			return m.handleEnter()

		case "tab", "shift+tab":
			m.cycleFocus(msg.String() == "shift+tab")

		case "up":
			m.moveFocus(-1)

		case "down":
			m.moveFocus(1)

		case "backspace":
			m.deleteChar()

		case "left":
			m.moveCursor(-1)

		case "right":
			m.moveCursor(1)

		case "home":
			m.moveCursorToStart()

		case "end":
			m.moveCursorToEnd()

		default:
			// Handle character input
			if len(msg.String()) == 1 {
				m.insertChar(msg.String())
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *CreateModalModel) View() string {
	return m.renderModal()
}

// renderModal renders the modal content
func (m *CreateModalModel) renderModal() string {
	var content strings.Builder

	// Title
	title := styles.TitleStyle.Render("Create New Session")
	content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(title))
	content.WriteString("\n\n")

	// Form fields
	for i, input := range m.inputs {
		content.WriteString(m.renderInput(input, i == m.focused))
		content.WriteString("\n")
	}

	// Error display
	if m.err != nil {
		content.WriteString("\n")
		content.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n")
	}

	// Instructions
	content.WriteString("\n")
	content.WriteString(m.renderInstructions())

	return content.String()
}

// renderInput renders a single input field
func (m *CreateModalModel) renderInput(input textInput, focused bool) string {
	var field strings.Builder

	// Label
	labelStyle := styles.LabelStyle
	if focused {
		labelStyle = styles.FocusedStyle
	}

	label := input.label
	if input.required {
		label += " *"
	}
	field.WriteString(labelStyle.Render(label))
	field.WriteString("\n")

	// Input field
	inputValue := input.value
	if inputValue == "" && !focused {
		inputValue = input.placeholder
	}

	// Add cursor if focused
	if focused && input.cursorVisible {
		if input.cursor >= len(inputValue) {
			inputValue += "│"
		} else {
			inputValue = inputValue[:input.cursor] + "│" + inputValue[input.cursor:]
		}
	}

	// Style the input field
	inputStyle := styles.InputStyle
	if focused {
		inputStyle = styles.InputFocusedStyle
	}

	// Ensure minimum width
	minWidth := 40
	if len(inputValue) < minWidth {
		inputValue += strings.Repeat(" ", minWidth-len(inputValue))
	}

	field.WriteString(inputStyle.Render(inputValue))

	return field.String()
}

// renderInstructions renders the help instructions
func (m *CreateModalModel) renderInstructions() string {
	instructions := []string{
		"Tab/Shift+Tab: Navigate fields",
		"Enter: Create session",
		"Esc: Cancel",
	}

	return styles.DimTextStyle.Render(strings.Join(instructions, " • "))
}

// handleEnter handles the enter key press
func (m *CreateModalModel) handleEnter() (tea.Model, tea.Cmd) {
	// Validate required fields
	if err := m.validateForm(); err != nil {
		m.err = err
		return m, nil
	}

	// Create session
	return m, m.createSession()
}

// validateForm validates the form inputs
func (m *CreateModalModel) validateForm() error {
	for _, input := range m.inputs {
		if input.required && strings.TrimSpace(input.value) == "" {
			return fmt.Errorf("%s is required", input.label)
		}
	}
	return nil
}

// createSession creates a new session
func (m *CreateModalModel) createSession() tea.Cmd {
	return func() tea.Msg {
		name := strings.TrimSpace(m.inputs[inputName].value)
		description := strings.TrimSpace(m.inputs[inputDescription].value)
		projectPath := strings.TrimSpace(m.inputs[inputProjectPath].value)

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
	m.focused += delta
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	} else if m.focused >= len(m.inputs) {
		m.focused = 0
	}

	// Reset cursor position for newly focused field
	if m.focused >= 0 && m.focused < len(m.inputs) {
		m.inputs[m.focused].cursor = len(m.inputs[m.focused].value)
	}
}

// Text input methods

func (m *CreateModalModel) insertChar(char string) {
	if m.focused < 0 || m.focused >= len(m.inputs) {
		return
	}

	input := &m.inputs[m.focused]
	if input.cursor >= len(input.value) {
		input.value += char
	} else {
		input.value = input.value[:input.cursor] + char + input.value[input.cursor:]
	}
	input.cursor++
}

func (m *CreateModalModel) deleteChar() {
	if m.focused < 0 || m.focused >= len(m.inputs) {
		return
	}

	input := &m.inputs[m.focused]
	if input.cursor > 0 && len(input.value) > 0 {
		if input.cursor >= len(input.value) {
			input.value = input.value[:len(input.value)-1]
		} else {
			input.value = input.value[:input.cursor-1] + input.value[input.cursor:]
		}
		input.cursor--
	}
}

func (m *CreateModalModel) moveCursor(delta int) {
	if m.focused < 0 || m.focused >= len(m.inputs) {
		return
	}

	input := &m.inputs[m.focused]
	input.cursor += delta
	if input.cursor < 0 {
		input.cursor = 0
	} else if input.cursor > len(input.value) {
		input.cursor = len(input.value)
	}
}

func (m *CreateModalModel) moveCursorToStart() {
	if m.focused >= 0 && m.focused < len(m.inputs) {
		m.inputs[m.focused].cursor = 0
	}
}

func (m *CreateModalModel) moveCursorToEnd() {
	if m.focused >= 0 && m.focused < len(m.inputs) {
		m.inputs[m.focused].cursor = len(m.inputs[m.focused].value)
	}
}

// State management

func (m *CreateModalModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *CreateModalModel) IsCompleted() bool {
	return m.completed
}

func (m *CreateModalModel) Reset() {
	m.completed = false
	m.err = nil
	m.focused = 0

	// Clear all input values
	for i := range m.inputs {
		m.inputs[i].value = ""
		m.inputs[i].cursor = 0
		m.inputs[i].cursorVisible = false
	}
}

// Message types
type CursorBlinkMsg struct{}

type SessionCreatedMsg struct {
	Session *api.Session
	Error   error
}
