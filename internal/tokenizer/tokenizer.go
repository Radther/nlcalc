package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TokenType string

const (
	NUMBER     TokenType = "NUMBER"
	OPERATOR   TokenType = "OPERATOR"
	PARENTHESIS TokenType = "PARENTHESIS"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Type, t.Value)
}

var (
	numberPattern     = regexp.MustCompile(`\d+\.?\d*`)
	operatorPattern   = regexp.MustCompile(`[+\-*/]`)
	parenthesisPattern = regexp.MustCompile(`[()]`)
	whitespacePattern = regexp.MustCompile(`\s+`)
)

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

		if numberMatch := numberPattern.FindString(input[i:]); numberMatch != "" && strings.HasPrefix(input[i:], numberMatch) {
			if _, err := strconv.ParseFloat(numberMatch, 64); err != nil {
				return nil, fmt.Errorf("invalid number: %s", numberMatch)
			}
			tokens = append(tokens, Token{Type: NUMBER, Value: numberMatch})
			i += len(numberMatch)
			continue
		}

		if operatorMatch := operatorPattern.FindString(input[i:]); operatorMatch != "" && strings.HasPrefix(input[i:], operatorMatch) {
			tokens = append(tokens, Token{Type: OPERATOR, Value: operatorMatch})
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
		case NUMBER:
			if i > 0 && tokens[i-1].Type == NUMBER {
				return fmt.Errorf("consecutive numbers not allowed: %s %s", tokens[i-1].Value, token.Value)
			}
		}
	}

	return nil
}