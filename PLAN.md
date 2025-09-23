# Inline Math Expression Parser - Implementation Plan

## Project Overview

A Go module and CLI tool that parses natural language math expressions and returns calculated results. The tool should handle various formats including written numbers, mathematical operations, and percentage calculations.

## Examples

| Input | Expected Output |
|-------|----------------|
| "ten plus fifteen" | 25 |
| "5 + 10" | 15 |
| "20% of 100" | 20 |
| "10 plus 10 times 4" | 50 |

## Architecture Overview

The parsing system will follow a multi-stage pipeline:

1. **Basic CLI Setup** - Create CLI with progressive output capability
2. **Text Normalization** - Convert written numbers and operations to symbols
3. **Token Cleaning** - Remove unnecessary filler words
4. **Tokenization** - Convert cleaned string to mathematical tokens
5. **Expression Evaluation** - Calculate the final result

## Commands

- Build: `go build -o mathparser cmd/main.go`
- Test full suite: `go test ./...`

## Stages of Implementation

### Stage 1: Basic CLI Setup

**Objective**: Create a basic CLI that accepts input and returns "not implemented yet" while setting up the foundation for progressive output.

**Acceptance Criteria**:
- Accept expression as command-line argument
- Accept expression via stdin for interactive mode
- Return "not implemented yet" for all inputs
- Set up module structure and basic pipeline functions

**Progressive Output Example**:
```
Input: "what is 10 plus fifteen"
Output: "not implemented yet"
```

**Test Cases**:
```bash
./mathparser "ten plus fifteen"  # Output: not implemented yet
./mathparser "5 + 10"           # Output: not implemented yet
echo "invalid" | ./mathparser   # Output: not implemented yet
```

### Stage 2: Text Normalization

**Objective**: Convert written numbers and mathematical operations to their numeric/symbolic equivalents.

**Acceptance Criteria**:
- Convert written numbers (one, two, ten, fifteen, etc.) to digits
- Convert operation words (plus, minus, times, divided by, etc.) to symbols (+, -, *, /)
- Handle percentage notation ("% of" → multiplication by 0.01)
- Preserve existing numeric values and symbols
- CLI now shows original input and normalized result

**Progressive Output Example**:
```
Input: "what is 10 plus fifteen"
Output:
"""
Original: "what is 10 plus fifteen"
Normalized: "what is 10 + 15"
"""
```

**Test Cases**:
```go
// Number conversion
"ten" → "10"
"fifteen" → "15"
"twenty-five" → "25"
"one hundred" → "100"

// Operation conversion
"plus" → "+"
"minus" → "-"
"times" → "*"
"divided by" → "/"
"percent of" → "* 0.01 *"

// Mixed examples
"ten plus fifteen" → "10 + 15"
"5 plus ten" → "5 + 10"
"twenty percent of one hundred" → "20 * 0.01 * 100"
```

### Stage 3: Token Cleaning

**Objective**: Remove unnecessary words that don't contribute to the mathematical expression.

**Acceptance Criteria**:
- Remove articles (a, an, the)
- Remove prepositions that aren't part of operations (of, by, etc.)
- Remove extra whitespace
- Preserve mathematical structure
- CLI now shows original, normalized, and cleaned results

**Progressive Output Example**:
```
Input: "what is 10 plus fifteen"
Output:
"""
Original: "what is 10 plus fifteen"
Normalized: "what is 10 + 15"
Cleaned: "10 + 15"
"""
```

**Test Cases**:
```go
"the sum of ten and fifteen" → "10 + 15"
"a total of five plus ten" → "5 + 10"
"multiply the number ten by four" → "10 * 4"
```

### Stage 4: Tokenization

**Objective**: Parse the cleaned string into mathematical tokens (numbers, operators, parentheses).

**Acceptance Criteria**:
- Identify and classify all tokens (NUMBER, OPERATOR, PARENTHESIS)
- Handle decimal numbers
- Validate token sequence for mathematical correctness
- Return error for invalid expressions
- CLI now shows original, normalized, cleaned, and tokenized results

**Progressive Output Example**:
```
Input: "what is 10 plus fifteen"
Output:
"""
Original: "what is 10 plus fifteen"
Normalized: "what is 10 + 15"
Cleaned: "10 + 15"
Tokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]
"""
```

**Test Cases**:
```go
"10 + 15" → [NUMBER(10), OPERATOR(+), NUMBER(15)]
"5 * 3.14" → [NUMBER(5), OPERATOR(*), NUMBER(3.14)]
"10 + + 5" → ERROR (invalid operator sequence)
"10 +" → ERROR (incomplete expression)
```

### Stage 5: Expression Evaluation

**Objective**: Calculate the result following proper order of operations.

**Acceptance Criteria**:
- Implement PEMDAS/BODMAS order of operations
- Handle parentheses correctly
- Support basic arithmetic operations (+, -, *, /)
- Return floating-point results when appropriate
- Handle division by zero gracefully
- CLI now shows complete pipeline with final result

**Progressive Output Example**:
```
Input: "what is 10 plus fifteen"
Output:
"""
Original: "what is 10 plus fifteen"
Normalized: "what is 10 + 15"
Cleaned: "10 + 15"
Tokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]
Result: 25
"""
```

**Test Cases**:
```go
"10 + 15" → 25.0
"10 + 10 * 4" → 50.0 (not 80.0)
"(10 + 10) * 4" → 80.0
"20 * 0.01 * 100" → 20.0
"10 / 0" → ERROR (division by zero)
```


## Module Structure

```
mathparser/
├── cmd/
│   └── main.go              # CLI entry point
├── internal/
│   ├── normalizer/
│   │   └── normalizer.go    # Text normalization
│   ├── cleaner/
│   │   └── cleaner.go       # Token cleaning
│   ├── tokenizer/
│   │   └── tokenizer.go     # String tokenization
│   └── evaluator/
│       └── evaluator.go     # Expression evaluation
├── pkg/
│   └── mathparser/
│       └── parser.go        # Public API
├── go.mod
├── go.sum
└── README.md
```

## Error Handling

The module should handle and return appropriate errors for:
- Invalid expressions
- Division by zero
- Unrecognized words/operations
- Malformed input
- Empty input

## Future Enhancements

- Support for more complex operations (square root, exponents)
- Variable support
- Function calls (sin, cos, log, etc.)
- Multiple expression support
- Configuration for different number systems (Roman numerals, etc.)
