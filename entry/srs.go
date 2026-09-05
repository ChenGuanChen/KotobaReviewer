package entry

import "time"

// GradeReview records the result of reviewing e right now, where quality
// is how well the review went on a 0-5 scale (below 3 counts as a lapse
// and resets the interval). It updates EaseFactor, IntervalDays,
// Repetitions, and schedules NextReviewAt accordingly.
func (e *Entry) GradeReview(quality int) {
	if e.EaseFactor == 0 {
		e.EaseFactor = 2.5
	}

	if quality < 3 {
		e.Repetitions = 0
		e.IntervalDays = 1
	} else {
		switch e.Repetitions {
		case 0:
			e.IntervalDays = 1
		case 1:
			e.IntervalDays = 6
		default:
			e.IntervalDays = int(float64(e.IntervalDays)*e.EaseFactor + 0.5)
		}
		e.Repetitions++
	}

	q := float64(quality)
	// SM-2
	e.EaseFactor += 0.1 - (5-q)*(0.08+(5-q)*0.02)
	if e.EaseFactor < 1.3 {
		e.EaseFactor = 1.3
	}

	now := time.Now()
	e.LastReviewedAt = now.Format(DateLayout)
	e.NextReviewAt = now.AddDate(0, 0, e.IntervalDays).Format(DateLayout)
}

// IsDue reports whether e is due for review: true if it's never been
// scheduled, its NextReviewAt is unparseable, or that date has arrived.
func (e *Entry) IsDue() bool {
	if e.NextReviewAt == "" {
		return true
	}
	due, err := time.Parse(DateLayout, e.NextReviewAt)
	if err != nil {
		return true
	}
	today := time.Now().Truncate(24 * time.Hour)
	return !due.After(today)
}
