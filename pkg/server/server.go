package server

import (
	"fmt"
	"os"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/geropl/linear-mcp-go/pkg/tools"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

const (
	// ServerName is the name of the MCP server
	ServerName = "Linear MCP Server"
	// ServerVersion is the version of the MCP server
	ServerVersion = "1.6.2"
)

// LinearMCPServer represents the Linear MCP server
type LinearMCPServer struct {
	mcpServer    *mcpserver.MCPServer
	linearClient *linear.LinearClient
	writeAccess  bool // Controls whether write operations are enabled
}

// NewLinearMCPServer creates a new Linear MCP server
func NewLinearMCPServer(writeAccess bool) (*LinearMCPServer, error) {
	// Create the Linear client
	linearClient, err := linear.NewLinearClientFromEnv(ServerVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create Linear client: %w", err)
	}

	// Create the MCP server
	mcpServer := mcpserver.NewMCPServer(ServerName, ServerVersion)

	// Create the Linear MCP server
	server := &LinearMCPServer{
		mcpServer:    mcpServer,
		linearClient: linearClient,
		writeAccess:  writeAccess,
	}

	// Register tools
	RegisterTools(mcpServer, linearClient, writeAccess)

	// Register resources
	RegisterResources(mcpServer, linearClient)

	return server, nil
}

// Start starts the Linear MCP server
func (s *LinearMCPServer) Start() error {
	// Check if the Linear API key is set
	if os.Getenv("LINEAR_API_KEY") == "" {
		return fmt.Errorf("LINEAR_API_KEY environment variable is required")
	}

	// Start the server
	fmt.Printf("Starting %s v%s\n", ServerName, ServerVersion)
	return mcpserver.ServeStdio(s.mcpServer)
}

// GetMCPServer returns the underlying MCP server
func (s *LinearMCPServer) GetMCPServer() *mcpserver.MCPServer {
	return s.mcpServer
}

// GetLinearClient returns the Linear client
func (s *LinearMCPServer) GetLinearClient() *linear.LinearClient {
	return s.linearClient
}

// GetReadOnlyToolNames returns the names of all read-only tools
func GetReadOnlyToolNames() map[string]bool {
	return map[string]bool{
		"linear_search_issues":      true,
		"linear_get_user_issues":    true,
		"linear_get_issue":          true,
		"linear_get_issue_comments": true,
		"linear_get_teams":          true,
	}
}

// RegisterTools registers all Linear tools with the MCP server
func RegisterTools(s *mcpserver.MCPServer, linearClient *linear.LinearClient, writeAccess bool) {
	// Register tools, based on writeAccess
	addTool := func(tool mcp.Tool, handler mcpserver.ToolHandlerFunc) {
		if !writeAccess {
			if readOnly := GetReadOnlyToolNames()[tool.Name]; !readOnly {
				// Skip registering write tools if write access is disabled
				return
			}
		}
		s.AddTool(tool, handler)
	}

	// Register each tool
	addTool(tools.SearchIssuesTool, tools.SearchIssuesHandler(linearClient))
	addTool(tools.GetUserIssuesTool, tools.GetUserIssuesHandler(linearClient))
	addTool(tools.GetIssueTool, tools.GetIssueHandler(linearClient))
	addTool(tools.GetIssueCommentsTool, tools.GetIssueCommentsHandler(linearClient))
	addTool(tools.GetTeamsTool, tools.GetTeamsHandler(linearClient))
	addTool(tools.CreateIssueTool, tools.CreateIssueHandler(linearClient))
	addTool(tools.UpdateIssueTool, tools.UpdateIssueHandler(linearClient))
	addTool(tools.AddCommentTool, tools.AddCommentHandler(linearClient))
}
