package cmd

import (
	"fmt"
	"os"
	"strings"

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
		writeAccessChanged := cmd.Flags().Changed("write-access")

		// Check LINEAR_WRITE_ACCESS environment variable if flag wasn't explicitly set
		if !writeAccessChanged {
			if envWriteAccess := os.Getenv("LINEAR_WRITE_ACCESS"); envWriteAccess != "" {
				envValue := strings.ToLower(strings.TrimSpace(envWriteAccess))
				if envValue == "true" {
					writeAccess = true
				} else if envValue == "false" {
					writeAccess = false
				}
				// If the env var is set to something other than "true" or "false", ignore it and use default
			}
		}

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
