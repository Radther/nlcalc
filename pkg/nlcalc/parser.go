// Package nlcalc parses natural language mathematical expressions and returns calculated results.
//
// The package handles:
//   - Written numbers: "ten", "twenty", "five hundred"
//   - Operation words: "plus", "minus", "times", "divided by"
//   - Mathematical symbols: +, -, *, /
//   - Percentages: "20% of 100", "50 percent"
//   - Parentheses for grouping: "(10 + 5) * 2"
//   - Variables: custom named values passed as a map
//   - Standard order of operations (PEMDAS/BODMAS)
//   - Thousands separators: "10,000" → 10000 (default), "10.000" → 10000 (comma decimal mode)
//   - Decimal separator word: "ten point two" → 10.2
//
// Example usage:
//
//	result, err := nlcalc.Parse("ten plus fifteen", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result) // Output: 25
//
//	result, err = nlcalc.Parse("20% of 100", nil)
//	fmt.Println(result) // Output: 20
//
//	// With variables
//	vars := map[string]float64{"price": 100, "tax": 0.15}
//	result, err = nlcalc.Parse("price plus price times tax", vars)
//	fmt.Println(result) // Output: 115
//
//	// With comma as decimal delimiter (European notation)
//	result, err = nlcalc.ParseWithOptions("1.000,83 + 2,5", nil, nlcalc.Options{DecimalDelimiter: ','})
//	fmt.Println(result) // Output: 1003.33
package nlcalc

import (
	"github.com/radther/nlcalc/internal/cleaner"
	"github.com/radther/nlcalc/internal/evaluator"
	"github.com/radther/nlcalc/internal/favourable"
	"github.com/radther/nlcalc/internal/normalizer"
	"github.com/radther/nlcalc/internal/tokenizer"
)

// Options configures parsing behavior.
type Options struct {
	// DecimalDelimiter specifies the character used as the decimal separator in numeric input.
	// Use '.' (default, zero value) for standard notation (e.g. "3.14") or ',' for European
	// notation (e.g. "3,14"). The alternate character is treated as a thousands separator and
	// stripped automatically (e.g. "1,000" → 1000 in default mode; "1.000" → 1000 in comma mode).
	// The word "point" is always recognized as a decimal separator regardless of this setting.
	DecimalDelimiter rune

	// Favourable enables post-tokenization heuristics that attempt to recover a
	// calculable expression from one that would otherwise be invalid. When true,
	// the parser applies a set of recovery rules after tokenization:
	//   - Leading/trailing operators and other invalid boundary tokens are stripped.
	//   - Consecutive binary operators are collapsed to the first one.
	//   - Unmatched parentheses are removed.
	//
	// Example: "* 10 +* 2 /" → recovered as "10 + 2" → result 12.
	// Consecutive numbers (e.g. "10 15") are not recoverable and still return an error.
	Favourable bool
}

// ParseWithOptions converts a natural language mathematical expression into a calculated result,
// using the provided options to control parsing behavior.
//
// See Parse for a full description of supported expression formats.
func ParseWithOptions(input string, variables map[string]float64, options Options) (float64, error) {
	normalized := normalizer.Normalize(input, variables, options.DecimalDelimiter)
	cleaned := cleaner.Clean(normalized)

	var tokens []tokenizer.Token
	var err error

	if options.Favourable {
		tokens, err = tokenizer.TokenizeRaw(cleaned)
		if err != nil {
			return 0, err
		}
		tokens = favourable.Apply(tokens)
	} else {
		tokens, err = tokenizer.Tokenize(cleaned)
		if err != nil {
			return 0, err
		}
	}

	result, err := evaluator.Evaluate(tokens)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Parse converts a natural language mathematical expression into a calculated result.
//
// The function accepts expressions in various formats:
//   - Written numbers: "five hundred twenty three"
//   - Words for operations: "ten plus five", "twenty minus eight"
//   - Standard math symbols: "5 + 10", "100 / 4"
//   - Percentages: "15% of 200", "50 percent of 80"
//   - Mixed formats: "twenty * 3", "5 plus ten"
//   - Parentheses: "(10 + 5) * 2"
//   - Variables: Named values provided in the variables map
//   - Thousands separators: "10,000" → 10000 (comma is stripped as thousands separator)
//   - Decimal word: "ten point two" → 10.2
//
// The variables parameter is optional (can be nil). Variable names in the input
// will be replaced with their corresponding numeric values before evaluation.
// Variable names are case-insensitive and matched using word boundaries.
//
// Order of operations follows standard PEMDAS/BODMAS rules.
// For example, "10 + 10 * 4" evaluates to 50 (not 80).
//
// Returns an error if the expression is invalid, contains unrecognized
// characters, has malformed syntax, or attempts division by zero.
func Parse(input string, variables map[string]float64) (float64, error) {
	return ParseWithOptions(input, variables, Options{})
}
