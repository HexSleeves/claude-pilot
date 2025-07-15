package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// TableConfig holds configuration for table rendering
type TableConfig struct {
	Width       int
	ShowHeaders bool
	Interactive bool
	SelectedRow int
	MaxRows     int
}

// TableData represents structured data for table rendering
type TableData struct {
	Headers []string
	Rows    [][]string
}

// SessionData represents session information for table display
type SessionData struct {
	ID          string
	Name        string
	Status      string
	Backend     string
	Created     time.Time
	LastActive  time.Time
	Messages    int
	ProjectPath string
}

// Table provides a unified table component for both CLI and TUI
type Table struct {
	config TableConfig
	data   TableData
}

// NewTable creates a new table instance
func NewTable(config TableConfig) *Table {
	return &Table{
		config: config,
	}
}

// SetData sets the table data
func (t *Table) SetData(data TableData) {
	t.data = data
}

// SetSessionData converts session data to table format
func (t *Table) SetSessionData(sessions []SessionData) {
	headers := []string{"ID", "Name", "Status", "Backend", "Created", "Last Active", "Messages", "Project"}
	rows := make([][]string, len(sessions))

	for i, session := range sessions {
		rows[i] = []string{
			styles.TruncateText(session.ID, 11),
			styles.TruncateText(session.Name, 19),
			session.Status,
			session.Backend,
			formatTime(session.Created),
			formatTimeAgo(session.LastActive),
			fmt.Sprintf("%d", session.Messages),
			formatProjectPath(session.ProjectPath, 29),
		}
	}

	t.data = TableData{
		Headers: headers,
		Rows:    rows,
	}
}

// RenderCLI renders the table for CLI output (static)
func (t *Table) RenderCLI() string {
	if len(t.data.Rows) == 0 {
		return styles.Dim("No data to display.")
	}

	// Create lipgloss table
	tbl := table.New()

	// Set border and styling
	tbl.Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.ClaudePrimary)).
		StyleFunc(t.createCLIStyleFunc())

	// Set headers if configured
	if t.config.ShowHeaders && len(t.data.Headers) > 0 {
		headers := make([]string, len(t.data.Headers))
		for i, header := range t.data.Headers {
			headers[i] = styles.Bold(header)
		}
		tbl.Headers(headers...)
	}

	// Add rows with limits
	maxRows := len(t.data.Rows)
	if t.config.MaxRows > 0 && t.config.MaxRows < maxRows {
		maxRows = t.config.MaxRows
	}

	for i := 0; i < maxRows; i++ {
		row := t.data.Rows[i]
		styledRow := make([]string, len(row))
		for j, cell := range row {
			styledRow[j] = t.styleCellForCLI(cell, j)
		}
		tbl.Row(styledRow...)
	}

	return tbl.String()
}

// RenderTUI renders the table for TUI output (can be interactive)
func (t *Table) RenderTUI() string {
	if len(t.data.Rows) == 0 {
		return styles.Dim("No data to display.")
	}

	var builder strings.Builder

	// Render headers if configured
	if t.config.ShowHeaders && len(t.data.Headers) > 0 {
		headerRow := make([]string, len(t.data.Headers))
		for i, header := range t.data.Headers {
			headerRow[i] = styles.TableHeaderStyle.Render(header)
		}
		builder.WriteString(strings.Join(headerRow, " ") + "\n")
		builder.WriteString(styles.HorizontalLine(t.config.Width) + "\n")
	}

	// Render data rows
	maxRows := len(t.data.Rows)
	if t.config.MaxRows > 0 && t.config.MaxRows < maxRows {
		maxRows = t.config.MaxRows
	}

	for i := 0; i < maxRows; i++ {
		row := t.data.Rows[i]
		isSelected := t.config.Interactive && i == t.config.SelectedRow

		styledRow := make([]string, len(row))
		for j, cell := range row {
			styledRow[j] = t.styleCellForTUI(cell, j, isSelected)
		}

		rowString := strings.Join(styledRow, " ")
		if isSelected {
			rowString = styles.TableSelectedRowStyle.Render(rowString)
		}
		builder.WriteString(rowString + "\n")
	}

	return builder.String()
}

// SetSelectedRow updates the selected row for interactive mode
func (t *Table) SetSelectedRow(row int) {
	if row >= 0 && row < len(t.data.Rows) {
		t.config.SelectedRow = row
	}
}

// SetWidth sets the table width
func (t *Table) SetWidth(width int) {
	t.config.Width = width
}

// SetMaxRows sets the maximum number of rows to display
func (t *Table) SetMaxRows(maxRows int) {
	t.config.MaxRows = maxRows
}

// GetSelectedRow returns the currently selected row index
func (t *Table) GetSelectedRow() int {
	return t.config.SelectedRow
}

// GetRowCount returns the number of rows in the table
func (t *Table) GetRowCount() int {
	return len(t.data.Rows)
}

// GetSelectedData returns the data for the currently selected row
func (t *Table) GetSelectedData() []string {
	if t.config.SelectedRow >= 0 && t.config.SelectedRow < len(t.data.Rows) {
		return t.data.Rows[t.config.SelectedRow]
	}
	return nil
}

// createCLIStyleFunc creates styling for CLI table rendering
func (t *Table) createCLIStyleFunc() table.StyleFunc {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.TextPrimary).
		Background(styles.BackgroundAccent)

	evenRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary)

	oddRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted)

	return func(row, col int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			return headerStyle
		case row%2 == 0:
			return evenRowStyle
		default:
			return oddRowStyle
		}
	}
}

// styleCellForCLI applies CLI-specific styling to table cells
func (t *Table) styleCellForCLI(cell string, colIndex int) string {
	switch colIndex {
	case 0: // ID column
		return styles.Highlight(cell)
	case 1: // Name column
		return styles.Title(cell)
	case 2: // Status column
		return formatStatus(cell)
	default:
		return cell
	}
}

// styleCellForTUI applies TUI-specific styling to table cells
func (t *Table) styleCellForTUI(cell string, colIndex int, isSelected bool) string {
	if isSelected {
		return cell // Selection styling handled at row level
	}

	switch colIndex {
	case 0: // ID column
		return styles.Highlight(cell)
	case 1: // Name column
		return styles.SessionNameStyle.Render(cell)
	case 2: // Status column
		return formatStatus(cell)
	default:
		return styles.TableCellStyle.Render(cell)
	}
}

// Utility functions

func formatStatus(status string) string {
	switch status {
	case "active":
		return styles.StatusActive("active")
	case "inactive":
		return styles.StatusInactive("inactive")
	case "connected":
		return styles.StatusConnected("connected")
	case "error":
		return styles.StatusError("error")
	default:
		return styles.Dim(status)
	}
}

func formatTime(t time.Time) string {
	return styles.Dim(t.Format("2006-01-02 15:04"))
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return styles.Success("just now")
	case duration < time.Hour:
		return styles.Info(fmt.Sprintf("%dm ago", int(duration.Minutes())))
	case duration < 24*time.Hour:
		return styles.Warning(fmt.Sprintf("%dh ago", int(duration.Hours())))
	default:
		return styles.Dim(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
	}
}

func formatProjectPath(path string, maxLen int) string {
	if path == "" {
		return styles.Dim("—")
	}

	if len(path) > maxLen {
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			truncated := ".../" + parts[len(parts)-1]
			if len(truncated) <= maxLen {
				return styles.Dim(truncated)
			}
		}
		return styles.Dim(path[:maxLen-3] + "...")
	}

	return styles.Dim(path)
}
