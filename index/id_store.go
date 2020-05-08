package index

import "gitlab.com/remipassmoilesel/gitsearch/utils"

type HashStore interface {
	Append(ids []string)
	FilterExisting(hashList []string) []string
}

func NewHashStore() HashStore {
	return &HashStoreImpl{
		utils: utils.NewUtils(),
	}
}

// TODO: use persisted storage
type HashStoreImpl struct {
	store []string
	utils utils.Utils
}

func (s *HashStoreImpl) Append(ids []string) {
	s.store = append(s.store, ids...)
}

func (s *HashStoreImpl) FilterExisting(hashList []string) []string {
	res := []string{}
	for _, hash := range hashList {
		if !s.utils.ContainsString(s.store, hash) {
			res = append(res, hash)
		}
	}
	return res
}
