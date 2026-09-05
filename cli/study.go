package cli

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ChenGuanChen/kotobaReviewer/db"
	"github.com/ChenGuanChen/kotobaReviewer/entry"
	"github.com/ChenGuanChen/kotobaReviewer/quiz"
)

// runStudy is the real spaced-repetition mode: only cards that are actually
// due get shown, and after each one you rate your own recall so SM-2 can
// (re)schedule when it comes back. Unlike `review`, this command updates
// and saves vocab.json. args: an optional count, or "all" for every due
// card regardless of count.
func RunStudy(args []string) {
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

	var dueIdx []int
	for i, e := range entries {
		if e.IsDue() {
			dueIdx = append(dueIdx, i)
		}
	}
	if len(dueIdx) == 0 {
		fmt.Println("Nothing due right now — you're all caught up.")
		return
	}

	limit := 20
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "all" {
			limit = len(dueIdx)
		} else if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > len(dueIdx) {
		limit = len(dueIdx)
	}

	rand.Shuffle(len(dueIdx), func(i, j int) { dueIdx[i], dueIdx[j] = dueIdx[j], dueIdx[i] })
	dueIdx = dueIdx[:limit]

	reader := bufio.NewReader(os.Stdin)
	pool := quiz.BuildFormPool(database)

	fmt.Printf("%d card(s) due today (%d total due). Press Enter to reveal, then grade yourself. 'q' to stop early — progress so far is still saved.\n\n", limit, len(dueIdx))

	studied := 0
	for i, idx := range dueIdx {
		e := &entries[idx]

		fmt.Printf("[%d/%d] ", i+1, limit)

		if e.Type == entry.Grammar {
			if matched, quit, quality := quiz.RunGrammarMCQ(reader, *e, pool); matched {
				if quit {
					fmt.Println("\nStopped early.")
					goto done
				}
				e.GradeReview(quality)
				studied++
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
			goto done
		}

		printAnswer(*e, e.Type == entry.Grammar)

		quality := promptGrade(reader)
		if quality < 0 {
			fmt.Println("\nStopped early.")
			goto done
		}
		e.GradeReview(quality)
		studied++
		fmt.Println()
	}

done:
	if err := database.Save(); err != nil {
		fmt.Println("Error saving progress:", err)
		return
	}
	fmt.Printf("Saved. %d card(s) reviewed and rescheduled.\n", studied)
}

// promptGrade asks how well you recalled a card and maps the answer onto
// SM-2's 0-5 quality scale. Returns -1 if you typed 'q' to quit instead.
func promptGrade(reader *bufio.Reader) int {
	fmt.Print("  How did that go? 1=Again 2=Hard 3=Good 4=Easy -> ")
	line, _ := reader.ReadString('\n')
	switch strings.TrimSpace(line) {
	case "q":
		return -1
	case "1":
		return 0
	case "2":
		return 3
	case "4":
		return 5
	default:
		return 4
	}
}
