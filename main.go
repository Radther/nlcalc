package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/radther/nlcalc/internal/evaluator"
	"github.com/radther/nlcalc/internal/normalizer"
	"github.com/radther/nlcalc/internal/tokenizer"
)

func main() {
	verbose := flag.Bool("verbose", false, "Show detailed pipeline stages")
	decimalDelimiterFlag := flag.String("decimal-delimiter", ".", "Decimal delimiter character ('.' or ',')")
	flag.Parse()

	// Parse and validate the decimal delimiter flag.
	if len(*decimalDelimiterFlag) != 1 {
		fmt.Fprintf(os.Stderr, "Error: --decimal-delimiter must be a single character, got %q\n", *decimalDelimiterFlag)
		os.Exit(1)
	}
	decimalDelimiter := rune((*decimalDelimiterFlag)[0])

	var input string

	if flag.NArg() > 0 {
		input = strings.Join(flag.Args(), " ")
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	}

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("not implemented yet")
		return
	}

	normalized := normalizer.Normalize(input, nil, decimalDelimiter)
	tokens, err := tokenizer.Tokenize(normalized)

	if *verbose {
		fmt.Printf("Original: %s\n", input)
		fmt.Printf("Normalized: %s\n", normalized)
		// Reconstruct cleaned view from tokens
		tokenValues := make([]string, len(tokens))
		for i, token := range tokens {
			tokenValues[i] = token.Value
		}
		fmt.Printf("Cleaned: %s\n", strings.Join(tokenValues, " "))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Tokenization error: %s\n", err)
		os.Exit(1)
	}

	if *verbose {
		tokenStrings := make([]string, len(tokens))
		for i, token := range tokens {
			tokenStrings[i] = token.String()
		}
		fmt.Printf("Tokens: [%s]\n", strings.Join(tokenStrings, ", "))
	}

	result, err := evaluator.Evaluate(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Evaluation error: %s\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Result: %g\n", result)
	} else {
		fmt.Printf("%g\n", result)
	}
}
