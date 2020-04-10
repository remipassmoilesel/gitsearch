package main

import (
	"github.com/pkg/errors"
	"os"
	"path"
)

type Config struct {
	DataRootPath string
	Repository   RepositoryContext
	Web          WebConfig
}

type RepositoryContext struct {
	Path     string
	MaxDepth int
}

type WebConfig struct {
	ListenAddress string
	StaticDir     string
}

var DEFAULT_LISTEN_ADDRESS = "127.0.0.1:7777"

func LoadConfig() (Config, error) {
	nearestGitRepo, err := getNearestRepoPath()
	if err != nil {
		return *new(Config), err
	}

	dataRootPath, err := getDataRootPath()
	if err != nil {
		return *new(Config), err
	}

	repoContext := RepositoryContext{
		Path:     nearestGitRepo,
		MaxDepth: 1000,
	}

	port := getListenAddress()
	webConfig := WebConfig{
		ListenAddress: port,
		StaticDir:     path.Join(dataRootPath, "web-client"),
	}

	config := Config{
		Repository:   repoContext,
		DataRootPath: dataRootPath,
		Web:          webConfig,
	}

	return config, err
}

// Return path of firectory where gitsearch write data
func getDataRootPath() (string, error) {
	envPath := os.Getenv("GITSEARCH_HOME")
	if len(envPath) > 0 {
		return envPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".gitsearch"), nil
}

// Searching from the current working directory, get path of nearest git repository
func getNearestRepoPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repoPath := currentDir
	for {
		if repoInfo, err := os.Stat(path.Join(repoPath, ".git")); err == nil && repoInfo.IsDir() {
			break
		}
		if path.Dir(repoPath) == "/" {
			return "", errors.New("no git repository found")
		}
		repoPath = path.Dir(repoPath)
	}
	return repoPath, nil
}

func getListenAddress() string {
	envAddr := os.Getenv("GITSEARCH_LISTEN_ADDRESS")

	if len(envAddr) > 0 {
		return envAddr
	}

	return DEFAULT_LISTEN_ADDRESS
}
