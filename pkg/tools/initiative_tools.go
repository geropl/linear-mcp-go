package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

var GetInitiativeTool = mcp.NewTool("linear_get_initiative",
	mcp.WithDescription("Get a single initiative by its identifier (ID or name)."),
	mcp.WithString("initiative", mcp.Required(), mcp.Description("The identifier of the initiative to get. Can be the initiative's ID or name.")),
)

func GetInitiativeHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		initiativeIdentifier, err := request.RequireString("initiative")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		initiative, err := linearClient.GetInitiative(initiativeIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get initiative: %v", err)}}}, nil
		}

		resultText := FormatInitiative(*initiative)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

func FormatInitiative(initiative linear.Initiative) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Initiative: %s\n", initiative.Name))
	builder.WriteString(fmt.Sprintf("  ID: %s\n", initiative.ID))
	if initiative.Description != "" {
		builder.WriteString(fmt.Sprintf("  Description: %s\n", initiative.Description))
	}
	builder.WriteString(fmt.Sprintf("  URL: %s\n", initiative.URL))
	return builder.String()
}

var CreateInitiativeTool = mcp.NewTool("linear_create_initiative",
	mcp.WithDescription("Create a new initiative."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the initiative.")),
	mcp.WithString("description", mcp.Description("The description of the initiative.")),
)

func CreateInitiativeHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		description := request.GetString("description", "")

		input := linear.InitiativeCreateInput{
			Name:        name,
			Description: description,
		}

		initiative, err := linearClient.CreateInitiative(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to create initiative: %v", err)}}}, nil
		}

		resultText := FormatInitiative(*initiative)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

var UpdateInitiativeTool = mcp.NewTool("linear_update_initiative",
	mcp.WithDescription("Update an existing initiative."),
	mcp.WithString("initiative", mcp.Required(), mcp.Description("The ID or name of the initiative to update.")),
	mcp.WithString("name", mcp.Description("The new name of the initiative.")),
	mcp.WithString("description", mcp.Description("The new description of the initiative.")),
)

func UpdateInitiativeHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		initiativeIdentifier, err := request.RequireString("initiative")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Get the initiative first to get its ID
		init, err := linearClient.GetInitiative(initiativeIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get initiative: %v", err)}}}, nil
		}

		name := request.GetString("name", "")
		description := request.GetString("description", "")

		input := linear.InitiativeUpdateInput{
			Name:        name,
			Description: description,
		}

		initiative, err := linearClient.UpdateInitiative(init.ID, input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to update initiative: %v", err)}}}, nil
		}

		resultText := FormatInitiative(*initiative)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
