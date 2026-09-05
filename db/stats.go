package db

import (
	"time"

	"github.com/ChenGuanChen/kotobaReviewer/entry"
)

// Stats summarizes review progress across a VocabDB at a point in time.
type Stats struct {
	Total        int
	NeverStudied int
	DueNow       int
	DueThisWeek  int
	AverageEase  float64
	HasEaseData  bool
}

// ComputeStats tallies review-progress numbers across every entry in db:
// how many are due now or within the coming week, how many have never
// been studied, and the average ease factor among entries that have.
func (db *VocabDB) ComputeStats() Stats {
	var s Stats
	s.Total = len(db.Database)

	var totalEase float64
	var reviewedCount int
	weekFromNow := time.Now().AddDate(0, 0, 7)

	for _, e := range db.Database {
		if e.NextReviewAt == "" {
			s.NeverStudied++
		}
		if e.IsDue() {
			s.DueNow++
		}
		if e.NextReviewAt != "" {
			if due, err := time.Parse(entry.DateLayout, e.NextReviewAt); err == nil && due.Before(weekFromNow) {
				s.DueThisWeek++
			}
		}
		if e.EaseFactor > 0 {
			totalEase += e.EaseFactor
			reviewedCount++
		}
	}

	if reviewedCount > 0 {
		s.AverageEase = totalEase / float64(reviewedCount)
		s.HasEaseData = true
	}
	return s
}
