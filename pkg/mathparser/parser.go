package mathparser

import (
	"nlcalc/internal/cleaner"
	"nlcalc/internal/evaluator"
	"nlcalc/internal/normalizer"
	"nlcalc/internal/tokenizer"
)

func Parse(input string) (float64, error) {
	normalized := normalizer.Normalize(input)
	cleaned := cleaner.Clean(normalized)
	tokens, err := tokenizer.Tokenize(cleaned)
	if err != nil {
		return 0, err
	}
	
	result, err := evaluator.Evaluate(tokens)
	if err != nil {
		return 0, err
	}
	
	return result, nil
}