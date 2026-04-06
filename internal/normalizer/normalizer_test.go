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
			expected: "20 * 0.01 * of 100",
		},
		{
			name:     "complex_expression",
			input:    "Twenty percent of one hundred",
			expected: "20 * 0.01 * of 100",
		},
		{
			name:     "complex_expression_percentage",
			input:    "Twenty % of one hundred",
			expected: "20 * 0.01 * of 100",
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
		// Power operations
		{
			name:     "power_keyword",
			input:    "2 power 3",
			expected: "2 ^ 3",
		},
		{
			name:     "raised_to",
			input:    "2 raised to 3",
			expected: "2 ^ to 3", // "to" left for cleaner to remove
		},
		{
			name:     "raised_to_the_power_of",
			input:    "2 raised to the power of 3",
			expected: "2 ^ 3", // Full phrase replaced, no filler words left
		},
		{
			name:     "to_the_power_of",
			input:    "2 to the power of 3",
			expected: "2 to the ^ of 3", // "power" replaced, filler words left for cleaner
		},
		{
			name:     "squared",
			input:    "5 squared",
			expected: "5 ^ 2",
		},
		{
			name:     "cubed",
			input:    "3 cubed",
			expected: "3 ^ 3",
		},
		{
			name:     "written_number_squared",
			input:    "ten squared",
			expected: "10 ^ 2",
		},
		{
			name:     "written_number_cubed",
			input:    "five cubed",
			expected: "5 ^ 3",
		},
		{
			name:     "power_with_written_numbers",
			input:    "two power three",
			expected: "2 ^ 3",
		},
		{
			name:     "raised_with_written_numbers",
			input:    "two raised to three",
			expected: "2 ^ to 3", // "to" left for cleaner to remove
		},
		{
			name:     "power_symbol",
			input:    "2 ^ 3",
			expected: "2 ^ 3",
		},
		// Function operations
		{
			name:     "square_root",
			input:    "square root 16",
			expected: "sqrt 16",
		},
		{
			name:     "absolute_value",
			input:    "absolute -5",
			expected: "abs -5",
		},
		{
			name:     "logarithm",
			input:    "logarithm 100",
			expected: "log 100",
		},
		{
			name:     "natural_logarithm",
			input:    "natural logarithm 10",
			expected: "ln 10",
		},
		{
			name:     "sine",
			input:    "sine 90",
			expected: "sin 90",
		},
		{
			name:     "cosine",
			input:    "cosine 0",
			expected: "cos 0",
		},
		{
			name:     "tangent",
			input:    "tangent 45",
			expected: "tan 45",
		},
		{
			name:     "nested_functions",
			input:    "square root squareroot 16",
			expected: "sqrt sqrt 16",
		},
		{
			name:     "function_with_written_number",
			input:    "square root sixteen",
			expected: "sqrt 16",
		},
		{
			name:     "function_in_expression",
			input:    "two times square root 16",
			expected: "2 * sqrt 16",
		},
		// Function operations with extra
		{
			name:     "square_root",
			input:    "square root of 16",
			expected: "sqrt of 16",
		},
		{
			name:     "absolute_value",
			input:    "absolute value of -5",
			expected: "abs value of -5",
		},
		{
			name:     "logarithm",
			input:    "logarithm of 100",
			expected: "log of 100",
		},
		{
			name:     "natural_logarithm",
			input:    "natural logarithm of 10",
			expected: "ln of 10",
		},
		{
			name:     "sine",
			input:    "sine of 90",
			expected: "sin of 90",
		},
		{
			name:     "cosine",
			input:    "cosine of 0",
			expected: "cos of 0",
		},
		{
			name:     "tangent",
			input:    "tangent of 45",
			expected: "tan of 45",
		},
		{
			name:     "nested_functions",
			input:    "square root of square root of 16",
			expected: "sqrt of sqrt of 16",
		},
		{
			name:     "function_with_written_number",
			input:    "square root of sixteen",
			expected: "sqrt of 16",
		},
		{
			name:     "function_in_expression",
			input:    "two times square root of 16",
			expected: "2 * sqrt of 16",
		},
		// Thousands separator (default dot-decimal mode, comma is thousands sep)
		{
			name:     "thousands_separator_simple",
			input:    "10,000",
			expected: "10000",
		},
		{
			name:     "thousands_separator_large",
			input:    "1,000,000",
			expected: "1000000",
		},
		{
			name:     "thousands_separator_with_decimal",
			input:    "1,000.83",
			expected: "1000.83",
		},
		// "point" as decimal separator word
		{
			name:     "point_as_decimal",
			input:    "ten point two",
			expected: "10.2",
		},
		{
			name:     "point_as_decimal_in_expression",
			input:    "ten point two plus five",
			expected: "10.2 + 5",
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
			result := Normalize(tt.input, tt.variables, '.')
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNormalizeWithBuiltInConstants(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "pi_constant",
			input:    "pi times two",
			expected: "3.141592653589793 * 2",
		},
		{
			name:     "e_constant",
			input:    "e plus one",
			expected: "2.718281828459045 + 1",
		},
		{
			name:     "tau_constant",
			input:    "tau divided by two",
			expected: "6.283185307179586 / 2",
		},
		{
			name:     "phi_constant",
			input:    "phi times phi",
			expected: "1.618033988749895 * 1.618033988749895",
		},
		{
			name:     "multiple_constants",
			input:    "pi plus e",
			expected: "3.141592653589793 + 2.718281828459045",
		},
		{
			name:     "constant_with_written_numbers",
			input:    "two times pi",
			expected: "2 * 3.141592653589793",
		},
		{
			name:     "constant_case_insensitive",
			input:    "PI plus E",
			expected: "3.141592653589793 + 2.718281828459045",
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

func TestNormalizeUserVariablesOverrideBuiltInConstants(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		variables map[string]float64
		expected  string
	}{
		{
			name:      "user_pi_overrides_builtin",
			input:     "pi times two",
			variables: map[string]float64{"pi": 3.14},
			expected:  "3.14 * 2",
		},
		{
			name:      "user_e_overrides_builtin",
			input:     "e plus one",
			variables: map[string]float64{"e": 2.72},
			expected:  "2.72 + 1",
		},
		{
			name:      "partial_override",
			input:     "pi plus e",
			variables: map[string]float64{"pi": 3},
			expected:  "3 + 2.718281828459045",
		},
		{
			name:      "user_variable_and_builtin_constant",
			input:     "x times pi",
			variables: map[string]float64{"x": 5},
			expected:  "5 * 3.141592653589793",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, tt.variables, '.')
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNormalizeCommaDecimalDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "dot_as_thousands_separator",
			input:    "10.000",
			expected: "10000",
		},
		{
			name:     "dot_thousands_with_comma_decimal",
			input:    "1.000,83",
			expected: "1000.83",
		},
		{
			name:     "comma_decimal_simple",
			input:    "3,14",
			expected: "3.14",
		},
		{
			name:     "comma_decimal_expression",
			input:    "3,14 + 2,5",
			expected: "3.14 + 2.5",
		},
		{
			name:     "point_word_still_works",
			input:    "ten point two",
			expected: "10.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, nil, ',')
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
