package analysis

import (
	"regexp"
)

type Token struct {
	Word     string
	Position int
}

var splitRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// This is a very simple tokenizer that splits on whitespace, punctuation, symbols etc.
func Tokenize(text string) []Token {
	tokens := splitRegex.Split(text, -1)
	position := 0 // track the position of each token in the original text to avoid whitespace items eventually produced by Split
	var result []Token
	for _, token := range tokens {
		if token != "" {
			result = append(result, Token{Word: token, Position: position})
			position++
		}
	}
	return result
}
