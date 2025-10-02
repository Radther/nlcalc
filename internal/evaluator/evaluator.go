// Package evaluator calculates the final result from a sequence of mathematical tokens.
// It implements standard order of operations (PEMDAS/BODMAS), handling parentheses,
// exponentiation, multiplication, division, addition, and subtraction with proper precedence.
package evaluator

import (
	"errors"
	"math"
	"strconv"

	"github.com/radther/nlcalc/internal/tokenizer"
)

// Evaluate calculates the result of a mathematical expression represented as tokens.
// It follows standard order of operations (PEMDAS/BODMAS):
//   - Parentheses (evaluated from innermost to outermost)
//   - Exponentiation (right to left, right-associative)
//   - Multiplication and Division (left to right)
//   - Addition and Subtraction (left to right)
//
// Returns an error if the expression is empty, invalid, has mismatched parentheses,
// or attempts division by zero.
func Evaluate(tokens []tokenizer.Token) (float64, error) {
	if len(tokens) == 0 {
		return 0.0, errors.New("empty expression")
	}

	if len(tokens) == 1 {
		if tokens[0].Type == tokenizer.NUMBER {
			return strconv.ParseFloat(tokens[0].Value, 64)
		}
		return 0.0, errors.New("invalid single token expression")
	}

	return evaluateExpression(tokens)
}

func evaluateExpression(tokens []tokenizer.Token) (float64, error) {
	expression := make([]tokenizer.Token, len(tokens))
	copy(expression, tokens)

	for {
		parenthesesIndex := findInnerParentheses(expression)
		if parenthesesIndex == -1 {
			break
		}

		start := parenthesesIndex
		end := findClosingParenthesis(expression, start)
		if end == -1 {
			return 0.0, errors.New("mismatched parentheses")
		}

		subExpr := expression[start+1 : end]
		result, err := evaluateSimpleExpression(subExpr)
		if err != nil {
			return 0.0, err
		}

		resultToken := tokenizer.Token{
			Type:  tokenizer.NUMBER,
			Value: strconv.FormatFloat(result, 'f', -1, 64),
		}

		newExpression := make([]tokenizer.Token, 0, len(expression)-(end-start)+1)
		newExpression = append(newExpression, expression[:start]...)
		newExpression = append(newExpression, resultToken)
		newExpression = append(newExpression, expression[end+1:]...)
		expression = newExpression
	}

	return evaluateSimpleExpression(expression)
}

// processUnaryOperators handles unary operators (+ and -) by applying them to numbers.
// Unary minus negates the value, unary plus is identity (returns value unchanged).
// It scans for [UNARY_OPERATOR, NUMBER] patterns and replaces them with modified NUMBER tokens.
// This function is called iteratively until no more unary operators remain, to handle
// cases like consecutive unary operators (e.g., -+5 or +-5).
func processUnaryOperators(tokens []tokenizer.Token) []tokenizer.Token {
	result := tokens

	// Keep processing until no more unary operators can be reduced
	for {
		changed := false
		newResult := make([]tokenizer.Token, 0, len(result))

		for i := 0; i < len(result); i++ {
			// Check for unary operator followed by a number
			if result[i].Type == tokenizer.UNARY_OPERATOR &&
			   i+1 < len(result) &&
			   result[i+1].Type == tokenizer.NUMBER {

				// Parse the number and apply the unary operator
				value, err := strconv.ParseFloat(result[i+1].Value, 64)
				if err == nil {
					var newValue float64
					if result[i].Value == "-" {
						newValue = -value // Negate
					} else {
						newValue = value // Identity (unary +)
					}

					newResult = append(newResult, tokenizer.Token{
						Type:  tokenizer.NUMBER,
						Value: strconv.FormatFloat(newValue, 'f', -1, 64),
					})
					i++ // Skip the number token since we've processed it
					changed = true
				} else {
					// If parsing fails, keep original tokens
					newResult = append(newResult, result[i])
				}
			} else {
				newResult = append(newResult, result[i])
			}
		}

		result = newResult
		if !changed {
			break
		}
	}

	return result
}

func evaluateSimpleExpression(tokens []tokenizer.Token) (float64, error) {
	if len(tokens) == 0 {
		return 0.0, errors.New("empty expression")
	}

	if len(tokens) == 1 {
		if tokens[0].Type == tokenizer.NUMBER {
			return strconv.ParseFloat(tokens[0].Value, 64)
		}
		return 0.0, errors.New("invalid single token")
	}

	// Process unary operators first
	expression := processUnaryOperators(tokens)

	for {
		opIndex := findHighestPrecedenceOperator(expression)
		if opIndex == -1 {
			break
		}

		if opIndex == 0 || opIndex == len(expression)-1 {
			return 0.0, errors.New("invalid operator position")
		}

		left, err := strconv.ParseFloat(expression[opIndex-1].Value, 64)
		if err != nil {
			return 0.0, errors.New("invalid left operand")
		}

		right, err := strconv.ParseFloat(expression[opIndex+1].Value, 64)
		if err != nil {
			return 0.0, errors.New("invalid right operand")
		}

		var result float64
		switch expression[opIndex].Value {
		case "^":
			result = math.Pow(left, right)
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return 0.0, errors.New("division by zero")
			}
			result = left / right
		case "+":
			result = left + right
		case "-":
			result = left - right
		default:
			return 0.0, errors.New("unknown operator")
		}

		resultToken := tokenizer.Token{
			Type:  tokenizer.NUMBER,
			Value: strconv.FormatFloat(result, 'f', -1, 64),
		}

		newExpression := make([]tokenizer.Token, 0, len(expression)-2)
		newExpression = append(newExpression, expression[:opIndex-1]...)
		newExpression = append(newExpression, resultToken)
		newExpression = append(newExpression, expression[opIndex+2:]...)
		expression = newExpression
	}

	if len(expression) != 1 || expression[0].Type != tokenizer.NUMBER {
		return 0.0, errors.New("invalid expression structure")
	}

	return strconv.ParseFloat(expression[0].Value, 64)
}

func findInnerParentheses(tokens []tokenizer.Token) int {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Type == tokenizer.PARENTHESIS && tokens[i].Value == "(" {
			return i
		}
	}
	return -1
}

func findClosingParenthesis(tokens []tokenizer.Token, start int) int {
	for i := start + 1; i < len(tokens); i++ {
		if tokens[i].Type == tokenizer.PARENTHESIS && tokens[i].Value == ")" {
			return i
		}
	}
	return -1
}

func findHighestPrecedenceOperator(tokens []tokenizer.Token) int {
	// Exponentiation - RIGHT-associative (scan right-to-left)
	// This ensures 2^3^2 evaluates as 2^(3^2) = 512, not (2^3)^2 = 64
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Type == tokenizer.OPERATOR && tokens[i].Value == "^" {
			return i
		}
	}

	// Multiplication/Division - left-associative (scan left-to-right)
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenizer.OPERATOR && (tokens[i].Value == "*" || tokens[i].Value == "/") {
			return i
		}
	}

	// Addition/Subtraction - left-associative (scan left-to-right)
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenizer.OPERATOR && (tokens[i].Value == "+" || tokens[i].Value == "-") {
			return i
		}
	}

	return -1
}