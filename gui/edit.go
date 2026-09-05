package gui

import (
	"fmt"
	"strings"

	"github.com/ChenGuanChen/kotobaReviewer/entry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// buildEditForm renders the Word/Reading/Meaning/Type/Level/Notes form for
// whichever entry is selected (or a blank one, for a new entry).
func buildEditForm(content *fyne.Container, st *manageState, query string) fyne.CanvasObject {
	var e entry.Entry
	isNew := st.editingIdx < 0
	if !isNew {
		e = st.database.Database[st.editingIdx]
	}

	wordEntry := widget.NewEntry()
	wordEntry.SetText(e.Word)
	wordEntry.SetPlaceHolder("Word / grammar pattern")

	readingEntry := widget.NewEntry()
	readingEntry.SetText(e.Reading)
	readingEntry.SetPlaceHolder("Reading (hiragana) — optional")

	meaningEntry := widget.NewEntry()
	meaningEntry.SetText(e.Meaning)
	meaningEntry.SetPlaceHolder("Meaning")

	tagEntry := widget.NewEntry()
	tagEntry.SetText(e.Tag)
	tagEntry.SetPlaceHolder("JLPT level, e.g. N3 — optional")

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetText(e.Notes)
	notesEntry.SetPlaceHolder("Notes / extra explanation — optional")

	typeSelect := widget.NewSelect([]string{string(entry.Vocab), string(entry.Grammar)}, nil)
	if e.Type != "" {
		typeSelect.SetSelected(string(e.Type))
	} else {
		typeSelect.SetSelected(string(entry.Vocab))
	}

	statusLabel := widget.NewLabel("")

	saveBtn := widget.NewButton("Save", func() {
		if strings.TrimSpace(wordEntry.Text) == "" {
			statusLabel.SetText("Word can't be empty.")
			return
		}

		updated := entry.Entry{
			Type:    entry.EntryType(typeSelect.Selected),
			Word:    wordEntry.Text,
			Reading: readingEntry.Text,
			Meaning: meaningEntry.Text,
			Tag:     tagEntry.Text,
			Notes:   notesEntry.Text,
		}

		if isNew {
			newEntry, err := entry.NewManual(
				updated.Type,
				updated.Word,
				updated.Reading,
				updated.Meaning,
				updated.Tag,
				updated.Notes,
			)
			if err != nil {
				fmt.Println("Error creating new entry in ShowManage: buildEditForm():", err)
				return
			}

			st.database.Database = append(st.database.Database, *newEntry)
			st.editingIdx = len(st.database.Database) - 1
		} else {
			// Preserve everything this form doesn't touch: ID, related
			// forms, the original raw line, and all SM-2 study progress.
			old := st.database.Database[st.editingIdx]
			updated.ID = old.ID
			updated.RelatedForms = old.RelatedForms
			updated.RawLine = old.RawLine
			updated.NeedsReview = old.NeedsReview
			updated.EaseFactor = old.EaseFactor
			updated.IntervalDays = old.IntervalDays
			updated.Repetitions = old.Repetitions
			updated.NextReviewAt = old.NextReviewAt
			updated.LastReviewedAt = old.LastReviewedAt
			st.database.Database[st.editingIdx] = updated
		}

		if err := st.database.Save(); err != nil {
			statusLabel.SetText("Error saving: " + err.Error())
			return
		}
		renderManage(content, st, query)
	})

	buttons := []fyne.CanvasObject{saveBtn}
	if !isNew {
		idx := st.editingIdx
		deleteBtn := widget.NewButton("Delete", func() {
			st.database.Database = append(st.database.Database[:idx], st.database.Database[idx+1:]...)
			st.editingIdx = -1
			if err := st.database.Save(); err != nil {
				statusLabel.SetText("Error saving: " + err.Error())
				return
			}
			renderManage(content, st, query)
		})
		buttons = append(buttons, deleteBtn)
	}

	title := "New entry"
	if !isNew {
		title = "Editing: " + e.Word
	}

	return container.NewVBox(
		widget.NewLabel(title),
		widget.NewForm(
			widget.NewFormItem("Type", typeSelect),
			widget.NewFormItem("Word", wordEntry),
			widget.NewFormItem("Reading", readingEntry),
			widget.NewFormItem("Meaning", meaningEntry),
			widget.NewFormItem("Level", tagEntry),
			widget.NewFormItem("Notes", notesEntry),
		),
		container.NewHBox(buttons...),
		statusLabel,
	)
}
