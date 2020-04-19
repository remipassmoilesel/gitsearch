package test

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func Test_shellCommand(t *testing.T) {
	var err = Helper.ShellCommand("which sh")
	assert.NoError(t, err)
}

func Test_sampleGitRepo(t *testing.T) {
	repoPath, err := Helper.SampleGitRepository(REPO_SMALL)
	assert.NoError(t, err)
	assert.NotEmpty(t, repoPath)

	assert.DirExists(t, path.Join(repoPath, ".git"))
}
