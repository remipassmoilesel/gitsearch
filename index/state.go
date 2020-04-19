package index

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/remipassmoilesel/gitsearch/utils"
	"io/ioutil"
	"os"
)

type IndexState struct {
	path           string
	IndexedCommits []string
}

func LoadIndexState(path string) (IndexState, error) {
	jsonContent, err := ioutil.ReadFile(path)

	if err != nil { // Index may not exists
		state := IndexState{
			path:           path,
			IndexedCommits: []string{},
		}
		return state, nil
	}

	var state IndexState
	err = json.Unmarshal(jsonContent, &state)
	state.path = path

	return state, errors.Wrap(err, "cannot unmarshall state file")
}

func (s *IndexState) Append(commit string) {
	s.IndexedCommits = append(s.IndexedCommits, commit)
}

func (s *IndexState) Write() error {
	jsonContent, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "cannot marshall state")
	}

	return ioutil.WriteFile(s.path, jsonContent, 0644)
}

func (s *IndexState) ContainsCommit(hash string) bool {
	return utils.ContainsString(s.IndexedCommits, hash)
}

func (s *IndexState) Clean() error {
	if _, ferr := os.Stat(s.path); ferr == nil {
		err := os.Remove(s.path)
		return errors.Wrap(err, "cannot clean state file")
	}
	return nil
}
