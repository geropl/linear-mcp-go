package main

import (
	"fmt"
	"os"

	"github.com/geropl/linear-mcp-go/pkg/server"
)

func main() {
	// Create the Linear MCP server
	linearServer, err := server.NewLinearMCPServer()
	if err != nil {
		fmt.Printf("Failed to create Linear MCP server: %v\n", err)
		os.Exit(1)
	}

	// Start the server
	if err := linearServer.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
