package main

import (
	"net/http"
)

type HttpServer struct {
	config   Config
	handlers HttpHandlers
}

func NewHttpServer(config Config, index Index) HttpServer {

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

	staticFiles := http.FileServer(http.Dir(s.config.Web.StaticDir))
	http.Handle("/", staticFiles)
}

func (s HttpServer) Start() error {
	return http.ListenAndServe(s.config.Web.ListenAddress, nil)
}
