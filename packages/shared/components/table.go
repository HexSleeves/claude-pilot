package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	lipglosstable "github.com/charmbracelet/lipgloss/table"
	"github.com/evertras/bubble-table/table"
)

// Column key constants for evertras table
const (
	columnKeyID         = "id"
	columnKeyName       = "name"
	columnKeyStatus     = "status"
	columnKeyBackend    = "backend"
	columnKeyCreated    = "created"
	columnKeyLastActive = "last_active"
	columnKeyMessages   = "messages"
	columnKeyProject    = "project"
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
	tbl := lipglosstable.New()

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

// createCLIStyleFunc creates styling for CLI table rendering - Enhanced with standardized theme
func (t *Table) createCLIStyleFunc() lipglosstable.StyleFunc {
	// Use enhanced table header style with Claude orange
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.TextPrimary).
		Background(styles.ClaudePrimary).
		Align(lipgloss.Center)

	// Enhanced row styles with better contrast
	evenRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary)

	oddRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted)

	// Selected row style for interactive tables
	selectedRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary).
		Background(styles.SelectedColor).
		Bold(true)

	return func(row, col int) lipgloss.Style {
		switch {
		case row == lipglosstable.HeaderRow:
			return headerStyle
		case t.config.Interactive && row-1 == t.config.SelectedRow: // Adjust for header row
			return selectedRowStyle
		case row%2 == 0:
			return evenRowStyle
		default:
			return oddRowStyle
		}
	}
}

// styleCellForCLI applies CLI-specific styling to table cells - Enhanced with theme-aware styling
func (t *Table) styleCellForCLI(cell string, colIndex int) string {
	switch colIndex {
	case 0: // ID column - Use muted style for IDs
		return styles.TableCellIDStyle.Render(cell)
	case 1: // Name column - Use primary text with bold
		return styles.TableCellNameStyle.Render(cell)
	case 2: // Status column - Use semantic status colors
		return formatStatusEnhanced(cell)
	case 3: // Backend column - Use secondary text
		return styles.TableCellStyle.Render(cell)
	case 4, 5: // Created/Last Active columns - Use timestamp style
		return styles.TableCellTimestampStyle.Render(cell)
	case 6: // Messages column - Use info style for numbers
		return styles.TableCellInfoStyle.Render(cell)
	case 7: // Project column - Use muted style for paths
		return cell // Already styled in formatProjectPath
	default:
		return styles.TableCellStyle.Render(cell)
	}
}

// styleCellForTUI applies TUI-specific styling to table cells - Enhanced with theme consistency
func (t *Table) styleCellForTUI(cell string, colIndex int, isSelected bool) string {
	if isSelected {
		return cell // Selection styling handled at row level
	}

	switch colIndex {
	case 0: // ID column - Use muted style for IDs
		return styles.TableCellIDStyle.Render(cell)
	case 1: // Name column - Use session name style
		return styles.SessionNameStyle.Render(cell)
	case 2: // Status column - Use enhanced status formatting
		return formatStatusEnhanced(cell)
	case 3: // Backend column - Use secondary text
		return styles.TableCellStyle.Render(cell)
	case 4, 5: // Created/Last Active columns - Use timestamp style
		return styles.TableCellTimestampStyle.Render(cell)
	case 6: // Messages column - Use info style for numbers
		return styles.TableCellInfoStyle.Render(cell)
	case 7: // Project column - Use muted style for paths
		return cell // Already styled in formatProjectPath
	default:
		return styles.TableCellStyle.Render(cell)
	}
}

// Utility functions - Enhanced with comprehensive status formatting

// formatStatusEnhanced provides enhanced status formatting with improved styling
func formatStatusEnhanced(status string) string {
	switch status {
	case "active":
		return styles.TableCellSuccessStyle.Render("● " + status)
	case "inactive":
		return styles.TableCellWarningStyle.Render("⏸ " + status)
	case "connected":
		return styles.TableCellInfoStyle.Render("🔗 " + status)
	case "error", "failed":
		return styles.TableCellErrorStyle.Render("✗ " + status)
	case "starting", "pending":
		return styles.TableCellWarningStyle.Render("⏳ " + status)
	case "stopped":
		return styles.TableCellStyle.Render("⏹ " + status)
	default:
		return styles.TableCellStyle.Render("? " + status)
	}
}

// formatTime formats timestamps with consistent styling
func formatTime(t time.Time) string {
	return styles.TableCellTimestampStyle.Render(t.Format("2006-01-02 15:04"))
}

// formatTimeAgo formats relative time with semantic colors based on recency
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return styles.TableCellSuccessStyle.Render("just now")
	case duration < time.Hour:
		return styles.TableCellInfoStyle.Render(fmt.Sprintf("%dm ago", int(duration.Minutes())))
	case duration < 24*time.Hour:
		return styles.TableCellWarningStyle.Render(fmt.Sprintf("%dh ago", int(duration.Hours())))
	case duration < 7*24*time.Hour:
		return styles.TableCellTimestampStyle.Render(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
	default:
		return styles.TableCellStyle.Render(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
	}
}

// formatProjectPath formats project paths with consistent styling and smart truncation
func formatProjectPath(path string, maxLen int) string {
	if path == "" {
		return styles.TableCellStyle.Render("—")
	}

	if len(path) > maxLen {
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			// Try to show the most relevant part (last directory + filename)
			truncated := ".../" + parts[len(parts)-1]
			if len(truncated) <= maxLen {
				return styles.TableCellStyle.Render(truncated)
			}
		}
		// Fallback to simple truncation
		return styles.TableCellStyle.Render(styles.TruncateText(path, maxLen))
	}

	return styles.TableCellStyle.Render(path)
}

// Evertras Integration Methods
// These methods provide compatibility with the evertras/bubble-table component

// ToEvertrasColumns returns table.Column definitions for evertras table
func (t *Table) ToEvertrasColumns() []table.Column {
	if len(t.data.Headers) == 0 {
		return nil
	}

	// Calculate column widths based on content and terminal width
	columnWidths := styles.GetTableColumnWidths(t.config.Width, len(t.data.Headers))

	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", columnWidths[0]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewFlexColumn(columnKeyName, "Name", 2).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextPrimary).Bold(true),
		),
		table.NewColumn(columnKeyStatus, "Status", columnWidths[2]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewColumn(columnKeyBackend, "Backend", columnWidths[3]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewColumn(columnKeyCreated, "Created", columnWidths[4]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewColumn(columnKeyLastActive, "Last Active", columnWidths[5]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewColumn(columnKeyMessages, "Messages", columnWidths[6]).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewFlexColumn(columnKeyProject, "Project", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
	}

	return columns
}

// ToEvertrasRows converts table data to table.Row format for evertras table
func (t *Table) ToEvertrasRows() []table.Row {
	if len(t.data.Rows) == 0 {
		return nil
	}

	rows := make([]table.Row, len(t.data.Rows))
	for i, row := range t.data.Rows {
		if len(row) >= 8 {
			rows[i] = table.NewRow(table.RowData{
				columnKeyID:         row[0],
				columnKeyName:       row[1],
				columnKeyStatus:     row[2],
				columnKeyBackend:    row[3],
				columnKeyCreated:    row[4],
				columnKeyLastActive: row[5],
				columnKeyMessages:   row[6],
				columnKeyProject:    row[7],
			})
		}
	}

	return rows
}

// ToEvertrasSessionRows converts session data directly to table.Row format
func ToEvertrasSessionRows(sessions []SessionData) []table.Row {
	if len(sessions) == 0 {
		return nil
	}

	rows := make([]table.Row, len(sessions))
	for i, session := range sessions {
		rows[i] = table.NewRow(table.RowData{
			columnKeyID:         styles.TruncateText(session.ID, 11),
			columnKeyName:       styles.TruncateText(session.Name, 19),
			columnKeyStatus:     session.Status,
			columnKeyBackend:    session.Backend,
			columnKeyCreated:    session.Created.Format("2006-01-02 15:04"),
			columnKeyLastActive: formatTimeAgoPlain(session.LastActive),
			columnKeyMessages:   fmt.Sprintf("%d", session.Messages),
			columnKeyProject:    formatProjectPathPlain(session.ProjectPath, 29),
		})
	}

	return rows
}

// GetEvertrasTableColumns returns predefined column definitions for session table
func GetEvertrasTableColumns() []table.Column {
	return []table.Column{
		table.NewColumn(columnKeyID, "ID", 12).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewFlexColumn(columnKeyName, "Name", 2).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextPrimary).Bold(true),
		),
		table.NewColumn(columnKeyStatus, "Status", 10).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewColumn(columnKeyBackend, "Backend", 8).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewColumn(columnKeyCreated, "Created", 16).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewColumn(columnKeyLastActive, "Last Active", 12).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
		table.NewColumn(columnKeyMessages, "Messages", 8).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary),
		),
		table.NewFlexColumn(columnKeyProject, "Project", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted),
		),
	}
}

// AsEvertrasModel configures an evertras table model with data and styling
func (t *Table) AsEvertrasModel(baseModel table.Model) table.Model {
	// Set columns and rows
	model := baseModel.
		WithColumns(t.ToEvertrasColumns()).
		WithRows(t.ToEvertrasRows())

	// Configure dimensions
	if t.config.Width > 0 {
		model = model.WithTargetWidth(t.config.Width)
	}
	if t.config.MaxRows > 0 {
		model = model.WithMinimumHeight(t.config.MaxRows)
	}

	return model
}

// Utility functions for plain text formatting (without lipgloss styling)
// These are used for Bubbles table content where styling is handled by the table component

// formatTimeAgoPlain formats relative time without styling
func formatTimeAgoPlain(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	case duration < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	default:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

// formatProjectPathPlain formats project paths without styling
func formatProjectPathPlain(path string, maxLen int) string {
	if path == "" {
		return "—"
	}

	if len(path) > maxLen {
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			// Try to show the most relevant part (last directory + filename)
			truncated := ".../" + parts[len(parts)-1]
			if len(truncated) <= maxLen {
				return truncated
			}
		}
		// Fallback to simple truncation
		return styles.TruncateText(path, maxLen)
	}

	return path
}
