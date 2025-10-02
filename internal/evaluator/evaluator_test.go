package evaluator

import (
	"math"
	"testing"

	"github.com/radther/nlcalc/internal/tokenizer"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []tokenizer.Token
		expected float64
		hasError bool
	}{
		{
			name:     "single_number",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "42"}},
			expected: 42.0,
			hasError: false,
		},
		{
			name:     "simple_addition",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "10"}, {Type: tokenizer.OPERATOR, Value: "+"}, {Type: tokenizer.NUMBER, Value: "15"}},
			expected: 25.0,
			hasError: false,
		},
		{
			name:     "simple_subtraction",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "20"}, {Type: tokenizer.OPERATOR, Value: "-"}, {Type: tokenizer.NUMBER, Value: "5"}},
			expected: 15.0,
			hasError: false,
		},
		{
			name:     "simple_multiplication",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "6"}, {Type: tokenizer.OPERATOR, Value: "*"}, {Type: tokenizer.NUMBER, Value: "7"}},
			expected: 42.0,
			hasError: false,
		},
		{
			name:     "simple_division",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "84"}, {Type: tokenizer.OPERATOR, Value: "/"}, {Type: tokenizer.NUMBER, Value: "2"}},
			expected: 42.0,
			hasError: false,
		},
		{
			name:     "decimal_numbers",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "5"}, {Type: tokenizer.OPERATOR, Value: "*"}, {Type: tokenizer.NUMBER, Value: "3.14"}},
			expected: 15.7,
			hasError: false,
		},
		{
			name: "order_of_operations_multiply_first",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "4"},
			},
			expected: 50.0,
			hasError: false,
		},
		{
			name: "order_of_operations_divide_first",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "100"},
				{Type: tokenizer.OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "80"},
				{Type: tokenizer.OPERATOR, Value: "/"},
				{Type: tokenizer.NUMBER, Value: "4"},
			},
			expected: 80.0,
			hasError: false,
		},
		{
			name: "parentheses_override_precedence",
			tokens: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "4"},
			},
			expected: 80.0,
			hasError: false,
		},
		{
			name: "percentage_calculation",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "20"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "0.01"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "100"},
			},
			expected: 20.0,
			hasError: false,
		},
		{
			name: "complex_expression",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "3"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "4"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 13.0,
			hasError: false,
		},
		{
			name:     "empty_tokens",
			tokens:   []tokenizer.Token{},
			expected: 0.0,
			hasError: true,
		},
		{
			name:     "division_by_zero",
			tokens:   []tokenizer.Token{{Type: tokenizer.NUMBER, Value: "10"}, {Type: tokenizer.OPERATOR, Value: "/"}, {Type: tokenizer.NUMBER, Value: "0"}},
			expected: 0.0,
			hasError: true,
		},
		{
			name:     "invalid_single_token",
			tokens:   []tokenizer.Token{{Type: tokenizer.OPERATOR, Value: "+"}},
			expected: 0.0,
			hasError: true,
		},
		{
			name: "mismatched_parentheses",
			tokens: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 0.0,
			hasError: true,
		},
		{
			name: "operator_at_start",
			tokens: []tokenizer.Token{
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "10"},
			},
			expected: 0.0,
			hasError: true,
		},
		{
			name: "operator_at_end",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
			},
			expected: 0.0,
			hasError: true,
		},
		{
			name: "negative_number_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "2000"},
			},
			expected: -2000.0,
			hasError: false,
		},
		{
			name: "negative_plus_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "20"},
				{Type: tokenizer.OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "10"},
			},
			expected: -30.0,
			hasError: false,
		},
		{
			name: "addition_with_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: -5.0,
			hasError: false,
		},
		{
			name: "positive_plus_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "negative_parentheses",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: -15.0,
			hasError: false,
		},
		{
			name: "addition_with_negative_parentheses",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "4"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: 4.0,
			hasError: false,
		},
		{
			name: "multiplication_with_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: -50.0,
			hasError: false,
		},
		{
			name: "double_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "unary_plus_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "unary_plus_preserves_negative_from_parentheses",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "100"},
				{Type: tokenizer.OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "200"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: -100.0,
			hasError: false,
		},
		{
			name: "unary_plus_after_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: -5.0,
			hasError: false,
		},
		{
			name: "mixed_unary_plus_and_minus_operators",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: -5.0,
			hasError: false,
		},
		{
			name: "unary_plus_with_multiplication",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "3"},
			},
			expected: 30.0,
			hasError: false,
		},
		{
			name: "double_unary_plus",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.UNARY_OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "power_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "3"},
			},
			expected: 8.0,
			hasError: false,
		},
		{
			name: "power_right_associative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "3"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 512.0, // 2^(3^2) = 2^9 = 512, NOT (2^3)^2 = 64
			hasError: false,
		},
		{
			name: "power_with_multiplication_precedence",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "3"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 18.0, // 2 * (3^2) = 2 * 9 = 18
			hasError: false,
		},
		{
			name: "power_with_parentheses",
			tokens: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "3"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 25.0, // (2+3)^2 = 5^2 = 25
			hasError: false,
		},
		{
			name: "negative_base_power",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "3"},
			},
			expected: -8.0, // (-2)^3 = -8
			hasError: false,
		},
		{
			name: "decimal_exponent",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "4"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 2.0, // 4^0.5 = sqrt(4) = 2
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.tokens)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}

			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.hasError && math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestEvaluateFunctions(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []tokenizer.Token
		expected float64
		hasError bool
	}{
		{
			name: "sqrt_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: 4.0,
			hasError: false,
		},
		{
			name: "sqrt_consecutive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: 2.0, // sqrt(sqrt(16)) = sqrt(4) = 2
			hasError: false,
		},
		{
			name: "abs_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "abs"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "abs_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "abs"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 5.0,
			hasError: false,
		},
		{
			name: "log_base10",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log"},
				{Type: tokenizer.NUMBER, Value: "100"},
			},
			expected: 2.0, // log10(100) = 2
			hasError: false,
		},
		{
			name: "ln_natural",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "ln"},
				{Type: tokenizer.NUMBER, Value: "2.718281828459045"},
			},
			expected: 1.0, // ln(e) ≈ 1
			hasError: false,
		},
		{
			name: "function_in_expression",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: 8.0, // 2 * sqrt(16) = 2 * 4 = 8
			hasError: false,
		},
		{
			name: "function_with_power",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
				{Type: tokenizer.OPERATOR, Value: "^"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 16.0, // (sqrt(16))^2 = 4^2 = 16
			hasError: false,
		},
		{
			name: "multiple_different_functions",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "10000"},
			},
			expected: 2.0, // log(sqrt(10000)) = log(100) = 2
			hasError: false,
		},
		{
			name: "negative_unary_with_function",
			tokens: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: -4.0, // -(sqrt(16)) = -4
			hasError: false,
		},
		{
			name: "sqrt_negative_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: 0.0,
			hasError: true, // sqrt(-16) is error
		},
		{
			name: "log_negative_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: 0.0,
			hasError: true, // log(-5) is error
		},
		{
			name: "log_zero_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0,
			hasError: true, // log(0) is error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.tokens)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}

			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.hasError && math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}