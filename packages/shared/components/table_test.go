package components

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func TestNewTable(t *testing.T) {
	config := TableConfig{
		Width:       80,
		ShowHeaders: true,
		Interactive: true,
		MaxRows:     10,
	}

	tbl := NewTable(config)

	if tbl.config.Width != 80 {
		t.Errorf("Expected width 80, got %d", tbl.config.Width)
	}

	if !tbl.config.ShowHeaders {
		t.Error("Expected ShowHeaders to be true")
	}

	if !tbl.config.Interactive {
		t.Error("Expected Interactive to be true")
	}

	if tbl.config.MaxRows != 10 {
		t.Errorf("Expected MaxRows 10, got %d", tbl.config.MaxRows)
	}
}

func TestSetData(t *testing.T) {
	tbl := NewTable(TableConfig{})

	data := TableData{
		Headers: []string{"ID", "Name", "Status"},
		Rows: [][]string{
			{"1", "session1", "active"},
			{"2", "session2", "inactive"},
		},
	}

	tbl.SetData(data)

	if len(tbl.data.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(tbl.data.Headers))
	}

	if len(tbl.data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(tbl.data.Rows))
	}

	if tbl.data.Headers[0] != "ID" {
		t.Errorf("Expected first header to be 'ID', got '%s'", tbl.data.Headers[0])
	}

	if tbl.data.Rows[0][0] != "1" {
		t.Errorf("Expected first row first cell to be '1', got '%s'", tbl.data.Rows[0][0])
	}
}

func TestSetSessionData(t *testing.T) {
	tbl := NewTable(TableConfig{})

	sessions := []SessionData{
		{
			ID:          "test-id-1",
			Name:        "test-session-1",
			Status:      "active",
			Backend:     "tmux",
			Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			LastActive:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			Messages:    5,
			ProjectPath: "/test/project",
		},
		{
			ID:          "test-id-2",
			Name:        "test-session-2",
			Status:      "inactive",
			Backend:     "zellij",
			Created:     time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			LastActive:  time.Date(2023, 1, 2, 12, 30, 0, 0, time.UTC),
			Messages:    3,
			ProjectPath: "/test/project2",
		},
	}

	tbl.SetSessionData(sessions)

	if len(tbl.data.Headers) != 8 {
		t.Errorf("Expected 8 headers, got %d", len(tbl.data.Headers))
	}

	if len(tbl.data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(tbl.data.Rows))
	}

	expectedHeaders := []string{"ID", "Name", "Status", "Backend", "Created", "Last Active", "Messages", "Project"}
	for i, header := range expectedHeaders {
		if tbl.data.Headers[i] != header {
			t.Errorf("Expected header %d to be '%s', got '%s'", i, header, tbl.data.Headers[i])
		}
	}

	// Check first row data
	if tbl.data.Rows[0][0] != "test-id-1" {
		t.Errorf("Expected first row ID to be 'test-id-1', got '%s'", tbl.data.Rows[0][0])
	}

	if tbl.data.Rows[0][1] != "test-session-1" {
		t.Errorf("Expected first row name to be 'test-session-1', got '%s'", tbl.data.Rows[0][1])
	}

	if tbl.data.Rows[0][2] != "active" {
		t.Errorf("Expected first row status to be 'active', got '%s'", tbl.data.Rows[0][2])
	}
}

func TestToBubblesColumns(t *testing.T) {
	tbl := NewTable(TableConfig{Width: 100})

	data := TableData{
		Headers: []string{"ID", "Name", "Status"},
		Rows:    [][]string{},
	}

	tbl.SetData(data)

	columns := tbl.ToBubblesColumns()

	if len(columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(columns))
	}

	if columns[0].Title != "ID" {
		t.Errorf("Expected first column title to be 'ID', got '%s'", columns[0].Title)
	}

	if columns[1].Title != "Name" {
		t.Errorf("Expected second column title to be 'Name', got '%s'", columns[1].Title)
	}

	if columns[2].Title != "Status" {
		t.Errorf("Expected third column title to be 'Status', got '%s'", columns[2].Title)
	}

	// Check that columns have widths
	for i, col := range columns {
		if col.Width <= 0 {
			t.Errorf("Column %d width should be > 0, got %d", i, col.Width)
		}
	}
}

func TestToBubblesRows(t *testing.T) {
	tbl := NewTable(TableConfig{})

	data := TableData{
		Headers: []string{"ID", "Name", "Status"},
		Rows: [][]string{
			{"1", "session1", "active"},
			{"2", "session2", "inactive"},
		},
	}

	tbl.SetData(data)

	rows := tbl.ToBubblesRows()

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	// Check first row
	if len(rows[0]) != 3 {
		t.Errorf("Expected first row to have 3 cells, got %d", len(rows[0]))
	}

	if rows[0][0] != "1" {
		t.Errorf("Expected first row first cell to be '1', got '%s'", rows[0][0])
	}

	if rows[0][1] != "session1" {
		t.Errorf("Expected first row second cell to be 'session1', got '%s'", rows[0][1])
	}

	if rows[0][2] != "active" {
		t.Errorf("Expected first row third cell to be 'active', got '%s'", rows[0][2])
	}

	// Check second row
	if len(rows[1]) != 3 {
		t.Errorf("Expected second row to have 3 cells, got %d", len(rows[1]))
	}

	if rows[1][0] != "2" {
		t.Errorf("Expected second row first cell to be '2', got '%s'", rows[1][0])
	}
}

func TestToBubblesSessionRows(t *testing.T) {
	sessions := []SessionData{
		{
			ID:          "test-id-1",
			Name:        "test-session-1",
			Status:      "active",
			Backend:     "tmux",
			Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			LastActive:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			Messages:    5,
			ProjectPath: "/test/project",
		},
	}

	rows := ToBubblesSessionRows(sessions)

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	if len(rows[0]) != 8 {
		t.Errorf("Expected row to have 8 cells, got %d", len(rows[0]))
	}

	// Check specific cells
	if rows[0][0] != "test-id-1" {
		t.Errorf("Expected ID cell to be 'test-id-1', got '%s'", rows[0][0])
	}

	if rows[0][1] != "test-session-1" {
		t.Errorf("Expected name cell to be 'test-session-1', got '%s'", rows[0][1])
	}

	if rows[0][2] != "active" {
		t.Errorf("Expected status cell to be 'active', got '%s'", rows[0][2])
	}

	if rows[0][3] != "tmux" {
		t.Errorf("Expected backend cell to be 'tmux', got '%s'", rows[0][3])
	}

	if rows[0][6] != "5" {
		t.Errorf("Expected messages cell to be '5', got '%s'", rows[0][6])
	}
}

func TestGetBubblesTableColumns(t *testing.T) {
	columns := GetBubblesTableColumns()

	if len(columns) != 8 {
		t.Errorf("Expected 8 columns, got %d", len(columns))
	}

	expectedTitles := []string{"ID", "Name", "Status", "Backend", "Created", "Last Active", "Messages", "Project"}
	for i, title := range expectedTitles {
		if columns[i].Title != title {
			t.Errorf("Expected column %d title to be '%s', got '%s'", i, title, columns[i].Title)
		}
	}

	// Check that all columns have positive widths
	for i, col := range columns {
		if col.Width <= 0 {
			t.Errorf("Column %d width should be > 0, got %d", i, col.Width)
		}
	}
}

func TestConfigureBubblesTable(t *testing.T) {
	tbl := NewTable(TableConfig{Width: 100, MaxRows: 20})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
		},
	}

	tbl.SetData(data)

	// Create a Bubbles table
	bubblesTable := table.New()

	// Configure it
	configuredTable := tbl.ConfigureBubblesTable(bubblesTable)

	// Check that the table was configured (this is a bit limited without accessing internals)
	if configuredTable.Width() != 100 {
		t.Errorf("Expected configured table width to be 100, got %d", configuredTable.Width())
	}

	if configuredTable.Height() != 20 {
		t.Errorf("Expected configured table height to be 20, got %d", configuredTable.Height())
	}
}

func TestSetSelectedRow(t *testing.T) {
	tbl := NewTable(TableConfig{Interactive: true})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
			{"3", "session3"},
		},
	}

	tbl.SetData(data)

	// Test setting valid row
	tbl.SetSelectedRow(1)
	if tbl.GetSelectedRow() != 1 {
		t.Errorf("Expected selected row to be 1, got %d", tbl.GetSelectedRow())
	}

	// Test setting invalid row (should be ignored)
	tbl.SetSelectedRow(10)
	if tbl.GetSelectedRow() != 1 {
		t.Errorf("Expected selected row to remain 1, got %d", tbl.GetSelectedRow())
	}

	// Test setting negative row (should be ignored)
	tbl.SetSelectedRow(-1)
	if tbl.GetSelectedRow() != 1 {
		t.Errorf("Expected selected row to remain 1, got %d", tbl.GetSelectedRow())
	}
}

func TestGetSelectedData(t *testing.T) {
	tbl := NewTable(TableConfig{Interactive: true})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
			{"3", "session3"},
		},
	}

	tbl.SetData(data)

	// Test getting selected data
	tbl.SetSelectedRow(1)
	selectedData := tbl.GetSelectedData()

	if len(selectedData) != 2 {
		t.Errorf("Expected selected data to have 2 cells, got %d", len(selectedData))
	}

	if selectedData[0] != "2" {
		t.Errorf("Expected first cell to be '2', got '%s'", selectedData[0])
	}

	if selectedData[1] != "session2" {
		t.Errorf("Expected second cell to be 'session2', got '%s'", selectedData[1])
	}
}

func TestGetRowCount(t *testing.T) {
	tbl := NewTable(TableConfig{})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
			{"3", "session3"},
		},
	}

	tbl.SetData(data)

	if tbl.GetRowCount() != 3 {
		t.Errorf("Expected row count to be 3, got %d", tbl.GetRowCount())
	}
}

func TestRenderCLI(t *testing.T) {
	tbl := NewTable(TableConfig{ShowHeaders: true, Width: 60})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
		},
	}

	tbl.SetData(data)

	output := tbl.RenderCLI()

	// Check that output contains expected content
	if output == "" {
		t.Error("Expected non-empty CLI output")
	}

	// Should contain headers
	if !contains(output, "ID") {
		t.Error("Expected output to contain 'ID' header")
	}

	if !contains(output, "Name") {
		t.Error("Expected output to contain 'Name' header")
	}

	// Should contain data
	if !contains(output, "session1") {
		t.Error("Expected output to contain 'session1'")
	}

	if !contains(output, "session2") {
		t.Error("Expected output to contain 'session2'")
	}
}

func TestRenderTUI(t *testing.T) {
	tbl := NewTable(TableConfig{ShowHeaders: true, Interactive: true, Width: 60})

	data := TableData{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "session1"},
			{"2", "session2"},
		},
	}

	tbl.SetData(data)

	output := tbl.RenderTUI()

	// Check that output contains expected content
	if output == "" {
		t.Error("Expected non-empty TUI output")
	}

	// Should contain data
	if !contains(output, "session1") {
		t.Error("Expected output to contain 'session1'")
	}

	if !contains(output, "session2") {
		t.Error("Expected output to contain 'session2'")
	}
}

func TestEmptyTable(t *testing.T) {
	tbl := NewTable(TableConfig{})

	// Test empty table CLI render
	cliOutput := tbl.RenderCLI()
	if !contains(cliOutput, "No data to display") {
		t.Error("Expected CLI output to show 'No data to display' for empty table")
	}

	// Test empty table TUI render
	tuiOutput := tbl.RenderTUI()
	if !contains(tuiOutput, "No data to display") {
		t.Error("Expected TUI output to show 'No data to display' for empty table")
	}

	// Test empty table Bubbles conversions
	columns := tbl.ToBubblesColumns()
	if columns != nil {
		t.Error("Expected nil columns for empty table")
	}

	rows := tbl.ToBubblesRows()
	if rows != nil {
		t.Error("Expected nil rows for empty table")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
