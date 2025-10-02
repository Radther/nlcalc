package normalizer

import (
	"testing"
)

// BenchmarkNormalize_Simple benchmarks basic number normalization
func BenchmarkNormalize_Simple(b *testing.B) {
	input := "ten plus fifteen"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_Complex benchmarks complex expression with all features
func BenchmarkNormalize_Complex(b *testing.B) {
	input := "twenty percent of square root of one hundred twenty five plus three squared"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_CompoundNumbers benchmarks compound number parsing
func BenchmarkNormalize_CompoundNumbers(b *testing.B) {
	input := "two thousand five hundred and seventy three plus one hundred and forty two"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_WithVariables benchmarks variable replacement
func BenchmarkNormalize_WithVariables(b *testing.B) {
	input := "price times quantity plus tax percent of total"
	vars := map[string]float64{
		"price":    25.50,
		"quantity": 10,
		"tax":      15,
		"total":    255,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, vars)
	}
}

// BenchmarkNormalize_Functions benchmarks function phrase normalization
func BenchmarkNormalize_Functions(b *testing.B) {
	input := "square root of absolute value of sine of ninety"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_PowerOperations benchmarks power phrase normalization
func BenchmarkNormalize_PowerOperations(b *testing.B) {
	input := "five squared plus two raised to the power of three"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_Percentage benchmarks percentage normalization
func BenchmarkNormalize_Percentage(b *testing.B) {
	input := "twenty percent of one hundred plus fifteen percentage of fifty"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Normalize(input, nil)
	}
}

// BenchmarkNormalize_Batch simulates real-world usage with varied inputs
func BenchmarkNormalize_Batch(b *testing.B) {
	inputs := []string{
		"ten plus fifteen",
		"twenty percent of one hundred",
		"five squared",
		"square root of sixteen",
		"two thousand five hundred",
		"price times quantity",
		"three raised to the power of four",
		"absolute value of negative five",
		"one hundred and twenty seven plus three hundred forty five",
		"natural logarithm of ten",
	}
	vars := map[string]float64{"price": 100, "quantity": 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_ = Normalize(input, vars)
		}
	}
}

// BenchmarkNormalize_LargeScale simulates high-volume usage
func BenchmarkNormalize_LargeScale(b *testing.B) {
	inputs := []string{
		"ten plus fifteen",
		"twenty percent of one hundred",
		"five squared plus two cubed",
		"square root of one hundred twenty five",
		"two thousand five hundred and seventy three",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate 100 calculations per iteration
		for j := 0; j < 100; j++ {
			_ = Normalize(inputs[j%len(inputs)], nil)
		}
	}
}
