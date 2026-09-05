// Package db loads and saves a VocabDB — the whole collection of vocab
// entries — as a single JSON file on disk.
//
// Init reads that file into memory (or returns an empty VocabDB if it
// doesn't exist yet), and Save writes it back:
//
//	database, err := db.Init()
//	...
//	err = database.Save()
//
// dataFilePath (dataFilePath.go) works out where that JSON file actually
// lives: next to the built executable normally, or next to go.mod when
// running via `go run`/`go build` during development, so the data file
// doesn't get written into a temporary go-build directory.
//
// ComputeStats (stats.go) derives review-progress numbers — due counts,
// average ease — from the entries already loaded, without touching disk.
package db
