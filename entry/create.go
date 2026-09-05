package entry

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func newID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// New creates a blank Entry with a freshly generated ID, ready to have
// its other fields filled in directly.
func New() (*Entry, error) {
	id, err := newID()
	if err != nil {
		return nil, fmt.Errorf("creating entry: %w", err)
	}
	return &Entry{
		ID: id,
	}, nil
}

// NewFromLine creates an Entry with a freshly generated ID and RawLine set
// to line, for parser code to fill in further as it parses that line.
func NewFromLine(line string) (*Entry, error) {
	id, err := newID()
	if err != nil {
		return nil, fmt.Errorf("creating entry with line %q: %w", line, err)
	}
	return &Entry{
		ID:      id,
		RawLine: line,
	}, nil
}

// NewManual creates an Entry from fields typed in directly (e.g. via a CLI
// mode), rather than parsed from a doc line. RawLine falls back to word,
// since manually-added entries have no original line to keep around.
func NewManual(entryType EntryType, word, reading, meaning, tag, notes string) (*Entry, error) {
	id, err := newID()
	if err != nil {
		return nil, fmt.Errorf("creating entry for %q using full info: %w", word, err)
	}
	return &Entry{
		ID:      id,
		Type:    entryType,
		Word:    word,
		Reading: reading,
		Meaning: meaning,
		Tag:     tag,
		Notes:   notes,
		RawLine: word,
	}, nil
}
