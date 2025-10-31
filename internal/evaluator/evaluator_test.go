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
		// Cube root tests
		{
			name: "cbrt_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "cbrt"},
				{Type: tokenizer.NUMBER, Value: "27"},
			},
			expected: 3.0, // cbrt(27) = 3
			hasError: false,
		},
		{
			name: "cbrt_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "cbrt"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "8"},
			},
			expected: -2.0, // cbrt(-8) = -2
			hasError: false,
		},
		// Inverse trigonometric tests
		{
			name: "asin_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asin"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 0.5235987755982989, // asin(0.5) = π/6 ≈ 0.524
			hasError: false,
		},
		{
			name: "asin_domain_error_high",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asin"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 0.0,
			hasError: true, // asin(2) is out of domain
		},
		{
			name: "asin_domain_error_low",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asin"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 0.0,
			hasError: true, // asin(-2) is out of domain
		},
		{
			name: "acos_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acos"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 1.0471975511965979, // acos(0.5) = π/3 ≈ 1.047
			hasError: false,
		},
		{
			name: "acos_domain_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acos"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 0.0,
			hasError: true, // acos(2) is out of domain
		},
		{
			name: "atan_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "atan"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.7853981633974483, // atan(1) = π/4 ≈ 0.785
			hasError: false,
		},
		// Hyperbolic function tests
		{
			name: "sinh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sinh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 1.1752011936438014, // sinh(1) ≈ 1.175
			hasError: false,
		},
		{
			name: "cosh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "cosh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 1.5430806348152437, // cosh(1) ≈ 1.543
			hasError: false,
		},
		{
			name: "tanh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "tanh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.7615941559557649, // tanh(1) ≈ 0.762
			hasError: false,
		},
		// Rounding function tests
		{
			name: "floor_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "floor"},
				{Type: tokenizer.NUMBER, Value: "3.7"},
			},
			expected: 3.0, // floor(3.7) = 3
			hasError: false,
		},
		{
			name: "floor_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "floor"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "3.7"},
			},
			expected: -4.0, // floor(-3.7) = -4
			hasError: false,
		},
		{
			name: "ceil_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "ceil"},
				{Type: tokenizer.NUMBER, Value: "3.2"},
			},
			expected: 4.0, // ceil(3.2) = 4
			hasError: false,
		},
		{
			name: "ceil_negative",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "ceil"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "3.2"},
			},
			expected: -3.0, // ceil(-3.2) = -3
			hasError: false,
		},
		{
			name: "round_up",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "round"},
				{Type: tokenizer.NUMBER, Value: "3.7"},
			},
			expected: 4.0, // round(3.7) = 4
			hasError: false,
		},
		{
			name: "round_down",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "round"},
				{Type: tokenizer.NUMBER, Value: "3.2"},
			},
			expected: 3.0, // round(3.2) = 3
			hasError: false,
		},
		{
			name: "round_half",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "round"},
				{Type: tokenizer.NUMBER, Value: "3.5"},
			},
			expected: 4.0, // round(3.5) = 4 (rounds to even)
			hasError: false,
		},
		// Combined new function tests
		{
			name: "cbrt_in_expression",
			tokens: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.FUNCTION, Value: "cbrt"},
				{Type: tokenizer.NUMBER, Value: "8"},
			},
			expected: 4.0, // 2 * cbrt(8) = 2 * 2 = 4
			hasError: false,
		},
		{
			name: "floor_ceil_nested",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "floor"},
				{Type: tokenizer.FUNCTION, Value: "ceil"},
				{Type: tokenizer.NUMBER, Value: "3.2"},
			},
			expected: 4.0, // floor(ceil(3.2)) = floor(4) = 4
			hasError: false,
		},
		// Exponential and logarithmic tests
		{
			name: "exp_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "exp"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 2.718281828459045, // exp(1) = e
			hasError: false,
		},
		{
			name: "exp_zero",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "exp"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 1.0, // exp(0) = 1
			hasError: false,
		},
		{
			name: "expm1_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "expm1"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0, // expm1(0) = 0
			hasError: false,
		},
		{
			name: "log2_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log2"},
				{Type: tokenizer.NUMBER, Value: "8"},
			},
			expected: 3.0, // log2(8) = 3
			hasError: false,
		},
		{
			name: "log2_negative_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log2"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.0,
			hasError: true, // log2(-1) is error
		},
		{
			name: "log1p_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log1p"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0, // log1p(0) = log(1) = 0
			hasError: false,
		},
		{
			name: "log1p_domain_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log1p"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 0.0,
			hasError: true, // log1p(-2) is error (must be > -1)
		},
		// Extended trigonometric tests
		{
			name: "sec_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sec"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 1.0, // sec(0) = 1/cos(0) = 1
			hasError: false,
		},
		{
			name: "csc_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "csc"},
				{Type: tokenizer.NUMBER, Value: "1.5707963267948966"}, // π/2
			},
			expected: 1.0, // csc(π/2) = 1/sin(π/2) = 1
			hasError: false,
		},
		{
			name: "cot_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "cot"},
				{Type: tokenizer.NUMBER, Value: "0.7853981633974483"}, // π/4
			},
			expected: 1.0, // cot(π/4) = 1/tan(π/4) = 1
			hasError: false,
		},
		// Inverse extended trigonometric tests
		{
			name: "asec_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asec"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 1.0471975511965976, // asec(2) = acos(1/2) = π/3
			hasError: false,
		},
		{
			name: "asec_domain_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asec"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 0.0,
			hasError: true, // asec(0.5) is error (must be >= 1 or <= -1)
		},
		{
			name: "acsc_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acsc"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 0.5235987755982989, // acsc(2) = asin(1/2) = π/6
			hasError: false,
		},
		{
			name: "acsc_domain_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acsc"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 0.0,
			hasError: true, // acsc(0.5) is error
		},
		{
			name: "acot_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acot"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.7853981633974483, // acot(1) = π/4
			hasError: false,
		},
		{
			name: "acot_zero",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acot"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 1.5707963267948966, // acot(0) = π/2
			hasError: false,
		},
		// Inverse hyperbolic tests
		{
			name: "asinh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asinh"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0, // asinh(0) = 0
			hasError: false,
		},
		{
			name: "asinh_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "asinh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.881373587019543, // asinh(1) ≈ 0.881
			hasError: false,
		},
		{
			name: "acosh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acosh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.0, // acosh(1) = 0
			hasError: false,
		},
		{
			name: "acosh_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acosh"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: 1.3169578969248166, // acosh(2) ≈ 1.317
			hasError: false,
		},
		{
			name: "acosh_domain_error",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "acosh"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 0.0,
			hasError: true, // acosh(0.5) is error (must be >= 1)
		},
		{
			name: "atanh_simple",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "atanh"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0, // atanh(0) = 0
			hasError: false,
		},
		{
			name: "atanh_positive",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "atanh"},
				{Type: tokenizer.NUMBER, Value: "0.5"},
			},
			expected: 0.5493061443340548, // atanh(0.5) ≈ 0.549
			hasError: false,
		},
		{
			name: "atanh_domain_error_high",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "atanh"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.0,
			hasError: true, // atanh(1) is error (must be in (-1, 1))
		},
		{
			name: "atanh_domain_error_low",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "atanh"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "1"},
			},
			expected: 0.0,
			hasError: true, // atanh(-1) is error
		},
		// Combined new function tests
		{
			name: "exp_log2_combined",
			tokens: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "log2"},
				{Type: tokenizer.FUNCTION, Value: "exp"},
				{Type: tokenizer.NUMBER, Value: "0"},
			},
			expected: 0.0, // log2(exp(0)) = log2(1) = 0
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