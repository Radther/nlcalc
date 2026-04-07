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

// builtInConstants defines mathematical constants available as built-in variables.
// These are lower priority than user-provided variables.
var builtInConstants = map[string]float64{
	"pi":  3.141592653589793,
	"e":   2.718281828459045,
	"tau": 6.283185307179586,
	"phi": 1.618033988749895,
}

// builtInConstantStrs holds pre-formatted string values of built-in constants.
var builtInConstantStrs = func() map[string]string {
	m := make(map[string]string, len(builtInConstants))
	for name, val := range builtInConstants {
		m[name] = strconv.FormatFloat(val, 'f', -1, 64)
	}
	return m
}()

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
		{words: []string{"to", "the", "power", "of"}, replacement: "^"},
	},
	"natural": {
		{words: []string{"logarithm"}, replacement: "ln"},
	},
	"square": {
		{words: []string{"root"}, replacement: "sqrt"},
	},
	"divided": {
		{words: []string{"by"}, replacement: "/"},
	},
	"added": {
		{words: []string{"to"}, replacement: "+"},
	},
}

// wordReplacements maps single keywords to their mathematical symbol replacements.
// Values are trimmed — the word scanner handles inter-word spacing.
var wordReplacements = map[string]string{
	"percentage": "* 0.01 *",
	"percent":    "* 0.01 *",
	"squareroot": "sqrt",
	"squared":    "^ 2",
	"cubed":      "^ 3",
	"raised":     "^",
	"power":      "^",
	"absolute":   "abs",
	"logarithm":  "log",
	"sine":       "sin",
	"cosine":     "cos",
	"tangent":    "tan",
	"plus":       "+",
	"add":        "+",
	"minus":      "-",
	"subtract":   "-",
	"times":      "*",
	"multiply":   "*",
	"divide":     "/",
}

// pointBetweenDigitsPattern matches the word "point" used as a decimal separator between digits.
var pointBetweenDigitsPattern = regexp.MustCompile(`(\d)\s*point\s*(\d)`)

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// parseNumberWords converts a slice of number/scale words into an integer.
// Words must already be lowercase. Returns false if the phrase is invalid.
func parseNumberWords(words []string) (int, bool) {
	if len(words) == 0 {
		return 0, false
	}

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
				current *= 100
			} else if scale >= 1000 {
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
			if i == 0 {
				return false
			}

			prevWord := words[i-1]
			if _, isScale := scaleWords[prevWord]; !isScale {
				return hasRecentScaleContext(words, i)
			}
		}
	}
	return true
}

// hasRecentScaleContext checks if there's a scale word in recent context
func hasRecentScaleContext(words []string, andPos int) bool {
	for i := andPos - 1; i >= 0; i-- {
		if _, isScale := scaleWords[words[i]]; isScale {
			return true
		}
		if words[i] == "and" || i < andPos-3 {
			break
		}
	}
	return false
}

// replaceUserVariables replaces user-provided variable names in the input string.
// Uses strings.Index for matching, which correctly handles multi-word variable names
// like "my var" and names adjacent to operators like "x+5".
// Variables are sorted by length (descending) to replace longest first.
func replaceUserVariables(input string, variables map[string]float64) string {
	if len(variables) == 0 {
		return input
	}

	// Build sorted name list and value map
	type varEntry struct {
		name  string
		value string
	}
	entries := make([]varEntry, 0, len(variables))
	for name, val := range variables {
		entries = append(entries, varEntry{
			name:  strings.ToLower(name),
			value: strconv.FormatFloat(val, 'f', -1, 64),
		})
	}
	// Sort by name length descending — longest match first
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && len(entries[j].name) > len(entries[j-1].name); j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}

	result := input
	for _, e := range entries {
		nameLen := len(e.name)

		// Lazy builder per variable
		var b strings.Builder
		changed := false
		lastWritten := 0

		i := 0
		for {
			idx := strings.Index(result[i:], e.name)
			if idx == -1 {
				break
			}
			pos := i + idx
			startOk := pos == 0 || !isWordChar(result[pos-1])
			endOk := pos+nameLen == len(result) || !isWordChar(result[pos+nameLen])
			if startOk && endOk {
				if !changed {
					b.Grow(len(result))
					changed = true
				}
				b.WriteString(result[lastWritten:pos])
				b.WriteString(e.value)
				lastWritten = pos + nameLen
				i = pos + nameLen
			} else {
				i = pos + 1
			}
		}

		if changed {
			b.WriteString(result[lastWritten:])
			result = b.String()
		}
	}

	return result
}

// replaceBuiltInsAndPercent performs a single O(N) pass through the string, replacing
// built-in constant names (pi, e, tau, phi) and the '%' symbol.
// Uses a lazy builder: returns the original string with zero allocations if nothing matches.
func replaceBuiltInsAndPercent(input string) string {
	var b strings.Builder
	lastWritten := 0
	changed := false

	i := 0
	for i < len(input) {
		c := input[i]

		if c == '%' {
			if !changed {
				b.Grow(len(input) + 32)
				changed = true
			}
			b.WriteString(input[lastWritten:i])
			b.WriteString(" * 0.01 * ")
			i++
			lastWritten = i
			continue
		}

		if isAlpha(c) || c == '_' {
			start := i
			i++
			for i < len(input) && isWordChar(input[i]) {
				i++
			}
			word := input[start:i]
			if replacement, ok := builtInConstantStrs[word]; ok {
				if !changed {
					b.Grow(len(input) + 32)
					changed = true
				}
				b.WriteString(input[lastWritten:start])
				b.WriteString(replacement)
				lastWritten = i
			}
			continue
		}

		if isDigit(c) {
			i++
			for i < len(input) && isWordChar(input[i]) {
				i++
			}
			continue
		}

		i++
	}

	if !changed {
		return input
	}
	b.WriteString(input[lastWritten:])
	return b.String()
}

// stripThousandsSep removes a thousands separator character that appears between digits.
// Returns the original string with zero allocations if no separator is found.
func stripThousandsSep(input string, sep byte) string {
	if !strings.ContainsRune(input, rune(sep)) {
		return input
	}
	var b strings.Builder
	b.Grow(len(input))
	for i := 0; i < len(input); i++ {
		if input[i] == sep && i > 0 && i < len(input)-1 && isDigit(input[i-1]) && isDigit(input[i+1]) {
			continue
		}
		b.WriteByte(input[i])
	}
	return b.String()
}

// normalizeDecimalDelim replaces a non-dot decimal delimiter with '.' between digits.
func normalizeDecimalDelim(input string, delim byte) string {
	if !strings.ContainsRune(input, rune(delim)) {
		return input
	}
	var b strings.Builder
	b.Grow(len(input))
	for i := 0; i < len(input); i++ {
		if input[i] == delim && i > 0 && i < len(input)-1 && isDigit(input[i-1]) && isDigit(input[i+1]) {
			b.WriteByte('.')
			continue
		}
		b.WriteByte(input[i])
	}
	return b.String()
}

// replaceWords performs a single-pass scan over whitespace-delimited words, replacing
// keywords with mathematical symbols and converting number words to digits.
func replaceWords(words []string) string {
	if len(words) == 0 {
		return ""
	}

	var b strings.Builder
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
				if num, ok := parseNumberWords(words[i:j]); ok {
					b.WriteString(strconv.Itoa(num))
					i = j
					continue
				}
			}

			// Single word or compound failed — convert just this word
			if num, ok := parseNumberWords(words[i : i+1]); ok {
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
func Normalize(input string, variables map[string]float64, decimalDelimiter rune) string {
	if decimalDelimiter == 0 {
		decimalDelimiter = '.'
	}

	result := strings.ToLower(input)

	// Strip thousands separators and normalize decimal delimiter.
	if decimalDelimiter == '.' {
		result = stripThousandsSep(result, ',')
	} else {
		result = stripThousandsSep(result, '.')
		result = normalizeDecimalDelim(result, byte(decimalDelimiter))
	}

	// Replace user variables first (highest priority, supports multi-word names),
	// then built-in constants and '%' symbol in a single pass.
	result = replaceUserVariables(result, variables)
	result = replaceBuiltInsAndPercent(result)

	// Single-pass word scan: replaces keywords and converts number words.
	words := strings.Fields(result)
	result = replaceWords(words)

	// Convert "point" as decimal separator (e.g. "10 point 2" → "10.2").
	if strings.Contains(result, "point") {
		result = pointBetweenDigitsPattern.ReplaceAllString(result, "${1}.${2}")
	}

	return strings.TrimSpace(result)
}
