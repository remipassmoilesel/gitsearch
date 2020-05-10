package index

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/test"
	"os"
	"testing"
)

func Test_Index_Initialize_unlocked(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)
	assert.NoError(t, index.Close())

	docs, err := index.DocumentCount()
	assert.NoError(t, err)

	_, err = NewIndex(index.config)
	assert.NoError(t, err)

	docsB, err := index.DocumentCount()
	assert.NoError(t, err)

	assert.Equal(t, docs, docsB)
}

func Test_Index_Search(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	res, err := index.Search("lorem", 5, OutputAnsi)

	assert.NoError(t, err)
	// assert.NotZero(t, res.TookMs)
	assert.NotZero(t, len(res.Matches))
	assert.Equal(t, res.Query, "lorem")
	assert.NotEmpty(t, res.Matches[0].File.Commit)
	assert.NotEmpty(t, res.Matches[0].File.Content)
	assert.NotEmpty(t, res.Matches[0].File.Hash)
	assert.NotEmpty(t, res.Matches[0].File.Path)
	assert.NotEmpty(t, res.Matches[0].File.Name)
	assert.NotZero(t, res.Matches[0].File.Date.Unix())
}

func Test_Index_FindDocumentById_fullHash(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	res, err := index.FindDocumentById("2e27e035b1b327134c62614fd764be8b352f51e5")

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Commit)
	assert.NotEmpty(t, res.Content)
	assert.NotEmpty(t, res.Hash)
	assert.NotEmpty(t, res.Name)
	assert.NotEmpty(t, res.Path)
}

func Test_Index_FindDocumentById_hashPrefix(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	res, err := index.FindDocumentById("2e27e035b1b32")

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Commit)
	assert.NotEmpty(t, res.Content)
	assert.NotEmpty(t, res.Hash)
	assert.NotEmpty(t, res.Name)
	assert.NotEmpty(t, res.Path)
}

func Test_Index_FindDocumentById_wrongHash(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	_, err = index.FindDocumentById("not-a-valid-hash")
	assert.Error(t, err)
}

func Test_Index_Clean(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)

	_, err := index.Clean()
	assert.NoError(t, err)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)

	_, err = os.Stat(index.state.Path())
	assert.True(t, os.IsNotExist(err))
}

func Test_Index_Clean_afterBuild(t *testing.T) {
	idx := testIndex(t, test.REPO_SMALL, 5)

	_, err := idx.Build()
	assert.NoError(t, err)

	_, err = idx.Clean()
	assert.NoError(t, err)

	_, err = os.Stat(idx.state.Path())
	assert.True(t, os.IsNotExist(err))

	err = idx.Close()
	assert.NoError(t, err)

	idx2, err := NewIndex(idx.config)
	assert.NoError(t, err)

	docCount, err := idx2.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)

	_, err = os.Stat(idx.state.Path())
	assert.NoError(t, nil)
}

// Create a fake git repository then initialize index in
func testIndex(t *testing.T, templateName string, batchSize int) IndexImpl {
	var path, err = test.Helper.SampleGitRepository(templateName)
	assert.NoError(t, err)

	dir, err := test.Helper.RandomDataDir()
	fmt.Println("Using data directory: " + dir)
	assert.NoError(t, err)

	cfg := config.Config{
		DataRootPath: dir,
		Repository: config.RepositoryContext{
			Path:     path,
			MaxDepth: 5,
		},
		Index: config.SearchConfig{
			Shards:    3,
			BatchSize: batchSize,
		},
	}
	index, err := NewIndex(cfg)
	assert.NoError(t, err)

	return *(index.(*IndexImpl))
}
