package gui

import (
	"fmt"
	"kotobaReviewer/entry"
	"kotobaReviewer/parser"
	"kotobaReviewer/quiz"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// showCard renders whatever should be on screen for the current position in
// the session: the next card (MCQ or plain), or the end-of-session summary
// once the queue is exhausted.
func showCard(w fyne.Window, content *fyne.Container, s *studySession) {
	e := s.current()
	if e == nil {
		if err := s.database.Save(); err != nil {
			content.Objects = []fyne.CanvasObject{widget.NewLabel("Error saving: " + err.Error())}
			content.Refresh()
			return
		}
		content.Objects = []fyne.CanvasObject{
			widget.NewLabel(fmt.Sprintf("Session complete! %d card(s) studied.", s.studied)),
			widget.NewButton("Back to menu", func() { ShowMenu(w, content) }),
		}
		content.Refresh()
		return
	}

	if e.Type == entry.Grammar {
		if forms, point, ok := parser.ExtractFormsAndPoint(e.RawLine); ok {
			showMCQCard(w, content, s, e, forms, point)
			return
		}
	}

	progress := widget.NewLabel(fmt.Sprintf("Card %d / %d", s.pos+1, len(s.queue)))

	front := e.Word
	if e.Type == entry.Grammar {
		front = e.RawLine
	}
	frontLabel := widget.NewLabel(front)
	frontLabel.Wrapping = fyne.TextWrapWord

	revealBtn := widget.NewButton("Reveal", func() {
		showAnswer(w, content, s, e)
	})

	content.Objects = []fyne.CanvasObject{progress, frontLabel, revealBtn}
	content.Refresh()
}

// showAnswer reveals everything on file for a vocab (or non-MCQ grammar)
// card, then offers the four SM-2 grading buttons.
func showAnswer(w fyne.Window, content *fyne.Container, s *studySession, e *entry.Entry) {
	progress := widget.NewLabel(fmt.Sprintf("Card %d / %d", s.pos+1, len(s.queue)))

	rawShown := e.Type == entry.Grammar
	front := e.Word
	if rawShown {
		front = e.RawLine
	}
	frontLabel := widget.NewLabel(front)
	frontLabel.Wrapping = fyne.TextWrapWord

	var lines []string
	if e.Reading != "" {
		lines = append(lines, "Reading: "+e.Reading)
	}
	if e.Meaning != "" {
		lines = append(lines, "Meaning: "+e.Meaning)
	}
	if e.Tag != "" {
		lines = append(lines, "Level: "+e.Tag)
	}
	for _, rf := range e.RelatedForms {
		if rf.Reading != "" {
			lines = append(lines, fmt.Sprintf("Related: %s (%s)", rf.Form, rf.Reading))
		} else {
			lines = append(lines, "Related: "+rf.Form)
		}
	}
	if notes := quiz.DisplayNotes(e.Notes, rawShown); notes != "" && notes != e.Meaning {
		lines = append(lines, "Notes: "+quiz.TruncateRunes(notes, 300))
	}
	if len(lines) == 0 {
		if rawShown {
			lines = append(lines, "(no additional notes)")
		} else {
			lines = append(lines, "(nothing recorded yet)")
		}
	}
	detailLabel := widget.NewLabel(strings.Join(lines, "\n"))
	detailLabel.Wrapping = fyne.TextWrapWord

	gradeRow := container.NewGridWithColumns(4,
		widget.NewButton("Again", func() { gradeAndNext(w, content, s, e, 0) }),
		widget.NewButton("Hard", func() { gradeAndNext(w, content, s, e, 3) }),
		widget.NewButton("Good", func() { gradeAndNext(w, content, s, e, 4) }),
		widget.NewButton("Easy", func() { gradeAndNext(w, content, s, e, 5) }),
	)

	content.Objects = []fyne.CanvasObject{
		progress, frontLabel, widget.NewSeparator(), detailLabel, gradeRow,
	}
	content.Refresh()
}
