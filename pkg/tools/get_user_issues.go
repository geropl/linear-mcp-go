package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetUserIssuesTool is the tool definition for getting user issues
var GetUserIssuesTool = mcp.NewTool("linear_get_user_issues",
	mcp.WithDescription("Retrieves issues assigned to a user."),
	mcp.WithString("user", mcp.Description("Optional user identifier (UUID, name, or email). If not provided, returns authenticated user's issues")),
	mcp.WithBoolean("includeArchived", mcp.Description("Include archived issues in results")),
	mcp.WithNumber("limit", mcp.Description("Maximum number of issues to return (default: 50)")),
)

// GetUserIssuesHandler handles the linear_get_user_issues tool
func GetUserIssuesHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Build input
		input := linear.GetUserIssuesInput{}

		if user, ok := args["user"].(string); ok && user != "" {
			// Resolve user identifier to a user ID
			userID, err := resolveUserIdentifier(linearClient, user)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve user: %v", err)), nil
			}
			input.UserID = userID
		}

		if includeArchived, ok := args["includeArchived"].(bool); ok {
			input.IncludeArchived = includeArchived
		}

		if limit, ok := args["limit"].(float64); ok {
			input.Limit = int(limit)
		}

		// Get user issues
		issues, err := linearClient.GetUserIssues(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user issues: %v", err)), nil
		}

		// Format the result
		resultText := fmt.Sprintf("Found %d issues:\n", len(issues))
		for _, issue := range issues {
			// Create a temporary Issue object to use with formatIssueIdentifier
			tempIssue := &linear.Issue{
				ID:         issue.ID,
				Identifier: issue.Identifier,
			}
			
			priorityStr := "None"
			if issue.Priority > 0 {
				priorityStr = strconv.Itoa(issue.Priority)
			}

			statusStr := "None"
			if issue.Status != "" {
				statusStr = issue.Status
			} else if issue.StateName != "" {
				statusStr = issue.StateName
			}

			resultText += fmt.Sprintf("- %s\n", formatIssueIdentifier(tempIssue))
			resultText += fmt.Sprintf("  Title: %s\n", issue.Title)
			resultText += fmt.Sprintf("  Priority: %s\n", priorityStr)
			resultText += fmt.Sprintf("  Status: %s\n", statusStr)
			resultText += fmt.Sprintf("  URL: %s\n", issue.URL)
		}

		return mcp.NewToolResultText(resultText), nil
	}
}
