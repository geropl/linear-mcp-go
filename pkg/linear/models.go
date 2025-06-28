package linear

import "time"

// Issue represents a Linear issue
type Issue struct {
	ID              string                   `json:"id"`
	Identifier      string                   `json:"identifier"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Priority        int                      `json:"priority"`
	Status          string                   `json:"status"`
	Assignee        *User                    `json:"assignee,omitempty"`
	Team            *Team                    `json:"team,omitempty"`
	Project         *Project                 `json:"project,omitempty"`
	ProjectMilestone *ProjectMilestone       `json:"projectMilestone,omitempty"`
	URL             string                   `json:"url"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	Labels          *LabelConnection         `json:"labels,omitempty"`
	State           *State                   `json:"state,omitempty"`
	Estimate        *float64                 `json:"estimate,omitempty"`
	Comments        *CommentConnection       `json:"comments,omitempty"`
	Relations       *IssueRelationConnection `json:"relations,omitempty"`
	InverseRelations *IssueRelationConnection `json:"inverseRelations,omitempty"`
	Attachments     *AttachmentConnection    `json:"attachments,omitempty"`
}

// User represents a Linear user
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

// Team represents a Linear team
type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Project represents a Linear project
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SlugID      string `json:"slugId"`
	State       string `json:"state"`
	Creator     *User  `json:"creator,omitempty"`
	Lead        *User  `json:"lead,omitempty"`
	// Members     *UserConnection `json:"members,omitempty"`
	// Teams       *TeamConnection `json:"teams,omitempty"`
	StartDate  *time.Time `json:"startDate,omitempty"`
	TargetDate *time.Time `json:"targetDate,omitempty"`
	Color      string     `json:"color"`
	Icon       string     `json:"icon,omitempty"`
	URL        string     `json:"url"`
}

// ProjectConnection represents a connection of projects
type ProjectConnection struct {
	Nodes []Project `json:"nodes"`
}

// ProjectMilestoneConnection represents a connection of project milestones.
type ProjectMilestoneConnection struct {
	Nodes []ProjectMilestone `json:"nodes"`
}

// ProjectMilestone represents a Linear project milestone
type ProjectMilestone struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	TargetDate  *string  `json:"targetDate,omitempty"`
	Project     *Project `json:"project,omitempty"`
	SortOrder   float64  `json:"sortOrder"`
}

// InitiativeConnection represents a connection of initiatives.
type InitiativeConnection struct {
	Nodes []Initiative `json:"nodes"`
}

// Initiative represents a Linear initiative
type Initiative struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       *User  `json:"owner,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	SlugID      string `json:"slugId"`
	URL         string `json:"url"`
}

// State represents a workflow state in Linear
type State struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LabelConnection represents a connection of labels
type LabelConnection struct {
	Nodes []Label `json:"nodes"`
}

// Label represents a Linear issue label
type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CommentConnection represents a connection of comments
type CommentConnection struct {
	Nodes []Comment `json:"nodes"`
}

// PageInfo represents pagination information
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// PaginatedCommentConnection represents a paginated connection of comments
type PaginatedCommentConnection struct {
	Nodes      []Comment `json:"nodes"`
	PageInfo   PageInfo  `json:"pageInfo"`
}

// Comment represents a comment on a Linear issue
type Comment struct {
	ID         string            `json:"id"`
	Body       string            `json:"body"`
	User       *User             `json:"user,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	URL        string            `json:"url,omitempty"`
	Parent     *Comment          `json:"parent,omitempty"`
	Children   *CommentConnection `json:"children,omitempty"`
}

// IssueRelationConnection represents a connection of issue relations
type IssueRelationConnection struct {
	Nodes []IssueRelation `json:"nodes"`
}

// IssueRelation represents a relation between two issues
type IssueRelation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	RelatedIssue *Issue `json:"relatedIssue,omitempty"`
	Issue        *Issue `json:"issue,omitempty"`
}

// AttachmentConnection represents a connection of attachments
type AttachmentConnection struct {
	Nodes []Attachment `json:"nodes"`
}

// Attachment represents an external resource linked to an issue
type Attachment struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Subtitle   string                 `json:"subtitle,omitempty"`
	URL        string                 `json:"url"`
	SourceType string                 `json:"sourceType,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
}

// Organization represents a Linear organization
type Organization struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	URLKey string  `json:"urlKey"`
	Teams  []Team  `json:"teams,omitempty"`
	Users  []User  `json:"users,omitempty"`
}

// LinearIssueResponse represents a simplified issue response
type LinearIssueResponse struct {
	ID               string            `json:"id"`
	Identifier       string            `json:"identifier"`
	Title            string            `json:"title"`
	Priority         int               `json:"priority"`
	Status           string            `json:"status,omitempty"`
	StateName        string            `json:"stateName,omitempty"`
	URL              string            `json:"url"`
	Project          *Project          `json:"project,omitempty"`
	ProjectMilestone *ProjectMilestone `json:"projectMilestone,omitempty"`
}

// APIMetrics represents metrics about API usage
type APIMetrics struct {
	RequestsInLastHour int    `json:"requestsInLastHour"`
	RemainingRequests  int    `json:"remainingRequests"`
	AverageRequestTime string `json:"averageRequestTime"`
	QueueLength        int    `json:"queueLength"`
	LastRequestTime    string `json:"lastRequestTime"`
}

// CreateIssueInput represents input for creating an issue
type CreateIssueInput struct {
	Title       string   `json:"title"`
	TeamID      string   `json:"teamId"`
	Description string   `json:"description,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Status      string   `json:"status,omitempty"`
	ParentID    *string  `json:"parentId,omitempty"`
	LabelIDs    []string `json:"labelIds,omitempty"`
}

// UpdateIssueInput represents input for updating an issue
type UpdateIssueInput struct {
	ID          string  `json:"id"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	Status      string  `json:"status,omitempty"`
	TeamID   	string  `json:"teamId,omitempty"`
	ProjectID   string  `json:"projectId,omitempty"`
	MilestoneID string  `json:"milestoneId,omitempty"`
}

// SearchIssuesInput represents input for searching issues
type SearchIssuesInput struct {
	Query           string   `json:"query,omitempty"`
	TeamID          string   `json:"teamId,omitempty"`
	Status          string   `json:"status,omitempty"`
	AssigneeID      string   `json:"assigneeId,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	Priority        *int     `json:"priority,omitempty"`
	Estimate        *float64 `json:"estimate,omitempty"`
	IncludeArchived bool     `json:"includeArchived,omitempty"`
	Limit           int      `json:"limit,omitempty"`
}

// GetUserIssuesInput represents input for getting user issues
type GetUserIssuesInput struct {
	UserID          string `json:"userId,omitempty"`
	IncludeArchived bool   `json:"includeArchived,omitempty"`
	Limit           int    `json:"limit,omitempty"`
}

// GetIssueCommentsInput represents input for getting issue comments
type GetIssueCommentsInput struct {
	IssueID    string `json:"issueId"`
	ParentID   string `json:"parentId,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	AfterCursor string `json:"afterCursor,omitempty"`
}

// AddCommentInput represents input for adding a comment
type AddCommentInput struct {
	IssueID      string `json:"issueId"`
	Body         string `json:"body"`
	CreateAsUser string `json:"createAsUser,omitempty"`
	ParentID     string `json:"parentId,omitempty"`
}

// ProjectCreateInput represents the input for creating a project.
type ProjectCreateInput struct {
	Name        string   `json:"name"`
	TeamIDs     []string `json:"teamIds"`
	Description string   `json:"description,omitempty"`
	LeadID      string   `json:"leadId,omitempty"`
	StartDate   string   `json:"startDate,omitempty"`
	TargetDate  string   `json:"targetDate,omitempty"`
}

// ProjectUpdateInput represents the input for updating a project.
type ProjectUpdateInput struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	LeadID      string   `json:"leadId,omitempty"`
	StartDate   string   `json:"startDate,omitempty"`
	TargetDate  string   `json:"targetDate,omitempty"`
	TeamIDs     []string `json:"teamIds,omitempty"`
}

// ProjectMilestoneCreateInput represents the input for creating a project milestone.
type ProjectMilestoneCreateInput struct {
	Name        string `json:"name"`
	ProjectID   string `json:"projectId"`
	Description string `json:"description,omitempty"`
	TargetDate  string `json:"targetDate,omitempty"`
}

// ProjectMilestoneUpdateInput represents the input for updating a project milestone.
type ProjectMilestoneUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TargetDate  string `json:"targetDate,omitempty"`
}

// InitiativeCreateInput represents the input for creating an initiative.
type InitiativeCreateInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// InitiativeUpdateInput represents the input for updating an initiative.
type InitiativeUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message string `json:"message"`
}
