package http

import (
	"encoding/json"
	"fmt"
	config "github.com/remipassmoilesel/gitsearch/config"
	index "github.com/remipassmoilesel/gitsearch/index"
	"net/http"
)

type HttpHandlers struct {
	config config.Config
	index  index.Index
}

func (s *HttpHandlers) RepositoryContext(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, s.config.Repository, nil)
}

func (s *HttpHandlers) BuildIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.index.Build()
	jsonResponse(w, res, err)
}

func (s *HttpHandlers) CleanIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.index.Clean()
	jsonResponse(w, res, err)
}

func (s *HttpHandlers) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := s.index.Search(query, 50, index.OutputHtml)
	jsonResponse(w, results, err)
}

func jsonResponse(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		jsonError(w, err)
		return
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		jsonError(w, err)
		return
	}

	_, err = w.Write(responseJson)
	if err != nil {
		jsonError(w, err)
		return
	}
}

func jsonError(w http.ResponseWriter, err error) {
	fmt.Println("Http error: ", err)

	errorObject := map[string]string{"error": err.Error()}
	responseJson, err := json.Marshal(errorObject)

	_, err = w.Write(responseJson)
	if err != nil {
		fmt.Println("Http write error: ", err)
	}
}
