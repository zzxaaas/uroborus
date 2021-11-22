package service

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"uroborus/model"
	"uroborus/store"
)

type ProjectService struct {
	projectStore *store.ProjectStore
	gitService   *GitService
}

func NewProjectService(projectStore *store.ProjectStore, gitService *GitService) *ProjectService {
	return &ProjectService{
		projectStore: projectStore,
		gitService:   gitService,
	}
}

func (s ProjectService) Save(project *model.Project) error {
	if has, err := s.Get(&model.Project{
		Name: project.Name,
	}); err != nil {
		return err
	} else if has {
		return errors.New("项目名重复")
	}
	project.LocalRepo = viper.GetString("user.rootPath") + project.UserName + "/" + project.Name

	if err := s.gitService.Clone(project.LocalRepo, false, &git.CloneOptions{
		URL:               project.RemoteRepo,
		ReferenceName:     plumbing.NewBranchReferenceName(project.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}); err != nil {
		return err
	}

	if err := os.MkdirAll(project.LocalRepo, os.ModePerm); err != nil {
		return err
	}
	return s.projectStore.Save(project)
}

func (s ProjectService) CheckOut(req *model.Project) error {
	if err := s.projectStore.Update(req); err != nil {
		return err
	}
	if _, err := s.Get(req); err != nil {
		return err
	}
	return s.gitService.Checkout(req.LocalRepo, req.Branch, req.RemoteRepo)
}

func (s ProjectService) Get(project *model.Project) (bool, error) {
	err := s.projectStore.Get(project)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
