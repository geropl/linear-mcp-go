package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

var GetMilestoneTool = mcp.NewTool("linear_get_milestone",
	mcp.WithDescription("Get a single project milestone by its ID or name."),
	mcp.WithString("milestone", mcp.Required(), mcp.Description("The ID or name of the project milestone to get.")),
)

func GetMilestoneHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		milestoneIdentifier, err := request.RequireString("milestone")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		milestone, err := linearClient.GetMilestone(milestoneIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get milestone: %v", err)}}}, nil
		}

		resultText := FormatMilestone(*milestone)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

func FormatMilestone(milestone linear.ProjectMilestone) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Milestone: %s\n", milestone.Name))
	builder.WriteString(fmt.Sprintf("  ID: %s\n", milestone.ID))
	if milestone.Description != "" {
		builder.WriteString(fmt.Sprintf("  Description: %s\n", milestone.Description))
	}
	if milestone.TargetDate != nil {
		builder.WriteString(fmt.Sprintf("  Target Date: %s\n", *milestone.TargetDate))
	}
	if milestone.Project != nil {
		builder.WriteString(fmt.Sprintf("  Project: %s (%s)\n", milestone.Project.Name, milestone.Project.ID))
	}
	return builder.String()
}

var UpdateMilestoneTool = mcp.NewTool("linear_update_milestone",
	mcp.WithDescription("Update an existing project milestone."),
	mcp.WithString("milestone", mcp.Required(), mcp.Description("The ID or name of the milestone to update.")),
	mcp.WithString("name", mcp.Description("The new name of the milestone.")),
	mcp.WithString("description", mcp.Description("The new description of the milestone.")),
	mcp.WithString("targetDate", mcp.Description("The new target date of the milestone (YYYY-MM-DD).")),
)

func UpdateMilestoneHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		milestoneIdentifier, err := request.RequireString("milestone")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Get the milestone first to get its ID
		mil, err := linearClient.GetMilestone(milestoneIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get milestone: %v", err)}}}, nil
		}

		name := request.GetString("name", "")
		description := request.GetString("description", "")
		targetDate := request.GetString("targetDate", "")

		input := linear.ProjectMilestoneUpdateInput{
			Name:        name,
			Description: description,
			TargetDate:  targetDate,
		}

		milestone, err := linearClient.UpdateMilestone(mil.ID, input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to update milestone: %v", err)}}}, nil
		}

		resultText := FormatMilestone(*milestone)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

var CreateMilestoneTool = mcp.NewTool("linear_create_milestone",
	mcp.WithDescription("Create a new project milestone."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the milestone.")),
	mcp.WithString("projectId", mcp.Required(), mcp.Description("The ID of the project to create the milestone in.")),
	mcp.WithString("description", mcp.Description("The description of the milestone.")),
	mcp.WithString("targetDate", mcp.Description("The target date of the milestone (YYYY-MM-DD).")),
)

func CreateMilestoneHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		projectID, err := request.RequireString("projectId")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		description := request.GetString("description", "")
		targetDate := request.GetString("targetDate", "")

		input := linear.ProjectMilestoneCreateInput{
			Name:        name,
			ProjectID:   projectID,
			Description: description,
			TargetDate:  targetDate,
		}

		milestone, err := linearClient.CreateMilestone(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to create milestone: %v", err)}}}, nil
		}

		resultText := FormatMilestone(*milestone)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}
