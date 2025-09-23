package main

import (
	"bufio"
	"fmt"
	"mathparser/internal/cleaner"
	"mathparser/internal/evaluator"
	"mathparser/internal/normalizer"
	"mathparser/internal/tokenizer"
	"os"
	"strings"
)

func main() {
	var input string

	if len(os.Args) > 1 {
		input = strings.Join(os.Args[1:], " ")
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

	normalized := normalizer.Normalize(input)
	cleaned := cleaner.Clean(normalized)
	tokens, err := tokenizer.Tokenize(cleaned)

	fmt.Printf("Original: %s\n", input)
	fmt.Printf("Normalized: %s\n", normalized)
	fmt.Printf("Cleaned: %s\n", cleaned)

	if err != nil {
		fmt.Printf("Tokenization error: %s\n", err)
		return
	}

	tokenStrings := make([]string, len(tokens))
	for i, token := range tokens {
		tokenStrings[i] = token.String()
	}
	fmt.Printf("Tokens: [%s]\n", strings.Join(tokenStrings, ", "))

	result, err := evaluator.Evaluate(tokens)
	if err != nil {
		fmt.Printf("Evaluation error: %s\n", err)
		return
	}

	fmt.Printf("Result: %g\n", result)
}