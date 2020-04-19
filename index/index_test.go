package index

import (
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Index_Search(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	res, err := index.Search("lorem", 5, OutputAnsi)

	assert.NoError(t, err)
	assert.NotZero(t, res.TookUs)
	assert.NotZero(t, len(res.Matches))
	assert.Equal(t, res.Query, "lorem")
	assert.NotEmpty(t, res.Matches[0].File.Commit)
	assert.NotEmpty(t, res.Matches[0].File.Content)
	assert.NotEmpty(t, res.Matches[0].File.Hash)
	assert.NotEmpty(t, res.Matches[0].File.Path)
}

func Test_Index_Clean(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)

	_, err := index.Clean()
	assert.NoError(t, err)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)

	_, err = os.Stat(index.state.path)
	assert.True(t, os.IsNotExist(err))
}

func Test_Index_Clean_afterBuild(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)

	_, err := index.Build()
	assert.NoError(t, err)

	_, err = index.Clean()
	assert.NoError(t, err)

	err = index.initialize()
	assert.NoError(t, err)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)

	_, err = os.Stat(index.state.path)
	assert.True(t, os.IsNotExist(err))
}
