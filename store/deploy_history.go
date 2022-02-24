package store

import (
	"uroborus/model"
)

// DeployHistoryStore -
type DeployHistoryStore struct {
	db *DB
}

// NewDeployHistoryStore -
func NewDeployHistoryStore(db *DB) *DeployHistoryStore {
	return &DeployHistoryStore{
		db: db,
	}
}

func (s *DeployHistoryStore) Save(body *model.DeployHistory) error {
	return s.db.Model(body).Create(body).Error
}

func (s *DeployHistoryStore) Update(body *model.DeployHistory) error {
	return s.db.Model(body).Where("id=?", body.ID).Updates(body).Error
}

func (s *DeployHistoryStore) Delete(body *model.DeployHistory) error {
	return s.db.Delete(body, body).Error
}

func (s *DeployHistoryStore) Get(body *model.DeployHistory) error {
	return s.db.First(&body, body).Error
}

func (s *DeployHistoryStore) Find(body *model.DeployHistory) ([]model.DeployHistory, error) {
	ans := make([]model.DeployHistory, 0)
	err := s.db.Order("id desc").Find(&ans, body).Error
	return ans, err
}
