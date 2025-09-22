package cleaner

import (
	"regexp"
	"strings"
)

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