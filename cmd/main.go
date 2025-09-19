package main

import (
	"bufio"
	"fmt"
	"mathparser/internal/normalizer"
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
	
	fmt.Printf("Original: %s\n", input)
	fmt.Printf("Normalized: %s\n", normalized)
	fmt.Println("Result: not implemented yet")
}