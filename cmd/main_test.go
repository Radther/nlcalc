package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "nlcalc_test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("nlcalc_test")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "written_numbers",
			input:    "ten plus fifteen",
			expected: "Original: ten plus fifteen\nNormalized: 10 + 15\nCleaned: 10 + 15\nTokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]\nResult: 25",
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			expected: "Original: 5 + 10\nNormalized: 5 + 10\nCleaned: 5 + 10\nTokens: [NUMBER(5), OPERATOR(+), NUMBER(10)]\nResult: 15",
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			expected: "Original: 20% of 100\nNormalized: 20 * 0.01 * 100\nCleaned: 20 * 0.01 * 100\nTokens: [NUMBER(20), OPERATOR(*), NUMBER(0.01), OPERATOR(*), NUMBER(100)]\nResult: 20",
		},
		{
			name:     "invalid_and_usage",
			input:    "the sum of fifteen and ten",
			expected: "Original: the sum of fifteen and ten\nNormalized: the sum of 15 and 10\nCleaned: 15 10\nTokenization error: consecutive numbers not allowed: 15 10",
		},
		{
			name:     "valid_compound_number",
			input:    "one hundred and twenty seven",
			expected: "Original: one hundred and twenty seven\nNormalized: 127\nCleaned: 127\nTokens: [NUMBER(127)]\nResult: 127",
		},
		{
			name:     "complex_billion",
			input:    "five billion three hundred million and twenty",
			expected: "Original: five billion three hundred million and twenty\nNormalized: 5300000020\nCleaned: 5300000020\nTokens: [NUMBER(5300000020)]\nResult: 5.30000002e+09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./nlcalc_test", tt.input)
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
