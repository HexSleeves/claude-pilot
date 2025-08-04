package main

import (
	"os"
	"path/filepath"
	"testing"

	goldentest "claude-pilot/internal/testing"
)

const (
	binaryName = "claude-pilot"
)

var (
	suite *goldentest.GoldenTestSuite
)

func TestMain(m *testing.M) {
	// Build binary for testing
	projectRoot, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	
	// Navigate to the project root (2 levels up from packages/claudepilot)
	projectRoot = filepath.Join(projectRoot, "..", "..")
	testdataDir := filepath.Join(projectRoot, "testdata")
	binaryPath := filepath.Join(projectRoot, binaryName)
	
	// Build the binary
	if err := goldentest.BuildBinary(projectRoot, binaryPath); err != nil {
		panic(err)
	}
	
	// Initialize test suite
	suite = goldentest.NewGoldenTestSuite(binaryPath, testdataDir)
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	os.Remove(binaryPath)
	os.Exit(code)
}

// TestHelpOutputs tests all help command variations
func TestHelpOutputs(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "root help",
			Args:       []string{"--help"},
			GoldenFile: "help/help-root.txt",
			WantCode:   0,
		},
		{
			Name:       "create help",
			Args:       []string{"create", "--help"},
			GoldenFile: "help/help-create.txt",
			WantCode:   0,
		},
		{
			Name:       "list help",
			Args:       []string{"list", "--help"},
			GoldenFile: "help/help-list.txt",
			WantCode:   0,
		},
		{
			Name:       "details help",
			Args:       []string{"details", "--help"},
			GoldenFile: "help/help-details.txt",
			WantCode:   0,
		},
		{
			Name:       "kill help",
			Args:       []string{"kill", "--help"},
			GoldenFile: "help/help-kill.txt",
			WantCode:   0,
		},
		{
			Name:       "attach help",
			Args:       []string{"attach", "--help"},
			GoldenFile: "help/help-attach.txt",
			WantCode:   0,
		},
		{
			Name:       "version output",
			Args:       []string{"--version"},
			GoldenFile: "help/help-version.txt",
			WantCode:   0,
		},
	}

	suite.RunTests(t, tests)
}

// TestListOutputFormats tests list command with different output formats
func TestListOutputFormats(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "list human format with sessions",
			Args:       []string{"list", "--output", "human"},
			GoldenFile: "output/list-human.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "list table format",
			Args:       []string{"list", "--output", "table"},
			GoldenFile: "output/list-table.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "list json format",
			Args:       []string{"list", "--output", "json"},
			GoldenFile: "output/list-json.json",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "list ndjson format",
			Args:       []string{"list", "--output", "ndjson"},
			GoldenFile: "output/list-ndjson.ndjson",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "list quiet format",
			Args:       []string{"list", "--output", "quiet"},
			GoldenFile: "output/list-quiet.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "list empty human format",
			Args:       []string{"list", "--output", "human"},
			GoldenFile: "output/list-empty-human.txt",
			WantCode:   0,
			// No setup - should return empty results
		},
		{
			Name:       "list empty json format",
			Args:       []string{"list", "--output", "json"},
			GoldenFile: "output/list-empty-json.json",
			WantCode:   0,
			// No setup - should return empty results
		},
	}

	suite.RunTests(t, tests)
}

// TestCreateCommand tests create command scenarios
func TestCreateCommand(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "create success human format",
			Args:       []string{"create", "test-session", "--output", "human"},
			GoldenFile: "output/create-success-human.txt",
			WantCode:   0,
			Cleanup:    cleanupTestSession,
		},
		{
			Name:       "create success json format",
			Args:       []string{"create", "test-session", "--output", "json"},
			GoldenFile: "output/create-success-json.json",
			WantCode:   0,
			Cleanup:    cleanupTestSession,
		},
		{
			Name:       "create success quiet format",
			Args:       []string{"create", "test-session", "--output", "quiet"},
			GoldenFile: "output/create-success-quiet.txt",
			WantCode:   0,
			Cleanup:    cleanupTestSession,
		},
	}

	suite.RunTests(t, tests)
}

// TestDetailsCommand tests details command scenarios
func TestDetailsCommand(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "details success human format",
			Args:       []string{"details", "--id", "d7e8f9a0-1234-5678-9abc-def012345678", "--output", "human"},
			GoldenFile: "output/details-success-human.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
		{
			Name:       "details success json format",
			Args:       []string{"details", "--id", "d7e8f9a0-1234-5678-9abc-def012345678", "--output", "json"},
			GoldenFile: "output/details-success-json.json",
			WantCode:   0,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
	}

	suite.RunTests(t, tests)
}

// TestKillCommand tests kill command scenarios
func TestKillCommand(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "kill success human format",
			Args:       []string{"kill", "example-sess", "--output", "human"},
			GoldenFile: "output/kill-success-human.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
		},
		{
			Name:       "kill success json format",
			Args:       []string{"kill", "example-sess", "--output", "json"},
			GoldenFile: "output/kill-success-json.json",
			WantCode:   0,
			Setup:      setupMockSessions,
		},
		{
			Name:       "kill all success",
			Args:       []string{"kill", "--all", "--output", "human"},
			GoldenFile: "output/kill-all-success-human.txt",
			WantCode:   0,
			Setup:      setupMockSessions,
		},
		{
			Name:       "kill partial failure",
			Args:       []string{"kill", "nonexistent-session", "--output", "human"},
			GoldenFile: "output/kill-partial-failure-human.txt",
			WantCode:   1,
			Setup:      setupMockSessions,
			Cleanup:    cleanupMockSessions,
		},
	}

	suite.RunTests(t, tests)
}

// TestErrorScenarios tests various error conditions
func TestErrorScenarios(t *testing.T) {
	tests := []goldentest.GoldenTestCase{
		{
			Name:       "session not found error",
			Args:       []string{"details", "--id", "nonexistent-session"},
			GoldenFile: "error/error-session-not-found.txt",
			WantCode:   3,
		},
		{
			Name:       "session not found error json",
			Args:       []string{"details", "--id", "nonexistent-session", "--output", "json"},
			GoldenFile: "error/error-session-not-found-json.json",
			WantCode:   3,
		},
		{
			Name:       "validation error",
			Args:       []string{"create", "", "--output", "human"},
			GoldenFile: "error/error-validation.txt",
			WantCode:   2,
		},
		{
			Name:       "validation error json",
			Args:       []string{"create", "", "--output", "json"},
			GoldenFile: "error/error-validation-json.json",
			WantCode:   2,
		},
	}

	suite.RunTests(t, tests)
}

// Helper functions for test setup and cleanup

func setupMockSessions() error {
	// Create mock session data for consistent testing
	// This would integrate with your session storage system
	// Implementation depends on your storage backend
	return nil
}

func cleanupMockSessions() error {
	// Clean up mock session data
	return nil
}

func cleanupTestSession() error {
	// Clean up the test session created during create command test
	return nil
}

// TestGoldenFileIntegrity validates that all golden files are present and valid
func TestGoldenFileIntegrity(t *testing.T) {
	expectedFiles := []string{
		// Help files
		"help/help-root.txt",
		"help/help-create.txt",
		"help/help-list.txt",
		"help/help-details.txt",
		"help/help-kill.txt",
		"help/help-attach.txt",
		"help/help-version.txt",
		
		// Output files
		"output/create-success-human.txt",
		"output/create-success-json.json",
		"output/create-success-quiet.txt",
		"output/details-success-human.txt",
		"output/details-success-json.json",
		"output/kill-all-success-human.txt",
		"output/kill-partial-failure-human.txt",
		"output/kill-success-human.txt",
		"output/kill-success-json.json",
		"output/list-empty-human.txt",
		"output/list-empty-json.json",
		"output/list-human.txt",
		"output/list-json.json",
		"output/list-ndjson.ndjson",
		"output/list-quiet.txt",
		"output/list-table.txt",
		
		// Error files
		"error/error-network-json.json",
		"error/error-network.txt",
		"error/error-permission-denied-json.json",
		"error/error-permission-denied.txt",
		"error/error-session-not-found-json.json",
		"error/error-session-not-found.txt",
		"error/error-validation-json.json",
		"error/error-validation.txt",
	}
	
	suite.ValidateGoldenFiles(t, expectedFiles)
}

// TestDiscoverAllGoldenFiles tests the golden file discovery functionality
func TestDiscoverAllGoldenFiles(t *testing.T) {
	categories := []string{"help", "output", "error"}
	
	for _, category := range categories {
		t.Run(category, func(t *testing.T) {
			files, err := suite.DiscoverGoldenFiles(category)
			if err != nil {
				t.Errorf("Failed to discover golden files in %s: %v", category, err)
			}
			
			if len(files) == 0 {
				t.Errorf("No golden files found in %s directory", category)
			}
			
			t.Logf("Found %d golden files in %s: %v", len(files), category, files)
		})
	}
}