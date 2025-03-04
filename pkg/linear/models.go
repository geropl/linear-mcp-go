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
	URL             string                   `json:"url"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	Labels          []Label                  `json:"labels,omitempty"`
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

// State represents a workflow state in Linear
type State struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// Comment represents a comment on a Linear issue
type Comment struct {
	ID        string            `json:"id"`
	Body      string            `json:"body"`
	User      *User             `json:"user,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	URL       string            `json:"url,omitempty"`
	Parent    *Comment          `json:"parent,omitempty"`
	Children  *CommentConnection `json:"children,omitempty"`
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
	ID         string  `json:"id"`
	Identifier string  `json:"identifier"`
	Title      string  `json:"title"`
	Priority   int     `json:"priority"`
	Status     string  `json:"status,omitempty"`
	StateName  string  `json:"stateName,omitempty"`
	URL        string  `json:"url"`
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
	Title       string  `json:"title"`
	TeamID      string  `json:"teamId"`
	Description string  `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	Status      string  `json:"status,omitempty"`
}

// UpdateIssueInput represents input for updating an issue
type UpdateIssueInput struct {
	ID          string  `json:"id"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	Status      string  `json:"status,omitempty"`
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

// AddCommentInput represents input for adding a comment
type AddCommentInput struct {
	IssueID        string `json:"issueId"`
	Body           string `json:"body"`
	CreateAsUser   string `json:"createAsUser,omitempty"`
	DisplayIconURL string `json:"displayIconUrl,omitempty"`
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
