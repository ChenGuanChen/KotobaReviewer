package quiz

import "strings"

// DisplayNotes strips out "[unparsed: ...]" bookkeeping segments once the
// full raw line has already been shown elsewhere (e.g. a grammar entry's
// front already showed RawLine) — at that point they're pure duplication,
// not new information. Genuine footnote text is always kept.
func DisplayNotes(notes string, suppressUnparsed bool) string {
	if !suppressUnparsed {
		return notes
	}
	parts := strings.Split(notes, " | ")
	var kept []string
	for _, p := range parts {
		if strings.HasPrefix(p, "[unparsed:") {
			continue
		}
		kept = append(kept, p)
	}
	return strings.Join(kept, " | ")
}

// TruncateRunes shortens s to at most max runes, appending "..." if cut.
func TruncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "..."
}
