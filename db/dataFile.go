package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const dataFileName = "vocab.json"

func GetDataFileName() string {
	return dataFileName
}

func dataFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	if strings.Contains(exePath, "go-build") {
		root, err := moduleRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, dataFileName), nil
	}
	return filepath.Join(filepath.Dir(exePath), dataFileName), nil
}

func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}
