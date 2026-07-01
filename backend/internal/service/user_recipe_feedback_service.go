package service

import (
	"errors"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

var ErrInvalidRecipeFeedbackType = errors.New("invalid recipe feedback type")

var allowedRecipeFeedbackTypes = map[string]bool{
	"cooked":  true,
	"like":    true,
	"dislike": true,
	"block":   true,
}

type UserRecipeFeedbackService struct {
	repo *repository.UserRecipeFeedbackRepo
}

func NewUserRecipeFeedbackService(repo *repository.UserRecipeFeedbackRepo) *UserRecipeFeedbackService {
	return &UserRecipeFeedbackService{repo: repo}
}

func (s *UserRecipeFeedbackService) Set(userID, recipeID uint, feedbackType, source string) error {
	feedbackType = strings.TrimSpace(feedbackType)
	if !allowedRecipeFeedbackTypes[feedbackType] {
		return ErrInvalidRecipeFeedbackType
	}
	source = strings.TrimSpace(source)
	if source == "" {
		source = "detail"
	}
	return s.repo.Upsert(&model.UserRecipeFeedback{
		UserID:       userID,
		RecipeID:     recipeID,
		FeedbackType: feedbackType,
		Source:       source,
	})
}

func (s *UserRecipeFeedbackService) Delete(userID, recipeID uint, feedbackType string) error {
	feedbackType = strings.TrimSpace(feedbackType)
	if !allowedRecipeFeedbackTypes[feedbackType] {
		return ErrInvalidRecipeFeedbackType
	}
	return s.repo.Delete(userID, recipeID, feedbackType)
}

func (s *UserRecipeFeedbackService) Status(userID, recipeID uint) (map[string]bool, error) {
	items, err := s.repo.FindByUserAndRecipe(userID, recipeID)
	if err != nil {
		return nil, err
	}
	return FeedbackStatusMap(items), nil
}

func (s *UserRecipeFeedbackService) List(userID uint) ([]model.UserRecipeFeedback, error) {
	return s.repo.FindByUser(userID)
}

func FeedbackStatusMap(items []model.UserRecipeFeedback) map[string]bool {
	status := map[string]bool{
		"cooked":  false,
		"like":    false,
		"dislike": false,
		"block":   false,
	}
	for _, item := range items {
		if _, ok := status[item.FeedbackType]; ok {
			status[item.FeedbackType] = true
		}
	}
	return status
}
