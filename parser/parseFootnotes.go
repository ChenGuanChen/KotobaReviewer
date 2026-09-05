package parser

import "strings"

// parseFootnotes groups the tail of the document into marker -> full text.
// A footnote can span multiple lines; it continues until the next line that
// starts a new [marker].
func parseFootnotes(lines []string) map[string]string {
	notes := map[string]string{}
	currentKey := ""
	var currentText []string

	flush := func() {
		if currentKey != "" {
			notes[currentKey] = strings.TrimSpace(strings.Join(currentText, " "))
		}
	}

	for _, l := range lines {
		if m := footnoteStartRe.FindStringSubmatch(l); m != nil {
			flush()
			currentKey = m[1]
			currentText = []string{strings.TrimSpace(m[2])}
		} else if currentKey != "" {
			trimmed := strings.TrimSpace(l)
			if trimmed != "" {
				currentText = append(currentText, trimmed)
			}
		}
	}
	flush()
	return notes
}
