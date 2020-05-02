package index

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/pkg/errors"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/utils"
	"time"
)

type IndexBuilder struct {
	index              *Index
	config             config.Config
	git                utils.GitReader
	state              *IndexState
	hashStore          hashStore
	repositoryMaxDepth int
}

func NewIndexBuilder(index *Index) IndexBuilder {
	return IndexBuilder{
		index:     index,
		config:    index.config,
		git:       index.git,
		state:     &index.state,
		hashStore: hashStore{},
	}
}

type BuildOperationResult struct {
	TookSeconds float64
	Files       int
	TotalFiles  int
}

type batchIndexResult struct {
	commitList []string
	hashList   []string
	err        error
	took       float64
}

func (s *IndexBuilder) Build() (BuildOperationResult, error) {
	start := time.Now()

	commits, err := s.git.GetLastsCommits(s.config.Repository.MaxDepth)
	if err != nil {
		return BuildOperationResult{}, errors.Wrap(err, "cannot build index")
	}

	indexedFiles := 0
	totalFiles := 0
	batchSize := s.config.Index.BatchSize
	batchNumber := 0
	buffer := []IndexedFile{}

	ch := make(chan batchIndexResult)

	state := *(s.state)
	for _, commit := range commits {
		if state.ContainsCommit(commit) {
			continue
		}
		state.AppendCommit(commit)

		commitFiles, err := s.git.GetCommitFiles(commit)
		if err != nil {
			return BuildOperationResult{}, err
		}
		indexedFiles := s.commitFilesToIndexedFiles(commitFiles)
		totalFiles += len(indexedFiles)

		batches := s.splitList(s.filterFiles(indexedFiles), batchSize)
		s.hashStore.append(hashListFromFiles(indexedFiles))

		for _, batch := range batches {
			batchWithBuffer := append(buffer, batch...)
			if len(batchWithBuffer) < batchSize {
				buffer = append(buffer, batch...)
				continue
			}
			buffer = []IndexedFile{}

			shardId := batchNumber % s.index.shards.Size()
			batchNumber++

			go s.batchIndex(ch, *s.index.shards.GetShard(shardId), batchWithBuffer)
		}
	}

	if len(buffer) > 0 {
		shardId := batchNumber % s.index.shards.Size()
		batchNumber++

		go s.batchIndex(ch, *s.index.shards.GetShard(shardId), buffer)
	}

	if batchNumber > 0 {
		i := 0
		for res := range ch {
			if res.err != nil {
				return BuildOperationResult{}, err
			}

			indexedFiles += len(res.hashList)

			i++
			if i == batchNumber {
				break
			}
		}
	}

	err = state.Write()
	if err != nil {
		return BuildOperationResult{}, err
	}

	tookSec := time.Since(start).Seconds()
	response := BuildOperationResult{TookSeconds: tookSec, Files: indexedFiles, TotalFiles: totalFiles}
	return response, err
}

func (s *IndexBuilder) batchIndex(ch chan batchIndexResult, index bleve.Index, files []IndexedFile) {
	batch := index.NewBatch()
	for _, file := range files {
		berr := batch.Index(file.Hash, file)
		if berr != nil {
			fmt.Println("Indexing error. Cannot batch file: ", berr)
		}
	}

	start := time.Now()
	err := index.Batch(batch)
	took := time.Since(start)

	res := batchIndexResult{
		commitList: commitListFromFiles(files),
		hashList:   hashListFromFiles(files),
		took:       took.Seconds(),
		err:        err,
	}

	ch <- res
}

func (s *IndexBuilder) filterFiles(files []IndexedFile) []IndexedFile {
	hashList := hashListFromFiles(files)
	filteredHashList := s.hashStore.filterExisting(hashList)

	filteredFiles := []IndexedFile{}
	for _, file := range files {
		if utils.ContainsString(filteredHashList, file.Hash) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func (s *IndexBuilder) splitList(files []IndexedFile, bundleSize int) [][]IndexedFile {
	processed := 0
	result := [][]IndexedFile{}
	for processed < len(files) {
		upperBound := processed + bundleSize
		if upperBound > len(files) {
			upperBound = len(files)
		}
		files := files[processed:upperBound]
		result = append(result, files)
		processed += bundleSize
	}
	return result
}

func (s *IndexBuilder) commitFilesToIndexedFiles(files []utils.CommitFile) []IndexedFile {
	res := []IndexedFile{}
	for _, fl := range files {
		res = append(res, IndexedFile{
			Hash:    fl.Hash,
			Commit:  fl.Commit,
			Date:    fl.Date,
			Content: fl.Content,
			Path:    fl.Path,
			Name:    fl.Name,
		})
	}
	return res
}

func hashListFromFiles(files []IndexedFile) []string {
	res := []string{}
	for _, file := range files {
		res = append(res, file.Hash)
	}
	return uniqueStrings(res)
}

func commitListFromFiles(files []IndexedFile) []string {
	res := []string{}
	for _, file := range files {
		res = append(res, file.Commit)
	}
	return uniqueStrings(res)
}

func uniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
