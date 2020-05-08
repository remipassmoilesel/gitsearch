package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Config_LoadConfig(t *testing.T) {
	config, err := LoadConfig()
	assert.NoError(t, err)

	assert.NotEmpty(t, config.DataRootPath)

	assert.NotZero(t, config.Index.Shards)
	assert.NotZero(t, config.Index.BatchSize)

	assert.NotEmpty(t, config.Repository.Path)
	assert.NotEmpty(t, config.Repository.Hostname)
	assert.NotZero(t, config.Repository.MaxDepth)
	assert.NotZero(t, config.Repository.Username)

	assert.NotEmpty(t, config.Web.ListenAddress)
}

func Test_Config_LoadConfig_DataRootPath(t *testing.T) {
	err := os.Setenv(ENV_HOME, "/tmp/gitsearch")
	assert.NoError(t, err)

	config, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "/tmp/gitsearch", config.DataRootPath)
}

func Test_Config_LoadConfig_ListenAddress(t *testing.T) {
	err := os.Setenv(ENV_LISTEN_ADDRESS, "0.0.0.0")
	assert.NoError(t, err)

	config, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "0.0.0.0", config.Web.ListenAddress)
}

func Test_Config_LoadConfig_ListenAddress_customPort(t *testing.T) {
	err := os.Setenv(ENV_PORT, "80")
	assert.NoError(t, err)

	config, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, 80, config.Web.Port)
}

func Test_Config_LoadConfig_ListenAddress_wrongCustomPort(t *testing.T) {
	err := os.Setenv(ENV_PORT, "ABC")
	assert.NoError(t, err)

	_, err = LoadConfig()
	assert.EqualError(t, err, "invalid port: ABC")
}
