package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
)

type GitHelper struct {
}

type FileIterator func(bundle *CommitBundle) error

type CommitBundle struct {
	Commit string
	Files  []CommitBundleFile
}

type CommitBundleFile struct {
	Path    string
	Hash    string
	Name    string
	Content string
}

// Return an object containing commit informations and files
// Objects returned are thread safe
// TODO: we should limit bundle size
func (s *GitHelper) ForEachCommitBundle(path string, maxDepth int, iterator FileIterator) error {
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

		bundleFiles := []CommitBundleFile{}
		err = tree.Files().ForEach(func(file *object.File) error {
			content, ferr := file.Contents()
			if ferr != nil {
				fmt.Println("Indexing error. Cannot read file: ", ferr)
				return nil
			}

			bundleFile := CommitBundleFile{
				Hash:    file.Hash.String(),
				Content: content,
				Path:    file.Name,
				Name:    filepath.Base(file.Name),
			}

			bundleFiles = append(bundleFiles, bundleFile)
			return nil
		})
		if err != nil {
			return err
		}

		commitBundle := CommitBundle{
			Commit: commit.Hash.String(),
			Files:  bundleFiles,
		}

		err = iterator(&commitBundle)
		if err != nil {
			return err
		}

		depth++
		commit, err = commitIter.Next()
	}

	return nil
}

func (s *GitHelper) GetLastCommitHash(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}
