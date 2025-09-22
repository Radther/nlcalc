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

func parseNumberPhrase(phrase string) (int, bool) {
	words := strings.Fields(strings.ToLower(phrase))
	if len(words) == 0 {
		return 0, false
	}

	total := 0
	current := 0

	for _, word := range words {
		if word == "and" {
			continue
		}

		if word == "hundred" {
			if current == 0 {
				current = 1
			}
			current *= 100
			total += current
			current = 0
		} else if word == "thousand" {
			if current == 0 {
				current = 1
			}
			total = (total + current) * 1000
			current = 0
		} else if val, exists := numberWords[word]; exists {
			current += val
		} else {
			return 0, false
		}
	}

	return total + current, true
}

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

	numberPattern := regexp.MustCompile(`\b(?:(?:one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|thirteen|fourteen|fifteen|sixteen|seventeen|eighteen|nineteen|twenty|thirty|forty|fifty|sixty|seventy|eighty|ninety|hundred|thousand|zero|and)\s*)+\b`)

	result = numberPattern.ReplaceAllStringFunc(result, func(match string) string {
		if num, ok := parseNumberPhrase(match); ok {
			return strconv.Itoa(num)
		}
		return match
	})

	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
