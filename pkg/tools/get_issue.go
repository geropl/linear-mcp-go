package tools

import (
	"context"
	"fmt"
	"slices"

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
		args := request.Params.Arguments

		// Validate required arguments
		issueIdentifier, ok := args["issue"].(string)
		if !ok || issueIdentifier == "" {
			return mcp.NewToolResultError("issue must be a non-empty string"), nil
		}

		// Resolve issue identifier to a UUID
		issueID, err := resolveIssueIdentifier(linearClient, issueIdentifier)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve issue: %v", err)), nil
		}

		// Get the issue
		issue, err := linearClient.GetIssue(issueID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get issue: %v", err)), nil
		}

		// Format the result using the full issue formatting
		resultText := formatIssue(issue)
		
		// Add assignee and team information using identifier formatting
		if issue.Assignee != nil {
			resultText += fmt.Sprintf("Assignee: %s\n", formatUserIdentifier(issue.Assignee))
		} else {
			resultText += "Assignee: None\n"
		}
		
		if issue.Team != nil {
			resultText += fmt.Sprintf("Team: %s\n", formatTeamIdentifier(issue.Team))
		} else {
			resultText += "Team: None\n"
		}

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

		// Add comments section
		if issue.Comments != nil && len(issue.Comments.Nodes) > 0 {
			resultText += "\nComments:\n"

			// Create a temporary slice to hold comments in reversed order (oldest first)
			reversedComments := slices.Clone(issue.Comments.Nodes)
			slices.Reverse(reversedComments)

			// Use the reversed comments for display
			for _, comment := range reversedComments {
				if comment.Parent != nil {
					// Skip nested comments (they will be displayed with their parent)
					continue
				}

				userName := "Unknown"
				if comment.User != nil {
					userName = comment.User.Name
				}
				createdAt := comment.CreatedAt.Format("2006-01-02 15:04:05")
				resultText += fmt.Sprintf("- [%s] %s: %s\n", createdAt, userName, comment.Body)

				// Add nested comments if there are any
				if comment.Children != nil && len(comment.Children.Nodes) > 0 {
					// Create a temporary slice to hold child comments in reversed order (oldest first)
					reversedChildComments := slices.Clone(comment.Children.Nodes)
					slices.Reverse(reversedChildComments)

					// Use the reversed child comments for display
					for _, childComment := range reversedChildComments {
						childUserName := "Unknown"
						if childComment.User != nil {
							childUserName = childComment.User.Name
						}
						childCreatedAt := childComment.CreatedAt.Format("2006-01-02 15:04:05")
						resultText += fmt.Sprintf("  - [%s] %s: %s\n", childCreatedAt, childUserName, childComment.Body)
					}
				}
			}
		} else {
			resultText += "\nComments: None\n"
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

		return mcp.NewToolResultText(resultText), nil
	}
}
