package domain

import "time"

type CleanOperationResult struct {
	TookMs int64
}

type BuildOperationResult struct {
	TookSeconds  float64
	Files        int
	TotalFiles   int
	OldestCommit time.Time
}
