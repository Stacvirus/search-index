package main

import (
	"fmt"
	"os"

	"github.com/Stacvirus/search-index/index"
)

const indexPath = "search.idx"

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	arg := os.Args[2]

	switch command {
	case "add":
		if err := addPath(arg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "search":
		if err := search(arg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}

}

func addPath(path string) error {
	idx, err := index.Load(indexPath)
	if err != nil {
		return fmt.Errorf("loading index: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("reading path: %w", err)
	}

	if info.IsDir() {
		extensions := []string{".txt", ".md"}
		if err := idx.AddDocuments(path, extensions); err != nil {
			return fmt.Errorf("indexing directory %s: %w", path, err)
		}
		fmt.Printf("indexed directory: %s\n", path)
	} else {
		if err := idx.AddDocument(path); err != nil {
			return fmt.Errorf("indexing %s: %w", path, err)
		}
		doc := idx.Docs[idx.NextDocID-1]
		fmt.Printf("indexed: %s (%d tokens)\n", path, doc.Length)
	}

	return idx.Save(indexPath)
}

func search(query string) error {
	idx, err := index.Load(indexPath)
	if err != nil {
		return fmt.Errorf("loading index: %w", err)
	}

	results := idx.Search(query)
	if len(results) == 0 {
		fmt.Printf("no results found for %q\n", query)
		return nil
	}

	fmt.Printf("Results for %q:\n", query)
	for i, doc := range results {
		fmt.Printf(" %d. %s\n", i+1, doc.FilePath)
	}
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" go run main.go add <filepath.extension>")
	fmt.Println(" go run main.go add <directory>")
	fmt.Println(" go run main.go search <query>")
}
