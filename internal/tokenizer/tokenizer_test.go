package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		hasError bool
	}{
		{
			name:  "simple_expression",
			input: "10 + 15",
			expected: []Token{
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "15"},
			},
			hasError: false,
		},
		{
			name:  "decimal_expression",
			input: "5 * 3.14",
			expected: []Token{
				{Type: NUMBER, Value: "5"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "3.14"},
			},
			hasError: false,
		},
		{
			name:  "complex_expression",
			input: "20 * 0.01 * 100",
			expected: []Token{
				{Type: NUMBER, Value: "20"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "0.01"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "100"},
			},
			hasError: false,
		},
		{
			name:  "expression_with_parentheses",
			input: "(10 + 15) * 2",
			expected: []Token{
				{Type: PARENTHESIS, Value: "("},
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "15"},
				{Type: PARENTHESIS, Value: ")"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "2"},
			},
			hasError: false,
		},
		{
			name:  "division_expression",
			input: "100 / 5",
			expected: []Token{
				{Type: NUMBER, Value: "100"},
				{Type: OPERATOR, Value: "/"},
				{Type: NUMBER, Value: "5"},
			},
			hasError: false,
		},
		{
			name:  "subtraction_expression",
			input: "50 - 25",
			expected: []Token{
				{Type: NUMBER, Value: "50"},
				{Type: OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "25"},
			},
			hasError: false,
		},
		{
			name:     "consecutive_operators_error",
			input:    "10 + + 5",
			expected: nil,
			hasError: true,
		},
		{
			name:     "incomplete_expression_error",
			input:    "10 +",
			expected: nil,
			hasError: true,
		},
		{
			name:     "starts_with_operator_error",
			input:    "+ 10",
			expected: nil,
			hasError: true,
		},
		{
			name:     "consecutive_numbers_error",
			input:    "10 15",
			expected: nil,
			hasError: true,
		},
		{
			name:     "empty_input_error",
			input:    "",
			expected: nil,
			hasError: true,
		},
		{
			name:     "whitespace_only_error",
			input:    "   ",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid_character_error",
			input:    "10 + abc",
			expected: nil,
			hasError: true,
		},
		{
			name:  "single_number",
			input: "42",
			expected: []Token{
				{Type: NUMBER, Value: "42"},
			},
			hasError: false,
		},
		{
			name:  "decimal_with_zero",
			input: "0.5 + 10",
			expected: []Token{
				{Type: NUMBER, Value: "0.5"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "10"},
			},
			hasError: false,
		},
		{
			name:  "multiple_decimals",
			input: "3.14 * 2.5",
			expected: []Token{
				{Type: NUMBER, Value: "3.14"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "2.5"},
			},
			hasError: false,
		},
		{
			name:  "negative_number_at_start",
			input: "-2000",
			expected: []Token{
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "2000"},
			},
			hasError: false,
		},
		{
			name:  "negative_after_operator",
			input: "10 + -5",
			expected: []Token{
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "+"},
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "5"},
			},
			hasError: false,
		},
		{
			name:  "binary_subtraction",
			input: "10 - 5",
			expected: []Token{
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "5"},
			},
			hasError: false,
		},
		{
			name:  "unary_minus_in_parentheses",
			input: "-(10 + 5)",
			expected: []Token{
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: PARENTHESIS, Value: "("},
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "5"},
				{Type: PARENTHESIS, Value: ")"},
			},
			hasError: false,
		},
		{
			name:  "unary_minus_after_opening_paren",
			input: "10 + -(2 + 4)",
			expected: []Token{
				{Type: NUMBER, Value: "10"},
				{Type: OPERATOR, Value: "+"},
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: PARENTHESIS, Value: "("},
				{Type: NUMBER, Value: "2"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "4"},
				{Type: PARENTHESIS, Value: ")"},
			},
			hasError: false,
		},
		{
			name:  "double_unary_minus",
			input: "-(-5)",
			expected: []Token{
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: PARENTHESIS, Value: "("},
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "5"},
				{Type: PARENTHESIS, Value: ")"},
			},
			hasError: false,
		},
		{
			name:  "mixed_unary_and_binary_minus",
			input: "-20 - 10",
			expected: []Token{
				{Type: UNARY_OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "20"},
				{Type: OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "10"},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Tokenize(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected tokens %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTokenString(t *testing.T) {
	tests := []struct {
		token    Token
		expected string
	}{
		{Token{Type: NUMBER, Value: "10"}, "NUMBER(10)"},
		{Token{Type: OPERATOR, Value: "+"}, "OPERATOR(+)"},
		{Token{Type: PARENTHESIS, Value: "("}, "PARENTHESIS(()"},
	}

	for _, tt := range tests {
		result := tt.token.String()
		if result != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, result)
		}
	}
}