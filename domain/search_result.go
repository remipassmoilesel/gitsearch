package domain

type SearchResult struct {
	// Executed query
	Query string
	// Search duration in milli seconds
	TookMs int64
	// Files matching query
	Matches []SearchMatch
}

type SearchMatch struct {
	File      IndexedFile
	Fragments []string
}
