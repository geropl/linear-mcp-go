package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// TeamsResource is the resource definition for Linear teams
var TeamsResource = mcp.NewResource(
	"linear://teams",
	"Linear Teams",
	mcp.WithResourceDescription("List of teams in Linear"),
	mcp.WithMIMEType("application/json"),
)

// TeamResource is the resource definition for a specific Linear team
var TeamResource = mcp.NewResource(
	"linear://team/{id}",
	"Linear Team",
	mcp.WithResourceDescription("Details of a specific team in Linear"),
	mcp.WithMIMEType("application/json"),
)

// RegisterResources registers all Linear resources with the MCP server
func RegisterResources(s *mcpserver.MCPServer, linearClient *linear.LinearClient) {
	// Register Teams resource
	s.AddResource(TeamsResource, TeamsResourceHandler(linearClient))

	// Register Team resource
	s.AddResource(TeamResource, TeamResourceHandler(linearClient))
}

// TeamsResourceHandler handles the linear://teams resource
func TeamsResourceHandler(linearClient *linear.LinearClient) mcpserver.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Get teams from Linear
		teams, err := linearClient.GetTeams("")
		if err != nil {
			return nil, fmt.Errorf("failed to get teams: %v", err)
		}

		// Create resource content
		results := []mcp.ResourceContents{}
		for _, t := range teams {
			teamJSON, err := json.Marshal(t)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal team: %v", err)
			}

			results = append(results, mcp.TextResourceContents{
				URI:      fmt.Sprintf("linear://team/%s", t.ID),
				MIMEType: "application/json",
				Text:     string(teamJSON),
			})
		}

		return results, nil
	}
}

// TeamResourceHandler handles the linear://team/{id} resource
func TeamResourceHandler(linearClient *linear.LinearClient) mcpserver.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract team ID from URI
		uri := request.Params.URI
		if !strings.HasPrefix(uri, "linear://team/") {
			return nil, fmt.Errorf("invalid team URI: %s", uri)
		}

		teamID := uri[len("linear://team/"):]
		if teamID == "" {
			return nil, fmt.Errorf("team ID is required")
		}

		// Resolve team ID (could be UUID, name, or key)
		resolvedTeamID, err := resolveTeamIdentifier(linearClient, teamID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve team identifier: %v", err)
		}

		// Get all teams and find the matching one
		teams, err := linearClient.GetTeams("")
		if err != nil {
			return nil, fmt.Errorf("failed to get teams: %v", err)
		}

		var team *linear.Team
		for i, t := range teams {
			if t.ID == resolvedTeamID {
				team = &teams[i]
				break
			}
		}

		if team == nil {
			return nil, fmt.Errorf("team not found: %s", teamID)
		}

		// Format team as JSON
		teamJSON, err := json.Marshal(team)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal team: %v", err)
		}

		// Create resource content
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      fmt.Sprintf("linear://team/%s", team.ID),
				MIMEType: "application/json",
				Text:     string(teamJSON),
			},
		}, nil
	}
}

// resolveTeamIdentifier resolves a team identifier (UUID, name, or key) to a team ID
func resolveTeamIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find a team by name or key
	teams, err := linearClient.GetTeams("")
	if err != nil {
		return "", fmt.Errorf("failed to get teams: %v", err)
	}

	// First try exact match on name or key
	for _, team := range teams {
		if team.Name == identifier || team.Key == identifier {
			return team.ID, nil
		}
	}

	// If no exact match, try case-insensitive match
	identifierLower := strings.ToLower(identifier)
	for _, team := range teams {
		if strings.ToLower(team.Name) == identifierLower || strings.ToLower(team.Key) == identifierLower {
			return team.ID, nil
		}
	}

	return "", fmt.Errorf("no team found with identifier '%s'", identifier)
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(uuidStr string) bool {
	// Simple UUID validation - check if it has the correct format
	// This is a simplified version and doesn't validate the UUID fully
	return len(uuidStr) == 36 && strings.Count(uuidStr, "-") == 4
}
