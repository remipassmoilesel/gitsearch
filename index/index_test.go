package index

import (
	"errors"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Index_Initialize_locked(t *testing.T) {
	lockedState := NewFakeState(true)
	index := Index{
		state: &lockedState,
	}

	err := index.initialize()
	assert.EqualError(t, err, "cannot initialize index: index locked")
}

func Test_Index_Initialize_unlocked(t *testing.T) {
	index := testIndex(t, test.REPO_SMALL, 5)
	_, err := index.Build()
	assert.NoError(t, err)

	assert.NoError(t, index.Close())

	cfg := config.Config{
		DataRootPath: index.config.DataRootPath,
		Repository: config.RepositoryContext{
			Path:     index.config.Repository.Path,
			MaxDepth: 5,
		},
		Index: config.SearchConfig{
			Shards:    3,
			BatchSize: index.config.Index.BatchSize,
		},
	}
	_, err = NewIndex(cfg)
	assert.NoError(t, err)
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
	index := testIndex(t, test.REPO_SMALL, 5)

	_, err := index.Build()
	assert.NoError(t, err)

	_, err = index.Clean()
	assert.NoError(t, err)

	_, err = os.Stat(index.state.Path())
	assert.True(t, os.IsNotExist(err))

	err = index.initialize()
	assert.NoError(t, err)

	docCount, err := index.DocumentCount()
	assert.NoError(t, err)
	assert.Zero(t, docCount)

	_, err = os.Stat(index.state.Path())
	assert.NoError(t, nil)
}

type FakeState struct {
	locked bool
}

func NewFakeState(locked bool) FakeState {
	return FakeState{locked}
}

func (s *FakeState) Path() string {
	return ""
}

func (s *FakeState) Content() *PersistedState {
	return &PersistedState{}
}

func (s *FakeState) TryLock() error {
	if s.locked {
		return errors.New("index locked")
	}
	return nil
}

func (s *FakeState) Unlock() error {
	s.locked = false
	return nil
}

func (s *FakeState) AppendCommit(commit string) {

}

func (s *FakeState) ContainsCommit(hash string) bool {
	return false
}

func (s *FakeState) Write() error {
	return nil
}

func (s *FakeState) Clean() error {
	return nil
}
