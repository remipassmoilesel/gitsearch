package utils

import (
	"github.com/remipassmoilesel/gitsearch/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitReader_GitReader_GetHeadHash(t *testing.T) {
	git, err := gitReaderOnSmallRepo()
	assert.NoError(t, err)

	hash, err := git.GetHeadHash()
	assert.Equal(t, "8158970149e4199242f0a428e4b7d32b10264dfc", hash)
}

func TestGitReader_GitReader_GetLastsCommits(t *testing.T) {
	git, err := gitReaderOnSmallRepo()
	assert.NoError(t, err)

	hashes, err := git.GetLastsCommits(2)
	assert.Equal(t, []string{"8158970149e4199242f0a428e4b7d32b10264dfc", "dca73bebbcbf5d3e13b11b87939fb203013dce58"}, hashes)
}

func TestGitReader_GitReader_GetCommitFiles(t *testing.T) {
	git, err := gitReaderOnSmallRepo()
	assert.NoError(t, err)

	head, err := git.GetHeadHash()
	assert.NoError(t, err)

	files, err := git.GetCommitFiles(head)
	assert.Len(t, files, 4)

	for _, file := range files {
		assert.NotEmpty(t, file.Path)
		assert.NotEmpty(t, file.Name)
		assert.NotEmpty(t, file.Content)
		assert.NotEmpty(t, file.Hash)
		assert.Equal(t, head, file.Commit)
	}
}

func gitReaderOnSmallRepo() (GitReader, error) {
	repoPath, err := test.Helper.SampleGitRepository(test.REPO_SMALL)
	if err != nil {
		return GitReader{}, err
	}

	git, err := NewGitReader(repoPath)
	return git, err
}
