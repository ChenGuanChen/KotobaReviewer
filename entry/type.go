package entry

// EntryType distinguishes a plain vocab word from a grammar pattern.
type EntryType string

const (
	Vocab   EntryType = "vocab"
	Grammar EntryType = "grammar"
)

// DateLayout is the format used for every stored date/time string on an
// Entry (NextReviewAt, LastReviewedAt), so parsing and formatting stay
// consistent across the package.
const DateLayout = "2006-01-02"
