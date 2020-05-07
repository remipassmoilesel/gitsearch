package cli

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	mock "gitlab.com/remipassmoilesel/gitsearch/mocks"
	"testing"
)

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
	cliDisplay.EXPECT().Search("query", gomock.Any(), true).Times(1)

	err := handlers.Search("query", 10, true)
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
	cliDisplay.EXPECT().ShowFile(gomock.Any(), true).Times(1)

	err := handlers.ShowFile("hash", true)
	assert.NoError(t, err)
}
