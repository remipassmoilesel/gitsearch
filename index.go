package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path"
	"time"
)

type Index struct {
	config        Config
	git           GitHelper
	indexDataPath string
	internalIndex bleve.Index
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
	Query       string
	TookSeconds float64
	Matches     []SearchMatch
}

type SearchMatch struct {
	File      IndexedFile
	Fragments []string
}

type IndexOperationResult struct {
	TookSeconds float64
}

func NewIndex(config Config) (Index, error) {
	indexDataPath := path.Join(config.DataRootPath, "index", config.Repository.Path)

	var internalIndex bleve.Index
	var err error
	if _, ferr := os.Stat(indexDataPath); ferr == nil {
		internalIndex, err = bleve.Open(indexDataPath)
	} else {
		mapping := bleve.NewIndexMapping()
		internalIndex, err = bleve.New(indexDataPath, mapping)
	}

	if err != nil {
		return *new(Index), errors.Wrap(err, "cannot initialize index")
	}

	index := Index{
		config:        config,
		git:           GitHelper{},
		indexDataPath: indexDataPath,
		internalIndex: internalIndex,
	}
	return index, nil
}

func (s *Index) Close() error {
	return s.internalIndex.Close()
}

func (s *Index) Build() (IndexOperationResult, error) {
	start := time.Now()
	fileHandler := func(commit *object.Commit, file *object.File) error {
		content, ierr := file.Contents()
		if ierr == nil {
			doc := IndexedFile{
				Id:      indexedFileId(file),
				Hash:    file.Hash.String(),
				Commit:  commit.Hash.String(),
				Content: content,
				Path:    file.Name,
			}
			ierr = s.internalIndex.Index(doc.Id, doc)
		}
		if ierr != nil {
			fmt.Println("indexing error: ", ierr)
		}
		return nil
	}
	err := s.git.ForEachFiles(s.config.Repository.Path, s.config.Repository.MaxDepth, fileHandler)

	tookSec := time.Since(start).Seconds()
	response := IndexOperationResult{TookSeconds: tookSec}
	return response, err
}

func (s *Index) Clean() (IndexOperationResult, error) {
	start := time.Now()
	var err error

	if _, ferr := os.Stat(s.indexDataPath); ferr == nil {
		err = os.RemoveAll(s.indexDataPath)
	}

	tookSec := time.Since(start).Seconds()
	response := IndexOperationResult{TookSeconds: tookSec}
	return response, err
}

func (s *Index) Search(textQuery string) (SearchResult, error) {
	start := time.Now()

	query := bleve.NewQueryStringQuery(textQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 100
	searchRequest.Fields = []string{"*"} // return all fields in results
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
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

	tookSec := time.Since(start).Seconds()
	response := SearchResult{Query: textQuery, TookSeconds: tookSec, Matches: resultMatches}
	return response, err
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
func indexedFileId(file *object.File) string {
	return fmt.Sprintf("%s:%s", file.Name, file.Hash.String())
}
