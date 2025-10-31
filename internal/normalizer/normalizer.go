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
	cubeRootRegex         = regexp.MustCompile(`\bcube root\b`)
	cuberootRegex         = regexp.MustCompile(`\bcuberoot\b`)
	absoluteRegex         = regexp.MustCompile(`\babsolute\b`)
	logarithmRegex        = regexp.MustCompile(`\blogarithm\b`)
	sineRegex             = regexp.MustCompile(`\bsine\b`)
	cosineRegex           = regexp.MustCompile(`\bcosine\b`)
	tangentRegex          = regexp.MustCompile(`\btangent\b`)
	// Inverse trigonometric functions
	arcsineRegex        = regexp.MustCompile(`\barc sine\b`)
	arcsinRegex         = regexp.MustCompile(`\barcsin\b`)
	arcsineAltRegex     = regexp.MustCompile(`\barcsine\b`)
	arccosineRegex      = regexp.MustCompile(`\barc cosine\b`)
	arccosRegex         = regexp.MustCompile(`\barccos\b`)
	arccosineAltRegex   = regexp.MustCompile(`\barccosine\b`)
	arctangentRegex     = regexp.MustCompile(`\barc tangent\b`)
	arctanRegex         = regexp.MustCompile(`\barctan\b`)
	arctangentAltRegex  = regexp.MustCompile(`\barctangent\b`)
	// Hyperbolic functions
	hyperbolicSineRegex     = regexp.MustCompile(`\bhyperbolic sine\b`)
	hyperbolicCosineRegex   = regexp.MustCompile(`\bhyperbolic cosine\b`)
	hyperbolicTangentRegex  = regexp.MustCompile(`\bhyperbolic tangent\b`)
	// Inverse hyperbolic functions
	inverseHyperbolicSineRegex    = regexp.MustCompile(`\binverse hyperbolic sine\b`)
	inverseHyperbolicCosineRegex  = regexp.MustCompile(`\binverse hyperbolic cosine\b`)
	inverseHyperbolicTangentRegex = regexp.MustCompile(`\binverse hyperbolic tangent\b`)
	arcsinhRegex                  = regexp.MustCompile(`\barcsinh\b`)
	arccoshRegex                  = regexp.MustCompile(`\barccosh\b`)
	arctanhRegex                  = regexp.MustCompile(`\barctanh\b`)
	// Extended trigonometric functions
	secantRegex    = regexp.MustCompile(`\bsecant\b`)
	cosecantRegex  = regexp.MustCompile(`\bcosecant\b`)
	cotangentRegex = regexp.MustCompile(`\bcotangent\b`)
	// Inverse extended trigonometric functions
	arcsecantRegex     = regexp.MustCompile(`\barc secant\b`)
	arcsecRegex        = regexp.MustCompile(`\barcsec\b`)
	arcsecantAltRegex  = regexp.MustCompile(`\barcsecant\b`)
	arccosecantRegex   = regexp.MustCompile(`\barc cosecant\b`)
	arccscRegex        = regexp.MustCompile(`\barccsc\b`)
	arccosecantAltRegex = regexp.MustCompile(`\barccosecant\b`)
	arccotangentRegex  = regexp.MustCompile(`\barc cotangent\b`)
	arccotRegex        = regexp.MustCompile(`\barccot\b`)
	arccotangentAltRegex = regexp.MustCompile(`\barccotangent\b`)
	// Exponential and logarithmic functions
	exponentialRegex = regexp.MustCompile(`\bexponential\b`)
	logarithmBase2Regex = regexp.MustCompile(`\blogarithm base 2\b`)
	log2WordRegex = regexp.MustCompile(`\blog two\b`)
	expm1Regex = regexp.MustCompile(`\bexpm1\b`)
	log1pRegex = regexp.MustCompile(`\blog1p\b`)
	// Rounding functions
	ceilingRegex = regexp.MustCompile(`\bceiling\b`)
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

// Normalize converts a raw input string into a normalized mathematical expression.
// It transforms written numbers ("ten" → "10"), operation words ("plus" → "+", "power" → "^"),
// percentage notations ("20% of" → "* 0.01 *"), and power shortcuts ("squared" → "^ 2") into symbolic form.
// Variables from the optional map are replaced with their numeric values before other transformations.
// Returns a normalized string ready for cleaning and tokenization.
func Normalize(input string, variables map[string]float64) string {
	result := strings.ToLower(input)

	// Replace variables first, before other transformations
	result = replaceVariables(result, variables)

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
	result = logarithmBase2Regex.ReplaceAllString(result, "log2 ")  // Must be before general "logarithm"
	result = log2WordRegex.ReplaceAllString(result, "log2 ")
	result = squareRootRegex.ReplaceAllString(result, "sqrt ")
	result = squarerootRegex.ReplaceAllString(result, "sqrt ")
	result = cubeRootRegex.ReplaceAllString(result, "cbrt ")
	result = cuberootRegex.ReplaceAllString(result, "cbrt ")
	result = absoluteRegex.ReplaceAllString(result, "abs ")
	result = logarithmRegex.ReplaceAllString(result, "log ")
	// Inverse trigonometric functions - longer phrases first
	result = arcsineRegex.ReplaceAllString(result, "asin ")        // "arc sine" before "arcsine"
	result = arcsineAltRegex.ReplaceAllString(result, "asin ")     // "arcsine" before "arcsin"
	result = arcsinRegex.ReplaceAllString(result, "asin ")
	result = arccosineRegex.ReplaceAllString(result, "acos ")      // "arc cosine" before "arccosine"
	result = arccosineAltRegex.ReplaceAllString(result, "acos ")   // "arccosine" before "arccos"
	result = arccosRegex.ReplaceAllString(result, "acos ")
	result = arctangentRegex.ReplaceAllString(result, "atan ")     // "arc tangent" before "arctangent"
	result = arctangentAltRegex.ReplaceAllString(result, "atan ")  // "arctangent" before "arctan"
	result = arctanRegex.ReplaceAllString(result, "atan ")
	// Exponential and logarithmic functions
	result = exponentialRegex.ReplaceAllString(result, "exp ")
	result = expm1Regex.ReplaceAllString(result, "expm1 ")
	result = log1pRegex.ReplaceAllString(result, "log1p ")
	// Inverse hyperbolic functions - longer phrases first
	result = inverseHyperbolicSineRegex.ReplaceAllString(result, "asinh ")
	result = inverseHyperbolicCosineRegex.ReplaceAllString(result, "acosh ")
	result = inverseHyperbolicTangentRegex.ReplaceAllString(result, "atanh ")
	result = arcsinhRegex.ReplaceAllString(result, "asinh ")
	result = arccoshRegex.ReplaceAllString(result, "acosh ")
	result = arctanhRegex.ReplaceAllString(result, "atanh ")
	// Hyperbolic functions - longer phrases first
	result = hyperbolicSineRegex.ReplaceAllString(result, "sinh ")
	result = hyperbolicCosineRegex.ReplaceAllString(result, "cosh ")
	result = hyperbolicTangentRegex.ReplaceAllString(result, "tanh ")
	// Inverse extended trigonometric functions - longer phrases first
	result = arcsecantRegex.ReplaceAllString(result, "asec ")
	result = arcsecantAltRegex.ReplaceAllString(result, "asec ")
	result = arcsecRegex.ReplaceAllString(result, "asec ")
	result = arccosecantRegex.ReplaceAllString(result, "acsc ")
	result = arccosecantAltRegex.ReplaceAllString(result, "acsc ")
	result = arccscRegex.ReplaceAllString(result, "acsc ")
	result = arccotangentRegex.ReplaceAllString(result, "acot ")
	result = arccotangentAltRegex.ReplaceAllString(result, "acot ")
	result = arccotRegex.ReplaceAllString(result, "acot ")
	// Extended trigonometric functions
	result = secantRegex.ReplaceAllString(result, "sec ")
	result = cosecantRegex.ReplaceAllString(result, "csc ")
	result = cotangentRegex.ReplaceAllString(result, "cot ")
	// Basic trigonometric functions
	result = sineRegex.ReplaceAllString(result, "sin ")
	result = cosineRegex.ReplaceAllString(result, "cos ")
	result = tangentRegex.ReplaceAllString(result, "tan ")
	// Rounding functions
	result = ceilingRegex.ReplaceAllString(result, "ceil ")

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

	result = whitespaceRegex.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
