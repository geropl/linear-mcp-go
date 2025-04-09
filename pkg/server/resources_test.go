package server

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func TestResourceHandlers(t *testing.T) {
	// Define test cases
	tests := []struct {
		handlerName string
		name        string
		uri         string
		handlerFunc func(*linear.LinearClient) mcpserver.ResourceHandlerFunc
	}{
		// TeamsResourceHandler test cases
		{
			handlerName: "TeamsResourceHandler",
			name:        "List All",
			uri:         "linear://teams",
			handlerFunc: TeamsResourceHandler,
		},
		// TeamResourceHandler test cases
		{
			handlerName: "TeamResourceHandler",
			name:        "Fetch By ID",
			uri:         "linear://team/" + TEAM_ID,
			handlerFunc: TeamResourceHandler,
		},
		{
			handlerName: "TeamResourceHandler",
			name:        "Fetch By Name",
			uri:         "linear://team/" + TEAM_NAME,
			handlerFunc: TeamResourceHandler,
		},
		{
			handlerName: "TeamResourceHandler",
			name:        "Fetch By Key",
			uri:         "linear://team/" + TEAM_KEY,
			handlerFunc: TeamResourceHandler,
		},
		{
			handlerName: "TeamResourceHandler",
			name:        "Invalid ID",
			uri:         "linear://team/invalid-identifier-does-not-exist", // Use a clearly invalid identifier
			handlerFunc: TeamResourceHandler,
		},
		{
			handlerName: "TeamResourceHandler",
			name:        "Missing ID",
			uri:         "linear://team/", // Test case where ID is missing from URI path
			handlerFunc: TeamResourceHandler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.handlerName+"_"+tt.name, func(t *testing.T) {
			// Generate fixture and golden file paths
			fixtureName := "resource_" + tt.handlerName + "_" + tt.name
			goldenPath := filepath.Join("../../testdata/golden", fixtureName+".golden")

			// Create test client with VCR
			// Use distinct flags for resource tests to avoid conflicts
			client, cleanup := linear.NewTestClient(t, fixtureName, *record || *recordWrites)
			defer cleanup()

			// Get the handler function
			handler := tt.handlerFunc(client)

			// Create the request
			request := mcp.ReadResourceRequest{}
			request.Params.URI = tt.uri

			// Call the handler
			contents, err := handler(context.Background(), request)

			// Extract the actual output and error
			var actualOutput, actualErr string
			if err != nil {
				actualErr = err.Error()
			} else {
				// Marshal the contents to JSON for comparison
				jsonBytes, jsonErr := json.MarshalIndent(contents, "", "  ") // Use indent for readability
				if jsonErr != nil {
					t.Fatalf("Failed to marshal resource contents to JSON: %v", jsonErr)
				}
				actualOutput = string(jsonBytes)
			}

			// If goldenResource flag is set, update the golden file
			if *golden {
				writeGoldenFile(t, goldenPath, expectation{
					Err:    actualErr,
					Output: actualOutput,
				})
				// Also update the VCR recording implicitly by running the test
				t.Logf("Updated golden file: %s", goldenPath)
				// We might need to re-run the test without the golden flag
				// after recording to ensure the comparison passes.
				// However, for now, just writing the golden file is the goal.
				return // Skip comparison when updating golden files
			}

			// Otherwise, read the golden file and compare
			expected := readGoldenFile(t, goldenPath)

			// Compare error
			if diff := cmp.Diff(expected.Err, actualErr); diff != "" {
				t.Errorf("Error mismatch (-want +got):\n%s", diff)
			}

			// Compare output (only if no error is expected)
			if expected.Err == "" && actualErr == "" {
				// Compare JSON strings directly
				if diff := cmp.Diff(expected.Output, actualOutput); diff != "" {
					// To make diffs easier to read, unmarshal and compare structures
					var expectedContents []mcp.ResourceContents
					var actualContents []mcp.ResourceContents
					json.Unmarshal([]byte(expected.Output), &expectedContents) // Ignore error for diffing
					json.Unmarshal([]byte(actualOutput), &actualContents)      // Ignore error for diffing
					t.Errorf("Output mismatch (-want +got):\n%s", cmp.Diff(expectedContents, actualContents))
					t.Logf("Expected JSON:\n%s", expected.Output)
					t.Logf("Actual JSON:\n%s", actualOutput)
				}
			} else if expected.Err == "" && actualErr != "" {
				t.Errorf("Expected no error, but got: %s", actualErr)
			} else if expected.Err != "" && actualErr == "" {
				t.Errorf("Expected error '%s', but got none", expected.Err)
			}
		})
	}
}

// readGoldenFile and writeGoldenFile are defined in test_helpers.go
