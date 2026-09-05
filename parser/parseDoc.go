package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/ChenGuanChen/kotobaReviewer/db"
)

// ParseDocFile reads your exported Google Doc text file and turns it into
// VocabEntry structs. It never throws away data it doesn't understand —
// anything ambiguous gets NeedsReview=true and keeps its RawLine so you can
// fix it up by hand afterward.
func ParseDocFile(path string) (db.VocabDB, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return db.VocabDB{}, err
	}

	text := strings.TrimPrefix(string(raw), "\ufeff") // strip BOM if present
	lines := strings.Split(text, "\n")

	footnoteStart := len(lines)
	for i, l := range lines {
		if footnoteStartRe.MatchString(l) {
			footnoteStart = i
			break
		}
	}

	footnotes := parseFootnotes(lines[footnoteStart:])

	var database db.VocabDB
	for i, l := range lines[:footnoteStart] {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		entry, err := parseEntryLine(l, footnotes)
		if err != nil {
			return db.VocabDB{}, fmt.Errorf("Error when parsing line %v: %w", i, err)
		}
		database.Database = append(database.Database, *entry)
	}
	return database, nil
}
