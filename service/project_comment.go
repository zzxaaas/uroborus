package service

import (
	"uroborus/model"
	"uroborus/store"
)

type ProjectCommentService struct {
	commentStore *store.ProjectCommentStore
}

func NewProjectCommentService(commentStore *store.ProjectCommentStore) *ProjectCommentService {
	return &ProjectCommentService{
		commentStore: commentStore,
	}
}

func (s ProjectCommentService) Register(req *model.ProjectComment) error {
	return s.commentStore.Save(req)
}

func (s ProjectCommentService) Find(req *model.ProjectComment) ([]model.ProjectComment, error) {
	return s.commentStore.Find(req)
}
