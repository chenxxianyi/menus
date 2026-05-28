package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type FavoriteService struct {
	favRepo    *repository.FavoriteRepo
	recipeRepo *repository.RecipeRepo
}

func NewFavoriteService(favRepo *repository.FavoriteRepo, recipeRepo *repository.RecipeRepo) *FavoriteService {
	return &FavoriteService{favRepo: favRepo, recipeRepo: recipeRepo}
}

func (s *FavoriteService) AddFavorite(userID, recipeID uint) error {
	if s.favRepo.Exists(userID, recipeID) {
		return nil
	}
	if err := s.favRepo.Add(userID, recipeID); err != nil {
		return err
	}
	return s.recipeRepo.IncrementFavoriteCount(recipeID)
}

func (s *FavoriteService) RemoveFavorite(userID, recipeID uint) error {
	if !s.favRepo.Exists(userID, recipeID) {
		return nil
	}
	if err := s.favRepo.Remove(userID, recipeID); err != nil {
		return err
	}
	return s.recipeRepo.DecrementFavoriteCount(recipeID)
}

func (s *FavoriteService) GetUserFavorites(userID uint, page, pageSize int) ([]model.Recipe, int64, error) {
	favs, total, err := s.favRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	recipeIDs := make([]uint, len(favs))
	for i, f := range favs {
		recipeIDs[i] = f.RecipeID
	}

	if len(recipeIDs) == 0 {
		return []model.Recipe{}, 0, nil
	}

	recipes, err := s.recipeRepo.FindByIDs(recipeIDs)
	return recipes, total, err
}
