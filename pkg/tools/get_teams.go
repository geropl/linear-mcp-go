package tools

import (
	"context"
	"fmt"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetTeamsTool is the tool definition for getting teams
var GetTeamsTool = mcp.NewTool("linear_get_teams",
	mcp.WithDescription("Retrieves Linear teams."),
	mcp.WithString("name", mcp.Description("Optional team name filter. Returns teams whose names contain this string.")),
)

// GetTeamsHandler handles the linear_get_teams tool
func GetTeamsHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := request.Params.Arguments

		// Extract optional name filter
		name := ""
		if nameArg, ok := args["name"].(string); ok {
			name = nameArg
		}

		// Get teams
		teams, err := linearClient.GetTeams(name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get teams: %v", err)), nil
		}

		// Format the result
		resultText := fmt.Sprintf("Found %d teams:\n", len(teams))
		for _, team := range teams {
			// Create a pointer to the team for formatTeamIdentifier
			teamPtr := &team
			resultText += fmt.Sprintf("- %s\n", formatTeamIdentifier(teamPtr))
			resultText += fmt.Sprintf("  Key: %s\n", team.Key)
		}

		return mcp.NewToolResultText(resultText), nil
	}
}
