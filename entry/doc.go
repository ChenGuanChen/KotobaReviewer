// Package entry defines a single vocab/grammar Entry and the behavior
// attached to it: constructing a fresh one, and grading spaced-repetition
// reviews over time.
//
// New and NewFromLine (create.go) are the two ways to construct an Entry;
// both generate a random ID via the unexported newID, so nothing outside
// this package ever needs to think about how IDs are made.
//
// GradeReview and IsDue (srs.go) implement a simplified SM-2-style spaced
// repetition schedule: GradeReview updates EaseFactor, IntervalDays, and
// NextReviewAt based on how well a review went, and IsDue reports whether
// an entry is due for review right now.
package entry
