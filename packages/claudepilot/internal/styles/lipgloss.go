package styles

import (
	"claude-pilot/shared/styles"

	"github.com/charmbracelet/lipgloss"
)

// Re-export colors from shared styles for backward compatibility - Enhanced with new theme elements
var (
	// Primary colors - Enhanced with light/dark variants
	ClaudePrimary        = styles.ClaudePrimary
	ClaudePrimaryLight   = styles.ClaudePrimaryLight
	ClaudePrimaryDark    = styles.ClaudePrimaryDark
	ClaudeSecondary      = styles.ClaudeSecondary
	ClaudeSecondaryLight = styles.ClaudeSecondaryLight
	ClaudeSecondaryDark  = styles.ClaudeSecondaryDark

	// Status colors - Enhanced with light variants
	SuccessColor      = styles.SuccessColor
	SuccessColorLight = styles.SuccessColorLight
	ErrorColor        = styles.ErrorColor
	ErrorColorLight   = styles.ErrorColorLight
	WarningColor      = styles.WarningColor
	WarningColorLight = styles.WarningColorLight
	InfoColor         = styles.InfoColor
	InfoColorLight    = styles.InfoColorLight

	// Neutral colors - Enhanced with additional variants
	TextPrimary   = styles.TextPrimary
	TextSecondary = styles.TextSecondary
	TextMuted     = styles.TextMuted
	TextDim       = styles.TextDim
	TextDisabled  = styles.TextDisabled

	// Background colors - Enhanced with surface variants
	BackgroundPrimary      = styles.BackgroundPrimary
	BackgroundSecondary    = styles.BackgroundSecondary
	BackgroundAccent       = styles.BackgroundAccent
	BackgroundSurface      = styles.BackgroundSurface
	BackgroundSurfaceLight = styles.BackgroundSurfaceLight

	// Interactive state colors
	HoverColor      = styles.HoverColor
	FocusColor      = styles.FocusColor
	ActiveColor     = styles.ActiveColor
	SelectedColor   = styles.SelectedColor
	SelectedBgColor = styles.SelectedBgColor
	DisabledColor   = styles.DisabledColor

	// Content colors
	ContentPrimary   = styles.ContentPrimary
	ContentSecondary = styles.ContentSecondary
	ContentMuted     = styles.ContentMuted
	ContentHighlight = styles.ContentHighlight
	ContentCode      = styles.ContentCode
	ContentLink      = styles.ContentLink

	// Action colors
	ActionPrimary   = styles.ActionPrimary
	ActionSecondary = styles.ActionSecondary
	ActionSuccess   = styles.ActionSuccess
	ActionWarning   = styles.ActionWarning
	ActionDanger    = styles.ActionDanger
	ActionNeutral   = styles.ActionNeutral
)

// Re-export styles from shared package for backward compatibility - Enhanced with new styles
var (
	// Title styles
	TitleStyle    = styles.TitleStyle
	SubtitleStyle = styles.SubtitleStyle
	HeaderStyle   = styles.HeaderStyle

	// Message styles
	SuccessStyle = styles.SuccessStyle
	ErrorStyle   = styles.ErrorStyle
	WarningStyle = styles.WarningStyle
	InfoStyle    = styles.InfoStyle

	// Text styles - Enhanced with additional variants
	BoldStyle          = styles.BoldStyle
	DimStyle           = styles.DimTextStyle
	HighlightStyle     = styles.HighlightStyle
	PrimaryTextStyle   = styles.PrimaryTextStyle
	SecondaryTextStyle = styles.SecondaryTextStyle
	MutedTextStyle     = styles.MutedTextStyle

	// Interactive elements - Enhanced with state styles
	ArrowStyle    = styles.ArrowStyle
	PromptStyle   = styles.PromptStyle
	SelectedStyle = styles.SelectedStyle
	FocusedStyle  = styles.FocusedStyle
	HoveredStyle  = styles.HoveredStyle
	ActiveStyle   = styles.ActiveStyle
	DisabledStyle = styles.DisabledStyle

	// Button styles - Enhanced with comprehensive variants
	ButtonPrimaryStyle   = styles.ButtonPrimaryStyle
	ButtonSecondaryStyle = styles.ButtonSecondaryStyle
	ButtonSuccessStyle   = styles.ButtonSuccessStyle
	ButtonWarningStyle   = styles.ButtonWarningStyle
	ButtonDangerStyle    = styles.ButtonDangerStyle
	ButtonFocusedStyle   = styles.ButtonFocusedStyle
	ButtonDisabledStyle  = styles.ButtonDisabledStyle
	ButtonStyle          = styles.ButtonStyle // Legacy compatibility

	// Form styles - Enhanced with state variants
	InputStyle         = styles.InputStyle
	InputFocusedStyle  = styles.InputFocusedStyle
	InputErrorStyle    = styles.InputErrorStyle
	InputDisabledStyle = styles.InputDisabledStyle
	LabelStyle         = styles.LabelStyle
	LabelRequiredStyle = styles.LabelRequiredStyle
	LabelErrorStyle    = styles.LabelErrorStyle

	// Link styles
	LinkStyle        = styles.LinkStyle
	LinkHoverStyle   = styles.LinkHoverStyle
	LinkVisitedStyle = styles.LinkVisitedStyle

	// Table styles - Enhanced with comprehensive cell types
	TableHeaderStyle        = styles.TableHeaderStyle
	TableCellStyle          = styles.TableCellStyle
	TableSelectedRowStyle   = styles.TableSelectedRowStyle
	TableRowStyle           = styles.TableRowStyle
	TableRowAlternateStyle  = styles.TableRowAlternateStyle
	TableRowHoverStyle      = styles.TableRowHoverStyle
	TableBorderStyle        = styles.TableBorderStyle
	TableCellSuccessStyle   = styles.TableCellSuccessStyle
	TableCellWarningStyle   = styles.TableCellWarningStyle
	TableCellErrorStyle     = styles.TableCellErrorStyle
	TableCellInfoStyle      = styles.TableCellInfoStyle
	TableCellIDStyle        = styles.TableCellIDStyle
	TableCellNameStyle      = styles.TableCellNameStyle
	TableCellTimestampStyle = styles.TableCellTimestampStyle

	// Status indicators
	StatusActiveStyle    = styles.SessionStatusActiveStyle
	StatusInactiveStyle  = styles.SessionStatusInactiveStyle
	StatusConnectedStyle = styles.SessionStatusConnectedStyle
	StatusErrorStyle     = styles.SessionStatusErrorStyle

	// Progress and loading
	SpinnerStyle  = styles.SpinnerStyle
	ProgressStyle = styles.ProgressStyle
)

// Re-export box styles from shared package - Enhanced with comprehensive container styles
var (
	// Basic box styles
	MainBoxStyle    = styles.MainBoxStyle
	InfoBoxStyle    = styles.InfoBoxStyle
	SuccessBoxStyle = styles.SuccessBoxStyle
	WarningBoxStyle = styles.WarningBoxStyle
	ErrorBoxStyle   = styles.ErrorBoxStyle

	// Card styles
	CardStyle        = styles.CardStyle
	CardHeaderStyle  = styles.CardHeaderStyle
	CardContentStyle = styles.CardContentStyle
	CardFooterStyle  = styles.CardFooterStyle

	// Modal and dialog styles
	ModalStyle        = styles.ModalStyle
	ModalHeaderStyle  = styles.ModalHeaderStyle
	ModalContentStyle = styles.ModalContentStyle
	ModalFooterStyle  = styles.ModalFooterStyle

	// Panel styles
	PanelStyle        = styles.PanelStyle
	PanelHeaderStyle  = styles.PanelHeaderStyle
	PanelActiveStyle  = styles.PanelActiveStyle
	PanelFocusedStyle = styles.PanelFocusedStyle

	// Sidebar and navigation styles
	SidebarStyle           = styles.SidebarStyle
	SidebarHeaderStyle     = styles.SidebarHeaderStyle
	SidebarItemStyle       = styles.SidebarItemStyle
	SidebarItemActiveStyle = styles.SidebarItemActiveStyle
	SidebarItemHoverStyle  = styles.SidebarItemHoverStyle

	// Container styles
	ContainerStyle = styles.ContainerStyle
	BorderStyle    = styles.BorderStyle
)

func Title(text string) string {
	return styles.Title(text)
}

func Subtitle(text string) string {
	return styles.Subtitle(text)
}

func Header(text string) string {
	return styles.Header(text)
}

func Success(text string) string {
	return styles.Success(text)
}

func Error(text string) string {
	return styles.Error(text)
}

func Warning(text string) string {
	return styles.Warning(text)
}

func Info(text string) string {
	return styles.Info(text)
}

func Bold(text string) string {
	return styles.Bold(text)
}

func BoldWithBackground(text string, background lipgloss.Color) string {
	return styles.WithBackground(styles.Bold(text), background)
}

func Dim(text string) string {
	return styles.Dim(text)
}

func Highlight(text string) string {
	return styles.Highlight(text)
}

func Arrow() string {
	return styles.Arrow()
}

func Prompt(text string) string {
	return styles.Prompt(text)
}

// Status formatting functions
func StatusActive(text string) string {
	return styles.StatusActive(text)
}

func StatusInactive(text string) string {
	return styles.StatusInactive(text)
}

func StatusConnected(text string) string {
	return styles.StatusConnected(text)
}

func StatusError(text string) string {
	return styles.StatusError(text)
}

func Spinner(text string) string {
	return styles.Spinner(text)
}

// Box rendering functions
func MainBox(content string) string {
	return styles.MainBox(content)
}

func InfoBox(content string) string {
	return styles.InfoBox(content)
}

func WarningBox(content string) string {
	return styles.WarningBox(content)
}

func ErrorBox(content string) string {
	return styles.ErrorBox(content)
}

// Horizontal line with styling
func HorizontalLine(width int) string {
	return styles.HorizontalLine(width)
}
