package ui

import (
	"fmt"
	"io"
	"os"

	"claude-pilot/internal/interfaces"
)

// RenderMode defines the UI rendering mode
type RenderMode string

const (
	ModeCLI RenderMode = "cli"
	ModeTUI RenderMode = "tui"
)

// RenderContext provides context for rendering elements
type RenderContext struct {
	Width  int
	Height int
	Mode   RenderMode
	Theme  *Theme
	Writer io.Writer
}

// Element represents a UI element that can be rendered
type Element interface {
	Render(ctx *RenderContext) string
}

// Renderer interface for different UI modes
type Renderer interface {
	Title(text string) Element
	Subtitle(text string) Element
	ErrorMsg(text string) Element
	SuccessMsg(text string) Element
	WarningMsg(text string) Element
	InfoMsg(text string) Element
	Prompt(text string) Element
	SessionTable(sessions []*interfaces.Session) Element
	SessionDetail(session *interfaces.Session) Element
}

// Theme defines color and styling configuration
type Theme struct {
	Name      string
	Primary   ColorFunc
	Success   ColorFunc
	Error     ColorFunc
	Warning   ColorFunc
	Info      ColorFunc
	Muted     ColorFunc
	Highlight ColorFunc
	Connected ColorFunc
	Active    ColorFunc
	Inactive  ColorFunc
}

// ColorFunc represents a color application function
type ColorFunc func(string) string

// GetDefaultTheme returns the default Claude theme
func GetDefaultTheme() *Theme {
	return &Theme{
		Name:      "claude-default",
		Highlight: Highlight,
		Primary:   func(s string) string { return ClaudePrimary.Sprint(s) },
		Success:   func(s string) string { return Success.Sprint(s) },
		Error:     func(s string) string { return Error.Sprint(s) },
		Warning:   func(s string) string { return Warning.Sprint(s) },
		Info:      func(s string) string { return Info.Sprint(s) },
		Muted:     func(s string) string { return TextMuted.Sprint(s) },
		Connected: func(s string) string { return StatusConnected.Sprint(s) },
		Active:    func(s string) string { return StatusActive.Sprint(s) },
		Inactive:  func(s string) string { return StatusInactive.Sprint(s) },
	}
}

// NewRenderContext creates a new render context
func NewRenderContext(mode RenderMode, writer io.Writer) *RenderContext {
	if writer == nil {
		writer = os.Stdout
	}

	return &RenderContext{
		Width:  80, // Default width
		Height: 24, // Default height
		Mode:   mode,
		Theme:  GetDefaultTheme(),
		Writer: writer,
	}
}

// CLI-specific elements
type textElement struct {
	text  string
	style func(string) string
}

func (e *textElement) Render(ctx *RenderContext) string {
	if e.style != nil {
		return e.style(e.text)
	}
	return e.text
}

type sessionTableElement struct {
	sessions []*interfaces.Session
}

func (e *sessionTableElement) Render(ctx *RenderContext) string {
	// For now, delegate to existing table implementation
	// In the future, this could render differently based on context.Mode
	switch ctx.Mode {
	case ModeTUI:
		// Future TUI implementation
		return "TUI table rendering not yet implemented"
	default:
		// Use existing CLI table implementation - temporarily disabled for compatibility
		return "CLI table rendering with new interfaces not yet implemented"
	}
}

type sessionDetailElement struct {
	session *interfaces.Session
}

func (e *sessionDetailElement) Render(ctx *RenderContext) string {
	// For now, delegate to existing detail implementation
	switch ctx.Mode {
	case ModeTUI:
		// Future TUI implementation
		return "TUI detail rendering not yet implemented"
	default:
		// Use existing CLI detail implementation - temporarily disabled for compatibility
		return "CLI detail rendering with new interfaces not yet implemented"
	}
}

// CLIRenderer implements the Renderer interface for CLI mode
type CLIRenderer struct {
	theme *Theme
}

// NewCLIRenderer creates a new CLI renderer
func NewCLIRenderer(theme *Theme) *CLIRenderer {
	if theme == nil {
		theme = GetDefaultTheme()
	}
	return &CLIRenderer{theme: theme}
}

func (r *CLIRenderer) Title(text string) Element {
	return &textElement{text: text, style: r.theme.Primary}
}

func (r *CLIRenderer) Subtitle(text string) Element {
	return &textElement{text: text, style: r.theme.Highlight}
}

func (r *CLIRenderer) ErrorMsg(text string) Element {
	return &textElement{text: fmt.Sprintf("✗ %s", text), style: r.theme.Error}
}

func (r *CLIRenderer) SuccessMsg(text string) Element {
	return &textElement{text: fmt.Sprintf("✓ %s", text), style: r.theme.Success}
}

func (r *CLIRenderer) WarningMsg(text string) Element {
	return &textElement{text: fmt.Sprintf("⚠ %s", text), style: r.theme.Warning}
}

func (r *CLIRenderer) InfoMsg(text string) Element {
	return &textElement{text: fmt.Sprintf("ℹ %s", text), style: r.theme.Info}
}

func (r *CLIRenderer) Prompt(text string) Element {
	return &textElement{text: text, style: r.theme.Highlight}
}

func (r *CLIRenderer) SessionTable(sessions []*interfaces.Session) Element {
	return &sessionTableElement{sessions: sessions}
}

func (r *CLIRenderer) SessionDetail(session *interfaces.Session) Element {
	return &sessionDetailElement{session: session}
}

// Output utility functions
func RenderToString(element Element, ctx *RenderContext) string {
	return element.Render(ctx)
}

func RenderAndPrint(element Element, ctx *RenderContext) {
	output := element.Render(ctx)
	fmt.Fprint(ctx.Writer, output)
}

// Legacy compatibility functions - these maintain backward compatibility
// while allowing gradual migration to the new renderer system

func RenderTitle(text string) string {
	renderer := NewCLIRenderer(nil)
	ctx := NewRenderContext(ModeCLI, nil)
	return renderer.Title(text).Render(ctx)
}

func RenderError(text string) string {
	renderer := NewCLIRenderer(nil)
	ctx := NewRenderContext(ModeCLI, nil)
	return renderer.ErrorMsg(text).Render(ctx)
}

func RenderSuccess(text string) string {
	renderer := NewCLIRenderer(nil)
	ctx := NewRenderContext(ModeCLI, nil)
	return renderer.SuccessMsg(text).Render(ctx)
}
