package server

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var record = flag.Bool("record", false, "Record HTTP interactions (excluding writes)")
var recordWrites = flag.Bool("recordWrites", false, "Record HTTP interactions (incl. writes)")
var golden = flag.Bool("golden", false, "Update all golden files and recordings")

// Shared constants for tests
const (
	TEAM_NAME        = "Test Team"
	TEAM_KEY         = "TEST"
	TEAM_ID          = "234c5451-a839-4c8f-98d9-da00973f1060"
	ISSUE_ID         = "TEST-10"
	COMMENT_ISSUE_ID = "TEST-12" // Used for testing add_comment handler
	USER_ID          = "cc24eee4-9edc-4bfe-b91b-fedde125ba85"
)

// expectation defines the expected output and error for a test case
// For resource tests, Output will store the JSON representation of []mcp.ResourceContents
type expectation struct {
	Err    string `yaml:"err"`          // Empty string means no error expected
	Output string `yaml:"output", flow` // Expected complete output
}

// readGoldenFile reads an expectation from a golden file
func readGoldenFile(t *testing.T, path string) expectation {
	t.Helper()

	// Check if the golden file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the file doesn't exist, return an empty expectation
		// This allows tests to pass initially when golden files are missing,
		// prompting the user to run with -golden* flags to create them.
		t.Logf("Golden file %s does not exist. Run with appropriate -golden* flag to create it.", path)
		return expectation{}
	}

	// Read the golden file
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", path, err)
	}

	// Parse the golden file
	var exp expectation
	if err := yaml.Unmarshal(data, &exp); err != nil {
		// If unmarshalling fails, treat it as an empty expectation
		// This handles cases where the golden file might be corrupted or empty
		t.Logf("Failed to parse golden file %s: %v. Treating as empty.", path, err)
		return expectation{}
	}

	return exp
}

// writeGoldenFile writes an expectation to a golden file
func writeGoldenFile(t *testing.T, path string, exp expectation) {
	t.Helper()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Marshal the YAML node
	data, err := yaml.Marshal(&exp)
	if err != nil {
		t.Fatalf("Failed to marshal expectation: %v", err)
	}

	// Write the golden file
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("Failed to write golden file %s: %v", path, err)
	}
	t.Logf("Successfully wrote golden file: %s", path)
}
