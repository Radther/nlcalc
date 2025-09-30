// Package normalizer converts written numbers and operation words into mathematical symbols.
// It handles number words (e.g., "ten", "twenty five"), operation words (e.g., "plus", "minus"),
// and percentage expressions (e.g., "20% of", "percent of") as part of the nlcalc parsing pipeline.
package normalizer

import (
	"regexp"
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

// Compile regex patterns once at package level
var individualWordPattern = regexp.MustCompile(`\b(?:one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|thirteen|fourteen|fifteen|sixteen|seventeen|eighteen|nineteen|twenty|thirty|forty|fifty|sixty|seventy|eighty|ninety|zero)\b`)

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

// Normalize converts a raw input string into a normalized mathematical expression.
// It transforms written numbers ("ten" → "10"), operation words ("plus" → "+"),
// and percentage notations ("20% of" → "* 0.01 *") into symbolic form.
// Returns a normalized string ready for cleaning and tokenization.
func Normalize(input string) string {
	result := strings.ToLower(input)

	result = regexp.MustCompile(`\bpercent of\b`).ReplaceAllString(result, " * 0.01 * ")
	result = regexp.MustCompile(`\bpercentage of\b`).ReplaceAllString(result, " * 0.01 * ")
	result = regexp.MustCompile(`% of\b`).ReplaceAllString(result, " * 0.01 * ")
	result = regexp.MustCompile(`\bpercent\b`).ReplaceAllString(result, " * 0.01 * ")
	result = regexp.MustCompile(`\bpercentage\b`).ReplaceAllString(result, " * 0.01 * ")
	result = regexp.MustCompile(`%`).ReplaceAllString(result, " * 0.01 * ")

	operationMap := map[string]string{
		"plus":       " + ",
		"add":        " + ",
		"added to":   " + ",
		"minus":      " - ",
		"subtract":   " - ",
		"times":      " * ",
		"multiply":   " * ",
		"divided by": " / ",
		"divide":     " / ",
	}

	for word, symbol := range operationMap {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		result = re.ReplaceAllString(result, symbol)
	}

	// Single number pattern - handles both compound and individual numbers
	numberPattern := regexp.MustCompile(`\b(?:(?:one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|thirteen|fourteen|fifteen|sixteen|seventeen|eighteen|nineteen|twenty|thirty|forty|fifty|sixty|seventy|eighty|ninety|hundred|thousand|million|billion|zero|and)\s*)+\b`)

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

	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
