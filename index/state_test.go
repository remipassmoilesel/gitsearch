package index

import (
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func Test_IndexState_Clean_ifNotExisting(t *testing.T) {
	dataPath, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	state := IndexState{path: dataPath}
	assert.NoError(t, state.Clean())
}

func Test_IndexState_Write(t *testing.T) {
	dir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dir, "gs-state.json")

	state := IndexState{path: statePath, IndexedCommits: []string{"a", "b", "c"}}
	assert.NoError(t, state.Write())
	assert.FileExists(t, statePath)

	loaded, err := LoadIndexState(statePath)
	assert.NoError(t, err)
	assert.Equal(t, statePath, loaded.path)
	assert.Equal(t, []string{"a", "b", "c"}, loaded.IndexedCommits)
}
