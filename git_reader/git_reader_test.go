package git_reader

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/remipassmoilesel/gitsearch/test"
	"testing"
	"time"
)

func TestGitReader_GitReader_GetHeadHash(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SMALL)
	assert.NoError(t, err)

	hash, err := git.GetHeadHash()
	assert.Equal(t, "8158970149e4199242f0a428e4b7d32b10264dfc", hash)
}

func TestGitReader_GitReader_GetLastsCommits(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SMALL)
	assert.NoError(t, err)

	commits, err := git.GetLastsCommits(2)
	expected := []Commit{
		commitParseDate("8158970149e4199242f0a428e4b7d32b10264dfc", "2020-04-19 10:07:02"),
		commitParseDate("dca73bebbcbf5d3e13b11b87939fb203013dce58", "2020-04-19 10:06:55"),
	}
	assert.Equal(t, expected, commits)
}

func TestGitReader_GitReader_GetLastsCommits_getMoreThanPresent(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SMALL)
	assert.NoError(t, err)

	commits, err := git.GetLastsCommits(10_000)
	expected := []Commit{
		commitParseDate("8158970149e4199242f0a428e4b7d32b10264dfc", "2020-04-19 10:07:02"),
		commitParseDate("dca73bebbcbf5d3e13b11b87939fb203013dce58", "2020-04-19 10:06:55"),
		commitParseDate("70a3f39312261042885d357f60b6eb5b6fd089d4", "2020-04-19 10:06:47"),
		commitParseDate("3d1cc0e431f9ebfd11ebf5ba9885df08b828f6fd", "2020-04-19 10:06:36"),
	}
	assert.Equal(t, expected, commits)
}

func TestGitReader_GitReader_GetCommitsSpacedBy(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SPACED_COMMITS)
	assert.NoError(t, err)

	hashes, err := git.GetCommitsSpacedBy(4, 25*60)
	expected := []Commit{
		commitParseDate("3fb250abfb6cdc82d357498e7c1ee1006537fac1", "2020-05-03 10:14:44"),
		commitParseDate("218816c947a5d4cfc1c287d6d3de4cfe421b3755", "2019-07-28 15:15:00"),
		commitParseDate("3663bbfbd10a9a3f9642165f3e388b212561742f", "2019-07-28 14:45:00"),
		commitParseDate("a7c2e7764f671e8f57ecffca413311dc7f9e46ee", "2019-07-28 14:15:00"),
	}
	assert.Equal(t, expected, hashes)
}

func TestGitReader_GitReader_GetCommitsSpacedBy_getMoreThanExisting(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SPACED_COMMITS)
	assert.NoError(t, err)

	hashes, err := git.GetCommitsSpacedBy(10_000, 25*60)
	expected := []Commit{
		commitParseDate("3fb250abfb6cdc82d357498e7c1ee1006537fac1", "2020-05-03 10:14:44"),
		commitParseDate("218816c947a5d4cfc1c287d6d3de4cfe421b3755", "2019-07-28 15:15:00"),
		commitParseDate("3663bbfbd10a9a3f9642165f3e388b212561742f", "2019-07-28 14:45:00"),
		commitParseDate("a7c2e7764f671e8f57ecffca413311dc7f9e46ee", "2019-07-28 14:15:00"),
		commitParseDate("cff4a64c6817d2d764cd9536a72d442917cdf630", "2019-07-28 13:45:00"),
		commitParseDate("1012388929d15112c42f92a2e2e829512071e9f0", "2019-07-28 13:15:00"),
	}
	assert.Equal(t, expected, hashes)
}

func TestGitReader_GitReader_GetCommitFiles(t *testing.T) {
	git, err := gitReaderOnRepo(test.REPO_SMALL)
	assert.NoError(t, err)

	head, err := git.GetHeadHash()
	assert.NoError(t, err)

	files, err := git.GetCommitFiles(head)
	assert.Len(t, files, 4)

	for _, file := range files {
		assert.NotEmpty(t, file.Path)
		assert.NotEmpty(t, file.Name)
		assert.NotEmpty(t, file.Content)
		assert.NotEmpty(t, file.Hash)
		assert.Equal(t, head, file.Commit)
	}
}

func gitReaderOnRepo(template string) (GitReader, error) {
	repoPath, err := test.Helper.SampleGitRepository(template)
	if err != nil {
		return &GitReaderImpl{}, err
	}

	git, err := NewGitReader(repoPath)
	return git, err
}

func commitParseDate(hash string, dateStr string) Commit {
	date, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		panic(errors.Wrap(err, "invalid date: "+dateStr))
	}
	return Commit{
		Hash: hash,
		Date: date,
	}
}
