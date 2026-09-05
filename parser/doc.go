// Package parser turns a raw study-doc text file into a db.VocabDB.
//
// The entry point is ParseDocFile: give it a file path, expecting .txt file
// exported by Google Doc, get back a populated database
// (or an error naming the offending line).
//
//	database, err := parser.ParseDocFile("vocab.txt")
//
// Internally, ParseDocFile splits the file into two regions — vocab/grammar
// lines followed by a footnotes section — and hands each region off to a
// dedicated step:
//
//   - parseFootnotes (parseFootnotes.go) collects the "[a] some text..."
//     tail of the document into a marker -> text map.
//   - parseEntryLine (parseEntry.go) turns one line into an *entry.Entry,
//     pulling out the word, reading, related forms, JLPT level, and
//     attaching any footnote text found via parseFootnotes.
//   - ClassifyType (classifier.go) decides whether a line is vocab or
//     grammar. It's exported separately from parseEntryLine because the
//     `reclassify` command also calls it directly, to patch the Type on
//     entries already saved in vocab.json without re-parsing them.
//
// reg.go holds the regexes and character sets the above steps share, with
// notes on what real-world line shapes each one is guarding against.
package parser
