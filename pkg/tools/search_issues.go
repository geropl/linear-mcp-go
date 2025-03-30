package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// SearchIssuesTool is the tool definition for searching issues
var SearchIssuesTool = mcp.NewTool("linear_search_issues",
	mcp.WithDescription("Searches Linear issues."),
	mcp.WithString("query", mcp.Description("Optional text to search in title and description")),
	mcp.WithString("team", mcp.Description("Filter by team identifier (UUID, name, or key)")),
	mcp.WithString("status", mcp.Description("Filter by status name (e.g., 'In Progress', 'Done')")),
	mcp.WithString("assignee", mcp.Description("Filter by assignee identifier (UUID, name, or email)")),
	mcp.WithString("labels", mcp.Description("Filter by label names (comma-separated)")),
	mcp.WithNumber("priority", mcp.Description("Filter by priority (1=urgent, 2=high, 3=normal, 4=low)")),
	mcp.WithNumber("estimate", mcp.Description("Filter by estimate points")),
	mcp.WithBoolean("includeArchived", mcp.Description("Include archived issues in results (default: false)")),
	mcp.WithNumber("limit", mcp.Description("Max results to return (default: 10)")),
)

// SearchIssuesHandler handles the linear_search_issues tool
func SearchIssuesHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Build search input
		input := linear.SearchIssuesInput{}

		if query, ok := args["query"].(string); ok {
			input.Query = query
		}

		if team, ok := args["team"].(string); ok && team != "" {
			// Resolve team identifier to a team ID
			teamID, err := resolveTeamIdentifier(linearClient, team)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve team: %v", err)), nil
			}
			input.TeamID = teamID
		}

		if status, ok := args["status"].(string); ok {
			input.Status = status
		}

		if assignee, ok := args["assignee"].(string); ok && assignee != "" {
			// Resolve assignee identifier to a user ID
			assigneeID, err := resolveUserIdentifier(linearClient, assignee)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve assignee: %v", err)), nil
			}
			input.AssigneeID = assigneeID
		}

		if labelsStr, ok := args["labels"].(string); ok && labelsStr != "" {
			// Split comma-separated labels
			labels := []string{}
			for _, label := range strings.Split(labelsStr, ",") {
				trimmedLabel := strings.TrimSpace(label)
				if trimmedLabel != "" {
					labels = append(labels, trimmedLabel)
				}
			}
			input.Labels = labels
		}

		if priority, ok := args["priority"].(float64); ok {
			pInt := int(priority)
			input.Priority = &pInt
		}

		if estimate, ok := args["estimate"].(float64); ok {
			input.Estimate = &estimate
		}

		if includeArchived, ok := args["includeArchived"].(bool); ok {
			input.IncludeArchived = includeArchived
		}

		if limit, ok := args["limit"].(float64); ok {
			input.Limit = int(limit)
		}

		// Search for issues
		issues, err := linearClient.SearchIssues(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search issues: %v", err)), nil
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
