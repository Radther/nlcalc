package cleaner

import "testing"

func TestClean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_expression",
			input:    "10 + 15",
			expected: "10 + 15",
		},
		{
			name:     "remove_articles_and_words",
			input:    "the sum of 10 + 15",
			expected: "10 + 15",
		},
		{
			name:     "remove_filler_words",
			input:    "what is 10 + 15",
			expected: "10 + 15",
		},
		{
			name:     "complex_percentage",
			input:    "20 * 0.01 * 100",
			expected: "20 * 0.01 * 100",
		},
		{
			name:     "preserve_decimal_numbers",
			input:    "5 * 3.14",
			expected: "5 * 3.14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clean(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}