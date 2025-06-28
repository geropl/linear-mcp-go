package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// CreateIssueTool is the tool definition for creating issues
var CreateIssueTool = mcp.NewTool("linear_create_issue",
	mcp.WithDescription("Creates a new Linear issue."),
	mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
	mcp.WithString("team", mcp.Required(), mcp.Description("Team identifier (key, UUID or name)")),
	mcp.WithString("description", mcp.Description("Issue description")),
	mcp.WithNumber("priority", mcp.Description("Priority (0-4)")),
	mcp.WithString("status", mcp.Description("Issue status")),
	mcp.WithString("parentIssue", mcp.Description("Optional parent issue ID or identifier (e.g., 'TEAM-123') to create a sub-issue")),
	mcp.WithString("labels", mcp.Description("Optional comma-separated list of label IDs or names to assign")),
)

// CreateIssueHandler handles the linear_create_issue tool
func CreateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		title, err := request.RequireString("title")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		team, err := request.RequireString("team")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Resolve team identifier to a team ID
		teamId, err := resolveTeamIdentifier(linearClient, team)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve team: %v", err)}}}, nil
		}

		// Extract optional arguments
		description := request.GetString("description", "")

		var priority *int
		if p, err := request.RequireInt("priority"); err == nil {
			priority = &p
		}

		status := request.GetString("status", "")

		// Extract parentIssue parameter and resolve it if needed
		var parentID *string
		if parentIssue, err := request.RequireString("parentIssue"); err == nil && parentIssue != "" {
			resolvedParentID, err := resolveIssueIdentifier(linearClient, parentIssue)
			if err != nil {
				return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve parent issue: %v", err)}}}, nil
			}
			parentID = &resolvedParentID
		}

		// Extract labels parameter and resolve them if needed
		var labelIDs []string
		if labelsStr, err := request.RequireString("labels"); err == nil && labelsStr != "" {
			// Split comma-separated labels
			var labelIdentifiers []string
			for _, label := range strings.Split(labelsStr, ",") {
				trimmedLabel := strings.TrimSpace(label)
				if trimmedLabel != "" {
					labelIdentifiers = append(labelIdentifiers, trimmedLabel)
				}
			}

			// Resolve label identifiers to UUIDs
			if len(labelIdentifiers) > 0 {
				resolvedLabelIDs, err := resolveLabelIdentifiers(linearClient, teamId, labelIdentifiers)
				if err != nil {
					return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve labels: %v", err)}}}, nil
				}
				labelIDs = resolvedLabelIDs
			}
		}

		// Create the issue
		input := linear.CreateIssueInput{
			Title:       title,
			TeamID:      teamId,
			Description: description,
			Priority:    priority,
			Status:      status,
			ParentID:    parentID,
			LabelIDs:    labelIDs,
		}

		issue, err := linearClient.CreateIssue(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to create issue: %v", err)}}}, nil
		}

		// Return the result
		resultText := fmt.Sprintf("Created %s", formatIssueIdentifier(issue))
		resultText += fmt.Sprintf("\nTitle: %s", issue.Title)
		resultText += fmt.Sprintf("\nURL: %s", issue.URL)
		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
