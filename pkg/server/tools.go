package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CreateIssueHandler handles the linear_create_issue tool
func CreateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Validate required arguments
		title, ok := args["title"].(string)
		if !ok || title == "" {
			return mcp.NewToolResultError("title must be a non-empty string"), nil
		}

		teamID, ok := args["teamId"].(string)
		if !ok || teamID == "" {
			return mcp.NewToolResultError("teamId must be a non-empty string"), nil
		}

		// Extract optional arguments
		description := ""
		if desc, ok := args["description"].(string); ok {
			description = desc
		}

		var priority *int
		if p, ok := args["priority"].(float64); ok {
			pInt := int(p)
			priority = &pInt
		}

		status := ""
		if s, ok := args["status"].(string); ok {
			status = s
		}

		// Create the issue
		input := linear.CreateIssueInput{
			Title:       title,
			TeamID:      teamID,
			Description: description,
			Priority:    priority,
			Status:      status,
		}

		issue, err := linearClient.CreateIssue(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create issue: %v", err)), nil
		}

		// Return the result
		resultText := fmt.Sprintf("Created issue %s: %s\nURL: %s", issue.Identifier, issue.Title, issue.URL)
		return mcp.NewToolResultText(resultText), nil
	}
}

// UpdateIssueHandler handles the linear_update_issue tool
func UpdateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Validate required arguments
		id, ok := args["id"].(string)
		if !ok || id == "" {
			return mcp.NewToolResultError("id must be a non-empty string"), nil
		}

		// Extract optional arguments
		title := ""
		if t, ok := args["title"].(string); ok {
			title = t
		}

		description := ""
		if desc, ok := args["description"].(string); ok {
			description = desc
		}

		var priority *int
		if p, ok := args["priority"].(float64); ok {
			pInt := int(p)
			priority = &pInt
		}

		status := ""
		if s, ok := args["status"].(string); ok {
			status = s
		}

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
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update issue: %v", err)), nil
		}

		// Return the result
		resultText := fmt.Sprintf("Updated issue %s\nURL: %s", issue.Identifier, issue.URL)
		return mcp.NewToolResultText(resultText), nil
	}
}

// SearchIssuesHandler handles the linear_search_issues tool
func SearchIssuesHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Build search input
		input := linear.SearchIssuesInput{}

		if query, ok := args["query"].(string); ok {
			input.Query = query
		}

		if teamID, ok := args["teamId"].(string); ok {
			input.TeamID = teamID
		}

		if status, ok := args["status"].(string); ok {
			input.Status = status
		}

		if assigneeID, ok := args["assigneeId"].(string); ok {
			input.AssigneeID = assigneeID
		}

		if labelsStr, ok := args["labels"].(string); ok && labelsStr != "" {
			// Split comma-separated labels
			labels := []string{}
			for _, label := range strings.Split(labelsStr, ",") {
				trimmedLabel := strings.TrimSpace(label)
				if trimmedLabel != "" {
					labels = append(labels, trimmedLabel)
				}
			}
			input.Labels = labels
		}

		if priority, ok := args["priority"].(float64); ok {
			pInt := int(priority)
			input.Priority = &pInt
		}

		if estimate, ok := args["estimate"].(float64); ok {
			input.Estimate = &estimate
		}

		if includeArchived, ok := args["includeArchived"].(bool); ok {
			input.IncludeArchived = includeArchived
		}

		if limit, ok := args["limit"].(float64); ok {
			input.Limit = int(limit)
		}

		// Search for issues
		issues, err := linearClient.SearchIssues(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search issues: %v", err)), nil
		}

		// Format the result
		resultText := fmt.Sprintf("Found %d issues:\n", len(issues))
		for _, issue := range issues {
			priorityStr := "None"
			if issue.Priority > 0 {
				priorityStr = strconv.Itoa(issue.Priority)
			}

			statusStr := "None"
			if issue.Status != "" {
				statusStr = issue.Status
			} else if issue.StateName != "" {
				statusStr = issue.StateName
			}

			resultText += fmt.Sprintf("- %s: %s\n  Priority: %s\n  Status: %s\n  %s\n",
				issue.Identifier, issue.Title, priorityStr, statusStr, issue.URL)
		}

		// Add API metrics
		metrics := linearClient.GetMetrics()
		resultText += fmt.Sprintf("\nAPI Metrics:\n  Requests in last hour: %d\n  Remaining requests: %d\n  Average request time: %s",
			metrics.RequestsInLastHour, metrics.RemainingRequests, metrics.AverageRequestTime)

		return mcp.NewToolResultText(resultText), nil
	}
}

// GetUserIssuesHandler handles the linear_get_user_issues tool
func GetUserIssuesHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Build input
		input := linear.GetUserIssuesInput{}

		if userID, ok := args["userId"].(string); ok {
			input.UserID = userID
		}

		if includeArchived, ok := args["includeArchived"].(bool); ok {
			input.IncludeArchived = includeArchived
		}

		if limit, ok := args["limit"].(float64); ok {
			input.Limit = int(limit)
		}

		// Get user issues
		issues, err := linearClient.GetUserIssues(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user issues: %v", err)), nil
		}

		// Format the result
		resultText := fmt.Sprintf("Found %d issues:\n", len(issues))
		for _, issue := range issues {
			priorityStr := "None"
			if issue.Priority > 0 {
				priorityStr = strconv.Itoa(issue.Priority)
			}

			statusStr := "None"
			if issue.Status != "" {
				statusStr = issue.Status
			} else if issue.StateName != "" {
				statusStr = issue.StateName
			}

			resultText += fmt.Sprintf("- %s: %s\n  Priority: %s\n  Status: %s\n  %s\n",
				issue.Identifier, issue.Title, priorityStr, statusStr, issue.URL)
		}

		// Add API metrics
		metrics := linearClient.GetMetrics()
		resultText += fmt.Sprintf("\nAPI Metrics:\n  Requests in last hour: %d\n  Remaining requests: %d\n  Average request time: %s",
			metrics.RequestsInLastHour, metrics.RemainingRequests, metrics.AverageRequestTime)

		return mcp.NewToolResultText(resultText), nil
	}
}

// GetIssueHandler handles the linear_get_issue tool
func GetIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Validate required arguments
		issueID, ok := args["issueId"].(string)
		if !ok || issueID == "" {
			return mcp.NewToolResultError("issueId must be a non-empty string"), nil
		}

		// Get the issue
		issue, err := linearClient.GetIssue(issueID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get issue: %v", err)), nil
		}

		// Format the result
		priorityStr := "None"
		if issue.Priority > 0 {
			priorityStr = strconv.Itoa(issue.Priority)
		}

		statusStr := "None"
		if issue.Status != "" {
			statusStr = issue.Status
		} else if issue.State != nil {
			statusStr = issue.State.Name
		}

		assigneeStr := "None"
		if issue.Assignee != nil {
			assigneeStr = issue.Assignee.Name
		}

		teamStr := "None"
		if issue.Team != nil {
			teamStr = issue.Team.Name
		}

		resultText := fmt.Sprintf("Issue %s: %s\n", issue.Identifier, issue.Title)
		resultText += fmt.Sprintf("URL: %s\n", issue.URL)
		resultText += fmt.Sprintf("Priority: %s\n", priorityStr)
		resultText += fmt.Sprintf("Status: %s\n", statusStr)
		resultText += fmt.Sprintf("Assignee: %s\n", assigneeStr)
		resultText += fmt.Sprintf("Team: %s\n", teamStr)
		
		if issue.Description != "" {
			resultText += fmt.Sprintf("\nDescription:\n%s\n", issue.Description)
		}

		// Add API metrics
		metrics := linearClient.GetMetrics()
		resultText += fmt.Sprintf("\nAPI Metrics:\n  Requests in last hour: %d\n  Remaining requests: %d\n  Average request time: %s",
			metrics.RequestsInLastHour, metrics.RemainingRequests, metrics.AverageRequestTime)

		return mcp.NewToolResultText(resultText), nil
	}
}

// AddCommentHandler handles the linear_add_comment tool
func AddCommentHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Validate required arguments
		issueID, ok := args["issueId"].(string)
		if !ok || issueID == "" {
			return mcp.NewToolResultError("issueId must be a non-empty string"), nil
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

		displayIconURL := ""
		if icon, ok := args["displayIconUrl"].(string); ok {
			displayIconURL = icon
		}

		// Add the comment
		input := linear.AddCommentInput{
			IssueID:        issueID,
			Body:           body,
			CreateAsUser:   createAsUser,
			DisplayIconURL: displayIconURL,
		}

		comment, issue, err := linearClient.AddComment(input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to add comment: %v", err)), nil
		}

		// Return the result
		resultText := fmt.Sprintf("Added comment to issue %s\nURL: %s", issue.Identifier, comment.URL)
		return mcp.NewToolResultText(resultText), nil
	}
}

// RegisterTools registers all Linear tools with the MCP server
func RegisterTools(s *server.MCPServer, linearClient *linear.LinearClient) {
	// Create Issue Tool
	createIssueTool := mcp.NewTool("linear_create_issue",
		mcp.WithDescription("Creates a new Linear issue with specified details. Use this to create tickets for tasks, bugs, or feature requests. Returns the created issue's identifier and URL. Required fields are title and teamId, with optional description, priority (0-4, where 0 is no priority and 1 is urgent), and status."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("teamId", mcp.Required(), mcp.Description("Team ID")),
		mcp.WithString("description", mcp.Description("Issue description")),
		mcp.WithNumber("priority", mcp.Description("Priority (0-4)")),
		mcp.WithString("status", mcp.Description("Issue status")),
	)
	s.AddTool(createIssueTool, CreateIssueHandler(linearClient))

	// Update Issue Tool
	updateIssueTool := mcp.NewTool("linear_update_issue",
		mcp.WithDescription("Updates an existing Linear issue's properties. Use this to modify issue details like title, description, priority, or status. Requires the issue ID and accepts any combination of updatable fields. Returns the updated issue's identifier and URL."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Issue ID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithNumber("priority", mcp.Description("New priority (0-4)")),
		mcp.WithString("status", mcp.Description("New status")),
	)
	s.AddTool(updateIssueTool, UpdateIssueHandler(linearClient))

	// Search Issues Tool
	searchIssuesTool := mcp.NewTool("linear_search_issues",
		mcp.WithDescription("Searches Linear issues using flexible criteria. Supports filtering by any combination of: title/description text, team, status, assignee, labels (comma-separated), priority (1=urgent, 2=high, 3=normal, 4=low), and estimate. Returns up to 10 issues by default (configurable via limit)."),
		mcp.WithString("query", mcp.Description("Optional text to search in title and description")),
		mcp.WithString("teamId", mcp.Description("Filter by team ID")),
		mcp.WithString("status", mcp.Description("Filter by status name (e.g., 'In Progress', 'Done')")),
		mcp.WithString("assigneeId", mcp.Description("Filter by assignee's user ID")),
		mcp.WithString("labels", mcp.Description("Filter by label names (comma-separated)")),
		mcp.WithNumber("priority", mcp.Description("Filter by priority (1=urgent, 2=high, 3=normal, 4=low)")),
		mcp.WithNumber("estimate", mcp.Description("Filter by estimate points")),
		mcp.WithBoolean("includeArchived", mcp.Description("Include archived issues in results (default: false)")),
		mcp.WithNumber("limit", mcp.Description("Max results to return (default: 10)")),
	)
	s.AddTool(searchIssuesTool, SearchIssuesHandler(linearClient))

	// Get User Issues Tool
	getUserIssuesTool := mcp.NewTool("linear_get_user_issues",
		mcp.WithDescription("Retrieves issues assigned to a specific user or the authenticated user if no userId is provided. Returns issues sorted by last updated, including priority, status, and other metadata. Useful for finding a user's workload or tracking assigned tasks."),
		mcp.WithString("userId", mcp.Description("Optional user ID. If not provided, returns authenticated user's issues")),
		mcp.WithBoolean("includeArchived", mcp.Description("Include archived issues in results")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of issues to return (default: 50)")),
	)
	s.AddTool(getUserIssuesTool, GetUserIssuesHandler(linearClient))

	// Get Issue Tool
	getIssueTool := mcp.NewTool("linear_get_issue",
		mcp.WithDescription("Retrieves a single Linear issue by its ID. Returns detailed information about the issue including title, description, priority, status, assignee, and team."),
		mcp.WithString("issueId", mcp.Required(), mcp.Description("ID of the issue to retrieve")),
	)
	s.AddTool(getIssueTool, GetIssueHandler(linearClient))

	// Add Comment Tool
	addCommentTool := mcp.NewTool("linear_add_comment",
		mcp.WithDescription("Adds a comment to an existing Linear issue. Supports markdown formatting in the comment body. Can optionally specify a custom user name and avatar for the comment. Returns the created comment's details including its URL."),
		mcp.WithString("issueId", mcp.Required(), mcp.Description("ID of the issue to comment on")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment text in markdown format")),
		mcp.WithString("createAsUser", mcp.Description("Optional custom username to show for the comment")),
		mcp.WithString("displayIconUrl", mcp.Description("Optional avatar URL for the comment")),
	)
	s.AddTool(addCommentTool, AddCommentHandler(linearClient))
}
