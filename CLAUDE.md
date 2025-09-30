# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

nlcalc is a Go module and CLI tool that parses natural language math expressions and returns calculated results. It handles written numbers ("ten"), mathematical operations in text ("plus"), symbols ("5 + 10"), and percentage calculations ("20% of 100").

## Build and Test Commands

- Build CLI: `go build -o nlcalc cmd/main.go`
- Run CLI: `./nlcalc "ten plus fifteen"` or `echo "5 + 10" | ./nlcalc`
- Test all: `go test ./...`
- Test specific package: `go test ./internal/normalizer` (or any other package)
- Run with go: `go run cmd/main.go "expression"`

## Architecture

The parser follows a multi-stage pipeline architecture:

1. **Normalization** (`internal/normalizer`) - Converts written numbers ("ten" → "10") and operation words ("plus" → "+") to symbols. Handles percentage notation ("% of" → "* 0.01 *").

2. **Cleaning** (`internal/cleaner`) - Removes filler words, articles, and unnecessary whitespace that don't contribute to the mathematical expression.

3. **Tokenization** (`internal/tokenizer`) - Parses cleaned string into mathematical tokens (NUMBER, OPERATOR, PARENTHESIS). Validates token sequences.

4. **Evaluation** (`internal/evaluator`) - Calculates final result following PEMDAS/BODMAS order of operations.

The public API is exposed through `pkg/nlcalc/parser.go` with a single `Parse(input string) (float64, error)` function that runs the complete pipeline.

The CLI (`cmd/main.go`) provides progressive output showing each pipeline stage for debugging and transparency.

## Order of Operations

The evaluator implements standard PEMDAS/BODMAS:
- "10 + 10 * 4" evaluates to 50 (not 80)
- "(10 + 10) * 4" evaluates to 80

## Error Handling

The parser returns errors for:
- Invalid token sequences ("10 + + 5")
- Incomplete expressions ("10 +")
- Division by zero
- Unrecognized input after cleaning