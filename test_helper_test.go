package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var testHelper = TestHelper{}

type TestHelper struct {
}

func (s *TestHelper) RandomDataPath() string {
	rUuid, _ := uuid.NewRandom()
	return path.Join(os.TempDir(), "gitsearch", "test-data-"+rUuid.String())
}

func (s *TestHelper) SampleGitRepository() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	testDataPath := path.Join(cwd, "test-data")
	repoPath := s.RandomDataPath()

	err = s.ShellCommand("mkdir -p " + repoPath)
	if err != nil {
		return "", err
	}

	err = s.ShellCommand("cp -R " + testDataPath + "/* " + repoPath)
	if err != nil {
		return "", err
	}

	err = s.ShellCommandOnDir("git init", repoPath)
	if err != nil {
		return "", err
	}

	err = s.ShellCommandOnDir("git config user.email gitsearch@test && git config user.name gitsearch@test", repoPath)
	if err != nil {
		return "", err
	}

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if strings.Index(path, ".git") != -1 || path == repoPath {
			return nil
		}
		return s.commitFile(path, repoPath)
	})

	return repoPath, err
}

func (s *TestHelper) commitFile(file string, repoPath string) error {
	return s.ShellCommandOnDir("git add "+file+" && git commit -m 'Commit "+file+"'", repoPath)
}

func (s *TestHelper) ShellCommand(cmd string) error {
	return s.ShellCommandOnDir(cmd, "")
}

func (s *TestHelper) ShellCommandOnDir(cmd string, workingDir string) error {
	sh, err := exec.LookPath("sh")
	if err != nil {
		return err
	}

	shCmd := &exec.Cmd{
		Path: sh,
		Args: append([]string{sh, "-c", cmd}),
		Dir:  workingDir,
		// Stdout: os.Stdout, // For debug purposes
		// Stderr: os.Stderr, // For debug purposes
	}

	// fmt.Println("Executing: " + shCmd.String()) 	// For debug purposes
	return shCmd.Run()
}

func Test_shellCommand(t *testing.T) {
	var err = testHelper.ShellCommand("which sh")
	assert.NoError(t, err)
}

func Test_sampleGitRepo(t *testing.T) {
	repoPath, err := testHelper.SampleGitRepository()
	assert.NoError(t, err)
	assert.NotEmpty(t, repoPath)
}
