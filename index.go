package main

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type Index struct {
	config        Config
	git           GitHelper
	indexDataPath string
	internalIndex bleve.Index
	state         IndexPersistedState
}

type IndexedFile struct {
	// See function below
	Id string
	// File hash
	Hash string
	// Commit hash
	Commit string
	// File content
	Content string
	// File path
	Path string
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

type CleanOperationResult struct {
	TookMillis int64
}

type BuildOperationResult struct {
	TookSeconds float64
	Files       int
}

func NewIndex(config Config) (Index, error) {
	indexDataPath := path.Join(config.DataRootPath, "index", config.Repository.Path)
	indexStatePath := path.Join(indexDataPath, "gs-index-state.json")

	state, err := LoadPersistedState(indexStatePath)
	if err != nil {
		return Index{}, err
	}

	index := Index{
		config:        config,
		git:           GitHelper{},
		indexDataPath: indexDataPath,
		state:         state,
	}

	err = index.initialize()
	return index, err
}

func (s *Index) initialize() error {
	var err error
	if _, ferr := os.Stat(s.indexDataPath); ferr == nil {
		s.internalIndex, err = bleve.Open(s.indexDataPath)
	} else {
		mapping := bleve.NewIndexMapping()
		s.internalIndex, err = bleve.New(s.indexDataPath, mapping)
	}

	if err != nil {
		return errors.Wrap(err, "cannot initialize index")
	}
	return nil
}

func (s *Index) Close() error {
	return s.internalIndex.Close()
}

func (s *Index) Build() (BuildOperationResult, error) {
	start := time.Now()

	files := 0
	var wg sync.WaitGroup
	err := s.git.ForEachCommitBundle(s.config.Repository.Path, s.config.Repository.MaxDepth, func(bundle *CommitBundle) error {
		if !s.state.ContainsCommit(bundle.Commit) {
			wg.Add(1)
			files += len(bundle.Files)

			go commitBundleIndexer(&wg, s.internalIndex, bundle)

			s.state.Append(bundle.Commit)
		}
		return nil
	})

	wg.Wait()

	err = s.state.Write()
	if err != nil {
		return BuildOperationResult{}, err
	}

	tookSec := time.Since(start).Seconds()
	response := BuildOperationResult{TookSeconds: tookSec, Files: files}
	return response, err
}

func (s *Index) Clean() (CleanOperationResult, error) {
	start := time.Now()
	var err error

	if _, ferr := os.Stat(s.indexDataPath); ferr == nil {
		err = os.RemoveAll(s.indexDataPath)
	}
	if err != nil {
		return CleanOperationResult{}, err
	}

	err = s.initialize()

	tookMs := time.Since(start).Milliseconds()
	response := CleanOperationResult{TookMillis: tookMs}
	return response, err
}

const (
	OutputHtml = "html"
	OutputAnsi = "ansi"
)

func (s *Index) Search(textQuery string, size int, output string) (SearchResult, error) {
	start := time.Now()

	query := bleve.NewQueryStringQuery(textQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = size
	searchRequest.Fields = []string{"*"} // return all fields in results
	searchRequest.Highlight = bleve.NewHighlightWithStyle(output)
	searchResult, err := s.internalIndex.Search(searchRequest)

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
	hash, err := s.git.GetLastCommitHash(s.config.Repository.Path)
	if err != nil {
		return false, err
	}

	return s.state.ContainsCommit(hash), nil
}

func (s *Index) DocumentCount() (uint64, error) {
	return s.internalIndex.DocCount()
}

func hitToIndexedFile(document *search.DocumentMatch) IndexedFile {
	return IndexedFile{
		Id:      document.Fields["Id"].(string),
		Hash:    document.Fields["Hash"].(string),
		Commit:  document.Fields["Commit"].(string),
		Content: document.Fields["Content"].(string),
		Path:    document.Fields["Path"].(string),
	}
}

// There is one id per unique file
func indexedFileId(fileName string, hash string) string {
	return fmt.Sprintf("%v:%v", fileName, hash)
}

func commitBundleIndexer(wg *sync.WaitGroup, index bleve.Index, commitBundle *CommitBundle) {
	defer wg.Done()

	batch := index.NewBatch()
	for _, file := range commitBundle.Files {
		doc := IndexedFile{
			Id:      indexedFileId(file.Name, file.Hash),
			Hash:    file.Hash,
			Commit:  commitBundle.Commit,
			Content: file.Content,
			Path:    file.Path,
		}
		berr := batch.Index(doc.Id, doc)
		if berr != nil {
			fmt.Println("Indexing error. Cannot batch file: ", berr)
		}
	}

	err := index.Batch(batch)
	if err != nil {
		fmt.Println(err, "Indexing error")
	}
}

type IndexPersistedState struct {
	Path           string
	IndexedCommits []string
}

func LoadPersistedState(path string) (IndexPersistedState, error) {
	jsonContent, err := ioutil.ReadFile(path)

	if err != nil { // Index may not exists
		state := IndexPersistedState{
			Path:           path,
			IndexedCommits: []string{},
		}
		return state, nil
	}

	var state IndexPersistedState
	err = json.Unmarshal(jsonContent, &state)

	return state, err
}

func (s *IndexPersistedState) Write() error {
	jsonContent, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "cannot marshall state as json")
	}

	return ioutil.WriteFile(s.Path, jsonContent, 0644)
}

func (s *IndexPersistedState) ContainsCommit(hash string) bool {
	for _, indexedHash := range s.IndexedCommits {
		if indexedHash == hash {
			return true
		}
	}
	return false
}

func (s *IndexPersistedState) Append(commit string) {
	s.IndexedCommits = append(s.IndexedCommits, commit)
}
