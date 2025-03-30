# Tool Standardization Implementation Guide

This document provides a detailed implementation plan for standardizing the Linear MCP Server tools according to the requirements outlined in [002-tool-standardization.md](./002-tool-standardization.md).

## Implementation Approach

We'll implement the standardization in phases, focusing on one rule at a time across all tools to ensure consistency:

1. First, create the necessary shared utility functions
2. Then, update each tool one by one, applying all three rules
3. Finally, update tests to verify the changes

## Shared Utility Functions

### 1. Identifier Resolution Functions

Create or update the following functions in `pkg/tools/common.go`:

```go
// resolveIssueIdentifier resolves an issue identifier (UUID or "TEAM-123") to a UUID
// This is an extension of the existing resolveParentIssueIdentifier function
func resolveIssueIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
    // Implementation similar to resolveParentIssueIdentifier
}

// resolveUserIdentifier resolves a user identifier (UUID, name, or email) to a UUID
func resolveUserIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
    // Implementation
}
```

### 2. Entity Rendering Functions

Create a new file `pkg/tools/rendering.go` with two types of formatting functions:

#### Full Entity Rendering Functions

```go
// formatIssue returns a consistently formatted full representation of an issue
func formatIssue(issue *linear.Issue) string {
    if issue == nil {
        return "Issue: Unknown"
    }
    
    var result strings.Builder
    result.WriteString(fmt.Sprintf("Issue: %s (UUID: %s)\n", issue.Identifier, issue.ID))
    result.WriteString(fmt.Sprintf("Title: %s\n", issue.Title))
    result.WriteString(fmt.Sprintf("URL: %s\n", issue.URL))
    
    // Add other required fields
    
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
    
    // Add other required fields
    
    return result.String()
}

// formatUser returns a consistently formatted full representation of a user
func formatUser(user *linear.User) string {
    if user == nil {
        return "User: Unknown"
    }
    
    var result strings.Builder
    result.WriteString(fmt.Sprintf("User: %s (UUID: %s)\n", user.Name, user.ID))
    
    // Add other required fields
    
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
    
    // Add other required fields
    
    return result.String()
}
```

#### Entity Identifier Rendering Functions

```go
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

## Detailed Implementation Tasks

### Phase 1: Create Shared Utility Functions

1. Update `pkg/tools/common.go`:
   - Refactor `resolveParentIssueIdentifier` to `resolveIssueIdentifier`
   - Add `resolveUserIdentifier`
   - Ensure all resolution functions follow the same pattern

2. Create `pkg/tools/rendering.go`:
   - Add formatting functions for each entity type
   - Ensure consistent formatting across all entity types

### Phase 2: Update Tools

For each tool, perform the following tasks:

1. Update the tool description to be concise
2. Update parameter handling to use the appropriate resolution functions
3. Update result formatting to use the rendering functions
4. Ensure retrieval methods include all fields that can be set in create/update methods

#### Implementing Rule 4: Field Superset for Retrieval Methods

For each entity type, follow these steps:

1. **Identify Modifiable Fields**:
   - Review all create and update methods to identify fields that can be set or modified
   - Create a comprehensive list of these fields for each entity type

2. **Categorize Retrieval Methods**:
   - **Detail Retrieval Methods** (e.g., `linear_get_issue`): Must include all fields
   - **Overview Retrieval Methods** (e.g., `linear_search_issues`, `linear_get_user_issues`): Only need metadata fields

3. **Update Detail Retrieval Methods**:
   - Ensure they include all fields that can be set in create/update methods
   - Modify the formatting functions to include all required fields

4. **Update Overview Retrieval Methods**:
   - Ensure they include key metadata fields (ID, title, status, priority, etc.)
   - No need to include full content like descriptions or comments

5. **Entity-Specific Considerations**:
   - **Issues**: 
     - `linear_get_issue` must include all fields from `linear_create_issue` and `linear_update_issue`
     - `linear_search_issues` and `linear_get_user_issues` only need metadata fields
   - **Comments**: 
     - Comments returned in `linear_get_issue` must include all necessary fields from `linear_add_comment`
     - Overview methods don't need to display comments
   - **Teams**: 
     - `linear_get_teams` should include all team fields that can be referenced

#### Tool-Specific Tasks

| Tool | Description Update | Identifier Resolution Update | Rendering Update |
|------|-------------------|----------------------------|-----------------|
| linear_create_issue | Remove parameter listing and result format explanation | Already uses resolution functions | Use formatIssue for result |
| linear_update_issue | Remove parameter listing and result format explanation | Update to use resolveIssueIdentifier | Use formatIssue for result |
| linear_search_issues | Remove parameter listing and result format explanation | Update to use resolveTeamIdentifier for teamId | Use formatIssue for each issue in results |
| linear_get_user_issues | Remove parameter listing and result format explanation | Add resolveUserIdentifier for userId | Use formatIssue for each issue in results |
| linear_get_issue | Remove parameter listing and result format explanation | Update to use resolveIssueIdentifier | Use formatIssue, formatTeam, formatUser, and formatComment |
| linear_add_comment | Remove parameter listing and result format explanation | Update to use resolveIssueIdentifier | Use formatIssue and formatComment |
| linear_get_teams | Remove parameter listing and result format explanation | No changes needed | Use formatTeam for each team in results |

### Phase 3: Update Tests

1. Update test fixtures to reflect the new formatting
2. Add tests for the new resolution functions
3. Verify that all tests pass with the new implementation

## Detailed Tracking Table

| Tool | Task | Status | Notes |
|------|------|--------|-------|
| Shared | Create resolveIssueIdentifier | Not Started | Extend from resolveParentIssueIdentifier |
| Shared | Create resolveUserIdentifier | Not Started | New function |
| Shared | Create rendering.go with formatting functions | Not Started | New file |
| linear_create_issue | Update description | Not Started | |
| linear_create_issue | Update result formatting | Not Started | |
| linear_update_issue | Update description | Not Started | |
| linear_update_issue | Add issue identifier resolution | Not Started | |
| linear_update_issue | Update result formatting | Not Started | |
| linear_search_issues | Update description | Not Started | |
| linear_search_issues | Add team identifier resolution | Not Started | |
| linear_search_issues | Update result formatting | Not Started | |
| linear_get_user_issues | Update description | Not Started | |
| linear_get_user_issues | Add user identifier resolution | Not Started | |
| linear_get_user_issues | Update result formatting | Not Started | |
| linear_get_issue | Update description | Not Started | |
| linear_get_issue | Add issue identifier resolution | Not Started | |
| linear_get_issue | Update result formatting | Not Started | |
| linear_get_issue | Ensure all fields from create/update are included | Not Started | Rule 4 implementation |
| linear_get_issue | Ensure all comment fields are included | Not Started | Rule 4 implementation |
| linear_get_user_issues | Ensure all relevant issue fields are included | Not Started | Rule 4 implementation |
| linear_search_issues | Ensure all relevant issue fields are included | Not Started | Rule 4 implementation |
| linear_get_teams | Ensure all team fields are included | Not Started | Rule 4 implementation |
| linear_add_comment | Update description | Not Started | |
| linear_add_comment | Add issue identifier resolution | Not Started | |
| linear_add_comment | Update result formatting | Not Started | |
| linear_get_teams | Update description | Not Started | |
| linear_get_teams | Update result formatting | Not Started | |
| Tests | Update test fixtures | Not Started | |
| Tests | Add tests for new resolution functions | Not Started | |
| Tests | Add tests for field superset compliance | Not Started | Rule 4 testing |

## Example Implementation: linear_create_issue

Here's an example of how the `linear_create_issue` tool would be updated:

### Before

```go
var CreateIssueTool = mcp.NewTool("linear_create_issue",
    mcp.WithDescription("Creates a new Linear issue with specified details. Use this to create tickets for tasks, bugs, or feature requests. Returns the created issue's identifier and URL. Supports creating sub-issues and assigning labels."),
    // ... parameters ...
)

// CreateIssueHandler handles the linear_create_issue tool
func CreateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... implementation ...
    
    // Return the result
    resultText := fmt.Sprintf("Created issue: %s\nTitle: %s\nURL: %s", issue.Identifier, issue.Title, issue.URL)
    return mcp.NewToolResultText(resultText), nil
}
```

### After

```go
var CreateIssueTool = mcp.NewTool("linear_create_issue",
    mcp.WithDescription("Creates a new Linear issue."),
    // ... parameters ...
)

// CreateIssueHandler handles the linear_create_issue tool
func CreateIssueHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... implementation ...
    
    // Return the result
    resultText := fmt.Sprintf("%s\nTitle: %s\nURL: %s", formatIssue(issue), issue.Title, issue.URL)
    return mcp.NewToolResultText(resultText), nil
}
```

## Timeline

| Phase | Estimated Duration | Dependencies |
|-------|-------------------|--------------|
| Phase 1: Create Shared Utility Functions | 1 day | None |
| Phase 2: Update Tools | 3 days | Phase 1 |
| Phase 3: Update Tests | 1 day | Phase 2 |

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking changes to API | High | Ensure backward compatibility or version the API |
| Test failures | Medium | Update test fixtures and add new tests |
| Inconsistent implementation | Medium | Review each tool implementation for consistency |

## Success Criteria

1. All tool descriptions are concise
2. All tools that reference Linear objects accept multiple identifier types
3. All tools render entities in a consistent format
4. Retrieval methods include all fields that can be set in create/update methods
5. All tests pass with the new implementation
6. Code review confirms consistency across all tools
