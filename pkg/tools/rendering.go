package tools

import (
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
)

// Full Entity Rendering Functions

// formatIssue returns a consistently formatted full representation of an issue
func formatIssue(issue *linear.Issue) string {
	if issue == nil {
		return "Issue: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Issue: %s (UUID: %s)\n", issue.Identifier, issue.ID))
	result.WriteString(fmt.Sprintf("Title: %s\n", issue.Title))
	result.WriteString(fmt.Sprintf("URL: %s\n", issue.URL))
	
	priorityStr := "None"
	if issue.Priority > 0 {
		priorityStr = fmt.Sprintf("%d", issue.Priority)
	}
	result.WriteString(fmt.Sprintf("Priority: %s\n", priorityStr))
	
	statusStr := "None"
	if issue.Status != "" {
		statusStr = issue.Status
	} else if issue.State != nil {
		statusStr = issue.State.Name
	}
	result.WriteString(fmt.Sprintf("Status: %s\n", statusStr))
	
	// Include labels if available
	if issue.Labels != nil && len(issue.Labels.Nodes) > 0 {
		labelNames := make([]string, 0, len(issue.Labels.Nodes))
		for _, label := range issue.Labels.Nodes {
			labelNames = append(labelNames, label.Name)
		}
		result.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(labelNames, ", ")))
	} else {
		result.WriteString("Labels: None\n")
	}
	
	// Include description
	if issue.Description != "" {
		result.WriteString(fmt.Sprintf("Description: %s\n", issue.Description))
	} else {
		result.WriteString("Description: None\n")
	}
	
	return result.String()
}

// formatTeam returns a consistently formatted full representation of a team
func formatTeam(team *linear.Team) string {
	if team == nil {
		return "Team: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Team: %s (UUID: %s)\n", team.Name, team.ID))
	result.WriteString(fmt.Sprintf("Key: %s\n", team.Key))
	
	return result.String()
}

// formatUser returns a consistently formatted full representation of a user
func formatUser(user *linear.User) string {
	if user == nil {
		return "User: Unknown"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("User: %s (UUID: %s)\n", user.Name, user.ID))
	
	if user.Email != "" {
		result.WriteString(fmt.Sprintf("Email: %s\n", user.Email))
	}
	
	return result.String()
}

// formatComment returns a consistently formatted full representation of a comment
func formatComment(comment *linear.Comment) string {
	if comment == nil {
		return "Comment: Unknown"
	}
	
	userName := "Unknown"
	if comment.User != nil {
		userName = comment.User.Name
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Comment by %s (UUID: %s)\n", userName, comment.ID))
	result.WriteString(fmt.Sprintf("Body: %s\n", comment.Body))
	
	if !comment.CreatedAt.IsZero() {
		result.WriteString(fmt.Sprintf("Created: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	
	return result.String()
}

// Entity Identifier Rendering Functions

// formatIssueIdentifier returns a consistently formatted identifier for an issue
func formatIssueIdentifier(issue *linear.Issue) string {
	if issue == nil {
		return "Issue: Unknown"
	}
	return fmt.Sprintf("Issue: %s (UUID: %s)", issue.Identifier, issue.ID)
}

// formatTeamIdentifier returns a consistently formatted identifier for a team
func formatTeamIdentifier(team *linear.Team) string {
	if team == nil {
		return "Team: Unknown"
	}
	return fmt.Sprintf("Team: %s (UUID: %s)", team.Name, team.ID)
}

// formatUserIdentifier returns a consistently formatted identifier for a user
func formatUserIdentifier(user *linear.User) string {
	if user == nil {
		return "User: Unknown"
	}
	return fmt.Sprintf("User: %s (UUID: %s)", user.Name, user.ID)
}

// formatCommentIdentifier returns a consistently formatted identifier for a comment
func formatCommentIdentifier(comment *linear.Comment) string {
	if comment == nil {
		return "Unknown"
	}
	
	// Extract hash from URL fragment (e.g., "#comment-ae3d62d6" -> "ae3d62d6")
	shortHash := comment.ID[:8] // fallback to first 8 chars of UUID
	if comment.URL != "" {
		// Look for "#comment-" in the URL
		if idx := strings.Index(comment.URL, "#comment-"); idx != -1 {
			hashPart := comment.URL[idx+9:] // Skip "#comment-"
			if hashPart != "" {
				shortHash = hashPart
			}
		}
	}
	
	return fmt.Sprintf("Comment: %s (UUID: %s)", shortHash, comment.ID)
}
