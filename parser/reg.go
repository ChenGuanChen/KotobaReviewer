package parser

import "regexp"

// A footnote line looks like "[a]some text..." at the very start of the line.
var footnoteStartRe = regexp.MustCompile(`^\[([a-z]+)\](.*)$`)

// A word/related-form can carry a trailing footnote marker glued onto it,
// e.g. "兼ねる[a]" or "筆（ふで）を折る[c]".
var trailingMarkerRe = regexp.MustCompile(`^(.*)\[([a-z]+)\]$`)

// JLPT level tokens: N3, N2, or full-width Ｎ３, Ｎ２, etc.
var levelRe = regexp.MustCompile(`^[NＮ][1-3１-３]$`)

// Plain page-number tokens, e.g. "47", "117".
var pageNumRe = regexp.MustCompile(`^[0-9]+$`)

// Characters that, when a line starts with one of them, mean "this is a
// grammar pattern" rather than a vocab word. Real Japanese words are never
// written starting with a bare romaji v/V/N or a tilde — only grammar
// placeholders (V/N as stand-ins for "verb"/"noun") and pattern markers are.
const grammarStartChars = "vVｖＶnNｎＮ～"

// Grammar meta-terms that show up *inside* a line rather than at the very
// start — e.g. "（動詞連用）がたい" opens with a parenthetical, not a bare
// v/V/N, so grammarStartChars alone misses it. Checked against the full
// dataset before adding: every vocab-labeled entry containing one of these
// really is a grammar note, and no genuinely-vocab entry (e.g. "（何階）建
// て", which just has a usage-hint parenthetical) contains them.
var grammarJargon = []string{"動詞連用", "動詞", "い形容詞", "な形容詞", "～"}

// formsAndPointRe matches grammar raw lines shaped like
// "formA/formB/... + point..." — connecting forms (slash-separated) on the
// left, a single grammar point on the right. Only lines with a SPACE on
// both sides of the +/＋ match; that's what distinguishes this clean shape
// from unrelated text that happens to contain a "+" with no surrounding
// space (like "N+からして...", a different kind of pattern entirely).
var formsAndPointRe = regexp.MustCompile(`^(.+?)\s[+＋]\s(.+)$`)
