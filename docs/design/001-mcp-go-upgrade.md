# Design Doc: `mcp-go` Library Upgrade

**Author:** Cline
**Date:** 2025-06-28
**Status:** Proposed

## 1. Abstract

This document outlines the plan to upgrade the `mcp-go` dependency from its current version (`v0.18.0`) to the latest version provided. The upgrade is necessary to leverage new features, performance improvements, and bug fixes in the library. This document details the scope, identifies breaking changes, and provides a phased implementation plan to ensure a smooth and successful migration.

## 2. Background and Motivation

The Linear MCP Server currently relies on `mcp-go v0.18.0`. A new version of the library is available and has been cloned to `context/mcp-go` for reference. Analysis of the new version reveals significant API improvements and breaking changes. Upgrading will align the project with the latest MCP specification, improve developer experience through a more robust API, and ensure long-term maintainability.

## 3. Key Breaking Changes

Analysis of the new `mcp-go` library reveals three primary areas of breaking changes that will drive the refactoring effort.

### 3.1. Tool Definition and Schema

The API for defining tool schemas has been completely redesigned from a monolithic `WithInputSchema` method to a more granular, fluent builder pattern.

-   **Old Approach:** A single `WithInputSchema` call with a nested `mcp.Object` map.
-   **New Approach:** A series of top-level builder functions (`mcp.WithString`, `mcp.WithNumber`, etc.) are passed directly to `mcp.NewTool`. Property attributes like descriptions and requirement constraints are now handled by `PropertyOption` functions (`mcp.Description`, `mcp.Required`).

### 3.2. Tool Argument Parsing

The new `mcp.CallToolRequest` struct introduces a suite of type-safe methods for accessing tool arguments, deprecating the need for manual map access and type assertions.

-   **Old Approach:** `request.GetArguments()` returned a `map[string]any`, requiring developers to perform manual key lookups and type assertions.
-   **New Approach:** Methods like `request.RequireString("key")`, `request.GetInt("key", 0)`, and `request.BindArguments(&myStruct)` provide robust, type-safe access to arguments.

### 3.3. Tool Result Construction

The mechanism for returning results from a tool handler has been fundamentally changed. The simple, single-type helper functions have been replaced by a more flexible, multi-content structure.

-   **Old Approach:** Helper functions like `mcp.NewToolResultText(...)` and `mcp.NewToolResultError(...)` were used to create the result.
-   **New Approach:** The `mcp.CallToolResult` struct now contains a `Content` slice. Successful results are returned by populating this slice with `mcp.Content` objects (e.g., `mcp.TextContent`). Errors are handled by populating the `Content` slice with an error message and setting the `IsError` boolean flag to `true`.

## 4. Proposed Implementation Plan

The upgrade will be performed in four distinct phases to ensure a controlled and verifiable migration.

### Phase 1: Dependency Update

The first step is to instruct the Go compiler to use the new version of the library.

1.  **Update `go.mod`:** A `replace` directive will be added to `go.mod` to point the `github.com/mark3labs/mcp-go` module to the local directory `context/mcp-go`.
2.  **Tidy Dependencies:** `go mod tidy` will be executed to resolve the new dependency tree and remove unused old dependencies.

### Phase 2: Refactor Tool Definitions

All tool registration sites will be updated to use the new schema definition API.

1.  **Locate Tool Registrations:** All calls to `server.AddTool` will be identified.
2.  **Rewrite Schemas:** The `WithInputSchema` calls will be replaced with the new fluent builder pattern (`mcp.WithString`, `mcp.WithNumber`, etc.).

**Example Transformation:**

```go
// OLD
mcp.NewTool("get_issue").
    WithDescription("...").
    WithInputSchema(mcp.Object{
        "issue": mcp.Required(mcp.String()).WithDescription("..."),
    })

// NEW
mcp.NewTool("get_issue",
    mcp.WithDescription("..."),
    mcp.WithString("issue",
        mcp.Required(),
        mcp.Description("..."),
    ),
)
```

### Phase 3: Refactor Tool Handlers

This is the most significant phase, requiring changes to the logic within every tool handler.

1.  **Update Argument Parsing:** All manual argument map access will be replaced with the new type-safe methods on `mcp.CallToolRequest`.
2.  **Update Result Construction:** All `return` statements will be refactored to construct the `mcp.CallToolResult` struct with the appropriate `Content` slice and `IsError` flag.

**Example Transformation (Success):**

```go
// OLD
return mcp.NewToolResultText("Success!"), nil

// NEW
return &mcp.CallToolResult{
    Content: []mcp.Content{
        mcp.TextContent{Type: "text", Text: "Success!"},
    },
}, nil
```

**Example Transformation (Error):**

```go
// OLD
return mcp.NewToolResultError("Error message"), nil

// NEW
return &mcp.CallToolResult{
    IsError: true,
    Content: []mcp.Content{
        mcp.TextContent{Type: "text", Text: "Error message"},
    },
}, nil
```

### Phase 4: Compilation and Testing

The final phase is to verify the correctness of the refactored code.

1.  **Compile Project:** The entire project will be compiled to ensure there are no build errors.
2.  **Update and Run Tests:** The existing test suite will be executed. It is anticipated that tests will fail due to the API changes. Test code, particularly mock object creation and result assertions, will be updated to align with the new library version. All tests must pass before the upgrade is considered complete.

## 5. Risks and Mitigation

-   **Risk:** Unforeseen breaking changes.
    -   **Mitigation:** The phased approach allows for isolating and addressing issues systematically. The initial file analysis was thorough, but compilation will be the ultimate verification.
-   **Risk:** Logic errors introduced during refactoring.
    -   **Mitigation:** The existing test suite provides a safety net. All tests will be run and updated to ensure existing functionality is preserved.

## 6. Success Criteria

-   The project successfully compiles against the new `mcp-go` library version.
-   All existing tests pass after being updated for the new API.
-   The server runs correctly and all tools function as expected.
-   The `go.mod` file correctly references the new library path.

## 7. Progress Tracking

-   [x] **Phase 1: Dependency Update**
    -   [x] Update `go.mod` to `v0.32.0`.
    -   [x] Run `go mod tidy`.
-   [x] **Phase 2: Refactor Tool Definitions**
    -   [x] Refactor `linear_create_issue`
    -   [x] Refactor `linear_update_issue`
    -   [x] Refactor `linear_add_comment`
    -   [x] Refactor `linear_get_issue`
    -   [x] Refactor `linear_get_issue_comments`
    -   [x] Refactor `linear_get_teams`
    -   [x] Refactor `linear_get_user_issues`
    -   [x] Refactor `linear_search_issues`
-   [x] **Phase 3: Refactor Tool Handlers**
    -   [x] Refactor `createIssueHandler`
    -   [x] Refactor `updateIssueHandler`
    -   [x] Refactor `addCommentHandler`
    -   [x] Refactor `getIssueHandler`
    -   [x] Refactor `getIssueCommentsHandler`
    -   [x] Refactor `getTeamsHandler`
    -   [x] Refactor `getUserIssuesHandler`
    -   [x] Refactor `searchIssuesHandler`
-   [x] **Phase 4: Compilation and Testing**
    -   [x] Compile project successfully.
    -   [x] Update and pass all tests.
