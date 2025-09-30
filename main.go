package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/radther/nlcalc/internal/cleaner"
	"github.com/radther/nlcalc/internal/evaluator"
	"github.com/radther/nlcalc/internal/normalizer"
	"github.com/radther/nlcalc/internal/tokenizer"
)

func main() {
	verbose := flag.Bool("verbose", false, "Show detailed pipeline stages")
	flag.Parse()

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

	normalized := normalizer.Normalize(input, nil)
	cleaned := cleaner.Clean(normalized)
	tokens, err := tokenizer.Tokenize(cleaned)

	if *verbose {
		fmt.Printf("Original: %s\n", input)
		fmt.Printf("Normalized: %s\n", normalized)
		fmt.Printf("Cleaned: %s\n", cleaned)
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
