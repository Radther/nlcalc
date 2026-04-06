package favourable

import (
	"reflect"
	"testing"

	"github.com/radther/nlcalc/internal/tokenizer"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		input    []tokenizer.Token
		expected []tokenizer.Token
	}{
		// ── Rule 1: strip leading invalid tokens ──────────────────────────────
		{
			name: "strip_leading_operator",
			// "* 10 + 2" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "strip_multiple_leading_operators",
			// "* / 10 + 2" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.OPERATOR, Value: "/"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "strip_leading_closing_paren",
			// ") 10 + 2" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "preserve_leading_unary_operator",
			// "-10 + 2" is valid — UNARY_OPERATOR at start is fine
			input: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "preserve_leading_function",
			// "sqrt 16" is valid — FUNCTION at start is fine
			input: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
		},

		// ── Rule 2: strip trailing invalid tokens ─────────────────────────────
		{
			name: "strip_trailing_operator",
			// "10 + 2 /" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "/"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "strip_trailing_unary_operator_cascades",
			// "10 + -" → strip "-" (UNARY_OPERATOR) → "10 +" → strip "+" (OPERATOR) → "10"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
			},
		},
		{
			name: "strip_trailing_function_cascades",
			// "10 + sqrt" → strip "sqrt" (FUNCTION) → "10 +" → strip "+" (OPERATOR) → "10"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
			},
		},
		{
			name: "strip_trailing_open_paren_cascades",
			// "10 + (" → strip "(" → "10 +" → strip "+" → "10"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.PARENTHESIS, Value: "("},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
			},
		},

		// ── Both leading and trailing ─────────────────────────────────────────
		{
			name: "strip_leading_and_trailing_operators",
			// "* 10 + 2 /" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "/"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},

		// ── Rule 3: collapse consecutive OPERATOR tokens ──────────────────────
		{
			name: "collapse_consecutive_operators_keep_first",
			// "10 +* 2" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "collapse_three_consecutive_operators",
			// "10 +** 2" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "preserve_consecutive_functions",
			// "sqrt sqrt 16" is valid — consecutive FUNCTIONs are intentional
			input: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.FUNCTION, Value: "sqrt"},
				{Type: tokenizer.NUMBER, Value: "16"},
			},
		},
		{
			name: "preserve_consecutive_unary_operators",
			// "--5" (double negation) is valid — consecutive UNARY_OPERATORs are intentional
			input: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.UNARY_OPERATOR, Value: "-"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
		},
		{
			name: "preserve_consecutive_numbers_unchanged",
			// "10 15" — consecutive numbers are left as-is (still invalid, no fix possible)
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.NUMBER, Value: "15"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.NUMBER, Value: "15"},
			},
		},

		// ── Rule 4: strip unmatched parentheses ───────────────────────────────
		{
			name: "strip_unmatched_opening_paren",
			// "(10 + 5" → "10 + 5"
			input: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
		},
		{
			name: "strip_unmatched_closing_paren",
			// "10 + 5)" → "10 + 5"
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
		},
		{
			name: "preserve_matched_parentheses",
			// "(10 + 5) * 2" — parens are matched, leave them alone
			input: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "strip_both_unmatched_parens",
			// "(10 + 5" with an extra stray ")" somewhere
			// "( 10 + 5 ) )" → strip outer unmatched ")" → "(10 + 5)"
			input: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.PARENTHESIS, Value: "("},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
				{Type: tokenizer.PARENTHESIS, Value: ")"},
			},
		},

		// ── Combined rules ────────────────────────────────────────────────────
		{
			name: "combined_leading_operator_and_consecutive_operators",
			// "* 10 +* 2 /" → "10 + 2"
			input: []tokenizer.Token{
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.OPERATOR, Value: "*"},
				{Type: tokenizer.NUMBER, Value: "2"},
				{Type: tokenizer.OPERATOR, Value: "/"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},

		// ── No-op cases: already valid expressions ───────────────────────────
		{
			name: "no_change_simple_expression",
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name: "no_change_single_number",
			input: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "42"},
			},
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "42"},
			},
		},
		{
			name: "empty_input_unchanged",
			input:    []tokenizer.Token{},
			expected: []tokenizer.Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Apply(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Apply():\n  got  %v\n  want %v", result, tt.expected)
			}
		})
	}
}

// TestApplyIntegration tests the favourable step combined with tokenization,
// covering the full pipeline from raw string to recovered token sequence.
func TestApplyIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []tokenizer.Token
	}{
		{
			name:  "leading_operator_in_string",
			input: "* 10 + 2",
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name:  "trailing_operator_in_string",
			input: "10 + 2 /",
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name:  "consecutive_operators_in_string",
			input: "10 +* 2",
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "2"},
			},
		},
		{
			name:  "unmatched_opening_paren_in_string",
			input: "(10 + 5",
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
		},
		{
			name:  "unmatched_closing_paren_in_string",
			input: "10 + 5)",
			expected: []tokenizer.Token{
				{Type: tokenizer.NUMBER, Value: "10"},
				{Type: tokenizer.OPERATOR, Value: "+"},
				{Type: tokenizer.NUMBER, Value: "5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.TokenizeRaw(tt.input)
			if err != nil {
				t.Fatalf("TokenizeRaw error: %v", err)
			}
			result := Apply(tokens)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Apply():\n  got  %v\n  want %v", result, tt.expected)
			}
		})
	}
}
