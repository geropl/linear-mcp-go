package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// ReplyToCommentTool is a specialized tool for replying to comments
var ReplyToCommentTool = mcp.NewTool("linear_reply_to_comment",
	mcp.WithDescription("Reply to an existing comment on a Linear issue. Convenience tool that automatically resolves the issue from the comment. Accepts comment URLs, UUIDs, or shorthand identifiers."),
	mcp.WithString("thread", mcp.Required(), mcp.Description("Comment to reply to. Accepts: full Linear comment URL, UUID, shorthand (comment-abc123), or hash (abc123).")),
	mcp.WithString("body", mcp.Required(), mcp.Description("Reply text in markdown format")),
	mcp.WithString("createAsUser", mcp.Description("Optional custom username to show for the reply")),
)

// ReplyToCommentHandler handles the linear_reply_to_comment tool
func ReplyToCommentHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		threadIdentifier, err := request.RequireString("thread")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		body, err := request.RequireString("body")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		createAsUser := request.GetString("createAsUser", "")

		// Resolve the parent comment to get its UUID
		parentCommentID, err := resolveCommentIdentifier(linearClient, threadIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to resolve comment: %v", err)}}}, nil
		}

		// Get the parent comment to find its issue
		parentComment, err := linearClient.GetComment(parentCommentID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get parent comment: %v", err)}}}, nil
		}

		if parentComment.Issue == nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Parent comment does not have an associated issue"}}}, nil
		}

		// Add the reply
		input := linear.AddCommentInput{
			IssueID:      parentComment.Issue.ID,
			Body:         body,
			CreateAsUser: createAsUser,
			ParentID:     parentCommentID,
		}

		comment, issue, err := linearClient.AddComment(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to add reply: %v", err)}}}, nil
		}

		// Return the result
		resultText := fmt.Sprintf("Added reply to comment on %s\n", formatIssueIdentifier(issue))
		resultText += fmt.Sprintf("In reply to thread: %s\n", parentCommentID)
		resultText += fmt.Sprintf("Comment ID: %s\n", comment.ID)
		resultText += fmt.Sprintf("Thread (for replies): %s\n", comment.ID)
		resultText += fmt.Sprintf("URL: %s", comment.URL)
		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
