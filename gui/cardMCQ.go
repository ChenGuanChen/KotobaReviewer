package gui

import (
	"fmt"
	"math/rand"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/ChenGuanChen/kotobaReviewer/entry"
	"github.com/ChenGuanChen/kotobaReviewer/quiz"
)

// showMCQCard renders a grammar entry that fits the "forms + point" shape
// as a multi-select checklist, mirroring the CLI's runGrammarMCQ — same
// distractor pool, same correctness check, just checkboxes instead of typed
// letters.
func showMCQCard(w fyne.Window, content *fyne.Container, s *studySession, e *entry.Entry, forms []string, point string) {
	progress := widget.NewLabel(fmt.Sprintf("Card %d / %d", s.pos+1, len(s.queue)))

	pointLabel := widget.NewLabel("+ " + point)
	pointLabel.Wrapping = fyne.TextWrapWord

	correct := map[string]bool{}
	for _, f := range forms {
		correct[f] = true
	}

	var candidates []string
	for _, f := range s.pool {
		if !correct[f] {
			candidates = append(candidates, f)
		}
	}
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	distractorsWanted := 4 - len(forms)
	if distractorsWanted < 1 {
		distractorsWanted = 1
	}
	if distractorsWanted > len(candidates) {
		distractorsWanted = len(candidates)
	}

	options := append([]string{}, forms...)
	options = append(options, candidates[:distractorsWanted]...)
	rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

	checks := widget.NewCheckGroup(options, nil)

	submitBtn := widget.NewButton("Submit", func() {
		chosen := map[string]bool{}
		for _, v := range checks.Selected {
			chosen[v] = true
		}

		quality := 1
		resultText := "Not quite — correct answer(s): " + strings.Join(sortedCopy(forms), ", ")
		if quiz.FormSetsEqual(chosen, correct) {
			quality = 4
			resultText = "Correct!"
		}

		resultLabel := widget.NewLabel(resultText)
		nextBtn := widget.NewButton("Next", func() {
			gradeAndNext(w, content, s, e, quality)
		})

		content.Objects = []fyne.CanvasObject{progress, pointLabel, checks, resultLabel, nextBtn}
		content.Refresh()
	})

	content.Objects = []fyne.CanvasObject{progress, pointLabel, checks, submitBtn}
	content.Refresh()
}
