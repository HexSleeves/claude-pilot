package styles

import "github.com/charmbracelet/lipgloss"

// Style rendering utility functions

func Title(text string) string {
	// Use shared theme title style instead of fatih/color
	return TitleStyle.Render(text)
}

func Subtitle(text string) string {
	// Use shared theme subtitle style
	return SubtitleStyle.Render(text)
}

func SuccessMsg(text string) string {
	// Use shared theme success style with icon
	return Success(text)
}

func ErrorMsg(text string) string {
	// Use shared theme error style with icon
	return Error(text)
}

func WarningMsg(text string) string {
	// Use shared theme warning style with icon
	return Warning(text)
}

func InfoMsg(text string) string {
	// Use shared theme info style with icon
	return Info(text)
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func Success(text string) string {
	return SuccessStyle.Render("✔ " + text)
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
