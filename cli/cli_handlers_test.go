package cli

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"testing"
)

func Test_CliHandlers_UpdateIndex_NewCliHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	server := mock.NewMockHttpServer(ctrl)

	_ = NewCliHandlers(config.Config{}, idx, server)
}

func Test_CliHandlers_UpdateIndex_shouldUpdateIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	idx.EXPECT().IsUpToDate().Return(false, nil).AnyTimes()
	idx.EXPECT().BuildWith(gomock.Any()).Times(1)
	cliDisplay.EXPECT().IndexBuild(gomock.Any()).Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)

	err := handlers.UpdateIndex()
	assert.NoError(t, err)
}

func Test_CliHandlers_UpdateIndex_shouldNotUpdateIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	idx.EXPECT().IsUpToDate().Return(true, nil).AnyTimes()
	idx.EXPECT().BuildWith(gomock.Any()).Times(0)
	cliDisplay.EXPECT().IndexBuild(gomock.Any()).Times(0)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(0)

	err := handlers.UpdateIndex()
	assert.NoError(t, err)
}

func Test_CliHandlers_CleanIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	idx.EXPECT().Clean().Times(1)
	cliDisplay.EXPECT().IndexClean(gomock.Any()).Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)

	err := handlers.CleanIndex()
	assert.NoError(t, err)
}

func Test_CliHandlers_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	idx.EXPECT().IsUpToDate().Return(true, nil).AnyTimes()
	idx.EXPECT().Search("query", 10, gomock.Any()).Times(1)
	cliDisplay.EXPECT().Search(gomock.Any()).Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)

	err := handlers.Search("query", 10, true)
	assert.NoError(t, err)
}

func Test_CliHandlers_Search_tooMuchResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	searchRes := index.SearchResult{
		Query:  "query",
		TookMs: 5,
		Matches: []index.SearchMatch{
			{
				File:      index.IndexedFile{},
				Fragments: []string{"match 1"},
			},
			{
				File:      index.IndexedFile{},
				Fragments: []string{"match 2"},
			},
		},
	}

	reducedRes := index.SearchResult{
		Query:  "query",
		TookMs: 5,
		Matches: []index.SearchMatch{
			{
				File:      index.IndexedFile{},
				Fragments: []string{"match 1"},
			},
		},
	}

	idx.EXPECT().IsUpToDate().Return(true, nil).AnyTimes()
	idx.EXPECT().Search("query", 1, gomock.Any()).Times(1).Return(searchRes, nil)
	cliDisplay.EXPECT().Search(gomock.Eq(reducedRes)).Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)

	err := handlers.Search("query", 1, true)
	assert.NoError(t, err)
}

func Test_CliHandlers_ShowFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
	}

	idx.EXPECT().IsUpToDate().Return(true, nil).AnyTimes()
	idx.EXPECT().FindDocumentById("hash").Times(1)
	cliDisplay.EXPECT().ShowFile(gomock.Any()).Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)

	err := handlers.ShowFile("hash", true)
	assert.NoError(t, err)
}

func Test_CliHandlers_StartServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	cliDisplay := mock.NewMockCliDisplay(ctrl)
	server := mock.NewMockHttpServer(ctrl)

	handlers := &CliHandlersImpl{
		config:  config.Config{},
		index:   idx,
		display: cliDisplay,
		server:  server,
	}

	idx.EXPECT().IsUpToDate().Return(true, nil).AnyTimes()
	server.EXPECT().GetAvailableAddress().Return("localhost:7778", nil)
	cliDisplay.EXPECT().StartServer("http://localhost:7778").Times(1)
	cliDisplay.EXPECT().Display(gomock.Any(), gomock.Any()).Times(1)
	server.EXPECT().Start("localhost:7778").Times(1)

	err := handlers.StartServer()
	assert.NoError(t, err)
}
