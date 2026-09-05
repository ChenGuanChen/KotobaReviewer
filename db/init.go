package db

import (
	"encoding/json"
	"os"
)

// Init loads the saved VocabDB from disk. If no data file exists yet
// (first run), it returns an empty VocabDB and a nil error rather than
// treating a missing file as a failure.
func Init() (VocabDB, error) {
	path, err := dataFilePath()
	if err != nil {
		return VocabDB{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return VocabDB{}, nil
		}
		return VocabDB{}, err
	}
	var v VocabDB
	if err := json.Unmarshal(data, &v); err != nil {
		return VocabDB{}, err
	}
	return v, nil
}
