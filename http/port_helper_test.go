package http

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
)

const TEST_PORT = 56_885

func Test_PortHelperImpl_IsPortAvailable_portAvailable(t *testing.T) {
	portH := PortHelperImpl{}
	res, err := portH.IsPortAvailable("localhost", TEST_PORT)
	assert.NoError(t, err)
	assert.True(t, res)
}

func Test_PortHelperImpl_IsPortAvailable_portNotAvailable(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:"+strconv.Itoa(TEST_PORT))
	assert.NoError(t, err)
	defer func() {
		err := ln.Close()
		if err != nil {
			panic(err)
		}
	}()

	portH := PortHelperImpl{}
	res, err := portH.IsPortAvailable("localhost", TEST_PORT)
	assert.NoError(t, err)
	assert.False(t, res)
}
