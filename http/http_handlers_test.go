package http

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	mock "gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_HttpHandlersImpl_RepositoryContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	handlers := &HttpHandlersImpl{
		config: config.Config{
			Repository: config.RepositoryContext{
				Path:     "/path/to/repo",
				MaxDepth: 5,
				Username: "username-value",
				Hostname: "hostname-value",
			},
		},
		index: idx,
	}

	req := newTestRequest(t, "GET", "/", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RepositoryContext)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"Path\":\"/path/to/repo\",\"MaxDepth\":5,\"Username\":\"username-value\",\"Hostname\":\"hostname-value\"}"
	assertJsonResponse(t, expectedRes, res)
}

func Test_HttpHandlersImpl_BuildIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	handlers := &HttpHandlersImpl{
		config: config.Config{},
		index:  idx,
	}

	req := newTestRequest(t, "GET", "/", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.BuildIndex)

	idx.EXPECT().BuildWith(gomock.Any()).Times(1).Return(domain.BuildOperationResult{
		TookSeconds:  1,
		Files:        2,
		TotalFiles:   3,
		OldestCommit: time.Time{},
	}, nil)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"TookSeconds\":1,\"Files\":2,\"TotalFiles\":3,\"OldestCommit\":\"0001-01-01T00:00:00Z\"}"
	assertJsonResponse(t, expectedRes, res)
}

func Test_HttpHandlersImpl_CleanIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	handlers := &HttpHandlersImpl{
		config: config.Config{},
		index:  idx,
	}

	req := newTestRequest(t, "GET", "/", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CleanIndex)

	idx.EXPECT().Clean().Times(1).Return(domain.CleanOperationResult{TookMs: 1}, nil)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"TookMs\":1}"
	assertJsonResponse(t, expectedRes, res)
}

func Test_HttpHandlersImpl_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	handlers := &HttpHandlersImpl{
		config: config.Config{},
		index:  idx,
	}

	req := newTestRequest(t, "GET", "/search?query=search%20query", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Search)

	idx.EXPECT().Search(gomock.Eq("search query"), gomock.Eq(50), gomock.Eq(index.OutputHtml)).Times(1).Return(domain.SearchResult{
		Query:  "search query",
		TookMs: 5,
		Matches: []domain.SearchMatch{
			{
				File: domain.IndexedFile{
					Hash:    "hash",
					Commit:  "commit",
					Date:    time.Date(2015, 05, 31, 0, 0, 0, 0, time.UTC),
					Content: "content",
					Name:    "",
					Path:    "",
				},
				Fragments: []string{"match 1", "match 2"},
			},
		},
	}, nil)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"Query\":\"search query\",\"TookMs\":5,\"Matches\":[{\"File\":{\"Hash\":\"hash\",\"Commit\":\"commit\",\"Date\":\"2015-05-31T00:00:00Z\",\"Content\":\"content\",\"Name\":\"\",\"Path\":\"\"},\"Fragments\":[\"match 1\",\"match 2\"]}]}"
	assertJsonResponse(t, expectedRes, res)
}

func Test_HttpHandlersImpl_FindDocumentById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	handlers := &HttpHandlersImpl{
		config: config.Config{},
		index:  idx,
	}

	req := newTestRequest(t, "GET", "/document?id=document-id", nil)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.FindDocumentById)

	idx.EXPECT().FindDocumentById(gomock.Eq("document-id")).Times(1).Return(domain.IndexedFile{
		Hash:    "hash",
		Commit:  "commit",
		Date:    time.Date(2015, 05, 31, 0, 0, 0, 0, time.UTC),
		Content: "content",
		Name:    "name",
		Path:    "path",
	}, nil)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"Hash\":\"hash\",\"Commit\":\"commit\",\"Date\":\"2015-05-31T00:00:00Z\",\"Content\":\"content\",\"Name\":\"name\",\"Path\":\"path\"}"
	assertJsonResponse(t, expectedRes, res)
}

func Test_HttpHandlersImpl_jsonError(t *testing.T) {
	err := errors.New("unexpected error")
	res := httptest.NewRecorder()
	jsonError(res, err)

	assert.Equal(t, res.Body.String(), "{\"error\":\"unexpected error\"}")
	assert.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, res.Header())
	assert.Equal(t, 500, res.Code)
}

func assertJsonResponse(t *testing.T, expected string, actual *httptest.ResponseRecorder) {
	assert.Equal(t, actual.Body.String(), expected)
	assert.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, actual.Header())
	assert.Equal(t, 200, actual.Code)
}

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
