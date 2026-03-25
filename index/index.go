package index

import (
	"bufio"
	"os"

	"github.com/Stacvirus/search-index/analysis"
)

type Index struct {
	postings  map[string]PostingList
	docs      map[int]Document
	nextDocID int
}

type Document struct {
	ID       int
	FilePath string
	Length   int
}

func NewIndex() *Index {
	return &Index{
		postings:  make(map[string]PostingList),
		docs:      make(map[int]Document),
		nextDocID: 1,
	}
}

func (idx *Index) AddDocument(filePath string) error {
	totalTokens := 0
	docID := idx.nextDocID

	// Get access to the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the file line by line and analyze the content
	for scanner.Scan() {
		line := scanner.Text()
		tokens := analysis.Analyze(line) // Get the tokens from the line
		totalTokens += len(tokens)
		idx.handleTokens(tokens, docID)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Add the document to the index and increment the document ID
	idx.docs[docID] = Document{
		ID:       docID,
		FilePath: filePath,
		Length:   totalTokens,
	}
	idx.nextDocID++

	return nil
}

// Update the index with the tokens from the document
func (idx *Index) handleTokens(tokens []analysis.Token, docID int) {
	for _, token := range tokens {
		// Update the posting list for the token
		postingList := idx.postings[token.Word] // If the token doesn't exist, this will return an empty list

		// Create a new posting for the document
		posting := Posting{
			DocID:     docID,
			Frequency: 1,
			Positions: []int{token.Position},
		}

		// Check if the document already exists in the posting list
		found := false
		for i, p := range postingList {
			if p.DocID == docID {
				postingList[i].Frequency++
				postingList[i].Positions = append(postingList[i].Positions, token.Position)
				found = true
				break
			}
		}
		if !found {
			postingList = append(postingList, posting)
		}
		idx.postings[token.Word] = postingList
	}
}
