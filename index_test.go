package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BuildIndex(t *testing.T) {
	testIndex(t)
}

func Test_Search(t *testing.T) {
	index := testIndex(t)
	res, err := index.Search("lorem", 5, OutputAnsi)

	assert.NoError(t, err)
	assert.NotZero(t, res.TookUs)
	assert.NotZero(t, len(res.Matches))
	assert.Equal(t, res.Query, "lorem")
	assert.NotEmpty(t, res.Matches[0].File.Id)
	assert.NotEmpty(t, res.Matches[0].File.Commit)
	assert.NotEmpty(t, res.Matches[0].File.Content)
	assert.NotEmpty(t, res.Matches[0].File.Hash)
	assert.NotEmpty(t, res.Matches[0].File.Path)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Equal(t, docCount, uint64(3))
}

func Test_Clean(t *testing.T) {
	index := testIndex(t)
	res, err := index.Clean()

	assert.NoError(t, err)
	assert.NotZero(t, res.TookMillis)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)
}

// Create a fake git repository, then index it and return index
func testIndex(t *testing.T) Index {
	var path, err = testHelper.SampleGitRepository()
	assert.NoError(t, err)

	config := Config{
		DataRootPath: testHelper.RandomDataPath(),
		Repository: RepositoryContext{
			Path:     path,
			MaxDepth: 5,
		},
	}
	index, err := NewIndex(config)
	assert.NoError(t, err)

	res, err := index.Build()
	assert.NoError(t, err)
	assert.NotZero(t, res.TookSeconds)
	assert.NotZero(t, res.Files)

	return index
}
