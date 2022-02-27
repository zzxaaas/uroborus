package service

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
	"uroborus/common/kafka"
	"uroborus/model"
)

type DeployService struct {
	deployHistoryService *DeployHistoryService
	projectService       *ProjectService
	gitService           *GitService
	containerService     *ContainerService
	kafkaCli             *kafka.Client
}

func NewDeployService(deployHistoryService *DeployHistoryService, projectService *ProjectService, gitService *GitService, containerService *ContainerService, kafkaCli *kafka.Client) *DeployService {
	return &DeployService{
		deployHistoryService: deployHistoryService,
		projectService:       projectService,
		gitService:           gitService,
		containerService:     containerService,
		kafkaCli:             kafkaCli,
	}
}

func (s DeployService) Deploy(body *model.DeployHistory) error {
	project := &model.Project{Model: model.Model{ID: body.Origin_ID}}
	if _, err := s.projectService.Get(project); err != nil {
		return err
	}
	body.CreatedAt = time.Now()
	body.Image = fmt.Sprintf("%s:%s-%s", project.Name, project.Branch, body.CreatedAt.Format("20060102.1504"))
	body.Status = model.DEPLOY_STATUS_RUNING
	if err := s.deployHistoryService.CreateDeploy(body); err != nil {
		return err
	}
	if err := s.kafkaCli.CreateTopic(strconv.Itoa(int(body.ID))); err != nil {
		return err
	}
	go s.doDeploy(project, body)
	return nil
}

func (s DeployService) doDeploy(project *model.Project, body *model.DeployHistory) {
	needPull := true
	if project.Type == "project" {
		// 1.pull code
		if err := s.Pull(project, body); err != nil {
			logrus.Error(err.Error())
			return
		}
		//2.build image
		if err := s.Build(project, body); err != nil {
			logrus.Error(err.Error())
			return
		}
		needPull = false
	}
	//3.run container
	if err := s.Run(project, body, needPull); err != nil {
		logrus.Error(err.Error())
		return
	}
	//4.success
	s.deployHistoryService.DeployStepInto(body)
	s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_SUCCESS, time.Since(body.CreatedAt))
}

func (s DeployService) Pull(project *model.Project, body *model.DeployHistory) error {
	topic := strconv.Itoa(int(body.ID))
	s.deployHistoryService.DeployStepInto(body)
	err := s.gitService.Pull(project.LocalRepo)
	if err == nil {
		s.kafkaCli.SendLog(kafka.PackMsg(topic, "git pull sucess", int32(body.Step)))
	} else {
		s.kafkaCli.SendLog(kafka.PackMsg(topic, err.Error(), int32(body.Step)))
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(body.CreatedAt))
			return err
		}
	}
	return nil
}

func (s DeployService) Build(project *model.Project, body *model.DeployHistory) error {
	s.deployHistoryService.DeployStepInto(body)
	if err := s.doBuild(project.LocalRepo, body.Image, body.ID); err != nil {
		s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(body.CreatedAt))
		return err
	}
	return nil
}

func (s DeployService) Run(project *model.Project, body *model.DeployHistory, needPull bool) error {
	s.deployHistoryService.DeployStepInto(body)
	if err := s.doRun(project, needPull, body.ID); err != nil {
		s.deployHistoryService.UpdateStatus(body.ID, model.DEPLOY_STATUS_FAILED, time.Since(body.CreatedAt))
		return err
	}
	return nil
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

func (s DeployService) doBuild(path string, image string, deployId uint) error {
	if err := s.containerService.BuildImage(model.BuildImageOption{
		Path:     path,
		Tag:      image,
		DeployID: deployId,
	}); err != nil {
		return err
	}
	return nil
}

func (s DeployService) doRun(req *model.Project, needPull bool, deployID uint) error {
	if req.Container != "" {
		if err := s.containerService.RemoveContainer(req.Container); err != nil {
			return err
		}
		s.kafkaCli.SendLog(
			kafka.PackMsg(strconv.Itoa(int(deployID)),
				fmt.Sprintf("remove old container {%s} success", req.Container),
				model.DEPLOY_STEP_RUN))
	}

	if id, err := s.containerService.StartContainerWithOption(model.ContainerOption{
		Name:      req.Name,
		Image:     req.Image,
		ProtoPort: req.Port,
		Port:      strconv.Itoa(req.BindPort),
		Env:       strings.Split(req.Env, ","),
		NeedPull:  needPull,
		DeployID:  deployID,
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
