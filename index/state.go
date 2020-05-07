package index

import (
	"encoding/json"
	"github.com/nightlyone/lockfile"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"io/ioutil"
	"os"
	"path"
)

type IndexState interface {
	Path() string
	TryLock() error
	Unlock() error
	AppendCommit(commit string)
	ContainsCommit(commit string) bool
	Content() *PersistedState
	Write() error
	Clean() error
}

type IndexStateImpl struct {
	path     string
	lockPath string
	lock     lockfile.Lockfile
	state    *PersistedState
}

type PersistedState struct {
	IndexedCommits []string
}

const (
	StateFileName = "gs-index-state.json"
	StateLockName = "gs-state.lock"
)

func LoadIndexState(stateDir string) (IndexState, error) {
	statePath := path.Join(stateDir, StateFileName)
	lockPath := path.Join(stateDir, StateLockName)
	jsonContent, err := ioutil.ReadFile(statePath)

	if err != nil { // Index may not exists
		state := IndexStateImpl{
			path:     statePath,
			lockPath: lockPath,
			state: &PersistedState{
				IndexedCommits: []string{},
			},
		}
		return &state, nil
	}

	var pState PersistedState
	err = json.Unmarshal(jsonContent, &pState)
	if err != nil {
		return &IndexStateImpl{}, errors.Wrap(err, "cannot unmarshall state file")
	}

	state := IndexStateImpl{
		path:     statePath,
		lockPath: lockPath,
		lock:     "",
		state:    &pState,
	}
	return &state, nil
}

func (s *IndexStateImpl) Path() string {
	return s.path
}

func (s *IndexStateImpl) Content() *PersistedState {
	return s.state
}

func (s *IndexStateImpl) TryLock() error {
	err := os.MkdirAll(path.Dir(s.lockPath), 0755)
	if err != nil {
		return errors.Wrap(err, "cannot lock index state")
	}

	lock, err := lockfile.New(s.lockPath)
	if err != nil {
		return errors.Wrap(err, "cannot lock index state")
	}

	s.lock = lock
	return lock.TryLock()
}

func (s *IndexStateImpl) Unlock() error {
	return s.lock.Unlock()
}

func (s *IndexStateImpl) AppendCommit(commit string) {
	s.state.IndexedCommits = append(s.state.IndexedCommits, commit)
}

func (s *IndexStateImpl) ContainsCommit(hash string) bool {
	return utils.ContainsString(s.state.IndexedCommits, hash)
}

func (s *IndexStateImpl) Write() error {
	err := os.MkdirAll(path.Dir(s.path), 0755)
	if err != nil {
		return errors.Wrap(err, "cannot write state")
	}

	jsonContent, err := json.Marshal(s.state)
	if err != nil {
		return errors.Wrap(err, "cannot marshall state")
	}

	return ioutil.WriteFile(s.path, jsonContent, 0644)
}

func (s *IndexStateImpl) Clean() error {
	if _, ferr := os.Stat(s.path); ferr == nil {
		err := os.Remove(s.path)
		return errors.Wrap(err, "cannot clean state file")
	}
	return nil
}
