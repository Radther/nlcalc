package evaluator

import (
	"mathparser/internal/tokenizer"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []tokenizer.Token
		expected float64
	}{
		{
			name:     "empty_tokens",
			tokens:   []tokenizer.Token{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.tokens)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}