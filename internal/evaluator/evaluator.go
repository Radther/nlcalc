package evaluator

import (
	"errors"
	"mathparser/internal/tokenizer"
	"strconv"
)

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

	expression := make([]tokenizer.Token, len(tokens))
	copy(expression, tokens)

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
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenizer.OPERATOR && (tokens[i].Value == "*" || tokens[i].Value == "/") {
			return i
		}
	}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == tokenizer.OPERATOR && (tokens[i].Value == "+" || tokens[i].Value == "-") {
			return i
		}
	}

	return -1
}