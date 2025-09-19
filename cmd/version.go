package cmd

import (
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/server"
	"github.com/spf13/cobra"
)

// Build information - these will be set at build time
var (
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information for the Linear MCP server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Linear MCP Server %s\n", server.ServerVersion)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}