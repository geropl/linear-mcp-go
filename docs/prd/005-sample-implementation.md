# Sample Implementation

This document provides sample implementations for the key components of the tool standardization effort. These samples can be used as references when implementing the actual changes.

## 1. Shared Utility Functions

### rendering.go

```go
package tools

import (
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
)

// Full Entity Rendering Functions

// formatIssue returns a consistently formatted full representation of an issue
func formatIssue(issue *linear.Issue) string {
	if issue == nil {
		return "Issue: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Issue: %s (UUID: %s)\n", issue.Identifier, issue.ID))
	result.WriteString(fmt.Sprintf("Title: %s\n", issue.Title))
	result.WriteString(fmt.Sprintf("URL: %s\n", issue.URL))
	
	if issue.Description != "" {
		result.WriteString(fmt.Sprintf("Description: %s\n", issue.Description))
	}
	
	priorityStr := "None"
	if issue.Priority > 0 {
		priorityStr = fmt.Sprintf("%d", issue.Priority)
	}
	result.WriteString(fmt.Sprintf("Priority: %s\n", priorityStr))
	
	statusStr := "None"
	if issue.Status != "" {
		statusStr = issue.Status
	} else if issue.State != nil {
		statusStr = issue.State.Name
	}
	result.WriteString(fmt.Sprintf("Status: %s\n", statusStr))
	
	if issue.Assignee != nil {
		result.WriteString(fmt.Sprintf("Assignee: %s\n", formatUserIdentifier(issue.Assignee)))
	}
	
	if issue.Team != nil {
		result.WriteString(fmt.Sprintf("Team: %s\n", formatTeamIdentifier(issue.Team)))
	}
	
	return result.String()
}

// formatTeam returns a consistently formatted full representation of a team
func formatTeam(team *linear.Team) string {
	if team == nil {
		return "Team: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Team: %s (UUID: %s)\n", team.Name, team.ID))
	result.WriteString(fmt.Sprintf("Key: %s\n", team.Key))
	
	return result.String()
}

// formatUser returns a consistently formatted full representation of a user
func formatUser(user *linear.User) string {
	if user == nil {
		return "User: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("User: %s (UUID: %s)\n", user.Name, user.ID))
	
	if user.Email != "" {
		result.WriteString(fmt.Sprintf("Email: %s\n", user.Email))
	}
	
	return result.String()
}

// formatComment returns a consistently formatted full representation of a comment
func formatComment(comment *linear.Comment) string {
	if comment == nil {
		return "Comment: Unknown"
	}
	
	userName := "Unknown"
	if comment.User != nil {
		userName = comment.User.Name
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Comment by %s (UUID: %s)\n", userName, comment.ID))
	result.WriteString(fmt.Sprintf("Body: %s\n", comment.Body))
	
	if comment.CreatedAt != nil {
		result.WriteString(fmt.Sprintf("Created: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	
	return result.String()
}

// Entity Identifier Rendering Functions

// formatIssueIdentifier returns a consistently formatted identifier for an issue
func formatIssueIdentifier(issue *linear.Issue) string {
	if issue == nil {
		return "Issue: Unknown"
	}
	return fmt.Sprintf("Issue: %s (UUID: %s)", issue.Identifier, issue.ID)
}

// formatTeamIdentifier returns a consistently formatted identifier for a team
func formatTeamIdentifier(team *linear.Team) string {
	if team == nil {
		return "Team: Unknown"
	}
	return fmt.Sprintf("Team: %s (UUID: %s)", team.Name, team.ID)
}

// formatUserIdentifier returns a consistently formatted identifier for a user
func formatUserIdentifier(user *linear.User) string {
	if user == nil {
		return "User: Unknown"
	}
	return fmt.Sprintf("User: %s (UUID: %s)", user.Name, user.ID)
}

// formatCommentIdentifier returns a consistently formatted identifier for a comment
func formatCommentIdentifier(comment *linear.Comment) string {
	if comment == nil {
		return "Comment: Unknown"
	}
	
	userName := "Unknown"
	if comment.User != nil {
		userName = comment.User.Name
	}
	
	return fmt.Sprintf("Comment by %s (UUID: %s)", userName, comment.ID)
}
```

### Updated common.go

```go
// resolveIssueIdentifier resolves an issue identifier (UUID or "TEAM-123") to a UUID
func resolveIssueIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find an issue by identifier
	issue, err := linearClient.GetIssueByIdentifier(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to resolve issue identifier '%s': %v", identifier, err)
	}

	return issue.ID, nil
}

// resolveUserIdentifier resolves a user identifier (UUID, name, or email) to a UUID
func resolveUserIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find a user by name or email
	users, err := linearClient.GetUsers()
	if err != nil {
		return "", fmt.Errorf("failed to get users: %v", err)
	}

	// First try exact match on name or email
	for _, user := range users {
		if user.Name == identifier || user.Email == identifier {
			return user.ID, nil
		}
	}

	// If no exact match, try case-insensitive match
	identifierLower := strings.ToLower(identifier)
	for _, user := range users {
		if strings.ToLower(user.Name) == identifierLower || strings.ToLower(user.Email) == identifierLower {
			return user.ID, nil
		}
	}

	return "", fmt.Errorf("no user found with identifier '%s'", identifier)
}
```

## 2. Updated Tool Examples

### linear_create_issue

```go
// Before
var CreateIssueTool = mcp.NewTool("linear_create_issue",
	mcp.WithDescription("Creates a new Linear issue with specified details. Use this to create tickets for tasks, bugs, or feature requests. Returns the created issue's identifier and URL. Supports creating sub-issues and assigning labels."),
	// ... parameters ...
)

// After
var CreateIssueTool = mcp.NewTool("linear_create_issue",
	mcp.WithDescription("Creates a new Linear issue."),
	// ... parameters ...
)

// Before (result formatting)
resultText := fmt.Sprintf("Created issue: %s\nTitle: %s\nURL: %s", issue.Identifier, issue.Title, issue.URL)

// After (result formatting)
resultText := fmt.Sprintf("Created %s\nTitle: %s\nURL: %s", formatIssue(issue), issue.Title, issue.URL)
```

### linear_update_issue

```go
// Before
var UpdateIssueTool = mcp.NewTool("linear_update_issue",
	mcp.WithDescription("Updates an existing Linear issue's properties. Use this to modify issue details like title, description, priority, or status. Requires the issue ID and accepts any combination of updatable fields. Returns the updated issue's identifier and URL."),
	// ... parameters ...
)

// After
var UpdateIssueTool = mcp.NewTool("linear_update_issue",
	mcp.WithDescription("Updates an existing Linear issue."),
	// ... parameters ...
)

// Before (parameter handling)
id, ok := args["id"].(string)
if !ok || id == "" {
	return mcp.NewToolResultError("id must be a non-empty string"), nil
}

// After (parameter handling)
issueID, ok := args["id"].(string)
if !ok || issueID == "" {
	return mcp.NewToolResultError("id must be a non-empty string"), nil
}

// Resolve the issue identifier
id, err := resolveIssueIdentifier(linearClient, issueID)
if err != nil {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve issue: %v", err)), nil
}

// Before (result formatting)
resultText := fmt.Sprintf("Updated issue %s\nURL: %s", issue.Identifier, issue.URL)

// After (result formatting)
resultText := fmt.Sprintf("Updated %s\nURL: %s", formatIssue(issue), issue.URL)
```

### linear_get_issue

```go
// Before
var GetIssueTool = mcp.NewTool("linear_get_issue",
	mcp.WithDescription("Retrieves a single Linear issue by its ID. Returns detailed information about the issue including title, description, priority, status, assignee, team, full comment history (including nested comments), related issues, and all attachments (pull requests, design files, documents, etc.)."),
	// ... parameters ...
)

// After
var GetIssueTool = mcp.NewTool("linear_get_issue",
	mcp.WithDescription("Retrieves a single Linear issue."),
	// ... parameters ...
)

// Before (parameter handling)
issueID, ok := args["issueId"].(string)
if !ok || issueID == "" {
	return mcp.NewToolResultError("issueId must be a non-empty string"), nil
}

// After (parameter handling)
issueIdentifier, ok := args["issueId"].(string)
if !ok || issueIdentifier == "" {
	return mcp.NewToolResultError("issueId must be a non-empty string"), nil
}

// Resolve the issue identifier
issueID, err := resolveIssueIdentifier(linearClient, issueIdentifier)
if err != nil {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve issue: %v", err)), nil
}

// Before (result formatting)
resultText := fmt.Sprintf("Issue %s: %s\n", issue.Identifier, issue.Title)
// ... more formatting ...

// After (result formatting)
// Use the full formatIssue function for the main entity
resultText := formatIssue(issue)

// When referencing related entities, use the identifier formatting functions
if issue.Assignee != nil {
    resultText += fmt.Sprintf("Assignee: %s\n", formatUserIdentifier(issue.Assignee))
}

if issue.Team != nil {
    resultText += fmt.Sprintf("Team: %s\n", formatTeamIdentifier(issue.Team))
}

// For related issues, use the identifier formatting
if issue.Relations != nil && len(issue.Relations.Nodes) > 0 {
    resultText += "\nRelated Issues:\n"
    for _, relation := range issue.Relations.Nodes {
        if relation.RelatedIssue != nil {
            resultText += fmt.Sprintf("- %s\n  RelationType: %s\n", 
                formatIssueIdentifier(relation.RelatedIssue),
                relation.Type)
        }
    }
}
```

## 3. Testing Examples

### Test for resolveIssueIdentifier

```go
func TestResolveIssueIdentifier(t *testing.T) {
	// Create test client
	client, cleanup := linear.NewTestClient(t, "resolve_issue_identifier", true)
	defer cleanup()

	// Test cases
	tests := []struct {
		name       string
		identifier string
		wantErr    bool
	}{
		{
			name:       "Valid UUID",
			identifier: "1c2de93f-4321-4015-bfde-ee893ef7976f",
			wantErr:    false,
		},
		{
			name:       "Valid identifier",
			identifier: "TEST-10",
			wantErr:    false,
		},
		{
			name:       "Invalid identifier",
			identifier: "NONEXISTENT-123",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveIssueIdentifier(client, tt.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveIssueIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("resolveIssueIdentifier() returned empty UUID")
			}
		})
	}
}
```

### Test for Formatting Functions

```go
func TestFormatIssue(t *testing.T) {
	tests := []struct {
		name  string
		issue *linear.Issue
		want  string
	}{
		{
			name: "Valid issue",
			issue: &linear.Issue{
				ID:         "1c2de93f-4321-4015-bfde-ee893ef7976f",
				Identifier: "TEST-10",
			},
			want: "Issue: TEST-10 (UUID: 1c2de93f-4321-4015-bfde-ee893ef7976f)",
		},
		{
			name:  "Nil issue",
			issue: nil,
			want:  "Issue: Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatIssue(tt.issue); got != tt.want {
				t.Errorf("formatIssue() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## 4. Implementation Strategy

1. Start with creating the shared utility functions in `rendering.go` and updating `common.go`
2. Implement the changes for one tool (e.g., `linear_create_issue`) as a reference
3. Review the reference implementation to ensure it meets all requirements
4. Apply the same patterns to all remaining tools
5. Update tests to verify the changes

## 5. Potential Challenges and Solutions

### Challenge: Backward Compatibility
Changing the output format of tools could break existing integrations that parse the output.

**Solution:** Consider versioning the API or providing a compatibility mode.

### Challenge: Test Fixtures
Updating the output format will require updating all test fixtures.

**Solution:** Use the `--golden` flag to update all golden files at once after implementing the changes.

### Challenge: Consistent Implementation
Ensuring consistency across all tools can be challenging.

**Solution:** Create a code review checklist to verify that each tool follows the same patterns.
