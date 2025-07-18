package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

	// Note: Removed unused multi-selection and display enhancement options
	// These features are not currently utilized in CLI or TUI
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

// Table provides a unified table component wrapping evertras/bubble-table
// This is the primary table interface for both CLI and TUI components
type Table struct {
	config TableConfig
	data   TableData

	// Evertras table model - this is the core wrapped component
	evertrasModel table.Model

	// Internal state
	initialized bool
	keyMap      TableKeyMap
}

// TableKeyMap defines key bindings for table interaction
type TableKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Enter    key.Binding
	Space    key.Binding
	Filter   key.Binding
	Sort     key.Binding
}

// DefaultTableKeyMap returns the default key bindings for table interaction
func DefaultTableKeyMap() TableKeyMap {
	return TableKeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		PageUp:   key.NewBinding(key.WithKeys("pgup", "b"), key.WithHelp("pgup/b", "page up")),
		PageDown: key.NewBinding(key.WithKeys("pgdown", "f"), key.WithHelp("pgdn/f", "page down")),
		Home:     key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("home/g", "top")),
		End:      key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("end/G", "bottom")),
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Space:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
		Filter:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Sort:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "sort")),
	}
}

// NewTable creates a new table instance with evertras wrapper
func NewTable(config TableConfig) *Table {
	t := &Table{
		config: config,
		keyMap: DefaultTableKeyMap(),
	}

	// Initialize the evertras model
	t.initializeEvertrasModel()

	return t
}

// initializeEvertrasModel sets up the underlying evertras table model
func (t *Table) initializeEvertrasModel() {
	// Create base model with default configuration
	t.evertrasModel = table.New([]table.Column{}).WithBaseStyle(
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.ClaudePrimary).
			Align(lipgloss.Left),
	)

	// Apply initial configuration
	if t.config.Width > 0 {
		t.evertrasModel = t.evertrasModel.WithTargetWidth(t.config.Width)
	}

	if t.config.MaxRows > 0 {
		t.evertrasModel = t.evertrasModel.WithMinimumHeight(t.config.MaxRows)
	}

	if t.config.PageSize > 0 {
		t.evertrasModel = t.evertrasModel.WithPageSize(t.config.PageSize)
	}

	// Configure for interactive use if enabled
	if t.config.Interactive {
		t.evertrasModel = t.evertrasModel.Focused(true)
	}

	t.initialized = true
}

// SetData sets the table data and updates the evertras model
func (t *Table) SetData(data TableData) {
	t.data = data
	t.refreshEvertrasModel()
}

// refreshEvertrasModel updates the evertras model with current data and config
func (t *Table) refreshEvertrasModel() {
	if !t.initialized {
		t.initializeEvertrasModel()
	}

	// Update columns
	columns := t.ToEvertrasColumns()
	if len(columns) > 0 {
		t.evertrasModel = t.evertrasModel.WithColumns(columns)
	}

	// Update rows with processed data
	rows := t.ToEvertrasRows()
	if len(rows) > 0 {
		t.evertrasModel = t.evertrasModel.WithRows(rows)
	}

	// Apply current configuration
	t.applyConfigToEvertrasModel()
}

// applyConfigToEvertrasModel applies current config to the evertras model
func (t *Table) applyConfigToEvertrasModel() {
	if t.config.Width > 0 {
		t.evertrasModel = t.evertrasModel.WithTargetWidth(t.config.Width)
	}

	if t.config.MaxRows > 0 {
		t.evertrasModel = t.evertrasModel.WithMinimumHeight(t.config.MaxRows)
	}

	if t.config.PageSize > 0 {
		t.evertrasModel = t.evertrasModel.WithPageSize(t.config.PageSize)
	}

	// Apply interactive state
	if t.config.Interactive {
		t.evertrasModel = t.evertrasModel.Focused(true)
	} else {
		t.evertrasModel = t.evertrasModel.Focused(false)
	}
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

// SetSelectedRow updates the selected row for interactive mode
func (t *Table) SetSelectedRow(row int) {
	if row >= 0 && row < len(t.data.Rows) {
		t.config.SelectedRow = row
		// Update evertras model if initialized
		if t.initialized {
			// Note: evertras handles row selection through its own state
			// This is maintained for compatibility with existing CLI rendering
		}
	}
}

// SetWidth sets the table width and updates the evertras model
func (t *Table) SetWidth(width int) {
	t.config.Width = width
	if t.initialized {
		t.evertrasModel = t.evertrasModel.WithTargetWidth(width)
	}
}

// SetMaxRows sets the maximum number of rows to display
func (t *Table) SetMaxRows(maxRows int) {
	t.config.MaxRows = maxRows
	if t.initialized {
		t.evertrasModel = t.evertrasModel.WithMinimumHeight(maxRows)
	}
}

// GetSelectedRow returns the currently selected row index
// For TUI mode, this returns the evertras highlighted row
// For CLI mode, this returns the internal config value
func (t *Table) GetSelectedRow() int {
	if t.initialized && t.config.Interactive {
		return t.evertrasModel.GetHighlightedRowIndex()
	}
	return t.config.SelectedRow
}

// GetHighlightedRowIndex returns the highlighted row from evertras model
// This is the preferred method for TUI interactions
func (t *Table) GetHighlightedRowIndex() int {
	if t.initialized {
		return t.evertrasModel.GetHighlightedRowIndex()
	}
	return -1
}

// GetRowCount returns the number of rows in the table
func (t *Table) GetRowCount() int {
	return len(t.data.Rows)
}

// GetSelectedData returns the data for the currently selected row
func (t *Table) GetSelectedData() []string {
	selectedRow := t.GetSelectedRow()
	if selectedRow >= 0 && selectedRow < len(t.data.Rows) {
		return t.data.Rows[selectedRow]
	}
	return nil
}

// Note: Removed GetSelectedRows as multi-select is not used

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

	// Apply to evertras model if initialized
	if t.initialized {
		t.SortByColumn(column, direction == "asc")
	}

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

// Note: Removed unused selection methods (SelectRow, DeselectRow, etc.)
// Multi-select functionality is not used in current implementation

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

// Note: Removed duplicate sortSessionData and filterSessionData functions
// These are redundant as the table already has sortTableData and filterTableData methods

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

	// Use default column configuration
	return t.createDefaultColumns(columnWidths)
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

// Note: Removed unused createFixedWidthColumns, createFlexColumns, and createAutoColumns methods
// These complex column creation methods were never used as ColumnWidthMode is not utilized

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

// Note: Removed unused ToEvertrasSessionRows and GetEvertrasTableColumns functions
// These are replaced by instance methods ToEvertrasRows and ToEvertrasColumns

// Note: Removed AsEvertrasModel - not used in current implementation

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

	// Note: Multi-selection configuration removed (not used)

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

			// Note: Row numbers feature removed (not used)

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

// ===== BUBBLE TEA INTERFACE METHODS =====
// These methods provide the Bubble Tea interface for TUI usage

// Init initializes the table for Bubble Tea
func (t *Table) Init() tea.Cmd {
	if !t.initialized {
		t.initializeEvertrasModel()
	}
	return t.evertrasModel.Init()
}

// Update handles Bubble Tea messages and key events
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.initialized {
		t.initializeEvertrasModel()
	}

	// Pass all key events directly to evertras for better performance

	// Pass message to evertras model
	var cmd tea.Cmd
	t.evertrasModel, cmd = t.evertrasModel.Update(msg)

	return t, cmd
}

// View renders the table for Bubble Tea
func (t *Table) View() string {
	if !t.initialized {
		t.initializeEvertrasModel()
	}
	return t.evertrasModel.View()
}

// handleKeyEvent processes key events for table navigation
// Simplified to let evertras handle navigation directly
func (t *Table) handleKeyEvent(msg tea.KeyMsg) tea.Cmd {
	// Let evertras table handle all key events natively
	return nil
}

// Note: Removed unused navigation methods (MoveUp, MoveDown, etc.)
// These wrapper methods are not used by the TUI, which handles navigation
// through the underlying evertras model directly

// ===== EVERTRAS MODEL ACCESS =====
// These methods provide direct access to the underlying evertras model

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

// GetEvertrasModel returns the underlying evertras table model
// Use this for advanced evertras-specific functionality
func (t *Table) GetEvertrasModel() table.Model {
	if !t.initialized {
		t.initializeEvertrasModel()
	}
	return t.evertrasModel
}

// SetEvertrasModel allows setting a custom evertras model
// Use this for advanced customization
func (t *Table) SetEvertrasModel(model table.Model) {
	t.evertrasModel = model
	t.initialized = true
}

// WithEvertrasModel applies a transformation to the evertras model
// This allows chaining evertras-specific methods
func (t *Table) WithEvertrasModel(fn func(table.Model) table.Model) *Table {
	if !t.initialized {
		t.initializeEvertrasModel()
	}
	t.evertrasModel = fn(t.evertrasModel)
	return t
}

// ===== FOCUS AND INTERACTION =====

// Focus sets the table as focused for interactive use
func (t *Table) Focus() {
	t.config.Interactive = true
	if t.initialized {
		t.evertrasModel = t.evertrasModel.Focused(true)
	}
}

// Blur removes focus from the table
func (t *Table) Blur() {
	t.config.Interactive = false
	if t.initialized {
		t.evertrasModel = t.evertrasModel.Focused(false)
	}
}

// IsFocused returns whether the table is currently focused
func (t *Table) IsFocused() bool {
	if t.initialized {
		return t.evertrasModel.GetFocused()
	}
	return t.config.Interactive
}

// Note: Removed unused selection methods (ToggleRowSelection, SelectAllRows, etc.)
// Multi-select functionality is not currently used in CLI or TUI

// ===== SORTING METHODS =====

// SortByColumn sorts the table by the specified column
func (t *Table) SortByColumn(columnKey string, ascending bool) {
	if t.initialized {
		if ascending {
			t.evertrasModel = t.evertrasModel.SortByAsc(columnKey)
		} else {
			t.evertrasModel = t.evertrasModel.SortByDesc(columnKey)
		}
	}

	// Update internal config
	t.config.SortColumn = columnKey
	if ascending {
		t.config.SortDirection = "asc"
	} else {
		t.config.SortDirection = "desc"
	}
}

// ===== FILTERING METHODS =====

// SetFilterText sets the filter text and applies it
func (t *Table) SetFilterText(text string) {
	t.SetFilter(text)
	if t.initialized {
		// Custom filter function that matches our existing logic
		t.evertrasModel = t.evertrasModel.WithFilterFunc(func(row table.Row, filterInput string) bool {
			if filterInput == "" {
				return true
			}

			lowerText := strings.ToLower(filterInput)
			// Check all column values
			for key := range row.Data {
				if val, ok := row.Data[key].(string); ok {
					if strings.Contains(strings.ToLower(val), lowerText) {
						return true
					}
				}
			}
			return false
		})
	}
}

// GetFilterText returns the current filter text
func (t *Table) GetFilterText() string {
	if t.initialized {
		return t.evertrasModel.GetCurrentFilter()
	}
	return t.config.FilterText
}

// CanFilter returns whether filtering is enabled
func (t *Table) CanFilter() bool {
	if t.initialized {
		return t.evertrasModel.GetCanFilter()
	}
	return t.config.FilterEnabled
}
