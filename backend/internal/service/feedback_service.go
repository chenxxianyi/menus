package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type FeedbackService struct {
	repo *repository.FeedbackRepo
}

func NewFeedbackService(repo *repository.FeedbackRepo) *FeedbackService {
	return &FeedbackService{repo: repo}
}

func (s *FeedbackService) Create(fb *model.Feedback) error {
	return s.repo.Create(fb)
}

func (s *FeedbackService) List(status *int8, page, pageSize int) ([]model.Feedback, int64, error) {
	return s.repo.List(status, page, pageSize)
}

func (s *FeedbackService) UpdateStatus(id uint, status int8, reply string) error {
	return s.repo.UpdateStatus(id, status, reply)
}
