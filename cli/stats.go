package cli

import (
	"fmt"
	"kotobaReviewer/db"
)

// runStats gives a quick overview of where your deck stands.
func RunStats() {
	database, err := db.Init()
	if err != nil {
		fmt.Println("Error loading:", err)
		return
	}

	entries := database.Database
	if len(entries) == 0 {
		fmt.Println("No entries yet.")
		return
	}

	s := database.ComputeStats()
	fmt.Printf("Total entries:      %d\n", s.Total)
	fmt.Printf("Never studied:      %d\n", s.NeverStudied)
	fmt.Printf("Due right now:      %d\n", s.DueNow)
	fmt.Printf("Due within 7 days:  %d\n", s.DueThisWeek)
	if s.HasEaseData {
		fmt.Printf("Average ease:       %.2f (2.50 = default; higher = easier for you)\n", s.AverageEase)
	}
}
