package analysis

import "strings"

// Normalize converts a token to lowercase.
// This is a simple normalization step that helps to ensure
// that different forms of the same word are treated as the same term in the index.
func Normalize(token Token) Token {
	return Token{
		Word:     strings.ToLower(token.Word),
		Position: token.Position,
	}
}

// NormalizeAll applies normalization to a slice of tokens.
func NormalizeAll(tokens []Token) []Token {
	var normalized = make([]Token, 0, len(tokens))
	for _, token := range tokens {
		normalized = append(normalized, Normalize(token))
	}
	return normalized
}
