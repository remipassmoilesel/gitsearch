package http

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	mock "gitlab.com/remipassmoilesel/gitsearch/mocks"
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

	idx.EXPECT().BuildWith(gomock.Any()).Times(1).Return(index.BuildOperationResult{
		TookSeconds:  1,
		Files:        2,
		TotalFiles:   3,
		OldestCommit: time.Time{},
	}, nil)

	handler.ServeHTTP(res, req)

	expectedRes := "{\"TookSeconds\":1,\"Files\":2,\"TotalFiles\":3,\"OldestCommit\":\"0001-01-01T00:00:00Z\"}"
	assertJsonResponse(t, expectedRes, res)
}

func assertJsonResponse(t *testing.T, expected string, actual *httptest.ResponseRecorder) {
	assert.Equal(t, expected, actual.Body.String())
	assert.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, actual.Header())
}

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
