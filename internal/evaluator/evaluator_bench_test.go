package evaluator

import (
	"testing"

	"github.com/radther/nlcalc/internal/tokenizer"
)

// BenchmarkEvaluate_Simple benchmarks basic evaluation
func BenchmarkEvaluate_Simple(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("10 + 15")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_Complex benchmarks complex expression evaluation
func BenchmarkEvaluate_Complex(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("20 * 0.01 * 100 + 3 ^ 2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_WithParentheses benchmarks evaluation with parentheses
func BenchmarkEvaluate_WithParentheses(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("( 10 + 15 ) * 2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_NestedParentheses benchmarks evaluation with nested parentheses
func BenchmarkEvaluate_NestedParentheses(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("( ( 10 + 5 ) * 2 ) + ( 3 ^ 2 )")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_Functions benchmarks evaluation with functions
func BenchmarkEvaluate_Functions(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("sqrt 16 + abs -5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_PowerOperations benchmarks evaluation with power operations
func BenchmarkEvaluate_PowerOperations(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("5 ^ 2 + 2 ^ 3")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_UnaryOperators benchmarks evaluation with unary operators
func BenchmarkEvaluate_UnaryOperators(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("-5 + -10 * +3")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_OrderOfOperations benchmarks PEMDAS evaluation
func BenchmarkEvaluate_OrderOfOperations(b *testing.B) {
	tokens, _ := tokenizer.Tokenize("10 + 10 * 4 - 5 / 2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(tokens)
	}
}

// BenchmarkEvaluate_Batch simulates real-world usage with varied inputs
func BenchmarkEvaluate_Batch(b *testing.B) {
	tokensList := [][]tokenizer.Token{}
	inputs := []string{
		"10 + 15",
		"20 * 0.01 * 100",
		"5 ^ 2",
		"sqrt 16",
		"2500",
		"100 * 5",
		"3 ^ 4",
		"abs -5",
		"127 + 345",
		"ln 10",
	}

	for _, input := range inputs {
		tokens, _ := tokenizer.Tokenize(input)
		tokensList = append(tokensList, tokens)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tokens := range tokensList {
			_, _ = Evaluate(tokens)
		}
	}
}

// BenchmarkEvaluate_LargeScale simulates high-volume usage
func BenchmarkEvaluate_LargeScale(b *testing.B) {
	tokensList := [][]tokenizer.Token{}
	inputs := []string{
		"10 + 15",
		"20 * 0.01 * 100",
		"5 ^ 2 + 2 ^ 3",
		"sqrt 125",
		"2573",
	}

	for _, input := range inputs {
		tokens, _ := tokenizer.Tokenize(input)
		tokensList = append(tokensList, tokens)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate 100 calculations per iteration
		for j := 0; j < 100; j++ {
			_, _ = Evaluate(tokensList[j%len(tokensList)])
		}
	}
}
