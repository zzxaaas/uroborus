package store

import "uroborus/model"

// ProjectCommentStore -
type ProjectCommentStore struct {
	db *DB
}

// NewProjectCommentStore -
func NewProjectCommentStore(db *DB) *ProjectCommentStore {
	return &ProjectCommentStore{
		db: db,
	}
}

func (s ProjectCommentStore) Save(body *model.ProjectComment) error {
	return s.db.Model(body).Create(body).Error
}

func (s ProjectCommentStore) Find(body *model.ProjectComment) ([]model.ProjectComment, error) {
	ans := make([]model.ProjectComment, 0)
	err := s.db.Where("origin_id=?", body.OriginId).Order("created_at").Find(&ans).Error
	return ans, err
}
