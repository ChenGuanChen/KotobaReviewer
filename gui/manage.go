package gui

import (
	"fmt"
	"strings"

	"github.com/ChenGuanChen/kotobaReviewer/db"
	"github.com/ChenGuanChen/kotobaReviewer/quiz"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// manageState holds everything the Manage tab needs across rebuilds.
// editingIdx is an index into entries, or -1 for "composing a new entry".
type manageState struct {
	database   db.VocabDB
	editingIdx int
}

func ShowManage(content *fyne.Container) {
	database, err := db.Init()
	if err != nil {
		content.Objects = []fyne.CanvasObject{widget.NewLabel("Error loading vocab.json: " + err.Error())}
		content.Refresh()
		return
	}

	st := &manageState{database: database, editingIdx: -1}
	renderManage(content, st, "")
}

// renderManage draws the search box, matching results, and the edit form
// for whichever entry (or blank new entry) is currently selected. The
// whole tab rebuilds on any change — see the note above about why.
func renderManage(content *fyne.Container, st *manageState, query string) {
	searchEntry := widget.NewEntry()
	searchEntry.SetText(query)
	searchEntry.SetPlaceHolder("Search by word, reading, or raw text...")
	searchEntry.OnChanged = func(q string) {
		renderManage(content, st, q)
	}

	newBtn := widget.NewButton("+ New Entry", func() {
		st.editingIdx = -1
		renderManage(content, st, query)
	})

	var resultObjs []fyne.CanvasObject
	if strings.TrimSpace(query) == "" {
		resultObjs = append(resultObjs, widget.NewLabel(fmt.Sprintf("%d entries — type to search", len(st.database.Database))))
	} else {
		q := strings.ToLower(query)
		matchCount := 0
		for i, e := range st.database.Database {
			haystack := strings.ToLower(e.Word + " " + e.Reading + " " + e.RawLine)
			if !strings.Contains(haystack, q) {
				continue
			}
			matchCount++
			if matchCount > 30 {
				resultObjs = append(resultObjs, widget.NewLabel("... more matches, narrow your search ..."))
				break
			}
			idx := i // capture per-iteration for the closure below
			label := e.Word
			if label == "" {
				label = quiz.TruncateRunes(e.RawLine, 40)
			}
			resultObjs = append(resultObjs, widget.NewButton(label, func() {
				st.editingIdx = idx
				renderManage(content, st, query)
			}))
		}
		if matchCount == 0 {
			resultObjs = append(resultObjs, widget.NewLabel("No matches."))
		}
	}
	resultsScroll := container.NewScroll(container.NewVBox(resultObjs...))

	form := buildEditForm(content, st, query)

	layout := container.NewVBox(
		searchEntry,
		newBtn,
		resultsScroll,
		widget.NewSeparator(),
		form,
	)

	content.Objects = []fyne.CanvasObject{layout}
	content.Refresh()
}
