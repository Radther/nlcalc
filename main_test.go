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
		verbose  bool
		expected string
	}{
		{
			name:     "written_numbers",
			input:    "ten plus fifteen",
			verbose:  false,
			expected: "25",
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			verbose:  false,
			expected: "15",
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			verbose:  false,
			expected: "20",
		},
		{
			name:     "written_numbers_verbose",
			input:    "ten plus fifteen",
			verbose:  true,
			expected: "Original: ten plus fifteen\nNormalized: 10 + 15\nCleaned: 10 + 15\nTokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]\nResult: 25",
		},
		{
			name:     "numeric_expression_verbose",
			input:    "5 + 10",
			verbose:  true,
			expected: "Original: 5 + 10\nNormalized: 5 + 10\nCleaned: 5 + 10\nTokens: [NUMBER(5), OPERATOR(+), NUMBER(10)]\nResult: 15",
		},
		{
			name:     "valid_compound_number",
			input:    "one hundred and twenty seven",
			verbose:  false,
			expected: "127",
		},
		{
			name:     "complex_billion",
			input:    "five billion three hundred million and twenty",
			verbose:  false,
			expected: "5.30000002e+09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{}
			if tt.verbose {
				args = append(args, "--verbose")
			}
			args = append(args, tt.input)

			cmd := exec.Command("./nlcalc_test", args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
			}

			result := strings.TrimSpace(string(output))
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
