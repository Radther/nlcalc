// Package cleaner removes filler words and non-mathematical text from normalized expressions.
// It extracts only valid mathematical tokens (numbers, operators, parentheses) from the input.
package cleaner

import (
	"regexp"
	"strings"
)

// Clean removes all non-mathematical text from a normalized expression.
// It extracts and preserves only numbers, operators (+, -, *, /), parentheses,
// and percentage conversion tokens. Filler words like "what is" are removed.
// Returns a cleaned string containing only mathematical tokens, ready for tokenization.
func Clean(input string) string {
	// Define what constitutes valid mathematical tokens
	// Numbers (including decimals), operators, and parentheses
	mathTokenPattern := regexp.MustCompile(`(\d+(?:\.\d+)?|[+\-*/()]|\*\s*0\.01\s*\*)`)

	// Find all mathematical tokens
	tokens := mathTokenPattern.FindAllString(input, -1)

	// Join them back with spaces
	result := strings.Join(tokens, " ")

	// Clean up extra whitespace
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}