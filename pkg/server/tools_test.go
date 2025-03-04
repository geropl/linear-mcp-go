package server

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
)

var record = flag.Bool("record", false, "Record HTTP interactions (excluding writes)")
var recordWrites = flag.Bool("recordWrites", false, "Record HTTP interactions (incl. writes)")
var golden = flag.Bool("golden", false, "Update all golden files and recordings")

const (
	TEAM_NAME = "Test Team"
	TEAM_ID   = "234c5451-a839-4c8f-98d9-da00973f1060"
	ISSUE_ID  = "TEST-10"
	// COMMENT ISSUE ID is used for testing the add_comment handler
	COMMENT_ISSUE_ID = "TEST-12"
	USER_ID          = "cc24eee4-9edc-4bfe-b91b-fedde125ba85"
)

// expectation defines the expected output and error for a test case
type expectation struct {
	Err    string `json:"err"`    // Empty string means no error expected
	Output string `json:"output"` // Expected complete output
}

func TestHandlers(t *testing.T) {
	// Define test cases
	tests := []struct {
		handler string
		name    string
		args    map[string]interface{}
		write   bool
	}{
		// GetTeamsHandler test cases
		{
			handler: "get_teams",
			name:    "Get Teams",
			args: map[string]interface{}{
				"name": TEAM_NAME,
			},
		},
		// CreateIssueHandler test cases
		{
			handler: "create_issue",
			name:    "Valid issue",
			args: map[string]interface{}{
				"title":  "Test Issue",
				"teamId": TEAM_ID,
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Missing title",
			args: map[string]interface{}{
				"teamId": TEAM_ID,
			},
		},
		{
			handler: "create_issue",
			name:    "Missing teamId",
			args: map[string]interface{}{
				"title": "Test Issue",
			},
		},

		// UpdateIssueHandler test cases
		{
			handler: "update_issue",
			name:    "Valid update",
			args: map[string]interface{}{
				"id":    ISSUE_ID,
				"title": "Updated Test Issue",
			},
			write: true,
		},
		{
			handler: "update_issue",
			name:    "Missing id",
			args: map[string]interface{}{
				"title": "Updated Test Issue",
			},
		},

		// SearchIssuesHandler test cases
		{
			handler: "search_issues",
			name:    "Search by team",
			args: map[string]interface{}{
				"teamId": TEAM_ID,
				"limit":  float64(5),
			},
		},
		{
			handler: "search_issues",
			name:    "Search by query",
			args: map[string]interface{}{
				"query": "test",
				"limit": float64(5),
			},
		},

		// GetUserIssuesHandler test cases
		{
			handler: "get_user_issues",
			name:    "Current user issues",
			args: map[string]interface{}{
				"limit": float64(5),
			},
		},
		{
			handler: "get_user_issues",
			name:    "Specific user issues",
			args: map[string]interface{}{
				"userId": USER_ID,
				"limit":  float64(5),
			},
		},

		// GetIssueHandler test cases
		{
			handler: "get_issue",
			name:    "Valid issue",
			args: map[string]interface{}{
				"issueId": ISSUE_ID,
			},
		},
		{
			handler: "get_issue",
			name:    "Get comment issue",
			args: map[string]interface{}{
				"issueId": COMMENT_ISSUE_ID,
			},
		},
		{
			handler: "get_issue",
			name:    "Missing issueId",
			args:    map[string]interface{}{},
		},

		// AddCommentHandler test cases
		{
			handler: "add_comment",
			name:    "Valid comment",
			write:   true,
			args: map[string]interface{}{
				"issueId": ISSUE_ID,
				"body":    "Test comment",
			},
		},
		{
			handler: "add_comment",
			name:    "Missing issueId",
			args: map[string]interface{}{
				"body": "Test comment",
			},
		},
		{
			handler: "add_comment",
			name:    "Missing body",
			args: map[string]interface{}{
				"issueId": ISSUE_ID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.handler+"_"+tt.name, func(t *testing.T) {
			if tt.write && *record && !*recordWrites {
				t.Skip("Skipping write test when recordWrites=false")
				return
			}

			// Generate golden file path
			goldenPath := filepath.Join("../../testdata/golden", tt.handler+"_handler_"+tt.name+".golden")

			// Create test client
			client, cleanup := linear.NewTestClient(t, tt.handler+"_handler_"+tt.name, *record || *recordWrites)
			defer cleanup()

			// Create the appropriate handler based on tt.handler
			var handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
			switch tt.handler {
			case "get_teams":
				handler = GetTeamsHandler(client)
			case "create_issue":
				handler = CreateIssueHandler(client)
			case "update_issue":
				handler = UpdateIssueHandler(client)
			case "search_issues":
				handler = SearchIssuesHandler(client)
			case "get_user_issues":
				handler = GetUserIssuesHandler(client)
			case "get_issue":
				handler = GetIssueHandler(client)
			case "add_comment":
				handler = AddCommentHandler(client)
			default:
				t.Fatalf("Unknown handler type: %s", tt.handler)
			}

			// Create the request
			request := mcp.CallToolRequest{}
			request.Params.Name = "linear_" + tt.handler
			request.Params.Arguments = tt.args

			// Call the handler
			result, err := handler(context.Background(), request)

			// Check for errors
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			// Extract the actual output and error
			var actualOutput, actualErr string
			if result.IsError {
				// For error results, the error message is in the text content
				for _, content := range result.Content {
					if textContent, ok := content.(mcp.TextContent); ok {
						actualErr = textContent.Text
						break
					}
				}
			} else {
				// For success results, the output is in the text content
				for _, content := range result.Content {
					if textContent, ok := content.(mcp.TextContent); ok {
						actualOutput = textContent.Text
						break
					}
				}
			}

			// If golden flag is set, update the golden file
			if *golden {
				writeGoldenFile(t, goldenPath, expectation{
					Err:    actualErr,
					Output: actualOutput,
				})
				return
			}

			// Otherwise, read the golden file and compare
			expected := readGoldenFile(t, goldenPath)

			// Compare error
			if diff := cmp.Diff(expected.Err, actualErr); diff != "" {
				t.Errorf("Error mismatch (-want +got):\n%s", diff)
			}

			// Compare output (only if no error is expected)
			if expected.Err == "" {
				if diff := cmp.Diff(expected.Output, actualOutput); diff != "" {
					t.Errorf("Output mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// readGoldenFile reads an expectation from a golden file
func readGoldenFile(t *testing.T, path string) expectation {
	t.Helper()

	// Check if the golden file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Golden file %s does not exist", path)
	}

	// Read the golden file
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", path, err)
	}

	// Parse the golden file
	var exp expectation
	if err := json.Unmarshal(data, &exp); err != nil {
		t.Fatalf("Failed to parse golden file %s: %v", path, err)
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

	// Marshal the expectation
	data, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal expectation: %v", err)
	}

	// Write the golden file
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("Failed to write golden file %s: %v", path, err)
	}
}
