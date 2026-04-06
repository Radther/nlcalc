// Package normalizer converts written numbers and operation words into mathematical symbols.
// It handles number words (e.g., "ten", "twenty five"), operation words (e.g., "plus", "minus"),
// and percentage expressions (e.g., "20% of", "percent of") as part of the nlcalc parsing pipeline.
package normalizer

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var numberWords = map[string]int{
	"zero":      0,
	"one":       1,
	"two":       2,
	"three":     3,
	"four":      4,
	"five":      5,
	"six":       6,
	"seven":     7,
	"eight":     8,
	"nine":      9,
	"ten":       10,
	"eleven":    11,
	"twelve":    12,
	"thirteen":  13,
	"fourteen":  14,
	"fifteen":   15,
	"sixteen":   16,
	"seventeen": 17,
	"eighteen":  18,
	"nineteen":  19,
	"twenty":    20,
	"thirty":    30,
	"forty":     40,
	"fifty":     50,
	"sixty":     60,
	"seventy":   70,
	"eighty":    80,
	"ninety":    90,
}

var scaleWords = map[string]int{
	"hundred":  100,
	"thousand": 1000,
	"million":  1000000,
	"billion":  1000000000,
}

// builtInConstants defines mathematical constants available as built-in variables.
// These are lower priority than user-provided variables.
var builtInConstants = map[string]float64{
	"pi":  3.141592653589793, // Ratio of circle's circumference to diameter
	"e":   2.718281828459045, // Euler's number, base of natural logarithms
	"tau": 6.283185307179586, // Full circle constant (2π)
	"phi": 1.618033988749895, // Golden ratio
}

// buildNumberWordsPattern creates a regex pattern from map keys
func buildNumberWordsPattern(words map[string]int) string {
	keys := make([]string, 0, len(words))
	for word := range words {
		keys = append(keys, word)
	}
	sort.Strings(keys) // Sort for deterministic pattern
	return `\b(?:` + strings.Join(keys, "|") + `)\b`
}

// buildCompoundNumberPattern creates a regex pattern for compound numbers
func buildCompoundNumberPattern() string {
	allWords := make([]string, 0, len(numberWords)+len(scaleWords)+1)
	for word := range numberWords {
		allWords = append(allWords, word)
	}
	for word := range scaleWords {
		allWords = append(allWords, word)
	}
	allWords = append(allWords, "and")
	sort.Strings(allWords) // Sort for deterministic pattern
	return `\b(?:(?:` + strings.Join(allWords, "|") + `)\s*)+\b`
}

// Compile regex patterns once at package level - built from maps for single source of truth
var individualWordPattern = regexp.MustCompile(buildNumberWordsPattern(numberWords))

// Pre-compiled regexes for percentage operations
var (
	percentageRegex = regexp.MustCompile(`\bpercentage\b`)
	percentRegex    = regexp.MustCompile(`\bpercent\b`)
	percentSymbol   = regexp.MustCompile(`%`)
)

// Pre-compiled regexes for power operations
var (
	raisedToThePowerOfRegex = regexp.MustCompile(`\braised to the power of\b`)
	squaredRegex            = regexp.MustCompile(`\bsquared\b`)
	cubedRegex              = regexp.MustCompile(`\bcubed\b`)
	raisedRegex             = regexp.MustCompile(`\braised\b`)
	powerRegex              = regexp.MustCompile(`\bpower\b`)
)

// Pre-compiled regexes for function operations
var (
	naturalLogarithmRegex = regexp.MustCompile(`\bnatural logarithm\b`)
	squareRootRegex       = regexp.MustCompile(`\bsquare root\b`)
	squarerootRegex       = regexp.MustCompile(`\bsquareroot\b`)
	absoluteRegex         = regexp.MustCompile(`\babsolute\b`)
	logarithmRegex        = regexp.MustCompile(`\blogarithm\b`)
	sineRegex             = regexp.MustCompile(`\bsine\b`)
	cosineRegex           = regexp.MustCompile(`\bcosine\b`)
	tangentRegex          = regexp.MustCompile(`\btangent\b`)
)

// Pre-compiled regexes for basic operations
var (
	plusRegex      = regexp.MustCompile(`\bplus\b`)
	addRegex       = regexp.MustCompile(`\badd\b`)
	addedToRegex   = regexp.MustCompile(`\badded to\b`)
	minusRegex     = regexp.MustCompile(`\bminus\b`)
	subtractRegex  = regexp.MustCompile(`\bsubtract\b`)
	timesRegex     = regexp.MustCompile(`\btimes\b`)
	multiplyRegex  = regexp.MustCompile(`\bmultiply\b`)
	dividedByRegex = regexp.MustCompile(`\bdivided by\b`)
	divideRegex    = regexp.MustCompile(`\bdivide\b`)
)

// Pre-compiled utility regexes
var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
	numberPattern   = regexp.MustCompile(buildCompoundNumberPattern())
)

func parseNumberPhrase(phrase string) (int, bool) {
	words := strings.Fields(strings.ToLower(phrase))
	if len(words) == 0 {
		return 0, false
	}

	// Check for invalid "and" usage (connecting same-level numbers)
	if !isValidAndUsage(words) {
		return 0, false
	}

	result := 0
	current := 0

	for _, word := range words {
		if word == "and" {
			continue
		}

		if scale, isScale := scaleWords[word]; isScale {
			if current == 0 {
				current = 1
			}

			if scale == 100 {
				// hundred: multiply current by 100
				current *= 100
			} else if scale >= 1000 {
				// For thousand, million, billion: add to result and reset current
				result += current * scale
				current = 0
			}
		} else if val, exists := numberWords[word]; exists {
			current += val
		} else {
			return 0, false
		}
	}

	return result + current, true
}

// isValidAndUsage checks if "and" is used correctly in compound numbers
func isValidAndUsage(words []string) bool {
	for i, word := range words {
		if word == "and" {
			// "and" must follow a scale word (hundred, thousand, etc.)
			if i == 0 {
				return false // "and" at beginning is invalid
			}

			prevWord := words[i-1]
			// Check if previous word is a scale word
			if _, isScale := scaleWords[prevWord]; !isScale {
				// If not a scale word, check if it's a number that could be part of a scale
				// e.g., "twenty" in "one hundred twenty and five" is valid
				// but "ten and fifteen" is not valid
				return hasRecentScaleContext(words, i)
			}
		}
	}
	return true
}

// hasRecentScaleContext checks if there's a scale word in recent context
func hasRecentScaleContext(words []string, andPos int) bool {
	// Look backwards for a scale word
	for i := andPos - 1; i >= 0; i-- {
		if _, isScale := scaleWords[words[i]]; isScale {
			return true
		}
		// If we hit another "and" or go too far back, no valid context
		if words[i] == "and" || i < andPos-3 {
			break
		}
	}
	return false
}

// replaceVariables replaces variable names in the input string with their numeric values.
// Variables are sorted by length (descending) to ensure longer variable names are replaced first,
// preventing partial replacements (e.g., "xx" should not become "valuex" when both "x" and "xx" exist).
// Uses word boundaries to ensure variables match complete words only.
func replaceVariables(input string, variables map[string]float64) string {
	if len(variables) == 0 {
		return input
	}

	// Sort variable names by length (descending) to replace longest first
	varNames := make([]string, 0, len(variables))
	for name := range variables {
		varNames = append(varNames, name)
	}
	sort.Slice(varNames, func(i, j int) bool {
		return len(varNames[i]) > len(varNames[j])
	})

	result := input
	for _, name := range varNames {
		value := variables[name]
		// Use word boundaries to match complete variable names only
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
		result = pattern.ReplaceAllString(result, strconv.FormatFloat(value, 'f', -1, 64))
	}

	return result
}

// pointBetweenDigitsPattern matches the word "point" used as a decimal separator between digits.
// Uses \s* (zero or more spaces) because the number word pattern may consume trailing whitespace,
// leaving "point" adjacent to the preceding digit (e.g. "ten point two" → "10point 2").
var pointBetweenDigitsPattern = regexp.MustCompile(`(\d)\s*point\s*(\d)`)

// Normalize converts a raw input string into a normalized mathematical expression.
// It transforms written numbers ("ten" → "10"), operation words ("plus" → "+", "power" → "^"),
// percentage notations ("20% of" → "* 0.01 *"), and power shortcuts ("squared" → "^ 2") into symbolic form.
// Variables from the optional map are replaced with their numeric values before other transformations.
// Built-in constants (pi, e, tau, phi) are also available but have lower priority than user variables.
//
// The decimalDelimiter parameter specifies the decimal separator character used in the input.
// Use '.' (or 0 for default) for standard notation (e.g. "3.14") or ',' for European notation (e.g. "3,14").
// The alternate character is treated as a thousands separator and stripped (e.g. "1,000" → "1000" in default mode,
// "1.000" → "1000" in comma mode). The word "point" is always recognized as a decimal separator.
//
// Returns a normalized string ready for cleaning and tokenization.
func Normalize(input string, variables map[string]float64, decimalDelimiter rune) string {
	if decimalDelimiter == 0 {
		decimalDelimiter = '.'
	}

	result := strings.ToLower(input)

	// Determine the thousands separator (whichever character is NOT the decimal delimiter).
	var thousandsSep rune
	if decimalDelimiter == '.' {
		thousandsSep = ','
	} else {
		thousandsSep = '.'
	}

	// Strip thousands separators (the alternate character when it appears between digits).
	thousandsSepStr := regexp.QuoteMeta(string(thousandsSep))
	thousandsPattern := regexp.MustCompile(`(\d)` + thousandsSepStr + `(\d)`)
	result = thousandsPattern.ReplaceAllString(result, "${1}${2}")

	// Normalize the decimal delimiter to '.' for standard downstream processing.
	if decimalDelimiter != '.' {
		delimStr := regexp.QuoteMeta(string(decimalDelimiter))
		decimalPattern := regexp.MustCompile(`(\d)` + delimStr + `(\d)`)
		result = decimalPattern.ReplaceAllString(result, "${1}.${2}")
	}

	// Replace user variables first (highest priority), then built-in constants (lower priority)
	// This ensures user-provided variables override built-in constants
	result = replaceVariables(result, variables)
	result = replaceVariables(result, builtInConstants)

	// Handle percentage operations - order matters (longer phrases first)
	result = percentageRegex.ReplaceAllString(result, " * 0.01 * ")
	result = percentRegex.ReplaceAllString(result, " * 0.01 * ")
	result = percentSymbol.ReplaceAllString(result, " * 0.01 * ")

	// Handle power operations - order matters to prevent double replacement
	result = raisedToThePowerOfRegex.ReplaceAllString(result, " ^ ") // Must be first to prevent double ^
	result = squaredRegex.ReplaceAllString(result, " ^ 2")            // Shortcut for ^2
	result = cubedRegex.ReplaceAllString(result, " ^ 3")              // Shortcut for ^3
	result = raisedRegex.ReplaceAllString(result, " ^ ")              // General "raised to"
	result = powerRegex.ReplaceAllString(result, " ^ ")               // General "power"

	// Handle function phrases - convert natural language to function names
	// Order matters: longer phrases must be matched before shorter ones
	result = naturalLogarithmRegex.ReplaceAllString(result, "ln ")  // Must be before "logarithm"
	result = squareRootRegex.ReplaceAllString(result, "sqrt ")
	result = squarerootRegex.ReplaceAllString(result, "sqrt ")
	result = absoluteRegex.ReplaceAllString(result, "abs ")
	result = logarithmRegex.ReplaceAllString(result, "log ")
	result = sineRegex.ReplaceAllString(result, "sin ")
	result = cosineRegex.ReplaceAllString(result, "cos ")
	result = tangentRegex.ReplaceAllString(result, "tan ")

	// Handle basic arithmetic operations
	result = addedToRegex.ReplaceAllString(result, " + ")   // Must be before "add"
	result = dividedByRegex.ReplaceAllString(result, " / ") // Must be before "divide"
	result = plusRegex.ReplaceAllString(result, " + ")
	result = addRegex.ReplaceAllString(result, " + ")
	result = minusRegex.ReplaceAllString(result, " - ")
	result = subtractRegex.ReplaceAllString(result, " - ")
	result = timesRegex.ReplaceAllString(result, " * ")
	result = multiplyRegex.ReplaceAllString(result, " * ")
	result = divideRegex.ReplaceAllString(result, " / ")

	// Single number pattern - handles both compound and individual numbers
	result = numberPattern.ReplaceAllStringFunc(result, func(match string) string {
		// Try compound number parsing first (handles complex cases and validation)
		if num, ok := parseNumberPhrase(match); ok {
			return strconv.Itoa(num)
		}

		// If compound parsing fails, convert individual words within the match
		// This handles cases like "five and ten" where compound parsing fails due to invalid "and"
		// but individual words should still be converted to "5 and 10"
		return individualWordPattern.ReplaceAllStringFunc(match, func(word string) string {
			if val, exists := numberWords[word]; exists {
				return strconv.Itoa(val)
			}
			return word
		})
	})

	// Convert "point" used as a decimal separator (e.g. "ten point two" → "10.2").
	// This runs after number word conversion so digit context is established.
	result = pointBetweenDigitsPattern.ReplaceAllString(result, "${1}.${2}")

	result = whitespaceRegex.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
