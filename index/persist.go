package index

import (
	"encoding/gob"
	"errors"
	"os"
)

// Save writes the index to a file at the specified path.
// If the file already exists, it will be overwritten.
func (idx *Index) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewEncoder(file).Encode(idx)
}

// Load reads the index from a file at the specified path.
// If the file does not exist, it returns a new empty index.
func Load(path string) (*Index, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewIndex(), nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	idx := &Index{}
	if err := gob.NewDecoder(file).Decode(idx); err != nil {
		return nil, err
	}
	return idx, nil
}
