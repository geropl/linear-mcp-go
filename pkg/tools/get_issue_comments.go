package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetIssueCommentsTool is the tool definition for getting paginated comments for an issue
var GetIssueCommentsTool = mcp.NewTool("linear_get_issue_comments",
	mcp.WithDescription("Retrieves paginated comments for a Linear issue. Supports navigating through comment threads and pagination."),
	mcp.WithString("issue", mcp.Required(), mcp.Description("ID or identifier (e.g., 'TEAM-123') of the issue to retrieve comments for")),
	mcp.WithString("thread", mcp.Description("Optional UUID of the parent comment to retrieve replies for. If not provided, returns top-level comments.")),
	mcp.WithNumber("limit", mcp.Description("Maximum number of comments to return (default: 10)")),
	mcp.WithString("after", mcp.Description("Cursor for pagination, to get comments after this point")),
)

// GetIssueCommentsHandler handles the linear_get_issue_comments tool
func GetIssueCommentsHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Validate required arguments
		issueIdentifier, ok := args["issue"].(string)
		if !ok || issueIdentifier == "" {
			return mcp.NewToolResultError("issue must be a non-empty string"), nil
		}

		// Extract optional arguments
		var parentID string
		if parentIDArg, ok := args["thread"].(string); ok {
			parentID = parentIDArg
		}

		var limit int
		if limitArg, ok := args["limit"].(float64); ok {
			limit = int(limitArg)
		}

		var afterCursor string
		if afterArg, ok := args["after"].(string); ok {
			afterCursor = afterArg
		}

		// Resolve issue identifier to a UUID
		issueID, err := resolveIssueIdentifier(linearClient, issueIdentifier)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve issue: %v", err)), nil
		}

		// Get the issue for basic information
		issue, err := linearClient.GetIssue(issueID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get issue: %v", err)), nil
		}

		// Get the comments
		commentsInput := linear.GetIssueCommentsInput{
			IssueID:     issueID,
			ParentID:    parentID,
			Limit:       limit,
			AfterCursor: afterCursor,
		}

		comments, err := linearClient.GetIssueComments(commentsInput)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get comments: %v", err)), nil
		}

		// Format the result
		var resultText string

		// Add issue information
		resultText += fmt.Sprintf("Issue: %s (UUID: %s)\n", issue.Identifier, issue.ID)

		// Add thread information
		if parentID == "" {
			resultText += "Thread: root (top-level comments)\n"
		} else {
			resultText += fmt.Sprintf("Thread: %s (replies to comment)\n", parentID)
		}

		resultText += "\n"

		// Add comments
		if len(comments.Nodes) > 0 {
			resultText += "Comments:\n"

			for _, comment := range comments.Nodes {
				createdAt := comment.CreatedAt.Format("2006-01-02 15:04:05")
				hasReplies := false
				if comment.Children != nil && len(comment.Children.Nodes) > 0 {
					hasReplies = true
				}

				resultText += fmt.Sprintf("- ID: %s\n  %s\n  CreatedAt: %s\n  HasReplies: %s\n  Body: %s\n",
					formatCommentIdentifier(&comment),
					formatUserIdentifier(comment.User),
					createdAt,
					formatBool(hasReplies),
					comment.Body)
			}
		} else {
			resultText += "Comments: None\n"
		}

		// Add pagination information
		resultText += "\nPagination:\n"
		resultText += fmt.Sprintf("Has more comments: %s\n", formatBool(comments.PageInfo.HasNextPage))

		if comments.PageInfo.HasNextPage {
			resultText += fmt.Sprintf("Next cursor: %s\n", comments.PageInfo.EndCursor)
		}

		return mcp.NewToolResultText(resultText), nil
	}
}

// formatBool formats a boolean value as "yes" or "no"
func formatBool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
