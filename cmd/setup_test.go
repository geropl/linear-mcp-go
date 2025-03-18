package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Define expectation types
type fileExpectation struct {
	path      string
	content   string
	mustExist bool
}

type expectations struct {
	files     map[string]fileExpectation
	errors    []string
	exitCode  int
}

// TestSetupCommand tests the setup command with various combinations of parameters
func TestSetupCommand(t *testing.T) {
	// Build the binary
	binaryPath, err := buildBinary()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(binaryPath))

	// Define test cases
	testCases := []struct {
		name        string
		toolParam   string
		writeAccess bool
		autoApprove string
		expect      expectations
	}{
		{
			name:        "Cline Only",
			toolParam:   "cline",
			writeAccess: true,
			autoApprove: "allow-read-only",
			expect: expectations{
				files: map[string]fileExpectation{
					"cline": {
						path:      "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"linear": {
									"args": ["serve", "--write-access=true"],
									"autoApprove": ["linear_search_issues", "linear_get_user_issues", "linear_get_issue", "linear_get_teams"],
									"disabled": false
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Roo Code Only",
			toolParam:   "roo-code",
			writeAccess: true,
			autoApprove: "allow-read-only",
			expect: expectations{
				files: map[string]fileExpectation{
					"roo-code": {
						path:      "home/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"linear": {
									"args": ["serve", "--write-access=true"],
									"autoApprove": ["linear_search_issues", "linear_get_user_issues", "linear_get_issue", "linear_get_teams"],
									"disabled": false
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Multiple Tools",
			toolParam:   "cline,roo-code",
			writeAccess: true,
			autoApprove: "allow-read-only",
			expect: expectations{
				files: map[string]fileExpectation{
					"cline": {
						path:      "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"linear": {
									"args": ["serve", "--write-access=true"],
									"autoApprove": ["linear_search_issues", "linear_get_user_issues", "linear_get_issue", "linear_get_teams"],
									"disabled": false
								}
							}
						}`,
					},
					"roo-code": {
						path:      "home/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"linear": {
									"args": ["serve", "--write-access=true"],
									"autoApprove": ["linear_search_issues", "linear_get_user_issues", "linear_get_issue", "linear_get_teams"],
									"disabled": false
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Invalid Tool",
			toolParam:   "invalid-tool,cline",
			writeAccess: true,
			autoApprove: "allow-read-only",
			expect: expectations{
				files: map[string]fileExpectation{
					"cline": {
						path:      "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"linear": {
									"args": ["serve", "--write-access=true"],
									"autoApprove": ["linear_search_issues", "linear_get_user_issues", "linear_get_issue", "linear_get_teams"],
									"disabled": false
								}
							}
						}`,
					},
				},
				errors:   []string{"Unsupported tool: invalid-tool"},
				exitCode: 1,
			},
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory
			rootDir, err := os.MkdirTemp("", "linear-mcp-go-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(rootDir)

			// Set up the directory structure
			homeDir := filepath.Join(rootDir, "home")

			// Copy the binary to the temp directory
			tempBinaryPath := filepath.Join(rootDir, "linear-mcp-go")
			if err := copyFile(binaryPath, tempBinaryPath); err != nil {
				t.Fatalf("Failed to copy binary: %v", err)
			}
			if err := os.Chmod(tempBinaryPath, 0755); err != nil {
				t.Fatalf("Failed to make binary executable: %v", err)
			}

			// Set the HOME environment variable
			oldHome := os.Getenv("HOME")
			os.Setenv("HOME", homeDir)
			defer os.Setenv("HOME", oldHome)

			// Set the LINEAR_API_KEY environment variable
			oldApiKey := os.Getenv("LINEAR_API_KEY")
			os.Setenv("LINEAR_API_KEY", "test-api-key")
			defer os.Setenv("LINEAR_API_KEY", oldApiKey)

			// Build the command
			args := []string{"setup", "--tool=" + tc.toolParam}
			if tc.writeAccess {
				args = append(args, "--write-access=true")
			}
			if tc.autoApprove != "" {
				args = append(args, "--auto-approve="+tc.autoApprove)
			}

			// Execute the command
			cmd := exec.Command(tempBinaryPath, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()

			// Check exit code
			exitCode := 0
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					t.Fatalf("Failed to run command: %v", err)
				}
			}

			// Verify exit code
			if exitCode != tc.expect.exitCode {
				t.Errorf("Expected exit code %d, got %d", tc.expect.exitCode, exitCode)
			}

			// Verify expected files
			verifyFileExpectations(t, rootDir, tc.expect.files)

			// Verify expected errors in output
			output := stdout.String() + stderr.String()
			for _, expectedError := range tc.expect.errors {
				if !strings.Contains(output, expectedError) {
					t.Errorf("Expected output to contain '%s', got: %s", expectedError, output)
				}
			}
		})
	}
}

// Helper function to verify file expectations
func verifyFileExpectations(t *testing.T, rootDir string, fileExpects map[string]fileExpectation) {
	for tool, expect := range fileExpects {
		filePath := filepath.Join(rootDir, expect.path)

		// Check if file exists
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			if expect.mustExist {
				t.Errorf("Expected file %s was not created for %s", filePath, tool)
			}
			continue
		}

		// File exists, verify content if expected
		if expect.content != "" {
			actualContent, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read configuration file %s: %v", filePath, err)
			}

			// Parse both expected and actual content as JSON for comparison
			var expectedJSON, actualJSON map[string]interface{}
			
			if err := json.Unmarshal([]byte(expect.content), &expectedJSON); err != nil {
				t.Fatalf("Failed to parse expected JSON for %s: %v", tool, err)
			}
			
			if err := json.Unmarshal(actualContent, &actualJSON); err != nil {
				t.Fatalf("Failed to parse actual JSON in file %s: %v", filePath, err)
			}
			
			// Process the JSON objects to make them comparable
			normalizeJSON(expectedJSON)
			normalizeJSON(actualJSON)
			
			// Compare the JSON objects
			if diff := cmp.Diff(expectedJSON, actualJSON); diff != "" {
				t.Errorf("File content mismatch for %s (-want +got):\n%s", tool, diff)
			}
		}
	}
}

// normalizeJSON processes a JSON object to make it comparable
// by removing fields that may vary and sorting arrays
func normalizeJSON(jsonObj map[string]interface{}) {
	if mcpServers, ok := jsonObj["mcpServers"].(map[string]interface{}); ok {
		if linear, ok := mcpServers["linear"].(map[string]interface{}); ok {
			// Remove the command field since it contains the full path
			delete(linear, "command")
			
			// Remove the env field since it contains the API key
			delete(linear, "env")
			
			// Sort the autoApprove array
			if autoApprove, ok := linear["autoApprove"].([]interface{}); ok {
				// Convert to strings and sort
				strSlice := make([]string, len(autoApprove))
				for i, v := range autoApprove {
					strSlice[i] = v.(string)
				}
				
				// Sort the strings
				sort.Strings(strSlice)
				
				// Convert back to []interface{}
				sortedSlice := make([]interface{}, len(strSlice))
				for i, v := range strSlice {
					sortedSlice[i] = v
				}
				
				// Replace the original array with the sorted one
				linear["autoApprove"] = sortedSlice
			}
		}
	}
}

// Helper function to build the binary
func buildBinary() (string, error) {
	// Create a temporary directory for the binary
	tempDir, err := os.MkdirTemp("", "linear-mcp-go-build-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Get the project root directory (parent of cmd directory)
	currentDir, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	
	// Ensure we're building from the project root
	projectRoot := filepath.Dir(currentDir)
	if filepath.Base(currentDir) != "cmd" {
		// If we're already in the project root, use the current directory
		projectRoot = currentDir
	}
	
	fmt.Printf("Building binary from project root: %s\n", projectRoot)
	
	// Build the binary
	binaryPath := filepath.Join(tempDir, "linear-mcp-go")
	cmd := exec.Command("go", "build", "-o", binaryPath)
	cmd.Dir = projectRoot // Set the working directory to the project root
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to build binary: %w\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	// Verify the binary exists and is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to stat binary: %w", err)
	}
	
	if info.Size() == 0 {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("binary file is empty")
	}

	// Make sure the binary is executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	fmt.Printf("Successfully built binary at %s (size: %d bytes)\n", binaryPath, info.Size())
	return binaryPath, nil
}

// Helper function to copy a file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}