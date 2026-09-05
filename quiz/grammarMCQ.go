package quiz

import (
	"bufio"
	"fmt"
	"kotobaReviewer/entry"
	"kotobaReviewer/parser"
	"math/rand"
	"sort"
	"strings"
)

// runGrammarMCQ quizzes on one entry as multi-select multiple choice: given
// the grammar point, pick every connecting form that's valid before it.
// matched=false means this entry doesn't fit the MCQ shape, so the caller
// should fall back to the plain free-recall format instead. quality is an
// SM-2 quality score (4 if correct, 1 if not) for callers doing real spaced
// repetition; plain `review` sessions can just ignore it.
func RunGrammarMCQ(reader *bufio.Reader, e entry.Entry, pool []string) (matched bool, quit bool, quality int) {
	forms, point, ok := parser.ExtractFormsAndPoint(e.RawLine)
	if !ok {
		return false, false, 0
	}

	correct := map[string]bool{}
	for _, f := range forms {
		correct[f] = true
	}

	var candidates []string
	for _, f := range pool {
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

	fmt.Printf("+ %s\n", point)
	const labels = "abcdefgh"
	labelOf := map[string]string{}
	for i, opt := range options {
		label := string(labels[i])
		labelOf[label] = opt
		fmt.Printf("  (%s) %s\n", label, opt)
	}
	fmt.Print("  Which form(s) fit? e.g. \"a,c\" (or 'q' to stop) -> ")

	line, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(line)
	if trimmed == "q" {
		return true, true, 0
	}

	chosen := map[string]bool{}
	for _, tok := range strings.FieldsFunc(strings.ToLower(trimmed), func(r rune) bool {
		return r == ',' || r == ' '
	}) {
		if opt, exists := labelOf[tok]; exists {
			chosen[opt] = true
		}
	}

	if FormSetsEqual(chosen, correct) {
		fmt.Println("  Correct!")
		quality = 4
	} else {
		correctList := append([]string{}, forms...)
		sort.Strings(correctList)
		fmt.Println("  Not quite — correct answer(s):", strings.Join(correctList, ", "))
		quality = 1
	}
	return true, false, quality
}
