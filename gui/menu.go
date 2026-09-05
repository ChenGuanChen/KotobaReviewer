package gui

import (
	"fmt"

	"github.com/ChenGuanChen/kotobaReviewer/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// showMenu renders the Study tab's starting screen: how many cards are due,
// and a button to begin a session.
func ShowMenu(w fyne.Window, content *fyne.Container) {
	database, err := db.Init()
	if err != nil {
		content.Objects = []fyne.CanvasObject{widget.NewLabel("Error loading vocab.json: " + err.Error())}
		content.Refresh()
		return
	}

	dueCount := 0
	entries := database.Database
	for _, e := range entries {
		if e.IsDue() {
			dueCount++
		}
	}

	status := widget.NewLabel(fmt.Sprintf("%d entries total, %d due today", len(entries), dueCount))

	startBtn := widget.NewButton("Start Study Session", func() {
		session := newStudySession(database)
		if session.done() {
			content.Objects = []fyne.CanvasObject{
				widget.NewLabel("Nothing due right now — you're all caught up."),
				widget.NewButton("Back", func() { ShowMenu(w, content) }),
			}
			content.Refresh()
			return
		}
		showCard(w, content, session)
	})

	content.Objects = []fyne.CanvasObject{status, startBtn}
	content.Refresh()
}
