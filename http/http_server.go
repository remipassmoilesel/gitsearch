package http

import (
	"github.com/markbates/pkger"
	config "github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/index"
	"net/http"
)

type HttpServer struct {
	config   config.Config
	handlers HttpHandlers
}

func NewHttpServer(config config.Config, index index.Index) HttpServer {

	httpHandlers := HttpHandlers{config, index}
	server := HttpServer{
		config:   config,
		handlers: httpHandlers,
	}
	server.configure()

	return server
}

func (s HttpServer) configure() {
	http.HandleFunc("/api/repository/context", s.handlers.RepositoryContext)
	http.HandleFunc("/api/index/build", s.handlers.BuildIndex)
	http.HandleFunc("/api/index/clean", s.handlers.CleanIndex)
	http.HandleFunc("/api/search", s.handlers.Search)

	staticFiles := http.FileServer(pkger.Dir("/web_client/dist"))
	http.Handle("/", staticFiles)
}

func (s HttpServer) Start() error {
	return http.ListenAndServe(s.config.Web.ListenAddress, nil)
}
