package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetIssueTool is the tool definition for getting an issue
var GetIssueTool = mcp.NewTool("linear_get_issue",
	mcp.WithDescription("Retrieves a single Linear issue."),
	mcp.WithString("issue", mcp.Required(), mcp.Description("ID or identifier (e.g., 'TEAM-123') of the issue to retrieve")),
)

// GetIssueHandler handles the linear_get_issue tool
func GetIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		issueIdentifier, err := request.RequireString("issue")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Resolve issue identifier to a UUID
		issueID, err := resolveIssueIdentifier(linearClient, issueIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve issue: %v", err)}}}, nil
		}

		// Get the issue
		issue, err := linearClient.GetIssue(issueID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get issue: %v", err)}}}, nil
		}

		// Format the result using the full issue formatting
		resultText := formatIssue(issue)

		// Add assignee and team information using identifier formatting
		if issue.Assignee != nil {
			resultText += fmt.Sprintf("Assignee: %s\n", formatUserIdentifier(issue.Assignee))
		} else {
			resultText += "Assignee: None\n"
		}

		resultText += fmt.Sprintf("%s\n", formatTeamIdentifier(issue.Team))

		// Add attachments section if there are attachments
		if issue.Attachments != nil && len(issue.Attachments.Nodes) > 0 {
			resultText += "\nAttachments:\n"

			// Display all attachments in a simple list without grouping by source type
			for _, attachment := range issue.Attachments.Nodes {
				resultText += fmt.Sprintf("- %s: %s\n", attachment.Title, attachment.URL)
				if attachment.Subtitle != "" {
					resultText += fmt.Sprintf("  %s\n", attachment.Subtitle)
				}
			}
		} else {
			resultText += "\nAttachments: None\n"
		}

		// Add related issues section
		if (issue.Relations != nil && len(issue.Relations.Nodes) > 0) ||
			(issue.InverseRelations != nil && len(issue.InverseRelations.Nodes) > 0) {
			resultText += "\nRelated Issues:\n"

			// Add direct relations
			if issue.Relations != nil && len(issue.Relations.Nodes) > 0 {
				for _, relation := range issue.Relations.Nodes {
					if relation.RelatedIssue != nil {
						resultText += fmt.Sprintf("- %s\n  Title: %s\n  RelationType: %s\n  URL: %s\n",
							formatIssueIdentifier(relation.RelatedIssue),
							relation.RelatedIssue.Title,
							relation.Type,
							relation.RelatedIssue.URL)
					}
				}
			}

			// Add inverse relations
			if issue.InverseRelations != nil && len(issue.InverseRelations.Nodes) > 0 {
				for _, relation := range issue.InverseRelations.Nodes {
					if relation.Issue != nil {
						resultText += fmt.Sprintf("- %s\n  Title: %s\n  RelationType: %s (inverse)\n  URL: %s\n",
							formatIssueIdentifier(relation.Issue),
							relation.Issue.Title,
							relation.Type,
							relation.Issue.URL)
					}
				}
			}
		} else {
			resultText += "\nRelated Issues: None\n"
		}

		// Note about comments
		resultText += "\nComments: Use the linear_get_issue_comments tool to retrieve comments for this issue.\n"

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
