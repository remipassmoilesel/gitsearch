package config

import (
	"github.com/pkg/errors"
	"os"
	"os/user"
	"path"
	"runtime"
	"strconv"
)

type Config struct {
	DataRootPath string
	Index        SearchConfig
	Repository   RepositoryContext
	Web          WebConfig
}

type SearchConfig struct {
	Shards    int
	BatchSize int
}

type RepositoryContext struct {
	// Path to the target repository
	Path string
	// Number of commit to inspect
	MaxDepth int
	// Current username
	Username string
	// Current host name
	Hostname string
}

type WebConfig struct {
	ListenAddress string
	Port          int
}

var ENV_HOME = "GITSEARCH_HOME"
var ENV_LISTEN_ADDRESS = "GITSEARCH_LISTEN_ADDRESS"
var ENV_PORT = "GITSEARCH_PORT"

var DEFAULT_LISTEN_ADDRESS = "127.0.0.1"
var DEFAULT_PORT = 7777

func LoadConfig() (Config, error) {
	dataRootPath, err := getDataRootPath()
	if err != nil {
		return *new(Config), err
	}

	repoContext, err := getRepositoryContext()
	if err != nil {
		return *new(Config), err
	}
	repoContext.MaxDepth = 50

	webConfig, err := getWebConfig()
	if err != nil {
		return *new(Config), err
	}

	seachConfig := SearchConfig{Shards: runtime.NumCPU() * 4, BatchSize: 10} // TODO improve

	config := Config{
		DataRootPath: dataRootPath,
		Index:        seachConfig,
		Repository:   repoContext,
		Web:          webConfig,
	}

	return config, err
}

func getWebConfig() (WebConfig, error) {
	address := getListenAddress()
	port, err := getListenPort()
	config := WebConfig{
		ListenAddress: address,
		Port:          port,
	}
	return config, err
}

// Return path of directory where gitsearch write data
func getDataRootPath() (string, error) {
	envPath := os.Getenv(ENV_HOME)
	if len(envPath) > 0 {
		return envPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".gitsearch"), nil
}

func getRepositoryContext() (RepositoryContext, error) {
	nearestGitRepo, err := getNearestRepoPath()
	if err != nil {
		return *new(RepositoryContext), err
	}

	currentUser, err := user.Current()
	if err != nil {
		return *new(RepositoryContext), err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return *new(RepositoryContext), err
	}

	context := RepositoryContext{
		Path:     nearestGitRepo,
		Username: currentUser.Username,
		Hostname: hostname,
	}
	return context, nil
}

// Searching from the current working directory, find then return path of nearest git repository
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
	envAddr := os.Getenv(ENV_LISTEN_ADDRESS)
	if len(envAddr) > 0 {
		return envAddr
	}
	return DEFAULT_LISTEN_ADDRESS
}

func getListenPort() (int, error) {
	envPort := os.Getenv(ENV_PORT)
	if len(envPort) < 1 {
		return DEFAULT_PORT, nil
	}
	port, err := strconv.Atoi(envPort)
	if err != nil {
		return 0, errors.New("invalid port: " + envPort)
	}
	return port, nil
}
