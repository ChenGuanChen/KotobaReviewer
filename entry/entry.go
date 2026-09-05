package entry

// Entry is a single vocab or grammar item, plus its spaced-repetition
// review state. Zero-value fields are omitted from JSON where marked
// `omitempty`, since most fields are optional depending on how the entry
// was created and how much it's been reviewed so far.
type Entry struct {
	ID           string        `json:"id"`
	Type         EntryType     `json:"type"`
	Word         string        `json:"word"`
	Reading      string        `json:"reading,omitempty"`
	Meaning      string        `json:"meaning,omitempty"`
	Notes        string        `json:"notes,omitempty"`
	RelatedForms []RelatedForm `json:"related_forms,omitempty"`
	Tag          string        `json:"tag,omitempty"`
	RawLine      string        `json:"raw_line,omitempty"`
	NeedsReview  bool          `json:"needs_review,omitempty"`

	EaseFactor     float64 `json:"ease_factor,omitempty"`
	IntervalDays   int     `json:"interval_days,omitempty"`
	Repetitions    int     `json:"repetitions,omitempty"`
	NextReviewAt   string  `json:"next_review_at,omitempty"`
	LastReviewedAt string  `json:"last_reviewed_at,omitempty"`
}
