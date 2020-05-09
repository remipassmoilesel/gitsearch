package domain

import "time"

type IndexedFile struct {
	Hash string
	// Commit hash
	Commit string
	// Date of youngest commit containing this version of the file
	Date    time.Time
	Content string
	Name    string
	Path    string
}
