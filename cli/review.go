package cli

import (
	"bufio"
	"fmt"
	"kotobaReviewer/db"
	"kotobaReviewer/entry"
	"kotobaReviewer/quiz"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// runReview picks a random subset of entries and quizzes you on them one at
// a time: show the word, you recall the meaning yourself, press Enter to
// reveal what's on file. No grading — that's `study`.
//
// args can include a count ("20") and/or a type filter ("vocab"/"grammar"),
// in any order — e.g. `review grammar`, `review 30 grammar`, `review 15`.
func RunReview(args []string) {
	database, err := db.Init()
	if err != nil {
		fmt.Println("Error loading:", err)
		return
	}

	count := 10
	var filterType entry.EntryType
	filtered := false
	for _, a := range args {
		switch strings.ToLower(a) {
		case "grammar":
			filterType, filtered = entry.Grammar, true
		case "vocab":
			filterType, filtered = entry.Vocab, true
		default:
			if n, err := strconv.Atoi(a); err == nil && n > 0 {
				count = n
			}
		}
	}

	entries := database.Database
	if filtered {
		var subset []entry.Entry
		for _, e := range entries {
			if e.Type == filterType {
				subset = append(subset, e)
			}
		}
		entries = subset
	}

	if len(entries) == 0 {
		fmt.Println("No entries yet — run an import first (or your filter matched nothing).")
		return
	}
	if count > len(entries) {
		count = len(entries)
	}

	order := rand.Perm(len(entries))[:count]
	reader := bufio.NewReader(os.Stdin)
	pool := quiz.BuildFormPool(database)

	fmt.Printf("Reviewing %d of %d entries. Press Enter to reveal each answer, 'q' then Enter to stop early.\n\n", count, len(entries))

	for i, idx := range order {
		e := entries[idx]

		fmt.Printf("[%d/%d] ", i+1, count)

		if e.Type == entry.Grammar {
			if matched, quit, _ := quiz.RunGrammarMCQ(reader, e, pool); matched {
				if quit {
					fmt.Println("\nStopped early.")
					return
				}
				fmt.Println()
				continue
			}
		}

		front := e.Word
		if e.Type == entry.Grammar {
			front = e.RawLine
		}

		fmt.Println(front)
		fmt.Print("  -- press Enter to reveal -- ")

		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) == "q" {
			fmt.Println("\nStopped early.")
			return
		}

		printAnswer(e, e.Type == entry.Grammar)
		fmt.Println()
	}

	fmt.Println("Session complete.")
}

// printAnswer shows whatever we actually have for an entry.
func printAnswer(e entry.Entry, rawAlreadyShown bool) {
	shown := false
	if e.Reading != "" {
		fmt.Println("  Reading:", e.Reading)
		shown = true
	}
	if e.Meaning != "" {
		fmt.Println("  Meaning:", e.Meaning)
		shown = true
	}
	if e.Tag != "" {
		fmt.Println("  Level:", e.Tag)
	}
	for _, rf := range e.RelatedForms {
		if rf.Reading != "" {
			fmt.Printf("  Related: %s (%s)\n", rf.Form, rf.Reading)
		} else {
			fmt.Printf("  Related: %s\n", rf.Form)
		}
		shown = true
	}
	if e.Notes != "" {
		if notes := quiz.DisplayNotes(e.Notes, rawAlreadyShown); notes != "" && notes != e.Meaning {
			fmt.Println("  Notes:", quiz.TruncateRunes(notes, 200))
			shown = true
		}
	}
	if !shown {
		if rawAlreadyShown {
			fmt.Println("  (no additional notes)")
		} else {
			fmt.Println("  (nothing recorded yet — raw line:", e.RawLine, ")")
		}
	}
}
