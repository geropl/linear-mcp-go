package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/mark3labs/mcp-go/mcp"
)

var GetProjectTool = mcp.NewTool("linear_get_project",
	mcp.WithDescription("Get a single project."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The identifier of the project, either ID, name or slug.")),
)

func GetProjectHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectIdentifier, err := request.RequireString("project")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		project, err := linearClient.GetProject(projectIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get project: %v", err)}}}, nil
		}

		resultText := FormatProject(*project)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

var SearchProjectsTool = mcp.NewTool("linear_search_projects",
	mcp.WithDescription("Search for projects."),
	mcp.WithString("query", mcp.Description("The search query.")),
)

func SearchProjectsHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")

		if query == "" {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Query parameter may not be empty."}}}, nil
		}

		projects, err := linearClient.SearchProjects(query)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to search projects: %v", err)}}}, nil
		}

		if len(projects) == 0 {
			return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No projects found."}}}, nil
		}

		var builder strings.Builder
		for _, project := range projects {
			builder.WriteString(FormatProject(project))
			builder.WriteString("\n")
		}

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: builder.String()}}}, nil
	}
}

var CreateProjectTool = mcp.NewTool("linear_create_project",
	mcp.WithDescription("Create a new project."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the project.")),
	mcp.WithString("teamIds", mcp.Required(), mcp.Description("A comma-separated list of team IDs.")),
	mcp.WithString("description", mcp.Description("The description of the project.")),
	mcp.WithString("leadId", mcp.Description("The ID of the project lead.")),
	mcp.WithString("startDate", mcp.Description("The start date of the project (YYYY-MM-DD).")),
	mcp.WithString("targetDate", mcp.Description("The target date of the project (YYYY-MM-DD).")),
)

func CreateProjectHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		teamIDsStr, err := request.RequireString("teamIds")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}
		teamIDs := strings.Split(teamIDsStr, ",")

		description := request.GetString("description", "")
		leadID := request.GetString("leadId", "")
		startDate := request.GetString("startDate", "")
		targetDate := request.GetString("targetDate", "")

		input := linear.ProjectCreateInput{
			Name:        name,
			TeamIDs:     teamIDs,
			Description: description,
			LeadID:      leadID,
			StartDate:   startDate,
			TargetDate:  targetDate,
		}

		project, err := linearClient.CreateProject(input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to create project: %v", err)}}}, nil
		}

		resultText := FormatProject(*project)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

var UpdateProjectTool = mcp.NewTool("linear_update_project",
	mcp.WithDescription("Update an existing project."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The identifier of the project to update.")),
	mcp.WithString("name", mcp.Description("The new name of the project.")),
	mcp.WithString("description", mcp.Description("The new description of the project.")),
	mcp.WithString("leadId", mcp.Description("The ID of the project lead.")),
	mcp.WithString("startDate", mcp.Description("The start date of the project (YYYY-MM-DD).")),
	mcp.WithString("targetDate", mcp.Description("The target date of the project (YYYY-MM-DD).")),
	mcp.WithString("teamIds", mcp.Description("A comma-separated list of team IDs.")),
)

func UpdateProjectHandler(linearClient *linear.LinearClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectIdentifier, err := request.RequireString("project")
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}}}, nil
		}

		// Get the project first to get its ID
		proj, err := linearClient.GetProject(projectIdentifier)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to get project: %v", err)}}}, nil
		}

		name := request.GetString("name", "")
		description := request.GetString("description", "")
		leadID := request.GetString("leadId", "")
		startDate := request.GetString("startDate", "")
		targetDate := request.GetString("targetDate", "")
		teamIDsStr := request.GetString("teamIds", "")
		var teamIDs []string
		if teamIDsStr != "" {
			teamIDs = strings.Split(teamIDsStr, ",")
		}

		input := linear.ProjectUpdateInput{
			Name:        name,
			Description: description,
			LeadID:      leadID,
			StartDate:   startDate,
			TargetDate:  targetDate,
			TeamIDs:     teamIDs,
		}

		project, err := linearClient.UpdateProject(proj.ID, input)
		if err != nil {
			return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to update project: %v", err)}}}, nil
		}

		resultText := FormatProject(*project)

		return &mcp.CallToolResult{Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}}}, nil
	}
}

func FormatProject(project linear.Project) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Project: %s\n", project.Name))
	builder.WriteString(fmt.Sprintf("  ID: %s\n", project.ID))
	builder.WriteString(fmt.Sprintf("  State: %s\n", project.State))
	builder.WriteString(fmt.Sprintf("  URL: %s\n", project.URL))
	if project.Description != "" {
		builder.WriteString(fmt.Sprintf("  Description: %s\n", project.Description))
	}
	if project.Lead == nil {
		builder.WriteString("  Lead: None\n")
	} else {
		builder.WriteString(fmt.Sprintf("  Lead: %s\n", project.Lead.Name))
	}
	if project.StartDate == nil {
		builder.WriteString("  Start Date: None\n")
	} else {
		builder.WriteString(fmt.Sprintf("  Start Date: %s\n", *project.StartDate))
	}
	if project.TargetDate == nil {
		builder.WriteString("  Target Date: None\n")
	} else {
		builder.WriteString(fmt.Sprintf("  Target Date: %s\n", *project.TargetDate))
	}
	if project.Initiatives != nil && len(project.Initiatives.Nodes) > 0 {
		builder.WriteString("  Initiatives:\n")
		for _, i := range project.Initiatives.Nodes {
			builder.WriteString(fmt.Sprintf("    - %s (ID: %s)\n", i.Name, i.ID))
		}
	} else {
		builder.WriteString("  Initiatives: None\n")
	}
	return builder.String()
}
