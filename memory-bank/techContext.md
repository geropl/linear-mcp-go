# Technical Context: Linear MCP Server

## Technologies Used

### Programming Language
- **Go**: Version 1.23.6
  - Chosen for its performance, strong typing, and concurrency support
  - Excellent standard library for HTTP requests and JSON handling

### Key Libraries
1. **github.com/mark3labs/mcp-go v0.8.5**
   - MCP protocol implementation for Go
   - Provides server, tool registration, and request/response handling

2. **gopkg.in/dnaeon/go-vcr.v4 v4.0.2**
   - HTTP interaction recording and playback for testing
   - Allows tests to run without actual API calls

### APIs
- **Linear API**
  - REST API for Linear issue tracking system
  - Requires API key authentication
  - Has rate limiting constraints

## Development Setup

### Prerequisites
- Go 1.23 or higher
- Linear API key

### Environment Variables
- `LINEAR_API_KEY`: Required for authentication with Linear API

### Build Process
```bash
# Build the server
go build

# Run the server in read-only mode (default)
./linear-mcp-go

# Run the server with write operations enabled
./linear-mcp-go --write-access
```

### Command-Line Flags
- `--write-access`: Controls whether write operations are enabled (default: false)
  - When false, write tools (`linear_create_issue`, `linear_update_issue`, `linear_add_comment`) are disabled
  - When true, all tools are available

### Testing
```bash
# Run tests with existing recordings
go test -v ./...

# Re-record tests (requires LINEAR_API_KEY)
go test -v -record=true ./...

# Re-record all tests including state-changing ones
go test -v -recordWrites=true ./...
```

## Technical Constraints

### Linear API Limitations
1. **Rate Limiting**
   - Linear API has rate limits that must be respected
   - The server implements rate limiting to prevent quota exhaustion

2. **Authentication**
   - Requires API key passed via environment variable
   - No support for OAuth or other authentication methods

### MCP Protocol Constraints
1. **Communication Channel**
   - MCP server communicates via stdin/stdout
   - No HTTP or other network protocols for MCP communication

2. **Tool Schema**
   - Tools must define their parameters using MCP schema format
   - Parameters can be required or optional with descriptions

### Deployment Constraints
1. **Binary Distribution**
   - Server is distributed as a compiled binary
   - Binaries should be available for major platforms (Linux, macOS, Windows)

2. **Environment**
   - Requires environment variables to be set
   - No configuration file support currently

## Dependencies

### Direct Dependencies
```
github.com/mark3labs/mcp-go v0.8.5
gopkg.in/dnaeon/go-vcr.v4 v4.0.2
```

### Indirect Dependencies
```
github.com/google/go-cmp v0.7.0
github.com/google/uuid v1.6.0
gopkg.in/yaml.v3 v3.0.1
```

## Version Information
- **Server Version**: 1.0.0 (defined in pkg/server/server.go)
- **Go Version**: 1.23.6 (defined in go.mod)
- **MCP SDK Version**: 0.8.5

## Build and Release Process
- GitHub Actions workflow for automated testing and releases
- Releases are created when tags matching "v*" are pushed
- Binaries are built for Linux, macOS, and Windows

## Future Technical Considerations
1. **Configuration File Support**
   - Could add support for configuration files instead of just environment variables

2. **Additional Linear API Features**
   - More Linear API endpoints could be exposed as MCP tools

3. **Improved Error Handling**
   - More detailed error messages and recovery strategies

4. **Metrics and Logging**
   - Add structured logging and metrics collection
