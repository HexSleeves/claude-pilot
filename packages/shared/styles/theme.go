package styles

import (
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
)

// ButtonStyleOptions defines optional style modifiers for button creation
type ButtonStyleOptions struct {
	Bold              bool
	Border            *lipgloss.Border
	BorderForeground  *lipgloss.Color
	PaddingVertical   int
	PaddingHorizontal int
	MarginVertical    int
	MarginHorizontal  int
}

// createButtonStyle creates a button style with the given foreground, background, and options
func createButtonStyle(foreground, background lipgloss.Color, options ButtonStyleOptions) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(foreground).
		Background(background)

	// Apply optional modifiers
	if options.Bold {
		style = style.Bold(true)
	}

	// Set padding (default to 0, 2 if not specified)
	if options.PaddingVertical == 0 && options.PaddingHorizontal == 0 {
		style = style.Padding(0, 2)
	} else {
		style = style.Padding(options.PaddingVertical, options.PaddingHorizontal)
	}

	// Set margin (default to 0, 1 if not specified)
	if options.MarginVertical == 0 && options.MarginHorizontal == 0 {
		style = style.Margin(0, 1)
	} else {
		style = style.Margin(options.MarginVertical, options.MarginHorizontal)
	}

	// Apply border if specified
	if options.Border != nil {
		style = style.Border(*options.Border)
		if options.BorderForeground != nil {
			style = style.BorderForeground(*options.BorderForeground)
		}
	}

	return style
}

// Input styles
func createInputStyle(foreground, background, borderForeground lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(foreground).
		Background(background).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderForeground).
		Padding(0, 1).
		Margin(0, 1)
}

func createLabelStyle(foreground, background lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(foreground).
		Background(background).
		Bold(true).
		Margin(0, 1)
}

func createLinkStyle(foreground, background lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(foreground).
		Background(background).
		Underline(true)
}

var (
	// Helper variable for border pointer
	normalBorder = lipgloss.NormalBorder()

	// Button styles with comprehensive state management
	ButtonPrimaryStyle   = createButtonStyle(TextPrimary, ActionPrimary, ButtonStyleOptions{Bold: true})
	ButtonSecondaryStyle = createButtonStyle(TextPrimary, ActionSecondary, ButtonStyleOptions{Bold: true})
	ButtonSuccessStyle   = createButtonStyle(TextPrimary, ActionSuccess, ButtonStyleOptions{Bold: true})
	ButtonWarningStyle   = createButtonStyle(TextPrimary, ActionWarning, ButtonStyleOptions{Bold: true})
	ButtonDangerStyle    = createButtonStyle(TextPrimary, ActionDanger, ButtonStyleOptions{Bold: true})
	ButtonFocusedStyle   = createButtonStyle(ActionPrimary, TextPrimary, ButtonStyleOptions{
		Bold:             true,
		Border:           &normalBorder,
		BorderForeground: &FocusColor,
	})
	ButtonDisabledStyle = createButtonStyle(TextDisabled, DisabledColor, ButtonStyleOptions{})

	// Form elements with enhanced styling
	InputStyle         = createInputStyle(TextPrimary, BackgroundSecondary, TextMuted)
	InputFocusedStyle  = createInputStyle(TextPrimary, BackgroundSecondary, FocusColor)
	InputErrorStyle    = createInputStyle(TextPrimary, BackgroundSecondary, ErrorColor)
	InputDisabledStyle = createInputStyle(TextDisabled, DisabledColor, DisabledColor)

	LabelStyle         = createLabelStyle(TextPrimary, BackgroundSecondary)
	LabelRequiredStyle = createLabelStyle(TextPrimary, BackgroundSecondary)
	LabelErrorStyle    = createLabelStyle(ErrorColor, BackgroundSecondary)

	// Link styles
	LinkStyle        = createLinkStyle(ContentLink, BackgroundSecondary)
	LinkHoverStyle   = createLinkStyle(ClaudeSecondaryLight, BackgroundSecondary)
	LinkVisitedStyle = createLinkStyle(ClaudeSecondaryDark, BackgroundSecondary)
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
