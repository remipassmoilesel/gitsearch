package index

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/test"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"strconv"
	"testing"
)

func Test_BuildIndex_emptyRepository(t *testing.T) {
	index := testIndex(t, test.REPO_EMPTY, 5)

	_, err := index.Build()
	assert.EqualError(t, err, "cannot build index: cannot get last commits, repository is empty: reference not found")
}

func Test_BuildIndex_batchSizeGreaterThanCommitSize(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)

	buildResult, err := index.Build()
	assert.NoError(t, err)
	assert.NotZero(t, buildResult.TookSeconds)
	assert.Equal(t, 4, buildResult.Files)
	assert.Equal(t, 10, buildResult.TotalFiles)

	count, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), count)
}

func Test_BuildIndex_batchSizeLessThanCommitSize(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 2)

	buildResult, err := index.Build()
	assert.NoError(t, err)
	assert.NotZero(t, buildResult.TookSeconds)
	assert.Equal(t, 4, buildResult.Files)
	assert.Equal(t, 10, buildResult.TotalFiles)

	count, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), count)
}

func Test_BuildIndex_buildShouldWriteState(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 2)

	_, err := index.Build()
	assert.NoError(t, err)

	statePath := index.state.Path()
	assert.FileExists(t, statePath)
}

func Test_BuildIndex_buildTwiceShouldNotIndex(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 2)

	_, err := index.Build()
	assert.NoError(t, err)

	err = index.Close()
	assert.NoError(t, err)
	err = index.initialize()
	assert.NoError(t, err)

	buildResult, err := index.Build()
	assert.NoError(t, err)

	assert.NotZero(t, buildResult.TookSeconds)
	assert.Equal(t, 0, buildResult.Files)
	assert.Equal(t, 0, buildResult.TotalFiles)

	count, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), count)
}

func Test_filterFiles(t *testing.T) {
	builder := IndexBuilderImpl{
		hashStore: &HashStoreImpl{
			store: []string{"a", "b"},
			utils: utils.NewUtils(),
		},
		utils: utils.NewUtils(),
	}

	files := []domain.IndexedFile{
		{
			Hash: "a",
		},
		{
			Hash: "b",
		},
		{
			Hash: "c",
		},
	}

	res := builder.filterFiles(files)
	assert.Equal(t, []domain.IndexedFile{files[2]}, res)
}

func Test_splitList(t *testing.T) {
	builder := IndexBuilderImpl{}

	files := []domain.IndexedFile{}
	for i := 0; i < 10; i++ {
		files = append(files, domain.IndexedFile{Hash: strconv.Itoa(i)})
	}

	res := builder.splitList(files, 20)
	assert.Len(t, res, 1)
	assert.Len(t, res[0], 10)
	assert.Equal(t, res[0][0].Hash, "0")
	assert.Equal(t, res[0][9].Hash, "9")
}

func Test_hashListFromFiles(t *testing.T) {
	files := []domain.IndexedFile{
		{
			Hash: "a",
		},
		{
			Hash: "b",
		},
		{
			Hash: "c",
		},
		{
			Hash: "c",
		},
	}

	res := hashListFromFiles(files)
	assert.Equal(t, []string{"a", "b", "c"}, res)
}

func Test_commitListFromFiles(t *testing.T) {
	files := []domain.IndexedFile{
		{
			Commit: "a",
		},
		{
			Commit: "b",
		},
		{
			Commit: "c",
		},
		{
			Commit: "c",
		},
	}

	res := commitListFromFiles(files)
	assert.Equal(t, []string{"a", "b", "c"}, res)
}

func Test_uniqueStrings(t *testing.T) {
	list := uniqueStrings([]string{"a", "b", "b", "c"})
	assert.Equal(t, []string{"a", "b", "c"}, list)
}
