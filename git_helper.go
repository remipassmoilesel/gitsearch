package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type GitHelper struct {
}

type FileIterator func(commit *object.Commit, file *object.File) error

func (s *GitHelper) ForEachFiles(path string, iterator FileIterator) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: commit.Hash})
	if err != nil {
		return err
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		tree, err := commit.Tree()
		if err != nil {
			return err
		}

		err = tree.Files().ForEach(func(file *object.File) error {
			err = iterator(commit, file)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
