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
			expected: "ten plus fifteen",
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			expected: "5 + 10",
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			expected: "20% of 100",
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