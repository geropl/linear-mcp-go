package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linear-mcp-go",
	Short: "Linear MCP Server - A Model Context Protocol server for Linear",
	Long: `Linear MCP Server is a Model Context Protocol (MCP) server for Linear.
It provides tools for interacting with the Linear API through the MCP protocol,
enabling AI assistants to manage Linear issues and workflows.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}