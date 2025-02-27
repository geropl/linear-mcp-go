package server

import (
	"fmt"
	"os"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

const (
	// ServerName is the name of the MCP server
	ServerName = "Linear MCP Server"
	// ServerVersion is the version of the MCP server
	ServerVersion = "1.0.0"
)

// LinearMCPServer represents the Linear MCP server
type LinearMCPServer struct {
	mcpServer    *mcpserver.MCPServer
	linearClient *linear.LinearClient
}

// NewLinearMCPServer creates a new Linear MCP server
func NewLinearMCPServer() (*LinearMCPServer, error) {
	// Create the Linear client
	linearClient, err := linear.NewLinearClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to create Linear client: %w", err)
	}

	// Create the MCP server
	mcpServer := mcpserver.NewMCPServer(ServerName, ServerVersion)

	// Create the Linear MCP server
	server := &LinearMCPServer{
		mcpServer:    mcpServer,
		linearClient: linearClient,
	}

	// Register tools
	RegisterTools(mcpServer, linearClient)

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
