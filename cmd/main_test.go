package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "mathparser_test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("mathparser_test")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "written_numbers",
			input:    "ten plus fifteen",
			expected: "Original: ten plus fifteen\nNormalized: 10 + 15\nResult: not implemented yet",
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			expected: "Original: 5 + 10\nNormalized: 5 + 10\nResult: not implemented yet",
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			expected: "Original: 20% of 100\nNormalized: 20 * 0.01 * 100\nResult: not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./mathparser_test", tt.input)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			result := strings.TrimSpace(string(output))
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
