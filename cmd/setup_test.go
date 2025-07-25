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

type preExistingFile struct {
	path    string
	content string
}

type expectations struct {
	files    map[string]fileExpectation
	errors   []string
	exitCode int
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
		name             string
		toolParam        string
		writeAccess      bool
		autoApprove      string
		projectPath      string
		preExistingFiles map[string]preExistingFile
		expect           expectations
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
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
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
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
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
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
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
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
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
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
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
		{
			name:        "Preserve Existing Arrays in Config",
			toolParam:   "cline",
			writeAccess: true,
			autoApprove: "allow-read-only",
			preExistingFiles: map[string]preExistingFile{
				"cline": {
					path: "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
					content: `{
						"mcpServers": {
							"existing-server": {
								"command": "/path/to/existing/server",
								"args": ["serve", "--option1", "--option2"],
								"autoApprove": ["tool1", "tool2", "tool3"],
								"env": {
									"API_KEY": "existing-key"
								},
								"disabled": false,
								"customArray": ["item1", "item2"],
								"nestedObject": {
									"arrayField": ["nested1", "nested2"]
								}
							}
						},
						"otherTopLevelArray": ["value1", "value2"],
						"otherConfig": {
							"someArray": [1, 2, 3]
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"cline": {
						path:      "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"existing-server": {
									"command": "/path/to/existing/server",
									"args": ["serve", "--option1", "--option2"],
									"autoApprove": ["tool1", "tool2", "tool3"],
									"env": {
										"API_KEY": "existing-key"
									},
									"disabled": false,
									"customArray": ["item1", "item2"],
									"nestedObject": {
										"arrayField": ["nested1", "nested2"]
									}
								},
								"linear": {
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
									"disabled": false
								}
							},
							"otherTopLevelArray": ["value1", "value2"],
							"otherConfig": {
								"someArray": [1, 2, 3]
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Complex Array Preservation Test",
			toolParam:   "cline",
			writeAccess: false,
			autoApprove: "linear_get_issue,linear_search_issues",
			preExistingFiles: map[string]preExistingFile{
				"cline": {
					path: "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
					content: `{
						"mcpServers": {
							"github": {
								"command": "npx",
								"args": ["-y", "@modelcontextprotocol/server-github"],
								"env": {
									"GITHUB_PERSONAL_ACCESS_TOKEN": "github_token"
								},
								"autoApprove": ["search_repositories", "get_file_contents"],
								"disabled": false
							},
							"filesystem": {
								"command": "npx",
								"args": ["-y", "@modelcontextprotocol/server-filesystem", "/workspace"],
								"autoApprove": [],
								"disabled": false
							}
						},
						"globalSettings": {
							"enabledFeatures": ["autocomplete", "syntax-highlighting"],
							"debugModes": ["verbose", "trace"]
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"cline": {
						path:      "home/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json",
						mustExist: true,
						content: `{
							"mcpServers": {
								"github": {
									"command": "npx",
									"args": ["-y", "@modelcontextprotocol/server-github"],
									"env": {
										"GITHUB_PERSONAL_ACCESS_TOKEN": "github_token"
									},
									"autoApprove": ["search_repositories", "get_file_contents"],
									"disabled": false
								},
								"filesystem": {
									"command": "npx",
									"args": ["-y", "@modelcontextprotocol/server-filesystem", "/workspace"],
									"autoApprove": [],
									"disabled": false
								},
								"linear": {
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": ["linear_get_issue", "linear_search_issues"],
									"disabled": false
								}
							},
							"globalSettings": {
								"enabledFeatures": ["autocomplete", "syntax-highlighting"],
								"debugModes": ["verbose", "trace"]
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code Only",
			toolParam:   "claude-code",
			projectPath: "/workspace/test-project",
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/test-project": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": [],
											"disabled": false
										}
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code with Existing File",
			toolParam:   "claude-code",
			projectPath: "/workspace/test-project",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {
							"/workspace/another-project": {
								"mcpServers": {
									"another-server": {
										"command": "/path/to/another/server"
									}
								}
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/another-project": {
									"mcpServers": {
										"another-server": {
											"command": "/path/to/another/server"
										}
									}
								},
								"/workspace/test-project": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": [],
											"disabled": false
										}
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:      "Claude Code No Existing Projects and No Project Path - User Scope Registration",
			toolParam: "claude-code",
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {},
							"mcpServers": {
								"linear": {
									"type": "stdio",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": [],
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
			name:      "Claude Code Register to User Scope with Existing Projects",
			toolParam: "claude-code",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {
							"/workspace/project1": {
								"mcpServers": {
									"existing-server": {
										"command": "/path/to/existing/server"
									}
								}
							},
							"/workspace/project2": {
								"someOtherConfig": "value"
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/project1": {
									"mcpServers": {
										"existing-server": {
											"command": "/path/to/existing/server"
										}
									}
								},
								"/workspace/project2": {
									"someOtherConfig": "value"
								}
							},
							"mcpServers": {
								"linear": {
									"type": "stdio",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									},
									"autoApprove": [],
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
			name:        "Claude Code Multiple Project Paths",
			toolParam:   "claude-code",
			projectPath: "/workspace/proj1,/workspace/proj2, /workspace/proj3 ",
			writeAccess: true,
			autoApprove: "linear_get_issue,linear_search_issues",
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/proj1": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve", "--write-access=true"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": ["linear_get_issue", "linear_search_issues"],
											"disabled": false
										}
									}
								},
								"/workspace/proj2": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve", "--write-access=true"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": ["linear_get_issue", "linear_search_issues"],
											"disabled": false
										}
									}
								},
								"/workspace/proj3": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve", "--write-access=true"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": ["linear_get_issue", "linear_search_issues"],
											"disabled": false
										}
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code Mixed Existing and New Projects",
			toolParam:   "claude-code",
			projectPath: "/workspace/existing-project,/workspace/new-project",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {
							"/workspace/existing-project": {
								"mcpServers": {
									"other-server": {
										"command": "/path/to/other/server"
									}
								},
								"customConfig": "value"
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/existing-project": {
									"mcpServers": {
										"other-server": {
											"command": "/path/to/other/server"
										},
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": [],
											"disabled": false
										}
									},
									"customConfig": "value"
								},
								"/workspace/new-project": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve"],
											"env": {
												"LINEAR_API_KEY": "test-api-key"
											},
											"autoApprove": [],
											"disabled": false
										}
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code Empty Project Path List",
			toolParam:   "claude-code",
			projectPath: " , , ",
			expect: expectations{
				errors:   []string{"no valid project paths provided"},
				exitCode: 1,
			},
		},
		{
			name:        "Claude Code Complex Settings Preservation",
			toolParam:   "claude-code",
			projectPath: "/workspace/new-project",
			writeAccess: true,
			autoApprove: "linear_get_issue,linear_search_issues",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"firstStartTime": "2025-06-11T14:49:28.932Z",
						"userID": "31553dcf54399f00daf126faf48dbb0e626926f50e9bf49c16cb05c06f65cfd8",
						"globalSettings": {
							"theme": "dark",
							"autoSave": true,
							"debugMode": false,
							"experimentalFeatures": ["feature1", "feature2", "feature3"],
							"limits": {
								"maxTokens": 4096,
								"timeout": 30000,
								"retries": 3
							},
							"customMappings": {
								"shortcuts": {
									"ctrl+s": "save",
									"ctrl+z": "undo"
								},
								"aliases": ["alias1", "alias2"]
							}
						},
						"recentProjects": ["/workspace/project1", "/workspace/project2", "/workspace/project3"],
						"projects": {
							"/workspace/existing-project": {
								"allowedTools": ["tool1", "tool2", "tool3"],
								"history": [
									{
										"timestamp": "2025-06-11T15:00:00.000Z",
										"action": "create_file",
										"details": {"filename": "test.js", "size": 1024}
									}
								],
								"dontCrawlDirectory": false,
								"mcpContextUris": ["file:///workspace/docs", "https://api.example.com/docs"],
								"mcpServers": {
									"github": {
										"type": "stdio",
										"command": "npx",
										"args": ["-y", "@modelcontextprotocol/server-github"],
										"env": {"GITHUB_PERSONAL_ACCESS_TOKEN": "github_token_123"},
										"autoApprove": ["search_repositories", "get_file_contents"],
										"disabled": false,
										"customConfig": {"rateLimit": 5000, "features": ["search", "read"]}
									},
									"filesystem": {
										"type": "stdio",
										"command": "npx", 
										"args": ["-y", "@modelcontextprotocol/server-filesystem", "/workspace"],
										"autoApprove": ["read_file", "list_directory"],
										"disabled": false,
										"permissions": {"read": true, "write": false, "execute": false}
									}
								},
								"enabledMcpjsonServers": ["server1", "server2"],
								"disabledMcpjsonServers": ["server3", "server4"],
								"hasTrustDialogAccepted": true,
								"projectOnboardingSeenCount": 3,
								"customProjectSettings": {
									"linting": {"enabled": true, "rules": ["rule1", "rule2"]},
									"formatting": {"tabSize": 2, "insertSpaces": true}
								}
							}
						},
						"analytics": {
							"enabled": true,
							"sessionId": "session_12345",
							"metrics": {"commandsExecuted": 42, "filesModified": 15}
						},
						"version": "1.2.3"
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"firstStartTime": "2025-06-11T14:49:28.932Z",
							"userID": "31553dcf54399f00daf126faf48dbb0e626926f50e9bf49c16cb05c06f65cfd8",
							"globalSettings": {
								"theme": "dark",
								"autoSave": true,
								"debugMode": false,
								"experimentalFeatures": ["feature1", "feature2", "feature3"],
								"limits": {
									"maxTokens": 4096,
									"timeout": 30000,
									"retries": 3
								},
								"customMappings": {
									"shortcuts": {
										"ctrl+s": "save",
										"ctrl+z": "undo"
									},
									"aliases": ["alias1", "alias2"]
								}
							},
							"recentProjects": ["/workspace/project1", "/workspace/project2", "/workspace/project3"],
							"projects": {
								"/workspace/existing-project": {
									"allowedTools": ["tool1", "tool2", "tool3"],
									"history": [
										{
											"timestamp": "2025-06-11T15:00:00.000Z",
											"action": "create_file",
											"details": {"filename": "test.js", "size": 1024}
										}
									],
									"dontCrawlDirectory": false,
									"mcpContextUris": ["file:///workspace/docs", "https://api.example.com/docs"],
									"mcpServers": {
										"github": {
											"type": "stdio",
											"command": "npx",
											"args": ["-y", "@modelcontextprotocol/server-github"],
											"env": {"GITHUB_PERSONAL_ACCESS_TOKEN": "github_token_123"},
											"autoApprove": ["search_repositories", "get_file_contents"],
											"disabled": false,
											"customConfig": {"rateLimit": 5000, "features": ["search", "read"]}
										},
										"filesystem": {
											"type": "stdio",
											"command": "npx", 
											"args": ["-y", "@modelcontextprotocol/server-filesystem", "/workspace"],
											"autoApprove": ["read_file", "list_directory"],
											"disabled": false,
											"permissions": {"read": true, "write": false, "execute": false}
										}
									},
									"enabledMcpjsonServers": ["server1", "server2"],
									"disabledMcpjsonServers": ["server3", "server4"],
									"hasTrustDialogAccepted": true,
									"projectOnboardingSeenCount": 3,
									"customProjectSettings": {
										"linting": {"enabled": true, "rules": ["rule1", "rule2"]},
										"formatting": {"tabSize": 2, "insertSpaces": true}
									}
								},
								"/workspace/new-project": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve", "--write-access=true"],
											"env": {"LINEAR_API_KEY": "test-api-key"},
											"autoApprove": ["linear_get_issue", "linear_search_issues"],
											"disabled": false
										}
									}
								}
							},
							"analytics": {
								"enabled": true,
								"sessionId": "session_12345",
								"metrics": {"commandsExecuted": 42, "filesModified": 15}
							},
							"version": "1.2.3"
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code Update Existing Linear Server",
			toolParam:   "claude-code",
			projectPath: "/workspace/existing-project",
			writeAccess: false,
			autoApprove: "allow-read-only",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {
							"/workspace/existing-project": {
								"mcpServers": {
									"linear": {
										"type": "stdio",
										"command": "/old/path/to/linear",
										"args": ["serve", "--old-flag"],
										"env": {"LINEAR_API_KEY": "old-key"},
										"autoApprove": ["old_tool"],
										"disabled": true
									},
									"other-server": {
										"command": "/path/to/other/server"
									}
								}
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/existing-project": {
									"mcpServers": {
										"linear": {
											"type": "stdio",
											"command": "home/mcp-servers/linear-mcp-go",
											"args": ["serve"],
											"env": {"LINEAR_API_KEY": "test-api-key"},
											"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
											"disabled": false
										},
										"other-server": {
											"command": "/path/to/other/server"
										}
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Claude Code User Scope with Existing User-Scoped Servers",
			toolParam:   "claude-code",
			writeAccess: true,
			autoApprove: "linear_get_issue,linear_search_issues",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {
							"/workspace/project1": {
								"mcpServers": {
									"project-specific-server": {
										"command": "/path/to/project/server"
									}
								}
							}
						},
						"mcpServers": {
							"existing-user-server": {
								"type": "stdio",
								"command": "/path/to/existing/user/server",
								"args": ["serve"],
								"env": {"API_KEY": "existing-key"},
								"autoApprove": ["existing_tool"],
								"disabled": false
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {
								"/workspace/project1": {
									"mcpServers": {
										"project-specific-server": {
											"command": "/path/to/project/server"
										}
									}
								}
							},
							"mcpServers": {
								"existing-user-server": {
									"type": "stdio",
									"command": "/path/to/existing/user/server",
									"args": ["serve"],
									"env": {"API_KEY": "existing-key"},
									"autoApprove": ["existing_tool"],
									"disabled": false
								},
								"linear": {
									"type": "stdio",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {"LINEAR_API_KEY": "test-api-key"},
									"autoApprove": ["linear_get_issue", "linear_search_issues"],
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
			name:        "Claude Code User Scope Update Existing User-Scoped Linear Server",
			toolParam:   "claude-code",
			writeAccess: false,
			autoApprove: "allow-read-only",
			preExistingFiles: map[string]preExistingFile{
				"claude-code": {
					path: "home/.claude.json",
					content: `{
						"projects": {},
						"mcpServers": {
							"linear": {
								"type": "stdio",
								"command": "/old/path/to/linear",
								"args": ["serve", "--old-flag"],
								"env": {"LINEAR_API_KEY": "old-key"},
								"autoApprove": ["old_tool"],
								"disabled": true
							},
							"other-user-server": {
								"command": "/path/to/other/user/server"
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"claude-code": {
						path:      "home/.claude.json",
						mustExist: true,
						content: `{
							"projects": {},
							"mcpServers": {
								"linear": {
									"type": "stdio",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve"],
									"env": {"LINEAR_API_KEY": "test-api-key"},
											"autoApprove": ["linear_get_initiative", "linear_get_issue", "linear_get_issue_comments", "linear_get_milestone", "linear_get_project", "linear_get_teams", "linear_get_user_issues", "linear_search_issues", "linear_search_projects"],
											"disabled": false
								},
								"other-user-server": {
									"command": "/path/to/other/user/server"
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Ona Only",
			toolParam:   "ona",
			writeAccess: true,
			expect: expectations{
				files: map[string]fileExpectation{
					"ona": {
						path:      ".gitpod/mcp-config.json",
						mustExist: true,
						content: `{
							"servers": {
								"linear": {
									"name": "linear",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Ona with Project Path",
			toolParam:   "ona",
			projectPath: "/workspace/test-project",
			writeAccess: false,
			expect: expectations{
				files: map[string]fileExpectation{
					"ona": {
						path:      "/workspace/test-project/.gitpod/mcp-config.json",
						mustExist: true,
						content: `{
							"servers": {
								"linear": {
									"name": "linear",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=false"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Ona with Existing Config",
			toolParam:   "ona",
			writeAccess: true,
			preExistingFiles: map[string]preExistingFile{
				"ona": {
					path: ".gitpod/mcp-config.json",
					content: `{
						"servers": {
							"playwright": {
								"name": "playwright",
								"command": "npx",
								"args": ["-y", "@executeautomation/playwright-mcp-server"]
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"ona": {
						path:      ".gitpod/mcp-config.json",
						mustExist: true,
						content: `{
							"servers": {
								"playwright": {
									"name": "playwright",
									"command": "npx",
									"args": ["-y", "@executeautomation/playwright-mcp-server"]
								},
								"linear": {
									"name": "linear",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=true"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
			},
		},
		{
			name:        "Ona with Complex Nested Config",
			toolParam:   "ona",
			writeAccess: false,
			preExistingFiles: map[string]preExistingFile{
				"ona": {
					path: ".gitpod/mcp-config.json",
					content: `{
						"version": "1.0.0",
						"metadata": {
							"created": "2024-01-01T00:00:00Z",
							"author": "test-user",
							"tags": ["development", "testing", "automation"],
							"config": {
								"nested": {
									"deeply": {
										"properties": ["value1", "value2", "value3"],
										"settings": {
											"enabled": true,
											"timeout": 30000,
											"retries": 3,
											"features": {
												"advanced": {
													"caching": true,
													"compression": false,
													"encryption": {
														"algorithm": "AES-256",
														"keySize": 256,
														"modes": ["CBC", "GCM", "CTR"]
													}
												}
											}
										}
									}
								}
							}
						},
						"servers": {
							"custom-server": {
								"name": "custom-server",
								"command": "/usr/local/bin/custom-mcp-server",
								"args": ["--mode", "production", "--verbose"],
								"env": {
									"CUSTOM_API_KEY": "secret-key-123",
									"CUSTOM_ENDPOINT": "https://api.example.com/v1"
								},
								"timeout": 60000,
								"retries": 5,
								"features": {
									"streaming": true,
									"batching": false,
									"compression": {
										"enabled": true,
										"algorithm": "gzip",
										"level": 6
									}
								},
								"customArrays": {
									"supportedFormats": ["json", "xml", "yaml"],
									"allowedOrigins": ["localhost", "*.example.com", "api.test.com"],
									"permissions": ["read", "write", "execute"]
								},
								"nestedConfig": {
									"database": {
										"connection": {
											"host": "localhost",
											"port": 5432,
											"ssl": {
												"enabled": true,
												"cert": "/path/to/cert.pem",
												"key": "/path/to/key.pem",
												"ca": "/path/to/ca.pem"
											}
										},
										"pool": {
											"min": 5,
											"max": 20,
											"idle": 300
										}
									}
								}
							}
						},
						"globalSettings": {
							"logLevel": "info",
							"enableMetrics": true,
							"metricsConfig": {
								"endpoint": "http://metrics.example.com:9090",
								"interval": 30,
								"labels": {
									"environment": "test",
									"service": "mcp-server",
									"version": "1.0.0"
								}
							}
						}
					}`,
				},
			},
			expect: expectations{
				files: map[string]fileExpectation{
					"ona": {
						path:      ".gitpod/mcp-config.json",
						mustExist: true,
						content: `{
							"version": "1.0.0",
							"metadata": {
								"created": "2024-01-01T00:00:00Z",
								"author": "test-user",
								"tags": ["development", "testing", "automation"],
								"config": {
									"nested": {
										"deeply": {
											"properties": ["value1", "value2", "value3"],
											"settings": {
												"enabled": true,
												"timeout": 30000,
												"retries": 3,
												"features": {
													"advanced": {
														"caching": true,
														"compression": false,
														"encryption": {
															"algorithm": "AES-256",
															"keySize": 256,
															"modes": ["CBC", "GCM", "CTR"]
														}
													}
												}
											}
										}
									}
								}
							},
							"servers": {
								"custom-server": {
									"name": "custom-server",
									"command": "/usr/local/bin/custom-mcp-server",
									"args": ["--mode", "production", "--verbose"],
									"env": {
										"CUSTOM_API_KEY": "secret-key-123",
										"CUSTOM_ENDPOINT": "https://api.example.com/v1"
									},
									"timeout": 60000,
									"retries": 5,
									"features": {
										"streaming": true,
										"batching": false,
										"compression": {
											"enabled": true,
											"algorithm": "gzip",
											"level": 6
										}
									},
									"customArrays": {
										"supportedFormats": ["json", "xml", "yaml"],
										"allowedOrigins": ["localhost", "*.example.com", "api.test.com"],
										"permissions": ["read", "write", "execute"]
									},
									"nestedConfig": {
										"database": {
											"connection": {
												"host": "localhost",
												"port": 5432,
												"ssl": {
													"enabled": true,
													"cert": "/path/to/cert.pem",
													"key": "/path/to/key.pem",
													"ca": "/path/to/ca.pem"
												}
											},
											"pool": {
												"min": 5,
												"max": 20,
												"idle": 300
											}
										}
									}
								},
								"linear": {
									"name": "linear",
									"command": "home/mcp-servers/linear-mcp-go",
									"args": ["serve", "--write-access=false"],
									"env": {
										"LINEAR_API_KEY": "test-api-key"
									}
								}
							},
							"globalSettings": {
								"logLevel": "info",
								"enableMetrics": true,
								"metricsConfig": {
									"endpoint": "http://metrics.example.com:9090",
									"interval": 30,
									"labels": {
										"environment": "test",
										"service": "mcp-server",
										"version": "1.0.0"
									}
								}
							}
						}`,
					},
				},
				exitCode: 0,
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

			// Create pre-existing files if specified
			for _, preFile := range tc.preExistingFiles {
				fullPath := filepath.Join(rootDir, preFile.path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("Failed to create directory for pre-existing file %s: %v", fullPath, err)
				}
				if err := os.WriteFile(fullPath, []byte(preFile.content), 0644); err != nil {
					t.Fatalf("Failed to create pre-existing file %s: %v", fullPath, err)
				}
			}

			// Build the command
			args := []string{"setup", "--tool=" + tc.toolParam}
			if tc.writeAccess {
				args = append(args, "--write-access=true")
			}
			if tc.autoApprove != "" {
				args = append(args, "--auto-approve="+tc.autoApprove)
			}
			if tc.projectPath != "" {
				args = append(args, "--project-path="+tc.projectPath)
			}

			// Execute the command
			cmd := exec.Command(tempBinaryPath, args...)
			cmd.Dir = rootDir // Set working directory to the test root
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
	normalizeCommandPaths(jsonObj)
	normalizeJSONRecursive(jsonObj)
}

// normalizeCommandPaths normalizes command paths in server configurations
func normalizeCommandPaths(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Look for server configuration containers
		if isServerContainer(v) {
			for _, serverConfig := range v {
				if serverMap, ok := serverConfig.(map[string]interface{}); ok {
					// Normalize the command field by stripping temporary directory prefix
					if command, ok := serverMap["command"].(string); ok {
						// Strip the temporary test directory prefix, keeping only the meaningful part
						// Pattern: /tmp/linear-mcp-go-test-*/home/... -> home/...
						if strings.Contains(command, "/home/") {
							parts := strings.Split(command, "/home/")
							if len(parts) > 1 {
								serverMap["command"] = "home/" + parts[1]
							}
						}
					}
				}
			}
		}

		// Process nested objects recursively
		for _, value := range v {
			normalizeCommandPaths(value)
		}

	case []interface{}:
		// Process array elements recursively
		for _, item := range v {
			normalizeCommandPaths(item)
		}
	}
}

// isServerContainer checks if a map contains server configurations
func isServerContainer(m map[string]interface{}) bool {
	// Check if this looks like a server container by examining its values
	for _, value := range m {
		if serverMap, ok := value.(map[string]interface{}); ok {
			// If it has command field, it's likely a server config container
			if _, hasCommand := serverMap["command"]; hasCommand {
				return true
			}
		}
	}
	return false
}

// normalizeJSONRecursive recursively processes JSON objects to normalize them for comparison
func normalizeJSONRecursive(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Process all map entries recursively
		for _, value := range v {
			normalizeJSONRecursive(value)
		}

	case []interface{}:
		// Sort arrays if they contain strings
		if len(v) > 0 {
			// Check if all elements are strings
			allStrings := true
			for _, item := range v {
				if _, ok := item.(string); !ok {
					allStrings = false
					break
				}
			}

			if allStrings {
				// Convert to string slice, sort, and convert back
				strSlice := make([]string, len(v))
				for i, item := range v {
					strSlice[i] = item.(string)
				}
				sort.Strings(strSlice)

				// Update the original slice in place
				for i, str := range strSlice {
					v[i] = str
				}
			}
		}

		// Process array elements recursively
		for _, item := range v {
			normalizeJSONRecursive(item)
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
