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
			result := Normalize(tt.input, nil)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNormalizeWithVariables(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		variables map[string]float64
		expected  string
	}{
		{
			name:      "single_variable",
			input:     "x plus five",
			variables: map[string]float64{"x": 10},
			expected:  "10 + 5",
		},
		{
			name:      "multiple_variables",
			input:     "price plus tax",
			variables: map[string]float64{"price": 100, "tax": 15},
			expected:  "100 + 15",
		},
		{
			name:      "overlapping_variable_names",
			input:     "x plus xx",
			variables: map[string]float64{"x": 5, "xx": 10},
			expected:  "5 + 10",
		},
		{
			name:      "variable_with_operations",
			input:     "price times quantity",
			variables: map[string]float64{"price": 25.5, "quantity": 3},
			expected:  "25.5 * 3",
		},
		{
			name:      "case_insensitive_variables",
			input:     "Price PLUS Tax",
			variables: map[string]float64{"price": 100, "tax": 20},
			expected:  "100 + 20",
		},
		{
			name:      "variable_not_partial_match",
			input:     "price plus pricey",
			variables: map[string]float64{"price": 100},
			expected:  "100 + pricey",
		},
		{
			name:      "empty_variables_map",
			input:     "ten plus five",
			variables: map[string]float64{},
			expected:  "10 + 5",
		},
		{
			name:      "variable_with_decimal",
			input:     "rate times hundred",
			variables: map[string]float64{"rate": 0.15},
			expected:  "0.15 * 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, tt.variables)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
