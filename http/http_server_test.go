package http

import (
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HttpServer_GetAvailableAddress(t *testing.T) {
	server := HttpServer{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: &TestPortHelper{port: 9990},
	}
	res, err := server.GetAvailableAddress()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:9990", res)
}

func Test_HttpServer_GetAvailableAddress_shouldFindAvailableAddr(t *testing.T) {
	server := HttpServer{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: &TestPortHelper{port: 9995},
	}
	res, err := server.GetAvailableAddress()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:9995", res)
}

func Test_HttpServer_GetAvailableAddress_shouldFail(t *testing.T) {
	server := HttpServer{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: &TestPortHelper{port: 20_000},
	}
	_, err := server.GetAvailableAddress()
	assert.EqualError(t, err, "no available port found between 9990 and 10190")
}

type TestPortHelper struct {
	port int
}

func (s *TestPortHelper) IsPortAvailable(addr string, port int) (bool, error) {
	return port == s.port, nil
}
