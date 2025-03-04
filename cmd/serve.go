package cmd

import (
	"fmt"
	"os"

	"github.com/geropl/linear-mcp-go/pkg/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Linear MCP server",
	Long: `Start the Linear MCP server that listens for MCP requests on stdin/stdout.
The server provides tools for interacting with the Linear API through the MCP protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeAccess, _ := cmd.Flags().GetBool("write-access")

		// Create the Linear MCP server
		linearServer, err := server.NewLinearMCPServer(writeAccess)
		if err != nil {
			fmt.Printf("Failed to create Linear MCP server: %v\n", err)
			os.Exit(1)
		}

		// Start the server
		if err := linearServer.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add flags to the serve command
	serveCmd.Flags().Bool("write-access", false, "Enable tools that modify Linear data (create/update issues, add comments)")
}
