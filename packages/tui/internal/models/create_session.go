package models

import (
	"strings"

	"claude-pilot/core/api"
	"claude-pilot/tui/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CreateSessionModel handles the create session form
type CreateSessionModel struct {
	client     *api.Client
	width      int
	height     int
	focusIndex int
	inputs     []string
	labels     []string
	creating   bool
	err        error
}

// SessionCreatedMsg is sent when a session is created
type SessionCreatedMsg struct {
	Session *api.Session
	Error   error
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Error error
}

// NewCreateSessionModel creates a new create session model
func NewCreateSessionModel(client *api.Client) *CreateSessionModel {
	return &CreateSessionModel{
		client: client,
		inputs: make([]string, 3), // name, description, project path
		labels: []string{
			"Session Name:",
			"Description:",
			"Project Path:",
		},
	}
}

// Init implements tea.Model
func (m *CreateSessionModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *CreateSessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.creating {
			return m, nil
		}

		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			// Navigate between inputs
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) - 1
				}
			} else {
				m.focusIndex++
				if m.focusIndex >= len(m.inputs) {
					m.focusIndex = 0
				}
			}

		case "enter":
			// Create session
			if m.inputs[0] != "" { // Name is required
				m.creating = true
				return m, m.createSession()
			}

		case "backspace":
			// Remove character from current input
			if len(m.inputs[m.focusIndex]) > 0 {
				m.inputs[m.focusIndex] = m.inputs[m.focusIndex][:len(m.inputs[m.focusIndex])-1]
			}

		default:
			// Add character to current input
			if len(msg.String()) == 1 {
				m.inputs[m.focusIndex] += msg.String()
			}
		}

	case SessionCreatedMsg:
		m.creating = false
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			// Reset form
			m.inputs = make([]string, 3)
			m.focusIndex = 0
			m.err = nil
		}
		return m, func() tea.Msg { return msg }
	}

	return m, nil
}

// View implements tea.Model
func (m *CreateSessionModel) View() string {
	if m.creating {
		return styles.InfoStyle.Render("Creating session...")
	}

	var content strings.Builder

	// Header
	content.WriteString(styles.HeaderStyle.Render("Create New Session"))
	content.WriteString("\n\n")

	// Form inputs
	for i, label := range m.labels {
		// Label
		content.WriteString(styles.LabelStyle.Render(label))
		content.WriteString("\n")

		// Input field
		inputValue := m.inputs[i]
		if i == 0 && inputValue == "" {
			inputValue = "session-name" // Placeholder for required field
		}

		var inputStyle lipgloss.Style
		if i == m.focusIndex {
			inputStyle = styles.InputFocusedStyle
		} else {
			inputStyle = styles.InputStyle
		}

		// Show cursor for focused input
		displayValue := inputValue
		if i == m.focusIndex {
			displayValue += "█"
		}

		content.WriteString(inputStyle.Render(displayValue))
		content.WriteString("\n\n")
	}

	// Error display
	if m.err != nil {
		content.WriteString(styles.ErrorStyle.Render(m.err.Error()))
		content.WriteString("\n\n")
	}

	// Instructions
	content.WriteString(styles.MutedTextStyle.Render("↑/↓: navigate fields • enter: create • esc: cancel"))

	return content.String()
}

// SetSize updates the model's size
func (m *CreateSessionModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// createSession creates a new session using the API
func (m *CreateSessionModel) createSession() tea.Cmd {
	return func() tea.Msg {
		// Use current directory if project path is empty
		projectPath := m.inputs[2]
		if projectPath == "" {
			projectPath = api.GetProjectPath("")
		}

		session, err := m.client.CreateSession(api.CreateSessionRequest{
			Name:        m.inputs[0],
			Description: m.inputs[1],
			ProjectPath: projectPath,
		})

		return SessionCreatedMsg{
			Session: session,
			Error:   err,
		}
	}
}
