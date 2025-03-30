package server

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// resolveParentIssueIdentifier resolves a parent issue identifier (UUID or "TEAM-123") to a UUID
func resolveParentIssueIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find an issue by identifier
	issue, err := linearClient.GetIssueByIdentifier(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to resolve parent issue identifier '%s': %v", identifier, err)
	}

	return issue.ID, nil
}

// resolveLabelIdentifiers resolves a list of label identifiers (UUIDs or names) to UUIDs
func resolveLabelIdentifiers(linearClient *linear.LinearClient, teamID string, labelIdentifiers []string) ([]string, error) {
	// Separate UUIDs and names
	var labelUUIDs []string
	var labelNames []string

	for _, identifier := range labelIdentifiers {
		if isValidUUID(identifier) {
			labelUUIDs = append(labelUUIDs, identifier)
		} else {
			labelNames = append(labelNames, identifier)
		}
	}

	// If there are no names to resolve, return the UUIDs directly
	if len(labelNames) == 0 {
		return labelUUIDs, nil
	}

	// Get labels by name
	labels, err := linearClient.GetLabelsByName(teamID, labelNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels by name: %v", err)
	}

	// Check if all label names were found
	if len(labels) < len(labelNames) {
		// Create a map of found label names for quick lookup
		foundLabels := make(map[string]bool)
		for _, label := range labels {
			foundLabels[label.Name] = true
		}

		// Find which label names were not found
		var missingLabels []string
		for _, name := range labelNames {
			if !foundLabels[name] {
				missingLabels = append(missingLabels, name)
			}
		}

		return nil, fmt.Errorf("label(s) not found: %s", strings.Join(missingLabels, ", "))
	}

	// Add the resolved label UUIDs to the result
	for _, label := range labels {
		labelUUIDs = append(labelUUIDs, label.ID)
	}

	return labelUUIDs, nil
}

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

		// Check for team parameter or teamId parameter
		team, ok := args["team"].(string)
		if !ok || team == "" {
			return mcp.NewToolResultError("team must be a non-empty string (UUID, name, or key)"), nil
		}

		// Resolve team identifier to a team ID
		teamId, err := resolveTeamIdentifier(linearClient, team)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve team: %v", err)), nil
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

		// Extract parentIssue parameter and resolve it if needed
		var parentID *string
		if parentIssue, ok := args["parentIssue"].(string); ok && parentIssue != "" {
			resolvedParentID, err := resolveParentIssueIdentifier(linearClient, parentIssue)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve parent issue: %v", err)), nil
			}
			parentID = &resolvedParentID
		}

		// Extract labels parameter and resolve them if needed
		var labelIDs []string
		if labelsStr, ok := args["labels"].(string); ok && labelsStr != "" {
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
					return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve labels: %v", err)), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create issue: %v", err)), nil
		}

		// Return the result
		resultText := fmt.Sprintf("Created issue: %s\nTitle: %s\nURL: %s", issue.Identifier, issue.Title, issue.URL)
		return mcp.NewToolResultText(resultText), nil
	}
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(uuidStr string) bool {
	return uuid.Validate(uuidStr) == nil
}

// resolveTeamIdentifier resolves a team identifier (UUID, name, or key) to a team ID
func resolveTeamIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find a team by name or key
	teams, err := linearClient.GetTeams("")
	if err != nil {
		return "", fmt.Errorf("failed to get teams: %v", err)
	}

	// First try exact match on name or key
	for _, team := range teams {
		if team.Name == identifier || team.Key == identifier {
			return team.ID, nil
		}
	}

	// If no exact match, try case-insensitive match
	identifierLower := strings.ToLower(identifier)
	for _, team := range teams {
		if strings.ToLower(team.Name) == identifierLower || strings.ToLower(team.Key) == identifierLower {
			return team.ID, nil
		}
	}

	return "", fmt.Errorf("no team found with identifier '%s'", identifier)
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
		} else {
			resultText += "\nDescription: None\n"
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
						resultText += fmt.Sprintf("- %s: %s\n  RelationType: %s\n  URL: %s\n",
							relation.RelatedIssue.Identifier,
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
						resultText += fmt.Sprintf("- %s: %s\n  RelationType: %s (inverse)\n  URL: %s\n",
							relation.Issue.Identifier,
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

// GetTeamsHandler handles the linear_get_teams tool
func GetTeamsHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Extract optional name filter
		name := ""
		if nameArg, ok := args["name"].(string); ok {
			name = nameArg
		}

		// Get teams
		teams, err := linearClient.GetTeams(name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get teams: %v", err)), nil
		}

		// Format the result
		resultText := fmt.Sprintf("Found %d teams:\n", len(teams))
		for _, team := range teams {
			resultText += fmt.Sprintf("- %s (Key: %s)\n  ID: %s\n", team.Name, team.Key, team.ID)
		}

		return mcp.NewToolResultText(resultText), nil
	}
}

// GetReadOnlyToolNames returns the names of all read-only tools
func GetReadOnlyToolNames() map[string]bool {
	return map[string]bool{
		"linear_search_issues":   true,
		"linear_get_user_issues": true,
		"linear_get_issue":       true,
		"linear_get_teams":       true,
	}
}

// RegisterTools registers all Linear tools with the MCP server
func RegisterTools(s *server.MCPServer, linearClient *linear.LinearClient, writeAccess bool) {
	// Register tools, based on writeAccess
	addTool := func(tool mcp.Tool, handler server.ToolHandlerFunc) {
		if !writeAccess {
			if readOnly := GetReadOnlyToolNames()[tool.Name]; !readOnly {
				// Skip registering write tools if write access is disabled
				return
			}
		}
		s.AddTool(tool, handler)
	}

	// Search Issues Tool (read-only)
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
	addTool(searchIssuesTool, SearchIssuesHandler(linearClient))

	// Get User Issues Tool (read-only)
	getUserIssuesTool := mcp.NewTool("linear_get_user_issues",
		mcp.WithDescription("Retrieves issues assigned to a specific user or the authenticated user if no userId is provided. Returns issues sorted by last updated, including priority, status, and other metadata. Useful for finding a user's workload or tracking assigned tasks."),
		mcp.WithString("userId", mcp.Description("Optional user ID. If not provided, returns authenticated user's issues")),
		mcp.WithBoolean("includeArchived", mcp.Description("Include archived issues in results")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of issues to return (default: 50)")),
	)
	addTool(getUserIssuesTool, GetUserIssuesHandler(linearClient))

	// Get Issue Tool (read-only)
	getIssueTool := mcp.NewTool("linear_get_issue",
		mcp.WithDescription("Retrieves a single Linear issue by its ID. Returns detailed information about the issue including title, description, priority, status, assignee, team, full comment history (including nested comments), related issues, and all attachments (pull requests, design files, documents, etc.)."),
		mcp.WithString("issueId", mcp.Required(), mcp.Description("ID of the issue to retrieve")),
	)
	addTool(getIssueTool, GetIssueHandler(linearClient))

	// Get Teams Tool (read-only)
	getTeamsTool := mcp.NewTool("linear_get_teams",
		mcp.WithDescription("Retrieves Linear teams with an optional name filter. If no name is provided, returns all teams. Returns team details including ID, name, and key."),
		mcp.WithString("name", mcp.Description("Optional team name filter. Returns teams whose names contain this string.")),
	)
	addTool(getTeamsTool, GetTeamsHandler(linearClient))

	// Create Issue Tool
	createIssueTool := mcp.NewTool("linear_create_issue",
		mcp.WithDescription("Creates a new Linear issue with specified details. Use this to create tickets for tasks, bugs, or feature requests. Returns the created issue's identifier and URL. Supports creating sub-issues and assigning labels."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("team", mcp.Required(), mcp.Description("Team identifier (key, UUID or name)")),
		mcp.WithString("description", mcp.Description("Issue description")),
		mcp.WithNumber("priority", mcp.Description("Priority (0-4)")),
		mcp.WithString("status", mcp.Description("Issue status")),
		mcp.WithString("parentIssue", mcp.Description("Optional parent issue ID or identifier (e.g., 'TEAM-123') to create a sub-issue")),
		mcp.WithString("labels", mcp.Description("Optional comma-separated list of label IDs or names to assign")),
	)
	addTool(createIssueTool, CreateIssueHandler(linearClient))

	// Update Issue Tool
	updateIssueTool := mcp.NewTool("linear_update_issue",
		mcp.WithDescription("Updates an existing Linear issue's properties. Use this to modify issue details like title, description, priority, or status. Requires the issue ID and accepts any combination of updatable fields. Returns the updated issue's identifier and URL."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Issue ID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithNumber("priority", mcp.Description("New priority (0-4)")),
		mcp.WithString("status", mcp.Description("New status")),
	)
	addTool(updateIssueTool, UpdateIssueHandler(linearClient))

	// Add Comment Tool
	addCommentTool := mcp.NewTool("linear_add_comment",
		mcp.WithDescription("Adds a comment to an existing Linear issue. Supports markdown formatting in the comment body. Can optionally specify a custom user name and avatar for the comment. Returns the created comment's details including its URL."),
		mcp.WithString("issueId", mcp.Required(), mcp.Description("ID of the issue to comment on")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment text in markdown format")),
		mcp.WithString("createAsUser", mcp.Description("Optional custom username to show for the comment")),
		mcp.WithString("displayIconUrl", mcp.Description("Optional avatar URL for the comment")),
	)
	addTool(addCommentTool, AddCommentHandler(linearClient))
}
