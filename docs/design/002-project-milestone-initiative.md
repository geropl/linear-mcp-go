# Design Doc: Project, Milestone, and Initiative Support

## 1. Overview

This document outlines the plan to extend the Linear MCP Server with tools to read and manipulate `Project`, `ProjectMilestone`, and `Initiative` entities. This will enhance the server's capabilities, allowing AI assistants to manage a broader range of Linear workflows.

## 2. Guiding Principles

*   **Consistency**: The new tools and code will follow the existing architecture and design patterns of the server.
*   **Modularity**: Each entity will be implemented in a modular way, with clear separation between models, client methods, and tool handlers.
*   **User-Friendliness**: Tools will accept user-friendly identifiers (names, slugs) in addition to UUIDs, similar to the existing `issue` and `team` parameters.
*   **Testability**: All new functionality will be covered by unit tests using the existing `go-vcr` framework.

## 3. Architecture Changes

The existing architecture is well-suited for extension. The primary changes will be:

*   **`pkg/linear/models.go`**: Add new structs for `Project`, `ProjectMilestone`, `Initiative`, and their related types (e.g., `ProjectConnection`, `ProjectCreateInput`).
*   **`pkg/linear/client.go`**: Add new methods to the `LinearClient` for interacting with the new entities.
*   **`pkg/tools/`**: Create new files for each entity's tools (e.g., `project_tools.go`, `milestone_tools.go`, `initiative_tools.go`).
*   **`pkg/server/server.go`**: Update `RegisterTools` to include the new tools.

No fundamental changes to the core server logic or command structure are anticipated.

## 4. Implementation Plan

This section details the sub-tasks for implementing each handler.

### 4.1. Project Entity

#### 4.1.1. `linear_get_project` Handler

-   [x] **Model**: Define `Project` and `ProjectConnection` structs in `pkg/linear/models.go`.
-   [x] **Client**: Implement `GetProject(identifier string)` in `pkg/linear/client.go`.
    -   [x] Add a resolver function to handle UUIDs, names, and slug IDs.
    -   [x] Implement the GraphQL query to fetch a single project.
-   [x] **Tool**: Create `GetProjectTool` in `pkg/tools/project_tools.go`.
    -   [x] Define the tool with the name `linear_get_project`.
    -   [x] Add a description: "Get a single project by its identifier (ID, name, or slug)."
    -   [x] Define a required `project` string parameter.
-   [x] **Handler**: Implement `GetProjectHandler` in `pkg/tools/project_tools.go`.
    -   [x] Extract the `project` identifier from the request.
    -   [x] Call `linearClient.GetProject()` with the identifier.
    -   [x] Format the returned `Project` object into a user-friendly string.
    -   [x] Handle errors for not found projects.
-   [x] **Server**: Register the tool in `pkg/server/server.go`.
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

#### 4.1.2. `linear_search_projects` Handler

-   [x] **Client**: Implement `SearchProjects(query string)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL query for searching projects.
-   [x] **Tool**: Create `SearchProjectsTool` in `pkg/tools/project_tools.go`.
    -   [x] Define the tool with the name `linear_search_projects`.
    -   [x] Add a description: "Search for projects."
    -   [x] Define an optional `query` string parameter.
-   [x] **Handler**: Implement `SearchProjectsHandler` in `pkg/tools/project_tools.go`.
    -   [x] Extract the `query` from the request.
    -   [x] Call `linearClient.SearchProjects()`.
    -   [x] Format the list of `Project` objects.
-   [x] **Server**: Register the tool in `pkg/server/server.go`.
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

#### 4.1.3. `linear_create_project` Handler

-   [x] **Model**: Define `ProjectCreateInput` struct in `pkg/linear/models.go`.
-   [x] **Client**: Implement `CreateProject(input ProjectCreateInput)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL mutation to create a project.
-   [x] **Tool**: Create `CreateProjectTool` in `pkg/tools/project_tools.go`.
    -   [x] Define the tool with the name `linear_create_project`.
    -   [x] Add a description: "Create a new project."
    -   [x] Define required parameters: `name`, `teamIds`.
    -   [x] Define optional parameters: `description`, `leadId`, `startDate`, `targetDate`, etc.
-   [x] **Handler**: Implement `CreateProjectHandler` in `pkg/tools/project_tools.go`.
    -   [x] Build the `ProjectCreateInput` from the request parameters.
    -   [x] Call `linearClient.CreateProject()`.
    -   [x] Format the newly created `Project` object.
-   [x] **Server**: Register the tool in `pkg/server/server.go` (respecting `writeAccess`).
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

#### 4.1.4. `linear_update_project` Handler

-   [x] **Model**: Define `ProjectUpdateInput` struct in `pkg/linear/models.go`.
-   [x] **Client**: Implement `UpdateProject(id string, input ProjectUpdateInput)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL mutation to update a project.
-   [x] **Tool**: Create `UpdateProjectTool` in `pkg/tools/project_tools.go`.
    -   [x] Define the tool with the name `linear_update_project`.
    -   [x] Add a description: "Update an existing project."
    -   [x] Define a required `project` parameter.
    -   [x] Define optional parameters for updatable fields (`name`, `description`, etc.).
-   [x] **Handler**: Implement `UpdateProjectHandler` in `pkg/tools/project_tools.go`.
    -   [x] Resolve the `project` identifier to a UUID.
    -   [x] Build the `ProjectUpdateInput`.
    -   [x] Call `linearClient.UpdateProject()`.
    -   [x] Format the updated `Project` object.
-   [x] **Server**: Register the tool in `pkg/server/server.go` (respecting `writeAccess`).
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

### 4.2. ProjectMilestone Entity

#### 4.2.1. `linear_get_milestone` Handler

-   [x] **Model**: Define `ProjectMilestone` and `ProjectMilestoneConnection` structs in `pkg/linear/models.go`.
-   [x] **Client**: Implement `GetMilestone(id string)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL query to fetch a single milestone.
-   [x] **Tool**: Create `GetMilestoneTool` in `pkg/tools/milestone_tools.go`.
    -   [x] Define the tool with the name `linear_get_milestone`.
    -   [x] Add a description: "Get a single project milestone by its ID."
    -   [x] Define a required `milestoneId` string parameter.
-   [x] **Handler**: Implement `GetMilestoneHandler` in `pkg/tools/milestone_tools.go`.
    -   [x] Call `linearClient.GetMilestone()`.
    -   [x] Format the returned `ProjectMilestone` object.
-   [x] **Server**: Register the tool in `pkg/server/server.go`.
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

#### 4.2.2. `linear_create_milestone` Handler

-   [x] **Model**: Define `ProjectMilestoneCreateInput` in `pkg/linear/models.go`.
-   [x] **Client**: Implement `CreateMilestone(input ProjectMilestoneCreateInput)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL mutation.
-   [x] **Tool**: Create `CreateMilestoneTool` in `pkg/tools/milestone_tools.go`.
    -   [x] Define the tool with the name `linear_create_milestone`.
    -   [x] Add a description: "Create a new project milestone."
    -   [x] Define required parameters: `name`, `projectId`.
    -   [x] Define optional parameters: `description`, `targetDate`.
-   [x] **Handler**: Implement `CreateMilestoneHandler` in `pkg/tools/milestone_tools.go`.
    -   [x] Resolve `projectId`.
    -   [x] Build and call `linearClient.CreateMilestone()`.
    -   [x] Format the new `ProjectMilestone`.
-   [x] **Server**: Register the tool in `pkg/server/server.go` (respecting `writeAccess`).
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

### 4.3. Initiative Entity

#### 4.3.1. `linear_get_initiative` Handler

-   [x] **Model**: Define `Initiative` and `InitiativeConnection` structs in `pkg/linear/models.go`.
-   [x] **Client**: Implement `GetInitiative(identifier string)` in `pkg/linear/client.go`.
    -   [x] Add a resolver for UUID and name.
    -   [x] Implement the GraphQL query.
-   [x] **Tool**: Create `GetInitiativeTool` in `pkg/tools/initiative_tools.go`.
    -   [x] Define the tool with the name `linear_get_initiative`.
    -   [x] Add a description: "Get a single initiative by its identifier (ID or name)."
    -   [x] Define a required `initiative` string parameter.
-   [x] **Handler**: Implement `GetInitiativeHandler` in `pkg/tools/initiative_tools.go`.
    -   [x] Call `linearClient.GetInitiative()`.
    -   [x] Format the returned `Initiative` object.
-   [x] **Server**: Register the tool in `pkg/server/server.go`.
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

#### 4.3.2. `linear_create_initiative` Handler

-   [x] **Model**: Define `InitiativeCreateInput` in `pkg/linear/models.go`.
-   [x] **Client**: Implement `CreateInitiative(input InitiativeCreateInput)` in `pkg/linear/client.go`.
    -   [x] Implement the GraphQL mutation.
-   [x] **Tool**: Create `CreateInitiativeTool` in `pkg/tools/initiative_tools.go`.
    -   [x] Define the tool with the name `linear_create_initiative`.
    -   [x] Add a description: "Create a new initiative."
    -   [x] Define a required `name` parameter.
    -   [x] Define optional parameters like `description`.
-   [x] **Handler**: Implement `CreateInitiativeHandler` in `pkg/tools/initiative_tools.go`.
    -   [x] Build and call `linearClient.CreateInitiative()`.
    -   [x] Format the new `Initiative`.
-   [x] **Server**: Register the tool in `pkg/server/server.go` (respecting `writeAccess`).
-   [x] **Test**: Add test cases to `pkg/server/tools_test.go`.

### 4.4.1 Re-record tests

-   [x] Coordinate with user to prepare test data and re-record requests
-   [x] Iterate over each new test cases and validate that it works

### 4.4.2 Issues detected during testing

-   [x] get_project: filters by slugId ('e1153169a428'), but takes slug ('created-test-project-e1153169a428') as input. Should split string by '-' and used the last element
-   [x] search_projects: Cannot find multiple projects (prefix search only?)
-   [x] issues:
    -   [x] display milestone and project association
    -   [x] allow to update milestone and project association
-   [x] projects:
    -   [x] display initiative association
    -   [x] allow to update initiative association
-   [x] create_milestone: date parsing issue

## 5. Testing Strategy

Recording new test requests and data is the final step of this effort. We rely on the relevant test cases being added beforehand during each step.

*   **Unit Tests**: Tests are implemented in `pkg/server/tools_test.go`.
*   **Fixtures**: Use `go-vcr` to record real API interactions for each new client method. Fixtures will be stored in `testdata/fixtures/`.
*   **Golden Files**: Expected outputs for each test case will be stored in `testdata/golden/`.
*   **Coverage**: Ensure that both success and error paths are tested for each handler. This includes testing with invalid identifiers, missing parameters, and API errors.

## 6. Future Considerations

*   **Full CRUD**: This plan covers the most common operations. Full CRUD (including delete) can be added later if needed.
*   **Relationships**: Add tools for managing relationships between these entities (e.g., adding a project to an initiative).
*   **Resources**: Expose these new entities as MCP resources for easy reference in other tools.

## 7. Critique and Improvements

Based on a review of the initial implementation, several areas for improvement have been identified to enhance tool symmetry and consistency.

### 7.1. Enhance `linear_update_project`

-   [x] **Symmetry**: The `linear_update_project` tool should support the same set of optional parameters as `linear_create_project`.
-   [x] **Task**: Add `leadId`, `startDate`, `targetDate`, and `teamIds` as optional parameters to the `UpdateProjectTool` definition in `pkg/tools/project_tools.go`.
-   [x] **Task**: Update the `UpdateProjectHandler` to handle these new parameters.
-   [x] **Task**: Update the `linearClient.UpdateProject` method and `ProjectUpdateInput` struct to include these fields.
-   [ ] **Test**: Add test cases to verify that each field can be updated individually and in combination.

### 7.2. Add `linear_update_milestone`

-   [x] **Symmetry**: Add a `linear_update_milestone` tool to provide full CRUD operations for milestones.
-   [x] **Model**: Define `ProjectMilestoneUpdateInput` struct in `pkg/linear/models.go`.
-   [x] **Client**: Implement `UpdateMilestone(id string, input ProjectMilestoneUpdateInput)` in `pkg/linear/client.go`.
-   [x] **Tool**: Create `UpdateMilestoneTool` in `pkg/tools/milestone_tools.go` with a required `milestone` parameter and optional `name`, `description`, and `targetDate` parameters.
-   [x] **Handler**: Implement `UpdateMilestoneHandler` in `pkg/tools/milestone_tools.go`.
-   [x] **Server**: Register the new tool in `pkg/server/server.go`.
-   [ ] **Test**: Add test cases for the new tool.

### 7.3. Add `linear_update_initiative`

-   [x] **Symmetry**: Add a `linear_update_initiative` tool to provide full CRUD operations for initiatives.
-   [x] **Model**: Define `InitiativeUpdateInput` struct in `pkg/linear/models.go`.
-   [x] **Client**: Implement `UpdateInitiative(id string, input InitiativeUpdateInput)` in `pkg/linear/client.go`.
-   [x] **Tool**: Create `UpdateInitiativeTool` in `pkg/tools/initiative_tools.go` with a required `initiative` parameter and optional `name` and `description` parameters.
-   [x] **Handler**: Implement `UpdateInitiativeHandler` in `pkg/tools/initiative_tools.go`.
-   [x] **Server**: Register the new tool in `pkg/server/server.go`.
-   [ ] **Test**: Add test cases for the new tool.

### 7.4. Standardize Identifier Parameters

-   [x] **Consistency**: Ensure all `get` and `update` tools use a consistent and user-friendly identifier parameter.
-   [x] **Task**: In `pkg/tools/milestone_tools.go`, rename the `milestoneId` parameter of `GetMilestoneTool` to `milestone`.
-   [x] **Task**: Update the `GetMilestoneHandler` to use the `milestone` parameter.
-   [x] **Task**: Enhance `linearClient.GetMilestone` to resolve the milestone by name in addition to ID, similar to how `GetProject` works.
-   [ ] **Test**: Update existing tests for `linear_get_milestone` and add new tests for name-based resolution.
