package quiz

import (
	"github.com/ChenGuanChen/kotobaReviewer/db"
	"github.com/ChenGuanChen/kotobaReviewer/entry"
	"github.com/ChenGuanChen/kotobaReviewer/parser"
)

func BuildFormPool(database db.VocabDB) []string {
	seen := map[string]bool{}
	var pool []string
	for _, e := range database.Database {
		if e.Type != entry.Grammar {
			continue
		}
		forms, _, ok := parser.ExtractFormsAndPoint(e.RawLine)
		if !ok {
			continue
		}
		for _, f := range forms {
			if !seen[f] {
				seen[f] = true
				pool = append(pool, f)
			}
		}
	}
	if !seen["~"] {
		seen["~"] = true
		pool = append(pool, "~")
	}
	return pool
}
