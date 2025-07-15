package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Responsive breakpoints for adaptive layouts
const (
	BreakpointSmall  = 80  // Small terminal
	BreakpointMedium = 120 // Medium terminal
	BreakpointLarge  = 160 // Large terminal
)

// Claude Orange Theme - Comprehensive Color Palette
// Primary brand color: Claude orange (#FF6B35) with complementary colors
// Designed for accessibility and visual hierarchy across CLI and TUI

// Primary Colors - Claude Orange Brand Identity
var (
	// Core brand colors
	ClaudePrimary      = lipgloss.Color("#FF6B35") // Claude orange - primary brand color
	ClaudePrimaryLight = lipgloss.Color("#FF8A65") // Lighter orange for hover states
	ClaudePrimaryDark  = lipgloss.Color("#E55A2B") // Darker orange for pressed states

	// Complementary colors that work well with orange
	ClaudeSecondary      = lipgloss.Color("#6BB6FF") // Blue accent for links and highlights
	ClaudeSecondaryLight = lipgloss.Color("#8FC7FF") // Light blue for hover states
	ClaudeSecondaryDark  = lipgloss.Color("#4A90E2") // Darker blue for pressed states
)

// Neutral Colors - Foundation palette for backgrounds and text
var (
	// Background colors
	BackgroundPrimary      = lipgloss.Color("#2C3E50") // Primary dark background
	BackgroundSecondary    = lipgloss.Color("#34495E") // Alternative background for cards/surfaces
	BackgroundAccent       = lipgloss.Color("#58D68D") // More readable green accent
	BackgroundSurface      = lipgloss.Color("#4A5568") // Elevated surface color
	BackgroundSurfaceLight = lipgloss.Color("#718096") // Light surface for subtle elevation

	// Text colors with proper contrast ratios
	TextPrimary   = lipgloss.Color("#FFFFFF") // Primary text (high contrast)
	TextSecondary = lipgloss.Color("#D5DBDB") // Secondary text (medium contrast)
	TextMuted     = lipgloss.Color("#AEB6BF") // Muted text (low contrast)
	TextDim       = lipgloss.Color("#85929E") // Lighter gray for better readability
	TextDisabled  = lipgloss.Color("#718096") // Disabled text
)

// Semantic Colors - Status and feedback colors
var (
	// Status indicators
	SuccessColor      = lipgloss.Color("#2ECC71") // Green for success states
	SuccessColorLight = lipgloss.Color("#58D68D") // Light green for hover
	WarningColor      = lipgloss.Color("#F39C12") // Orange for warnings (complements brand)
	WarningColorLight = lipgloss.Color("#F6AD55") // Light orange for hover
	ErrorColor        = lipgloss.Color("#E74C3C") // Red for errors
	ErrorColorLight   = lipgloss.Color("#FC8181") // Light red for hover
	InfoColor         = lipgloss.Color("#5DADE2") // Blue for informational messages
	InfoColorLight    = lipgloss.Color("#63B3ED") // Light blue for hover
)

// Interactive State Colors - For buttons, links, and interactive elements
var (
	// Interactive states
	HoverColor      = ClaudePrimaryLight        // Hover state (lighter orange)
	FocusColor      = ClaudePrimary             // Focus state (primary orange)
	ActiveColor     = ClaudePrimaryDark         // Active/pressed state (darker orange)
	DisabledColor   = lipgloss.Color("#4A5568") // Disabled state (muted gray)
	SelectedColor   = ClaudePrimary             // Selected state (primary orange)
	SelectedBgColor = lipgloss.Color("#2A1810") // Selected background (dark orange tint)
)

// Navigation Colors - For menus, tabs, and navigation elements
var (
	NavBackground = BackgroundSecondary // Navigation background
	NavText       = TextPrimary         // Navigation text
	NavActive     = ClaudePrimary       // Active navigation item
	NavHover      = ClaudePrimaryLight  // Hovered navigation item
	NavBorder     = BackgroundSurface   // Navigation borders
)

// Content Colors - For different types of content
var (
	ContentPrimary   = TextPrimary               // Primary content text
	ContentSecondary = TextSecondary             // Secondary content text
	ContentMuted     = TextMuted                 // Muted content text
	ContentHighlight = ClaudePrimary             // Highlighted content
	ContentCode      = lipgloss.Color("#E53E3E") // Code snippets (red for contrast)
	ContentLink      = ClaudeSecondary           // Links and references
)

// Action Colors - For buttons and call-to-action elements
var (
	ActionPrimary   = ClaudePrimary     // Primary action buttons
	ActionSecondary = ClaudeSecondary   // Secondary action buttons
	ActionSuccess   = SuccessColor      // Success actions
	ActionWarning   = WarningColor      // Warning actions
	ActionDanger    = ErrorColor        // Dangerous actions
	ActionNeutral   = BackgroundSurface // Neutral actions
)

// Legacy color aliases for backward compatibility
var (
	// Session status colors (legacy aliases)
	StatusActiveColor    = SuccessColor // Green
	StatusInactiveColor  = WarningColor // Yellow
	StatusConnectedColor = InfoColor    // Blue
	StatusErrorColor     = ErrorColor   // Red
)

// Typography styles
var (
	// Title and headers
	TitleStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true).
			Underline(true)

	// Text styles
	PrimaryTextStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	SecondaryTextStyle = lipgloss.NewStyle().
				Foreground(TextSecondary)

	MutedTextStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(TextDim)

	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)
)

// Status and state styles
var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true)

	// Session status styles
	SessionStatusActiveStyle = lipgloss.NewStyle().
					Foreground(StatusActiveColor).
					Bold(true)

	SessionStatusInactiveStyle = lipgloss.NewStyle().
					Foreground(StatusInactiveColor).
					Bold(true)

	SessionStatusConnectedStyle = lipgloss.NewStyle().
					Foreground(StatusConnectedColor).
					Bold(true)

	SessionStatusErrorStyle = lipgloss.NewStyle().
				Foreground(StatusErrorColor).
				Bold(true)
)

// Interactive element styles - Enhanced with comprehensive state management
var (
	// Selection and focus states
	SelectedStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(SelectedColor).
			Bold(true).
			Padding(0, 1)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1)

	FocusedStyle = lipgloss.NewStyle().
			Foreground(FocusColor).
			Bold(true)

	BlurredStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	HoveredStyle = lipgloss.NewStyle().
			Foreground(HoverColor).
			Bold(true)

	ActiveStyle = lipgloss.NewStyle().
			Foreground(ActiveColor).
			Bold(true)

	DisabledStyle = lipgloss.NewStyle().
			Foreground(TextDisabled)

	// Button styles with comprehensive state management
	ButtonPrimaryStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ActionPrimary).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	ButtonSecondaryStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ActionSecondary).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	ButtonSuccessStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ActionSuccess).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	ButtonWarningStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ActionWarning).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	ButtonDangerStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ActionDanger).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	ButtonFocusedStyle = lipgloss.NewStyle().
				Foreground(ActionPrimary).
				Background(TextPrimary).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1).
				Border(lipgloss.NormalBorder()).
				BorderForeground(FocusColor)

	ButtonDisabledStyle = lipgloss.NewStyle().
				Foreground(TextDisabled).
				Background(DisabledColor).
				Padding(0, 2).
				Margin(0, 1)

	// Legacy button style for backward compatibility
	ButtonStyle = ButtonPrimaryStyle

	// Form elements with enhanced styling
	InputStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(BackgroundSecondary).
			Border(lipgloss.NormalBorder()).
			BorderForeground(TextMuted).
			Padding(0, 1).
			Margin(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(BackgroundSecondary).
				Border(lipgloss.NormalBorder()).
				BorderForeground(FocusColor).
				Padding(0, 1).
				Margin(0, 1)

	InputErrorStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(BackgroundSecondary).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ErrorColor).
			Padding(0, 1).
			Margin(0, 1)

	InputDisabledStyle = lipgloss.NewStyle().
				Foreground(TextDisabled).
				Background(DisabledColor).
				Border(lipgloss.NormalBorder()).
				BorderForeground(DisabledColor).
				Padding(0, 1).
				Margin(0, 1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true).
			Margin(0, 1)

	LabelRequiredStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Bold(true).
				Margin(0, 1)

	LabelErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			Margin(0, 1)

	// Link styles
	LinkStyle = lipgloss.NewStyle().
			Foreground(ContentLink).
			Underline(true)

	LinkHoverStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondaryLight).
			Underline(true).
			Bold(true)

	LinkVisitedStyle = lipgloss.NewStyle().
				Foreground(ClaudeSecondaryDark).
				Underline(true)
)

// Layout and container styles - Enhanced with comprehensive component styling
var (
	// Basic containers
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextMuted).
			Padding(1, 2)

	// Enhanced box styles with semantic coloring
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ClaudePrimary).
			Padding(1, 2).
			Margin(1, 0)

	InfoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(InfoColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(WarningColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ErrorColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	// Card styles for content organization
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BackgroundSurface).
			Padding(0, 1).
			Margin(1, 0)

	CardHeaderStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Padding(0, 0, 1, 0)

	CardContentStyle = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Padding(0)

	CardFooterStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(1, 0, 0, 0)

	// Modal and dialog styles
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ClaudePrimary).
			Background(BackgroundPrimary).
			Padding(2, 4).
			Margin(2, 4)

	ModalHeaderStyle = lipgloss.NewStyle().
				Foreground(ClaudePrimary).
				Bold(true).
				Align(lipgloss.Center).
				Padding(0, 0, 1, 0)

	ModalContentStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Padding(1, 0)

	ModalFooterStyle = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Align(lipgloss.Center).
				Padding(1, 0, 0, 0)

	// Panel styles for dashboard layout
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BackgroundSurface).
			Background(BackgroundSecondary).
			Padding(1).
			Margin(0, 1)

	PanelHeaderStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ClaudePrimary).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	PanelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ClaudePrimary).
				Background(BackgroundSecondary).
				Padding(1).
				Margin(0, 1)

	PanelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(FocusColor).
				Background(BackgroundSecondary).
				Padding(1).
				Margin(0, 1)

	// Sidebar and navigation styles
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(NavBorder).
			Background(NavBackground).
			Padding(1, 2).
			Width(20)

	SidebarHeaderStyle = lipgloss.NewStyle().
				Foreground(NavActive).
				Bold(true).
				Padding(0, 0, 1, 0).
				Align(lipgloss.Center)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(NavText).
				Padding(0, 1)

	SidebarItemActiveStyle = lipgloss.NewStyle().
				Foreground(NavActive).
				Background(SelectedBgColor).
				Bold(true).
				Padding(0, 1)

	SidebarItemHoverStyle = lipgloss.NewStyle().
				Foreground(NavHover).
				Bold(true).
				Padding(0, 1)
)

// Table styles - Enhanced with comprehensive table theming
var (
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ClaudePrimary).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1).
			Align(lipgloss.Left)

	TableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(SelectedColor).
				Bold(true)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	TableRowAlternateStyle = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Background(BackgroundSurface)

	TableRowHoverStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(HoverColor).
				Bold(true)

	TableBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(BackgroundSurface)

	// Status-specific table cell styles
	TableCellSuccessStyle = lipgloss.NewStyle().
				Foreground(SuccessColor).
				Bold(true).
				Padding(0, 1)

	TableCellWarningStyle = lipgloss.NewStyle().
				Foreground(WarningColor).
				Bold(true).
				Padding(0, 1)

	TableCellErrorStyle = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Bold(true).
				Padding(0, 1)

	TableCellInfoStyle = lipgloss.NewStyle().
				Foreground(InfoColor).
				Bold(true).
				Padding(0, 1)

	// Table column type styles
	TableCellIDStyle = lipgloss.NewStyle().
				Foreground(TextMuted).
				Padding(0, 1)

	TableCellNameStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Bold(true).
				Padding(0, 1)

	TableCellTimestampStyle = lipgloss.NewStyle().
				Foreground(TextDim).
				Padding(0, 1)
)

// Footer and navigation
var (
	FooterStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(1, 2).
			Align(lipgloss.Left)

	ArrowStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)
)

// Session-specific styles
var (
	SessionNameStyle = lipgloss.NewStyle().
				Foreground(ClaudePrimary).
				Bold(true)

	SessionIDStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	SessionDescriptionStyle = lipgloss.NewStyle().
				Foreground(TextSecondary)
)

// Progress and loading styles
var (
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)
)

// Text utility functions
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

// Responsive utilities - Enhanced with comprehensive adaptive behavior
func GetResponsiveWidth(terminalWidth int) (int, string) {
	switch {
	case terminalWidth < BreakpointSmall:
		return terminalWidth - 4, "small"
	case terminalWidth < BreakpointMedium:
		return terminalWidth - 8, "medium"
	default:
		return terminalWidth - 12, "large"
	}
}

// Adaptive width helper for different component types
func AdaptiveWidth(base lipgloss.Style, width int) lipgloss.Style {
	switch {
	case width < BreakpointSmall:
		return base.Width(width - 4)
	case width < BreakpointMedium:
		return base.Width(width - 8)
	default:
		return base.Width(width - 12)
	}
}

// Adaptive height helper for different component types
func AdaptiveHeight(base lipgloss.Style, height int) lipgloss.Style {
	switch {
	case height < 24:
		return base.Height(height - 2)
	case height < 40:
		return base.Height(height - 4)
	default:
		return base.Height(height - 6)
	}
}

// Responsive padding based on terminal size
func ResponsivePadding(width int) (horizontal, vertical int) {
	switch {
	case width < BreakpointSmall:
		return 1, 0
	case width < BreakpointMedium:
		return 2, 1
	default:
		return 3, 1
	}
}

// Responsive margin based on terminal size
func ResponsiveMargin(width int) (horizontal, vertical int) {
	switch {
	case width < BreakpointSmall:
		return 0, 0
	case width < BreakpointMedium:
		return 1, 0
	default:
		return 2, 1
	}
}

// Get appropriate table column widths based on terminal size
func GetTableColumnWidths(terminalWidth int, numColumns int) []int {
	availableWidth := terminalWidth - (numColumns * 3) // Account for borders and padding
	baseWidth := availableWidth / numColumns

	switch {
	case terminalWidth < BreakpointSmall:
		// Compact layout for small terminals
		return []int{baseWidth, baseWidth, baseWidth}[:numColumns]
	case terminalWidth < BreakpointMedium:
		// Balanced layout for medium terminals
		if numColumns >= 3 {
			return []int{baseWidth * 2, baseWidth, baseWidth}
		}
		return []int{baseWidth, baseWidth}[:numColumns]
	default:
		// Full layout for large terminals
		if numColumns >= 4 {
			return []int{baseWidth * 2, baseWidth, baseWidth, baseWidth}
		} else if numColumns >= 3 {
			return []int{baseWidth * 2, baseWidth, baseWidth}
		}
		return []int{baseWidth, baseWidth}[:numColumns]
	}
}

// Theme-aware color selection based on context
func GetContextualColor(context string, state string) lipgloss.Color {
	switch context {
	case "status":
		switch state {
		case "success", "active", "connected":
			return SuccessColor
		case "warning", "inactive":
			return WarningColor
		case "error", "failed":
			return ErrorColor
		case "info", "pending":
			return InfoColor
		default:
			return TextSecondary
		}
	case "action":
		switch state {
		case "primary":
			return ActionPrimary
		case "secondary":
			return ActionSecondary
		case "success":
			return ActionSuccess
		case "warning":
			return ActionWarning
		case "danger":
			return ActionDanger
		default:
			return ActionNeutral
		}
	case "interactive":
		switch state {
		case "hover":
			return HoverColor
		case "focus":
			return FocusColor
		case "active":
			return ActiveColor
		case "selected":
			return SelectedColor
		case "disabled":
			return DisabledColor
		default:
			return TextPrimary
		}
	default:
		return TextPrimary
	}
}

// Style rendering utility functions
func Title(text string) string {
	return TitleStyle.Render(text)
}

func Subtitle(text string) string {
	return SubtitleStyle.Render(text)
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func Success(text string) string {
	return SuccessStyle.Render("✓ " + text)
}

func Error(text string) string {
	return ErrorStyle.Render("✗ " + text)
}

func Warning(text string) string {
	return WarningStyle.Render("⚠ " + text)
}

func Info(text string) string {
	return InfoStyle.Render("ℹ " + text)
}

func Bold(text string) string {
	return BoldStyle.Render(text)
}

func Dim(text string) string {
	return MutedTextStyle.Render(text)
}

func Highlight(text string) string {
	return HighlightStyle.Render(text)
}

func Arrow() string {
	return ArrowStyle.Render("→")
}

func Prompt(text string) string {
	return PromptStyle.Render(text)
}

// Status formatting functions with icons
func StatusActive(text string) string {
	return SessionStatusActiveStyle.Render("● " + text)
}

func StatusInactive(text string) string {
	return SessionStatusInactiveStyle.Render("⏸ " + text)
}

func StatusConnected(text string) string {
	return SessionStatusConnectedStyle.Render("🔗 " + text)
}

func StatusError(text string) string {
	return SessionStatusErrorStyle.Render("✗ " + text)
}

func Spinner(text string) string {
	return SpinnerStyle.Render("⠋ " + text)
}

// Box rendering functions
func MainBox(content string) string {
	return MainBoxStyle.Render(content)
}

func InfoBox(content string) string {
	return InfoBoxStyle.Render(content)
}

func WarningBox(content string) string {
	return WarningBoxStyle.Render(content)
}

func ErrorBox(content string) string {
	return ErrorBoxStyle.Render(content)
}

// Panel rendering for dashboard layouts
func Panel(content string) string {
	return PanelStyle.Render(content)
}

func PanelWithHeader(header, content string) string {
	headerRendered := PanelHeaderStyle.Render(header)
	return PanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, headerRendered, content))
}

// Horizontal line with styling
func HorizontalLine(width int) string {
	return lipgloss.NewStyle().
		Foreground(TextMuted).
		Render(lipgloss.PlaceHorizontal(width, lipgloss.Center, "─"))
}

// Create a banner with title and subtitle
func Banner(title, subtitle string) string {
	titleRendered := TitleStyle.Render(title)
	subtitleRendered := SubtitleStyle.Render(subtitle)

	content := lipgloss.JoinVertical(lipgloss.Left, titleRendered, subtitleRendered)

	return MainBoxStyle.
		Width(60).
		Align(lipgloss.Center).
		Render(content)
}

// Helper for background color application
func WithBackground(text string, background lipgloss.Color) string {
	return lipgloss.NewStyle().Background(background).Render(text)
}

// Bubbles Component Styling Functions
// These functions provide styled configurations for Bubbles components
// using the existing Claude orange theme

// GetBubblesTableStyles returns a table.Styles configuration with Claude theme
func GetBubblesTableStyles() table.Styles {
	return table.Styles{
		Header: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(ClaudePrimary).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center),
		Cell: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1).
			Align(lipgloss.Left),
		Selected: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(SelectedColor).
			Bold(true),
	}
}

// ConfigureBubblesTextInputStyles applies styling to a textinput model
func ConfigureBubblesTextInputStyles(ti textinput.Model) textinput.Model {
	// Configure cursor style
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(ClaudePrimary)

	// Configure text styles
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextPrimary)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(TextMuted)

	// Configure prompt style
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ClaudePrimary)

	return ti
}

// ConfigureBubblesViewportStyles applies styling to a viewport model
func ConfigureBubblesViewportStyles(vp viewport.Model) viewport.Model {
	// Note: Viewport styling is limited in bubbles v0.21.0
	// Most styling is handled through the content
	return vp
}

// ConfigureBubblesHelpStyles applies styling to a help model
func ConfigureBubblesHelpStyles(h help.Model) help.Model {
	// Configure help styles
	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(TextSecondary)
	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(TextMuted)
	h.Styles.FullKey = lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true)
	h.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(TextSecondary)
	h.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(TextMuted)

	return h
}

// ConfigureBubblesTable applies Claude theme to a table model
func ConfigureBubblesTable(t table.Model) table.Model {
	t.SetStyles(GetBubblesTableStyles())
	return t
}

// ConfigureBubblesTextInput applies Claude theme to a textinput model
func ConfigureBubblesTextInput(ti textinput.Model) textinput.Model {
	return ConfigureBubblesTextInputStyles(ti)
}

// ConfigureBubblesViewport applies Claude theme to a viewport model
func ConfigureBubblesViewport(vp viewport.Model) viewport.Model {
	return ConfigureBubblesViewportStyles(vp)
}

// ConfigureBubblesHelp applies Claude theme to a help model
func ConfigureBubblesHelp(h help.Model) help.Model {
	return ConfigureBubblesHelpStyles(h)
}
