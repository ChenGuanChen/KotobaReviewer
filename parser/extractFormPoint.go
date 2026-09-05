package parser

import "strings"

// ExtractFormsAndPoint tries to pull "connecting forms" + "grammar point"
// out of a grammar entry's raw line. Not every grammar line fits this shape
// (many are single fused phrases) — ok=false means "no MCQ for this one,"
// not an error.
func ExtractFormsAndPoint(raw string) (forms []string, point string, ok bool) {
	m := formsAndPointRe.FindStringSubmatch(raw)
	if m == nil {
		return nil, "", false
	}
	left := strings.TrimSpace(m[1])
	point = strings.TrimSpace(m[2])
	if left == "" || point == "" {
		return nil, "", false
	}
	for _, f := range strings.FieldsFunc(left, func(r rune) bool { return r == '/' || r == '／' }) {
		if f = strings.TrimSpace(f); f != "" {
			forms = append(forms, normalizeFormToken(f))
		}
	}
	if len(forms) == 0 {
		return nil, "", false
	}
	return forms, point, true
}

// normalizeFormToken collapses halfwidth/fullwidth variants of the v/V and
// n/N placeholder letters (v原, Ｖた, ｖる, Nの, Ｎである, ...) down to one
// canonical form. These all mean the same thing — "verb" / "noun" — and the
// width/case difference is just inconsistent IME input in the original
// notes, not a meaningful distinction.
func normalizeFormToken(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'v', 'V', 'ｖ', 'Ｖ':
			b.WriteRune('V')
		case 'n', 'N', 'ｎ', 'Ｎ':
			b.WriteRune('N')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
