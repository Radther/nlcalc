// Package favourable applies heuristics to recover a calculable token sequence
// from one that would otherwise fail validation or evaluation.
package favourable

import (
	"github.com/radther/nlcalc/internal/tokenizer"
)

// Apply applies favourable parsing heuristics to a token sequence, attempting to
// recover a valid calculable expression. Rules are applied iteratively until stable:
//
//  1. Strip unmatched parentheses (closing with no opening, or opening with no closing).
//  2. Strip leading tokens that cannot validly start an expression: OPERATOR and ")".
//  3. Strip trailing tokens that cannot validly end an expression: OPERATOR, UNARY_OPERATOR,
//     FUNCTION, and "(". These strip rules cascade — once one invalid token is removed
//     the next may become invalid too.
//  4. Collapse consecutive OPERATOR tokens, keeping the first. This does not apply to
//     consecutive FUNCTION or UNARY_OPERATOR tokens, which can be legitimately adjacent
//     (e.g. "sqrt sqrt 16" or "--5").
//
// Consecutive NUMBER tokens are left untouched; two adjacent numbers remain an error
// since there is no unambiguous way to combine them.
func Apply(tokens []tokenizer.Token) []tokenizer.Token {
	if len(tokens) == 0 {
		return tokens
	}

	for {
		before := len(tokens)
		tokens = stripUnmatchedParentheses(tokens)
		tokens = stripLeadingInvalidTokens(tokens)
		tokens = stripTrailingInvalidTokens(tokens)
		tokens = collapseConsecutiveOperators(tokens)
		if len(tokens) == before {
			break
		}
	}

	return tokens
}

// stripLeadingInvalidTokens removes tokens from the start that cannot begin an expression:
// OPERATOR tokens and closing parentheses.
func stripLeadingInvalidTokens(tokens []tokenizer.Token) []tokenizer.Token {
	for len(tokens) > 0 {
		first := tokens[0]
		if first.Type == tokenizer.OPERATOR ||
			(first.Type == tokenizer.PARENTHESIS && first.Value == ")") {
			tokens = tokens[1:]
		} else {
			break
		}
	}
	return tokens
}

// stripTrailingInvalidTokens removes tokens from the end that cannot end an expression:
// OPERATOR, UNARY_OPERATOR, FUNCTION, and opening parentheses.
func stripTrailingInvalidTokens(tokens []tokenizer.Token) []tokenizer.Token {
	for len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		if last.Type == tokenizer.OPERATOR ||
			last.Type == tokenizer.UNARY_OPERATOR ||
			last.Type == tokenizer.FUNCTION ||
			(last.Type == tokenizer.PARENTHESIS && last.Value == "(") {
			tokens = tokens[:len(tokens)-1]
		} else {
			break
		}
	}
	return tokens
}

// collapseConsecutiveOperators removes the second of any two adjacent OPERATOR tokens,
// keeping the first. Applied repeatedly within a single pass to handle runs of three or more.
func collapseConsecutiveOperators(tokens []tokenizer.Token) []tokenizer.Token {
	if len(tokens) < 2 {
		return tokens
	}

	result := make([]tokenizer.Token, 0, len(tokens))
	result = append(result, tokens[0])

	for i := 1; i < len(tokens); i++ {
		prev := result[len(result)-1]
		curr := tokens[i]
		if prev.Type == tokenizer.OPERATOR && curr.Type == tokenizer.OPERATOR {
			// Keep prev (first occurrence), skip curr
			continue
		}
		result = append(result, curr)
	}

	return result
}

// stripUnmatchedParentheses removes parentheses that lack a matching counterpart:
// any ")" with no preceding "(" and any "(" with no following ")".
func stripUnmatchedParentheses(tokens []tokenizer.Token) []tokenizer.Token {
	remove := make([]bool, len(tokens))

	// Forward pass: mark unmatched ")"
	openStack := []int{} // indices of unmatched "(" tokens so far
	for i, tok := range tokens {
		if tok.Type != tokenizer.PARENTHESIS {
			continue
		}
		if tok.Value == "(" {
			openStack = append(openStack, i)
		} else { // ")"
			if len(openStack) > 0 {
				openStack = openStack[:len(openStack)-1] // matched pair
			} else {
				remove[i] = true // no opening paren to match
			}
		}
	}
	// Any remaining unmatched "(" indices
	for _, idx := range openStack {
		remove[idx] = true
	}

	result := make([]tokenizer.Token, 0, len(tokens))
	for i, tok := range tokens {
		if !remove[i] {
			result = append(result, tok)
		}
	}
	return result
}
