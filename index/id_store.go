package index

import "github.com/remipassmoilesel/gitsearch/utils"

// TODO: use persisted store
type hashStore struct {
	store []string
}

func (s *hashStore) append(ids []string) {
	s.store = append(s.store, ids...)
}

func (s *hashStore) filterExisting(hashList []string) []string {
	res := []string{}
	for _, hash := range hashList {
		if !utils.ContainsString(s.store, hash) {
			res = append(res, hash)
		}
	}
	return res
}
