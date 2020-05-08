//go:generate mockgen -package mock -destination ../test/mock/mocks_PortHelper.go gitlab.com/remipassmoilesel/gitsearch/http PortHelper
package http

import (
	"net"
	"strconv"
)

type PortHelper interface {
	IsPortAvailable(addr string, port int) (bool, error)
}

type PortHelperImpl struct {
}

func (s *PortHelperImpl) IsPortAvailable(addr string, port int) (bool, error) {
	ln, err := net.Listen("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return false, nil
	}

	err = ln.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}
