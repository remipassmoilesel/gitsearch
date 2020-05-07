package test

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path"
)

type TestHelper struct {
}

var Helper = TestHelper{}

func (s *TestHelper) ProjectRoot() (string, error) {
	return os.Getwd()
}

func (s *TestHelper) TestDataRoot() (string, error) {
	root, err := s.ProjectRoot()
	if err != nil {
		return "", err
	}
	return path.Join(root, "test", "data"), nil
}

func (s *TestHelper) RandomDataDir() (string, error) {
	rUuid, _ := uuid.NewRandom()
	dirPath := path.Join(os.TempDir(), "gitsearch", "test-"+rUuid.String())
	err := s.ShellCommand(fmt.Sprintf("mkdir -p %s", dirPath))
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

const (
	REPO_SMALL          = "small-repo.tar"
	REPO_EMPTY          = "empty-repo.tar"
	REPO_SPACED_COMMITS = "spaced-commits.tar"
)

func allTemplates() []string {
	return []string{REPO_EMPTY, REPO_SMALL, REPO_SPACED_COMMITS}
}

func (s *TestHelper) SampleGitRepository(templateName string) (string, error) {
	if !s.isTemplateKnown(templateName) {
		return "", errors.New("Unknwon template: " + templateName)
	}

	testData, err := s.TestDataRoot()
	if err != nil {
		return "", err
	}

	templatePath := path.Join(testData, templateName)
	repoPath, err := s.RandomDataDir()
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("Using repository: %s", repoPath))

	// Here we extract template and we remove template directory
	err = s.ShellCommandOnDir(fmt.Sprintf("tar -xvf %s --strip-components=1", templatePath), repoPath)
	if err != nil {
		return "", err
	}

	return repoPath, err
}

func (s *TestHelper) isTemplateKnown(templateName string) bool {
	allTemplates := allTemplates()
	for _, n := range allTemplates {
		if templateName == n {
			return true
		}
	}
	return false
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
