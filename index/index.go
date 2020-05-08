//go:generate mockgen -package mock -destination ../test/mock/mocks_Index.go gitlab.com/remipassmoilesel/gitsearch/index Index
package index

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/git_reader"
	"path"
	"time"
)

type Index interface {
	Build() (BuildOperationResult, error)
	BuildWith(options BuildOptions) (BuildOperationResult, error)
	Close() error
	Clean() (CleanOperationResult, error)
	Search(textQuery string, size int, output string) (SearchResult, error)
	FindDocumentById(hash string) (IndexedFile, error)
	IsUpToDate() (bool, error)
	DocumentCount() (uint64, error)
}

type IndexImpl struct {
	config config.Config
	// Data is split between a variable number of shards
	shards ShardGroup
	// Path to where is stored index data
	indexDataRoot string
	// All data is read from git repositories
	git git_reader.GitReader
	// State contains a list of git commits that have already been processed
	state IndexState
}

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
	shards := NewShardGroup(shardsRootPath, config.Index.Shards)

	gitReader, err := git_reader.NewGitReader(config.Repository.Path)
	if err != nil {
		return &IndexImpl{}, err
	}

	state, err := LoadIndexState(indexDataRoot)
	if err != nil {
		return &IndexImpl{}, errors.Wrap(err, "cannot initialize index")
	}

	index := IndexImpl{
		config:        config,
		shards:        shards,
		indexDataRoot: indexDataRoot,
		git:           gitReader,
		state:         state,
	}

	err = index.initialize()
	if err != nil {
		return &IndexImpl{}, err
	}

	return &index, err
}

func (s *IndexImpl) initialize() error {
	err := s.state.TryLock()
	if err != nil {
		return errors.Wrap(err, "cannot initialize index")
	}

	_, err = s.shards.Initialize()
	if err != nil {
		return errors.Wrap(err, "cannot initialize index")
	}

	return nil
}

func (s *IndexImpl) Build() (BuildOperationResult, error) {
	builder := NewIndexBuilder(s)
	return builder.Build(DefaultBuildOptions())
}

func (s *IndexImpl) BuildWith(options BuildOptions) (BuildOperationResult, error) {
	builder := NewIndexBuilder(s)
	return builder.Build(options)
}

func (s *IndexImpl) Close() error {
	defer func() {
		err := s.state.Unlock()
		if err != nil {
			fmt.Println(errors.Wrap(err, "cannot unlock index state"))
		}
	}()

	err := s.state.Write()
	if err != nil {
		return errors.Wrap(err, "cannot unlock index state")
	}

	return s.shards.Close()
}

func (s *IndexImpl) Clean() (CleanOperationResult, error) {
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

func (s *IndexImpl) Search(textQuery string, size int, output string) (SearchResult, error) {
	start := time.Now()

	query := bleve.NewQueryStringQuery(textQuery)
	req := bleve.NewSearchRequest(query)
	req.Size = size
	req.SortBy([]string{"Date", "_score"})
	req.Fields = []string{"*"} // return all fields in results
	req.Highlight = bleve.NewHighlightWithStyle(output)
	searchResult, err := s.shards.searchIndex.Search(req)

	if err != nil {
		return *new(SearchResult), errors.Wrap(err, "search error")
	}

	var resultMatches []SearchMatch
	for _, hit := range searchResult.Hits {
		indexedFile, err := docMatchToIndexedFile(hit)
		if err != nil {
			return SearchResult{}, errors.Wrap(err, "cannot parse document")
		}

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

	tookMs := time.Since(start).Milliseconds()
	response := SearchResult{Query: textQuery, TookMs: tookMs, Matches: resultMatches}
	return response, err
}

func (s *IndexImpl) FindDocumentById(hash string) (IndexedFile, error) {
	query := bleve.NewPrefixQuery(hash)
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"*"} // return all fields in results
	searchResult, err := s.shards.searchIndex.Search(req)

	if err != nil {
		return IndexedFile{}, errors.Wrap(err, "cannot search document "+hash)
	}

	if len(searchResult.Hits) < 1 {
		return IndexedFile{}, errors.New("not found " + hash)
	}

	return docMatchToIndexedFile(searchResult.Hits[0])
}

func (s *IndexImpl) IsUpToDate() (bool, error) {
	commit, err := s.git.GetHeadHash()
	if err != nil {
		return false, err
	}

	return s.state.ContainsCommit(commit), nil
}

func (s *IndexImpl) DocumentCount() (uint64, error) {
	return s.shards.searchIndex.DocCount()
}

func docMatchToIndexedFile(document *search.DocumentMatch) (IndexedFile, error) {
	date, err := stringToDate(document.Fields["Date"].(string))
	if err != nil {
		return IndexedFile{}, err
	}
	file := IndexedFile{
		Hash:    document.Fields["Hash"].(string),
		Commit:  document.Fields["Commit"].(string),
		Content: document.Fields["Content"].(string),
		Name:    document.Fields["Name"].(string),
		Path:    document.Fields["Path"].(string),
		Date:    date,
	}
	return file, nil
}

func stringToDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "cannot parse date "+dateStr)
	}
	return date, nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
