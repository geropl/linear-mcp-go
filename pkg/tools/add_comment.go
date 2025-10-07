package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// AddCommentTool is the tool definition for adding a comment
var AddCommentTool = mcp.NewTool("linear_add_comment",
	mcp.WithDescription("Add/post a comment to a Linear issue. To reply to an existing comment, use linear_get_issue_comments to get the comment ID, then pass it in 'thread'."),
	mcp.WithString("issue", mcp.Required(), mcp.Description("ID or identifier (e.g., 'TEAM-123') of the issue to comment on")),
	mcp.WithString("thread", mcp.Description("Optional UUID of parent comment to reply to. Discover from linear_get_issue_comments output or comment URL. Creates a threaded reply instead of top-level comment.")),
	mcp.WithString("body", mcp.Required(), mcp.Description("Comment text in markdown format")),
	mcp.WithString("createAsUser", mcp.Description("Optional custom username to show for the comment")),
)

// AddCommentHandler handles the linear_add_comment tool
func AddCommentHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		body, err := request.RequireString("body")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Extract optional arguments
		createAsUser := request.GetString("createAsUser", "")
		parentID := request.GetString("thread", "")

		// Add the comment
		input := linear.AddCommentInput{
			IssueID:      issueID,
			Body:         body,
			CreateAsUser: createAsUser,
			ParentID:     parentID,
		}

		comment, issue, err := linearClient.AddComment(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to add comment: %v", err)}}}, nil
		}

		// Return the result
		resultText := fmt.Sprintf("Added comment to %s\n", formatIssueIdentifier(issue))
		if parentID != "" {
			resultText += fmt.Sprintf("Thread: %s\n", parentID)
		}
		resultText += fmt.Sprintf("Comment: %s\n", formatCommentIdentifier(comment))
		resultText += fmt.Sprintf("URL: %s", comment.URL)
		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
