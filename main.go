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
		if err := addDocument(arg); err != nil {
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

func addDocument(filePath string) error {
	idx, err := index.Load(indexPath)
	if err != nil {
		return fmt.Errorf("loading index: %w", err)
	}

	if err := idx.AddDocument(filePath); err != nil {
		return fmt.Errorf("indexing %s: %w", filePath, err)
	}

	if err := idx.Save(indexPath); err != nil {
		return fmt.Errorf("saving index: %w", err)
	}

	doc := idx.Docs[idx.NextDocID-1]
	fmt.Printf("Indexed: %s (%d tokens)\n", filePath, doc.Length)
	return nil
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
	fmt.Println(" go run main.go add <filepath>")
	fmt.Println(" go run main.go search <query>")
}
