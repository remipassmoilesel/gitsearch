package index

import (
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_ShardGroup_Initialize(t *testing.T) {
	path, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	group := ShardGroup{
		dataRootPath: path,
		shardNumber:  8,
	}

	res, err := group.Initialize()
	assert.NoError(t, err)
	assert.Equal(t, 8, res.number)
}

func Test_ShardGroup_Clean(t *testing.T) {
	dataPath, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	group := ShardGroup{
		dataRootPath: dataPath,
		shardNumber:  8,
	}
	assert.Len(t, group.shards, 0)

	err = group.Clean()
	assert.NoError(t, err)
	assert.Len(t, group.shards, 0)

	_, err = os.Stat(dataPath)
	assert.True(t, os.IsNotExist(err))
}

func Test_ShardGroup_InitializeThenClean(t *testing.T) {
	dataPath, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	group := ShardGroup{
		dataRootPath: dataPath,
		shardNumber:  8,
	}

	res, err := group.Initialize()
	assert.NoError(t, err)
	assert.Equal(t, 8, res.number)
	assert.Len(t, group.shards, 8)

	err = group.Clean()
	assert.NoError(t, err)
	assert.Len(t, group.shards, 0)

	_, err = os.Stat(dataPath)
	assert.True(t, os.IsNotExist(err))
}
