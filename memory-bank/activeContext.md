# Active Context: Linear MCP Server

## Current Work Focus
The current focus is on implementing a GitHub Actions workflow for automated testing and release creation. This workflow will:
1. Run tests on all pushes to ensure code quality
2. Create releases automatically when version tags are pushed
3. Build binaries for multiple platforms (Linux, macOS, Windows)

## Recent Changes
1. Implemented write access control feature:
   - Added a `--write-access` flag (default: false) to control whether write operations are enabled
   - Modified server to conditionally register write tools based on this flag
   - Updated documentation and install script to support the new flag
   - Write operations (`linear_create_issue`, `linear_update_issue`, `linear_add_comment`) are disabled by default

2. Created a GitHub Actions workflow file (`.github/workflows/release.yml`) that:
   - Runs on push events to main branch and tags matching "v*"
   - Tests the build on all pushes
   - Creates a GitHub release when a tag matching "v*" is pushed
   - Builds binaries for Linux, macOS, and Windows
   - Uploads the binaries as release assets

3. Fixed GitHub Actions workflow permissions issue:
   - Added explicit `contents: write` permission to the `create-release` job
   - This resolves the "Resource not accessible by integration" error
   - Allows the workflow to create releases and upload assets

## Next Steps
1. **Testing the Workflow**:
   - Push the changes to GitHub
   - Create a test tag (e.g., v0.1.0) to verify the release process
   - Verify that binaries are correctly built and attached to the release

2. **Documentation Updates**:
   - Update README.md with information about the release process
   - Add a CONTRIBUTING.md file with development guidelines

3. **Future Enhancements**:
   - Add more Linear API features as MCP tools
   - Improve error handling and reporting
   - Add configuration file support

## Active Decisions and Considerations

### Write Access Control Design
- **Decision**: Implement write access control with a command-line flag
  - **Rationale**: Provides a simple way to control access to write operations
  - **Alternatives Considered**: Environment variable, configuration file
  - **Implications**: Users can easily control whether write operations are enabled

### Release Workflow Design
- **Decision**: Use GitHub Actions for automated releases
  - **Rationale**: GitHub Actions provides a simple, integrated way to automate the build and release process
  - **Alternatives Considered**: CircleCI, Travis CI, custom scripts
  - **Implications**: Requires GitHub repository, relies on GitHub's infrastructure

### Binary Distribution
- **Decision**: Build binaries for Linux, macOS, and Windows
  - **Rationale**: These are the major platforms where users might run the server
  - **Considerations**: ARM architectures not currently supported but could be added in the future

### Version Tagging
- **Decision**: Use semantic versioning with tags starting with "v" (e.g., v1.0.0)
  - **Rationale**: Standard practice for Go projects, easy to parse in automation
  - **Implications**: Need to ensure version in code (ServerVersion) matches release tags

### Testing Strategy
- **Decision**: Run tests on all pushes, including PRs
  - **Rationale**: Ensures code quality before merging
  - **Considerations**: Tests use recorded HTTP interactions, so they don't require a Linear API key

## Open Questions
1. Should we add Docker container builds to the release process?
2. Do we need to add any additional validation steps before creating releases?
3. Should we implement automatic changelog generation based on commit messages?
