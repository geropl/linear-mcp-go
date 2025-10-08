package tools

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	PriorityNone   = 0
	PriorityUrgent = 1
	PriorityHigh   = 2
	PriorityMedium = 3
	PriorityLow    = 4
)

var priorityNames = map[int]string{
	PriorityNone:   "No priority",
	PriorityUrgent: "Urgent",
	PriorityHigh:   "High",
	PriorityMedium: "Medium",
	PriorityLow:    "Low",
}

var priorityFromName = map[string]int{
	"no priority": PriorityNone,
	"none":        PriorityNone,
	"urgent":      PriorityUrgent,
	"high":        PriorityHigh,
	"medium":      PriorityMedium,
	"low":         PriorityLow,
}

// priorityToString converts numeric priority to textual representation
func priorityToString(priority int) string {
	if name, ok := priorityNames[priority]; ok {
		return name
	}
	return "Unknown"
}

// parsePriority accepts both numeric (0-4) and textual representations
// Returns the numeric value and an error if invalid
func parsePriority(input string) (int, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Try parsing as number first
	if num, err := strconv.Atoi(input); err == nil {
		if num >= PriorityNone && num <= PriorityLow {
			return num, nil
		}
		return 0, fmt.Errorf("priority number must be between 0 and 4, got %d", num)
	}

	// Try parsing as text
	if priority, ok := priorityFromName[input]; ok {
		return priority, nil
	}

	return 0, fmt.Errorf("invalid priority: %s (valid values: 0-4, no priority, urgent, high, medium, low)", input)
}

// getPriorityOptions returns the property options for priority parameters
func getPriorityOptions() []mcp.PropertyOption {
	return []mcp.PropertyOption{
		mcp.Description("Priority"),
		mcp.Enum("no priority", "urgent", "high", "medium", "low"),
	}
}
