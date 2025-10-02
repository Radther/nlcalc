package cleaner

import (
	"testing"
)

// BenchmarkClean_Simple benchmarks basic cleaning
func BenchmarkClean_Simple(b *testing.B) {
	input := "10 + 15"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_Complex benchmarks complex expression cleaning
func BenchmarkClean_Complex(b *testing.B) {
	input := "20 * 0.01 * sqrt 125 + 3 ^ 2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_WithFillerWords benchmarks cleaning with filler words
func BenchmarkClean_WithFillerWords(b *testing.B) {
	input := "what is 20 * 0.01 * of 100 plus some 3 ^ 2"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_Functions benchmarks cleaning with functions
func BenchmarkClean_Functions(b *testing.B) {
	input := "sqrt of abs value of sin 90"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_PowerOperations benchmarks cleaning with power operations
func BenchmarkClean_PowerOperations(b *testing.B) {
	input := "5 ^ 2 + 2 ^ to the 3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_Percentage benchmarks cleaning percentage expressions
func BenchmarkClean_Percentage(b *testing.B) {
	input := "20 * 0.01 * of 100 + 15 * 0.01 * of 50"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clean(input)
	}
}

// BenchmarkClean_Batch simulates real-world usage with varied inputs
func BenchmarkClean_Batch(b *testing.B) {
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
			_ = Clean(input)
		}
	}
}

// BenchmarkClean_LargeScale simulates high-volume usage
func BenchmarkClean_LargeScale(b *testing.B) {
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
			_ = Clean(inputs[j%len(inputs)])
		}
	}
}
