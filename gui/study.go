package gui

import (
	"math/rand"
	"sort"

	"fyne.io/fyne/v2"

	"github.com/ChenGuanChen/kotobaReviewer/db"
	"github.com/ChenGuanChen/kotobaReviewer/entry"
	"github.com/ChenGuanChen/kotobaReviewer/quiz"
)

// studySession holds one Study-tab session's state: which entries are due,
// where we are in the queue, and the shared distractor pool for grammar
// MCQs. allEntries is the FULL loaded deck (not just the due ones) because
// we save the whole thing back at the end, same as the CLI does.
type studySession struct {
	database db.VocabDB
	queue    []int // indices into allEntries
	pos      int
	pool     []string
	studied  int
}

func newStudySession(database db.VocabDB) *studySession {
	var due []int
	for i, e := range database.Database {
		if e.IsDue() {
			due = append(due, i)
		}
	}
	rand.Shuffle(len(due), func(i, j int) { due[i], due[j] = due[j], due[i] })

	const sessionLimit = 20
	if len(due) > sessionLimit {
		due = due[:sessionLimit]
	}

	return &studySession{
		database: database,
		queue:    due,
		pool:     quiz.BuildFormPool(database),
	}
}

func (s *studySession) current() *entry.Entry {
	if s.pos >= len(s.queue) {
		return nil
	}
	return &s.database.Database[s.queue[s.pos]]
}

func (s *studySession) advance() {
	s.pos++
}

func (s *studySession) done() bool {
	return s.pos >= len(s.queue)
}

func gradeAndNext(w fyne.Window, content *fyne.Container, s *studySession, e *entry.Entry, quality int) {
	e.GradeReview(quality)
	s.studied++
	s.advance()
	showCard(w, content, s)
}

func sortedCopy(forms []string) []string {
	out := append([]string{}, forms...)
	sort.Strings(out)
	return out
}
