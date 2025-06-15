# Development Workflows: Linear MCP Server

## Git Workflow

### Branch Management
- **Main Branch**: Protected branch that contains stable, production-ready code
- **Feature Branches**: Development happens in feature branches created from main
- **Branch Naming**: Use descriptive names (e.g., `feature/setup-command`, `fix/rate-limiting`)
- **Pull Requests**: All changes must go through pull request review process
- **Branch Protection**: Main branch requires PR approval before merging

### Version Control Practices
- **Commit Messages**: Use clear, descriptive commit messages
- **Atomic Commits**: Each commit should represent a single logical change
- **Code Review**: All PRs require review before merging
- **Testing**: All tests must pass before merging

## Release Process

### Semantic Versioning
- **Version Format**: Follow semantic versioning with "v" prefix (e.g., v1.0.0, v1.2.3)
- **Version Components**:
  - Major: Breaking changes
  - Minor: New features (backward compatible)
  - Patch: Bug fixes (backward compatible)

### Release Steps
1. **Update Version**: Update the version in `pkg/server/server.go`
2. **Commit Changes**: Commit version update with descriptive message
3. **Create Tag**: Create and push a tag matching the version:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
4. **Automated Release**: GitHub Actions workflow automatically:
   - Builds binaries for Linux, macOS, and Windows
   - Runs all tests
   - Creates GitHub release
   - Uploads release assets

### GitHub Actions Workflow
- **Trigger**: Activated when tags matching "v*" pattern are pushed
- **Build Matrix**: Builds for multiple platforms simultaneously
- **Testing**: Runs full test suite before creating release
- **Asset Upload**: Automatically uploads compiled binaries to GitHub releases

## Development Commands

### Project Setup
```bash
# Clone the repository
git clone <repository-url>
cd linear-mcp-go

# Install dependencies (Go modules)
go mod download

# Verify setup
go build
```

### Building
```bash
# Build for current platform
go build

# Build with custom output name
go build -o linear-mcp-server

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o linear-mcp-server-linux

# Build for all platforms (manual)
GOOS=linux GOARCH=amd64 go build -o dist/linear-mcp-server-linux
GOOS=darwin GOARCH=amd64 go build -o dist/linear-mcp-server-darwin
GOOS=windows GOARCH=amd64 go build -o dist/linear-mcp-server-windows.exe
```

### Testing

#### Running Tests
```bash
# Run all tests with existing recordings
go test -v ./...

# Run tests for specific package
go test -v ./pkg/server

# Run specific test function
go test -v -run TestCreateIssueHandler ./pkg/server

# Run tests with coverage
go test -v -cover ./...
```

#### Recording New Test Fixtures
```bash
# Re-record tests (requires LINEAR_API_KEY environment variable)
go test -v -record=true ./...

# Re-record all tests including state-changing operations
go test -v -recordWrites=true ./...

# Re-record specific test
LINEAR_API_KEY=your_key go test -v -record=true -run TestSpecificFunction ./...
```

#### Test Environment Setup
```bash
# Set up environment for recording tests
export LINEAR_API_KEY=your_linear_api_key

# Run tests that modify state (use with caution)
go test -v -recordWrites=true ./...
```

### Running the Server

#### Development Mode
```bash
# Run server in read-only mode (safe for development)
./linear-mcp-go server

# Run server with write access (use with caution)
./linear-mcp-go server --write-access

# Run with environment variable
LINEAR_API_KEY=your_key ./linear-mcp-go server
```

#### Setup for AI Assistants
```bash
# Set up for Cline (default)
./linear-mcp-go setup --api-key=your_linear_api_key

# Set up with write access enabled
./linear-mcp-go setup --api-key=your_linear_api_key --write-access

# Set up for specific AI assistant (future)
./linear-mcp-go setup --api-key=your_key --tool=assistant_name
```

### Code Quality

#### Formatting
```bash
# Format all Go code
go fmt ./...

# Check formatting
gofmt -l .

# Format specific file
go fmt pkg/server/server.go
```

#### Linting
```bash
# Run go vet (built-in static analysis)
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

#### Dependencies
```bash
# Update dependencies
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

## Debugging and Troubleshooting

### Common Issues

#### API Key Problems
```bash
# Verify API key is set
echo $LINEAR_API_KEY

# Test API key manually
curl -H "Authorization: Bearer $LINEAR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"query": "{ viewer { id name } }"}' \
     https://api.linear.app/graphql
```

#### Build Issues
```bash
# Clean build cache
go clean -cache

# Rebuild from scratch
go clean && go build
```

#### Test Issues
```bash
# Clean test cache
go clean -testcache

# Run tests with verbose output
go test -v -x ./...
```

### Development Tips

#### Working with Test Fixtures
- Test fixtures are stored in `testdata/fixtures/`
- Golden files (expected outputs) are in `testdata/golden/`
- Use `-record=true` flag sparingly to avoid API quota exhaustion
- Always review recorded fixtures before committing

#### MCP Tool Development
- Register new tools in `pkg/server/tools.go`
- Follow existing patterns for parameter validation
- Add comprehensive test coverage for new tools
- Update documentation when adding new tools

#### Linear API Integration
- All API calls go through the Linear client in `pkg/linear/client.go`
- Add new API methods to the client rather than calling API directly from tools
- Handle GraphQL errors consistently
- Respect rate limits in all API interactions

## Continuous Integration

### GitHub Actions
- **Workflow File**: `.github/workflows/release.yml`
- **Triggers**: Push to main branch, pull requests, version tags
- **Jobs**: Test, build, release (for tags)
- **Platforms**: Linux, macOS, Windows

### Quality Gates
- All tests must pass
- Code must be properly formatted
- No linting errors
- Build must succeed on all target platforms
