package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path"
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

type SearchMatch struct {
	File    IndexedFile
	Matches map[string][]string
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

func (s *Index) Build() error {
	err := s.git.ForEachFiles(s.config.Repository.Path, func(commit *object.Commit, file *object.File) error {
		content, err := file.Contents()
		if err != nil {
			content = fmt.Sprintf("Error %s", err)
		}
		doc := IndexedFile{
			Id:      indexedFileId(file),
			Hash:    file.Hash.String(),
			Commit:  commit.Hash.String(),
			Content: content,
			Path:    file.Name,
		}
		err = s.internalIndex.Index(doc.Id, doc)
		if err != nil {
			fmt.Println("indexing error: ", err)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "indexing error")
	}

	return err
}

func (s *Index) Search(query string) ([]SearchMatch, error) {
	bQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(bQuery)
	searchRequest.Highlight = bleve.NewHighlight() // can be html highlight
	searchResult, err := s.internalIndex.Search(searchRequest)

	if err != nil {
		return []SearchMatch{}, errors.Wrap(err, "search error")
	}

	matches := make([]SearchMatch, searchResult.Hits.Len())
	for index, hit := range searchResult.Hits {
		doc, err := s.internalIndex.Document(hit.ID)
		if err != nil {
			fmt.Println("error while fetching document " + hit.ID)
			continue
		}
		indexedFile := bleveDocumentToIndexedFile(doc)
		match := SearchMatch{
			File:    indexedFile,
			Matches: hit.Fragments,
		}
		matches[index] = match
	}

	return matches, nil
}

func (s *Index) Clean() error {
	if _, ferr := os.Stat(s.indexDataPath); ferr == nil {
		return os.RemoveAll(s.indexDataPath)
	}
	return nil
}

func bleveDocumentToIndexedFile(document *document.Document) IndexedFile {
	return IndexedFile{
		Id:      string(document.Fields[0].Value()),
		Hash:    string(document.Fields[1].Value()),
		Commit:  string(document.Fields[2].Value()),
		Content: string(document.Fields[3].Value()),
		Path:    string(document.Fields[4].Value()),
	}
}

// There is one id per unique file
func indexedFileId(file *object.File) string {
	return fmt.Sprintf("%s:%s", file.Name, file.Hash.String())
}
