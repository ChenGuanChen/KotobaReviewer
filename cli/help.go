package cli

import "fmt"

func PrintHelp() {
	fmt.Println(`Kotoba Trainer — commands:

  (no command)              Quick status: how many entries you have.

  import <file>             Re-parse an exported Google Doc (.txt) and
                             overwrite vocab.json with a fresh import.

  add                       Interactively add a single new entry. Prefer
                             this over the GUI's Manage tab for typing
                             Japanese — Fyne's text fields don't support IME
                             input correctly (a known Fyne limitation), a
                             normal terminal does.

  review [count] [type]     Casual, ungraded practice. Doesn't touch your
                             spaced-repetition schedule.
                               count : how many cards this session (default 10)
                               type  : "vocab" or "grammar" to filter
                             Examples:
                               review
                               review 20
                               review grammar
                               review vocab 30

  study [count|all]         The real spaced-repetition session. Only shows
                             cards that are actually due, grades your recall,
                             reschedules via SM-2, and saves progress.
                               count : how many due cards this session (default 20)
                               all   : every due card, regardless of count
                             Examples:
                               study
                               study 10
                               study all

  stats                     Overview of your deck: total, never studied,
                             due now, due this week, average ease.

  reclassify                Re-check vocab/grammar classification on your
                             existing entries and fix any that are wrong,
                             WITHOUT touching study progress or IDs. Safe
                             to re-run any time.

  help                      Show this message.

During review/study sessions:
  Enter    reveal the answer (vocab) / move on (after an MCQ)
  1-4      grade yourself during study: 1=Again 2=Hard 3=Good 4=Easy
  a, b...  pick multiple-choice answers, e.g. "a,c"
  q        stop the session early (study still saves progress so far)`)
}
