package cli

import (
	"bufio"
	"fmt"
	"kotobaReviewer/db"
	"kotobaReviewer/entry"
	"os"
	"strings"
)

// runAdd interactively prompts for a brand-new entry and appends it.
// Exists mainly because typing Japanese into Fyne's Entry widget doesn't
// work correctly with IME input (a known Fyne limitation, not specific to
// this app) — a normal terminal has no such problem.
func RunAdd() {
	database, err := db.Init()
	if err != nil {
		fmt.Println("Error loading:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Adding a new entry. Press Enter to leave any optional field blank.")

	typeStr := promptLine(reader, "Type (vocab/grammar) [vocab]: ")
	entryType := entry.Vocab
	if strings.ToLower(strings.TrimSpace(typeStr)) == "grammar" {
		entryType = entry.Grammar
	}

	word := strings.TrimSpace(promptLine(reader, "Word / pattern: "))
	for word == "" {
		fmt.Println("Word can't be empty.")
		word = strings.TrimSpace(promptLine(reader, "Word / pattern: "))
	}

	reading := strings.TrimSpace(promptLine(reader, "Reading (optional): "))
	meaning := strings.TrimSpace(promptLine(reader, "Meaning (optional): "))
	tag := strings.TrimSpace(promptLine(reader, "Level, e.g. N3 (optional): "))
	notes := strings.TrimSpace(promptLine(reader, "Notes (optional): "))

	newEntry, err := entry.NewManual(entryType, word, reading, meaning, tag, notes)
	if err != nil {
		fmt.Println("Error creating new entry in RunAdd():", err)
		return
	}

	database.Database = append(database.Database, *newEntry)
	if err := database.Save(); err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	fmt.Printf("Added %q. You now have %d entries.\n", newEntry.Word, len(database.Database))
}

func promptLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}
