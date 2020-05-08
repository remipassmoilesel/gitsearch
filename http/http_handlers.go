package http

import (
	"encoding/json"
	"fmt"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"net/http"
)

type HttpHandlers interface {
	RepositoryContext(w http.ResponseWriter, r *http.Request)
	BuildIndex(w http.ResponseWriter, r *http.Request)
	CleanIndex(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	FindDocumentById(w http.ResponseWriter, r *http.Request)
}

type HttpHandlersImpl struct {
	config config.Config
	index  index.Index
}

func NewHttpHandlers(cfg config.Config, idx index.Index) HttpHandlers {
	return &HttpHandlersImpl{
		config: cfg,
		index:  idx,
	}
}

func (s *HttpHandlersImpl) RepositoryContext(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, s.config.Repository, nil)
}

func (s *HttpHandlersImpl) BuildIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.index.BuildWith(index.BuildOptionsSpacedBy())
	jsonResponse(w, res, err)
}

func (s *HttpHandlersImpl) CleanIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.index.Clean()
	jsonResponse(w, res, err)
}

func (s *HttpHandlersImpl) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := s.index.Search(query, 50, index.OutputHtml)
	jsonResponse(w, results, err)
}

func (s *HttpHandlersImpl) FindDocumentById(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	results, err := s.index.FindDocumentById(hash)
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

	w.Header().Add("Content-Type", "application/json")

	errorObject := map[string]string{"error": err.Error()}
	responseJson, err := json.Marshal(errorObject)

	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(responseJson)
	if err != nil {
		fmt.Println("Http write error: ", err)
	}
}
