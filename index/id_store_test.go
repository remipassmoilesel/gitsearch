package index

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"testing"
)

func Test_hashStore_filterExisting(t *testing.T) {
	store := HashStoreImpl{
		store: []string{"a", "b"},
		utils: utils.NewUtils(),
	}

	res := store.FilterExisting([]string{"a", "b", "c"})
	assert.Equal(t, []string{"c"}, res)
}
