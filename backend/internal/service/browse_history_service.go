package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type BrowseHistoryService struct {
	historyRepo *repository.BrowseHistoryRepo
}

func NewBrowseHistoryService(historyRepo *repository.BrowseHistoryRepo) *BrowseHistoryService {
	return &BrowseHistoryService{historyRepo: historyRepo}
}

func (s *BrowseHistoryService) Record(userID, recipeID uint) error {
	return s.historyRepo.Record(userID, recipeID)
}

func (s *BrowseHistoryService) List(userID uint, page, pageSize int) ([]model.BrowseHistory, int64, error) {
	return s.historyRepo.List(userID, page, pageSize)
}

func (s *BrowseHistoryService) Clear(userID uint) error {
	return s.historyRepo.Clear(userID)
}
