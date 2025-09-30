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
go build -o nlcalc cmd/main.go
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
    result, err := nlcalc.Parse("ten plus fifteen")
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
func Parse(input string) (float64, error)
```

It handles the complete parsing pipeline automatically and returns the calculated result or an error.

## Features

- **Natural language**: Parses written numbers (one through nine hundred ninety-nine)
- **Operation words**: Supports "plus", "minus", "times", "divided by"
- **Symbols**: Standard math operators (+, -, *, /)
- **Percentages**: "X% of Y" notation
- **Order of operations**: Follows PEMDAS/BODMAS rules
- **Parentheses**: Supports grouping with ()
- **Error handling**: Clear error messages for invalid expressions

## License

MIT
