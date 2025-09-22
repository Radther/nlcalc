package normalizer

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "written_numbers",
			input:    "ten plus fifteen",
			expected: "10 + 15",
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			expected: "5 + 10",
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			expected: "20 * 0.01 * 100",
		},
		{
			name:     "complex_expression",
			input:    "Twenty percent of one hundred",
			expected: "20 * 0.01 * 100",
		},
		{
			name:     "complex_expression_percentage",
			input:    "Twenty % of one hundred",
			expected: "20 * 0.01 * 100",
		},
		{
			name:     "word_operations",
			input:    "five times three",
			expected: "5 * 3",
		},
		{
			name:     "compound_numbers",
			input:    "two hundred plus fifty",
			expected: "200 + 50",
		},
		{
			name:     "complex_compound_numbers",
			input:    "two hundred and fifty",
			expected: "250",
		},
		{
			name:     "very_complex_numbers",
			input:    "one hundred and twenty seven plus three hundred forty five",
			expected: "127 + 345",
		},
		{
			name:     "thousand_numbers",
			input:    "two thousand five hundred",
			expected: "2500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
