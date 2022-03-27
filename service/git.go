package service

import (
	"github.com/go-git/go-git/v5"
)

//import "github.com/go-git/go-git/v5"

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

func (s GitService) Clone(directory string, isBare bool, options *git.CloneOptions) error {
	r, err := git.PlainClone(directory, isBare, options)
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	_, err = r.CommitObject(ref.Hash())
	return err
}

func (s GitService) Pull(directory string) error {
	r, err := git.PlainOpen(directory)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return err
	}
	w, err := r.Worktree()
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	return nil
}
