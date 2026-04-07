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
			return numVal(tokens[0]), nil
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

		resultToken := makeNumberToken(result)

		newExpression := make([]tokenizer.Token, 0, len(expression)-(end-start)+1)
		newExpression = append(newExpression, expression[:start]...)
		newExpression = append(newExpression, resultToken)
		newExpression = append(newExpression, expression[end+1:]...)
		expression = newExpression
	}

	return evaluateSimpleExpression(expression)
}

// makeNumberToken creates a NUMBER token with a pre-computed float64 value.
// Value is left empty since the evaluator only uses NumValue for intermediate results.
func makeNumberToken(value float64) tokenizer.Token {
	return tokenizer.Token{
		Type:     tokenizer.NUMBER,
		NumValue: value,
	}
}

// numVal extracts the float64 value from a NUMBER token. It uses the pre-parsed
// NumValue when available, falling back to parsing the Value string for tokens
// not created by the tokenizer (e.g., test fixtures).
func numVal(t tokenizer.Token) float64 {
	if t.NumValue != 0 {
		return t.NumValue
	}
	v, _ := strconv.ParseFloat(t.Value, 64)
	return v
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

				newValue := numVal(result[i+1])
				if result[i].Value == "-" {
					newValue = -newValue
				}

				newResult = append(newResult, makeNumberToken(newValue))
				i++ // Skip the number token since we've processed it
				changed = true
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
			return numVal(tokens[0]), nil
		}
		return 0.0, errors.New("invalid single token")
	}

	expression := tokens

	// Process functions and unary operators in a loop until both are fully resolved
	// This handles cases like:
	// - "abs -5" → process unary first → "abs (-5)" → "5"
	// - "-sqrt 16" → process function first → "-(4)" → "-4"
	for {
		madeProgress := false

		// 1. Try to process unary operators (converts UNARY_OP + NUMBER to NUMBER)
		newExpression := processUnaryOperators(expression)
		if len(newExpression) != len(expression) {
			expression = newExpression
			madeProgress = true
			continue
		}

		// 2. Try to process functions right-to-left
		funcIndex := findRightmostFunction(expression)
		if funcIndex != -1 {
			if funcIndex >= len(expression)-1 {
				return 0.0, errors.New("function requires argument")
			}

			// Next token must be NUMBER (unary operators have been processed)
			if expression[funcIndex+1].Type != tokenizer.NUMBER {
				// Can't process this function yet, might need more unary op processing
				// This shouldn't happen if unary ops are processed correctly
				return 0.0, errors.New("function argument must be a number")
			}

			arg := numVal(expression[funcIndex+1])

			result, err := applyFunction(expression[funcIndex].Value, arg)
			if err != nil {
				return 0.0, err
			}

			// Replace [FUNCTION, NUMBER] with result NUMBER
			resultToken := makeNumberToken(result)

			newExpression := make([]tokenizer.Token, 0, len(expression)-1)
			newExpression = append(newExpression, expression[:funcIndex]...)
			newExpression = append(newExpression, resultToken)
			newExpression = append(newExpression, expression[funcIndex+2:]...)
			expression = newExpression
			madeProgress = true
		}

		if !madeProgress {
			break
		}
	}

	for {
		opIndex := findHighestPrecedenceOperator(expression)
		if opIndex == -1 {
			break
		}

		if opIndex == 0 || opIndex == len(expression)-1 {
			return 0.0, errors.New("invalid operator position")
		}

		if expression[opIndex-1].Type != tokenizer.NUMBER {
			return 0.0, errors.New("invalid left operand")
		}
		if expression[opIndex+1].Type != tokenizer.NUMBER {
			return 0.0, errors.New("invalid right operand")
		}

		left := numVal(expression[opIndex-1])
		right := numVal(expression[opIndex+1])

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

		resultToken := makeNumberToken(result)

		newExpression := make([]tokenizer.Token, 0, len(expression)-2)
		newExpression = append(newExpression, expression[:opIndex-1]...)
		newExpression = append(newExpression, resultToken)
		newExpression = append(newExpression, expression[opIndex+2:]...)
		expression = newExpression
	}

	if len(expression) != 1 || expression[0].Type != tokenizer.NUMBER {
		return 0.0, errors.New("invalid expression structure")
	}

	return numVal(expression[0]), nil
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

// findRightmostFunction scans right-to-left to find the rightmost FUNCTION token.
// This enables right-associative function evaluation, so "sqrt sqrt 16" evaluates
// as "sqrt (sqrt 16)" naturally.
func findRightmostFunction(tokens []tokenizer.Token) int {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Type == tokenizer.FUNCTION {
			return i
		}
	}
	return -1
}

// applyFunction applies the named function to a value.
// Returns an error for invalid function names or domain errors (e.g., sqrt of negative).
func applyFunction(name string, value float64) (float64, error) {
	switch name {
	case "sqrt":
		if value < 0 {
			return 0, errors.New("square root of negative number")
		}
		return math.Sqrt(value), nil
	case "abs":
		return math.Abs(value), nil
	case "log":
		if value <= 0 {
			return 0, errors.New("logarithm of non-positive number")
		}
		return math.Log10(value), nil
	case "ln":
		if value <= 0 {
			return 0, errors.New("natural logarithm of non-positive number")
		}
		return math.Log(value), nil
	case "sin":
		return math.Sin(value), nil
	case "cos":
		return math.Cos(value), nil
	case "tan":
		return math.Tan(value), nil
	default:
		return 0, errors.New("unknown function: " + name)
	}
}
