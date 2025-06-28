package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// UpdateIssueTool is the tool definition for updating issues
var UpdateIssueTool = mcp.NewTool("linear_update_issue",
	mcp.WithDescription("Updates an existing Linear issue."),
	mcp.WithString("issue", mcp.Required(), mcp.Description("Issue ID or identifier (e.g., 'TEAM-123')")),
	mcp.WithString("title", mcp.Description("New title")),
	mcp.WithString("description", mcp.Description("New description")),
	mcp.WithNumber("priority", mcp.Description("New priority (0-4)")),
	mcp.WithString("status", mcp.Description("New status")),
)

// UpdateIssueHandler handles the linear_update_issue tool
func UpdateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		issueIdentifier, err := request.RequireString("issue")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Resolve issue identifier to a UUID
		id, err := resolveIssueIdentifier(linearClient, issueIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve issue: %v", err)}}}, nil
		}

		// Extract optional arguments
		title := request.GetString("title", "")
		description := request.GetString("description", "")

		var priority *int
		if p, err := request.RequireInt("priority"); err == nil {
			priority = &p
		}

		status := request.GetString("status", "")

		// Update the issue
		input := linear.UpdateIssueInput{
			ID:          id,
			Title:       title,
			Description: description,
			Priority:    priority,
			Status:      status,
		}

		issue, err := linearClient.UpdateIssue(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to update issue: %v", err)}}}, nil
		}

		// Return the result
		resultText := fmt.Sprintf("Updated %s", formatIssueIdentifier(issue))
		resultText += fmt.Sprintf("\nURL: %s", issue.URL)
		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
