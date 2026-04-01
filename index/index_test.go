package index

import (
	"math"
	"testing"
)

func TestAddDocument_FileAccess(t *testing.T) {
	idx := NewIndex()

	t.Run("file does not exist", func(t *testing.T) {
		err := idx.AddDocument("non-existent-file.txt")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		path := makeTempFile(t, "")
		err := idx.AddDocument(path)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doc := idx.Docs[1]
		if doc.Length != 0 {
			t.Errorf("expected length 0, got %d", doc.Length)
		}

		if len(idx.Postings) != 0 {
			t.Errorf("expected no postings, got %v", idx.Postings)
		}
	})

	t.Run("only delimiters", func(t *testing.T) {
		path := makeTempFile(t, "   \n!!!,,,\n")
		err := idx.AddDocument(path)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doc := idx.Docs[2]
		if doc.Length != 0 {
			t.Errorf("expected length 0, got %d", doc.Length)
		}
	})
}

func TestAddDocument_Postings(t *testing.T) {
	t.Run("single term single document", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "golang")

		_ = idx.AddDocument(path)

		postings := idx.Postings["golang"]

		if len(postings) != 1 {
			t.Fatalf("expected 1 posting, got %d", len(postings))
		}

		p := postings[0]
		if p.DocID != 1 || p.Frequency != 1 || len(p.Positions) != 1 || p.Positions[0] != 0 {
			t.Errorf("unexpected posting: %+v", p)
		}
	})

	t.Run("same term multiple times", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "golang is golang")

		_ = idx.AddDocument(path)

		postings := idx.Postings["golang"]
		p := postings[0]

		if p.Frequency != 2 {
			t.Errorf("expected frequency 2, got %d", p.Frequency)
		}

		expected := []int{0, 2}
		if len(p.Positions) != 2 || p.Positions[0] != expected[0] || p.Positions[1] != expected[1] {
			t.Errorf("expected positions %v, got %v", expected, p.Positions)
		}
	})

	t.Run("multiple distinct terms", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "go fast")

		_ = idx.AddDocument(path)

		if len(idx.Postings) != 2 {
			t.Errorf("expected 2 terms, got %d", len(idx.Postings))
		}
	})

	t.Run("terms across multiple lines", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "golang\nsearch")

		_ = idx.AddDocument(path)

		if _, ok := idx.Postings["golang"]; !ok {
			t.Errorf("missing term golang")
		}

		if _, ok := idx.Postings["search"]; !ok {
			t.Errorf("missing term search")
		}
	})
}

func TestAddDocument_MultipleDocs(t *testing.T) {
	idx := NewIndex()

	path1 := makeTempFile(t, "golang")
	path2 := makeTempFile(t, "golang")

	_ = idx.AddDocument(path1)
	_ = idx.AddDocument(path2)

	postings := idx.Postings["golang"]

	if len(postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(postings))
	}

	if postings[0].DocID != 1 || postings[1].DocID != 2 {
		t.Errorf("docIDs incorrect: %+v", postings)
	}
}

func TestAddDocument_DocStore(t *testing.T) {
	idx := NewIndex()

	path := makeTempFile(t, "golang is fast")
	_ = idx.AddDocument(path)

	doc := idx.Docs[1]

	if doc.FilePath != path {
		t.Errorf("expected filepath %s, got %s", path, doc.FilePath)
	}

	if doc.Length != 2 {
		t.Errorf("expected length 2, got %d", doc.Length)
	}

	if idx.NextDocID != 2 {
		t.Errorf("expected nextDocID 2, got %d", idx.NextDocID)
	}
}

func TestSearch(t *testing.T) {

	t.Run("single term query", func(t *testing.T) {
		idx := NewIndex()

		path := makeTempFile(t, "golang is fast")
		_ = idx.AddDocument(path)

		results := idx.Search("golang")

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}

		if results[0].ID != 1 {
			t.Errorf("expected docID 1, got %d", results[0].ID)
		}
	})

	t.Run("two term query both exist", func(t *testing.T) {
		idx := NewIndex()

		path := makeTempFile(t, "golang is fast")
		_ = idx.AddDocument(path)

		results := idx.Search("golang fast")

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}

		if results[0].ID != 1 {
			t.Errorf("expected docID 1, got %d", results[0].ID)
		}
	})

	t.Run("two term query one missing", func(t *testing.T) {
		idx := NewIndex()

		path := makeTempFile(t, "golang is fast")
		_ = idx.AddDocument(path)

		results := idx.Search("golang python")

		if len(results) != 0 {
			t.Errorf("expected no results, got %v", results)
		}
	})

	t.Run("empty query", func(t *testing.T) {
		idx := NewIndex()

		path := makeTempFile(t, "golang is fast")
		_ = idx.AddDocument(path)

		results := idx.Search("")

		if len(results) != 0 {
			t.Errorf("expected empty result, got %v", results)
		}
	})
}

func TestTermFrequency(t *testing.T) {
	tests := []struct {
		name    string
		posting Posting
		doc     Document
		want    float64
	}{
		{
			name:    "basic tf",
			posting: Posting{Frequency: 2},
			doc:     Document{Length: 4},
			want:    0.5,
		},
		{
			name:    "higher frequency",
			posting: Posting{Frequency: 5},
			doc:     Document{Length: 10},
			want:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := termFrequency(tt.posting, tt.doc)

			if got != tt.want {
				t.Errorf("termFrequency() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestInverseDocumentFrequency(t *testing.T) {
	tests := []struct {
		name      string
		listLen   int
		totalDocs int
		want      float64
	}{
		{
			name:      "basic idf",
			listLen:   2,
			totalDocs: 10,
			want:      math.Log(5),
		},
		{
			name:      "term in all docs",
			listLen:   10,
			totalDocs: 10,
			want:      0.0, // log(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inverseDocumentFrequency(tt.listLen, tt.totalDocs)

			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("idf() = %f, want %f", got, tt.want)
			}
		})
	}
}
