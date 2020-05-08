//go:generate mockgen -package mock -destination ../test/mock/mocks_HttpServer.go gitlab.com/remipassmoilesel/gitsearch/http HttpServer
package http

import (
	"errors"
	"fmt"
	"github.com/markbates/pkger"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"net/http"
	"strconv"
)

type HttpServer interface {
	Start(addr string) error
	// Returns an available host:port string
	GetAvailableAddress() (string, error)
}

type HttpServerImpl struct {
	config     config.Config
	handlers   HttpHandlers
	portHelper PortHelper
}

func NewHttpServer(cfg config.Config, idx index.Index) HttpServer {
	handlers := NewHttpHandlers(cfg, idx)
	server := HttpServerImpl{
		config:     cfg,
		handlers:   handlers,
		portHelper: &PortHelperImpl{},
	}
	server.configure()
	return &server
}

func (s *HttpServerImpl) configure() {
	http.HandleFunc("/api/repository/context", s.handlers.RepositoryContext)
	http.HandleFunc("/api/repository/document", s.handlers.FindDocumentById)
	http.HandleFunc("/api/index/build", s.handlers.BuildIndex)
	http.HandleFunc("/api/index/clean", s.handlers.CleanIndex)
	http.HandleFunc("/api/search", s.handlers.Search)

	staticFiles := http.FileServer(pkger.Dir("/web_client/dist"))
	http.Handle("/", staticFiles)
}

func (s *HttpServerImpl) Start(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func (s *HttpServerImpl) GetAvailableAddress() (string, error) {
	port, err := s.getAvailablePort(s.config.Web.ListenAddress)
	if err != nil {
		return "", err
	}

	addr := s.config.Web.ListenAddress + ":" + strconv.Itoa(port)
	return addr, nil
}

func (s *HttpServerImpl) getAvailablePort(addr string) (int, error) {
	port := s.config.Web.Port
	limit := port + 200
	available, err := s.portHelper.IsPortAvailable(addr, port)
	if err != nil {
		return 0, err
	}
	if available {
		return port, nil
	}

	for !available && port < limit {
		port++

		available, err := s.portHelper.IsPortAvailable(addr, port)
		if err != nil {
			return 0, err
		}
		if available {
			return port, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("no available port found between %v and %v", s.config.Web.Port, limit))
}
