package tools

import (
	"fmt"
	"strings"

	"github.com/geropl/linear-mcp-go/pkg/linear"
	"github.com/google/uuid"
)

// resolveIssueIdentifier resolves an issue identifier (UUID or "TEAM-123") to a UUID
func resolveIssueIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find an issue by identifier
	issue, err := linearClient.GetIssueByIdentifier(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to resolve issue identifier '%s': %v", identifier, err)
	}

	return issue.ID, nil
}

// resolveParentIssueIdentifier is an alias for resolveIssueIdentifier for backward compatibility
func resolveParentIssueIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	return resolveIssueIdentifier(linearClient, identifier)
}

// resolveUserIdentifier resolves a user identifier (UUID, name, or email) to a UUID
func resolveUserIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Otherwise, try to find a user by name or email
	// Get the organization to access all users
	org, err := linearClient.GetOrganization()
	if err != nil {
		return "", fmt.Errorf("failed to get organization: %v", err)
	}

	// First try exact match on name or email
	for _, user := range org.Users {
		if user.Name == identifier || user.Email == identifier {
			return user.ID, nil
		}
	}

	// If no exact match, try case-insensitive match
	identifierLower := strings.ToLower(identifier)
	for _, user := range org.Users {
		if strings.ToLower(user.Name) == identifierLower || strings.ToLower(user.Email) == identifierLower {
			return user.ID, nil
		}
	}

	return "", fmt.Errorf("no user found with identifier '%s'", identifier)
}

// resolveLabelIdentifiers resolves a list of label identifiers (UUIDs or names) to UUIDs
func resolveLabelIdentifiers(linearClient *linear.LinearClient, teamID string, labelIdentifiers []string) ([]string, error) {
	// Separate UUIDs and names
	var labelUUIDs []string
	var labelNames []string

	for _, identifier := range labelIdentifiers {
		if isValidUUID(identifier) {
			labelUUIDs = append(labelUUIDs, identifier)
		} else {
			labelNames = append(labelNames, identifier)
		}
	}

	// If there are no names to resolve, return the UUIDs directly
	if len(labelNames) == 0 {
		return labelUUIDs, nil
	}

	// Get labels by name
	labels, err := linearClient.GetLabelsByName(teamID, labelNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels by name: %v", err)
	}

	// Check if all label names were found
	if len(labels) < len(labelNames) {
		// Create a map of found label names for quick lookup
		foundLabels := make(map[string]bool)
		for _, label := range labels {
			foundLabels[label.Name] = true
		}

		// Find which label names were not found
		var missingLabels []string
		for _, name := range labelNames {
			if !foundLabels[name] {
				missingLabels = append(missingLabels, name)
			}
		}

		return nil, fmt.Errorf("label(s) not found: %s", strings.Join(missingLabels, ", "))
	}

	// Add the resolved label UUIDs to the result
	for _, label := range labels {
		labelUUIDs = append(labelUUIDs, label.ID)
	}

	return labelUUIDs, nil
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(uuidStr string) bool {
	return uuid.Validate(uuidStr) == nil
}

// resolveCommentIdentifier resolves a comment identifier (UUID or shorthand like "comment-53099b37") to a UUID
func resolveCommentIdentifier(linearClient *linear.LinearClient, identifier string) (string, error) {
	// If it's a valid UUID, use it directly
	if isValidUUID(identifier) {
		return identifier, nil
	}

	// Handle shorthand format like "comment-53099b37" or just "53099b37"
	var hash string
	if strings.HasPrefix(identifier, "comment-") {
		hash = strings.TrimPrefix(identifier, "comment-")
	} else {
		// Assume it's already just the hash part
		hash = identifier
	}

	// Try to get the comment by hash
	comment, err := linearClient.GetCommentByHash(hash)
	if err != nil {
		return "", fmt.Errorf("failed to resolve comment identifier '%s': %v", identifier, err)
	}

	return comment.ID, nil
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
