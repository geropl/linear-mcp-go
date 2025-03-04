# Linear MCP Server

A Model Context Protocol (MCP) server for Linear, written in Go. This server provides tools for interacting with the Linear API through the MCP protocol.

## Features

- Create, update, and search Linear issues
- Get issues assigned to a user
- Add comments to issues
- Retrieve team information
- Rate-limited API requests to respect Linear's API limits

## Prerequisites

- Go 1.23 or higher
- Linear API key

## Installation

### From Releases

Pre-built binaries are available for Linux, macOS, and Windows on the [GitHub Releases page](https://github.com/geropl/linear-mcp-go/releases).

1. Download the appropriate binary for your platform
2. Make it executable (Linux/macOS):

```bash
chmod +x linear-mcp-go-*
```

3. Run the binary as described in the Usage section

### Automated

```
# Download linux binary for the latest release
RELEASE=$(curl -s https://api.github.com/repos/geropl/linear-mcp-go/releases/latest)
DOWNLOAD_URL=$(echo $RELEASE | jq -r '.assets[] | select(.name | contains("linux")) | .browser_download_url')
curl -L -o ./linear-mcp-go $DOWNLOAD_URL

# Setup the mcp server (.gitpod.yml, dotfiles repo, etc.)
./linear-mcp-go setup --tool=cline
```

## Usage

### Running the Server

1. Set your Linear API key as an environment variable:

```bash
export LINEAR_API_KEY=your_linear_api_key
```

2. Run the server:

```bash
# Run in read-only mode (default)
./linear-mcp-go serve

# Run with write access enabled
./linear-mcp-go serve --write-access
```

The server will start and listen for MCP requests on stdin/stdout.

### Setting Up for AI Assistants

The `setup` command automates the installation and configuration process for various AI assistants:

```bash
# Set your Linear API key as an environment variable
export LINEAR_API_KEY=your_linear_api_key

# Set up for Cline (default)
./linear-mcp-go setup

# Set up with write access enabled
./linear-mcp-go setup --write-access

# Set up with auto-approval for read-only tools
./linear-mcp-go setup --auto-approve=allow-read-only

# Set up with specific tools auto-approved
./linear-mcp-go setup --auto-approve=linear_get_issue,linear_search_issues

# Set up with write access and auto-approval for read-only tools
./linear-mcp-go setup --write-access --auto-approve=allow-read-only

# Set up for a different tool (only "cline" supported for now)
./linear-mcp-go setup --tool=cline
```

This command:
1. Checks if the Linear MCP binary is already installed
2. Copies the current binary to the installation directory if needed
3. Configures the AI assistant to use the Linear MCP server
4. Sets up auto-approval for specified tools if requested

The `--auto-approve` flag can be used to specify which tools should be auto-approved in the Cline configuration:
- `--auto-approve=allow-read-only`: Auto-approves all read-only tools (`linear_search_issues`, `linear_get_user_issues`, `linear_get_issue`, `linear_get_teams`)
- `--auto-approve=tool1,tool2,...`: Auto-approves the specified comma-separated list of tools

Currently supported AI assistants:
- Cline (VSCode extension)

By default, the server runs in read-only mode, which means the following tools are disabled:
- `linear_create_issue`
- `linear_update_issue`
- `linear_add_comment`

To enable these tools, use the `--write-access=true` flag.

## Available Tools

### linear_create_issue

Creates a new Linear issue with specified details.

**Parameters:**
- `title` (required): Issue title
- `teamId` (required): Team ID
- `description`: Issue description
- `priority`: Priority (0-4)
- `status`: Issue status

### linear_update_issue

Updates an existing Linear issue's properties.

**Parameters:**
- `id` (required): Issue ID
- `title`: New title
- `description`: New description
- `priority`: New priority (0-4)
- `status`: New status

### linear_search_issues

Searches Linear issues using flexible criteria.

**Parameters:**
- `query`: Optional text to search in title and description
- `teamId`: Filter by team ID
- `status`: Filter by status name (e.g., 'In Progress', 'Done')
- `assigneeId`: Filter by assignee's user ID
- `labels`: Filter by label names (comma-separated)
- `priority`: Filter by priority (1=urgent, 2=high, 3=normal, 4=low)
- `estimate`: Filter by estimate points
- `includeArchived`: Include archived issues in results (default: false)
- `limit`: Max results to return (default: 10)

### linear_get_user_issues

Retrieves issues assigned to a specific user or the authenticated user.

**Parameters:**
- `userId`: Optional user ID. If not provided, returns authenticated user's issues
- `includeArchived`: Include archived issues in results
- `limit`: Maximum number of issues to return (default: 50)

### linear_get_issue

Retrieves a single Linear issue by its ID.

**Parameters:**
- `issueId` (required): ID of the issue to retrieve

### linear_add_comment

Adds a comment to an existing Linear issue.

**Parameters:**
- `issueId` (required): ID of the issue to comment on
- `body` (required): Comment text in markdown format
- `createAsUser`: Optional custom username to show for the comment
- `displayIconUrl`: Optional avatar URL for the comment

### linear_get_teams

Retrieves Linear teams with an optional name filter.

**Parameters:**
- `name`: Optional team name filter. Returns teams whose names contain this string.

## Test
Tests are implemented using [`go-vcr`](https://github.com/dnaeon/go-vcr), and executed against https://linear.app/linear-mcp-go-test.

### Execute tests

Using the existing recordings (cassettes):
```
go test -v ./...
```

#### Re-recording test:

Requires `LINEAR_API_KEY` to be set.

```
go test -v -record=true ./...
```
This will update all tests that don't alter remote state.


```
go test -v -recordWrites=true ./...
```
This will re-run all tests, including some that might alter the outcome of other tests cases, which might require further manual work to adjust.

```
go test -v -golden=true ./...
```
Updates all .golden fields.

## Release Process

The project uses GitHub Actions for automated testing and releases:

1. All pushes to the main branch and pull requests are automatically tested
2. When a tag matching the pattern `v*` (e.g., `v1.0.0`) is pushed, a new release is automatically created
3. Binaries for Linux, macOS, and Windows are built and attached to the release

To create a new release:

1. Update the version in `pkg/server/server.go`
2. Commit the changes
3. Create and push a tag matching the version:
```bash
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will automatically create a release with the appropriate binaries.

## License

MIT
