package nlcalc

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "written_numbers",
			input:    "ten plus fifteen",
			expected: 25.0,
		},
		{
			name:     "numeric_expression",
			input:    "5 + 10",
			expected: 15.0,
		},
		{
			name:     "percentage",
			input:    "20% of 100",
			expected: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}