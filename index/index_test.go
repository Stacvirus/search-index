package index

import (
	"testing"
)

func TestAddDocument_FileAccess(t *testing.T) {
	idx := New()

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

		doc := idx.docs[1]
		if doc.Length != 0 {
			t.Errorf("expected length 0, got %d", doc.Length)
		}

		if len(idx.postings) != 0 {
			t.Errorf("expected no postings, got %v", idx.postings)
		}
	})

	t.Run("only delimiters", func(t *testing.T) {
		path := makeTempFile(t, "   \n!!!,,,\n")
		err := idx.AddDocument(path)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doc := idx.docs[2]
		if doc.Length != 0 {
			t.Errorf("expected length 0, got %d", doc.Length)
		}
	})
}

func TestAddDocument_Postings(t *testing.T) {
	t.Run("single term single document", func(t *testing.T) {
		idx := New()
		path := makeTempFile(t, "golang")

		_ = idx.AddDocument(path)

		postings := idx.postings["golang"]

		if len(postings) != 1 {
			t.Fatalf("expected 1 posting, got %d", len(postings))
		}

		p := postings[0]
		if p.DocID != 1 || p.Frequency != 1 || len(p.Positions) != 1 || p.Positions[0] != 0 {
			t.Errorf("unexpected posting: %+v", p)
		}
	})

	t.Run("same term multiple times", func(t *testing.T) {
		idx := New()
		path := makeTempFile(t, "golang is golang")

		_ = idx.AddDocument(path)

		postings := idx.postings["golang"]
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
		idx := New()
		path := makeTempFile(t, "go fast")

		_ = idx.AddDocument(path)

		if len(idx.postings) != 2 {
			t.Errorf("expected 2 terms, got %d", len(idx.postings))
		}
	})

	t.Run("terms across multiple lines", func(t *testing.T) {
		idx := New()
		path := makeTempFile(t, "golang\nsearch")

		_ = idx.AddDocument(path)

		if _, ok := idx.postings["golang"]; !ok {
			t.Errorf("missing term golang")
		}

		if _, ok := idx.postings["search"]; !ok {
			t.Errorf("missing term search")
		}
	})
}

func TestAddDocument_MultipleDocs(t *testing.T) {
	idx := New()

	path1 := makeTempFile(t, "golang")
	path2 := makeTempFile(t, "golang")

	_ = idx.AddDocument(path1)
	_ = idx.AddDocument(path2)

	postings := idx.postings["golang"]

	if len(postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(postings))
	}

	if postings[0].DocID != 1 || postings[1].DocID != 2 {
		t.Errorf("docIDs incorrect: %+v", postings)
	}
}

func TestAddDocument_DocStore(t *testing.T) {
	idx := New()

	path := makeTempFile(t, "golang is fast")
	_ = idx.AddDocument(path)

	doc := idx.docs[1]

	if doc.FilePath != path {
		t.Errorf("expected filepath %s, got %s", path, doc.FilePath)
	}

	if doc.Length != 2 {
		t.Errorf("expected length 2, got %d", doc.Length)
	}

	if idx.nextDocID != 2 {
		t.Errorf("expected nextDocID 2, got %d", idx.nextDocID)
	}
}
