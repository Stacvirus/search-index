package analysis

// StopWords is a set of common words that are often filtered out during text analysis
// because they are considered to have little semantic value in search queries.
var stopWords = map[string]struct{}{
	"the": {}, "is": {}, "at": {}, "which": {}, "on": {},
	"and": {}, "a": {}, "an": {}, "in": {}, "with": {},
	"to": {}, "of": {}, "for": {}, "by": {}, "that": {},
	"this": {}, "it": {}, "as": {}, "are": {}, "was": {},
	"be": {}, "from": {}, "or": {}, "but": {}, "not": {},
	"over": {}, "under": {}, "between": {}, "into": {}, "out": {},
	"can": {}, "its": {}, "have": {}, "has": {},
	"had": {}, "do": {}, "did": {}, "will": {}, "would": {},
	"could": {}, "should": {}, "may": {}, "might": {},
}

// checks if a given word is a stop word.
func IsStopWord(word string) bool {
	_, exists := stopWords[word]
	return exists
}

// removes stop words from a slice of tokens.
func FilterStopWords(tokens []Token) []Token {
	var filtered []Token
	for _, token := range tokens {
		if !IsStopWord(token.Word) {
			filtered = append(filtered, token)
		}
	}
	return filtered
}
