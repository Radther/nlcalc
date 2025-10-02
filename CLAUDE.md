# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

nlcalc is a Go module and CLI tool that parses natural language math expressions and returns calculated results. It handles written numbers ("ten"), compound numbers ("five billion"), mathematical operations in text ("plus"), symbols ("5 + 10"), percentage calculations ("20% of 100"), and custom variables.

## Build and Test Commands

- Build CLI: `go build -o nlcalc main.go` or `go build -o nlcalc`
- Run CLI: `./nlcalc "ten plus fifteen"` or `echo "5 + 10" | ./nlcalc`
- Run with verbose output: `./nlcalc --verbose "expression"`
- Test all: `go test ./...`
- Test specific package: `go test ./internal/normalizer` (or any other package)
- Run with go: `go run main.go "expression"` or `go run . "expression"`

## Architecture

The parser follows a multi-stage pipeline architecture:

1. **Normalization** (`internal/normalizer`) - Replaces custom variables with their numeric values, converts written numbers ("ten" → "10", "five billion" → "5000000000") including compound numbers with "and" ("one hundred and twenty" → "120"), transforms operation words ("plus" → "+") to symbols, and handles percentage notation ("% of" → "* 0.01 *"). Supports numbers from zero up to billions.

2. **Cleaning** (`internal/cleaner`) - Removes filler words, articles, and unnecessary whitespace that don't contribute to the mathematical expression.

3. **Tokenization** (`internal/tokenizer`) - Parses cleaned string into mathematical tokens (NUMBER, OPERATOR, PARENTHESIS). Validates token sequences.

4. **Evaluation** (`internal/evaluator`) - Calculates final result following PEMDAS/BODMAS order of operations.

The public API is exposed through `pkg/nlcalc/parser.go` with a single `Parse(input string, variables map[string]float64) (float64, error)` function that runs the complete pipeline. The optional `variables` parameter (can be nil) allows passing custom named values that will be replaced in the expression before evaluation.

The CLI (`main.go`) accepts an optional `--verbose` flag that provides progressive output showing each pipeline stage (Original, Normalized, Cleaned, Tokens, Result) for debugging and transparency.

## Order of Operations

The evaluator implements standard PEMDAS/BODMAS:
- "10 + 10 * 4" evaluates to 50 (not 80)
- "(10 + 10) * 4" evaluates to 80

## Features

### Variables Support
The Parse function accepts an optional variables map for custom named values:
```go
vars := map[string]float64{"price": 100, "tax": 0.15}
result, err := nlcalc.Parse("price plus price times tax", vars)
// Result: 115
```

Key characteristics:
- Variable names are case-insensitive
- Supports complex expressions with multiple variables
- Uses word boundaries for matching (e.g., "x" won't partially match "xx")
- Longer variable names are replaced first to prevent partial replacements
- Variables can be used with percentages: "discount% of price"

### Number Parsing
Supports written numbers from zero up to billions:
- Basic numbers: "zero" through "ninety nine"
- Compound numbers: "one hundred and twenty seven", "two thousand five hundred"
- Large numbers: "five billion three hundred million and twenty"
- Proper "and" validation in compound numbers

### CLI Verbose Mode
The `--verbose` flag shows detailed pipeline processing:
```bash
./nlcalc --verbose "ten plus fifteen"
# Original: ten plus fifteen
# Normalized: 10 + 15
# Cleaned: 10 + 15
# Tokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]
# Result: 25
```

## Error Handling

The parser returns errors for:
- Invalid token sequences ("10 + + 5")
- Incomplete expressions ("10 +")
- Division by zero
- Unrecognized input after cleaning