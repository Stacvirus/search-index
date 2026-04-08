package index

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTempFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp("", "testdoc-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { os.Remove(f.Name()) })

	return f.Name()
}

func makeTempDir(t *testing.T, files map[string]string) string {
	t.Helper()

	root, err := os.MkdirTemp("", "testdir-*")
	if err != nil {
		t.Fatal(err)
	}

	for path, content := range files {
		fullPath := filepath.Join(root, path)

		// create parent dirs if needed
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() { os.RemoveAll(root) })

	return root
}
