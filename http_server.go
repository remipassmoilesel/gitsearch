package main

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	config Config
}

func NewHttpServer(config Config) (HttpServer, error) {
	server := HttpServer{
		config,
	}

	return server, nil
}

func (s HttpServer) Start() error {
	http.HandleFunc("/api/search", onSearch)

	staticFiles := http.FileServer(http.Dir(s.config.Web.StaticDir))
	http.Handle("/", staticFiles)

	return http.ListenAndServe(s.config.Web.ListenAddress, nil)
}

func onSearch(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello !")
	if err != nil {
		fmt.Println(err)
	}
}
