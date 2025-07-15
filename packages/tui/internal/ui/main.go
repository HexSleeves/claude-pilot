package ui

import (
	"claude-pilot/core/api"
	"claude-pilot/tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// MainModel is the root model for the TUI application
type MainModel struct {
	client    *api.Client
	dashboard *models.DashboardModel
	width     int
	height    int
}

// NewMainModel creates a new main TUI model
func NewMainModel(client *api.Client) *MainModel {
	return &MainModel{
		client:    client,
		dashboard: models.NewDashboardModel(client),
	}
}

// Init implements tea.Model
func (m *MainModel) Init() tea.Cmd {
	return m.dashboard.Init()
}

// Update implements tea.Model
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard.SetSize(msg.Width, msg.Height)
	}

	// Forward all messages to the dashboard
	newDashboard, cmd := m.dashboard.Update(msg)
	if dashboard, ok := newDashboard.(*models.DashboardModel); ok {
		m.dashboard = dashboard
	}

	return m, cmd
}

// View implements tea.Model
func (m *MainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading dashboard..."
	}

	return m.dashboard.View()
}
