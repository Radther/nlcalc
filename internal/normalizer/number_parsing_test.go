package normalizer

import (
	"strings"
	"testing"
)

func TestParseNumberPhrase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		valid    bool
	}{
		// Valid compound numbers
		{
			name:     "one_hundred_and_twenty_seven",
			input:    "one hundred and twenty seven",
			expected: 127,
			valid:    true,
		},
		{
			name:     "simpler_billion_test",
			input:    "two billion",
			expected: 2000000000,
			valid:    true,
		},
		{
			name:     "billion_with_thousand",
			input:    "two billion three thousand",
			expected: 2000003000,
			valid:    true,
		},
		{
			name:     "simple_hundred",
			input:    "one hundred",
			expected: 100,
			valid:    true,
		},
		{
			name:     "hundred_and_one",
			input:    "hundred and one",
			expected: 101,
			valid:    true,
		},
		{
			name:     "thousand_and_fifty",
			input:    "thousand and fifty",
			expected: 1050,
			valid:    true,
		},
		{
			name:     "million_and_five",
			input:    "one million and five",
			expected: 1000005,
			valid:    true,
		},
		{
			name:     "complex_billion",
			input:    "five billion three hundred million and twenty",
			expected: 5300000020,
			valid:    true,
		},

		// Invalid "and" usage (should be rejected)
		{
			name:     "fifteen_and_ten",
			input:    "fifteen and ten",
			expected: 0,
			valid:    false,
		},
		{
			name:     "five_and_seven",
			input:    "five and seven",
			expected: 0,
			valid:    false,
		},
		{
			name:     "twenty_and_thirty",
			input:    "twenty and thirty",
			expected: 0,
			valid:    false,
		},

		// Edge cases
		{
			name:     "just_and",
			input:    "and",
			expected: 0,
			valid:    false,
		},
		{
			name:     "and_at_beginning",
			input:    "and fifty",
			expected: 0,
			valid:    false,
		},
		{
			name:     "single_number",
			input:    "twenty",
			expected: 20,
			valid:    true,
		},
		{
			name:     "compound_tens",
			input:    "twenty seven",
			expected: 27,
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := parseNumberWords(strings.Fields(strings.ToLower(tt.input)))
			if valid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.valid, valid)
			}
			if valid && result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestNormalizeWithFixedAndLogic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid_compound_number",
			input:    "one hundred and twenty seven",
			expected: "127",
		},
		{
			name:     "invalid_and_stays_separate",
			input:    "fifteen and ten",
			expected: "15 and 10",
		},
		{
			name:     "complex_expression_with_valid_and",
			input:    "the sum of one hundred and twenty plus fifteen",
			expected: "the sum of 120 + 15",
		},
		{
			name:     "complex_expression_with_invalid_and",
			input:    "the sum of ten and fifteen",
			expected: "the sum of 10 and 15",
		},
		{
			name:     "mixed_numbers_and_operations",
			input:    "two hundred and fifty plus ten and twenty",
			expected: "250 + 10 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, nil, '.')
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}