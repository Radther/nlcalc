# nlcalc

A Go CLI tool and package that parses natural language math expressions and returns calculated results.

## Installation

Install using `go install`:

```bash
go install github.com/radther/nlcalc@latest
```

Or build from source:

```bash
git clone https://github.com/radther/nlcalc.git
cd nlcalc
go build -o nlcalc
```

## CLI Usage

### Basic Examples

```bash
# Written numbers and words
nlcalc "ten plus fifteen"
# Output: 25

# Symbols
nlcalc "5 + 10"
# Output: 15

# Mixed format
nlcalc "twenty * 3"
# Output: 60

# Percentages
nlcalc "20% of 100"
# Output: 20

# Order of operations (PEMDAS/BODMAS)
nlcalc "10 + 10 * 4"
# Output: 50

nlcalc "(10 + 10) * 4"
# Output: 80
```

### Input from stdin

```bash
echo "5 + 10" | nlcalc
# Output: 15
```

### Verbose Mode

Use the `--verbose` flag to see the processing pipeline stages:

```bash
nlcalc --verbose "what is ten plus fifteen"
# Output:
# Original: what is ten plus fifteen
# Normalized: what is 10 + 15
# Cleaned: 10 + 15
# Tokens: [NUMBER(10), OPERATOR(+), NUMBER(15)]
# Result: 25
```

The verbose output shows:
- **Original**: Your input expression
- **Normalized**: Written numbers and operation words converted to symbols
- **Cleaned**: Filler words and articles removed
- **Tokens**: Mathematical tokens identified by the parser
- **Result**: Final calculated value

## Using as a Package

Import nlcalc in your Go applications:

```go
package main

import (
    "fmt"
    "log"

    "github.com/radther/nlcalc/pkg/nlcalc"
)

func main() {
    result, err := nlcalc.Parse("ten plus fifteen", nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Result: %g\n", result)
    // Output: Result: 25
}
```

### Add to your project

```bash
go get github.com/radther/nlcalc
```

The package exposes a single function:

```go
func Parse(input string, variables map[string]float64) (float64, error)
```

The `variables` parameter is optional (can be nil) and allows you to pass custom named values that will be replaced in the expression before evaluation.

## Features

### Written Numbers

Parses written numbers from zero up to billions with full support for compound numbers:

```bash
nlcalc "five"
# Output: 5

nlcalc "twenty three"
# Output: 23

nlcalc "one hundred and forty two"
# Output: 142

nlcalc "five thousand seven hundred"
# Output: 5700

nlcalc "two million three hundred thousand"
# Output: 2300000

nlcalc "five billion three hundred million and twenty"
# Output: 5300000020
```

### Basic Arithmetic Operations

Supports both operation words and mathematical symbols:

| Words | Symbol | Example |
|-------|--------|---------|
| plus, add, added to | + | "ten plus five" → 15 |
| minus, subtract | - | "twenty minus eight" → 12 |
| times, multiply | * | "seven times three" → 21 |
| divide, divided by | / | "hundred divided by four" → 25 |

```bash
nlcalc "five plus ten"
# Output: 15

nlcalc "twenty minus eight"
# Output: 12

nlcalc "seven times three"
# Output: 21

nlcalc "hundred divided by four"
# Output: 25
```

### Power Operations

Supports exponentiation with multiple notations:

```bash
# Symbol notation
nlcalc "2 ^ 3"
# Output: 8

# Natural language
nlcalc "2 power 3"
# Output: 8

nlcalc "2 raised to 3"
# Output: 8

nlcalc "2 raised to the power of 3"
# Output: 8

# Shortcuts
nlcalc "5 squared"
# Output: 25

nlcalc "3 cubed"
# Output: 27

# Right-associative (evaluates right-to-left)
nlcalc "2 ^ 3 ^ 2"
# Output: 512  (2^(3^2) = 2^9)
```

### Mathematical Functions

Built-in support for common mathematical functions:

```bash
# Square root
nlcalc "sqrt 16"
# Output: 4

nlcalc "square root of 25"
# Output: 5

# Absolute value
nlcalc "abs -10"
# Output: 10

nlcalc "absolute value of minus fifteen"
# Output: 15

# Logarithms
nlcalc "log 100"
# Output: 2  (base 10)

nlcalc "logarithm of 1000"
# Output: 3

nlcalc "natural logarithm of 2.718281828459045"
# Output: 1

# Trigonometric functions
nlcalc "sine of 0"
# Output: 0

nlcalc "cosine of 0"
# Output: 1

# Nested functions
nlcalc "sqrt sqrt 16"
# Output: 2  (sqrt(sqrt(16)) = sqrt(4) = 2)

nlcalc "log sqrt 10000"
# Output: 2  (log(sqrt(10000)) = log(100) = 2)
```

Supported functions:
- `sqrt` / `square root` - Square root
- `abs` / `absolute` - Absolute value
- `log` / `logarithm` - Base 10 logarithm
- `ln` / `natural logarithm` - Natural logarithm
- `sin` / `sine` - Sine (radians)
- `cos` / `cosine` - Cosine (radians)
- `tan` / `tangent` - Tangent (radians)

### Percentages

Handles percentage calculations with multiple notations:

```bash
nlcalc "20% of 100"
# Output: 20

nlcalc "50 percent of 80"
# Output: 40

nlcalc "15 percentage of 200"
# Output: 30
```

### Variables

Pass custom named values to use in expressions:

```go
// Simple variable
vars := map[string]float64{"x": 5}
result, err := nlcalc.Parse("x plus ten", vars)
// Result: 15

// Multiple variables
vars := map[string]float64{"price": 100, "tax": 15}
result, err := nlcalc.Parse("price plus tax", vars)
// Result: 115

// Variables with order of operations
vars := map[string]float64{"price": 100, "tax": 0.15}
result, err := nlcalc.Parse("price plus price times tax", vars)
// Result: 115  (100 + (100 * 0.15))

// Variables with percentages
vars := map[string]float64{"discount": 20, "price": 100}
result, err := nlcalc.Parse("discount% of price", vars)
// Result: 20
```

Variable features:
- Case-insensitive matching
- Word boundary matching (prevents partial replacements)
- Longer variable names replaced first
- Works with all operations and functions

### Order of Operations

Follows standard PEMDAS/BODMAS rules:

1. **Parentheses** - `(` `)`
2. **Exponentiation** - `^` (right-associative)
3. **Functions** - `sqrt`, `abs`, `log`, etc. (right-associative)
4. **Multiplication/Division** - `*` `/` (left to right)
5. **Addition/Subtraction** - `+` `-` (left to right)

```bash
nlcalc "2 + 3 * 4"
# Output: 14  (not 20)

nlcalc "(2 + 3) * 4"
# Output: 20

nlcalc "2 * 3 ^ 2"
# Output: 18  (2 * (3^2) = 2 * 9)

nlcalc "10 + sqrt 16 * 2"
# Output: 18  (10 + (sqrt(16) * 2) = 10 + 8)
```

### Negative Numbers

Supports unary minus and plus operators:

```bash
# Note: Use -- before expressions starting with - to prevent flag parsing
nlcalc -- "-10"
# Output: -10

nlcalc "minus twenty"
# Output: -20

nlcalc -- "-10 + 5"
# Output: -5

nlcalc "10 + -5"
# Output: 5

nlcalc -- "-(10 + 5)"
# Output: -15

nlcalc -- "-(-5)"
# Output: 5

nlcalc -- "-sqrt 16"
# Output: -4
```

### Parentheses

Use parentheses to override order of operations and group expressions:

```bash
nlcalc "(10 + 5) * 2"
# Output: 30

nlcalc -- "-(2 + 3) * 4"
# Output: -20

nlcalc "sqrt((8 * 3) - 8)"
# Output: 4  (sqrt(16))
```

### Natural Language Support

Mix and match formats for readability:

```bash
nlcalc "what is ten plus fifteen"
# Output: 25

nlcalc "calculate five times twenty"
# Output: 100

nlcalc "twenty percent of one hundred"
# Output: 20

nlcalc "square root of sixteen plus two"
# Output: 6  (sqrt(16) + 2)
```

The parser automatically removes filler words like "what is", "calculate", "the", etc.

## Error Handling

The parser returns clear error messages for:

```bash
# Invalid token sequences
nlcalc "10 + + 5"
# Error: Tokenization error: ...

# Incomplete expressions
nlcalc "10 +"
# Error: Tokenization error: ...

# Division by zero
nlcalc "10 / 0"
# Error: Evaluation error: division by zero

# Invalid function arguments
nlcalc "sqrt -4"
# Error: Evaluation error: square root of negative number

nlcalc "log 0"
# Error: Evaluation error: logarithm of non-positive number
```

## Architecture

The parser follows a multi-stage pipeline:

1. **Normalization** - Replaces variables, converts written numbers to digits, operation words to symbols, handles percentages and power notation
2. **Cleaning** - Removes filler words and unnecessary whitespace
3. **Tokenization** - Parses into mathematical tokens (NUMBER, OPERATOR, FUNCTION, etc.)
4. **Evaluation** - Calculates final result following order of operations

The `--verbose` flag reveals this pipeline for debugging and understanding how expressions are processed.

## Examples

### Simple Calculations

```bash
nlcalc "five plus ten"          # 15
nlcalc "twenty minus eight"      # 12
nlcalc "seven times three"       # 21
nlcalc "hundred divided by four" # 25
```

### Complex Expressions

```bash
nlcalc "five squared plus three cubed"
# Output: 52  ((5^2) + (3^3) = 25 + 27)

nlcalc "square root of sixteen times five"
# Output: 20  (sqrt(16) * 5 = 4 * 5)

nlcalc "(ten plus five) times two minus three"
# Output: 27  ((10 + 5) * 2 - 3 = 30 - 3)

nlcalc "absolute value of minus ten plus five"
# Output: 15  (abs(-10) + 5 = 10 + 5)
```

### Real-World Applications

```bash
# Sales tax calculation
nlcalc "15 percent of 299.99"
# Output: 44.9985

# Compound interest (simplified)
nlcalc "1000 times 1.05 ^ 10"
# Output: 1628.89...

# Pythagorean theorem
nlcalc "sqrt(3 ^ 2 + 4 ^ 2)"
# Output: 5

# Temperature conversion approximation
nlcalc "32 plus 9 divided by 5 times 25"
# Output: 77  (Celsius to Fahrenheit)
```

### Using Variables in Code

```go
// Shopping cart total
vars := map[string]float64{
    "quantity": 3,
    "price":    25.50,
    "shipping": 10,
}
result, _ := nlcalc.Parse("quantity times price plus shipping", vars)
// Result: 86.5

// Loan payment calculation (simplified)
vars := map[string]float64{
    "principal": 100000,
    "rate":      0.05,
    "years":     30,
}
result, _ := nlcalc.Parse("principal times rate", vars)
// Result: 5000
```

## Performance

nlcalc is optimized for both simple and complex expressions:

- Pre-compiled regex patterns for efficiency
- Single-pass normalization
- Efficient token-based evaluation
- Right-associative evaluation for functions and exponentiation

See the benchmark tests in the repository for detailed performance metrics.

## License

MIT
