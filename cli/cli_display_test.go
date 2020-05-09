package cli

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ttacon/chalk"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/domain"
	"gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"testing"
	"time"
)

func Test_CliDisplayImpl_NewCliDisplay(t *testing.T) {
	_ = NewCliDisplay(config.Config{})
}

func Test_CliDisplayImpl_Display(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	utils.EXPECT().PrintLn("output").Times(1)

	err := cliDisplay.Display("output", false)
	assert.NoError(t, err)
}

func Test_CliDisplayImpl_Display_WithPager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	utils.EXPECT().Pager("output").Times(1)

	err := cliDisplay.Display("output", true)
	assert.NoError(t, err)
}

func Test_CliDisplayImpl_Display_WithPager_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	utils.EXPECT().Pager("output").Times(1).Return(errors.New("unexpected error"))
	p1 := utils.EXPECT().PrintLn("*** You should install a pager and set $PAGER variable in your shell ***").Times(1)
	utils.EXPECT().PrintLn("output").Times(1).After(p1)

	err := cliDisplay.Display("output", true)
	assert.Error(t, err)
}

func Test_CliDisplayImpl_IndexBuild(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.IndexBuild(domain.BuildOperationResult{
		TookSeconds:  1,
		Files:        2,
		TotalFiles:   3,
		OldestCommit: time.Time{},
	})
	assert.Equal(t, "Indexed 2/3 files in 1 seconds. Oldest commit: 0001-01-01 00:00:00 +0000 UTC", output)
}

func Test_CliDisplayImpl_IndexClean(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.IndexClean(domain.CleanOperationResult{
		TookMs: 1,
	})
	assert.Equal(t, "Index clean took 1 ms", output)
}

func Test_CliDisplayImpl_StartServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.StartServer("http://localhost:8889")
	assert.Equal(t, "Listening on http://localhost:8889\nPress CTRL+C to quit", output)
}

func Test_CliDisplayImpl_ShowFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.ShowFile(domain.IndexedFile{
		Hash:    "508eaccae7b32fa2407d78757cdc850ab908ff8d",
		Commit:  "94e058175755f4a5a479a94d4a87492a0c8ca9bc",
		Date:    time.Time{},
		Content: "content",
		Name:    "name",
		Path:    "path",
	})

	expected := "\n" + chalk.Cyan.String() + "path - Commit: 94e058175755f4a - Date: 0001-01-01 00:00:00 +0000 UTC - Id: 508eaccae7b32fa" + chalk.Reset.String() + "\ncontent\n"
	assert.Equal(t, expected, output)
}

func Test_CliDisplayImpl_Search_2results(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.Search(domain.SearchResult{
		Query:  "query",
		TookMs: 5,
		Matches: []domain.SearchMatch{
			{
				File: domain.IndexedFile{
					Hash:    "1-508eaccae7b32fa2407d78757cdc850ab908ff8d",
					Commit:  "1-94e058175755f4a5a479a94d4a87492a0c8ca9bc",
					Date:    time.Time{},
					Content: "1-content",
					Name:    "1-name",
					Path:    "1-path",
				},
				Fragments: []string{
					"1. fragment 1",
					"1. fragment 2",
				},
			},
			{
				File: domain.IndexedFile{
					Hash:    "2-508eaccae7b32fa2407d78757cdc850ab908ff8d",
					Commit:  "2-94e058175755f4a5a479a94d4a87492a0c8ca9bc",
					Date:    time.Time{},
					Content: "2-content",
					Name:    "2-name",
					Path:    "2-path",
				},
				Fragments: []string{
					"2. fragment 1",
					"2. fragment 2",
				},
			},
		},
	})
	expected := "Query: query\n\n\x1b[36m🎀️  1: 1-path - Commit: 1-94e058175755f - Date: 0001-01-01 00:00:00 +0000 UTC - Id: 1-508eaccae7b32\x1b[49m\x1b[39m\n    1. fragment 1\n    \n    ---\n    \n    1. fragment 2\n\n\x1b[36m➡️  2: 2-path - Commit: 2-94e058175755f - Date: 0001-01-01 00:00:00 +0000 UTC - Id: 2-508eaccae7b32\x1b[49m\x1b[39m\n    2. fragment 1\n    \n    ---\n    \n    2. fragment 2\nSearch took 5 ms"
	assert.Equal(t, expected, output)
}

func Test_CliDisplayImpl_Search_noResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils := mock.NewMockUtils(ctrl)
	cliDisplay := CliDisplayImpl{
		config: config.Config{},
		utils:  utils,
	}

	output := cliDisplay.Search(domain.SearchResult{
		Query:   "query",
		TookMs:  5,
		Matches: []domain.SearchMatch{},
	})
	expected := "Query: query\n\nNothing found.Search took 5 ms"
	assert.Equal(t, expected, output)
}
