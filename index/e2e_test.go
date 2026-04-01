package index

import "testing"

func TestEndToEndSearch(t *testing.T) {

	t.Run("one document contains term", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "golang is fast")

		_ = idx.AddDocument(path)

		results := idx.Search("golang")

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("one document does not contain term", func(t *testing.T) {
		idx := NewIndex()
		path := makeTempFile(t, "golang is fast")

		_ = idx.AddDocument(path)

		results := idx.Search("python")

		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("two documents, term in only one", func(t *testing.T) {
		idx := NewIndex()

		p1 := makeTempFile(t, "golang is fast")
		p2 := makeTempFile(t, "python is popular")

		_ = idx.AddDocument(p1)
		_ = idx.AddDocument(p2)

		results := idx.Search("golang")

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}

		if results[0].ID != 1 {
			t.Errorf("expected doc 1, got %d", results[0].ID)
		}
	})

	t.Run("two documents, both contain term, higher score first", func(t *testing.T) {
		idx := NewIndex()

		// doc1: 2 occurrences
		p1 := makeTempFile(t, "golang golang fast")
		// doc2: 1 occurrence
		p2 := makeTempFile(t, "golang is great")

		_ = idx.AddDocument(p1)
		_ = idx.AddDocument(p2)

		results := idx.Search("golang")

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		if results[0].ID != 1 {
			t.Errorf("expected doc 1 ranked first, got %d", results[0].ID)
		}
	})

	t.Run("multi-term query, only one doc matches all terms", func(t *testing.T) {
		idx := NewIndex()

		p1 := makeTempFile(t, "golang is fast")
		p2 := makeTempFile(t, "golang is popular")

		_ = idx.AddDocument(p1)
		_ = idx.AddDocument(p2)

		results := idx.Search("golang fast")

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}

		if results[0].ID != 1 {
			t.Errorf("expected doc 1, got %d", results[0].ID)
		}
	})

	t.Run("multi-term query, no document matches all terms", func(t *testing.T) {
		idx := NewIndex()

		p1 := makeTempFile(t, "golang fast")
		p2 := makeTempFile(t, "python slow")

		_ = idx.AddDocument(p1)
		_ = idx.AddDocument(p2)

		results := idx.Search("golang slow")

		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("empty query returns no results", func(t *testing.T) {
		idx := NewIndex()

		p1 := makeTempFile(t, "golang fast")
		_ = idx.AddDocument(p1)

		results := idx.Search("")

		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})
}
