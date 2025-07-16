package layout

import (
	"claude-pilot/shared/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Package-level constants for layout constraints
const (
	MinHeaderHeight      = 2
	MaxHeaderHeight      = 5
	MinFooterHeight      = 1
	MaxFooterHeight      = 2
	MinMainContentHeight = 6
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

// JustifyContent represents the main-axis alignment
type JustifyContent int

const (
	FlexStart JustifyContent = iota
	FlexEnd
	Center
	SpaceBetween
	SpaceAround
	SpaceEvenly
)

// AlignItems represents the cross-axis alignment
type AlignItems int

const (
	AlignStart AlignItems = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// FlexWrap represents the wrap behavior
type FlexWrap int

const (
	NoWrap FlexWrap = iota
	Wrap
	WrapReverse
)

// FlexItem represents an individual flex item with CSS flexbox properties
type FlexItem struct {
	Content    string      // The content to render
	FlexGrow   int         // How much the item should grow (default: 0)
	FlexShrink int         // How much the item should shrink (default: 1)
	FlexBasis  int         // The initial main size (default: auto/-1)
	AlignSelf  *AlignItems // Override container's align-items for this item
	Order      int         // Display order (default: 0)
}

// FlexContainer represents a flexible container layout with CSS flexbox properties
type FlexContainer struct {
	config         LayoutConfig
	direction      FlexDirection
	justifyContent JustifyContent
	alignItems     AlignItems
	flexWrap       FlexWrap
	items          []FlexItem
	// Backward compatibility fields
	children []string
	weights  []int
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
		config:         config,
		direction:      direction,
		justifyContent: FlexStart,
		alignItems:     AlignStretch,
		flexWrap:       NoWrap,
		items:          make([]FlexItem, 0),
		children:       make([]string, 0),
		weights:        make([]int, 0),
	}
}

// SetJustifyContent sets the justify-content property for chainable configuration
func (f *FlexContainer) SetJustifyContent(justify JustifyContent) *FlexContainer {
	f.justifyContent = justify
	return f
}

// SetAlignItems sets the align-items property for chainable configuration
func (f *FlexContainer) SetAlignItems(align AlignItems) *FlexContainer {
	f.alignItems = align
	return f
}

// SetFlexWrap sets the flex-wrap property for chainable configuration
func (f *FlexContainer) SetFlexWrap(wrap FlexWrap) *FlexContainer {
	f.flexWrap = wrap
	return f
}

// AddItem adds a FlexItem to the container using the new API
func (f *FlexContainer) AddItem(item FlexItem) {
	f.items = append(f.items, item)
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

// AddChild adds a child element to the flex container (backward compatibility wrapper)
func (f *FlexContainer) AddChild(content string, weight int) {
	// Convert to FlexItem for new system
	item := FlexItem{
		Content:    content,
		FlexGrow:   weight,
		FlexShrink: 1,
		FlexBasis:  -1, // auto
		Order:      0,
	}
	f.items = append(f.items, item)

	// Keep old fields for any legacy code that might access them directly
	f.children = append(f.children, content)
	f.weights = append(f.weights, weight)
}

// Render renders the flex container using CSS flexbox algorithms
func (f *FlexContainer) Render() string {
	if len(f.items) == 0 {
		return ""
	}

	// Calculate available space
	availableWidth := f.config.Width - (f.config.Padding * 2) - (f.config.Margin * 2)
	availableHeight := f.config.Height - (f.config.Padding * 2) - (f.config.Margin * 2)

	// Sort items by Order property
	sortedItems := make([]FlexItem, len(f.items))
	copy(sortedItems, f.items)
	f.sortItemsByOrder(sortedItems)

	// Calculate gap space
	gapSpace := f.config.Gap * (len(sortedItems) - 1)

	var result string
	switch f.direction {
	case FlexRow:
		result = f.renderFlexRow(sortedItems, availableWidth-gapSpace, availableHeight)
	case FlexColumn:
		result = f.renderFlexColumn(sortedItems, availableWidth, availableHeight-gapSpace)
	default:
		return ""
	}

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Padding(f.config.Padding).
		Margin(f.config.Margin)

	return containerStyle.Render(result)
}

// sortItemsByOrder sorts flex items by their Order property
func (f *FlexContainer) sortItemsByOrder(items []FlexItem) {
	// Simple bubble sort by Order property
	n := len(items)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if items[j].Order > items[j+1].Order {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}
}

// renderFlexRow renders items in a row layout using CSS flexbox algorithms
func (f *FlexContainer) renderFlexRow(items []FlexItem, width, height int) string {
	if len(items) == 0 {
		return ""
	}

	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 1 {
		height = 1
	}

	// Calculate flex-basis for each item
	itemWidths := make([]int, len(items))
	totalBasis := 0
	totalGrow := 0
	totalShrink := 0

	for i, item := range items {
		if item.FlexBasis > 0 {
			itemWidths[i] = item.FlexBasis
		} else {
			// Auto basis - use content width or equal distribution
			itemWidths[i] = width / len(items)
		}
		totalBasis += itemWidths[i]
		totalGrow += item.FlexGrow
		if item.FlexShrink > 0 {
			totalShrink += item.FlexShrink
		} else {
			totalShrink += 1 // Default shrink is 1
		}
	}

	// Distribute remaining space using flex-grow or flex-shrink
	remainingSpace := width - totalBasis
	if remainingSpace > 0 && totalGrow > 0 {
		// Distribute extra space using flex-grow
		for i, item := range items {
			if item.FlexGrow > 0 {
				extraWidth := (remainingSpace * item.FlexGrow) / totalGrow
				itemWidths[i] += extraWidth
			}
		}
	} else if remainingSpace < 0 && totalShrink > 0 {
		// Shrink items using flex-shrink
		shrinkSpace := -remainingSpace
		for i, item := range items {
			shrinkFactor := item.FlexShrink
			if shrinkFactor == 0 {
				shrinkFactor = 1
			}
			shrinkAmount := (shrinkSpace * shrinkFactor) / totalShrink
			itemWidths[i] = max(0, itemWidths[i]-shrinkAmount)
		}
	}

	// Apply minimum panel width constraint
	for i := range itemWidths {
		if itemWidths[i] < f.config.MinPanelWidth && f.config.MinPanelWidth > 0 {
			itemWidths[i] = f.config.MinPanelWidth
		}
	}

	return f.renderRowWithJustifyContent(items, itemWidths, width, height)
}

// renderRowWithJustifyContent applies justify-content spacing to row items
func (f *FlexContainer) renderRowWithJustifyContent(items []FlexItem, itemWidths []int, containerWidth, height int) string {
	// Calculate total item width
	totalItemWidth := 0
	for _, width := range itemWidths {
		totalItemWidth += width
	}

	// Calculate available space for justify-content
	availableSpace := containerWidth - totalItemWidth

	// Render each item with proper alignment
	renderedItems := make([]string, len(items))
	for i, item := range items {
		// Determine cross-axis alignment for this item
		align := f.alignItems
		if item.AlignSelf != nil {
			align = *item.AlignSelf
		}

		// Calculate cross-axis positioning
		var hAlign, vAlign lipgloss.Position
		switch align {
		case AlignStart:
			vAlign = lipgloss.Top
		case AlignCenter:
			vAlign = lipgloss.Center
		case AlignEnd:
			vAlign = lipgloss.Bottom
		case AlignStretch:
			vAlign = lipgloss.Top // Will use full height
		}
		hAlign = lipgloss.Left

		// Apply width and height constraints
		itemStyle := lipgloss.NewStyle().Width(itemWidths[i])
		if align == AlignStretch {
			itemStyle = itemStyle.Height(height)
		}

		// Render the item content with alignment
		styledContent := itemStyle.Render(item.Content)
		if align != AlignStretch {
			// Use lipgloss.Place for non-stretch alignment
			renderedItems[i] = lipgloss.Place(itemWidths[i], height, hAlign, vAlign, styledContent)
		} else {
			renderedItems[i] = styledContent
		}
	}

	// Apply justify-content spacing
	return f.applyJustifyContentSpacing(renderedItems, availableSpace)
}

// applyJustifyContentSpacing applies justify-content spacing to rendered items
func (f *FlexContainer) applyJustifyContentSpacing(renderedItems []string, availableSpace int) string {
	if len(renderedItems) == 0 {
		return ""
	}

	switch f.justifyContent {
	case FlexStart:
		// Items at start, gaps between items
		return f.joinItemsWithGaps(renderedItems)

	case FlexEnd:
		// Items at end, add space at beginning
		if availableSpace > 0 {
			leadingSpace := lipgloss.NewStyle().Width(availableSpace).Render("")
			return leadingSpace + f.joinItemsWithGaps(renderedItems)
		}
		return f.joinItemsWithGaps(renderedItems)

	case Center:
		// Items centered, equal space on both sides
		if availableSpace > 0 {
			sideSpace := availableSpace / 2
			leadingSpace := lipgloss.NewStyle().Width(sideSpace).Render("")
			return leadingSpace + f.joinItemsWithGaps(renderedItems)
		}
		return f.joinItemsWithGaps(renderedItems)

	case SpaceBetween:
		// Equal space between items, no space at edges
		if len(renderedItems) == 1 {
			return renderedItems[0]
		}
		if availableSpace > 0 {
			spaceBetween := availableSpace / (len(renderedItems) - 1)
			spaceStyle := lipgloss.NewStyle().Width(spaceBetween)
			space := spaceStyle.Render("")

			var result strings.Builder
			for i, item := range renderedItems {
				if i > 0 {
					result.WriteString(space)
				}
				if f.config.Gap > 0 && i > 0 {
					gapStyle := lipgloss.NewStyle().Width(f.config.Gap)
					result.WriteString(gapStyle.Render(""))
				}
				result.WriteString(item)
			}
			return result.String()
		}
		return f.joinItemsWithGaps(renderedItems)

	case SpaceAround:
		// Equal space around each item
		if availableSpace > 0 {
			spaceAround := availableSpace / len(renderedItems)
			halfSpace := spaceAround / 2
			spaceStyle := lipgloss.NewStyle().Width(halfSpace)
			space := spaceStyle.Render("")

			var result strings.Builder
			for i, item := range renderedItems {
				result.WriteString(space)
				if f.config.Gap > 0 && i > 0 {
					gapStyle := lipgloss.NewStyle().Width(f.config.Gap)
					result.WriteString(gapStyle.Render(""))
				}
				result.WriteString(item)
				result.WriteString(space)
			}
			return result.String()
		}
		return f.joinItemsWithGaps(renderedItems)

	case SpaceEvenly:
		// Equal space between and around items
		if availableSpace > 0 {
			spaceEvenly := availableSpace / (len(renderedItems) + 1)
			spaceStyle := lipgloss.NewStyle().Width(spaceEvenly)
			space := spaceStyle.Render("")

			var result strings.Builder
			result.WriteString(space)
			for i, item := range renderedItems {
				if i > 0 {
					result.WriteString(space)
					if f.config.Gap > 0 {
						gapStyle := lipgloss.NewStyle().Width(f.config.Gap)
						result.WriteString(gapStyle.Render(""))
					}
				}
				result.WriteString(item)
			}
			result.WriteString(space)
			return result.String()
		}
		return f.joinItemsWithGaps(renderedItems)

	default:
		return f.joinItemsWithGaps(renderedItems)
	}
}

// joinItemsWithGaps joins rendered items with configured gaps
func (f *FlexContainer) joinItemsWithGaps(renderedItems []string) string {
	if len(renderedItems) == 0 {
		return ""
	}

	if f.config.Gap == 0 {
		return strings.Join(renderedItems, "")
	}

	gapStyle := lipgloss.NewStyle().Width(f.config.Gap)
	gap := gapStyle.Render("")
	return strings.Join(renderedItems, gap)
}

// renderFlexColumn renders items in a column layout using CSS flexbox algorithms
func (f *FlexContainer) renderFlexColumn(items []FlexItem, width, height int) string {
	if len(items) == 0 {
		return ""
	}

	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 3 {
		height = 3
	}

	// Calculate flex-basis for each item
	itemHeights := make([]int, len(items))
	totalBasis := 0
	totalGrow := 0
	totalShrink := 0

	for i, item := range items {
		if item.FlexBasis > 0 {
			itemHeights[i] = item.FlexBasis
		} else {
			// Auto basis - use equal distribution
			itemHeights[i] = height / len(items)
		}
		totalBasis += itemHeights[i]
		totalGrow += item.FlexGrow
		if item.FlexShrink > 0 {
			totalShrink += item.FlexShrink
		} else {
			totalShrink += 1 // Default shrink is 1
		}
	}

	// Distribute remaining space using flex-grow or flex-shrink
	remainingSpace := height - totalBasis
	if remainingSpace > 0 && totalGrow > 0 {
		// Distribute extra space using flex-grow
		for i, item := range items {
			if item.FlexGrow > 0 {
				extraHeight := (remainingSpace * item.FlexGrow) / totalGrow
				itemHeights[i] += extraHeight
			}
		}
	} else if remainingSpace < 0 && totalShrink > 0 {
		// Shrink items using flex-shrink
		shrinkSpace := -remainingSpace
		for i, item := range items {
			shrinkFactor := item.FlexShrink
			if shrinkFactor == 0 {
				shrinkFactor = 1
			}
			shrinkAmount := (shrinkSpace * shrinkFactor) / totalShrink
			itemHeights[i] = max(1, itemHeights[i]-shrinkAmount)
		}
	}

	return f.renderColumnWithJustifyContent(items, itemHeights, width, height)
}

// renderColumnWithJustifyContent applies justify-content spacing to column items
func (f *FlexContainer) renderColumnWithJustifyContent(items []FlexItem, itemHeights []int, width, containerHeight int) string {
	// Calculate total item height
	totalItemHeight := 0
	for _, height := range itemHeights {
		totalItemHeight += height
	}

	// Calculate available space for justify-content
	availableSpace := containerHeight - totalItemHeight

	// Render each item with proper alignment
	renderedItems := make([]string, len(items))
	for i, item := range items {
		// Determine cross-axis alignment for this item
		align := f.alignItems
		if item.AlignSelf != nil {
			align = *item.AlignSelf
		}

		// Calculate cross-axis positioning
		var hAlign, vAlign lipgloss.Position
		switch align {
		case AlignStart:
			hAlign = lipgloss.Left
		case AlignCenter:
			hAlign = lipgloss.Center
		case AlignEnd:
			hAlign = lipgloss.Right
		case AlignStretch:
			hAlign = lipgloss.Left // Will use full width
		}
		vAlign = lipgloss.Top

		// Apply width and height constraints
		itemStyle := lipgloss.NewStyle().Height(itemHeights[i])
		if align == AlignStretch {
			itemStyle = itemStyle.Width(width)
		}

		// Render the item content with alignment
		styledContent := itemStyle.Render(item.Content)
		if align != AlignStretch {
			// Use lipgloss.Place for non-stretch alignment
			renderedItems[i] = lipgloss.Place(width, itemHeights[i], hAlign, vAlign, styledContent)
		} else {
			renderedItems[i] = styledContent
		}
	}

	// Apply justify-content spacing for column layout
	return f.applyJustifyContentSpacingColumn(renderedItems, availableSpace)
}

// applyJustifyContentSpacingColumn applies justify-content spacing to column items
func (f *FlexContainer) applyJustifyContentSpacingColumn(renderedItems []string, availableSpace int) string {
	if len(renderedItems) == 0 {
		return ""
	}

	var result strings.Builder

	switch f.justifyContent {
	case FlexStart:
		// Items at start, gaps between items
		return f.joinItemsWithGapsVertical(renderedItems)

	case FlexEnd:
		// Items at end, add space at beginning
		if availableSpace > 0 {
			for i := 0; i < availableSpace; i++ {
				result.WriteString("\n")
			}
		}
		result.WriteString(f.joinItemsWithGapsVertical(renderedItems))
		return result.String()

	case Center:
		// Items centered, equal space on both sides
		if availableSpace > 0 {
			sideSpace := availableSpace / 2
			for i := 0; i < sideSpace; i++ {
				result.WriteString("\n")
			}
		}
		result.WriteString(f.joinItemsWithGapsVertical(renderedItems))
		return result.String()

	case SpaceBetween:
		// Equal space between items, no space at edges
		if len(renderedItems) == 1 {
			return renderedItems[0]
		}
		if availableSpace > 0 {
			spaceBetween := availableSpace / (len(renderedItems) - 1)

			for i, item := range renderedItems {
				if i > 0 {
					// Add space between items
					for j := 0; j < spaceBetween; j++ {
						result.WriteString("\n")
					}
					// Add configured gap
					for j := 0; j < f.config.Gap; j++ {
						result.WriteString("\n")
					}
				}
				result.WriteString(item)
			}
			return result.String()
		}
		return f.joinItemsWithGapsVertical(renderedItems)

	case SpaceAround:
		// Equal space around each item
		if availableSpace > 0 {
			spaceAround := availableSpace / len(renderedItems)
			halfSpace := spaceAround / 2

			for i, item := range renderedItems {
				// Add half space before item
				for j := 0; j < halfSpace; j++ {
					result.WriteString("\n")
				}
				if i > 0 {
					// Add configured gap
					for j := 0; j < f.config.Gap; j++ {
						result.WriteString("\n")
					}
				}
				result.WriteString(item)
				// Add half space after item
				for j := 0; j < halfSpace; j++ {
					result.WriteString("\n")
				}
			}
			return result.String()
		}
		return f.joinItemsWithGapsVertical(renderedItems)

	case SpaceEvenly:
		// Equal space between and around items
		if availableSpace > 0 {
			spaceEvenly := availableSpace / (len(renderedItems) + 1)

			// Add space at beginning
			for i := 0; i < spaceEvenly; i++ {
				result.WriteString("\n")
			}

			for i, item := range renderedItems {
				if i > 0 {
					// Add space between items
					for j := 0; j < spaceEvenly; j++ {
						result.WriteString("\n")
					}
					// Add configured gap
					for j := 0; j < f.config.Gap; j++ {
						result.WriteString("\n")
					}
				}
				result.WriteString(item)
			}

			// Add space at end
			for i := 0; i < spaceEvenly; i++ {
				result.WriteString("\n")
			}
			return result.String()
		}
		return f.joinItemsWithGapsVertical(renderedItems)

	default:
		return f.joinItemsWithGapsVertical(renderedItems)
	}
}

// joinItemsWithGapsVertical joins rendered items vertically with configured gaps
func (f *FlexContainer) joinItemsWithGapsVertical(renderedItems []string) string {
	if len(renderedItems) == 0 {
		return ""
	}

	var result strings.Builder
	for i, item := range renderedItems {
		if i > 0 {
			// Add gap lines
			for j := 0; j < f.config.Gap; j++ {
				result.WriteString("\n")
			}
		}
		result.WriteString(item)
	}

	return result.String()
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

	widthRemainder := (availableWidth - (g.config.Gap * (g.columns - 1))) % g.columns
	heightRemainder := (availableHeight - (g.config.Gap * (g.rows - 1))) % g.rows

	// Distribute remainders
	colWidths := make([]int, g.columns)
	for j := 0; j < g.columns; j++ {
		colWidths[j] = cellWidth
	}
	for j := range widthRemainder {
		colWidths[j]++
	}

	rowHeights := make([]int, g.rows)
	for i := 0; i < g.rows; i++ {
		rowHeights[i] = cellHeight
	}
	for i := range heightRemainder {
		rowHeights[i]++
	}

	// Render each row
	renderedRows := make([]string, g.rows)
	for i := 0; i < g.rows; i++ {
		// Render cells in this row
		renderedCells := make([]string, g.columns)
		for j := 0; j < g.columns; j++ {
			cellStyle := lipgloss.NewStyle().
				Width(colWidths[j]).
				Height(rowHeights[i])

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
	if headerHeight < MinHeaderHeight {
		headerHeight = MinHeaderHeight
	}
	if headerHeight > MaxHeaderHeight {
		headerHeight = MaxHeaderHeight
	}
	if footerHeight < MinFooterHeight {
		footerHeight = MinFooterHeight
	}
	if footerHeight > MaxFooterHeight {
		footerHeight = MaxFooterHeight
	}

	// Calculate remaining height for main content
	mainHeight := height - headerHeight - footerHeight
	if mainHeight < MinMainContentHeight {
		// If not enough space, reduce header height
		headerHeight = max(MinHeaderHeight, height-footerHeight-MinMainContentHeight)
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
