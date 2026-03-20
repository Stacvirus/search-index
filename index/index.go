package index

type Index struct {
	postings  map[string]PostingList
	docs      map[string]Document
	nextDocID int
}

type Document struct {
	ID       int
	FilePath string
	Length   int
}
