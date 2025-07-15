package layout

import (
	"claude-pilot/shared/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LayoutConfig holds configuration for responsive layouts
type LayoutConfig struct {
	Width         int
	Height        int
	Padding       int
	Margin        int
	Gap           int
	MinPanelWidth int
}

// FlexDirection represents the direction of flex layout
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

// FlexContainer represents a flexible container layout
type FlexContainer struct {
	config    LayoutConfig
	direction FlexDirection
	children  []string
	weights   []int
}

// GridContainer represents a grid layout
type GridContainer struct {
	config  LayoutConfig
	columns int
	rows    int
	cells   [][]string
}

// Panel represents a single panel with responsive behavior
type Panel struct {
	config  LayoutConfig
	title   string
	content string
	border  bool
	focused bool
}

// NewFlexContainer creates a new flex container
func NewFlexContainer(config LayoutConfig, direction FlexDirection) *FlexContainer {
	return &FlexContainer{
		config:    config,
		direction: direction,
		children:  make([]string, 0),
		weights:   make([]int, 0),
	}
}

// NewGridContainer creates a new grid container
func NewGridContainer(config LayoutConfig, columns, rows int) *GridContainer {
	cells := make([][]string, rows)
	for i := range cells {
		cells[i] = make([]string, columns)
	}

	return &GridContainer{
		config:  config,
		columns: columns,
		rows:    rows,
		cells:   cells,
	}
}

// NewPanel creates a new panel
func NewPanel(config LayoutConfig, title, content string, border bool) *Panel {
	return &Panel{
		config:  config,
		title:   title,
		content: content,
		border:  border,
		focused: false,
	}
}

// AddChild adds a child element to the flex container
func (f *FlexContainer) AddChild(content string, weight int) {
	f.children = append(f.children, content)
	f.weights = append(f.weights, weight)
}

// Render renders the flex container
func (f *FlexContainer) Render() string {
	if len(f.children) == 0 {
		return ""
	}

	// Calculate available space
	availableWidth := f.config.Width - (f.config.Padding * 2) - (f.config.Margin * 2)
	availableHeight := f.config.Height - (f.config.Padding * 2) - (f.config.Margin * 2)

	// Calculate gap space
	gapSpace := f.config.Gap * (len(f.children) - 1)

	switch f.direction {
	case FlexRow:
		return f.renderRow(availableWidth-gapSpace, availableHeight)
	case FlexColumn:
		return f.renderColumn(availableWidth, availableHeight-gapSpace)
	default:
		return ""
	}
}

// renderRow renders children in a row layout
func (f *FlexContainer) renderRow(width, height int) string {
	if len(f.children) == 0 {
		return ""
	}

	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 1 {
		height = 1
	}

	// Calculate total weight
	totalWeight := 0
	for _, weight := range f.weights {
		totalWeight += weight
	}

	if totalWeight == 0 {
		// Equal distribution
		for i := range f.weights {
			f.weights[i] = 1
		}
		totalWeight = len(f.weights)
	}

	// Calculate widths for each child, distributing remainder pixels
	childWidths := make([]int, len(f.children))
	totalDistributed := 0

	// First pass: calculate base widths using integer division
	for i := range f.children {
		childWidths[i] = (width * f.weights[i]) / totalWeight
		totalDistributed += childWidths[i]
	}

	// Calculate remainder pixels and distribute them
	remainderPixels := width - totalDistributed
	for i := 0; i < remainderPixels && i < len(childWidths); i++ {
		childWidths[i]++
	}

	// Render each child with calculated width
	renderedChildren := make([]string, len(f.children))
	for i, child := range f.children {
		childWidth := childWidths[i]
		if childWidth < f.config.MinPanelWidth && f.config.MinPanelWidth > 0 {
			childWidth = f.config.MinPanelWidth
		}

		// Apply width constraint to child content
		childStyle := lipgloss.NewStyle().
			Width(childWidth).
			Height(height)

		renderedChildren[i] = childStyle.Render(child)
	}

	// Join with gaps
	gapStyle := lipgloss.NewStyle().Width(f.config.Gap)
	gap := gapStyle.Render("")

	result := strings.Join(renderedChildren, gap)

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Padding(f.config.Padding).
		Margin(f.config.Margin)

	return containerStyle.Render(result)
}

// renderColumn renders children in a column layout
func (f *FlexContainer) renderColumn(width, height int) string {
	if len(f.children) == 0 {
		return ""
	}

	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 3 {
		height = 3
	}

	// Calculate total weight
	totalWeight := 0
	for _, weight := range f.weights {
		totalWeight += weight
	}

	if totalWeight == 0 {
		// Equal distribution
		for i := range f.weights {
			f.weights[i] = 1
		}
		totalWeight = len(f.weights)
	}

	// Calculate heights for each child, distributing remainder pixels
	childHeights := make([]int, len(f.children))
	totalDistributed := 0

	// First pass: calculate base heights using integer division
	for i := range f.children {
		childHeights[i] = (height * f.weights[i]) / totalWeight
		totalDistributed += childHeights[i]
	}

	// Calculate remainder pixels and distribute them
	remainderPixels := height - totalDistributed
	for i := 0; i < remainderPixels && i < len(childHeights); i++ {
		childHeights[i]++
	}

	// Render each child with calculated height
	renderedChildren := make([]string, len(f.children))
	for i, child := range f.children {
		childHeight := childHeights[i]

		// Apply height constraint to child content
		childStyle := lipgloss.NewStyle().
			Width(width).
			Height(childHeight)

		renderedChildren[i] = childStyle.Render(child)
	}

	// Join vertically with gaps
	var result strings.Builder
	for i, child := range renderedChildren {
		if i > 0 {
			// Add gap lines
			for j := 0; j < f.config.Gap; j++ {
				result.WriteString("\n")
			}
		}
		result.WriteString(child)
	}

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Padding(f.config.Padding).
		Margin(f.config.Margin)

	return containerStyle.Render(result.String())
}

// SetCell sets content for a specific grid cell
func (g *GridContainer) SetCell(row, col int, content string) {
	if row >= 0 && row < g.rows && col >= 0 && col < g.columns {
		g.cells[row][col] = content
	}
}

// Render renders the grid container
func (g *GridContainer) Render() string {
	if g.rows == 0 || g.columns == 0 {
		return ""
	}

	// Calculate cell dimensions
	availableWidth := g.config.Width - (g.config.Padding * 2) - (g.config.Margin * 2)
	availableHeight := g.config.Height - (g.config.Padding * 2) - (g.config.Margin * 2)

	cellWidth := (availableWidth - (g.config.Gap * (g.columns - 1))) / g.columns
	cellHeight := (availableHeight - (g.config.Gap * (g.rows - 1))) / g.rows

	// Render each row
	renderedRows := make([]string, g.rows)
	for i := 0; i < g.rows; i++ {
		// Render cells in this row
		renderedCells := make([]string, g.columns)
		for j := 0; j < g.columns; j++ {
			cellStyle := lipgloss.NewStyle().
				Width(cellWidth).
				Height(cellHeight)

			renderedCells[j] = cellStyle.Render(g.cells[i][j])
		}

		// Join cells with horizontal gaps
		gapStyle := lipgloss.NewStyle().Width(g.config.Gap)
		gap := gapStyle.Render("")
		renderedRows[i] = strings.Join(renderedCells, gap)
	}

	// Join rows with vertical gaps
	var result strings.Builder
	for i, row := range renderedRows {
		if i > 0 {
			// Add gap lines
			for j := 0; j < g.config.Gap; j++ {
				result.WriteString("\n")
			}
		}
		result.WriteString(row)
	}

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Padding(g.config.Padding).
		Margin(g.config.Margin)

	return containerStyle.Render(result.String())
}

// SetFocused sets the focused state of the panel
func (p *Panel) SetFocused(focused bool) {
	p.focused = focused
}

// SetContent updates the panel content
func (p *Panel) SetContent(content string) {
	p.content = content
}

// SetTitle updates the panel title
func (p *Panel) SetTitle(title string) {
	p.title = title
}

// Render renders the panel
func (p *Panel) Render() string {
	// Prepare content
	var content strings.Builder

	// Add title if provided
	if p.title != "" {
		titleStyle := styles.PanelHeaderStyle
		if p.focused {
			titleStyle = titleStyle.Foreground(styles.ClaudePrimary)
		}
		content.WriteString(titleStyle.Render(p.title) + "\n")
	}

	// Add main content
	content.WriteString(p.content)

	// Apply panel styling
	panelStyle := lipgloss.NewStyle().
		Padding(p.config.Padding).
		Margin(p.config.Margin)

	if p.config.Width > 0 {
		panelStyle = panelStyle.Width(p.config.Width)
	}

	if p.config.Height > 0 {
		panelStyle = panelStyle.Height(p.config.Height)
	}

	if p.border {
		borderColor := styles.TextMuted
		if p.focused {
			borderColor = styles.ClaudePrimary
		}
		panelStyle = panelStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	return panelStyle.Render(content.String())
}

// ResponsiveLayout creates an adaptive layout based on terminal size
func ResponsiveLayout(width, height int, panels []string) string {
	layoutWidth, size := styles.GetResponsiveWidth(width)

	switch size {
	case "small":
		// Stack panels vertically for small screens
		container := NewFlexContainer(
			LayoutConfig{Width: layoutWidth, Height: height, Padding: 1, Gap: 1},
			FlexColumn,
		)
		for _, panel := range panels {
			container.AddChild(panel, 1)
		}
		return container.Render()

	case "medium":
		// Two-column layout for medium screens
		if len(panels) >= 2 {
			container := NewFlexContainer(
				LayoutConfig{Width: layoutWidth, Height: height, Padding: 1, Gap: 2, MinPanelWidth: 30},
				FlexRow,
			)
			// Main content (wider) and sidebar (narrower)
			container.AddChild(panels[0], 2)
			if len(panels) > 1 {
				// Combine remaining panels in the second column
				sidebarContainer := NewFlexContainer(
					LayoutConfig{Width: layoutWidth / 3, Height: height - 4, Padding: 0, Gap: 1},
					FlexColumn,
				)
				for i := 1; i < len(panels); i++ {
					sidebarContainer.AddChild(panels[i], 1)
				}
				container.AddChild(sidebarContainer.Render(), 1)
			}
			return container.Render()
		}
		return panels[0]

	default: // large
		// Three-column layout for large screens
		if len(panels) >= 3 {
			container := NewFlexContainer(
				LayoutConfig{Width: layoutWidth, Height: height, Padding: 1, Gap: 2, MinPanelWidth: 25},
				FlexRow,
			)
			container.AddChild(panels[0], 1) // Left panel
			container.AddChild(panels[1], 2) // Center panel (wider)
			container.AddChild(panels[2], 1) // Right panel
			return container.Render()
		} else if len(panels) >= 2 {
			// Fall back to two-column layout
			container := NewFlexContainer(
				LayoutConfig{Width: layoutWidth, Height: height, Padding: 1, Gap: 2, MinPanelWidth: 30},
				FlexRow,
			)
			container.AddChild(panels[0], 2)
			container.AddChild(panels[1], 1)
			return container.Render()
		}
		return panels[0]
	}
}

// DashboardLayout creates a dashboard-style layout with header, main content, and footer
func DashboardLayout(width, height int, header, main, footer string) string {
	// Calculate fixed heights for header and footer
	headerHeight := calculateContentHeight(header)
	footerHeight := calculateContentHeight(footer)

	// Ensure reasonable minimum and maximum heights (more compact)
	if headerHeight < 2 {
		headerHeight = 2
	}
	if headerHeight > 5 { // Reduced from 6 to 5
		headerHeight = 5
	}
	if footerHeight < 1 {
		footerHeight = 1
	}
	if footerHeight > 2 { // Reduced from 3 to 2
		footerHeight = 2
	}

	// Calculate remaining height for main content
	mainHeight := height - headerHeight - footerHeight
	if mainHeight < 6 { // Reduced from 5 to 6 for better fit
		// If not enough space, reduce header height
		headerHeight = max(2, height-footerHeight-6)
		mainHeight = height - headerHeight - footerHeight
	}

	// Ensure main content doesn't get too tall
	if mainHeight > height*3/4 {
		mainHeight = height * 3 / 4
	}

	// Create styled sections with proper heights
	headerSection := lipgloss.NewStyle().
		Width(width).
		Height(headerHeight).
		Render(header)

	mainSection := lipgloss.NewStyle().
		Width(width).
		Height(mainHeight).
		Render(main)

	footerSection := lipgloss.NewStyle().
		Width(width).
		Height(footerHeight).
		Render(footer)

	// Join vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerSection,
		mainSection,
		footerSection,
	)
}

// calculateContentHeight estimates the height needed for content
func calculateContentHeight(content string) int {
	if content == "" {
		return 0
	}

	lines := strings.Split(content, "\n")
	return len(lines)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SidebarLayout creates a layout with a sidebar and main content area
func SidebarLayout(width, height int, sidebar, main string, sidebarWidth int) string {
	container := NewFlexContainer(
		LayoutConfig{Width: width, Height: height, Padding: 1, Gap: 2, MinPanelWidth: sidebarWidth},
		FlexRow,
	)

	mainWidth := width - sidebarWidth - 6 // Account for padding and gaps
	sidebarWeight := (sidebarWidth * 10) / width
	mainWeight := (mainWidth * 10) / width

	container.AddChild(sidebar, sidebarWeight)
	container.AddChild(main, mainWeight)

	return container.Render()
}
