package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/layout"
	"claude-pilot/shared/styles"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionStats represents session statistics
type SessionStats struct {
	Total    int
	Active   int
	Inactive int
	Error    int
	Backend  string
}

// SummaryPanelModel represents the summary panel showing metrics
type SummaryPanelModel struct {
	client   *api.Client
	width    int
	sessions []*api.Session
	stats    SessionStats
}

// NewSummaryPanelModel creates a new summary panel model
func NewSummaryPanelModel(client *api.Client) *SummaryPanelModel {
	return &SummaryPanelModel{
		client: client,
		stats: SessionStats{
			Backend: client.GetBackend(),
		},
	}
}

// Init implements tea.Model
func (m *SummaryPanelModel) Init() tea.Cmd {
	return tea.Batch(
		m.refreshMetrics(),
		m.startPeriodicRefresh(),
	)
}

// Update implements tea.Model
func (m *SummaryPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case SessionsLoadedMsg:
		if msg.Error == nil {
			m.SetSessions(msg.Sessions)
		}

	case time.Time:
		// Auto-refresh timer - refresh data and restart timer
		return m, tea.Batch(
			m.refreshMetrics(),
			m.startPeriodicRefresh(),
		)
	}

	return m, nil
}

// View implements tea.Model
func (m *SummaryPanelModel) View() string {
	if m.width == 0 {
		return ""
	}

	return m.renderSummaryCards()
}

// renderSummaryCards renders summary cards with responsive layout using shared layout components
func (m *SummaryPanelModel) renderSummaryCards() string {
	// Calculate responsive width and size
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

// renderCompactView renders a compact view for small screens
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
		fmt.Sprintf("%d", m.stats.Total),
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
		fmt.Sprintf("%d", m.stats.Active),
		"active",
		styles.SuccessColor,
	)

	inactiveCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   15,
			Title:   "",
			Icon:    "⏸️",
			Border:  false,
			Compact: true,
		},
		fmt.Sprintf("%d", m.stats.Inactive),
		"inactive",
		styles.InfoColor,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		totalCard.Render(),
		" • ",
		activeCard.Render(),
		" • ",
		inactiveCard.Render(),
		" • ",
		styles.DimTextStyle.Render(fmt.Sprintf("%s backend", m.stats.Backend)),
	)
}

// renderMediumView renders a medium view for medium screens using flexbox
func (m *SummaryPanelModel) renderMediumView() string {
	cardWidth := (m.width - 10) / 4 // 4 cards with spacing
	if cardWidth < 12 {
		return m.renderCompactView()
	}

	// Create cards
	totalCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Sessions",
			Icon:    "📊",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Total),
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
		fmt.Sprintf("%d", m.stats.Active),
		"running",
		styles.SuccessColor,
	)

	inactiveCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Inactive",
			Icon:    "⏸️",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Inactive),
		"inactive",
		styles.InfoColor,
	)

	errorCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Errors",
			Icon:    "❌",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Error),
		"error",
		styles.ErrorColor,
	)

	// Use flexbox container for better card distribution
	cardContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: m.width, Height: 0, Padding: 0, Gap: 2},
		layout.FlexRow,
	).SetJustifyContent(layout.SpaceEvenly).SetAlignItems(layout.AlignStretch)

	// Add cards as flex items with equal distribution
	cardContainer.AddItem(layout.FlexItem{
		Content:    totalCard.Render(),
		FlexGrow:   1,
		FlexShrink: 1,
		Order:      1,
	})

	cardContainer.AddItem(layout.FlexItem{
		Content:    activeCard.Render(),
		FlexGrow:   1,
		FlexShrink: 1,
		Order:      2,
	})

	cardContainer.AddItem(layout.FlexItem{
		Content:    inactiveCard.Render(),
		FlexGrow:   1,
		FlexShrink: 1,
		Order:      3,
	})

	cardContainer.AddItem(layout.FlexItem{
		Content:    errorCard.Render(),
		FlexGrow:   1,
		FlexShrink: 1,
		Order:      4,
	})

	return cardContainer.Render()
}

// renderFullView renders the full view for large screens
func (m *SummaryPanelModel) renderFullView() string {
	cardWidth := (m.width - 12) / 4 // 4 cards with spacing
	if cardWidth < 12 {
		return m.renderMediumView()
	}

	// Create comprehensive summary cards
	totalCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Total",
			Icon:    "📊",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Total),
		"sessions",
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
		fmt.Sprintf("%d", m.stats.Active),
		"running",
		styles.SuccessColor,
	)

	inactiveCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Inactive",
			Icon:    "⏸️",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Inactive),
		"stopped",
		styles.WarningColor,
	)

	errorCard := components.NewSummaryCard(
		components.CardConfig{
			Width:   cardWidth,
			Title:   "Errors",
			Icon:    "❌",
			Border:  true,
			Compact: false,
		},
		fmt.Sprintf("%d", m.stats.Error),
		"failed",
		styles.ErrorColor,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		totalCard.Render(),
		" ",
		activeCard.Render(),
		" ",
		inactiveCard.Render(),
		" ",
		errorCard.Render(),
	)
}

// SetSessions updates the sessions and recalculates metrics
func (m *SummaryPanelModel) SetSessions(sessions []*api.Session) {
	m.sessions = sessions
	m.calculateStats()
}

// SetSize updates the panel size
func (m *SummaryPanelModel) SetSize(width, height int) {
	m.width = width
}

// calculateStats calculates summary statistics from sessions
func (m *SummaryPanelModel) calculateStats() {
	m.stats.Total = len(m.sessions)
	m.stats.Active = 0
	m.stats.Inactive = 0
	m.stats.Error = 0

	for _, session := range m.sessions {
		switch session.Status {
		case api.StatusActive:
			m.stats.Active++
		case api.StatusInactive:
			m.stats.Inactive++
		case api.StatusError:
			m.stats.Error++
		default:
			// Handle any other status as inactive
			m.stats.Inactive++
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

// startPeriodicRefresh starts a periodic refresh timer for real-time updates
func (m *SummaryPanelModel) startPeriodicRefresh() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return t
	})
}

// GetStats returns current statistics
func (m *SummaryPanelModel) GetStats() SessionStats {
	return m.stats
}
