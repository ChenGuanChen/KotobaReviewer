package db

import (
	"encoding/json"
	"os"
)

// Save writes db to disk as indented JSON, overwriting any previous copy.
func (db *VocabDB) Save() error {
	path, err := dataFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
