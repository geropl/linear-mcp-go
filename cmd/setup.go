package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/server"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up the Linear MCP server for use with an AI assistant",
	Long: `Set up the Linear MCP server for use with an AI assistant.
This command installs the Linear MCP server and configures it for use with the specified AI assistant tool(s).
Currently supported tools: cline, roo-code, claude-code, ona`,
	Run: func(cmd *cobra.Command, args []string) {
		toolParam, _ := cmd.Flags().GetString("tool")
		writeAccess, _ := cmd.Flags().GetBool("write-access")
		autoApprove, _ := cmd.Flags().GetString("auto-approve")
		projectPath, _ := cmd.Flags().GetString("project-path")

		// Check if the Linear API key is provided in the environment
		apiKey := os.Getenv("LINEAR_API_KEY")
		if apiKey == "" {
			fmt.Println("Error: LINEAR_API_KEY environment variable is required")
			fmt.Println("Please set it before running the setup command:")
			fmt.Println("export LINEAR_API_KEY=your_linear_api_key")
			os.Exit(1)
		}

		// Create the MCP servers directory if it doesn't exist
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			os.Exit(1)
		}

		mcpServersDir := filepath.Join(homeDir, "mcp-servers")
		if err := os.MkdirAll(mcpServersDir, 0755); err != nil {
			fmt.Printf("Error creating MCP servers directory: %v\n", err)
			os.Exit(1)
		}

		// Check if the Linear MCP binary is already on the path
		binaryPath, found := checkBinary()
		if !found {
			fmt.Printf("Linear MCP binary not found on path, copying current binary to '%s'...\n", binaryPath)
			err := copySelfToBinaryPath(binaryPath)
			if err != nil {
				fmt.Printf("Error copying Linear MCP binary: %v\n", err)
				os.Exit(1)
			}
		}

		// Process each tool
		tools := strings.Split(toolParam, ",")
		hasErrors := false

		for _, t := range tools {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}

			fmt.Printf("Setting up tool: %s\n", t)

			// Set up the tool-specific configuration
			var err error
			switch strings.ToLower(t) {
			case "cline":
				err = setupCline(binaryPath, apiKey, writeAccess, autoApprove)
			case "roo-code":
				err = setupRooCode(binaryPath, apiKey, writeAccess, autoApprove)
			case "claude-code":
				err = setupClaudeCode(binaryPath, apiKey, writeAccess, autoApprove, projectPath)
			case "ona":
				err = setupOna(binaryPath, apiKey, writeAccess, autoApprove, projectPath)
			default:
				fmt.Printf("Unsupported tool: %s\n", t)
				fmt.Println("Currently supported tools: cline, roo-code, claude-code, ona")
				hasErrors = true
				continue
			}

			if err != nil {
				fmt.Printf("Error setting up %s: %v\n", t, err)
				hasErrors = true
			} else {
				fmt.Printf("Linear MCP server successfully set up for %s\n", t)
			}
		}

		if hasErrors {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Add flags to the setup command
	setupCmd.Flags().String("tool", "cline", "The AI assistant tool(s) to set up for (comma-separated, e.g., cline,roo-code,claude-code,ona)")
	setupCmd.Flags().Bool("write-access", false, "Enable write operations (default: false)")
	setupCmd.Flags().String("auto-approve", "", "Comma-separated list of tool names to auto-approve, or 'allow-read-only' to auto-approve all read-only tools")
	setupCmd.Flags().String("project-path", "", "The project path(s) for claude-code project-scoped configuration (comma-separated for multiple projects, or empty to register to user scope for all projects)")
}

// checkBinary checks if the Linear MCP binary is already on the path
func checkBinary() (string, bool) {
	// Try to find the binary on the path
	path, err := exec.LookPath("linear-mcp-go")
	if err == nil {
		fmt.Printf("Found Linear MCP binary at %s\n", path)
		return path, true
	}

	// Check if the binary exists in the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}

	binaryPath := filepath.Join(homeDir, "mcp-servers", "linear-mcp-go")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	if _, err := os.Stat(binaryPath); err == nil {
		fmt.Printf("Found Linear MCP binary at %s\n", binaryPath)
		return binaryPath, true
	}

	return binaryPath, false
}

// copySelfToBinaryPath copies the current executable to the specified path
func copySelfToBinaryPath(binaryPath string) error {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check if the destination is the same as the source
	absExecPath, _ := filepath.Abs(execPath)
	absDestPath, _ := filepath.Abs(binaryPath)
	if absExecPath == absDestPath {
		return nil // Already in the right place
	}

	// Copy the file
	sourceFile, err := os.Open(execPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(binaryPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Make the binary executable
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	fmt.Printf("Linear MCP server installed successfully at %s\n", binaryPath)
	return nil
}

// getOnaConfigPath determines the configuration file path for Ona
func getOnaConfigPath(projectPath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Default to current working directory
	baseDir := cwd

	// If project path is specified, use the first one
	if projectPath != "" {
		paths := strings.Split(projectPath, ",")
		trimmedPath := strings.TrimSpace(paths[0])
		if trimmedPath != "" {
			baseDir = trimmedPath
			// If absolute path doesn't exist, treat as relative to current directory
			if filepath.IsAbs(trimmedPath) {
				if _, err := os.Stat(trimmedPath); os.IsNotExist(err) {
					relativePath := strings.TrimPrefix(trimmedPath, "/")
					baseDir = filepath.Join(cwd, relativePath)
				}
			}
		}
	}

	return filepath.Join(baseDir, ".gitpod", "mcp-config.json"), nil
}

// setupOna sets up the Linear MCP server for Ona
func setupOna(binaryPath, apiKey string, writeAccess bool, autoApprove, projectPath string) error {
	configPath, err := getOnaConfigPath(projectPath)
	if err != nil {
		return err
	}

	// Create the .gitpod directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create .gitpod directory: %w", err)
	}

	// Prepare server arguments
	serverArgs := []string{"serve"}
	serverArgs = append(serverArgs, fmt.Sprintf("--write-access=%t", writeAccess))

	// Create the linear server configuration
	linearServerConfig := map[string]interface{}{
		"name":    "linear",
		"command": binaryPath,
		"args":    serverArgs,
	}

	// Add environment variables
	if apiKey != "" {
		linearServerConfig["env"] = map[string]string{
			"LINEAR_API_KEY": apiKey,
		}
	}

	// Read existing configuration or create new one
	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read existing ona config: %w", err)
		}
		// Initialize with empty structure if file doesn't exist
		config = map[string]interface{}{
			"servers": map[string]interface{}{},
		}
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing ona config: %w", err)
		}
		// Ensure servers field exists
		if config["servers"] == nil {
			config["servers"] = map[string]interface{}{}
		}
	}

	// Get or create servers map
	servers, ok := config["servers"].(map[string]interface{})
	if !ok {
		servers = map[string]interface{}{}
		config["servers"] = servers
	}

	// Add/update the linear server configuration
	servers["linear"] = linearServerConfig

	// Write the updated configuration
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ona config: %w", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write ona config: %w", err)
	}

	fmt.Printf("Ona MCP configuration updated at %s\n", configPath)
	return nil
}

// setupTool sets up the Linear MCP server for a specific tool
func setupTool(toolName string, binaryPath, apiKey string, writeAccess bool, autoApprove string, configDir string) error {
	// Create the config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	serverArgs := []string{"serve"}
	if writeAccess {
		serverArgs = append(serverArgs, "--write-access=true")
	}

	// Process auto-approve flag
	autoApproveTools := []string{}
	if autoApprove != "" {
		if autoApprove == "allow-read-only" {
			// Get the list of read-only tools
			for k := range server.GetReadOnlyToolNames() {
				autoApproveTools = append(autoApproveTools, k)
			}
		} else {
			// Split comma-separated list
			for _, tool := range strings.Split(autoApprove, ",") {
				trimmedTool := strings.TrimSpace(tool)
				if trimmedTool != "" {
					autoApproveTools = append(autoApproveTools, trimmedTool)
				}
			}
		}
	}

	// Create the MCP settings file
	settingsPath := filepath.Join(configDir, "cline_mcp_settings.json")
	newSettings := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"linear": map[string]interface{}{
				"command":     binaryPath,
				"args":        serverArgs,
				"env":         map[string]string{"LINEAR_API_KEY": apiKey},
				"disabled":    false,
				"autoApprove": autoApproveTools,
			},
		},
	}

	// Check if the settings file already exists
	var settings map[string]interface{}
	if _, err := os.Stat(settingsPath); err == nil {
		// Read the existing settings
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return fmt.Errorf("failed to read existing settings: %w", err)
		}

		// Parse the existing settings
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse existing settings: %w", err)
		}

		// Merge the new settings with the existing settings
		if mcpServers, ok := settings["mcpServers"].(map[string]interface{}); ok {
			mcpServers["linear"] = newSettings["mcpServers"].(map[string]interface{})["linear"]
		} else {
			settings["mcpServers"] = newSettings["mcpServers"]
		}
	} else {
		// Use the new settings
		settings = newSettings
	}

	// Write the settings to the file
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	fmt.Printf("%s MCP settings updated at %s\n", toolName, settingsPath)
	return nil
}

// setupCline sets up the Linear MCP server for Cline
func setupCline(binaryPath, apiKey string, writeAccess bool, autoApprove string) error {
	// Determine the Cline config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	case "linux":
		configDir = filepath.Join(homeDir, ".vscode-server", "data", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	case "windows":
		configDir = filepath.Join(homeDir, "AppData", "Roaming", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return setupTool("Cline", binaryPath, apiKey, writeAccess, autoApprove, configDir)
}

// setupRooCode sets up the Linear MCP server for Roo Code
func setupRooCode(binaryPath, apiKey string, writeAccess bool, autoApprove string) error {
	// Determine the Roo Code config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	case "linux":
		configDir = filepath.Join(homeDir, ".vscode-server", "data", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	case "windows":
		configDir = filepath.Join(homeDir, "AppData", "Roaming", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return setupTool("Roo Code", binaryPath, apiKey, writeAccess, autoApprove, configDir)
}

// setupClaudeCode sets up the Linear MCP server for Claude Code
func setupClaudeCode(binaryPath, apiKey string, writeAccess bool, autoApprove, projectPath string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("claude-code is only supported on Linux")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".claude.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", filepath.Dir(configPath), err)
	}

	// Use flexible map structure to preserve all existing settings
	var settings map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read claude code settings: %w", err)
		}
		// Initialize with empty structure if file doesn't exist
		settings = map[string]interface{}{
			"projects": map[string]interface{}{},
		}
	} else {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse claude code settings: %w", err)
		}
		// Ensure projects field exists
		if settings["projects"] == nil {
			settings["projects"] = map[string]interface{}{}
		}
	}

	serverArgs := []string{"serve"}
	if writeAccess {
		serverArgs = append(serverArgs, "--write-access=true")
	}

	autoApproveTools := []string{}
	if autoApprove != "" {
		if autoApprove == "allow-read-only" {
			for k := range server.GetReadOnlyToolNames() {
				autoApproveTools = append(autoApproveTools, k)
			}
		} else {
			for _, tool := range strings.Split(autoApprove, ",") {
				trimmedTool := strings.TrimSpace(tool)
				if trimmedTool != "" {
					autoApproveTools = append(autoApproveTools, trimmedTool)
				}
			}
		}
	}

	linearServerConfig := map[string]interface{}{
		"type":        "stdio",
		"command":     binaryPath,
		"args":        serverArgs,
		"env":         map[string]string{"LINEAR_API_KEY": apiKey},
		"disabled":    false,
		"autoApprove": autoApproveTools,
	}

	if projectPath == "" {
		// Register to user-scoped mcpServers (applies to all projects)
		if err := registerLinearToUserScope(settings, linearServerConfig); err != nil {
			return fmt.Errorf("failed to register Linear MCP server to user scope: %w", err)
		}
		fmt.Printf("Registered Linear MCP server to user scope (applies to all projects)\n")
	} else {
		// Parse comma-separated project paths and register to specific projects
		var targetProjects []string
		for _, path := range strings.Split(projectPath, ",") {
			trimmedPath := strings.TrimSpace(path)
			if trimmedPath != "" {
				targetProjects = append(targetProjects, trimmedPath)
			}
		}
		if len(targetProjects) == 0 {
			return fmt.Errorf("no valid project paths provided")
		}

		fmt.Printf("Registering Linear MCP server to %d specified projects\n", len(targetProjects))
		for _, projPath := range targetProjects {
			if err := registerLinearToProject(settings, projPath, linearServerConfig); err != nil {
				return fmt.Errorf("failed to register Linear MCP server to project '%s': %w", projPath, err)
			}
			fmt.Printf("  - Registered to project: %s\n", projPath)
		}
	}

	updatedData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal claude code settings: %w", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write claude code settings: %w", err)
	}

	fmt.Printf("Claude Code MCP settings updated at %s\n", configPath)
	return nil
}

// registerLinearToUserScope registers the Linear MCP server to user-scoped mcpServers (applies to all projects)
func registerLinearToUserScope(settings map[string]interface{}, linearServerConfig map[string]interface{}) error {
	// Get or create user-scoped mcpServers
	var mcpServers map[string]interface{}
	if existingMcpServers, exists := settings["mcpServers"]; exists {
		if mcpServersMap, ok := existingMcpServers.(map[string]interface{}); ok {
			mcpServers = mcpServersMap
		} else {
			// If existing mcpServers is not a map, create a new one
			mcpServers = map[string]interface{}{}
		}
	} else {
		mcpServers = map[string]interface{}{}
	}

	// Add/update the linear server configuration
	mcpServers["linear"] = linearServerConfig
	settings["mcpServers"] = mcpServers

	return nil
}

// registerLinearToProject registers the Linear MCP server to a specific project
func registerLinearToProject(settings map[string]interface{}, projectPath string, linearServerConfig map[string]interface{}) error {
	// Get projects map
	projects, ok := settings["projects"].(map[string]interface{})
	if !ok {
		projects = map[string]interface{}{}
		settings["projects"] = projects
	}

	// Get or create the specific project
	var project map[string]interface{}
	if existingProject, exists := projects[projectPath]; exists {
		if projectMap, ok := existingProject.(map[string]interface{}); ok {
			project = projectMap
		} else {
			// If existing project is not a map, create a new one
			project = map[string]interface{}{}
		}
	} else {
		project = map[string]interface{}{}
	}

	// Get or create mcpServers for this project
	var mcpServers map[string]interface{}
	if existingMcpServers, exists := project["mcpServers"]; exists {
		if mcpServersMap, ok := existingMcpServers.(map[string]interface{}); ok {
			mcpServers = mcpServersMap
		} else {
			// If existing mcpServers is not a map, create a new one
			mcpServers = map[string]interface{}{}
		}
	} else {
		mcpServers = map[string]interface{}{}
	}

	// Add/update the linear server configuration
	mcpServers["linear"] = linearServerConfig
	project["mcpServers"] = mcpServers
	projects[projectPath] = project

	return nil
}
