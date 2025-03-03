package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/geropl/linear-mcp-go/pkg/server"
)

func main() {
	// Parse command-line flags
	writeAccess := flag.Bool("write-access", false, "Enable tools that modify Linear data (create/update issues, add comments)")
	flag.Parse()

	// Create the Linear MCP server
	linearServer, err := server.NewLinearMCPServer(*writeAccess)
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
