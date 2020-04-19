package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_hashStore_filterExisting(t *testing.T) {
	store := hashStore{
		store: []string{"a", "b"},
	}

	res := store.filterExisting([]string{"a", "b", "c"})
	assert.Equal(t, []string{"c"}, res)
}
