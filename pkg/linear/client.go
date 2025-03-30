package linear

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LinearAPIEndpoint = "https://api.linear.app/graphql"
	UserAgent         = "linear-mcp-go/1.0.0"
)

// LinearClient is a client for the Linear API
type LinearClient struct {
	apiKey      string
	httpClient  *http.Client
	rateLimiter *RateLimiter
}

// NewLinearClient creates a new Linear API client
func NewLinearClient(apiKey string) (*LinearClient, error) {
	if apiKey == "" {
		return nil, errors.New("LINEAR_API_KEY environment variable is required")
	}

	return &LinearClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: NewRateLimiter(1400), // Linear API limit is 1400 requests per hour
	}, nil
}

// NewLinearClientFromEnv creates a new Linear API client from environment variables
func NewLinearClientFromEnv() (*LinearClient, error) {
	apiKey := os.Getenv("LINEAR_API_KEY")
	return NewLinearClient(apiKey)
}

// executeGraphQL executes a GraphQL query against the Linear API
func (c *LinearClient) executeGraphQL(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	// Create the request body
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Marshal the request body to JSON
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", LinearAPIEndpoint, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("User-Agent", UserAgent)

	// Execute the request with rate limiting
	var resp *http.Response
	err = c.rateLimiter.Enqueue(func() error {
		var reqErr error
		resp, reqErr = c.httpClient.Do(req)
		return reqErr
	}, "graphql")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var graphQLResp GraphQLResponse
	if err := json.Unmarshal(respBody, &graphQLResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for GraphQL errors
	if len(graphQLResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", graphQLResp.Errors[0].Message)
	}

	return &graphQLResp, nil
}

// GetIssue gets an issue by ID
func (c *LinearClient) GetIssue(issueID string) (*Issue, error) {
	query := `
		query GetIssue($id: String!) {
			issue(id: $id) {
				id
				identifier
				title
				description
				priority
				url
				createdAt
				updatedAt
				state {
					id
					name
				}
				assignee {
					id
					name
					email
				}
				team {
					id
					name
					key
				}
				comments(first: 100) {
					nodes {
						id
						body
						createdAt
						parent {
							id
						}
						user {
							id
							name
						}
						children(first: 100) {
							nodes {
								id
								body
								createdAt
								user {
									id
									name
								}
							}
						}
					}
				}
				relations(first: 20) {
					nodes {
						id
						type
						relatedIssue {
							id
							identifier
							title
							url
						}
					}
				}
				inverseRelations(first: 20) {
					nodes {
						id
						type
						issue {
							id
							identifier
							title
							url
						}
					}
				}
				attachments(first: 50) {
					nodes {
						id
						title
						subtitle
						url
						sourceType
						metadata
						createdAt
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": issueID,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the issue from the response
	issueData, ok := resp.Data["issue"].(map[string]interface{})
	if !ok || issueData == nil {
		return nil, fmt.Errorf("issue %s not found", issueID)
	}

	// Parse the issue data
	var issue Issue
	issueBytes, err := json.Marshal(issueData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue data: %w", err)
	}

	if err := json.Unmarshal(issueBytes, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return &issue, nil
}

// GetIssueByIdentifier gets an issue by its identifier (e.g., "TEAM-123")
func (c *LinearClient) GetIssueByIdentifier(identifier string) (*Issue, error) {
	// Split the identifier into team key and number parts
	parts := strings.Split(identifier, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid issue identifier format: %s (expected format: TEAM-123)", identifier)
	}

	teamKey := parts[0]
	numberStr := parts[1]

	// Convert the number part to an integer
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("invalid issue number in identifier: %s", identifier)
	}

	// Use the issues query with filters for team key and number
	query := `
		query GetIssueByIdentifier($teamKey: String!, $number: Float!) {
			issues(filter: { team: { key: { eq: $teamKey } }, number: { eq: $number } }, first: 1) {
				nodes {
					id
					identifier
					title
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamKey": teamKey,
		"number":  float64(number),
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the issues from the response
	issuesData, ok := resp.Data["issues"].(map[string]interface{})
	if !ok || issuesData == nil {
		return nil, fmt.Errorf("issue search failed for identifier %s", identifier)
	}

	nodesData, ok := issuesData["nodes"].([]interface{})
	if !ok || nodesData == nil || len(nodesData) == 0 {
		return nil, fmt.Errorf("no issue found with identifier %s", identifier)
	}

	// Get the first issue
	issueData, ok := nodesData[0].(map[string]interface{})
	if !ok || issueData == nil {
		return nil, fmt.Errorf("invalid issue data for identifier %s", identifier)
	}

	// Parse the issue data
	var issue Issue
	issueBytes, err := json.Marshal(issueData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue data: %w", err)
	}

	if err := json.Unmarshal(issueBytes, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return &issue, nil
}

// GetLabelsByName gets labels by name for a team
func (c *LinearClient) GetLabelsByName(teamID string, labelNames []string) ([]Label, error) {
	query := `
		query GetLabelsByName($teamId: String!, $names: [String!]!) {
			team(id: $teamId) {
				labels(filter: { name: { in: $names } }) {
					nodes {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId": teamID,
		"names":  labelNames,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the team from the response
	teamData, ok := resp.Data["team"].(map[string]interface{})
	if !ok || teamData == nil {
		return nil, fmt.Errorf("team %s not found", teamID)
	}

	// Extract the labels
	labelsData, ok := teamData["labels"].(map[string]interface{})
	if !ok || labelsData == nil {
		return []Label{}, nil
	}

	nodesData, ok := labelsData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []Label{}, nil
	}

	// Parse the labels data
	labels := make([]Label, 0, len(nodesData))
	for _, nodeData := range nodesData {
		labelData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		label := Label{
			ID:   getStringValue(labelData, "id"),
			Name: getStringValue(labelData, "name"),
		}

		labels = append(labels, label)
	}

	return labels, nil
}

// CreateIssue creates a new issue
func (c *LinearClient) CreateIssue(input CreateIssueInput) (*Issue, error) {
	query := `
		mutation CreateIssue($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				success
				issue {
					id
					identifier
					title
					description
					priority
					url
					createdAt
					updatedAt
					state {
						id
						name
					}
					team {
						id
						name
						key
					}
					labels {
						nodes {
							id
							name
						}
					}
				}
			}
		}
	`

	// Prepare variables
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":       input.Title,
			"teamId":      input.TeamID,
			"description": input.Description,
		},
	}

	if input.Priority != nil {
		variables["input"].(map[string]interface{})["priority"] = *input.Priority
	}

	if input.Status != "" {
		variables["input"].(map[string]interface{})["stateId"] = input.Status
	}

	if input.ParentID != nil && *input.ParentID != "" {
		variables["input"].(map[string]interface{})["parentId"] = *input.ParentID
	}

	if len(input.LabelIDs) > 0 {
		variables["input"].(map[string]interface{})["labelIds"] = input.LabelIDs
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the issue from the response
	issueCreateData, ok := resp.Data["issueCreate"].(map[string]interface{})
	if !ok || issueCreateData == nil {
		return nil, errors.New("failed to create issue")
	}

	success, ok := issueCreateData["success"].(bool)
	if !ok || !success {
		return nil, errors.New("failed to create issue")
	}

	issueData, ok := issueCreateData["issue"].(map[string]interface{})
	if !ok || issueData == nil {
		return nil, errors.New("failed to create issue")
	}

	// Parse the issue data
	var issue Issue
	issueBytes, err := json.Marshal(issueData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue data: %w", err)
	}

	if err := json.Unmarshal(issueBytes, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return &issue, nil
}

// UpdateIssue updates an existing issue
func (c *LinearClient) UpdateIssue(input UpdateIssueInput) (*Issue, error) {
	query := `
		mutation UpdateIssue($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				success
				issue {
					id
					identifier
					title
					description
					priority
					url
					createdAt
					updatedAt
					state {
						id
						name
					}
					team {
						id
						name
						key
					}
				}
			}
		}
	`

	// Prepare variables
	updateInput := map[string]interface{}{}

	if input.Title != "" {
		updateInput["title"] = input.Title
	}

	if input.Description != "" {
		updateInput["description"] = input.Description
	}

	if input.Priority != nil {
		updateInput["priority"] = *input.Priority
	}

	if input.Status != "" {
		updateInput["stateId"] = input.Status
	}

	variables := map[string]interface{}{
		"id":    input.ID,
		"input": updateInput,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the issue from the response
	issueUpdateData, ok := resp.Data["issueUpdate"].(map[string]interface{})
	if !ok || issueUpdateData == nil {
		return nil, errors.New("failed to update issue")
	}

	success, ok := issueUpdateData["success"].(bool)
	if !ok || !success {
		return nil, errors.New("failed to update issue")
	}

	issueData, ok := issueUpdateData["issue"].(map[string]interface{})
	if !ok || issueData == nil {
		return nil, errors.New("failed to update issue")
	}

	// Parse the issue data
	var issue Issue
	issueBytes, err := json.Marshal(issueData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue data: %w", err)
	}

	if err := json.Unmarshal(issueBytes, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return &issue, nil
}

// SearchIssues searches for issues with filters
func (c *LinearClient) SearchIssues(input SearchIssuesInput) ([]LinearIssueResponse, error) {
	query := `
		query SearchIssues($filter: IssueFilter, $first: Int, $includeArchived: Boolean) {
			issues(filter: $filter, first: $first, includeArchived: $includeArchived) {
				nodes {
					id
					identifier
					title
					description
					priority
					url
					state {
						id
						name
					}
					assignee {
						id
						name
					}
					labels {
						nodes {
							id
							name
						}
					}
				}
			}
		}
	`

	// Build the filter
	filter := map[string]interface{}{}

	if input.Query != "" {
		filter["or"] = []map[string]interface{}{
			{"title": map[string]interface{}{"contains": input.Query}},
			{"description": map[string]interface{}{"contains": input.Query}},
		}
	}

	if input.TeamID != "" {
		filter["team"] = map[string]interface{}{
			"id": map[string]interface{}{"eq": input.TeamID},
		}
	}

	if input.Status != "" {
		filter["state"] = map[string]interface{}{
			"name": map[string]interface{}{"eq": input.Status},
		}
	}

	if input.AssigneeID != "" {
		filter["assignee"] = map[string]interface{}{
			"id": map[string]interface{}{"eq": input.AssigneeID},
		}
	}

	if len(input.Labels) > 0 {
		filter["labels"] = map[string]interface{}{
			"some": map[string]interface{}{
				"name": map[string]interface{}{"in": input.Labels},
			},
		}
	}

	if input.Priority != nil {
		filter["priority"] = map[string]interface{}{"eq": *input.Priority}
	}

	if input.Estimate != nil {
		filter["estimate"] = map[string]interface{}{"eq": *input.Estimate}
	}

	// Set default limit if not provided
	limit := 10
	if input.Limit > 0 {
		limit = input.Limit
	}

	variables := map[string]interface{}{
		"filter":          filter,
		"first":           limit,
		"includeArchived": input.IncludeArchived,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the issues from the response
	issuesData, ok := resp.Data["issues"].(map[string]interface{})
	if !ok || issuesData == nil {
		return []LinearIssueResponse{}, nil
	}

	nodesData, ok := issuesData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []LinearIssueResponse{}, nil
	}

	// Parse the issues data
	issues := make([]LinearIssueResponse, 0, len(nodesData))
	for _, nodeData := range nodesData {
		issueData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract state name
		var stateName string
		if stateData, ok := issueData["state"].(map[string]interface{}); ok && stateData != nil {
			if name, ok := stateData["name"].(string); ok {
				stateName = name
			}
		}

		// Create the issue response
		issue := LinearIssueResponse{
			ID:         getStringValue(issueData, "id"),
			Identifier: getStringValue(issueData, "identifier"),
			Title:      getStringValue(issueData, "title"),
			URL:        getStringValue(issueData, "url"),
			StateName:  stateName,
		}

		// Extract priority
		if priority, ok := issueData["priority"].(float64); ok {
			issue.Priority = int(priority)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// GetUserIssues gets issues assigned to a user
func (c *LinearClient) GetUserIssues(input GetUserIssuesInput) ([]LinearIssueResponse, error) {
	var userID string
	var err error

	if input.UserID == "" {
		// Get the current user's ID
		userID, err = c.getCurrentUserID()
		if err != nil {
			return nil, err
		}
	} else {
		userID = input.UserID
	}

	query := `
		query GetUserIssues($userId: String!, $first: Int, $includeArchived: Boolean) {
			user(id: $userId) {
				assignedIssues(first: $first, includeArchived: $includeArchived) {
					nodes {
						id
						identifier
						title
						description
						priority
						url
						state {
							id
							name
						}
					}
				}
			}
		}
	`

	// Set default limit if not provided
	limit := 50
	if input.Limit > 0 {
		limit = input.Limit
	}

	variables := map[string]interface{}{
		"userId":          userID,
		"first":           limit,
		"includeArchived": input.IncludeArchived,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the user from the response
	userData, ok := resp.Data["user"].(map[string]interface{})
	if !ok || userData == nil {
		return nil, fmt.Errorf("user %s not found", userID)
	}

	// Extract the assigned issues
	assignedIssuesData, ok := userData["assignedIssues"].(map[string]interface{})
	if !ok || assignedIssuesData == nil {
		return []LinearIssueResponse{}, nil
	}

	nodesData, ok := assignedIssuesData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []LinearIssueResponse{}, nil
	}

	// Parse the issues data
	issues := make([]LinearIssueResponse, 0, len(nodesData))
	for _, nodeData := range nodesData {
		issueData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract state name
		var stateName string
		if stateData, ok := issueData["state"].(map[string]interface{}); ok && stateData != nil {
			if name, ok := stateData["name"].(string); ok {
				stateName = name
			}
		}

		// Create the issue response
		issue := LinearIssueResponse{
			ID:         getStringValue(issueData, "id"),
			Identifier: getStringValue(issueData, "identifier"),
			Title:      getStringValue(issueData, "title"),
			URL:        getStringValue(issueData, "url"),
			StateName:  stateName,
		}

		// Extract priority
		if priority, ok := issueData["priority"].(float64); ok {
			issue.Priority = int(priority)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// AddComment adds a comment to an issue
func (c *LinearClient) AddComment(input AddCommentInput) (*Comment, *Issue, error) {
	query := `
		mutation AddComment($input: CommentCreateInput!) {
			commentCreate(input: $input) {
				success
				comment {
					id
					body
					url
					createdAt
					user {
						id
						name
					}
					issue {
						id
						identifier
						title
						url
					}
				}
			}
		}
	`

	// Prepare variables
	commentInput := map[string]interface{}{
		"issueId": input.IssueID,
		"body":    input.Body,
	}

	if input.CreateAsUser != "" {
		commentInput["createAsUser"] = input.CreateAsUser
	}

	if input.DisplayIconURL != "" {
		commentInput["displayIconUrl"] = input.DisplayIconURL
	}

	variables := map[string]interface{}{
		"input": commentInput,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, nil, err
	}

	// Extract the comment from the response
	commentCreateData, ok := resp.Data["commentCreate"].(map[string]interface{})
	if !ok || commentCreateData == nil {
		return nil, nil, errors.New("failed to create comment")
	}

	success, ok := commentCreateData["success"].(bool)
	if !ok || !success {
		return nil, nil, errors.New("failed to create comment")
	}

	commentData, ok := commentCreateData["comment"].(map[string]interface{})
	if !ok || commentData == nil {
		return nil, nil, errors.New("failed to create comment")
	}

	issueData, ok := commentData["issue"].(map[string]interface{})
	if !ok || issueData == nil {
		return nil, nil, errors.New("failed to get issue for comment")
	}

	// Parse the comment data
	var comment Comment
	commentBytes, err := json.Marshal(commentData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal comment data: %w", err)
	}

	if err := json.Unmarshal(commentBytes, &comment); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal comment data: %w", err)
	}

	// Parse the issue data
	var issue Issue
	issueBytes, err := json.Marshal(issueData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal issue data: %w", err)
	}

	if err := json.Unmarshal(issueBytes, &issue); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal issue data: %w", err)
	}

	return &comment, &issue, nil
}

// GetTeamIssues gets issues for a team
func (c *LinearClient) GetTeamIssues(teamID string) ([]LinearIssueResponse, error) {
	query := `
		query GetTeamIssues($teamId: ID!) {
			team(id: $teamId) {
				issues {
					nodes {
						id
						identifier
						title
						description
						priority
						url
						state {
							id
							name
						}
						assignee {
							id
							name
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId": teamID,
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the team from the response
	teamData, ok := resp.Data["team"].(map[string]interface{})
	if !ok || teamData == nil {
		return nil, fmt.Errorf("team %s not found", teamID)
	}

	// Extract the issues
	issuesData, ok := teamData["issues"].(map[string]interface{})
	if !ok || issuesData == nil {
		return []LinearIssueResponse{}, nil
	}

	nodesData, ok := issuesData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []LinearIssueResponse{}, nil
	}

	// Parse the issues data
	issues := make([]LinearIssueResponse, 0, len(nodesData))
	for _, nodeData := range nodesData {
		issueData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract state name
		var stateName string
		if stateData, ok := issueData["state"].(map[string]interface{}); ok && stateData != nil {
			if name, ok := stateData["name"].(string); ok {
				stateName = name
			}
		}

		// Create the issue response
		issue := LinearIssueResponse{
			ID:         getStringValue(issueData, "id"),
			Identifier: getStringValue(issueData, "identifier"),
			Title:      getStringValue(issueData, "title"),
			URL:        getStringValue(issueData, "url"),
			StateName:  stateName,
		}

		// Extract priority
		if priority, ok := issueData["priority"].(float64); ok {
			issue.Priority = int(priority)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// GetViewer gets the current user
func (c *LinearClient) GetViewer() (*User, []Team, *Organization, error) {
	query := `
		query GetViewer {
			viewer {
				id
				name
				email
				admin
				teams {
					nodes {
						id
						name
						key
					}
				}
				organization {
					id
					name
					urlKey
				}
			}
		}
	`

	resp, err := c.executeGraphQL(query, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	// Extract the viewer from the response
	viewerData, ok := resp.Data["viewer"].(map[string]interface{})
	if !ok || viewerData == nil {
		return nil, nil, nil, errors.New("failed to get viewer")
	}

	// Parse the user data
	var user User
	user.ID = getStringValue(viewerData, "id")
	user.Name = getStringValue(viewerData, "name")
	user.Email = getStringValue(viewerData, "email")
	if admin, ok := viewerData["admin"].(bool); ok {
		user.Admin = admin
	}

	// Extract teams
	var teams []Team
	if teamsData, ok := viewerData["teams"].(map[string]interface{}); ok && teamsData != nil {
		if nodesData, ok := teamsData["nodes"].([]interface{}); ok && nodesData != nil {
			teams = make([]Team, 0, len(nodesData))
			for _, nodeData := range nodesData {
				teamData, ok := nodeData.(map[string]interface{})
				if !ok {
					continue
				}

				team := Team{
					ID:   getStringValue(teamData, "id"),
					Name: getStringValue(teamData, "name"),
					Key:  getStringValue(teamData, "key"),
				}
				teams = append(teams, team)
			}
		}
	}

	// Extract organization
	var org Organization
	if orgData, ok := viewerData["organization"].(map[string]interface{}); ok && orgData != nil {
		org.ID = getStringValue(orgData, "id")
		org.Name = getStringValue(orgData, "name")
		org.URLKey = getStringValue(orgData, "urlKey")
	}

	return &user, teams, &org, nil
}

// GetOrganization gets the organization
func (c *LinearClient) GetOrganization() (*Organization, error) {
	query := `
		query GetOrganization {
			organization {
				id
				name
				urlKey
				teams {
					nodes {
						id
						name
						key
					}
				}
				users {
					nodes {
						id
						name
						email
						admin
						active
					}
				}
			}
		}
	`

	resp, err := c.executeGraphQL(query, nil)
	if err != nil {
		return nil, err
	}

	// Extract the organization from the response
	orgData, ok := resp.Data["organization"].(map[string]interface{})
	if !ok || orgData == nil {
		return nil, errors.New("failed to get organization")
	}

	// Parse the organization data
	var org Organization
	org.ID = getStringValue(orgData, "id")
	org.Name = getStringValue(orgData, "name")
	org.URLKey = getStringValue(orgData, "urlKey")

	// Extract teams
	if teamsData, ok := orgData["teams"].(map[string]interface{}); ok && teamsData != nil {
		if nodesData, ok := teamsData["nodes"].([]interface{}); ok && nodesData != nil {
			org.Teams = make([]Team, 0, len(nodesData))
			for _, nodeData := range nodesData {
				teamData, ok := nodeData.(map[string]interface{})
				if !ok {
					continue
				}

				team := Team{
					ID:   getStringValue(teamData, "id"),
					Name: getStringValue(teamData, "name"),
					Key:  getStringValue(teamData, "key"),
				}
				org.Teams = append(org.Teams, team)
			}
		}
	}

	// Extract users
	if usersData, ok := orgData["users"].(map[string]interface{}); ok && usersData != nil {
		if nodesData, ok := usersData["nodes"].([]interface{}); ok && nodesData != nil {
			org.Users = make([]User, 0, len(nodesData))
			for _, nodeData := range nodesData {
				userData, ok := nodeData.(map[string]interface{})
				if !ok {
					continue
				}

				user := User{
					ID:    getStringValue(userData, "id"),
					Name:  getStringValue(userData, "name"),
					Email: getStringValue(userData, "email"),
				}

				if admin, ok := userData["admin"].(bool); ok {
					user.Admin = admin
				}

				org.Users = append(org.Users, user)
			}
		}
	}

	return &org, nil
}

// ListIssues lists issues
func (c *LinearClient) ListIssues() ([]LinearIssueResponse, error) {
	query := `
		query ListIssues {
			issues(first: 50, orderBy: updatedAt) {
				nodes {
					id
					identifier
					title
					priority
					url
					state {
						name
					}
					assignee {
						name
					}
					team {
						name
					}
				}
			}
		}
	`

	resp, err := c.executeGraphQL(query, nil)
	if err != nil {
		return nil, err
	}

	// Extract the issues from the response
	issuesData, ok := resp.Data["issues"].(map[string]interface{})
	if !ok || issuesData == nil {
		return []LinearIssueResponse{}, nil
	}

	nodesData, ok := issuesData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []LinearIssueResponse{}, nil
	}

	// Parse the issues data
	issues := make([]LinearIssueResponse, 0, len(nodesData))
	for _, nodeData := range nodesData {
		issueData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract state name
		var stateName string
		if stateData, ok := issueData["state"].(map[string]interface{}); ok && stateData != nil {
			if name, ok := stateData["name"].(string); ok {
				stateName = name
			}
		}

		// Create the issue response
		issue := LinearIssueResponse{
			ID:         getStringValue(issueData, "id"),
			Identifier: getStringValue(issueData, "identifier"),
			Title:      getStringValue(issueData, "title"),
			URL:        getStringValue(issueData, "url"),
			StateName:  stateName,
		}

		// Extract priority
		if priority, ok := issueData["priority"].(float64); ok {
			issue.Priority = int(priority)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// getCurrentUserID gets the current user's ID
func (c *LinearClient) getCurrentUserID() (string, error) {
	query := `
		query GetCurrentUser {
			viewer {
				id
			}
		}
	`

	resp, err := c.executeGraphQL(query, nil)
	if err != nil {
		return "", err
	}

	// Extract the viewer from the response
	viewerData, ok := resp.Data["viewer"].(map[string]interface{})
	if !ok || viewerData == nil {
		return "", errors.New("failed to get current user")
	}

	// Extract the ID
	id, ok := viewerData["id"].(string)
	if !ok || id == "" {
		return "", errors.New("failed to get current user ID")
	}

	return id, nil
}

// GetTeams gets teams by name (optional filter)
func (c *LinearClient) GetTeams(name string) ([]Team, error) {
	query := `
		query GetTeams($filter: TeamFilter) {
			teams(filter: $filter) {
				nodes {
					id
					name
					key
					description
					states {
						nodes {
							id
							name
						}
					}
				}
			}
		}
	`

	// Build the filter
	variables := map[string]interface{}{}

	if name != "" {
		variables["filter"] = map[string]interface{}{
			"name": map[string]interface{}{
				"contains": name,
			},
		}
	}

	resp, err := c.executeGraphQL(query, variables)
	if err != nil {
		return nil, err
	}

	// Extract the teams from the response
	teamsData, ok := resp.Data["teams"].(map[string]interface{})
	if !ok || teamsData == nil {
		return []Team{}, nil
	}

	nodesData, ok := teamsData["nodes"].([]interface{})
	if !ok || nodesData == nil {
		return []Team{}, nil
	}

	// Parse the teams data
	teams := make([]Team, 0, len(nodesData))
	for _, nodeData := range nodesData {
		teamData, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		team := Team{
			ID:   getStringValue(teamData, "id"),
			Name: getStringValue(teamData, "name"),
			Key:  getStringValue(teamData, "key"),
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// GetMetrics returns metrics about the API usage
func (c *LinearClient) GetMetrics() APIMetrics {
	metrics := c.rateLimiter.GetMetrics()

	return APIMetrics{
		RequestsInLastHour: metrics.RequestsInLastHour,
		RemainingRequests:  c.rateLimiter.requestsPerHour - metrics.RequestsInLastHour,
		AverageRequestTime: fmt.Sprintf("%dms", metrics.AverageRequestTime),
		QueueLength:        metrics.QueueLength,
		LastRequestTime:    time.Unix(0, metrics.LastRequestTime*int64(time.Millisecond)).Format(time.RFC3339),
	}
}

// Helper function to safely extract string values from maps
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
