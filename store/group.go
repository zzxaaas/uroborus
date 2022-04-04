package store

import (
	"uroborus/model"
)

// GroupStore -
type GroupStore struct {
	db        *DB
	projStore *ProjectStore
}

// NewGroupStore -
func NewGroupStore(db *DB, projStore *ProjectStore) *GroupStore {
	return &GroupStore{
		db:        db,
		projStore: projStore,
	}
}

func (s GroupStore) Save(body *model.Group) error {
	return s.db.Model(body).Create(body).Error
}

func (s GroupStore) Update(body *model.Group) error {
	baseSql := s.db.Model(body)
	if body.Name != "" {
		baseSql = baseSql.Where("name=?", body.Name)
	}
	if body.ID != 0 {
		baseSql = baseSql.Where("id=?", body.ID)
	}
	return baseSql.Updates(body).Error
}

func (s GroupStore) FindCreateGroup(body *model.Group) ([]model.Group, error) {
	ans := make([]model.Group, 0)
	err := s.db.Find(&ans, body).Error
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func (s GroupStore) FindJoinGroup(body *model.Group) ([]model.Group, error) {
	ans := make([]model.Group, 0)
	groups := make([]uint, 0)
	err := s.db.Model(model.Project{}).
		Select("group_id").Where("user_name=?", body.CreateUser).
		Group("group_id").Find(&groups).Error
	if err != nil {
		return nil, err
	}
	err = s.db.Where("id in (?)", groups).Where("create_user != ?", body.CreateUser).Find(&ans).Error
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func (s GroupStore) Delete(body *model.Group) error {
	return s.db.Delete(body, body).Error
}
