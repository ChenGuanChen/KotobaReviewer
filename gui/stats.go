package gui

import (
	"fmt"
	"kotobaReviewer/db"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// showStats renders a one-time snapshot on the Stats tab. It doesn't
// auto-refresh if you study in the same window session — restart the app
// (or switch tabs back and forth once we wire up a refresh hook) to see
// updated numbers. Fine for v1.
func ShowStats(content *fyne.Container) {
	database, err := db.Init()
	if err != nil {
		content.Objects = []fyne.CanvasObject{widget.NewLabel("Error loading vocab.json: " + err.Error())}
		content.Refresh()
		return
	}

	entries := database.Database
	if len(entries) == 0 {
		content.Objects = []fyne.CanvasObject{widget.NewLabel("No entries yet.")}
		content.Refresh()
		return
	}

	s := database.ComputeStats()
	lines := []string{
		fmt.Sprintf("Total entries:      %d", s.Total),
		fmt.Sprintf("Never studied:      %d", s.NeverStudied),
		fmt.Sprintf("Due right now:      %d", s.DueNow),
		fmt.Sprintf("Due within 7 days:  %d", s.DueThisWeek),
	}
	if s.HasEaseData {
		lines = append(lines, fmt.Sprintf("Average ease:       %.2f (2.50 = default; higher = easier for you)", s.AverageEase))
	}

	label := widget.NewLabel(strings.Join(lines, "\n"))
	content.Objects = []fyne.CanvasObject{label}
	content.Refresh()
}
