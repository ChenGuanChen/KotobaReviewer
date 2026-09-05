package db

import (
	"time"

	"github.com/ChenGuanChen/kotobaReviewer/entry"
)

// VocabDB is the whole saved collection: every entry, plus metadata about
// the last time it was synced to a remote copy.
type VocabDB struct {
	Database []entry.Entry `json:"database"`
	LastSync time.Time     `json:"lastSync,omitempty"`
	SyncUrl  string        `json:"syncUrl,omitempty"`
}
