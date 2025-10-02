package nlcalc

import (
	"testing"
)

// =============================================================================
// PART 1: STANDARD ADAPTIVE BENCHMARKS (auto-adjusting iterations)
// =============================================================================

// -----------------------------------------------------------------------------
// Full Pipeline Benchmarks - Parse() function
// -----------------------------------------------------------------------------

func BenchmarkParse_Simple(b *testing.B) {
	input := "ten plus fifteen"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_Complex(b *testing.B) {
	input := "twenty percent of square root of one hundred twenty five plus three squared"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_CompoundNumbers(b *testing.B) {
	input := "two thousand five hundred and seventy three plus one hundred and forty two"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_WithVariables(b *testing.B) {
	input := "price times quantity plus tax percent of total"
	vars := map[string]float64{
		"price":    25.50,
		"quantity": 10,
		"tax":      15,
		"total":    255,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, vars)
	}
}

func BenchmarkParse_Functions(b *testing.B) {
	input := "square root of absolute value of sine of ninety"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_PowerOperations(b *testing.B) {
	input := "five squared plus two raised to the power of three"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_Percentage(b *testing.B) {
	input := "twenty percent of one hundred plus fifteen percentage of fifty"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input, nil)
	}
}

func BenchmarkParse_Batch(b *testing.B) {
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
			_, _ = Parse(input, vars)
		}
	}
}

// =============================================================================
// PART 2: FIXED ITERATION BENCHMARKS (consistent runs for tracking)
// =============================================================================

// -----------------------------------------------------------------------------
// Fixed 1000 Iterations - Quick Performance Snapshot
// -----------------------------------------------------------------------------

func BenchmarkParse_Fixed1000_Simple(b *testing.B) {
	input := "ten plus fifteen"
	const iterations = 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, nil)
		}
	}
}

func BenchmarkParse_Fixed1000_Complex(b *testing.B) {
	input := "twenty percent of square root of one hundred twenty five plus three squared"
	const iterations = 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, nil)
		}
	}
}

func BenchmarkParse_Fixed1000_WithVariables(b *testing.B) {
	input := "price times quantity plus tax percent of total"
	vars := map[string]float64{
		"price":    25.50,
		"quantity": 10,
		"tax":      15,
		"total":    255,
	}
	const iterations = 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, vars)
		}
	}
}

// -----------------------------------------------------------------------------
// Fixed 10000 Iterations - Detailed Performance Measurement
// -----------------------------------------------------------------------------

func BenchmarkParse_Fixed10000_Simple(b *testing.B) {
	input := "ten plus fifteen"
	const iterations = 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, nil)
		}
	}
}

func BenchmarkParse_Fixed10000_Complex(b *testing.B) {
	input := "twenty percent of square root of one hundred twenty five plus three squared"
	const iterations = 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, nil)
		}
	}
}

func BenchmarkParse_Fixed10000_WithVariables(b *testing.B) {
	input := "price times quantity plus tax percent of total"
	vars := map[string]float64{
		"price":    25.50,
		"quantity": 10,
		"tax":      15,
		"total":    255,
	}
	const iterations = 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, vars)
		}
	}
}

func BenchmarkParse_Fixed10000_FullPipeline(b *testing.B) {
	input := "twenty percent of square root of one hundred twenty five plus three squared"
	const iterations = 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			_, _ = Parse(input, nil)
		}
	}
}
