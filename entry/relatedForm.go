package entry

// RelatedForm is an alternate written form of a word (e.g. a kanji
// variant), with its own optional reading.
type RelatedForm struct {
	Form    string `json:"form"`
	Reading string `json:"reading,omitempty"`
}
