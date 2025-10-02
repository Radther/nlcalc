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

func TestParseUnaryPlus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "unary_plus_simple",
			input:    "+5",
			expected: 5.0,
		},
		{
			name:     "unary_plus_preserves_negative",
			input:    "+(100 - 200)",
			expected: -100.0,
		},
		{
			name:     "unary_plus_with_addition",
			input:    "+10 + 5",
			expected: 15.0,
		},
		{
			name:     "binary_plus_with_unary_plus",
			input:    "10 + +5",
			expected: 15.0,
		},
		{
			name:     "unary_plus_natural_language",
			input:    "plus five",
			expected: 5.0,
		},
		{
			name:     "unary_plus_then_subtraction",
			input:    "plus ten minus five",
			expected: 5.0,
		},
		{
			name:     "unary_plus_with_multiplication",
			input:    "+5 * 3",
			expected: 15.0,
		},
		{
			name:     "mixed_unary_plus_and_minus",
			input:    "+10 - 5",
			expected: 5.0,
		},
		{
			name:     "unary_plus_after_unary_minus",
			input:    "-+5",
			expected: -5.0,
		},
		{
			name:     "double_unary_plus",
			input:    "+(+5)",
			expected: 5.0,
		},
		{
			name:     "unary_plus_complex_expression",
			input:    "+(10 + 5) * 2",
			expected: 30.0,
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

func TestParsePowerOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Symbol
		{
			name:     "power_symbol",
			input:    "2 ^ 3",
			expected: 8.0,
		},
		{
			name:     "power_symbol_chained",
			input:    "2 ^ 3 ^ 2",
			expected: 512.0, // 2^(3^2) = 2^9 = 512 (right-associative)
		},

		// Natural language - basic
		{
			name:     "power_keyword",
			input:    "2 power 3",
			expected: 8.0,
		},
		{
			name:     "raised_to",
			input:    "2 raised to 3",
			expected: 8.0,
		},
		{
			name:     "raised_to_the_power_of",
			input:    "2 raised to the power of 3",
			expected: 8.0,
		},
		{
			name:     "to_the_power_of",
			input:    "2 to the power of 3",
			expected: 8.0,
		},

		// Shortcuts
		{
			name:     "squared",
			input:    "5 squared",
			expected: 25.0,
		},
		{
			name:     "cubed",
			input:    "3 cubed",
			expected: 27.0,
		},
		{
			name:     "written_number_squared",
			input:    "ten squared",
			expected: 100.0,
		},

		// Precedence
		{
			name:     "power_before_multiply",
			input:    "2 * 3 ^ 2",
			expected: 18.0, // 2 * (3^2) = 2 * 9 = 18
		},
		{
			name:     "power_before_add",
			input:    "2 + 3 ^ 2",
			expected: 11.0, // 2 + (3^2) = 2 + 9 = 11
		},
		{
			name:     "parentheses_override",
			input:    "(2 + 3) ^ 2",
			expected: 25.0, // (2+3)^2 = 5^2 = 25
		},

		// Right-associativity
		{
			name:     "right_associative",
			input:    "2 ^ 3 ^ 2",
			expected: 512.0, // 2^(3^2) = 2^9 = 512
		},

		// Complex
		{
			name:     "complex_power_expression",
			input:    "10 power 2 plus 5 times 3",
			expected: 115.0, // (10^2) + (5*3) = 100 + 15 = 115
		},
		{
			name:     "negative_base",
			input:    "-2 ^ 3",
			expected: -8.0, // (-2)^3 = -8
		},
		{
			name:     "decimal_exponent",
			input:    "4 ^ 0.5",
			expected: 2.0, // 4^0.5 = sqrt(4) = 2
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