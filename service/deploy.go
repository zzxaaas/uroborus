package service

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
	"uroborus/model"
)

type DeployService struct {
	deployHistoryService *DeployHistoryService
	projectService       *ProjectService
	gitService           *GitService
	containerService     *ContainerService
}

func NewDeployService(deployHistoryService *DeployHistoryService, projectService *ProjectService, gitService *GitService, containerService *ContainerService) *DeployService {
	return &DeployService{
		deployHistoryService: deployHistoryService,
		projectService:       projectService,
		gitService:           gitService,
		containerService:     containerService,
	}
}

func (s DeployService) Deploy(body *model.DeployHistory) error {

	project := &model.Project{Model: model.Model{ID: body.Origin_ID}}
	if _, err := s.projectService.Get(project); err != nil {
		return err
	}
	now := time.Now()
	body.Image = fmt.Sprintf("%s:%s-%s", project.Name, project.Branch, now.Format("20060102.1504"))
	body.Status = model.DEPLOY_STATUS_RUNING
	if err := s.deployHistoryService.CreateDeploy(body); err != nil {
		return err
	}

	needPull := true
	if project.Type == "project" {
		// 1.pull code
		s.deployHistoryService.DeployStepInto(body)
		if err := s.gitService.Pull(project.LocalRepo); err != nil {
			s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(now))
			return err
		}

		//2.build image
		s.deployHistoryService.DeployStepInto(body)
		if err := s.Build(project.LocalRepo, body.Image); err != nil {
			s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(now))
			return err
		}
		needPull = false
	}

	//3.run container
	s.deployHistoryService.DeployStepInto(body)
	if err := s.Run(project, needPull); err != nil {
		s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(now))
		return err
	}
	s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_SUCCESS, time.Since(now))

	project.AccessUrl = fmt.Sprintf("http://121.196.214.245:%d", project.BindPort)
	//project.AccessUrl = fmt.Sprintf("http://%s-%s.%s:%d",project.Branch,project.UserName,viper.GetString("baseUrl"),project.BindPort)
	return s.projectService.Update(project)
}

func (s DeployService) Clone(project *model.Project) error {
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

func (s DeployService) Build(path string, image string) error {
	if err := s.containerService.BuildImage(model.BuildImageOption{
		Path: path,
		Tag:  image,
	}); err != nil {
		return err
	}
	return nil
}

func (s DeployService) Run(req *model.Project, needPull bool) error {
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
		NeedPull:  needPull,
	}); err != nil {
		return err
	} else {
		req.Container = id
	}

	if err := s.projectService.Update(req); err != nil {
		return err
	}
	return nil
}
