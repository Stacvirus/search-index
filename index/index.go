package index

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/Stacvirus/search-index/analysis"
)

type Index struct {
	Postings  map[string]PostingList
	Docs      map[int]Document
	NextDocID int
}

type Document struct {
	ID       int
	FilePath string
	Length   int
}

type scoredDocs struct {
	doc   Document
	score float64
}

func NewIndex() *Index {
	return &Index{
		Postings:  make(map[string]PostingList),
		Docs:      make(map[int]Document),
		NextDocID: 1,
	}
}

func (idx *Index) AddDocuments(root string, extensions []string) error {
	extSet := make(map[string]bool)
	for _, ext := range extensions {
		extSet[ext] = true
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // Skip directories, we only want to index files
		}
		if !extSet[filepath.Ext(path)] {
			return nil // Skip files that don't have the specified extensions
		}

		// store path relative to the root directory in the index, so that we can move
		// the index to a different location without breaking the file paths
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		return idx.AddDocument(filepath.Join(root, relPath))
	})
}

func (idx *Index) AddDocument(filePath string) error {
	totalTokens := 0
	docID := idx.NextDocID

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
		idx.buildIndex(tokens, docID)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Add the document to the index and increment the document ID
	idx.Docs[docID] = Document{
		ID:       docID,
		FilePath: filePath,
		Length:   totalTokens,
	}
	idx.NextDocID++

	return nil
}

// Update the index with the tokens from the document
func (idx *Index) buildIndex(tokens []analysis.Token, docID int) {
	for _, token := range tokens {
		// Update the posting list for the token
		postingList := idx.Postings[token.Word] // If the token doesn't exist, this will return an empty list

		// Create a new posting for the document
		posting := Posting{
			DocID:     docID,
			Frequency: 1,
			Positions: []int{token.Position},
		}

		// Check if the document already exists in the posting list
		postingListLength := len(postingList)
		if postingListLength > 0 && postingList[postingListLength-1].DocID == docID {
			last := postingListLength - 1
			postingList[last].Frequency++
			postingList[last].Positions = append(postingList[last].Positions, token.Position)
		} else {
			postingList = append(postingList, posting)
		}
		idx.Postings[token.Word] = postingList
	}
}

func (idx *Index) Search(query string) []Document {
	// Analyze the query to get the tokens
	tokens := analysis.Analyze(query)
	lists := idx.getPostings(tokens)
	result := idx.unionPostings(lists)
	scoredList := idx.rankDocuments(result, tokens)

	var documents []Document
	for _, scored := range scoredList {
		documents = append(documents, scored.doc)
	}
	return documents
}

func (idx *Index) getPostings(tokens []analysis.Token) []PostingList {
	if len(tokens) == 0 {
		return nil
	}

	var result []PostingList

	for _, token := range tokens {
		list := idx.Postings[token.Word]
		if len(list) == 0 {
			continue // We can skip tokens that don't exist in the index, as they won't contribute to the search results
		}
		result = append(result, list)
	}
	return result
}

// Union posting lists merges all the postings from the lists,
// while intersection only keeps the postings that are present in all lists.
// Union is used for OR queries, while intersection is used for AND queries.
func (idx *Index) unionPostings(lists []PostingList) PostingList {
	if len(lists) == 0 {
		return nil
	}
	seen := make(map[int]struct{})
	var result PostingList
	for _, list := range lists {
		for _, posting := range list {
			if _, exists := seen[posting.DocID]; !exists {
				result = append(result, posting)
				seen[posting.DocID] = struct{}{}
			}
		}
	}
	return result
}

func (idx *Index) rankDocuments(postings PostingList, tokens []analysis.Token) []scoredDocs {
	var scoredList []scoredDocs

	for _, posting := range postings {
		doc := idx.Docs[posting.DocID]
		totalScore := 0.0
		for _, token := range tokens {
			list := idx.Postings[token.Word]
			totalScore += Score(searchPosting(list, posting.DocID), doc, len(list), idx.NextDocID-1)
		}
		scoredList = append(scoredList, scoredDocs{doc, totalScore})
	}

	// Sort the scored documents by score in descending order
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	return scoredList
}

func searchPosting(list PostingList, docID int) Posting {
	last, first := len(list), 0
	for first < last {
		mid := (first + last) / 2
		if list[mid].DocID == docID {
			return list[mid]
		} else if list[mid].DocID < docID {
			first = mid + 1
		} else {
			last = mid
		}
	}
	return Posting{} // Branch suppose to never be reached if the document is guaranteed to be in the list
}
