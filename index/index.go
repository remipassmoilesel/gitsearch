package index

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/pkg/errors"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/utils"
	"path"
	"time"
)

type Index struct {
	config config.Config
	// Data is split between a variable number of shards
	shards ShardGroup
	// Path to where is stored index data
	indexDataRoot string
	// All data is read from git repositories
	git utils.GitReader
	// State contains a list of git commits that have already been processed
	state IndexState
}

type IndexedFile struct {
	// File hash
	Hash string
	// Commit hash
	Commit string
	// File content
	Content string
	// File path
	Path string
	// File name
	Name string
}

type SearchResult struct {
	// Executed query
	Query string
	// Search duration in micro seconds
	TookUs int64
	// Files matching query
	Matches []SearchMatch
}

type SearchMatch struct {
	File      IndexedFile
	Fragments []string
}

const (
	OutputHtml = "html"
	OutputAnsi = "ansi"
)

type CleanOperationResult struct {
	TookMs int64
}

func NewIndex(config config.Config) (Index, error) {
	indexDataRoot := path.Join(config.DataRootPath, "index", config.Repository.Path)
	shardsRootPath := path.Join(indexDataRoot, "shards")

	shards := NewShardGroup(shardsRootPath, config.Search.Shards)

	gitReader, err := utils.NewGitReader(config.Repository.Path)
	if err != nil {
		return Index{}, err
	}

	index := Index{
		config:        config,
		shards:        shards,
		indexDataRoot: indexDataRoot,
		git:           gitReader,
	}

	err = index.initialize()
	if err != nil {
		return Index{}, err
	}

	return index, err
}

func (s *Index) initialize() error {
	_, err := s.shards.Initialize()
	if err != nil {
		return err
	}

	indexStatePath := path.Join(s.indexDataRoot, "gs-index-state.json")
	state, err := LoadIndexState(indexStatePath)
	if err != nil {
		return err
	}
	s.state = state

	return err
}

func (s *Index) Build() (BuildOperationResult, error) {
	builder := NewIndexBuilder(s)
	return builder.Build()
}

func (s *Index) Close() error {
	err := s.shards.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Index) Clean() (CleanOperationResult, error) {
	start := time.Now()

	err := s.shards.Clean()
	if err != nil {
		return CleanOperationResult{}, err
	}

	err = s.state.Clean()
	if err != nil {
		return CleanOperationResult{}, err
	}

	tookMs := time.Since(start).Milliseconds()
	response := CleanOperationResult{TookMs: tookMs}
	return response, err
}

func (s *Index) Search(textQuery string, size int, output string) (SearchResult, error) {
	start := time.Now()

	query := bleve.NewQueryStringQuery(textQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = size
	searchRequest.Fields = []string{"*"} // return all fields in results
	searchRequest.Highlight = bleve.NewHighlightWithStyle(output)
	searchResult, err := s.shards.searchIndex.Search(searchRequest)

	if err != nil {
		return *new(SearchResult), errors.Wrap(err, "search error")
	}

	var resultMatches []SearchMatch
	for _, hit := range searchResult.Hits {
		indexedFile := hitToIndexedFile(hit)

		fragments := []string{}
		for _, frags := range hit.Fragments {
			fragments = append(fragments, frags...)
		}

		match := SearchMatch{
			File:      indexedFile,
			Fragments: fragments,
		}

		resultMatches = append(resultMatches, match)
	}

	tookUs := time.Since(start).Microseconds()
	response := SearchResult{Query: textQuery, TookUs: tookUs, Matches: resultMatches}
	return response, err
}

func (s *Index) IsUpToDate() (bool, error) {
	hash, err := s.git.GetHeadHash()
	if err != nil {
		return false, err
	}

	return s.state.ContainsCommit(hash), nil
}

func (s *Index) DocumentCount() (uint64, error) {
	return s.shards.searchIndex.DocCount()
}

func hitToIndexedFile(document *search.DocumentMatch) IndexedFile {
	return IndexedFile{
		Hash:    document.Fields["Hash"].(string),
		Commit:  document.Fields["Commit"].(string),
		Content: document.Fields["Content"].(string),
		Path:    document.Fields["Path"].(string),
	}
}
