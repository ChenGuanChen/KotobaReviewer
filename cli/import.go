package cli

import (
	"fmt"
	"kotobaReviewer/db"
	"kotobaReviewer/entry"
	"kotobaReviewer/parser"
)

func RunImport(path string) {
	database, err := parser.ParseDocFile(path)
	if err != nil {
		fmt.Println("Error parsing:", err)
		return
	}

	if err := database.Save(); err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	vocabCount, grammarCount, needsReview, hasMeaning := 0, 0, 0, 0
	entries := database.Database
	for _, e := range entries {
		if e.Type == entry.Vocab {
			vocabCount++
		} else {
			grammarCount++
		}
		if e.NeedsReview {
			needsReview++
		}
		if e.Meaning != "" {
			hasMeaning++
		}
	}

	fmt.Printf("Imported %d entries -> %s\n", len(entries), db.GetDataFileName())
	fmt.Printf("  vocab: %d, grammar: %d\n", vocabCount, grammarCount)
	fmt.Printf("  have a starting meaning already: %d\n", hasMeaning)
	fmt.Printf("  still need a meaning typed in: %d\n", len(entries)-hasMeaning)
	fmt.Printf("  flagged needs_review (parser was unsure): %d\n", needsReview)
}
