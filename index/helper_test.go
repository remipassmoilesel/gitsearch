package index

import (
	"fmt"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Create a fake git repository then initialize index in
func testIndex(t *testing.T, templateName string, batchSize int) Index {
	var path, err = test.Helper.SampleGitRepository(templateName)
	assert.NoError(t, err)

	dir, err := test.Helper.RandomDataDir()
	fmt.Println("Using data directory: " + dir)
	assert.NoError(t, err)

	cfg := config.Config{
		DataRootPath: dir,
		Repository: config.RepositoryContext{
			Path:     path,
			MaxDepth: 5,
		},
		Search: config.SearchConfig{
			Shards:    3,
			BatchSize: batchSize,
		},
	}
	index, err := NewIndex(cfg)
	assert.NoError(t, err)

	return index
}
