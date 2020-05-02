package utils

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"time"
)

// Warning: methods are not thread safe
type GitReader struct {
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

func NewGitReader(path string) (GitReader, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return GitReader{}, errors.Wrap(err, "cannot open git repository")
	}

	reader := GitReader{
		repo: repo,
	}
	return reader, nil
}

func (s *GitReader) GetHeadHash() (string, error) {
	head, err := s.repo.Head()
	if err != nil {
		return "", errors.Wrap(err, "cannot read head hash")
	}

	return head.Hash().String(), nil
}

func (s *GitReader) GetCommitFiles(commitStr string) ([]CommitFile, error) {
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
			Date:    commit.Committer.When.Local(),
			Content: content,
			Path:    file.Name,
			Name:    filepath.Base(file.Name),
		}

		files = append(files, commitFile)
		return nil
	})

	return files, errors.Wrap(err, "cannot read git history")
}

// Reading MUST be done from the most recent to the oldest
func (s *GitReader) GetLastsCommits(commitNbr int) ([]string, error) {
	head, err := s.repo.Head()
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot get last commits, repository is empty")
	}

	firstCommit, err := s.repo.CommitObject(head.Hash())
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot get last commits")
	}

	commitIter, err := s.repo.Log(&git.LogOptions{From: firstCommit.Hash})
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot get last commits")
	}
	defer commitIter.Close()

	res := []string{}
	for i := 0; i < commitNbr; i++ {
		commit, err := commitIter.Next()
		if commit == nil {
			break
		}
		if err != nil {
			return []string{}, errors.Wrap(err, "cannot get last commits")
		}
		res = append(res, commit.Hash.String())
	}
	return res, nil
}
