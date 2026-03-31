package index

import "math"

// Score calculates the relevance score of a document for a given posting using TF-IDF.
func Score(posting Posting, doc Document, docListLen int, totalDocs int) float64 {
	if doc.Length == 0 || docListLen == 0 {
		return 0.0
	}
	tf := termFrequency(posting, doc)
	idf := inverseDocumentFrequency(docListLen, totalDocs)
	return tf * idf
}

func termFrequency(posting Posting, doc Document) float64 {
	return float64(posting.Frequency) / float64(doc.Length)
}

func inverseDocumentFrequency(listLen int, totalDocs int) float64 {
	return math.Log(float64(totalDocs) / float64(listLen))
}
