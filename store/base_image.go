package store

import (
	"fmt"
	"gorm.io/gorm/clause"
	"uroborus/model"
)

// BaseImageStore -
type BaseImageStore struct {
	db *DB
}

// NewProjectStore -
func NewBaseImageStore(db *DB) *BaseImageStore {
	return &BaseImageStore{
		db: db,
	}
}

func (s BaseImageStore) Save(body model.BaseImage) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"tags", "port", "updated_at"}),
	}).Create(&body).Error
}

func (s BaseImageStore) Find(body model.BaseImage) ([]model.BaseImage, error) {
	ans := make([]model.BaseImage, 0)
	baseSQL := s.db.Model(body)
	if body.Name != "" {
		name := fmt.Sprintf("%s", "%"+body.Name+"%")
		baseSQL = baseSQL.Where("name like ?", name)
		body.Name = ""
	}
	err := baseSQL.Find(&ans, body).Error
	return ans, err
}
