package tokenizer

import (
	"testing"
)

// BenchmarkTokenize_Simple benchmarks basic tokenization
func BenchmarkTokenize_Simple(b *testing.B) {
	input := "10 + 15"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_Complex benchmarks complex expression tokenization
func BenchmarkTokenize_Complex(b *testing.B) {
	input := "20 * 0.01 * 100 + 3 ^ 2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_WithParentheses benchmarks tokenization with parentheses
func BenchmarkTokenize_WithParentheses(b *testing.B) {
	input := "( 10 + 15 ) * 2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_Functions benchmarks tokenization with functions
func BenchmarkTokenize_Functions(b *testing.B) {
	input := "sqrt 16 + abs -5"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_PowerOperations benchmarks tokenization with power operations
func BenchmarkTokenize_PowerOperations(b *testing.B) {
	input := "5 ^ 2 + 2 ^ 3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_UnaryOperators benchmarks tokenization with unary operators
func BenchmarkTokenize_UnaryOperators(b *testing.B) {
	input := "-5 + -10 * +3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_Decimals benchmarks tokenization with decimal numbers
func BenchmarkTokenize_Decimals(b *testing.B) {
	input := "25.5 * 10 + 0.15 * 255"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Tokenize(input)
	}
}

// BenchmarkTokenize_Batch simulates real-world usage with varied inputs
func BenchmarkTokenize_Batch(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_, _ = Tokenize(input)
		}
	}
}

// BenchmarkTokenize_LargeScale simulates high-volume usage
func BenchmarkTokenize_LargeScale(b *testing.B) {
	inputs := []string{
		"10 + 15",
		"20 * 0.01 * 100",
		"5 ^ 2 + 2 ^ 3",
		"sqrt 125",
		"2573",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate 100 calculations per iteration
		for j := 0; j < 100; j++ {
			_, _ = Tokenize(inputs[j%len(inputs)])
		}
	}
}
