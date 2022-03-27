package store

import (
	"uroborus/model"
)

// GroupStore -
type GroupStore struct {
	db *DB
}

// NewGroupStore -
func NewGroupStore(db *DB) *GroupStore {
	return &GroupStore{
		db: db,
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

func (s GroupStore) Find(body *model.Group) ([]model.Group, error) {
	ans := make([]model.Group, 0)
	err := s.db.Find(&ans, body).Error
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func (s GroupStore) Delete(body *model.Group) error {
	return s.db.Delete(body, body).Error
}
