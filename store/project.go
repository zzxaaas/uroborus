package store

import (
	"gorm.io/gorm/clause"
	"uroborus/model"
)

// ProjectStore -
type ProjectStore struct {
	db *DB
}

// NewProjectStore -
func NewProjectStore(db *DB) *ProjectStore {
	return &ProjectStore{
		db: db,
	}
}

func (s *ProjectStore) Save(body *model.Project) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"branch", "env", "container", "command", "status", "updated_at", "deleted_at"}),
	}).Create(&body).Error
}

func (s ProjectStore) Update(body *model.Project) error {
	baseSql := s.db.Model(body)
	if body.Name != "" {
		baseSql = baseSql.Where("name=?", body.Name)
	}
	if body.ID != 0 {
		baseSql = baseSql.Where("id=?", body.ID)
	}
	return baseSql.Updates(body).Error
}

func (s ProjectStore) Delete(body *model.Project) error {
	return s.db.Delete(body, body).Error
}

func (s *ProjectStore) Get(project *model.Project) error {
	return s.db.First(&project, project).Error
}

func (s *ProjectStore) Find(project model.Project) ([]model.Project, error) {
	ans := make([]model.Project, 0)
	err := s.db.Find(&ans, project).Error
	return ans, err
}
