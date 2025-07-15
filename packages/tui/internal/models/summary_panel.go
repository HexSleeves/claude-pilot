package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/styles"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SummaryPanelModel represents the summary panel showing metrics
type SummaryPanelModel struct {
	client   *api.Client
	width    int
	height   int
	sessions []*api.Session

	// Metrics
	totalSessions     int
	activeSessions    int
	connectedSessions int
	backend           string
	lastUpdated       time.Time
}

// NewSummaryPanelModel creates a new summary panel model
func NewSummaryPanelModel(client *api.Client) *SummaryPanelModel {
	return &SummaryPanelModel{
		client:      client,
		backend:     client.GetBackend(),
		lastUpdated: time.Now(),
	}
}

// Init implements tea.Model
func (m *SummaryPanelModel) Init() tea.Cmd {
	return m.refreshMetrics()
}

// Update implements tea.Model
func (m *SummaryPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case SessionsLoadedMsg:
		if msg.Error == nil {
			m.SetSessions(msg.Sessions)
		}

	case time.Time:
		// Auto-refresh timer
		return m, m.refreshMetrics()
	}

	return m, nil
}

// View implements tea.Model
func (m *SummaryPanelModel) View() string {
	if m.width == 0 {
		return ""
	}

	// Calculate card width for responsive layout
	_, size := styles.GetResponsiveWidth(m.width)

	switch size {
	case "small":
		return m.renderCompactView()
	case "medium":
		return m.renderMediumView()
	default:
		return m.renderFullView()
	}
}

// RenderCompactView renders a compact view for small screens
func (m *SummaryPanelModel) renderCompactView() string {
	// Single row with key metrics
	totalCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   15,
			Title:   "",
			Icon:    "📊",
			Border:  false,
			Compact: true,
		},
		fmt.Sprintf("%d", m.totalSessions),
		"total",
		styles.ClaudePrimary,
	)

	activeCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   15,
			Title:   "",
			Icon:    "✅",
			Border:  false,
			Compact: true,
		},
		fmt.Sprintf("%d", m.activeSessions),
		"active",
		styles.SuccessColor,
	)

	connectedCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   15,
			Title:   "",
			Icon:    "🔗",
			Border:  false,
			Compact: true,
		},
		fmt.Sprintf("%d", m.connectedSessions),
		"connected",
		styles.InfoColor,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		totalCard.Render(),
		" • ",
		activeCard.Render(),
		" • ",
		connectedCard.Render(),
		" • ",
		styles.DimTextStyle.Render(fmt.Sprintf("%s backend", m.backend)),
	)
}

// renderMediumView renders a medium view for medium screens
func (m *SummaryPanelModel) renderMediumView() string {
	cardWidth := (m.width - 10) / 3 // 3 cards with spacing
	if cardWidth < 12 {
		return m.renderCompactView()
	}

	totalCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Sessions",
			Icon:    "📊",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.totalSessions),
		"total",
		styles.ClaudePrimary,
	)

	activeCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Active",
			Icon:    "✅",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.activeSessions),
		"running",
		styles.SuccessColor,
	)

	systemCard := components.NewStatusCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "System",
			Icon:    "⚙️",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%s Ready", m.backend),
		"●",
		[]string{
			fmt.Sprintf("Connected: %d", m.connectedSessions),
			fmt.Sprintf("Updated: %s", m.lastUpdated.Format("15:04:05")),
		},
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		totalCard.Render(),
		"  ",
		activeCard.Render(),
		"  ",
		systemCard.Render(),
	)
}

// renderFullView renders the full view for large screens
func (m *SummaryPanelModel) renderFullView() string {
	// Create comprehensive summary cards
	cards := components.CreateSessionSummaryCards(
		m.totalSessions,
		m.activeSessions,
		m.connectedSessions,
		m.backend,
		m.width,
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}

// SetSessions updates the sessions and recalculates metrics
func (m *SummaryPanelModel) SetSessions(sessions []*api.Session) {
	m.sessions = sessions
	m.calculateMetrics()
	m.lastUpdated = time.Now()
}

// SetSize updates the panel size
func (m *SummaryPanelModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// calculateMetrics calculates summary metrics from sessions
func (m *SummaryPanelModel) calculateMetrics() {
	m.totalSessions = len(m.sessions)
	m.activeSessions = 0
	m.connectedSessions = 0

	for _, session := range m.sessions {
		switch session.Status {
		case api.StatusActive:
			m.activeSessions++
		case api.StatusConnected:
			m.connectedSessions++
			m.activeSessions++ // Connected sessions are also active
		}
	}
}

// refreshMetrics refreshes the metrics by loading sessions
func (m *SummaryPanelModel) refreshMetrics() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.client.ListSessions()
		return SessionsLoadedMsg{Sessions: sessions, Error: err}
	}
}

// GetMetrics returns current metrics
func (m *SummaryPanelModel) GetMetrics() (total, active, connected int) {
	return m.totalSessions, m.activeSessions, m.connectedSessions
}

// GetLastUpdated returns the last update time
func (m *SummaryPanelModel) GetLastUpdated() time.Time {
	return m.lastUpdated
}
