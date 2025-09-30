// Package nlcalc parses natural language mathematical expressions and returns calculated results.
//
// The package handles:
//   - Written numbers: "ten", "twenty", "five hundred"
//   - Operation words: "plus", "minus", "times", "divided by"
//   - Mathematical symbols: +, -, *, /
//   - Percentages: "20% of 100", "50 percent"
//   - Parentheses for grouping: "(10 + 5) * 2"
//   - Standard order of operations (PEMDAS/BODMAS)
//
// Example usage:
//
//	result, err := nlcalc.Parse("ten plus fifteen")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result) // Output: 25
//
//	result, err = nlcalc.Parse("20% of 100")
//	fmt.Println(result) // Output: 20
package nlcalc

import (
	"github.com/radther/nlcalc/internal/cleaner"
	"github.com/radther/nlcalc/internal/evaluator"
	"github.com/radther/nlcalc/internal/normalizer"
	"github.com/radther/nlcalc/internal/tokenizer"
)

// Parse converts a natural language mathematical expression into a calculated result.
//
// The function accepts expressions in various formats:
//   - Written numbers: "five hundred twenty three"
//   - Words for operations: "ten plus five", "twenty minus eight"
//   - Standard math symbols: "5 + 10", "100 / 4"
//   - Percentages: "15% of 200", "50 percent of 80"
//   - Mixed formats: "twenty * 3", "5 plus ten"
//   - Parentheses: "(10 + 5) * 2"
//
// Order of operations follows standard PEMDAS/BODMAS rules.
// For example, "10 + 10 * 4" evaluates to 50 (not 80).
//
// Returns an error if the expression is invalid, contains unrecognized
// characters, has malformed syntax, or attempts division by zero.
func Parse(input string) (float64, error) {
	normalized := normalizer.Normalize(input)
	cleaned := cleaner.Clean(normalized)
	tokens, err := tokenizer.Tokenize(cleaned)
	if err != nil {
		return 0, err
	}
	
	result, err := evaluator.Evaluate(tokens)
	if err != nil {
		return 0, err
	}
	
	return result, nil
}