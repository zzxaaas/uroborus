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
	"strconv"
	"strings"
	"uroborus/model"
	"uroborus/store"
)

type ProjectService struct {
	projectStore     *store.ProjectStore
	gitService       *GitService
	imageService     *BaseImageService
	containerService *ContainerService
}

func NewProjectService(
	projectStore *store.ProjectStore,
	gitService *GitService,
	imageService *BaseImageService,
	containerService *ContainerService) *ProjectService {
	return &ProjectService{
		projectStore:     projectStore,
		gitService:       gitService,
		imageService:     imageService,
		containerService: containerService,
	}
}

func (s ProjectService) Save(req *model.RegisterProjectReq) error {
	if req.Type == "tool" || req.Name == "" {
		req.Name = fmt.Sprintf("%s-%s", strings.Split(req.Image, ":")[0], req.UserName)
	}
	if req.Env != nil {
		req.Project.Env = strings.Join(req.Env, ",")
	}

	imageName := strings.Split(req.Image, ":")[0]
	if port, err := s.getImagePort(imageName); err != nil {
		return err
	} else {
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
			go s.cloneFromGit(req.Project)
		}
		if err := s.generatePort(&req.Project); err != nil {
			return err
		}
	}

	return s.projectStore.Save(&req.Project)
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

func (s ProjectService) cloneFromGit(project model.Project) {
	if err := s.gitService.Clone(project.LocalRepo, false, &git.CloneOptions{
		URL:               project.RemoteRepo,
		ReferenceName:     plumbing.NewBranchReferenceName(project.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}); err != nil {
		logrus.Error(err)
	}
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
	basePath := fmt.Sprintf("%s/%s/%s", root, project.UserName, project.Name)

	project.LocalRepo = basePath + model.RepoBasePath
	dockerfilePath := basePath + model.DockerfileBasePath

	for _, path := range []string{project.LocalRepo, dockerfilePath} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (s ProjectService) prepareDockerfile() {

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

func (s ProjectService) Build(req model.Project) error {
	if err := s.projectStore.Get(&req); err != nil {
		return err
	}

	if req.Type == "project" {
		if err := s.containerService.BuildImage(model.BuildImageOption{
			Path:       req.LocalRepo,
			Dockerfile: req.Dockerfile,
			Tag:        req.Name,
			Version:    req.Version,
		}); err != nil {
			return err
		}
	}

	if req.Container != "" {
		if err := s.containerService.RemoveContainer(req.Container); err != nil {
			return err
		}
	}

	if id, err := s.containerService.StartContainerWithOption(model.ContainerOption{
		Name:      req.Name,
		Image:     req.Image,
		ProtoPort: req.Port,
		Port:      strconv.Itoa(req.BindPort),
		Env:       strings.Split(req.Env, ","),
	}); err != nil {
		return err
	} else {
		req.Container = id
	}
	if err := s.projectStore.Save(&req); err != nil {
		return err
	}
	return nil
}
