# Product Requirements Document: Linear MCP Server Tool Standardization

## Overview
This document outlines the requirements for standardizing the Linear MCP Server tools according to a set of consistent rules. These rules aim to improve user experience, code maintainability, and consistency across all tools.

## Background
The Linear MCP Server currently provides several tools for interacting with the Linear API through the Model Context Protocol (MCP). While these tools are functional, they lack consistency in their descriptions, parameter handling, and result formatting. This inconsistency can lead to confusion for users and maintenance challenges for developers.

## Requirements

### 1. Concise Tool Descriptions
**Current State:** Tool descriptions are verbose and often contain parameter listings and result format explanations.

**Requirement:** Tool descriptions should be concise and focus only on the tool's purpose and functionality. They should not:
- List parameters (these are already defined in the schema)
- Explain the result format (this should be consistent across tools)

### 2. Flexible Object Identifier Resolution
**Current State:** Some tools (like `create_issue`) already support resolving different types of identifiers (UUID, name, key) to the underlying UUID, but this is not consistent across all tools.

**Requirement:** All input arguments that reference Linear objects should:
- Accept multiple forms of identification (UUID, name, key)
- Resolve these identifiers to the underlying UUID using appropriate resolution functions
- Use consistent resolution methods across all tools

### 3. Consistent Entity Rendering
**Current State:** Entity rendering in results varies across tools, with inconsistent formatting and information hierarchy.

**Requirement:** Tools fetching the same entities should:
- Emit results using the same format
- Add any tool-specific additional fields at the bottom of the result
- Use code reuse between different tools to ensure consistency

This requirement has two distinct parts:

1. **Full Entity Rendering**:
   - When displaying an entity as the primary subject of a response, use a consistent format with all required fields
   - Implement shared formatting functions (e.g., `formatIssue`, `formatTeam`) that include all necessary information
   - Example: When getting an issue with `linear_get_issue`, display the full issue details

2. **Entity Identifier Rendering**:
   - When referencing an entity from another entity, use a consistent, concise identifier format
   - Implement shared identifier formatting functions (e.g., `formatIssueIdentifier`, `formatTeamIdentifier`)
   - Always display entity identifiers in the format: "[Most descriptive field] (UUID: ...)"
     - For issues: "Issue: TEST-10 (UUID: ...)"
     - For teams: "Team: Test Team (UUID: ...)"

### 4. Field Superset for Retrieval Methods
**Current State:** Retrieval methods may not include all fields that can be set through create and update methods.

**Requirement:** The fields rendered on retrieval methods should follow these rules:

1. **Detail Retrieval Methods** (e.g., `linear_get_issue`):
   - Must include the complete superset of all fields that can be set through create and update methods
   - Ensures users can always view any field they can modify
   - Prevents "hidden" fields that can be set but not retrieved

2. **Overview Retrieval Methods** (e.g., `linear_search_issues`, `linear_get_user_issues`):
   - Only need to include key metadata fields (ID, title, status, priority, etc.)
   - Not required to display full content like descriptions or comments
   - Focus on providing sufficient information for selection and identification

This distinction ensures each tool provides an appropriate level of detail for its purpose while maintaining consistency.

For example:
- `linear_get_issue` must include all fields from `linear_create_issue` and `linear_update_issue`
- `linear_search_issues` only needs to show metadata fields, not full descriptions or comments

## Implementation Plan

### Phase 1: Analysis and Documentation
1. Review all existing tools to identify:
   - Current description formats
   - Object identifier resolution methods
   - Entity rendering patterns

2. Create a detailed tracking table for all tools, documenting:
   - Current state for each rule
   - Required changes
   - Implementation status

### Phase 2: Implementation
1. Create shared utility functions for:
   - Entity rendering
   - Identifier resolution (extending existing functions as needed)

2. Update each tool to:
   - Revise descriptions to be concise
   - Use shared identifier resolution functions
   - Implement consistent entity rendering

3. Update tests to verify:
   - Identifier resolution works correctly
   - Entity rendering is consistent

### Phase 3: Validation and Documentation
1. Verify all tools meet the new standards
2. Update documentation to reflect the changes
3. Create examples demonstrating the new consistent behavior

## Tool Standardization Tracking

| Tool | Rule 1: Concise Description | Rule 2: Flexible Identifiers | Rule 3: Consistent Rendering | Status |
|------|----------------------------|------------------------------|------------------------------|--------|
| linear_create_issue | ❌ Too verbose | ✅ Supports team, parent issue, and label resolution | ❌ Custom format | Not Started |
| linear_update_issue | ❌ Too verbose | ❌ Only accepts issue ID | ❌ Custom format | Not Started |
| linear_search_issues | ❌ Too verbose | ❌ Only accepts teamId | ❌ Custom format | Not Started |
| linear_get_user_issues | ❌ Too verbose | ❌ Only accepts userId | ❌ Custom format | Not Started |
| linear_get_issue | ❌ Too verbose | ❌ Only accepts issueId | ❌ Custom format | Not Started |
| linear_add_comment | ❌ Too verbose | ❌ Only accepts issueId | ❌ Custom format | Not Started |
| linear_get_teams | ❌ Too verbose | ✅ No identifiers needed | ❌ Custom format | Not Started |

## Implementation Details

### Common Identifier Resolution Functions
We'll need to create or extend the following resolution functions:
1. `resolveTeamIdentifier` (already exists)
2. `resolveIssueIdentifier` (extend from `resolveParentIssueIdentifier`)
3. `resolveUserIdentifier` (new)
4. `resolveLabelIdentifiers` (already exists)

### Common Entity Rendering Functions
We'll need to create the following rendering functions:
1. `formatIssue` - For consistent issue rendering
2. `formatTeam` - For consistent team rendering
3. `formatUser` - For consistent user rendering
4. `formatComment` - For consistent comment rendering

### Code Structure Changes
1. Move common functions to a shared package
2. Create a new `rendering.go` file for entity formatting functions
3. Update all tool handlers to use these shared functions

## Success Criteria
1. All tool descriptions are concise and focused on functionality
2. All tools that reference Linear objects accept multiple identifier types
3. All tools render entities in a consistent format
4. Code reuse is maximized through shared functions
5. All tests pass with the new implementation

## Next Steps
1. Begin with updating one tool to serve as a reference implementation
2. Review the reference implementation to ensure it meets all requirements
3. Apply the same patterns to all remaining tools
4. Update tests and documentation

## Conclusion
Implementing these standardization rules will improve the user experience, make the codebase more maintainable, and ensure consistency across all tools. This will make the Linear MCP Server more professional and easier to use.
