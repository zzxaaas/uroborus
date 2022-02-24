package service

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strings"
	"uroborus/model"
	"uroborus/store"
)

type ProjectService struct {
	projectStore *store.ProjectStore
	imageService *BaseImageService
	gitService   *GitService
}

func NewProjectService(
	projectStore *store.ProjectStore,
	imageService *BaseImageService,
	gitService *GitService,
) *ProjectService {
	return &ProjectService{
		projectStore: projectStore,
		imageService: imageService,
		gitService:   gitService,
	}
}

func (s ProjectService) Save(req *model.RegisterProjectReq) error {
	if req.Type == "tool" || req.Name == "" {
		req.Name = fmt.Sprintf("%s-%s", strings.Split(req.Image, ":")[0], req.UserName)
	}
	if req.Env != nil {
		req.Project.Env = strings.Join(req.Env, ",")
	}
	if req.Cmd != nil {
		req.Project.Command = strings.Join(req.Cmd, ",")
	}

	imageName := strings.Split(req.Image, ":")[0]
	if port, err := s.getImagePort(imageName); err != nil {
		return err
	} else if port != "" {
		req.Port = port
	}

	has, err := s.Get(&model.Project{
		Name: req.Name,
	})
	if err != nil {
		return err
	}

	if !has {
		if req.Type != "tool" {
			if err := s.initProjectPath(&req.Project); err != nil {
				return err
			}
			if err := s.clone(&req.Project); err != nil {
				return err
			}
		}
		if err := s.generatePort(&req.Project); err != nil {
			return err
		}
	}

	return s.projectStore.Save(&req.Project)
}

func (s ProjectService) Delete(req *model.Project) error {
	return s.projectStore.Delete(req)
}
func (s ProjectService) Update(req *model.Project) error {
	return s.projectStore.Update(req)
}

func (s ProjectService) generatePort(project *model.Project) error {
	for {
		project.BindPort = rand.Intn(65535)
		err := s.projectStore.Get(&model.Project{
			BindPort: project.BindPort,
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return err
		}
	}
	return nil
}

func (s ProjectService) getImagePort(name string) (string, error) {
	resp, err := s.imageService.Get(model.BaseImage{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", errors.New("image not found")
	}
	return resp[0].Port, nil
}

func (s ProjectService) initProjectPath(project *model.Project) error {
	root := viper.GetString("root")
	basePath := fmt.Sprintf("%s/%s/%s/%s", root, project.UserName, project.Branch, project.Name)

	project.LocalRepo = basePath + model.RepoBasePath
	dockerfilePath := basePath + model.DockerfileBasePath

	for _, path := range []string{project.LocalRepo, dockerfilePath} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

//func (s ProjectService) CheckOut(req *model.Project) error {
//	if err := s.projectStore.Update(req); err != nil {
//		return err
//	}
//	if _, err := s.Get(req); err != nil {
//		return err
//	}
//	return s.gitService.Checkout(req.LocalRepo, req.Branch, req.RemoteRepo)
//}

func (s ProjectService) Find(project model.Project) ([]model.Project, error) {
	return s.projectStore.Find(project)
}

func (s ProjectService) clone(project *model.Project) error {
	if err := s.gitService.Clone(project.LocalRepo, false, &git.CloneOptions{
		URL:               project.RemoteRepo,
		ReferenceName:     plumbing.NewBranchReferenceName(project.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
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
