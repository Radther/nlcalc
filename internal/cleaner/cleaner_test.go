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
			name:     "with_articles",
			input:    "the sum of ten and fifteen",
			expected: "the sum of ten and fifteen",
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