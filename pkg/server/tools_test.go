package server

import (
	"context"
	"flag"
	"testing"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
)

var record = flag.Bool("record", false, "Record HTTP interactions (excluding writes)")
var recordWrites = flag.Bool("recordWrites", false, "Record HTTP interactions (incl. writes)")

const (
	TEAM_NAME = "Test Team"
	TEAM_ID   = "234c5451-a839-4c8f-98d9-da00973f1060"
	ISSUE_ID  = "TEST-10"
	USER_ID   = "cc24eee4-9edc-4bfe-b91b-fedde125ba85" //a6ac53b5-d345-449f-bba2-e0b3b09ed2e3"
)

// expectation defines the expected output and error for a test case
type expectation struct {
	err    string // Empty string means no error expected
	output string // Expected complete output
}

func TestHandlers(t *testing.T) {
	// Define test cases
	tests := []struct {
		handler     string
		name        string
		args        map[string]interface{}
		write       bool
		expectation expectation
	}{
		// GetTeamsHandler test cases
		{
			handler: "get_teams",
			name:    "Get Teams",
			args: map[string]interface{}{
				"name": TEAM_NAME,
			},
			expectation: expectation{
				err:    "",
				output: "Found 1 teams:\n- Test Team (Key: TEST)\n  ID: 234c5451-a839-4c8f-98d9-da00973f1060\n",
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
			expectation: expectation{
				err:    "",
				output: "Created issue TEST-10: Test Issue\nURL: https://linear.app/linear-mcp-go-test/issue/TEST-10/test-issue",
			},
		},
		{
			handler: "create_issue",
			name:    "Missing title",
			args: map[string]interface{}{
				"teamId": TEAM_ID,
			},
			expectation: expectation{
				err:    "title must be a non-empty string",
				output: "",
			},
		},
		{
			handler: "create_issue",
			name:    "Missing teamId",
			args: map[string]interface{}{
				"title": "Test Issue",
			},
			expectation: expectation{
				err:    "teamId must be a non-empty string",
				output: "",
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
			expectation: expectation{
				err:    "",
				output: "Updated issue TEST-10\nURL: https://linear.app/linear-mcp-go-test/issue/TEST-10/updated-test-issue",
			},
		},
		{
			handler: "update_issue",
			name:    "Missing id",
			args: map[string]interface{}{
				"title": "Updated Test Issue",
			},
			expectation: expectation{
				err:    "id must be a non-empty string",
				output: "",
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
			expectation: expectation{
				err:    "",
				output: "Found 5 issues:\n- TEST-11: Test Issue\n  Priority: None\n  Status: Backlog\n  https://linear.app/linear-mcp-go-test/issue/TEST-11/test-issue\n- TEST-10: Updated Test Issue\n  Priority: None\n  Status: Backlog\n  https://linear.app/linear-mcp-go-test/issue/TEST-10/updated-test-issue\n- TEST-9: Next steps\n  Priority: 4\n  Status: Todo\n  https://linear.app/linear-mcp-go-test/issue/TEST-9/next-steps\n- TEST-3: Connect to Slack\n  Priority: 3\n  Status: Todo\n  https://linear.app/linear-mcp-go-test/issue/TEST-3/connect-to-slack\n- TEST-1: Welcome to Linear 👋\n  Priority: 2\n  Status: Todo\n  https://linear.app/linear-mcp-go-test/issue/TEST-1/welcome-to-linear\n",
			},
		},
		{
			handler: "search_issues",
			name:    "Search by query",
			args: map[string]interface{}{
				"query": "test",
				"limit": float64(5),
			},
			expectation: expectation{
				err:    "",
				output: "Found 0 issues:\n",
			},
		},

		// GetUserIssuesHandler test cases
		{
			handler: "get_user_issues",
			name:    "Current user issues",
			args: map[string]interface{}{
				"limit": float64(5),
			},
			expectation: expectation{
				err:    "",
				output: "Found 1 issues:\n- TEST-1: Welcome to Linear 👋\n  Priority: 2\n  Status: Todo\n  https://linear.app/linear-mcp-go-test/issue/TEST-1/welcome-to-linear\n",
			},
		},
		{
			handler: "get_user_issues",
			name:    "Specific user issues",
			args: map[string]interface{}{
				"userId": USER_ID,
				"limit":  float64(5),
			},
			expectation: expectation{
				err:    "",
				output: "Found 1 issues:\n- TEST-1: Welcome to Linear 👋\n  Priority: 2\n  Status: Todo\n  https://linear.app/linear-mcp-go-test/issue/TEST-1/welcome-to-linear\n",
			},
		},

		// GetIssueHandler test cases
		{
			handler: "get_issue",
			name:    "Valid issue",
			args: map[string]interface{}{
				"issueId": ISSUE_ID, // Using constant instead of hardcoded string
			},
			expectation: expectation{
				err:    "",
				output: "Issue TEST-10: Updated Test Issue\nURL: https://linear.app/linear-mcp-go-test/issue/TEST-10/updated-test-issue\nPriority: None\nStatus: Backlog\nAssignee: None\nTeam: Test Team\n",
			},
		},
		{
			handler: "get_issue",
			name:    "Missing issueId",
			args:    map[string]interface{}{},
			expectation: expectation{
				err:    "issueId must be a non-empty string",
				output: "",
			},
		},

		// AddCommentHandler test cases
		{
			handler: "add_comment",
			name:    "Valid comment",
			args: map[string]interface{}{
				"issueId": ISSUE_ID, // Using constant instead of hardcoded string
				"body":    "Test comment",
			},
			expectation: expectation{
				err:    "",
				output: "Added comment to issue TEST-10\nURL: https://linear.app/linear-mcp-go-test/issue/TEST-10/updated-test-issue#comment-4339c849",
			},
		},
		{
			handler: "add_comment",
			name:    "Missing issueId",
			args: map[string]interface{}{
				"body": "Test comment",
			},
			expectation: expectation{
				err:    "issueId must be a non-empty string",
				output: "",
			},
		},
		{
			handler: "add_comment",
			name:    "Missing body",
			args: map[string]interface{}{
				"issueId": ISSUE_ID, // Using constant instead of hardcoded string
			},
			expectation: expectation{
				err:    "body must be a non-empty string",
				output: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.handler+"_"+tt.name, func(t *testing.T) {
			if tt.write && *record && !*recordWrites {
				t.Skip("Skipping write test when recordWrites=false")
				return
			}

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

			// Compare error
			if diff := cmp.Diff(tt.expectation.err, actualErr); diff != "" {
				t.Errorf("Error mismatch (-want +got):\n%s", diff)
			}

			// Compare output (only if no error is expected)
			if tt.expectation.err == "" {
				if diff := cmp.Diff(tt.expectation.output, actualOutput); diff != "" {
					t.Errorf("Output mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
