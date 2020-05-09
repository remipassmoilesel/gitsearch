package index

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/test"
	"os"
	"path"
	"testing"
)

func Test_IndexState_Write(t *testing.T) {
	dataDir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dataDir, StateFileName)

	state := IndexStateImpl{
		path: statePath,
		state: &domain.PersistedState{
			IndexedCommits: []string{"a", "b", "c"},
		},
	}
	assert.NoError(t, state.Write())
	assert.FileExists(t, statePath)

	loaded, err := LoadIndexState(dataDir)
	assert.NoError(t, err)
	assert.Equal(t, statePath, loaded.Path())
	assert.Equal(t, state.state.IndexedCommits, loaded.Content().IndexedCommits)
}

func Test_IndexState_LoadIndexState_ifNotExisting(t *testing.T) {
	dataDir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dataDir, StateFileName)

	state, err := LoadIndexState(dataDir)
	assert.NoError(t, err)

	assert.Equal(t, statePath, state.Path())

	_, err = os.Stat(statePath)
	assert.True(t, os.IsNotExist(err))
}

func Test_IndexState_LoadIndexState_ifExisting(t *testing.T) {
	dataDir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dataDir, StateFileName)

	state, err := LoadIndexState(dataDir)
	assert.NoError(t, err)

	_, err = os.Stat(statePath)
	assert.True(t, os.IsNotExist(err))

	state.AppendCommit("a")
	state.AppendCommit("b")
	state.AppendCommit("c")

	assert.NoError(t, state.Write())

	state2, err := LoadIndexState(dataDir)
	assert.NoError(t, err)
	assert.Equal(t, state2.Path(), state.Path())
	assert.Equal(t, state2.Content(), state.Content())
}

func Test_IndexState_Clean_ifNotExisting(t *testing.T) {
	dataDir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dataDir, StateFileName)

	state := IndexStateImpl{path: statePath}
	assert.NoError(t, state.Clean())
}

func Test_IndexState_Clean_ifExisting(t *testing.T) {
	dataDir, err := test.Helper.RandomDataDir()
	assert.NoError(t, err)

	statePath := path.Join(dataDir, StateFileName)

	state, err := LoadIndexState(dataDir)
	assert.NoError(t, err)

	state.AppendCommit("a")
	state.AppendCommit("b")
	state.AppendCommit("c")
	err = state.Write()
	assert.NoError(t, err)

	assert.NoError(t, state.Clean())
	_, err = os.Stat(statePath)
	assert.True(t, os.IsNotExist(err))
}
