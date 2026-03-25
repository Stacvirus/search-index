package index

import (
	"os"
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
