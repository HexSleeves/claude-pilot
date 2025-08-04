package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// GoldenTestCase represents a single golden file test case
type GoldenTestCase struct {
	Name       string   // Test case name
	Args       []string // Command arguments
	GoldenFile string   // Path to golden file relative to testdata/golden/
	WantCode   int      // Expected exit code
	Setup      func() error // Optional setup function
	Cleanup    func() error // Optional cleanup function
}

// GoldenTestSuite manages and runs golden file tests
type GoldenTestSuite struct {
	BinaryPath string // Path to the CLI binary
	TestdataDir string // Path to testdata directory
}

// NewGoldenTestSuite creates a new golden test suite
func NewGoldenTestSuite(binaryPath, testdataDir string) *GoldenTestSuite {
	return &GoldenTestSuite{
		BinaryPath:  binaryPath,
		TestdataDir: testdataDir,
	}
}

// RunTests executes all provided test cases
func (gts *GoldenTestSuite) RunTests(t *testing.T, tests []GoldenTestCase) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gts.runSingleTest(t, tt)
		})
	}
}

// runSingleTest executes a single golden test case
func (gts *GoldenTestSuite) runSingleTest(t *testing.T, tc GoldenTestCase) {
	// Setup if provided
	if tc.Setup != nil {
		if err := tc.Setup(); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
	}

	// Cleanup if provided
	if tc.Cleanup != nil {
		defer func() {
			if err := tc.Cleanup(); err != nil {
				t.Errorf("Cleanup failed: %v", err)
			}
		}()
	}

	// Execute command
	output, exitCode := gts.executeCommand(tc.Args)

	// Check exit code
	if exitCode != tc.WantCode {
		t.Errorf("Expected exit code %d, got %d", tc.WantCode, exitCode)
	}

	// Load golden file and compare
	gts.assertGoldenOutput(t, output, tc.GoldenFile)
}

// executeCommand runs the CLI binary with given arguments
func (gts *GoldenTestSuite) executeCommand(args []string) (string, int) {
	cmd := exec.Command(gts.BinaryPath, args...)
	
	// Set environment to ensure consistent output
	cmd.Env = append(os.Environ(),
		"NO_COLOR=1",           // Disable colors for consistent output
		"TERM=xterm-256color",  // Consistent terminal type
		"COLUMNS=100",          // Consistent width
		"LINES=24",             // Consistent height
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// Combine stdout and stderr for complete output
	output := stdout.String()
	if stderr.Len() > 0 {
		if len(output) > 0 {
			output += "\n"
		}
		output += stderr.String()
	}

	// Filter out common terminal errors that appear in CI environments
	output = FilterOutputErrors(output)

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return output, exitCode
}

// assertGoldenOutput compares actual output with golden file
func (gts *GoldenTestSuite) assertGoldenOutput(t *testing.T, actual, goldenFile string) {
	goldenPath := filepath.Join(gts.TestdataDir, "golden", goldenFile)
	
	// Read golden file
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", goldenPath, err)
	}

	// Normalize line endings and trim whitespace
	expectedStr := normalizeOutput(string(expected))
	actualStr := normalizeOutput(actual)

	// Special handling for JSON files
	if strings.HasSuffix(goldenFile, ".json") {
		gts.assertJSONOutput(t, actualStr, expectedStr, goldenFile)
		return
	}

	// Line-by-line comparison for better error messages
	expectedLines := strings.Split(expectedStr, "\n")
	actualLines := strings.Split(actualStr, "\n")

	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	for i := 0; i < maxLines; i++ {
		var expectedLine, actualLine string
		
		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}

		if expectedLine != actualLine {
			t.Errorf("Output mismatch at line %d in %s:\nExpected: %q\nActual:   %q\n\nFull diff:\nExpected:\n%s\n\nActual:\n%s",
				i+1, goldenFile, expectedLine, actualLine, expectedStr, actualStr)
			return
		}
	}
}

// assertJSONOutput compares JSON output with proper formatting
func (gts *GoldenTestSuite) assertJSONOutput(t *testing.T, actual, expected, goldenFile string) {
	// Parse both JSON strings
	var actualJSON, expectedJSON interface{}
	
	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		t.Fatalf("Actual output is not valid JSON in %s: %v\nOutput:\n%s", goldenFile, err, actual)
	}
	
	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		t.Fatalf("Golden file contains invalid JSON in %s: %v", goldenFile, err)
	}

	// Pretty print both for comparison
	actualPretty, _ := json.MarshalIndent(actualJSON, "", "  ")
	expectedPretty, _ := json.MarshalIndent(expectedJSON, "", "  ")

	if string(actualPretty) != string(expectedPretty) {
		t.Errorf("JSON output mismatch in %s:\nExpected:\n%s\n\nActual:\n%s",
			goldenFile, string(expectedPretty), string(actualPretty))
	}
}

// normalizeOutput normalizes output for consistent comparison
func normalizeOutput(output string) string {
	// Split into lines for processing
	lines := strings.Split(output, "\n")
	var normalized []string

	for _, line := range lines {
		// Trim trailing whitespace but preserve leading whitespace for formatting
		line = strings.TrimRight(line, " \t")
		normalized = append(normalized, line)
	}

	// Join back and trim final whitespace
	result := strings.Join(normalized, "\n")
	result = strings.TrimRight(result, "\n")
	
	return result
}

// UpdateGoldenFile updates a golden file with new content (for maintenance)
func (gts *GoldenTestSuite) UpdateGoldenFile(goldenFile, content string) error {
	goldenPath := filepath.Join(gts.TestdataDir, "golden", goldenFile)
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write normalized content
	normalizedContent := normalizeOutput(content)
	if err := os.WriteFile(goldenPath, []byte(normalizedContent), 0644); err != nil {
		return fmt.Errorf("failed to write golden file: %v", err)
	}

	return nil
}

// LoadGoldenFile loads content from a golden file
func (gts *GoldenTestSuite) LoadGoldenFile(goldenFile string) (string, error) {
	goldenPath := filepath.Join(gts.TestdataDir, "golden", goldenFile)
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read golden file %s: %v", goldenPath, err)
	}
	return string(content), nil
}

// ValidateGoldenFiles validates that all golden files exist and are readable
func (gts *GoldenTestSuite) ValidateGoldenFiles(t *testing.T, goldenFiles []string) {
	for _, goldenFile := range goldenFiles {
		goldenPath := filepath.Join(gts.TestdataDir, "golden", goldenFile)
		if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
			t.Errorf("Golden file does not exist: %s", goldenPath)
		} else if err != nil {
			t.Errorf("Error accessing golden file %s: %v", goldenPath, err)
		}
	}
}

// DiscoverGoldenFiles discovers all golden files in a directory
func (gts *GoldenTestSuite) DiscoverGoldenFiles(subdir string) ([]string, error) {
	goldenDir := filepath.Join(gts.TestdataDir, "golden", subdir)
	
	var files []string
	err := filepath.Walk(goldenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			relPath, err := filepath.Rel(filepath.Join(gts.TestdataDir, "golden"), path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		
		return nil
	})
	
	return files, err
}

// FilterOutputErrors filters out known terminal/environment errors from output
func FilterOutputErrors(output string) string {
	lines := strings.Split(output, "\n")
	var filtered []string
	
	for _, line := range lines {
		// Skip common terminal environment errors that appear in CI
		if strings.Contains(line, "Error getting terminal size: operation not supported on socket") ||
		   strings.Contains(line, "inappropriate ioctl for device") ||
		   strings.Contains(line, "not a terminal") {
			continue
		}
		filtered = append(filtered, line)
	}
	
	return strings.Join(filtered, "\n")
}

// BuildBinary builds the CLI binary for testing
func BuildBinary(projectRoot, outputPath string) error {
	cmd := exec.Command("make", "build")
	cmd.Dir = projectRoot
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build binary: %v\nOutput: %s", err, stderr.String())
	}
	
	return nil
}