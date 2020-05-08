package http

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"testing"
)

func Test_HttpServer_NewHttpServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	_ = NewHttpServer(config.Config{}, idx)
}

func Test_HttpServer_GetAvailableAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portHelper := mock.NewMockPortHelper(ctrl)
	server := HttpServerImpl{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: portHelper,
	}

	portHelper.EXPECT().IsPortAvailable("localhost", gomock.Any()).DoAndReturn(func(host string, port int) (bool, error) {
		return port == 9990, nil
	}).Times(1)

	res, err := server.GetAvailableAddress()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:9990", res)
}

func Test_HttpServer_GetAvailableAddress_shouldFindAvailableAddr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portHelper := mock.NewMockPortHelper(ctrl)
	server := HttpServerImpl{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: portHelper,
	}

	portHelper.EXPECT().IsPortAvailable("localhost", gomock.Any()).DoAndReturn(func(host string, port int) (bool, error) {
		return port == 9995, nil
	}).Times(6)

	res, err := server.GetAvailableAddress()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:9995", res)
}

func Test_HttpServer_GetAvailableAddress_shouldFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portHelper := mock.NewMockPortHelper(ctrl)
	server := HttpServerImpl{
		config: config.Config{
			Web: config.WebConfig{
				ListenAddress: "localhost",
				Port:          9990,
			},
		},
		portHelper: portHelper,
	}

	portHelper.EXPECT().IsPortAvailable("localhost", gomock.Any()).DoAndReturn(func(host string, port int) (bool, error) {
		return port == 20_000, nil
	}).Times(201)

	_, err := server.GetAvailableAddress()
	assert.EqualError(t, err, "no available port found between 9990 and 10190")
}
