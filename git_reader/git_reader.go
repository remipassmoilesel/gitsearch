//go:generate mockgen -package mock -destination ../test/mock/mocks_GitReader.go gitlab.com/remipassmoilesel/gitsearch/git_reader GitReader
package git_reader

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"time"
)

type GitReader interface {
	GetHeadHash() (string, error)
	GetCommitFiles(commitStr string) ([]CommitFile, error)
	GetLastsCommits(commitNbr int) ([]Commit, error)
	GetCommitsSpacedBy(commitNbr int, intervalSec float64) ([]Commit, error)
}

// Warning: methods are not thread safe
type GitReaderImpl struct {
	repo *git.Repository
}

type CommitFile struct {
	// File hash
	Hash string
	// Commit hash
	Commit string
	// Date of commit
	Date time.Time
	// File content
	Content string
	// File path
	Path string
	// File name
	Name string
}

type Commit struct {
	Hash string
	Date time.Time
}

func NewGitReader(path string) (GitReader, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return &GitReaderImpl{}, errors.Wrap(err, "cannot open git repository")
	}

	reader := GitReaderImpl{
		repo: repo,
	}
	return &reader, nil
}

func (s *GitReaderImpl) GetHeadHash() (string, error) {
	head, err := s.repo.Head()
	if err != nil {
		return "", errors.Wrap(err, "cannot read head hash")
	}

	return head.Hash().String(), nil
}

func (s *GitReaderImpl) GetCommitFiles(commitStr string) ([]CommitFile, error) {
	commitIter, err := s.repo.Log(&git.LogOptions{From: plumbing.NewHash(commitStr)})
	if err != nil {
		return []CommitFile{}, errors.Wrap(err, "cannot access git history")
	}
	defer commitIter.Close()

	commit, err := commitIter.Next()
	if err != nil {
		return []CommitFile{}, errors.Wrap(err, "cannot access git history")
	}

	tree, err := commit.Tree()
	if err != nil {
		return []CommitFile{}, errors.Wrap(err, "cannot access git history")
	}

	files := []CommitFile{}
	err = tree.Files().ForEach(func(file *object.File) error {
		content, ferr := file.Contents()
		if ferr != nil {
			return errors.Wrap(ferr, "cannot read commit file")
		}

		commitFile := CommitFile{
			Hash:    file.Hash.String(),
			Commit:  commitStr,
			Date:    commit.Committer.When.UTC(),
			Content: content,
			Path:    file.Name,
			Name:    filepath.Base(file.Name),
		}

		files = append(files, commitFile)
		return nil
	})

	return files, errors.Wrap(err, "cannot read git history")
}

// Reading is done from the most recent to the oldest commit
func (s *GitReaderImpl) GetLastsCommits(commitNbr int) ([]Commit, error) {
	head, err := s.repo.Head()
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits, repository is empty")
	}

	firstCommit, err := s.repo.CommitObject(head.Hash())
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits")
	}

	commitIter, err := s.repo.Log(&git.LogOptions{From: firstCommit.Hash})
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits")
	}
	defer commitIter.Close()

	res := []Commit{}
	for i := 0; i < commitNbr; i++ {
		commit, err := commitIter.Next()
		if commit == nil {
			break
		}
		if err != nil {
			return []Commit{}, errors.Wrap(err, "cannot get last commits")
		}
		res = append(res, Commit{
			Hash: commit.Hash.String(),
			Date: commit.Committer.When.UTC(),
		})
	}
	return res, nil
}

// Reading is done from the most recent to the oldest commit
func (s *GitReaderImpl) GetCommitsSpacedBy(commitNbr int, intervalSec float64) ([]Commit, error) {
	head, err := s.repo.Head()
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits, repository is empty")
	}

	firstCommit, err := s.repo.CommitObject(head.Hash())
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits")
	}

	commitIter, err := s.repo.Log(&git.LogOptions{From: firstCommit.Hash})
	if err != nil {
		return []Commit{}, errors.Wrap(err, "cannot get last commits")
	}
	defer commitIter.Close()

	lastDate := time.Time{}
	res := []Commit{}
	for i := 0; len(res) < commitNbr; i++ {
		commit, err := commitIter.Next()
		if commit == nil {
			break
		}
		if err != nil {
			return []Commit{}, errors.Wrap(err, "cannot get last commits")
		}
		commitDate := commit.Committer.When
		if lastDate.Equal(time.Time{}) || lastDate.Sub(commitDate).Seconds() > intervalSec {
			res = append(res, Commit{
				Hash: commit.Hash.String(),
				Date: commit.Committer.When.UTC(),
			})
			lastDate = commit.Committer.When.UTC()
		}
	}
	return res, nil
}
