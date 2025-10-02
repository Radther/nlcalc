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
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestParseWithVariables(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		variables map[string]float64
		expected  float64
	}{
		{
			name:      "simple_variable",
			input:     "x plus ten",
			variables: map[string]float64{"x": 5},
			expected:  15.0,
		},
		{
			name:      "multiple_variables",
			input:     "price plus tax",
			variables: map[string]float64{"price": 100, "tax": 15},
			expected:  115.0,
		},
		{
			name:      "variable_with_order_of_operations",
			input:     "price plus price times tax",
			variables: map[string]float64{"price": 100, "tax": 0.15},
			expected:  115.0,
		},
		{
			name:      "overlapping_variable_names",
			input:     "x times xx",
			variables: map[string]float64{"x": 2, "xx": 5},
			expected:  10.0,
		},
		{
			name:      "variable_with_parentheses",
			input:     "(base plus bonus) times rate",
			variables: map[string]float64{"base": 100, "bonus": 20, "rate": 1.5},
			expected:  180.0,
		},
		{
			name:      "variable_with_percentage",
			input:     "discount% of price",
			variables: map[string]float64{"discount": 20, "price": 100},
			expected:  20.0,
		},
		{
			name:      "case_insensitive_variables",
			input:     "Price PLUS Tax",
			variables: map[string]float64{"price": 100, "tax": 25},
			expected:  125.0,
		},
		{
			name:      "complex_expression_with_variables",
			input:     "quantity times price plus shipping",
			variables: map[string]float64{"quantity": 3, "price": 25.5, "shipping": 10},
			expected:  86.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, tt.variables)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestParseNegativeNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "negative_number_symbol",
			input:    "-2000",
			expected: -2000.0,
		},
		{
			name:     "negative_written_twice",
			input:    "minus twenty minus 10",
			expected: -30.0,
		},
		{
			name:     "subtraction",
			input:    "10 - 5",
			expected: 5.0,
		},
		{
			name:     "negative_plus_positive",
			input:    "-10 + 5",
			expected: -5.0,
		},
		{
			name:     "positive_plus_negative",
			input:    "10 + -5",
			expected: 5.0,
		},
		{
			name:     "negative_parentheses",
			input:    "-(10 + 5)",
			expected: -15.0,
		},
		{
			name:     "addition_with_negative_parentheses",
			input:    "10 + -(2 + 4)",
			expected: 4.0,
		},
		{
			name:     "negative_written_number",
			input:    "minus ten",
			expected: -10.0,
		},
		{
			name:     "multiplication_with_negative",
			input:    "10 * -5",
			expected: -50.0,
		},
		{
			name:     "double_negative",
			input:    "-(-5)",
			expected: 5.0,
		},
		{
			name:     "complex_negative_expression",
			input:    "-10 + 5 * -2",
			expected: -20.0,
		},
		{
			name:     "negative_with_parentheses_and_precedence",
			input:    "-(2 + 3) * 4",
			expected: -20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}