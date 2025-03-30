package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// AddCommentTool is the tool definition for adding a comment
var AddCommentTool = mcp.NewTool("linear_add_comment",
	mcp.WithDescription("Adds a comment to a Linear issue."),
	mcp.WithString("issue", mcp.Required(), mcp.Description("ID or identifier (e.g., 'TEAM-123') of the issue to comment on")),
	mcp.WithString("body", mcp.Required(), mcp.Description("Comment text in markdown format")),
	mcp.WithString("thread", mcp.Description("Optional ID of a parent comment / thread to reply to")),
	mcp.WithString("createAsUser", mcp.Description("Optional custom username to show for the comment")),
)

// AddCommentHandler handles the linear_add_comment tool
func AddCommentHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		body, ok := args["body"].(string)
		if !ok || body == "" {
			return mcp.NewToolResultError("body must be a non-empty string"), nil
		}

		// Extract optional arguments
		createAsUser := ""
		if user, ok := args["createAsUser"].(string); ok {
			createAsUser = user
		}

		parentID := ""
		if parent, ok := args["thread"].(string); ok {
			parentID = parent
		}

		// Add the comment
		input := linear.AddCommentInput{
			IssueID:      issueID,
			Body:         body,
			CreateAsUser: createAsUser,
			ParentID:     parentID,
		}

		comment, issue, err := linearClient.AddComment(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to add comment: %v", err)), nil
		}

		// Return the result
		resultText := fmt.Sprintf("Added comment to %s\n", formatIssueIdentifier(issue))
		if parentID != "" {
			resultText += fmt.Sprintf("Thread: %s\n", parentID)
		}
		resultText += fmt.Sprintf("Comment: %s\n", formatCommentIdentifier(comment))
		resultText += fmt.Sprintf("URL: %s", comment.URL)
		return mcp.NewToolResultText(resultText), nil
	}
}
