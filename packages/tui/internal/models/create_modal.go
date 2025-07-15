package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	keys      KeyMap
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
	inputs[inputName].Width = 40
	inputs[inputName] = styles.ConfigureBubblesTextInput(inputs[inputName])

	// Description input (optional)
	inputs[inputDescription] = textinput.New()
	inputs[inputDescription].Placeholder = "Optional session description"
	inputs[inputDescription].Width = 40
	inputs[inputDescription] = styles.ConfigureBubblesTextInput(inputs[inputDescription])

	// Project path input (optional)
	inputs[inputProjectPath] = textinput.New()
	inputs[inputProjectPath].Placeholder = "/path/to/project (optional)"
	inputs[inputProjectPath].Width = 40
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

		// Update input widths if needed
		for i := range m.inputs {
			m.inputs[i].Width = min(40, m.width-10)
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
			}
		}
	}

	// Update all inputs for cursor blinking
	for i := range m.inputs {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
	for i := range m.inputs {
		content.WriteString(m.renderInput(i))
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

// renderInput renders a single input field using Bubbles textinput
func (m *CreateModalModel) renderInput(index int) string {
	var field strings.Builder

	// Label
	labelStyle := styles.LabelStyle
	if index == m.focused {
		labelStyle = styles.FocusedStyle
	}

	field.WriteString(labelStyle.Render(m.labels[index]))
	field.WriteString("\n")

	// Input field - rendered by Bubbles textinput
	field.WriteString(m.inputs[index].View())

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
	// Check if session name is provided (required field)
	if strings.TrimSpace(m.inputs[inputName].Value()) == "" {
		return fmt.Errorf("Session Name is required")
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
}

func (m *CreateModalModel) IsCompleted() bool {
	return m.completed
}

func (m *CreateModalModel) Reset() {
	m.completed = false
	m.err = nil
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

// Message types
type SessionCreatedMsg struct {
	Session *api.Session
	Error   error
}
