package evaluator

import (
	"nlcalc/internal/tokenizer"
	"math"
	"testing"
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