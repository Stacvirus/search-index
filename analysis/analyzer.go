package analysis

func Analyze(text string) []Token {
	tokens := Tokenize(text)
	normalized := NormalizeAll(tokens)
	filtered := FilterStopWords(normalized)
	return filtered
}
