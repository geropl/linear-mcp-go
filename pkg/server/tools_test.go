package server

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/geropl/linear-mcp-go/pkg/tools"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Shared constants and expectation struct are defined in test_helpers.go

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
			name:    "Valid issue with team",
			args: map[string]interface{}{
				"title": "Test Issue",
				"team":  TEAM_ID,
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Valid issue with team UUID",
			args: map[string]interface{}{
				"title": "Test Issue with team UUID",
				"team":  TEAM_ID,
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Valid issue with team name",
			args: map[string]interface{}{
				"title": "Test Issue with team name",
				"team":  TEAM_NAME,
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Valid issue with team key",
			args: map[string]interface{}{
				"title": "Test Issue with team key",
				"team":  TEAM_KEY,
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Create sub issue",
			args: map[string]interface{}{
				"title":       "Sub Issue",
				"team":        TEAM_ID,
				"parentIssue": "1c2de93f-4321-4015-bfde-ee893ef7976f", // UUID for TEST-10
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Create sub issue from identifier",
			args: map[string]interface{}{
				"title":       "Sub Issue",
				"team":        TEAM_ID,
				"parentIssue": "TEST-10",
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Create issue with labels",
			args: map[string]interface{}{
				"title":  "Issue with Labels",
				"team":   TEAM_ID,
				"labels": "team label 1",
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Create sub issue with labels",
			args: map[string]interface{}{
				"title":       "Sub Issue with Labels",
				"team":        TEAM_ID,
				"parentIssue": "1c2de93f-4321-4015-bfde-ee893ef7976f", // UUID for TEST-10
				"labels":      "ws-label 2,Feature",
			},
			write: true,
		},
		{
			handler: "create_issue",
			name:    "Missing title",
			args: map[string]interface{}{
				"team": TEAM_ID,
			},
		},
		{
			handler: "create_issue",
			name:    "Missing team",
			args: map[string]interface{}{
				"title": "Test Issue",
			},
		},
		{
			handler: "create_issue",
			name:    "Invalid team",
			args: map[string]interface{}{
				"title": "Test Issue",
				"team":  "NonExistentTeam",
			},
		},

		// UpdateIssueHandler test cases
		{
			handler: "update_issue",
			name:    "Valid update",
			args: map[string]interface{}{
				"issue": ISSUE_ID,
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
				"team":  TEAM_ID,
				"limit": float64(5),
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
				"user":  USER_ID,
				"limit": float64(5),
			},
		},

		// GetIssueHandler test cases
		{
			handler: "get_issue",
			name:    "Valid issue",
			args: map[string]interface{}{
				"issue": ISSUE_ID,
			},
		},
		{
			handler: "get_issue",
			name:    "Get comment issue",
			args: map[string]interface{}{
				"issue": COMMENT_ISSUE_ID,
			},
		},
		{
			handler: "get_issue",
			name:    "Missing issue",
			args:    map[string]interface{}{},
		},
		{
			handler: "get_issue",
			name:    "Missing issueId",
			args: map[string]interface{}{
				"issue": "NONEXISTENT-123",
			},
		},

		// GetIssueCommentsHandler test cases
		{
			handler: "get_issue_comments",
			name:    "Valid issue",
			args: map[string]interface{}{
				"issue": ISSUE_ID,
			},
		},
		{
			handler: "get_issue_comments",
			name:    "Missing issue",
			args:    map[string]interface{}{},
		},
		{
			handler: "get_issue_comments",
			name:    "Invalid issue",
			args: map[string]interface{}{
				"issue": "NONEXISTENT-123",
			},
		},
		{
			handler: "get_issue_comments",
			name:    "With limit",
			args: map[string]interface{}{
				"issue": ISSUE_ID,
				"limit": float64(3),
			},
		},
		{
			handler: "get_issue_comments",
			name:    "With_thread_parameter",
			args: map[string]interface{}{
				"issue":  ISSUE_ID,
				"thread": "ae3d62d6-3f40-4990-867b-5c97dd265a40", // ID of a comment to get replies for
			},
		},
		{
			handler: "get_issue_comments",
			name:    "Thread_with_pagination",
			args: map[string]interface{}{
				"issue":  ISSUE_ID,
				"thread": "ae3d62d6-3f40-4990-867b-5c97dd265a40", // ID of a comment to get replies for
				"limit":  float64(2),
			},
		},

		// AddCommentHandler test cases
		{
			handler: "add_comment",
			name:    "Valid comment",
			write:   true,
			args: map[string]interface{}{
				"issue": ISSUE_ID,
				"body":  "Test comment",
			},
		},
		{
			handler: "add_comment",
			name:    "Reply_to_comment",
			write:   true,
			args: map[string]interface{}{
				"issue":  ISSUE_ID,
				"body":   "This is a reply to the comment",
				"thread": "ae3d62d6-3f40-4990-867b-5c97dd265a40", // ID of the comment to reply to
			},
		},
		{
			handler: "add_comment",
			name:    "Missing issue",
			args: map[string]interface{}{
				"body": "Test comment",
			},
		},
		{
			handler: "add_comment",
			name:    "Missing body",
			args: map[string]interface{}{
				"issue": ISSUE_ID,
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
				handler = tools.GetTeamsHandler(client)
			case "create_issue":
				handler = tools.CreateIssueHandler(client)
			case "update_issue":
				handler = tools.UpdateIssueHandler(client)
			case "search_issues":
				handler = tools.SearchIssuesHandler(client)
			case "get_user_issues":
				handler = tools.GetUserIssuesHandler(client)
			case "get_issue":
				handler = tools.GetIssueHandler(client)
			case "get_issue_comments":
				handler = tools.GetIssueCommentsHandler(client)
			case "add_comment":
				handler = tools.AddCommentHandler(client)
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

// readGoldenFile and writeGoldenFile are defined in test_helpers.go
