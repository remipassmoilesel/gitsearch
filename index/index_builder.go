package index

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/git_reader"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"time"
)

type IndexBuilderImpl struct {
	index              *IndexImpl
	config             config.Config
	git                git_reader.GitReader
	state              IndexState
	hashStore          HashStore
	utils              utils.Utils
	repositoryMaxDepth int
}

func NewIndexBuilder(cfg config.Config, state IndexState, index *IndexImpl) (IndexBuilderImpl, error) {
	gitReader, err := git_reader.NewGitReader(cfg.Repository.Path)
	if err != nil {
		return IndexBuilderImpl{}, err
	}

	return IndexBuilderImpl{
		index:     index,
		config:    cfg,
		git:       gitReader,
		state:     state,
		utils:     utils.NewUtils(),
		hashStore: NewHashStore(),
	}, nil
}

type batchIndexResult struct {
	commitList []string
	hashList   []string
	err        error
	took       float64
}

const (
	BuildModeLastCommits     = "BuildModeLastCommits"
	BuildModeCommitsSpacedBy = "BuildModeCommitsSpacedBy"
)

func DefaultBuildOptions() domain.BuildOptions {
	return domain.BuildOptions{
		Mode: BuildModeLastCommits,
	}
}

func BuildOptionsSpacedBy() domain.BuildOptions {
	return domain.BuildOptions{
		Mode:            BuildModeCommitsSpacedBy,
		SpacedBySeconds: 24 * 3600,
	}
}

func (s *IndexBuilderImpl) Build(options domain.BuildOptions) (domain.BuildOperationResult, error) {
	start := time.Now()

	var commits []git_reader.Commit
	var err error
	if options.Mode == BuildModeLastCommits {
		commits, err = s.git.GetLastsCommits(s.config.Repository.MaxDepth)
	} else if options.Mode == BuildModeCommitsSpacedBy {
		commits, err = s.git.GetCommitsSpacedBy(s.config.Repository.MaxDepth, options.SpacedBySeconds)
	}
	if err != nil {
		return domain.BuildOperationResult{}, errors.Wrap(err, "cannot build index")
	}

	indexedFiles := 0
	totalFiles := 0
	batchSize := s.config.Index.BatchSize
	batchNumber := 0
	buffer := []domain.IndexedFile{}

	ch := make(chan batchIndexResult)

	lastCommit := git_reader.Commit{}
	for _, commit := range commits {
		if s.state.ContainsCommit(commit.Hash) {
			continue
		}
		s.state.AppendCommit(commit.Hash)
		lastCommit = commit

		commitFiles, err := s.git.GetCommitFiles(commit.Hash)
		if err != nil {
			return domain.BuildOperationResult{}, err
		}
		indexedFiles := s.commitFilesToIndexedFiles(commitFiles)
		totalFiles += len(indexedFiles)

		batches := s.splitList(s.filterFiles(indexedFiles), batchSize)
		s.hashStore.Append(hashListFromFiles(indexedFiles))

		for _, batch := range batches {
			batchWithBuffer := append(buffer, batch...)
			if len(batchWithBuffer) < batchSize {
				buffer = append(buffer, batch...)
				continue
			}
			buffer = []domain.IndexedFile{}

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
				return domain.BuildOperationResult{}, err
			}

			indexedFiles += len(res.hashList)

			i++
			if i == batchNumber {
				break
			}
		}
	}

	err = s.state.Write()
	if err != nil {
		return domain.BuildOperationResult{}, err
	}

	tookSec := time.Since(start).Seconds()
	response := domain.BuildOperationResult{TookSeconds: tookSec, Files: indexedFiles, TotalFiles: totalFiles, OldestCommit: lastCommit.Date}
	return response, err
}

func (s *IndexBuilderImpl) batchIndex(ch chan batchIndexResult, index bleve.Index, files []domain.IndexedFile) {
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

func (s *IndexBuilderImpl) filterFiles(files []domain.IndexedFile) []domain.IndexedFile {
	hashList := hashListFromFiles(files)
	filteredHashList := s.hashStore.FilterExisting(hashList)

	filteredFiles := []domain.IndexedFile{}
	for _, file := range files {
		if s.utils.ContainsString(filteredHashList, file.Hash) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func (s *IndexBuilderImpl) splitList(files []domain.IndexedFile, bundleSize int) [][]domain.IndexedFile {
	processed := 0
	result := [][]domain.IndexedFile{}
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

func (s *IndexBuilderImpl) commitFilesToIndexedFiles(files []git_reader.CommitFile) []domain.IndexedFile {
	res := []domain.IndexedFile{}
	for _, fl := range files {
		res = append(res, domain.IndexedFile{
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

func hashListFromFiles(files []domain.IndexedFile) []string {
	res := []string{}
	for _, file := range files {
		res = append(res, file.Hash)
	}
	return uniqueStrings(res)
}

func commitListFromFiles(files []domain.IndexedFile) []string {
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
