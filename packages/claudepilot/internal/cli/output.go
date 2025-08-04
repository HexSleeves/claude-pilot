package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
)

// OutputFormat represents the supported output formats
type OutputFormat string

const (
	OutputFormatHuman  OutputFormat = "human"
	OutputFormatTable  OutputFormat = "table"
	OutputFormatJSON   OutputFormat = "json"
	OutputFormatNDJSON OutputFormat = "ndjson"
	OutputFormatQuiet  OutputFormat = "quiet"
)

// String returns the string representation of OutputFormat
func (f OutputFormat) String() string {
	return string(f)
}

// IsValid checks if the output format is valid
func (f OutputFormat) IsValid() bool {
	switch f {
	case OutputFormatHuman, OutputFormatTable, OutputFormatJSON, OutputFormatNDJSON, OutputFormatQuiet:
		return true
	default:
		return false
	}
}

// SessionData represents a session for output formatting
type SessionData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Project     string    `json:"project,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AttachedTo  string    `json:"attachedTo,omitempty"`
	WindowCount int       `json:"windowCount,omitempty"`
	PaneCount   int       `json:"paneCount,omitempty"`
}

// ErrorData represents error information for output formatting
type ErrorData struct {
	Code      string            `json:"code"`
	Category  string            `json:"category"`
	Message   string            `json:"message"`
	Hint      string            `json:"hint,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	RequestID string            `json:"requestId,omitempty"`
}

// OperationResult represents the result of a command operation
type OperationResult struct {
	Success  bool              `json:"success"`
	Message  string            `json:"message"`
	Data     interface{}       `json:"data,omitempty"`
	Errors   []ErrorData       `json:"errors,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// OutputWriter interface defines methods for writing different types of output
type OutputWriter interface {
	WriteSessionList(sessions []SessionData, metadata map[string]string) error
	WriteSession(session SessionData, metadata map[string]string) error
	WriteError(err ErrorData) error
	WriteOperationResult(result OperationResult) error
	WriteString(s string) error
	SetWriter(w io.Writer)
	GetFormat() OutputFormat
}

// outputWriter implements OutputWriter interface
type outputWriter struct {
	format        OutputFormat
	writer        io.Writer
	isTTY         bool
	colorEnabled  bool
	schemaVersion string
}

// NewOutputWriter creates a new OutputWriter instance
func NewOutputWriter(format OutputFormat, isTTY bool, colorEnabled bool) OutputWriter {
	return &outputWriter{
		format:        format,
		writer:        os.Stdout,
		isTTY:         isTTY,
		colorEnabled:  colorEnabled,
		schemaVersion: "v1",
	}
}

// SetWriter sets the output writer
func (w *outputWriter) SetWriter(writer io.Writer) {
	w.writer = writer
}

// GetFormat returns the current output format
func (w *outputWriter) GetFormat() OutputFormat {
	return w.format
}

// WriteSessionList writes a list of sessions in the specified format
func (w *outputWriter) WriteSessionList(sessions []SessionData, metadata map[string]string) error {
	switch w.format {
	case OutputFormatJSON:
		return w.writeJSONSessionList(sessions, metadata)
	case OutputFormatNDJSON:
		return w.writeNDJSONSessionList(sessions, metadata)
	case OutputFormatTable:
		return w.writeTableSessionList(sessions)
	case OutputFormatQuiet:
		return w.writeQuietSessionList(sessions)
	case OutputFormatHuman:
		fallthrough
	default:
		return w.writeHumanSessionList(sessions)
	}
}

// WriteSession writes a single session in the specified format
func (w *outputWriter) WriteSession(session SessionData, metadata map[string]string) error {
	switch w.format {
	case OutputFormatJSON:
		return w.writeJSONSession(session, metadata)
	case OutputFormatNDJSON:
		return w.writeNDJSONSession(session, metadata)
	case OutputFormatTable:
		return w.writeTableSession(session)
	case OutputFormatQuiet:
		return w.writeQuietSession(session)
	case OutputFormatHuman:
		fallthrough
	default:
		return w.writeHumanSession(session)
	}
}

// WriteError writes an error in the specified format
func (w *outputWriter) WriteError(err ErrorData) error {
	switch w.format {
	case OutputFormatJSON, OutputFormatNDJSON:
		return w.writeJSONError(err)
	case OutputFormatQuiet:
		// In quiet mode, only output essential error information
		fmt.Fprintf(w.writer, "Error: %s\n", err.Message)
		return nil
	case OutputFormatTable, OutputFormatHuman:
		fallthrough
	default:
		return w.writeHumanError(err)
	}
}

// WriteOperationResult writes an operation result in the specified format
func (w *outputWriter) WriteOperationResult(result OperationResult) error {
	switch w.format {
	case OutputFormatJSON:
		return w.writeJSONOperationResult(result)
	case OutputFormatNDJSON:
		return w.writeNDJSONOperationResult(result)
	case OutputFormatQuiet:
		if !result.Success {
			fmt.Fprintf(w.writer, "Failed: %s\n", result.Message)
		}
		return nil
	case OutputFormatTable, OutputFormatHuman:
		fallthrough
	default:
		return w.writeHumanOperationResult(result)
	}
}

// WriteString writes a plain string (respects color settings)
func (w *outputWriter) WriteString(s string) error {
	if !w.colorEnabled {
		// Strip ANSI color codes if color is disabled
		s = stripANSI(s)
	}
	_, err := fmt.Fprint(w.writer, s)
	return err
}

// JSON output methods
func (w *outputWriter) writeJSONSessionList(sessions []SessionData, metadata map[string]string) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "SessionList",
		"metadata":      metadata,
		"items":         sessions,
		"count":         len(sessions),
	}
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func (w *outputWriter) writeJSONSession(session SessionData, metadata map[string]string) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "Session",
		"metadata":      metadata,
		"item":          session,
	}
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func (w *outputWriter) writeJSONError(err ErrorData) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "Error",
		"error":         err,
	}
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func (w *outputWriter) writeJSONOperationResult(result OperationResult) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "OperationResult",
		"result":        result,
	}
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// NDJSON output methods
func (w *outputWriter) writeNDJSONSessionList(sessions []SessionData, metadata map[string]string) error {
	encoder := json.NewEncoder(w.writer)

	// Write metadata header
	header := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "SessionListHeader",
		"metadata":      metadata,
		"count":         len(sessions),
	}
	if err := encoder.Encode(header); err != nil {
		return err
	}

	// Write each session
	for _, session := range sessions {
		item := map[string]interface{}{
			"schemaVersion": w.schemaVersion,
			"kind":          "Session",
			"item":          session,
		}
		if err := encoder.Encode(item); err != nil {
			return err
		}
	}
	return nil
}

func (w *outputWriter) writeNDJSONSession(session SessionData, metadata map[string]string) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "Session",
		"metadata":      metadata,
		"item":          session,
	}
	encoder := json.NewEncoder(w.writer)
	return encoder.Encode(response)
}

func (w *outputWriter) writeNDJSONOperationResult(result OperationResult) error {
	response := map[string]interface{}{
		"schemaVersion": w.schemaVersion,
		"kind":          "OperationResult",
		"result":        result,
	}
	encoder := json.NewEncoder(w.writer)
	return encoder.Encode(response)
}

// Table output methods
func (w *outputWriter) writeTableSessionList(sessions []SessionData) error {
	if len(sessions) == 0 {
		fmt.Fprintf(w.writer, "No sessions found.\n")
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(w.writer)

	// Configure table style
	if w.colorEnabled {
		t.SetStyle(table.StyleColoredBright)
	} else {
		t.SetStyle(table.StyleDefault)
	}

	t.AppendHeader(table.Row{"ID", "Name", "Status", "Project", "Created", "Windows", "Panes"})

	for _, session := range sessions {
		createdAt := session.CreatedAt.Format("2006-01-02 15:04")
		windowCount := ""
		if session.WindowCount > 0 {
			windowCount = fmt.Sprintf("%d", session.WindowCount)
		}
		paneCount := ""
		if session.PaneCount > 0 {
			paneCount = fmt.Sprintf("%d", session.PaneCount)
		}

		t.AppendRow(table.Row{
			session.ID,
			session.Name,
			session.Status,
			session.Project,
			createdAt,
			windowCount,
			paneCount,
		})
	}

	t.Render()
	return nil
}

func (w *outputWriter) writeTableSession(session SessionData) error {
	t := table.NewWriter()
	t.SetOutputMirror(w.writer)

	if w.colorEnabled {
		t.SetStyle(table.StyleColoredBright)
	} else {
		t.SetStyle(table.StyleDefault)
	}

	t.AppendHeader(table.Row{"Field", "Value"})
	t.AppendRow(table.Row{"ID", session.ID})
	t.AppendRow(table.Row{"Name", session.Name})
	if session.Description != "" {
		t.AppendRow(table.Row{"Description", session.Description})
	}
	if session.Project != "" {
		t.AppendRow(table.Row{"Project", session.Project})
	}
	t.AppendRow(table.Row{"Status", session.Status})
	t.AppendRow(table.Row{"Created", session.CreatedAt.Format("2006-01-02 15:04:05")})
	t.AppendRow(table.Row{"Updated", session.UpdatedAt.Format("2006-01-02 15:04:05")})
	if session.AttachedTo != "" {
		t.AppendRow(table.Row{"Attached To", session.AttachedTo})
	}
	if session.WindowCount > 0 {
		t.AppendRow(table.Row{"Windows", fmt.Sprintf("%d", session.WindowCount)})
	}
	if session.PaneCount > 0 {
		t.AppendRow(table.Row{"Panes", fmt.Sprintf("%d", session.PaneCount)})
	}

	t.Render()
	return nil
}

// Quiet output methods (minimal output)
func (w *outputWriter) writeQuietSessionList(sessions []SessionData) error {
	for _, session := range sessions {
		fmt.Fprintf(w.writer, "%s\n", session.ID)
	}
	return nil
}

func (w *outputWriter) writeQuietSession(session SessionData) error {
	fmt.Fprintf(w.writer, "%s\n", session.ID)
	return nil
}

// Human-readable output methods
func (w *outputWriter) writeHumanSessionList(sessions []SessionData) error {
	if len(sessions) == 0 {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if !w.colorEnabled {
			style = lipgloss.NewStyle()
		}
		fmt.Fprintf(w.writer, "%s\n", style.Render("No sessions found."))
		return nil
	}

	// Use existing table formatting but with human-friendly styling
	return w.writeTableSessionList(sessions)
}

func (w *outputWriter) writeHumanSession(session SessionData) error {
	// Create styled output for human readability
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	fieldStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("33"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	if !w.colorEnabled {
		headerStyle = lipgloss.NewStyle().Bold(true)
		fieldStyle = lipgloss.NewStyle().Bold(true)
		valueStyle = lipgloss.NewStyle()
	}

	fmt.Fprintf(w.writer, "%s\n", headerStyle.Render("Session Details"))
	fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("ID:"), valueStyle.Render(session.ID))
	fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Name:"), valueStyle.Render(session.Name))

	if session.Description != "" {
		fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Description:"), valueStyle.Render(session.Description))
	}
	if session.Project != "" {
		fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Project:"), valueStyle.Render(session.Project))
	}

	fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Status:"), valueStyle.Render(session.Status))
	fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Created:"), valueStyle.Render(session.CreatedAt.Format("2006-01-02 15:04:05")))
	fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Updated:"), valueStyle.Render(session.UpdatedAt.Format("2006-01-02 15:04:05")))

	if session.AttachedTo != "" {
		fmt.Fprintf(w.writer, "%s %s\n", fieldStyle.Render("Attached To:"), valueStyle.Render(session.AttachedTo))
	}
	if session.WindowCount > 0 {
		fmt.Fprintf(w.writer, "%s %d\n", fieldStyle.Render("Windows:"), session.WindowCount)
	}
	if session.PaneCount > 0 {
		fmt.Fprintf(w.writer, "%s %d\n", fieldStyle.Render("Panes:"), session.PaneCount)
	}

	return nil
}

func (w *outputWriter) writeHumanError(err ErrorData) error {
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	categoryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	if !w.colorEnabled {
		errorStyle = lipgloss.NewStyle().Bold(true)
		categoryStyle = lipgloss.NewStyle()
		hintStyle = lipgloss.NewStyle()
	}

	fmt.Fprintf(w.writer, "%s %s\n", errorStyle.Render("Error:"), err.Message)
	if err.Category != "" {
		fmt.Fprintf(w.writer, "%s %s\n", categoryStyle.Render("Category:"), err.Category)
	}
	if err.Code != "" {
		fmt.Fprintf(w.writer, "%s %s\n", categoryStyle.Render("Code:"), err.Code)
	}
	if err.Hint != "" {
		fmt.Fprintf(w.writer, "%s %s\n", hintStyle.Render("Hint:"), err.Hint)
	}

	return nil
}

func (w *outputWriter) writeHumanOperationResult(result OperationResult) error {
	if result.Success {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
		if !w.colorEnabled {
			successStyle = lipgloss.NewStyle()
		}
		fmt.Fprintf(w.writer, "%s %s\n", successStyle.Render("Success:"), result.Message)
	} else {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		if !w.colorEnabled {
			errorStyle = lipgloss.NewStyle()
		}
		fmt.Fprintf(w.writer, "%s %s\n", errorStyle.Render("Failed:"), result.Message)

		if len(result.Errors) > 0 {
			for _, err := range result.Errors {
				w.writeHumanError(err)
			}
		}
	}

	return nil
}

// IsTTY detects if the output is a terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// stripANSI removes ANSI escape sequences from a string
func stripANSI(s string) string {
	// Simple ANSI escape sequence removal
	// This is a basic implementation - you might want to use a library like github.com/acarl005/stripansi
	result := strings.Builder{}
	inEscape := false

	for i, r := range s {
		if r == '\033' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}
