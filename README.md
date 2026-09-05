# Kotoba Trainer

A personal Japanese vocab and grammar trainer with spaced repetition, available
as both a terminal app and a desktop GUI. Import vocab/grammar lists exported
from a study doc, then review and study them with an SM-2-style spaced
repetition schedule.

## Features

- **Import** — parse a `.txt` export of a study document (vocab/grammar lines
  plus footnotes) straight into your deck.
- **Add** — add a single entry by hand from the terminal.
- **Review** — casual, ungraded practice that doesn't touch your schedule.
- **Study** — the real spaced-repetition session: shows only due cards,
  grades your recall, and reschedules with SM-2.
- **Stats** — deck overview: total entries, due now, due this week, average
  ease.
- **GUI** — a desktop app (built with [Fyne](https://fyne.io)) with Study,
  Manage, and Stats tabs.

## Prerequisites

- Go 1.27 or later (see `go.mod`).
- To build the **GUI**, Fyne requires CGO and a C compiler, plus OS graphics
  libraries (e.g. OpenGL dev headers on Linux). See
  [Fyne's getting started guide](https://docs.fyne.io/started/) for
  OS-specific setup — the CLI build doesn't need any of this.

## Building and running

```sh
# Terminal app
go build -o kotoba-cli ./cmd/cli
./kotoba-cli help

# Desktop app
go build -o kotoba-gui ./cmd/gui
./kotoba-gui
```

Or run either directly without a separate build step:
```sh
go run ./cmd/cli help
go run ./cmd/gui
```

## CLI usage

```
kotoba-cli                    Quick status: how many entries you have.

kotoba-cli import <file>      Re-parse an exported study doc (.txt) and
                               overwrite vocab.json with a fresh import.

kotoba-cli add                Interactively add a single new entry. Prefer
                               this over the GUI's Manage tab for typing
                               Japanese — Fyne's text fields don't support
                               IME input correctly, a normal terminal does.

kotoba-cli review [count] [type]
                               Casual, ungraded practice.
                                 count : how many cards this session (default 10)
                                 type  : "vocab" or "grammar" to filter

kotoba-cli study [count|all]  The real spaced-repetition session. Only shows
                               cards that are due, grades recall via SM-2,
                               and saves progress.

kotoba-cli stats              Overview of your deck.

kotoba-cli reclassify         Re-check vocab/grammar classification on
                               existing entries without touching progress.

kotoba-cli help               Show the full help message.
```

During a review/study session: Enter reveals the answer or moves on, `1`-`4`
grades recall (Again/Hard/Good/Easy), letters pick multiple-choice answers,
and `q` stops the session early (progress so far is still saved).

## Data storage

All entries are stored in a single `vocab.json` file, saved next to the
built executable (or next to `go.mod` when running via `go run`). There's no
separate database to set up.

## Project layout

```
cmd/cli/     entry point for the terminal app
cmd/gui/     entry point for the desktop app
cli/         terminal commands (import, add, review, study, stats, ...)
gui/         Fyne UI screens (study, manage, stats)
entry/       the Entry type and its spaced-repetition logic
db/          loading/saving the vocab.json database
parser/      turns a raw study-doc export into entries
quiz/        multiple-choice question generation for grammar review
```

## License

MIT License

Copyright (c) 2026 GUAN-CHEN, CHEN

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

