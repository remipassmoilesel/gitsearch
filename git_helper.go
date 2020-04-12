package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type GitHelper struct {
}

type FileIterator func(commit *object.Commit, file *object.File) error

func (s *GitHelper) ForEachFiles(path string, maxDepth int, iterator FileIterator) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	firstCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: firstCommit.Hash})
	if err != nil {
		return err
	}

	depth := 0
	commit := firstCommit
	for commit != nil && depth <= maxDepth {
		tree, err := commit.Tree()
		if err != nil {
			return err
		}

		err = tree.Files().ForEach(func(file *object.File) error {
			return iterator(commit, file)
		})

		depth++
		commit, err = commitIter.Next()
	}
	return nil
}
