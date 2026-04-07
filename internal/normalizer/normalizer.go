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

// phraseReplacement represents a multi-word phrase and its replacement.
// The words field contains the remaining words after the first (the map key).
type phraseReplacement struct {
	words       []string
	replacement string
}

// multiWordPhrases maps the first word of a phrase to its possible completions.
// Entries with more words are listed first for greedy (longest) matching.
var multiWordPhrases = map[string][]phraseReplacement{
	"raised": {
		{words: []string{"to", "the", "power", "of"}, replacement: " ^ "},
	},
	"natural": {
		{words: []string{"logarithm"}, replacement: "ln "},
	},
	"square": {
		{words: []string{"root"}, replacement: "sqrt "},
	},
	"divided": {
		{words: []string{"by"}, replacement: " / "},
	},
	"added": {
		{words: []string{"to"}, replacement: " + "},
	},
}

// wordReplacements maps single keywords to their mathematical symbol replacements.
var wordReplacements = map[string]string{
	"percentage": " * 0.01 * ",
	"percent":    " * 0.01 * ",
	"squareroot": "sqrt ",
	"squared":    " ^ 2",
	"cubed":      " ^ 3",
	"raised":     " ^ ",
	"power":      " ^ ",
	"absolute":   "abs ",
	"logarithm":  "log ",
	"sine":       "sin ",
	"cosine":     "cos ",
	"tangent":    "tan ",
	"plus":       " + ",
	"add":        " + ",
	"minus":      " - ",
	"subtract":   " - ",
	"times":      " * ",
	"multiply":   " * ",
	"divide":     " / ",
}

// Pre-compiled utility regexes
var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// Pre-compiled delimiter patterns for thousands separator and decimal normalization.
var (
	digitCommaDigit = regexp.MustCompile(`(\d),(\d)`)
	digitDotDigit   = regexp.MustCompile(`(\d)\.(\d)`)
)

// pointBetweenDigitsPattern matches the word "point" used as a decimal separator between digits.
// Uses \s* (zero or more spaces) because the number word pattern may consume trailing whitespace,
// leaving "point" adjacent to the preceding digit (e.g. "ten point two" → "10point 2").
var pointBetweenDigitsPattern = regexp.MustCompile(`(\d)\s*point\s*(\d)`)

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

// isWordChar returns true for characters that constitute word boundaries in regex terms.
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
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
		value := strconv.FormatFloat(variables[name], 'f', -1, 64)
		nameLen := len(name)

		var b strings.Builder
		b.Grow(len(result))
		i := 0
		for {
			idx := strings.Index(result[i:], name)
			if idx == -1 {
				b.WriteString(result[i:])
				break
			}
			pos := i + idx
			// Check word boundaries
			startOk := pos == 0 || !isWordChar(result[pos-1])
			endOk := pos+nameLen == len(result) || !isWordChar(result[pos+nameLen])
			if startOk && endOk {
				b.WriteString(result[i:pos])
				b.WriteString(value)
				i = pos + nameLen
			} else {
				b.WriteString(result[i : pos+1])
				i = pos + 1
			}
		}
		result = b.String()
	}

	return result
}

// replaceWords performs a single-pass scan over whitespace-delimited words, replacing
// keywords with mathematical symbols and converting number words to digits.
// It handles multi-word phrases (e.g., "raised to the power of"), single-word keywords
// (e.g., "plus"), and compound number phrases (e.g., "one hundred and twenty").
func replaceWords(words []string) string {
	if len(words) == 0 {
		return ""
	}

	var b strings.Builder
	// Pre-allocate: most replacements are similar length to originals
	b.Grow(len(words) * 6)

	i := 0
	for i < len(words) {
		if i > 0 {
			b.WriteByte(' ')
		}

		word := words[i]

		// 1. Try multi-word phrase match (longest first)
		if phrases, ok := multiWordPhrases[word]; ok {
			matched := false
			for _, phrase := range phrases {
				phraseLen := len(phrase.words)
				if i+1+phraseLen <= len(words) {
					match := true
					for j, pw := range phrase.words {
						if words[i+1+j] != pw {
							match = false
							break
						}
					}
					if match {
						b.WriteString(phrase.replacement)
						i += 1 + phraseLen
						matched = true
						break
					}
				}
			}
			if matched {
				continue
			}
		}

		// 2. Single-word keyword replacement
		if replacement, ok := wordReplacements[word]; ok {
			b.WriteString(replacement)
			i++
			continue
		}

		// 3. Number word accumulation and compound number parsing
		_, isNum := numberWords[word]
		_, isScale := scaleWords[word]
		if isNum || isScale {
			// Accumulate consecutive number/scale words and "and"
			j := i + 1
			for j < len(words) {
				w := words[j]
				_, wIsNum := numberWords[w]
				_, wIsScale := scaleWords[w]
				if wIsNum || wIsScale || w == "and" {
					j++
				} else {
					break
				}
			}

			// Try compound number parsing on the full accumulated phrase
			if j > i+1 {
				phrase := strings.Join(words[i:j], " ")
				if num, ok := parseNumberPhrase(phrase); ok {
					b.WriteString(strconv.Itoa(num))
					i = j
					continue
				}
			}

			// Single word or compound failed — convert just this word
			if num, ok := parseNumberPhrase(word); ok {
				b.WriteString(strconv.Itoa(num))
			} else {
				b.WriteString(word)
			}
			i++
			continue
		}

		// 4. Pass through unchanged
		b.WriteString(word)
		i++
	}

	return b.String()
}

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

	// Strip thousands separators and normalize decimal delimiter using pre-compiled patterns.
	if decimalDelimiter == '.' {
		// Default: comma is thousands separator, dot is decimal (already correct)
		result = digitCommaDigit.ReplaceAllString(result, "${1}${2}")
	} else {
		// European: dot is thousands separator, comma is decimal
		result = digitDotDigit.ReplaceAllString(result, "${1}${2}")
		result = digitCommaDigit.ReplaceAllString(result, "${1}.${2}")
	}

	// Replace user variables first (highest priority), then built-in constants (lower priority)
	// This ensures user-provided variables override built-in constants
	result = replaceVariables(result, variables)
	result = replaceVariables(result, builtInConstants)

	// Handle % symbol before word scan (not a word boundary, needs character-level replacement)
	result = strings.ReplaceAll(result, "%", " * 0.01 * ")

	// Single-pass word scan: replaces keywords and converts number words in one pass
	words := strings.Fields(result)
	result = replaceWords(words)

	// Convert "point" used as a decimal separator (e.g. "ten point two" → "10.2").
	// This runs after number word conversion so digit context is established.
	result = pointBetweenDigitsPattern.ReplaceAllString(result, "${1}.${2}")

	result = whitespaceRegex.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
