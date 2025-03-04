# Active Context: Linear MCP Server

## Current Work Focus
The current focus is on improving the user experience by adding a setup command that simplifies the installation and configuration process. This includes:
1. Implementing a CLI framework with subcommands
2. Creating a setup command that automates the installation and configuration process
3. Supporting multiple AI assistants (starting with Cline)
4. Ensuring cross-platform compatibility

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

## Open Questions
1. Should we add support for more AI assistants in the setup command?
2. Do we need to add any additional validation steps for the API key?
3. Should we implement automatic updates for the binary?
4. How can we improve the error handling for network and file system operations?
