// Package tokenizer converts cleaned mathematical expressions into tokens for evaluation.
package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TokenType represents the category of a mathematical token.
type TokenType string

const (
	// NUMBER represents numeric values (integers and decimals).
	NUMBER TokenType = "NUMBER"
	// OPERATOR represents mathematical operators (+, -, *, /, ^).
	OPERATOR TokenType = "OPERATOR"
	// UNARY_OPERATOR represents unary operators (unary minus, unary plus).
	UNARY_OPERATOR TokenType = "UNARY_OPERATOR"
	// PARENTHESIS represents grouping symbols (parentheses).
	PARENTHESIS TokenType = "PARENTHESIS"
	// FUNCTION represents mathematical functions (sqrt, abs, log, etc.).
	FUNCTION TokenType = "FUNCTION"
)

// Token represents a single parsed element of a mathematical expression.
type Token struct {
	Type  TokenType // The category of the token
	Value string    // The string value of the token
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Type, t.Value)
}

var (
	numberPattern      = regexp.MustCompile(`\d+\.?\d*`)
	operatorPattern    = regexp.MustCompile(`[+\-*/^]`)
	parenthesisPattern = regexp.MustCompile(`[()]`)
	whitespacePattern  = regexp.MustCompile(`\s+`)
	functionPattern    = regexp.MustCompile(`^(sqrt|abs|log|ln|sin|cos|tan)`)
)

// isUnaryContext determines if the current position represents a unary operator context.
// An operator (+/-) is unary rather than binary when it appears:
// - At the start of the expression (no previous tokens)
// - After another operator
// - After a function
// - After an opening parenthesis
func isUnaryContext(tokens []Token) bool {
	if len(tokens) == 0 {
		return true // At start of expression
	}

	lastToken := tokens[len(tokens)-1]
	return lastToken.Type == OPERATOR ||
	       lastToken.Type == UNARY_OPERATOR ||
	       lastToken.Type == FUNCTION ||
	       (lastToken.Type == PARENTHESIS && lastToken.Value == "(")
}

// Tokenize parses a cleaned mathematical expression string into a slice of tokens.
// It validates that numbers are valid, operators are recognized, and the token
// sequence is syntactically correct (e.g., no consecutive operators or numbers).
// Returns an error if the input is empty, contains unrecognized characters, or
// has an invalid token sequence.
func Tokenize(input string) ([]Token, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("empty input")
	}

	input = strings.TrimSpace(input)
	tokens := []Token{}
	i := 0

	for i < len(input) {
		if whitespacePattern.MatchString(string(input[i])) {
			i++
			continue
		}

		if funcMatch := functionPattern.FindString(input[i:]); funcMatch != "" {
			tokens = append(tokens, Token{Type: FUNCTION, Value: funcMatch})
			i += len(funcMatch)
			continue
		}

		if numberMatch := numberPattern.FindString(input[i:]); numberMatch != "" && strings.HasPrefix(input[i:], numberMatch) {
			if _, err := strconv.ParseFloat(numberMatch, 64); err != nil {
				return nil, fmt.Errorf("invalid number: %s", numberMatch)
			}
			tokens = append(tokens, Token{Type: NUMBER, Value: numberMatch})
			i += len(numberMatch)
			continue
		}

		if operatorMatch := operatorPattern.FindString(input[i:]); operatorMatch != "" && strings.HasPrefix(input[i:], operatorMatch) {
			// Check if +/- is unary or binary based on context
			if (operatorMatch == "-" || operatorMatch == "+") && isUnaryContext(tokens) {
				tokens = append(tokens, Token{Type: UNARY_OPERATOR, Value: operatorMatch})
			} else {
				tokens = append(tokens, Token{Type: OPERATOR, Value: operatorMatch})
			}
			i += len(operatorMatch)
			continue
		}

		if parenthesisMatch := parenthesisPattern.FindString(input[i:]); parenthesisMatch != "" && strings.HasPrefix(input[i:], parenthesisMatch) {
			tokens = append(tokens, Token{Type: PARENTHESIS, Value: parenthesisMatch})
			i += len(parenthesisMatch)
			continue
		}

		return nil, fmt.Errorf("unrecognized character: %c at position %d", input[i], i)
	}

	if err := validateTokenSequence(tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func validateTokenSequence(tokens []Token) error {
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens found")
	}

	for i, token := range tokens {
		switch token.Type {
		case OPERATOR:
			if i == 0 {
				return fmt.Errorf("expression cannot start with operator: %s", token.Value)
			}
			if i == len(tokens)-1 {
				return fmt.Errorf("expression cannot end with operator: %s", token.Value)
			}
			if i > 0 && tokens[i-1].Type == OPERATOR {
				return fmt.Errorf("consecutive operators not allowed: %s %s", tokens[i-1].Value, token.Value)
			}
			if i < len(tokens)-1 && tokens[i+1].Type == OPERATOR {
				return fmt.Errorf("consecutive operators not allowed: %s %s", token.Value, tokens[i+1].Value)
			}
		case UNARY_OPERATOR:
			// Unary operators must be at start, after operator, after unary operator, after function, or after "("
			if i > 0 {
				prev := tokens[i-1]
				if prev.Type != OPERATOR && prev.Type != UNARY_OPERATOR && prev.Type != FUNCTION &&
				   !(prev.Type == PARENTHESIS && prev.Value == "(") {
					return fmt.Errorf("unary operator in invalid position after %s", prev.Value)
				}
			}
			// Unary operators must be followed by a number, function, opening parenthesis, or another unary operator
			if i == len(tokens)-1 {
				return fmt.Errorf("expression cannot end with unary operator: %s", token.Value)
			}
			if i < len(tokens)-1 {
				next := tokens[i+1]
				if next.Type != NUMBER &&
				   next.Type != UNARY_OPERATOR &&
				   next.Type != FUNCTION &&
				   !(next.Type == PARENTHESIS && next.Value == "(") {
					return fmt.Errorf("unary operator must be followed by number, function, '(', or another unary operator, got %s", next.Value)
				}
			}
		case NUMBER:
			if i > 0 && tokens[i-1].Type == NUMBER {
				return fmt.Errorf("consecutive numbers not allowed: %s %s", tokens[i-1].Value, token.Value)
			}
		case FUNCTION:
			// Functions must be followed by NUMBER, another FUNCTION, UNARY_OPERATOR, or opening parenthesis
			if i == len(tokens)-1 {
				return fmt.Errorf("function must be followed by a value: %s", token.Value)
			}
			if i < len(tokens)-1 {
				next := tokens[i+1]
				if next.Type != NUMBER &&
				   next.Type != FUNCTION &&
				   next.Type != UNARY_OPERATOR &&
				   !(next.Type == PARENTHESIS && next.Value == "(") {
					return fmt.Errorf("function must be followed by number, function, unary operator, or '(', got %s", next.Value)
				}
			}
		}
	}

	return nil
}