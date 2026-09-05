package cli

import (
	"fmt"

	"github.com/ChenGuanChen/kotobaReviewer/db"
	"github.com/ChenGuanChen/kotobaReviewer/parser"
)

func RunReclassify() {
	database, err := db.Init()
	if err != nil {
		fmt.Println("Error loading:", err)
		return
	}

	entries := database.Database
	if len(entries) == 0 {
		fmt.Println("No entries yet — run an import first.")
		return
	}

	changed := 0
	for i := range entries {
		if newType := parser.ClassifyType(entries[i].RawLine); newType != entries[i].Type {
			fmt.Printf("  %s -> %s: %s\n", entries[i].Type, newType, entries[i].RawLine)
			entries[i].Type = newType
			changed++
		}
	}

	if changed == 0 {
		fmt.Println("Nothing to change — classification already up to date.")
		return
	}

	if err := database.Save(); err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	fmt.Printf("\nReclassified %d entr(ies). Study progress and everything else untouched.\n", changed)
}
