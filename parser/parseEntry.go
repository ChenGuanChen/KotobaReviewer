package parser

import (
	"fmt"
	"kotobaReviewer/entry"
	"strings"
	"unicode"
)

type relatedNote struct{ form, marker string }

func parseEntryLine(line string, footnotes map[string]string) (*entry.Entry, error) {
	newEntry, err := entry.NewFromLine(line)
	if err != nil {
		return nil, fmt.Errorf("parseEntryLine: %w", err)
	}

	// Normalize full-width spaces to regular ones so Fields() splits cleanly.
	normalized := strings.ReplaceAll(line, "\u3000", " ")
	tokens := strings.Fields(normalized)
	if len(tokens) == 0 {
		newEntry.NeedsReview = true
		return newEntry, nil
	}

	// --- classify vocab vs grammar ---
	newEntry.Type = ClassifyType(line)

	// --- main word (+ possible footnote marker) ---
	word, mainMarker := stripMarker(tokens[0])
	newEntry.Word = word
	tokens = tokens[1:]

	// --- reading, if the next token is kana-only ---
	if len(tokens) > 0 && !strings.HasPrefix(tokens[0], "＊") && !strings.HasPrefix(tokens[0], "*") && isKanaOnly(tokens[0]) {
		newEntry.Reading = tokens[0]
		tokens = tokens[1:]
	}

	// --- everything else: related forms, level, page ref, or unknown ---
	var relatedMarkers []relatedNote
	var leftover []string
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch {
		case strings.HasPrefix(t, "＊") || strings.HasPrefix(t, "*"):
			form, fMarker := stripMarker(strings.TrimPrefix(strings.TrimPrefix(t, "＊"), "*"))
			if fMarker != "" {
				relatedMarkers = append(relatedMarkers, relatedNote{form, fMarker})
			}
			rf := entry.RelatedForm{Form: form}
			if i+1 < len(tokens) && isKanaOnly(tokens[i+1]) {
				rf.Reading = tokens[i+1]
				i++
			}
			newEntry.RelatedForms = append(newEntry.RelatedForms, rf)

		case levelRe.MatchString(t):
			newEntry.Tag = normalizeLevel(t)

		case pageNumRe.MatchString(t):
			// page reference from your study book; not modeled explicitly,
			// safe to drop. Flip to `leftover = append(leftover, t)` if you'd
			// rather keep it visible for review instead.

		default:
			leftover = append(leftover, t)
		}
	}

	if len(leftover) > 0 {
		newEntry.NeedsReview = true
		newEntry.Notes = appendNote(newEntry.Notes, "[unparsed: "+strings.Join(leftover, " ")+"]")
	}

	// --- attach footnote text (your doc's "comments") ---
	// Only the headword's own footnote can seed Meaning — a footnote glued to
	// a related form describes that related form, not the entry itself, so
	// it only ever goes into Notes (labeled with which form it explains).
	if mainMarker != "" {
		if text, ok := footnotes[mainMarker]; ok {
			newEntry.Notes = appendNote(newEntry.Notes, text)
			if len([]rune(text)) <= 60 {
				newEntry.Meaning = text
			}
		}
	}
	for _, rm := range relatedMarkers {
		if text, ok := footnotes[rm.marker]; ok {
			newEntry.Notes = appendNote(newEntry.Notes, "("+rm.form+") "+text)
		}
	}

	return newEntry, nil
}

func stripMarker(token string) (clean string, marker string) {
	if m := trailingMarkerRe.FindStringSubmatch(token); m != nil {
		return m[1], m[2]
	}
	return token, ""
}

func normalizeLevel(t string) string {
	replacer := strings.NewReplacer("Ｎ", "N", "１", "1", "２", "2", "３", "3")
	return replacer.Replace(t)
}

func isKanaOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		isHiragana := r >= 0x3040 && r <= 0x309F
		isKatakana := r >= 0x30A0 && r <= 0x30FF
		if !isHiragana && !isKatakana && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func appendNote(existing, addition string) string {
	if existing == "" {
		return addition
	}
	return existing + " | " + addition
}
