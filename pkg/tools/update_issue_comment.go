package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// UpdateCommentTool is the tool definition for updating a comment
var UpdateCommentTool = mcp.NewTool("linear_update_issue_comment",
	mcp.WithDescription("Updates an existing comment on a Linear issue."),
	mcp.WithString("comment", mcp.Required(), mcp.Description("ID or shorthand identifier (e.g., 'comment-53099b37') of the comment to update")),
	mcp.WithString("body", mcp.Required(), mcp.Description("New comment text in markdown format")),
)

// UpdateCommentHandler handles the linear_update_issue_comment tool
func UpdateCommentHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		commentIdentifier, err := request.RequireString("comment")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Resolve comment identifier to a UUID
		commentID, err := resolveCommentIdentifier(linearClient, commentIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve comment: %v", err)}}}, nil
		}

		body, err := request.RequireString("body")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Update the comment
		input := linear.UpdateCommentInput{
			CommentID: commentID,
			Body:      body,
		}

		comment, issue, err := linearClient.UpdateComment(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to update comment: %v", err)}}}, nil
		}

		// Return the result
		resultText := fmt.Sprintf("Updated comment on %s\n", formatIssueIdentifier(issue))
		resultText += fmt.Sprintf("Comment: %s\n", formatCommentIdentifier(comment))
		resultText += fmt.Sprintf("URL: %s", comment.URL)
		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
