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
		// UpdateCommentHandler test cases
		{
			handler: "update_comment",
			name:    "Valid comment update",
			write:   true,
			args: map[string]interface{}{
				"comment": "ae3d62d6-3f40-4990-867b-5c97dd265a40",
				"body":    "Updated comment text",
			},
		},
		{
			handler: "update_comment",
			name:    "Valid comment update with shorthand",
			write:   true,
			args: map[string]interface{}{
				"comment": "comment-ae3d62d6",
				"body":    "Updated comment text via shorthand",
			},
		},
		{
			handler: "update_comment",
			name:    "Valid comment update with hash only",
			write:   true,
			args: map[string]interface{}{
				"comment": "ae3d62d6",
				"body":    "Updated comment text via hash",
			},
		},
		{
			handler: "update_comment",
			name:    "Missing comment",
			args: map[string]interface{}{
				"body": "Updated comment text",
			},
		},
		{
			handler: "update_comment",
			name:    "Missing body",
			args: map[string]interface{}{
				"comment": "ae3d62d6-3f40-4990-867b-5c97dd265a40",
			},
		},
		{
			handler: "update_comment",
			name:    "Invalid comment identifier",
			write:   true,
			args: map[string]interface{}{
				"comment": "invalid-comment-id",
				"body":    "Updated comment text",
			},
		},
		// GetProjectHandler test cases
		{
			handler: "get_project",
			name:    "By ID",
			args: map[string]interface{}{
				"project": "01bff2dd-ab7f-4464-b425-97073862013f",
			},
		},
		{
			handler: "get_project",
			name:    "Missing project param",
			args:    map[string]interface{}{},
		},
		{
			handler: "get_project",
			name:    "Invalid project",
			args: map[string]interface{}{
				"project": "NONEXISTENT-PROJECT",
			},
		},
		{
			handler: "get_project",
			name:    "By slug",
			args: map[string]interface{}{
				"project": "mcp-tool-investigation-ae44897e42a7",
			},
		},
		{
			handler: "get_project",
			name:    "By name",
			args: map[string]interface{}{
				"project": "MCP tool investigation",
			},
		},
		{
			handler: "get_project",
			name:    "Non-existent slug",
			args: map[string]interface{}{
				"project": "non-existent-slug",
			},
		},
		// SearchProjectsHandler test cases
		{
			handler: "search_projects",
			name:    "Empty query",
			args: map[string]interface{}{
				"query": "",
			},
		},
		{
			handler: "search_projects",
			name:    "No results",
			args: map[string]interface{}{
				"query": "non-existent-project-query",
			},
		},
		{
			handler: "search_projects",
			name:    "Multiple results",
			args: map[string]interface{}{
				"query": "MCP",
			},
		},
		// CreateProjectHandler test cases
		{
			handler: "create_project",
			name:    "Valid project",
			args: map[string]interface{}{
				"name":    "Created Test Project",
				"teamIds": TEAM_ID,
			},
			write: true,
		},
		{
			handler: "create_project",
			name:    "With all optional fields",
			args: map[string]interface{}{
				"name":        "Test Project 2",
				"teamIds":     TEAM_ID,
				"description": "Test Description",
				"leadId":      USER_ID,
				"startDate":   "2024-01-01",
				"targetDate":  "2024-12-31",
			},
			write: true,
		},
		{
			handler: "create_project",
			name:    "Missing name",
			args: map[string]interface{}{
				"teamIds": TEAM_ID,
			},
			write: true,
		},
		{
			handler: "create_project",
			name:    "Invalid team ID",
			args: map[string]interface{}{
				"name":    "Test Project 3",
				"teamIds": "invalid-team-id",
			},
			write: true,
		},
		// UpdateProjectHandler test cases
		{
			handler: "update_project",
			name:    "Valid update",
			args: map[string]interface{}{
				"project": UPDATE_PROJECT_ID,
				"name":    "Updated Project Name",
			},
			write: true,
		},
		{
			handler: "update_project",
			name:    "Update name and description",
			args: map[string]interface{}{
				"project":     UPDATE_PROJECT_ID,
				"name":        "Updated Project Name 2",
				"description": "Updated Description",
			},
			write: true,
		},
		{
			handler: "update_project",
			name:    "Non-existent project",
			args: map[string]interface{}{
				"project": "non-existent-project",
				"name":    "Updated Project Name",
			},
			write: true,
		},
		{
			handler: "update_project",
			name:    "Update only description",
			args: map[string]interface{}{
				"project":     UPDATE_PROJECT_ID,
				"description": "Updated Description Only",
			},
			write: true,
		},
		// GetMilestoneHandler test cases
		{
			handler: "get_milestone",
			name:    "Valid milestone",
			args: map[string]interface{}{
				"milestone": MILESTONE_ID,
			},
		},
		{
			handler: "get_milestone",
			name:    "By name",
			args: map[string]interface{}{
				"milestone": "Test Milestone 2",
			},
		},
		{
			handler: "get_milestone",
			name:    "Non-existent milestone",
			args: map[string]interface{}{
				"milestone": "non-existent-milestone",
			},
		},
		// CreateMilestoneHandler test cases
		{
			handler: "create_milestone",
			name:    "Valid milestone",
			args: map[string]interface{}{
				"name":      "Test Milestone 2.2",
				"projectId": UPDATE_PROJECT_ID,
			},
			write: true,
		},
		{
			handler: "create_milestone",
			name:    "With all optional fields",
			args: map[string]interface{}{
				"name":        "Test Milestone 3.2",
				"projectId":   UPDATE_PROJECT_ID,
				"description": "Test Description",
				"targetDate":  "2024-12-31",
			},
			write: true,
		},
		{
			handler: "create_milestone",
			name:    "Missing name",
			args: map[string]interface{}{
				"projectId": UPDATE_PROJECT_ID,
			},
			write: true,
		},
		{
			handler: "create_milestone",
			name:    "Invalid project ID",
			args: map[string]interface{}{
				"name":      "Test Milestone 3.1",
				"projectId": "invalid-project-id",
			},
			write: true,
		},
		// UpdateMilestoneHandler test cases
		{
			handler: "update_milestone",
			name:    "Valid update",
			args: map[string]interface{}{
				"milestone":   UPDATE_MILESTONE_ID,
				"name":        "Updated Milestone Name 22",
				"description": "Updated Description",
				"targetDate":  "2025-01-01",
			},
			write: true,
		},
		{
			handler: "update_milestone",
			name:    "Non-existent milestone",
			args: map[string]interface{}{
				"milestone": "non-existent-milestone",
				"name":      "Updated Milestone Name",
			},
			write: true,
		},
		// GetInitiativeHandler test cases
		{
			handler: "get_initiative",
			name:    "Valid initiative",
			args: map[string]interface{}{
				"initiative": INITIATIVE_ID,
			},
		},
		{
			handler: "get_initiative",
			name:    "By name",
			args: map[string]interface{}{
				"initiative": "Push for MCP",
			},
		},
		{
			handler: "get_initiative",
			name:    "Non-existent name",
			args: map[string]interface{}{
				"initiative": "non-existent-name",
			},
		},
		// CreateInitiativeHandler test cases
		{
			handler: "create_initiative",
			name:    "Valid initiative",
			args: map[string]interface{}{
				"name": "Created Test Initiative",
			},
			write: true,
		},
		{
			handler: "create_initiative",
			name:    "With description",
			args: map[string]interface{}{
				"name":        "Created Test Initiative 2",
				"description": "Test Description",
			},
			write: true,
		},
		{
			handler: "create_initiative",
			name:    "Missing name",
			args:    map[string]interface{}{},
			write:   true,
		},
		// UpdateInitiativeHandler test cases
		{
			handler: "update_initiative",
			name:    "Valid update",
			args: map[string]interface{}{
				"initiative":  UPDATE_INITIATIVE_ID,
				"name":        "Updated Initiative Name",
				"description": "Updated Description",
			},
			write: true,
		},
		{
			handler: "update_initiative",
			name:    "Non-existent initiative",
			args: map[string]interface{}{
				"initiative": "non-existent-initiative",
				"name":       "Updated Initiative Name",
			},
			write: true,
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
			case "update_comment":
				handler = tools.UpdateCommentHandler(client)
			case "get_project":
				handler = tools.GetProjectHandler(client)
			case "search_projects":
				handler = tools.SearchProjectsHandler(client)
			case "create_project":
				handler = tools.CreateProjectHandler(client)
			case "update_project":
				handler = tools.UpdateProjectHandler(client)
			case "get_milestone":
				handler = tools.GetMilestoneHandler(client)
			case "create_milestone":
				handler = tools.CreateMilestoneHandler(client)
			case "update_milestone":
				handler = tools.UpdateMilestoneHandler(client)
			case "get_initiative":
				handler = tools.GetInitiativeHandler(client)
			case "create_initiative":
				handler = tools.CreateInitiativeHandler(client)
			case "update_initiative":
				handler = tools.UpdateInitiativeHandler(client)
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
