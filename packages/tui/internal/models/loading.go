package models

import (
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingModel represents a loading state with spinner
type LoadingModel struct {
	spinner spinner.Model
	message string
	width   int
	height  int
}

// NewLoadingModel creates a new loading model
func NewLoadingModel() *LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(styles.ClaudePrimary).
		Bold(true)

	return &LoadingModel{
		spinner: s,
		message: "Loading...",
	}
}

// NewLoadingModelWithMessage creates a new loading model with custom message
func NewLoadingModelWithMessage(message string) *LoadingModel {
	model := NewLoadingModel()
	model.message = message
	return model
}

// Init implements tea.Model
func (m *LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model
func (m *LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Loading screen doesn't handle key input by default
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View implements tea.Model
func (m *LoadingModel) View() string {
	if m.width == 0 || m.height == 0 {
		return m.spinner.View() + " " + m.message
	}

	// Create centered loading display
	spinnerAndMessage := m.spinner.View() + " " + m.message

	// Center horizontally and vertically
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		spinnerAndMessage,
	)
}

// SetMessage updates the loading message
func (m *LoadingModel) SetMessage(message string) {
	m.message = message
}

// SetSize updates the loading screen size
func (m *LoadingModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// LoadingOverlay creates a loading overlay over existing content
func LoadingOverlay(background, message string, width, height int) string {
	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(styles.ClaudePrimary).
		Bold(true)

	// Create loading message with spinner
	loadingContent := s.View() + " " + message

	// Create overlay box
	overlay := lipgloss.NewStyle().
		Width(40).
		Height(5).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ClaudePrimary).
		Background(styles.BackgroundPrimary).
		Padding(1).
		Align(lipgloss.Center)

	// Place overlay on top of background
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		overlay.Render(loadingContent),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.BackgroundSecondary),
	)
}

// LoadingMessage types for different loading states
type LoadingMessage struct {
	Message string
	Type    LoadingType
}

type LoadingType int

const (
	LoadingTypeGeneral LoadingType = iota
	LoadingTypeSessions
	LoadingTypeCreate
	LoadingTypeAttach
	LoadingTypeKill
	LoadingTypeRefresh
)

// GetLoadingMessage returns an appropriate loading message for the type
func GetLoadingMessage(loadingType LoadingType) string {
	switch loadingType {
	case LoadingTypeSessions:
		return "Loading sessions..."
	case LoadingTypeCreate:
		return "Creating session..."
	case LoadingTypeAttach:
		return "Attaching to session..."
	case LoadingTypeKill:
		return "Terminating session..."
	case LoadingTypeRefresh:
		return "Refreshing data..."
	default:
		return "Loading..."
	}
}

// LoadingStates for different operations
var LoadingStates = map[string]string{
	"sessions":   "Loading sessions...",
	"create":     "Creating session...",
	"attach":     "Attaching to session...",
	"kill":       "Terminating session...",
	"refresh":    "Refreshing data...",
	"connecting": "Connecting to backend...",
	"starting":   "Starting session...",
	"stopping":   "Stopping session...",
	"saving":     "Saving configuration...",
	"loading":    "Loading configuration...",
}

// CreateLoadingSpinner creates a styled spinner for consistent loading states
func CreateLoadingSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(styles.ClaudePrimary).
		Bold(true)
	return s
}

// LoadingBox creates a styled loading box with message
func LoadingBox(message string) string {
	s := CreateLoadingSpinner()

	content := fmt.Sprintf("%s %s", s.View(), message)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ClaudePrimary).
		Background(styles.BackgroundSecondary).
		Padding(1, 2).
		Align(lipgloss.Center).
		Render(content)
}

// InlineLoading creates a simple inline loading indicator
func InlineLoading(message string) string {
	return styles.SpinnerStyle.Render("⠋") + " " + styles.SecondaryTextStyle.Render(message)
}

// LoadingBar creates a simple loading bar (animated dots)
func LoadingBar(step int) string {
	dots := strings.Repeat("●", step%4) + strings.Repeat("○", 4-(step%4))
	return lipgloss.NewStyle().Foreground(styles.ClaudePrimary).Render(dots)
}
