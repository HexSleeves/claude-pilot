package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"sort"
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
	columnKeyProject    = "project"
	columnKeyMessages   = "messages"
)

// TableConfig holds configuration for table rendering
type TableConfig struct {
	Width       int
	ShowHeaders bool
	Interactive bool
	SelectedRow int
	MaxRows     int

	// Sorting configuration
	SortEnabled   bool
	SortColumn    string
	SortDirection string // "asc" or "desc"

	// Pagination configuration
	PageSize    int
	CurrentPage int

	// Filtering configuration
	FilterEnabled bool
	FilterText    string

	// Multi-selection configuration
	MultiSelectEnabled bool
	SelectedRows       []int

	// Display enhancements
	HoverEnabled    bool
	ShowRowNumbers  bool
	ColumnWidthMode string // "fixed", "flex", "auto"
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
	headers := []string{"ID", "Name", "Status", "Backend", "Created", "Last Active", "Project", "Messages"}
	rows := make([][]string, len(sessions))

	for i, session := range sessions {
		rows[i] = []string{
			styles.TruncateText(session.ID, 11),
			styles.TruncateText(session.Name, 19),
			session.Status,
			session.Backend,
			formatTime(session.Created),
			formatTimeAgo(session.LastActive),
			formatProjectPath(session.ProjectPath, 29),
			fmt.Sprintf("%d", session.Messages),
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

// Sorting Methods

// EnableSorting enables sorting functionality for the table
func (t *Table) EnableSorting() {
	t.config.SortEnabled = true
}

// DisableSorting disables sorting functionality for the table
func (t *Table) DisableSorting() {
	t.config.SortEnabled = false
}

// SetSort sets the sort column and direction
func (t *Table) SetSort(column, direction string) error {
	if !t.validateSortColumn(column) {
		return fmt.Errorf("invalid sort column: %s", column)
	}
	if direction != "asc" && direction != "desc" {
		return fmt.Errorf("invalid sort direction: %s (must be 'asc' or 'desc')", direction)
	}

	t.config.SortColumn = column
	t.config.SortDirection = direction
	return nil
}

// Pagination Methods

// EnablePagination enables pagination with the specified page size
func (t *Table) EnablePagination(pageSize int) {
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	t.config.PageSize = pageSize
	t.config.CurrentPage = 1
}

// NextPage moves to the next page if available
func (t *Table) NextPage() bool {
	if t.config.PageSize <= 0 {
		return false
	}

	totalPages := t.GetTotalPages()
	if t.config.CurrentPage < totalPages {
		t.config.CurrentPage++
		return true
	}
	return false
}

// PrevPage moves to the previous page if available
func (t *Table) PrevPage() bool {
	if t.config.PageSize <= 0 || t.config.CurrentPage <= 1 {
		return false
	}

	t.config.CurrentPage--
	return true
}

// GoToPage navigates to a specific page
func (t *Table) GoToPage(page int) error {
	if t.config.PageSize <= 0 {
		return fmt.Errorf("pagination not enabled")
	}

	if !t.validatePageNumber(page) {
		return fmt.Errorf("invalid page number: %d", page)
	}

	t.config.CurrentPage = page
	return nil
}

// Filtering Methods

// SetFilter sets the filter text for table data
func (t *Table) SetFilter(text string) {
	t.config.FilterEnabled = true
	t.config.FilterText = text
}

// ClearFilter removes the current filter
func (t *Table) ClearFilter() {
	t.config.FilterEnabled = false
	t.config.FilterText = ""
}

// Selection Methods

// SelectRow selects a row by index (for multi-select mode)
func (t *Table) SelectRow(index int) error {
	if !t.config.MultiSelectEnabled {
		return fmt.Errorf("multi-select not enabled")
	}

	if index < 0 || index >= len(t.data.Rows) {
		return fmt.Errorf("invalid row index: %d", index)
	}

	// Check if already selected
	for _, selected := range t.config.SelectedRows {
		if selected == index {
			return nil // Already selected
		}
	}

	t.config.SelectedRows = append(t.config.SelectedRows, index)
	return nil
}

// DeselectRow deselects a row by index
func (t *Table) DeselectRow(index int) error {
	if !t.config.MultiSelectEnabled {
		return fmt.Errorf("multi-select not enabled")
	}

	for i, selected := range t.config.SelectedRows {
		if selected == index {
			// Remove from slice
			t.config.SelectedRows = append(t.config.SelectedRows[:i], t.config.SelectedRows[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("row %d not selected", index)
}

// SelectAll selects all visible rows
func (t *Table) SelectAll() {
	if !t.config.MultiSelectEnabled {
		return
	}

	t.config.SelectedRows = make([]int, len(t.data.Rows))
	for i := range t.data.Rows {
		t.config.SelectedRows[i] = i
	}
}

// ClearSelection clears all selected rows
func (t *Table) ClearSelection() {
	t.config.SelectedRows = []int{}
}

// Data Retrieval Methods

// GetSortedData returns the table data sorted according to current sort settings
func (t *Table) GetSortedData() [][]string {
	if !t.config.SortEnabled || t.config.SortColumn == "" {
		return t.data.Rows
	}

	// Create a copy to avoid modifying original data
	sortedRows := make([][]string, len(t.data.Rows))
	copy(sortedRows, t.data.Rows)

	return t.sortTableData(sortedRows, t.config.SortColumn, t.config.SortDirection)
}

// GetFilteredData returns the table data filtered according to current filter settings
func (t *Table) GetFilteredData() [][]string {
	if !t.config.FilterEnabled || t.config.FilterText == "" {
		return t.data.Rows
	}

	return t.filterTableData(t.data.Rows, t.config.FilterText)
}

// GetPagedData returns the current page of data
func (t *Table) GetPagedData() [][]string {
	if t.config.PageSize <= 0 {
		return t.data.Rows
	}

	return t.paginateData(t.data.Rows, t.config.CurrentPage, t.config.PageSize)
}

// GetTotalPages returns the total number of pages based on current data and page size
func (t *Table) GetTotalPages() int {
	if t.config.PageSize <= 0 {
		return 1
	}

	totalRows := len(t.data.Rows)
	return (totalRows + t.config.PageSize - 1) / t.config.PageSize
}

// GetCurrentPageInfo returns information about the current page
func (t *Table) GetCurrentPageInfo() (currentPage, totalPages, startRow, endRow, totalRows int) {
	totalRows = len(t.data.Rows)

	if t.config.PageSize <= 0 {
		return 1, 1, 1, totalRows, totalRows
	}

	currentPage = t.config.CurrentPage
	totalPages = t.GetTotalPages()

	startRow = (currentPage-1)*t.config.PageSize + 1
	endRow = currentPage * t.config.PageSize
	if endRow > totalRows {
		endRow = totalRows
	}

	return currentPage, totalPages, startRow, endRow, totalRows
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

// Data Manipulation Utility Methods

// validateSortColumn validates if the given column is valid for sorting
func (t *Table) validateSortColumn(column string) bool {
	validColumns := []string{"id", "name", "status", "backend", "created", "last_active", "project", "messages"}
	for _, valid := range validColumns {
		if column == valid {
			return true
		}
	}
	return false
}

// validatePageNumber validates if the given page number is valid
func (t *Table) validatePageNumber(page int) bool {
	if t.config.PageSize <= 0 {
		return false
	}

	totalPages := t.GetTotalPages()
	return page >= 1 && page <= totalPages
}

// sortTableData sorts table data by the specified column and direction
func (t *Table) sortTableData(rows [][]string, column, direction string) [][]string {
	if len(rows) == 0 {
		return rows
	}

	// Get column index
	columnIndex := t.getColumnIndex(column)
	if columnIndex == -1 {
		return rows // Invalid column, return unsorted
	}

	// Sort the rows
	sortedRows := make([][]string, len(rows))
	copy(sortedRows, rows)

	sort.Slice(sortedRows, func(i, j int) bool {
		if columnIndex >= len(sortedRows[i]) || columnIndex >= len(sortedRows[j]) {
			return false
		}

		val1 := sortedRows[i][columnIndex]
		val2 := sortedRows[j][columnIndex]

		// Handle different column types
		switch column {
		case "created", "last_active":
			// Parse time values for proper sorting
			time1, err1 := time.Parse("2006-01-02 15:04", val1)
			time2, err2 := time.Parse("2006-01-02 15:04", val2)
			if err1 == nil && err2 == nil {
				if direction == "desc" {
					return time1.After(time2)
				}
				return time1.Before(time2)
			}
		case "messages":
			// Parse numeric values
			var num1, num2 int
			fmt.Sscanf(val1, "%d", &num1)
			fmt.Sscanf(val2, "%d", &num2)
			if direction == "desc" {
				return num1 > num2
			}
			return num1 < num2
		}

		// Default string comparison
		if direction == "desc" {
			return val1 > val2
		}
		return val1 < val2
	})

	return sortedRows
}

// filterTableData filters table data based on the filter text
func (t *Table) filterTableData(rows [][]string, filterText string) [][]string {
	if filterText == "" {
		return rows
	}

	filterText = strings.ToLower(filterText)
	var filteredRows [][]string

	for _, row := range rows {
		// Check if any cell in the row contains the filter text
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), filterText) {
				filteredRows = append(filteredRows, row)
				break
			}
		}
	}

	return filteredRows
}

// paginateData returns a specific page of data
func (t *Table) paginateData(rows [][]string, page, pageSize int) [][]string {
	if pageSize <= 0 || page <= 0 {
		return rows
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(rows) {
		return [][]string{} // Page beyond available data
	}

	if endIndex > len(rows) {
		endIndex = len(rows)
	}

	return rows[startIndex:endIndex]
}

// getColumnIndex returns the index of a column by name
func (t *Table) getColumnIndex(column string) int {
	columnMap := map[string]int{
		"id":          0,
		"name":        1,
		"status":      2,
		"backend":     3,
		"created":     4,
		"last_active": 5,
		"project":     6,
		"messages":    7,
	}

	if index, exists := columnMap[column]; exists {
		return index
	}
	return -1
}

// sortSessionData sorts session data by the specified column and direction
func sortSessionData(sessions []SessionData, column, direction string) []SessionData {
	if len(sessions) == 0 {
		return sessions
	}

	sortedSessions := make([]SessionData, len(sessions))
	copy(sortedSessions, sessions)

	sort.Slice(sortedSessions, func(i, j int) bool {
		switch column {
		case "id":
			if direction == "desc" {
				return sortedSessions[i].ID > sortedSessions[j].ID
			}
			return sortedSessions[i].ID < sortedSessions[j].ID
		case "name":
			if direction == "desc" {
				return sortedSessions[i].Name > sortedSessions[j].Name
			}
			return sortedSessions[i].Name < sortedSessions[j].Name
		case "status":
			if direction == "desc" {
				return sortedSessions[i].Status > sortedSessions[j].Status
			}
			return sortedSessions[i].Status < sortedSessions[j].Status
		case "backend":
			if direction == "desc" {
				return sortedSessions[i].Backend > sortedSessions[j].Backend
			}
			return sortedSessions[i].Backend < sortedSessions[j].Backend
		case "created":
			if direction == "desc" {
				return sortedSessions[i].Created.After(sortedSessions[j].Created)
			}
			return sortedSessions[i].Created.Before(sortedSessions[j].Created)
		case "last_active":
			if direction == "desc" {
				return sortedSessions[i].LastActive.After(sortedSessions[j].LastActive)
			}
			return sortedSessions[i].LastActive.Before(sortedSessions[j].LastActive)
		case "project":
			if direction == "desc" {
				return sortedSessions[i].ProjectPath > sortedSessions[j].ProjectPath
			}
			return sortedSessions[i].ProjectPath < sortedSessions[j].ProjectPath
		case "messages":
			if direction == "desc" {
				return sortedSessions[i].Messages > sortedSessions[j].Messages
			}
			return sortedSessions[i].Messages < sortedSessions[j].Messages
		default:
			return false
		}
	})

	return sortedSessions
}

// filterSessionData filters session data based on the filter text
func filterSessionData(sessions []SessionData, filterText string) []SessionData {
	if filterText == "" {
		return sessions
	}

	filterText = strings.ToLower(filterText)
	var filteredSessions []SessionData

	for _, session := range sessions {
		// Check if any field contains the filter text
		if strings.Contains(strings.ToLower(session.ID), filterText) ||
			strings.Contains(strings.ToLower(session.Name), filterText) ||
			strings.Contains(strings.ToLower(session.Status), filterText) ||
			strings.Contains(strings.ToLower(session.Backend), filterText) ||
			strings.Contains(strings.ToLower(session.ProjectPath), filterText) {
			filteredSessions = append(filteredSessions, session)
		}
	}

	return filteredSessions
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

// ToEvertrasColumns returns table.Column definitions for evertras table with enhanced features
func (t *Table) ToEvertrasColumns() []table.Column {
	if len(t.data.Headers) == 0 {
		return nil
	}

	// Calculate column widths based on content and terminal width
	columnWidths := styles.GetTableColumnWidths(t.config.Width, len(t.data.Headers))

	// Apply column width mode
	switch t.config.ColumnWidthMode {
	case "fixed":
		return t.createFixedWidthColumns(columnWidths)
	case "flex":
		return t.createFlexColumns()
	case "auto":
		return t.createAutoColumns(columnWidths)
	default:
		return t.createDefaultColumns(columnWidths)
	}
}

// createDefaultColumns creates columns with default configuration
func (t *Table) createDefaultColumns(columnWidths []int) []table.Column {
	return []table.Column{
		t.createEvertrasColumn(columnKeyID, "ID", columnWidths[0], styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyName, "Name", 2, styles.TextPrimary, true),
		t.createEvertrasColumn(columnKeyStatus, "Status", columnWidths[2], styles.TextSecondary, false),
		t.createEvertrasColumn(columnKeyBackend, "Backend", columnWidths[3], styles.TextSecondary, false),
		t.createEvertrasColumn(columnKeyCreated, "Created", columnWidths[4], styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyLastActive, "Last Active", columnWidths[5], styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyProject, "Project", 1, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyMessages, "Messages", columnWidths[6], styles.TextSecondary, false),
	}
}

// createFixedWidthColumns creates columns with fixed widths
func (t *Table) createFixedWidthColumns(columnWidths []int) []table.Column {
	return []table.Column{
		t.createEvertrasColumn(columnKeyID, "ID", 12, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyName, "Name", 20, styles.TextPrimary, true),
		t.createEvertrasColumn(columnKeyStatus, "Status", 10, styles.TextSecondary, false),
		t.createEvertrasColumn(columnKeyBackend, "Backend", 8, styles.TextSecondary, false),
		t.createEvertrasColumn(columnKeyCreated, "Created", 16, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyLastActive, "Last Active", 12, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyProject, "Project", 30, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyMessages, "Messages", 8, styles.TextSecondary, false),
	}
}

// createFlexColumns creates columns with flexible widths
func (t *Table) createFlexColumns() []table.Column {
	return []table.Column{
		t.createEvertrasFlexColumn(columnKeyID, "ID", 1, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyName, "Name", 3, styles.TextPrimary, true),
		t.createEvertrasFlexColumn(columnKeyStatus, "Status", 1, styles.TextSecondary, false),
		t.createEvertrasFlexColumn(columnKeyBackend, "Backend", 1, styles.TextSecondary, false),
		t.createEvertrasFlexColumn(columnKeyCreated, "Created", 2, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyLastActive, "Last Active", 2, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyProject, "Project", 3, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyMessages, "Messages", 1, styles.TextSecondary, false),
	}
}

// createAutoColumns creates columns with automatic width calculation
func (t *Table) createAutoColumns(columnWidths []int) []table.Column {
	// Use a mix of fixed and flex columns for optimal display
	return []table.Column{
		t.createEvertrasColumn(columnKeyID, "ID", 12, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyName, "Name", 2, styles.TextPrimary, true),
		t.createEvertrasColumn(columnKeyStatus, "Status", 10, styles.TextSecondary, false),
		t.createEvertrasColumn(columnKeyBackend, "Backend", 8, styles.TextSecondary, false),
		t.createEvertrasFlexColumn(columnKeyCreated, "Created", 1, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyLastActive, "Last Active", 1, styles.TextMuted, false),
		t.createEvertrasFlexColumn(columnKeyProject, "Project", 2, styles.TextMuted, false),
		t.createEvertrasColumn(columnKeyMessages, "Messages", 8, styles.TextSecondary, false),
	}
}

// createEvertrasColumn creates a standard column with enhanced styling
func (t *Table) createEvertrasColumn(key, title string, width int, color lipgloss.Color, bold bool) table.Column {
	style := lipgloss.NewStyle().Foreground(color)
	if bold {
		style = style.Bold(true)
	}

	return table.NewColumn(key, title, width).WithStyle(style)
}

// createEvertrasFlexColumn creates a flex column with enhanced styling
func (t *Table) createEvertrasFlexColumn(key, title string, flex int, color lipgloss.Color, bold bool) table.Column {
	style := lipgloss.NewStyle().Foreground(color)
	if bold {
		style = style.Bold(true)
	}

	return table.NewFlexColumn(key, title, flex).WithStyle(style)
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
				columnKeyProject:    row[6],
				columnKeyMessages:   row[7],
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
		table.NewFlexColumn(columnKeyName, "Name", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextPrimary).Bold(true).Align(lipgloss.Center),
		),
		table.NewColumn(columnKeyStatus, "Status", 10).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary).Align(lipgloss.Center),
		),
		table.NewColumn(columnKeyBackend, "Backend", 8).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary).Align(lipgloss.Center),
		),
		table.NewFlexColumn(columnKeyCreated, "Created", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted).Align(lipgloss.Center),
		),
		table.NewFlexColumn(columnKeyLastActive, "Last Active", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted).Align(lipgloss.Center),
		),
		table.NewFlexColumn(columnKeyProject, "Project", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextMuted).Align(lipgloss.Center),
		),
		table.NewFlexColumn(columnKeyMessages, "Messages", 1).WithStyle(
			lipgloss.NewStyle().Foreground(styles.TextSecondary).Align(lipgloss.Center),
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

	// Apply interactive features
	model = t.ConfigureInteractiveFeatures(model)

	return model
}

// ConfigureInteractiveFeatures applies all interactive settings to the evertras model
func (t *Table) ConfigureInteractiveFeatures(model table.Model) table.Model {
	// Configure pagination if enabled
	if t.config.PageSize > 0 {
		model = model.WithPageSize(t.config.PageSize)
		if t.config.CurrentPage > 1 {
			// Note: evertras pagination is 0-based
			model = model.WithCurrentPage(t.config.CurrentPage - 1)
		}
	}

	// Configure multi-selection if enabled
	if t.config.MultiSelectEnabled {
		model = model.WithMultiline(true)
		// Note: Individual row selection must be handled through key events
		// in the consuming TUI application
	}

	// Apply base styling
	model = model.WithBaseStyle(styles.GetEvertrasTableStyles())

	return model
}

// UpdateEvertrasModel refreshes the evertras model with current state
func (t *Table) UpdateEvertrasModel(model table.Model) table.Model {
	// Get processed data based on current filters, sorting, and pagination
	processedData := t.GetProcessedData()

	// Update rows with processed data
	rows := t.convertToEvertrasRows(processedData)
	model = model.WithRows(rows)

	// Update interactive features
	model = t.ConfigureInteractiveFeatures(model)

	return model
}

// GetProcessedData returns data with all current processing applied (filter, sort, paginate)
func (t *Table) GetProcessedData() [][]string {
	data := t.data.Rows

	// Apply filtering first
	if t.config.FilterEnabled && t.config.FilterText != "" {
		data = t.filterTableData(data, t.config.FilterText)
	}

	// Apply sorting
	if t.config.SortEnabled && t.config.SortColumn != "" {
		data = t.sortTableData(data, t.config.SortColumn, t.config.SortDirection)
	}

	// Apply pagination
	if t.config.PageSize > 0 {
		data = t.paginateData(data, t.config.CurrentPage, t.config.PageSize)
	}

	return data
}

// convertToEvertrasRows converts processed data to evertras rows
func (t *Table) convertToEvertrasRows(data [][]string) []table.Row {
	if len(data) == 0 {
		return nil
	}

	rows := make([]table.Row, len(data))
	for i, row := range data {
		if len(row) >= 8 {
			rowData := table.RowData{
				columnKeyID:         row[0],
				columnKeyName:       row[1],
				columnKeyStatus:     row[2],
				columnKeyBackend:    row[3],
				columnKeyCreated:    row[4],
				columnKeyLastActive: row[5],
				columnKeyProject:    row[6],
				columnKeyMessages:   row[7],
			}

			// Add row numbers if enabled
			if t.config.ShowRowNumbers {
				rowData["row_number"] = fmt.Sprintf("%d", i+1)
			}

			rows[i] = table.NewRow(rowData)
		}
	}

	return rows
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
