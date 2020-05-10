//go:generate mockgen -package mock -destination ../test/mock/mocks_Index.go gitlab.com/remipassmoilesel/gitsearch/index Index
package index

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/git_reader"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"path"
	"time"
)

type Index interface {
	Build() (domain.BuildOperationResult, error)
	BuildWith(options domain.BuildOptions) (domain.BuildOperationResult, error)
	Close() error
	Clean() (domain.CleanOperationResult, error)
	Search(textQuery string, size int, output string) (domain.SearchResult, error)
	FindDocumentById(hash string) (domain.IndexedFile, error)
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
	utils utils.Utils
}

const (
	OutputHtml = "html"
	OutputAnsi = "ansi"
)

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
		utils:         utils.NewUtils(),
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

func (s *IndexImpl) Build() (domain.BuildOperationResult, error) {
	return s.BuildWith(DefaultBuildOptions())
}

func (s *IndexImpl) BuildWith(options domain.BuildOptions) (domain.BuildOperationResult, error) {
	builder, err := NewIndexBuilder(s.config, s.state, s)
	if err != nil {
		return domain.BuildOperationResult{}, err
	}
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

func (s *IndexImpl) Clean() (domain.CleanOperationResult, error) {
	start := time.Now()

	err := s.shards.Clean()
	if err != nil {
		return domain.CleanOperationResult{}, err
	}

	err = s.state.Clean()
	if err != nil {
		return domain.CleanOperationResult{}, err
	}

	tookMs := time.Since(start).Milliseconds()
	response := domain.CleanOperationResult{TookMs: tookMs}
	return response, err
}

func (s *IndexImpl) Search(textQuery string, size int, output string) (domain.SearchResult, error) {
	start := time.Now()

	query := bleve.NewQueryStringQuery(textQuery)
	req := bleve.NewSearchRequest(query)
	req.Size = size
	req.SortBy([]string{"Date", "_score"})
	req.Fields = []string{"*"} // return all fields in results
	req.Highlight = bleve.NewHighlightWithStyle(output)

	searchResult, err := s.shards.searchIndex.Search(req)
	if err != nil {
		return *new(domain.SearchResult), errors.Wrap(err, "search error")
	}

	var resultMatches []domain.SearchMatch
	for _, hit := range searchResult.Hits {
		indexedFile, err := s.docMatchToIndexedFile(hit)
		if err != nil {
			return domain.SearchResult{}, errors.Wrap(err, "cannot parse document")
		}

		fragments := []string{}
		for _, frags := range hit.Fragments {
			fragments = append(fragments, frags...)
		}

		match := domain.SearchMatch{
			File:      indexedFile,
			Fragments: fragments,
		}

		resultMatches = append(resultMatches, match)
	}

	tookMs := time.Since(start).Milliseconds()
	response := domain.SearchResult{Query: textQuery, TookMs: tookMs, Matches: resultMatches}
	return response, err
}

func (s *IndexImpl) FindDocumentById(hash string) (domain.IndexedFile, error) {
	query := bleve.NewPrefixQuery(hash)
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"*"} // return all fields in results

	searchResult, err := s.shards.searchIndex.Search(req)
	if err != nil {
		return domain.IndexedFile{}, errors.Wrap(err, "cannot search document "+hash)
	}

	if len(searchResult.Hits) < 1 {
		return domain.IndexedFile{}, errors.New("not found " + hash)
	}

	return s.docMatchToIndexedFile(searchResult.Hits[0])
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

func (s *IndexImpl) docMatchToIndexedFile(document *search.DocumentMatch) (domain.IndexedFile, error) {
	date, err := s.utils.StringToDate(document.Fields["Date"].(string))
	if err != nil {
		return domain.IndexedFile{}, err
	}
	file := domain.IndexedFile{
		Hash:    document.Fields["Hash"].(string),
		Commit:  document.Fields["Commit"].(string),
		Content: document.Fields["Content"].(string),
		Name:    document.Fields["Name"].(string),
		Path:    document.Fields["Path"].(string),
		Date:    date,
	}
	return file, nil
}
