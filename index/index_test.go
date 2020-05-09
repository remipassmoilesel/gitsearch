package index

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/test/mock"
	"testing"
)

func Test_Index_Initialize_locked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	state := mock.NewMockIndexState(ctrl)
	index := IndexImpl{
		state: state,
	}

	state.EXPECT().TryLock().Return(errors.New("index locked"))

	err := index.initialize()
	assert.EqualError(t, err, "cannot initialize index: index locked")
}

func Test_Index_IsUpToDate_upToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	state := mock.NewMockIndexState(ctrl)
	reader := mock.NewMockGitReader(ctrl)
	index := IndexImpl{
		state: state,
		git:   reader,
	}

	reader.EXPECT().GetHeadHash().Return("c", nil)
	state.EXPECT().ContainsCommit("c").Return(true)

	upToDate, _ := index.IsUpToDate()
	assert.True(t, upToDate)
}

func Test_Index_IsUpToDate_notUpToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	state := mock.NewMockIndexState(ctrl)
	reader := mock.NewMockGitReader(ctrl)
	index := IndexImpl{
		state: state,
		git:   reader,
	}

	reader.EXPECT().GetHeadHash().Return("d", nil)
	state.EXPECT().ContainsCommit("d").Return(false)

	upToDate, _ := index.IsUpToDate()
	assert.False(t, upToDate)
}
