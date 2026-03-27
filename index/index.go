package index

import (
	"bufio"
	"os"
	"sort"

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

func (idx *Index) Search(query string) []Document {
	// Analyze the query to get the tokens
	tokens := analysis.Analyze(query)
	lists := idx.getPostings(tokens)
	result := idx.intersectPostings(lists)

	var documents []Document
	for _, docID := range result {
		if doc, ok := idx.docs[docID]; ok {
			documents = append(documents, doc)
		}
	}
	return documents
}

func (idx *Index) getPostings(tokens []analysis.Token) []PostingList {
	if len(tokens) == 0 {
		return nil
	}

	var result []PostingList

	for _, token := range tokens {
		list := idx.postings[token.Word]
		if len(list) == 0 {
			return nil // If any token has no postings, the intersection will be empty
		}
		result = append(result, list)
	}
	return result
}

func (idx *Index) intersectPostings(lists []PostingList) []int {
	if len(lists) == 0 {
		return nil
	}
	// sort the posting lists in ascending order to optimize intersection
	sort.Slice(lists, func(i, j int) bool {
		return len(lists[i]) < len(lists[j])
	})

	result := lists[0]
	for _, list := range lists[1:] {
		result = intersectTwo(result, list)
		if len(result) == 0 {
			return nil // early exit
		}
	}

	var docIDs []int
	for _, posting := range result {
		docIDs = append(docIDs, posting.DocID)
	}
	return docIDs
}

func intersectTwo(list1, list2 PostingList) PostingList {
	var result PostingList
	p1, p2 := 0, 0
	for p1 < len(list1) && p2 < len(list2) {
		if list1[p1].DocID == list2[p2].DocID {
			result = append(result, list1[p1])
			p1++
			p2++
		} else if list1[p1].DocID < list2[p2].DocID {
			p1++
		} else {
			p2++
		}
	}
	return result
}
