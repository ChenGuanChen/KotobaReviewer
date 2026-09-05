package parser

import (
	"kotobaReviewer/entry"
	"strings"
)

// ClassifyType decides vocab vs. grammar from a raw doc line. Shared by the
// importer (for fresh imports) and the `reclassify` command (for patching
// entries already sitting in vocab.json without disturbing anything else
// about them — IDs, SRS progress, meanings, all untouched).
func ClassifyType(raw string) entry.EntryType {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return entry.Vocab
	}
	if firstRune := []rune(trimmed)[0]; strings.ContainsRune(grammarStartChars, firstRune) {
		return entry.Grammar
	}
	for _, j := range grammarJargon {
		if strings.Contains(trimmed, j) {
			return entry.Grammar
		}
	}
	return entry.Vocab
}
