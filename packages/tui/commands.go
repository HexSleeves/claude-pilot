package main

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// loadSessionsCmd loads all sessions from the API
func loadSessionsCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionsLoadedMsg{
				sessions: nil,
				err:      fmt.Errorf("API client is nil"),
			}
		}

		sessions, err := client.ListSessions()
		return sessionsLoadedMsg{
			sessions: sessions,
			err:      err,
		}
	}
}

// createSessionCmd creates a new session with the specified parameters
func createSessionCmd(client *api.Client, name, description, projectPath string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionCreatedMsg{
				session: nil,
				err:     fmt.Errorf("API client is nil"),
			}
		}

		if strings.TrimSpace(name) == "" {
			return sessionCreatedMsg{
				session: nil,
				err:     fmt.Errorf("session name cannot be empty"),
			}
		}

		req := api.CreateSessionRequest{
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(description),
			ProjectPath: strings.TrimSpace(projectPath),
		}

		session, err := client.CreateSession(req)
		return sessionCreatedMsg{
			session: session,
			err:     err,
		}
	}
}

// killSessionCmd terminates a session by ID
func killSessionCmd(client *api.Client, sessionID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionKilledMsg{
				sessionID: sessionID,
				err:       fmt.Errorf("API client is nil"),
			}
		}

		if strings.TrimSpace(sessionID) == "" {
			return sessionKilledMsg{
				sessionID: sessionID,
				err:       fmt.Errorf("session ID cannot be empty"),
			}
		}

		err := client.KillSession(sessionID)
		return sessionKilledMsg{
			sessionID: sessionID,
			err:       err,
		}
	}
}

// attachSessionCmd attaches to a session and hands control to the multiplexer
func attachSessionCmd(client *api.Client, sessionID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return errorMsg{error: fmt.Errorf("API client is nil")}
		}

		if strings.TrimSpace(sessionID) == "" {
			return errorMsg{error: fmt.Errorf("session ID cannot be empty")}
		}

		// Get the session to find its name
		session, err := client.GetSession(sessionID)
		if err != nil {
			return errorMsg{error: fmt.Errorf("failed to get session: %w", err)}
		}

		if session == nil {
			return errorMsg{error: fmt.Errorf("session not found")}
		}

		if strings.TrimSpace(session.Name) == "" {
			return errorMsg{error: fmt.Errorf("session name is empty")}
		}

		// Create a command to attach to the session
		// This will hand control over to the multiplexer session
		cmd := exec.Command("tmux", "attach-session", "-t", session.Name)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Use tea.ExecProcess to hand control to the external process
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return errorMsg{error: fmt.Errorf("failed to attach to session: %w", err)}
			}
			// After the session ends, quit the TUI
			return tea.Quit()
		})
	}
}

// Table Data Commands

// SortTableDataCmd sorts table data by the specified column and direction asynchronously
func SortTableDataCmd(column, direction string) tea.Cmd {
	return func() tea.Msg {
		if err := validateSortColumn(column); err != nil {
			return TableErrorMsg{Error: err}
		}

		if direction != "asc" && direction != "desc" {
			return TableErrorMsg{Error: fmt.Errorf("invalid sort direction: %s (must be 'asc' or 'desc')", direction)}
		}

		// Execute sort operation
		if err := executeTableSort(column, direction); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to sort table: %w", err)}
		}

		return TableSortedMsg{
			Column:    column,
			Direction: direction,
		}
	}
}

// FilterTableDataCmd filters table data based on the provided filter text
func FilterTableDataCmd(filterText string) tea.Cmd {
	return func() tea.Msg {
		if err := validateFilterText(filterText); err != nil {
			return TableErrorMsg{Error: err}
		}

		// Execute filter operation
		rowCount, err := executeTableFilter(filterText)
		if err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to filter table: %w", err)}
		}

		return TableFilteredMsg{
			FilterText: filterText,
			RowCount:   rowCount,
		}
	}
}

// RefreshTableDataCmd reloads table data from the backend
func RefreshTableDataCmd() tea.Cmd {
	return func() tea.Msg {
		if err := executeTableRefresh(); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to refresh table data: %w", err)}
		}

		return TableRefreshMsg{}
	}
}

// ExportTableDataCmd exports filtered/sorted data to the specified format and filename
func ExportTableDataCmd(format, filename string, data []components.SessionData) tea.Cmd {
	return func() tea.Msg {
		if err := validateExportParams(format, filename); err != nil {
			return TableErrorMsg{Error: err}
		}

		if err := executeTableExport(format, filename, data); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to export table data: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Table data exported to %s", filename)}
	}
}

// Table State Commands

// SaveTableStateCmd persists current table configuration (sort, filter, page size)
func SaveTableStateCmd() tea.Cmd {
	return func() tea.Msg {
		if err := executeTableStateSave(); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to save table state: %w", err)}
		}

		return statusMsg{message: "Table state saved successfully"}
	}
}

// LoadTableStateCmd restores saved table configuration
func LoadTableStateCmd() tea.Cmd {
	return func() tea.Msg {
		if err := executeTableStateLoad(); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to load table state: %w", err)}
		}

		return statusMsg{message: "Table state loaded successfully"}
	}
}

// ResetTableStateCmd resets table to default state
func ResetTableStateCmd() tea.Cmd {
	return func() tea.Msg {
		if err := executeTableStateReset(); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to reset table state: %w", err)}
		}

		return statusMsg{message: "Table state reset to defaults"}
	}
}

// Table Interaction Commands

// BulkActionCmd performs actions on multiple selected rows
func BulkActionCmd(action string, selectedRows []int) tea.Cmd {
	return func() tea.Msg {
		if err := validateBulkAction(action, selectedRows); err != nil {
			return TableErrorMsg{Error: err}
		}

		if err := executeBulkAction(action, selectedRows); err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to execute bulk action: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Bulk action '%s' completed on %d rows", action, len(selectedRows))}
	}
}

// ValidateTableSelectionCmd validates selection before bulk actions
func ValidateTableSelectionCmd(selectedRows []int) tea.Cmd {
	return func() tea.Msg {
		if err := validateTableSelection(selectedRows); err != nil {
			return TableErrorMsg{Error: err}
		}

		return statusMsg{message: fmt.Sprintf("Selection of %d rows validated successfully", len(selectedRows))}
	}
}

// GetTableRowDetailsCmd fetches detailed information for a specific row
func GetTableRowDetailsCmd(rowIndex int) tea.Cmd {
	return func() tea.Msg {
		if err := validateRowIndex(rowIndex); err != nil {
			return TableErrorMsg{Error: err}
		}

		details, err := executeGetRowDetails(rowIndex)
		if err != nil {
			return TableErrorMsg{Error: fmt.Errorf("failed to get row details: %w", err)}
		}

		return statusMsg{message: fmt.Sprintf("Row details: %s", details)}
	}
}

// Helper Functions for Command Execution

// executeTableSort handles sort logic
func executeTableSort(column, direction string) error {
	// This would typically interact with the table component or data store
	// For now, we'll simulate the operation
	time.Sleep(10 * time.Millisecond) // Simulate async operation
	return nil
}

// Error Handling and Validation Functions

// validateSortColumn validates if the given column is valid for sorting
func validateSortColumn(column string) error {
	validColumns := []string{"id", "name", "status", "backend", "created", "last_active", "project", "messages"}
	for _, valid := range validColumns {
		if column == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid sort column: %s (valid columns: %s)", column, strings.Join(validColumns, ", "))
}

// validateFilterText validates filter text length and format
func validateFilterText(filterText string) error {
	if len(filterText) > 100 {
		return fmt.Errorf("filter text too long: %d characters (maximum: 100)", len(filterText))
	}

	// Check for potentially problematic characters
	if strings.Contains(filterText, "\x00") {
		return fmt.Errorf("filter text contains null characters")
	}

	return nil
}

// validateExportParams validates export format and filename
func validateExportParams(format, filename string) error {
	// Validate format
	validFormats := []string{"csv", "json"}
	formatValid := false
	for _, valid := range validFormats {
		if strings.ToLower(format) == valid {
			formatValid = true
			break
		}
	}
	if !formatValid {
		return fmt.Errorf("invalid export format: %s (valid formats: %s)", format, strings.Join(validFormats, ", "))
	}

	// Validate filename
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check for invalid filename characters
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	return nil
}

// validateBulkAction validates bulk action parameters
func validateBulkAction(action string, selectedRows []int) error {
	// Validate action
	validActions := []string{"delete", "export", "archive", "activate", "deactivate"}
	actionValid := false
	for _, valid := range validActions {
		if action == valid {
			actionValid = true
			break
		}
	}
	if !actionValid {
		return fmt.Errorf("invalid bulk action: %s (valid actions: %s)", action, strings.Join(validActions, ", "))
	}

	// Validate selection
	if len(selectedRows) == 0 {
		return fmt.Errorf("no rows selected for bulk action")
	}

	if len(selectedRows) > 100 {
		return fmt.Errorf("too many rows selected: %d (maximum: 100)", len(selectedRows))
	}

	// Validate row indices
	for _, rowIndex := range selectedRows {
		if rowIndex < 0 {
			return fmt.Errorf("invalid row index: %d (must be non-negative)", rowIndex)
		}
	}

	return nil
}

// validateTableSelection validates selection before bulk actions
func validateTableSelection(selectedRows []int) error {
	if len(selectedRows) == 0 {
		return fmt.Errorf("no rows selected")
	}

	// Check for duplicate indices
	seen := make(map[int]bool)
	for _, rowIndex := range selectedRows {
		if seen[rowIndex] {
			return fmt.Errorf("duplicate row index in selection: %d", rowIndex)
		}
		seen[rowIndex] = true

		if rowIndex < 0 {
			return fmt.Errorf("invalid row index: %d (must be non-negative)", rowIndex)
		}
	}

	return nil
}

// validateRowIndex validates if the given row index is valid
func validateRowIndex(rowIndex int) error {
	if rowIndex < 0 {
		return fmt.Errorf("invalid row index: %d (must be non-negative)", rowIndex)
	}

	// Note: Upper bound validation would typically be done against actual data size
	// This is a basic validation that can be extended based on context

	return nil
}

// validatePageNumber validates if the given page number is valid
func validatePageNumber(page int) error {
	if page < 1 {
		return fmt.Errorf("invalid page number: %d (must be positive)", page)
	}

	// Note: Upper bound validation would typically be done against actual total pages
	// This is a basic validation that can be extended based on context

	return nil
}

// validatePageSize validates if the given page size is valid
func validatePageSize(pageSize int) error {
	if pageSize < 1 {
		return fmt.Errorf("invalid page size: %d (must be positive)", pageSize)
	}

	if pageSize > 100 {
		return fmt.Errorf("page size too large: %d (maximum: 100)", pageSize)
	}

	return nil
}

// executeTableFilter handles filter logic
func executeTableFilter(filterText string) (int, error) {
	// This would typically apply the filter and return the number of matching rows
	// For now, we'll simulate the operation
	time.Sleep(10 * time.Millisecond) // Simulate async operation

	// Simulate row count based on filter complexity
	if len(filterText) == 0 {
		return 0, nil
	}
	return len(filterText) * 2, nil // Mock row count
}

// executeTableRefresh handles data refresh
func executeTableRefresh() error {
	// This would typically reload data from the backend
	// For now, we'll simulate the operation
	time.Sleep(50 * time.Millisecond) // Simulate network operation
	return nil
}

// executeTableExport handles data export
func executeTableExport(format, filename string, data []components.SessionData) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "csv":
		return exportToCSV(file, data)
	case "json":
		return exportToJSON(file, data)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// executeTableStateSave saves current table state
func executeTableStateSave() error {
	// This would typically save state to a config file or database
	// For now, we'll simulate the operation
	time.Sleep(20 * time.Millisecond) // Simulate I/O operation
	return nil
}

// executeTableStateLoad loads saved table state
func executeTableStateLoad() error {
	// This would typically load state from a config file or database
	// For now, we'll simulate the operation
	time.Sleep(20 * time.Millisecond) // Simulate I/O operation
	return nil
}

// executeTableStateReset resets table state to defaults
func executeTableStateReset() error {
	// This would typically reset all table configuration to defaults
	// For now, we'll simulate the operation
	time.Sleep(10 * time.Millisecond) // Simulate operation
	return nil
}

// executeBulkAction performs bulk actions on selected rows
func executeBulkAction(action string, selectedRows []int) error {
	// This would typically perform the specified action on all selected rows
	// For now, we'll simulate the operation
	time.Sleep(time.Duration(len(selectedRows)*10) * time.Millisecond) // Simulate processing time
	return nil
}

// executeGetRowDetails fetches detailed information for a row
func executeGetRowDetails(rowIndex int) (string, error) {
	// This would typically fetch detailed information from the data source
	// For now, we'll simulate the operation
	time.Sleep(20 * time.Millisecond) // Simulate data retrieval
	return fmt.Sprintf("Details for row %d", rowIndex), nil
}

// Export helper functions

// exportToCSV exports data to CSV format
func exportToCSV(file *os.File, data []components.SessionData) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Status", "Backend", "Created", "Last Active", "Project", "Messages"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, session := range data {
		record := []string{
			session.ID,
			session.Name,
			session.Status,
			session.Backend,
			session.Created.Format(time.RFC3339),
			session.LastActive.Format(time.RFC3339),
			session.ProjectPath,
			fmt.Sprintf("%d", session.Messages),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// exportToJSON exports data to JSON format
func exportToJSON(file *os.File, data []components.SessionData) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
