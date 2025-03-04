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

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up the Linear MCP server for use with an AI assistant",
	Long: `Set up the Linear MCP server for use with an AI assistant.
This command installs the Linear MCP server and configures it for use with the specified AI assistant tool.
Currently supported tools: cline`,
	Run: func(cmd *cobra.Command, args []string) {
		tool, _ := cmd.Flags().GetString("tool")
		writeAccess, _ := cmd.Flags().GetBool("write-access")

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
			fmt.Println("Linear MCP binary not found on path, copying current binary...")
			err := copySelfToBinaryPath(binaryPath)
			if err != nil {
				fmt.Printf("Error copying Linear MCP binary: %v\n", err)
				os.Exit(1)
			}
		}

		// Set up the tool-specific configuration
		switch strings.ToLower(tool) {
		case "cline":
			if err := setupCline(binaryPath, apiKey, writeAccess); err != nil {
				fmt.Printf("Error setting up Cline: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unsupported tool: %s\n", tool)
			fmt.Println("Currently supported tools: cline")
			os.Exit(1)
		}

		fmt.Printf("Linear MCP server successfully set up for %s\n", tool)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Add flags to the setup command
	setupCmd.Flags().String("tool", "cline", "The AI assistant tool to set up for (default: cline)")
	setupCmd.Flags().Bool("write-access", false, "Enable write operations (default: false)")
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

	return "", false
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

// setupCline sets up the Linear MCP server for Cline
func setupCline(binaryPath, apiKey string, writeAccess bool) error {
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

	// Create the config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	serverArgs := []string{"serve"}
	if writeAccess {
		serverArgs = append(serverArgs, "--write-access=true")
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
				"autoApprove": []string{},
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

	fmt.Printf("Cline MCP settings updated at %s\n", settingsPath)
	return nil
}
