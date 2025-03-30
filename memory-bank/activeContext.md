# Active Context: Linear MCP Server

## Current Work Focus
The current focus is on enhancing the functionality and user experience of the Linear MCP Server. This includes:
1. Improving the user experience by adding a setup command that simplifies the installation and configuration process
2. Enhancing the Linear API integration with support for more advanced features
3. Supporting multiple AI assistants (starting with Cline)
4. Ensuring cross-platform compatibility
5. Expanding the capabilities of existing MCP tools

## Recent Changes
1. Implemented CLI framework with subcommands:
   - Added the Cobra library for command-line handling
   - Restructured the main.go file to support subcommands
   - Created a root command that serves as the base for all subcommands
   - Moved the existing server functionality to a server subcommand

2. Created a setup command:
   - Implemented a setup command that automates the installation and configuration process
   - Added support for the Cline AI assistant
   - Implemented binary discovery and download functionality
   - Added configuration file management for AI assistants

3. Updated documentation:
   - Updated README.md with information about the new setup command
   - Added examples of how to use the setup command
   - Clarified the usage of the server command

4. Enhanced `linear_create_issue` tool:
   - Added support for creating sub-issues by specifying a parent issue ID or identifier (e.g., "TEAM-123")
   - Added support for assigning labels during issue creation using label IDs or names
   - Implemented resolution functions for parent issue identifiers and label names
   - Updated the Linear client to handle these new parameters
   - Added test cases and fixtures for the new functionality
   - Updated documentation to reflect the new capabilities

5. Fixed JSON unmarshaling issue with Labels field:
   - Updated the `Issue` struct in `models.go` to change the `Labels` field from `[]Label` to `*LabelConnection`
   - Added a new `LabelConnection` struct to match the structure returned by the Linear API
   - Updated test fixtures and golden files to reflect the changes
   - Added a new test case for creating issues with team key

6. Fixed label resolution issue:
   - Updated the GraphQL query in `GetLabelsByName` function to change the `$teamId` parameter type from `ID!` to `String!`
   - Re-recorded test fixtures for label-related tests
   - Updated golden files to reflect the new error messages
   - All tests now pass successfully

7. Fixed parent issue identifier resolution:
   - Updated the `GetIssueByIdentifier` function to split the identifier (e.g., "TEAM-123") into team key and number parts
   - Modified the GraphQL query to use the team key and number in the filter instead of the full identifier
   - Added proper error handling for invalid identifier formats
   - Added a new test case for creating sub-issues using human-readable identifiers
   - All tests now pass successfully

## Next Steps
1. **Testing the Setup Command**:
   - Test the setup command on different platforms (Linux, macOS, Windows)
   - Verify that the configuration files are correctly created
   - Ensure that the binary download works correctly

2. **Adding Support for More AI Assistants**:
   - Research other AI assistants that could benefit from Linear integration
   - Implement support for these assistants in the setup command
   - Update documentation with information about the new assistants

3. **Future Enhancements**:
   - Add more Linear API features as MCP tools
   - Improve error handling and reporting
   - Add configuration file support for the server

## Active Decisions and Considerations

### CLI Framework Selection
- **Decision**: Use the Cobra library for command-line handling
  - **Rationale**: Cobra is a widely used library for Go CLI applications with good documentation and community support
  - **Alternatives Considered**: urfave/cli, flag package
  - **Implications**: Provides a consistent way to handle subcommands and flags

### Setup Command Design
- **Decision**: Implement a setup command that automates the installation and configuration process
  - **Rationale**: Simplifies the user experience by automating manual steps
  - **Alternatives Considered**: Keeping the bash script, creating a separate tool
  - **Implications**: Users can easily set up the server for use with AI assistants

### AI Assistant Support
- **Decision**: Start with Cline support and design for extensibility
  - **Rationale**: Cline is the primary target, but the design should allow for adding more assistants
  - **Alternatives Considered**: Supporting only Cline, supporting multiple assistants from the start
  - **Implications**: The code is structured to easily add support for more assistants in the future

### Binary Management
- **Decision**: Check for existing binary before downloading
  - **Rationale**: Avoids unnecessary downloads if the binary is already installed
  - **Alternatives Considered**: Always downloading the latest version
  - **Implications**: Faster setup process for users who already have the binary

### Configuration File Management
- **Decision**: Merge new settings with existing settings
  - **Rationale**: Preserves user's existing configuration while adding the Linear MCP server
  - **Alternatives Considered**: Overwriting the entire file
  - **Implications**: Users can have multiple MCP servers configured

### Linear Issue Creation Enhancement
- **Decision**: Enhance the `linear_create_issue` tool with support for user-friendly identifiers
  - **Rationale**: Provides more flexibility and better user experience when creating issues
  - **Alternatives Considered**: Requiring UUIDs only, creating separate tools for different identifier types
  - **Implications**: Users can create issues with more intuitive identifiers without needing to look up UUIDs

### Identifier Resolution Implementation
- **Decision**: Implement separate resolution functions for parent issues and labels
  - **Rationale**: Keeps the code modular and easier to maintain
  - **Alternatives Considered**: Implementing a generic resolution function, handling resolution in the handler directly
  - **Implications**: Code is more maintainable and easier to extend for future enhancements

### JSON Structure Handling
- **Decision**: Update the `Issue` struct to match the nested structure returned by the Linear API
  - **Rationale**: Ensures proper JSON unmarshaling of API responses
  - **Alternatives Considered**: Custom unmarshaling logic, flattening the structure in the client
  - **Implications**: More robust handling of API responses and fewer unmarshaling errors

### GraphQL Parameter Type Correction
- **Decision**: Update the GraphQL query parameter types to match the Linear API expectations
  - **Rationale**: Ensures proper validation of GraphQL queries by the Linear API
  - **Alternatives Considered**: Custom error handling for API validation errors
  - **Implications**: More reliable API requests and fewer validation errors

### Parent Issue Identifier Resolution
- **Decision**: Split the identifier into team key and number parts for the GraphQL query
  - **Rationale**: The Linear API doesn't support searching for issues by the full identifier directly
  - **Alternatives Considered**: Using a different API endpoint, implementing a custom search function
  - **Implications**: More reliable resolution of human-readable identifiers to UUIDs

## Open Questions
1. Should we add support for more AI assistants in the setup command?
2. Do we need to add any additional validation steps for the API key?
3. Should we implement automatic updates for the binary?
4. How can we improve the error handling for network and file system operations?
