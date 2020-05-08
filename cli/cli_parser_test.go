package cli

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"testing"
)

func Test_CliParserImpl_NewCliParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idx := mock.NewMockIndex(ctrl)
	server := mock.NewMockHttpServer(ctrl)

	_ = NewCliParser(config.Config{}, idx, server)
}

func Test_CliParserImpl_ApplyCommand_noCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().StartServer().Times(1)

	err := parser.ApplyCommand(fakeArgs(""))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().Search(gomock.Eq("query"), gomock.Eq(10), gomock.Eq(true)).Times(1)

	err := parser.ApplyCommand(fakeArgs("search", "-q", "query", "-n", "10"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_search_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().Search(gomock.Eq("query"), gomock.Eq(10), gomock.Eq(true)).Times(1).Return(errors.New("unexpected error"))

	err := parser.ApplyCommand(fakeArgs("search", "-q", "query", "-n", "10"))
	assert.EqualError(t, err, "unexpected error")
}

func Test_CliParserImpl_ApplyCommand_search_short(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().Search(gomock.Eq("query"), gomock.Eq(10), gomock.Eq(true)).Times(1)

	err := parser.ApplyCommand(fakeArgs("s", "-q", "query", "-n", "10"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_showFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().ShowFile(gomock.Eq("hash"), gomock.Eq(true)).Times(1)

	err := parser.ApplyCommand(fakeArgs("show-file", "-ha", "hash"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_showFile_short(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().ShowFile(gomock.Eq("hash"), gomock.Eq(true)).Times(1)

	err := parser.ApplyCommand(fakeArgs("f", "-ha", "hash"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_webUi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().StartServer().Times(1)

	err := parser.ApplyCommand(fakeArgs("web-ui"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_webUi_short(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().StartServer().Times(1)

	err := parser.ApplyCommand(fakeArgs("w"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_index_clean(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().CleanIndex().Times(1)

	err := parser.ApplyCommand(fakeArgs("index", "clean"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_index_clean_short(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().CleanIndex().Times(1)

	err := parser.ApplyCommand(fakeArgs("i", "c"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_index_update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().UpdateIndex().Times(1)

	err := parser.ApplyCommand(fakeArgs("index", "update"))
	assert.NoError(t, err)
}

func Test_CliParserImpl_ApplyCommand_index_update_short(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handlers := mock.NewMockCliHandlers(ctrl)
	parser := CliParserImpl{
		handlers: handlers,
	}

	handlers.EXPECT().UpdateIndex().Times(1)

	err := parser.ApplyCommand(fakeArgs("i", "u"))
	assert.NoError(t, err)
}

func fakeArgs(args ...string) []string {
	return append([]string{"/tmp/executable"}, args...)
}
