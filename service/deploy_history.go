package service

import (
	"time"
	"uroborus/common/kafka"
	"uroborus/model"
	"uroborus/store"
)

type DeployHistoryService struct {
	deployHistoryStore *store.DeployHistoryStore
	kafkaCli           *kafka.Client
}

func NewDeployHistoryService(deployHistoryStore *store.DeployHistoryStore, kafkaCli *kafka.Client) *DeployHistoryService {
	return &DeployHistoryService{
		deployHistoryStore: deployHistoryStore,
		kafkaCli:           kafkaCli,
	}
}

func (s DeployHistoryService) CreateDeploy(body *model.DeployHistory) error {
	return s.deployHistoryStore.Save(body)
}

func (s DeployHistoryService) Get(body *model.DeployHistory) error {
	return s.deployHistoryStore.Get(body)
}

func (s DeployHistoryService) Find(body *model.DeployHistory) ([]model.DeployHistory, error) {
	return s.deployHistoryStore.Find(body)
}

func (s DeployHistoryService) UpdateStatus(id uint, status int, duration time.Duration) {
	s.deployHistoryStore.Update(&model.DeployHistory{Model: model.Model{ID: id}, Status: status, Duration: duration})
}

func (s DeployHistoryService) Update(body *model.DeployHistory, status int) {
	body.Status = status
	s.deployHistoryStore.Update(body)
}

func (s DeployHistoryService) DeployStepInto(body *model.DeployHistory) {
	if body.Step < model.DEPLOY_STEP_END {
		body.Step += 1
	}
	s.deployHistoryStore.Update(body)
}
