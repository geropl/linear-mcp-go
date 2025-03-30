# Product Requirements Document: Issue Comments Pagination

## Overview

This document outlines the requirements for splitting the comment functionality from the `get_issue` tool and implementing a new `get_issue_comments` tool with proper pagination support.

## Background

Currently, the `get_issue` tool returns all comments for an issue as part of its response. This approach has several limitations:

1. For issues with many comments, the response can become very large
2. There's no way to paginate through comments
3. There's no way to specifically retrieve replies to a particular comment
4. The client receives more data than might be needed if they're only interested in the issue details

## Requirements

### Functional Requirements

1. Create a new tool named `linear_get_issue_comments` that retrieves comments for a Linear issue
2. Implement pagination for comments to handle potentially long conversations
3. Support retrieving comments from specific threads (top-level or replies to a specific comment)
4. Modify the existing `get_issue` tool to remove the comments section and add a reference to the new tool
5. Ensure backward compatibility by maintaining the same comment formatting style

### Technical Requirements

1. The new tool should accept the following parameters:
   - `issue`: (required) ID or identifier of the issue to retrieve comments for
   - `thread`: (optional) UUID of the parent comment to retrieve replies for
   - `limit`: (optional) Maximum number of comments to return (default: 10)
   - `after`: (optional) Cursor for pagination to get comments after this point

2. The response should include:
   - Basic issue information (identifier, UUID)
   - Thread information (root or parent comment UUID)
   - List of comments with user, timestamp, and content
   - Pagination information (has more comments, next cursor)
   - Indication if comments have replies

3. Update the Linear client to support the new functionality:
   - Add a `GetIssueComments` method that supports pagination and thread filtering
   - Define a `GetIssueCommentsInput` struct for the parameters
   - Define a `PaginatedCommentConnection` struct for the response

4. Register the new tool in the server and add it to the read-only tools list

## Implementation Details

### API Changes

The Linear GraphQL API already supports pagination and filtering for comments. We'll use the following query structure:

```graphql
query GetIssueComments($issueId: String!, $parentId: String, $first: Int!, $after: String) {
  issue(id: $issueId) {
    comments(
      first: $first,
      after: $after,
      filter: { parent: { id: { eq: $parentId } } }
    ) {
      nodes {
        id
        body
        createdAt
        user {
          id
          name
        }
        childCount
      }
      pageInfo {
        hasNextPage
        endCursor
      }
      totalCount
    }
  }
}
```

### Implementation Steps

| Step | Description | Status |
|------|-------------|--------|
| 1 | Add `GetIssueCommentsInput` and `PaginatedCommentConnection` structs to models.go | ✅ |
| 2 | Implement `GetIssueComments` method in the Linear client | ✅ |
| 3 | Create the `get_issue_comments.go` file with tool definition and handler | ✅ |
| 4 | Update `get_issue.go` to remove comments section and add reference to new tool | ✅ |
| 5 | Register the new tool in server.go | ✅ |
| 6 | Add the new tool to the read-only tools list | ✅ |
| 7 | Test the implementation | ✅ |

## Usage Examples

### Example 1: Get top-level comments for an issue

```json
{
  "issue": "TEAM-123",
  "limit": 5
}
```

### Example 2: Get replies to a specific comment

```json
{
  "issue": "TEAM-123",
  "thread": "comment-uuid-here",
  "limit": 10
}
```

### Example 3: Paginate through comments

```json
{
  "issue": "TEAM-123",
  "after": "cursor-from-previous-response",
  "limit": 10
}
```

## Benefits

1. **Improved Performance**: Clients can request only the comments they need, reducing payload size
2. **Better User Experience**: Support for pagination allows handling large comment threads efficiently
3. **More Flexibility**: Ability to navigate through specific comment threads
4. **Cleaner API**: Separation of concerns between issue details and comments

## Conclusion

The implementation of the `linear_get_issue_comments` tool with pagination support will significantly improve the handling of issue comments, especially for issues with extensive discussions. This change aligns with best practices for API design by providing more granular control over data retrieval and reducing unnecessary data transfer.
