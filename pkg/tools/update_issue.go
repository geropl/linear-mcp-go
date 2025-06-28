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
	mcp.WithString("team", mcp.Description("New team (UUID, name, or key)")),
	mcp.WithString("projectId", mcp.Description("New project ID")),
	mcp.WithString("milestoneId", mcp.Description("New milestone ID")),
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

		// Resolve team identifier to a team ID
		var teamID string
		team := request.GetString("team", "")
		if team != "" {
			// Resolve team identifier to a team ID
			teamID, err = resolveTeamIdentifier(linearClient, team)
			if err != nil {
				return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve team: %v", err)}}}, nil
			}
		}

		status := request.GetString("status", "")
		projectID := request.GetString("projectId", "")
		milestoneID := request.GetString("milestoneId", "")

		// Update the issue
		input := linear.UpdateIssueInput{
			ID:          id,
			Title:       title,
			Description: description,
			Priority:    priority,
			Status:      status,
			TeamID:      teamID,
			ProjectID:   projectID,
			MilestoneID: milestoneID,
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
