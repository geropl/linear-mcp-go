package tools

import (
	"testing"
)

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		// Numeric inputs
		{"zero", "0", 0, false},
		{"one", "1", 1, false},
		{"two", "2", 2, false},
		{"three", "3", 3, false},
		{"four", "4", 4, false},
		{"invalid number", "5", 0, true},
		{"negative", "-1", 0, true},

		// Textual inputs (lowercase)
		{"no priority", "no priority", 0, false},
		{"none", "none", 0, false},
		{"urgent", "urgent", 1, false},
		{"high", "high", 2, false},
		{"medium", "medium", 3, false},
		{"low", "low", 4, false},

		// Textual inputs (mixed case)
		{"Urgent", "Urgent", 1, false},
		{"HIGH", "HIGH", 2, false},
		{"MeDiUm", "MeDiUm", 3, false},

		// Whitespace handling
		{"with spaces", "  urgent  ", 1, false},
		{"with tabs", "\thigh\t", 2, false},

		// Invalid inputs
		{"invalid text", "super-urgent", 0, true},
		{"empty", "", 0, true},
		{"random", "xyz", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePriority(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePriority(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPriorityToString(t *testing.T) {
	tests := []struct {
		name     string
		priority int
		want     string
	}{
		{"zero", 0, "No priority"},
		{"urgent", 1, "Urgent"},
		{"high", 2, "High"},
		{"medium", 3, "Medium"},
		{"low", 4, "Low"},
		{"invalid", 5, "Unknown"},
		{"negative", -1, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := priorityToString(tt.priority); got != tt.want {
				t.Errorf("priorityToString(%d) = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}
