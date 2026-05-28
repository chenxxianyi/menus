package service

import (
	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type RecipeService struct {
	recipeRepo   *repository.RecipeRepo
	categoryRepo *repository.CategoryRepo
	favRepo      *repository.FavoriteRepo
}

func NewRecipeService(recipeRepo *repository.RecipeRepo, categoryRepo *repository.CategoryRepo, favRepo *repository.FavoriteRepo) *RecipeService {
	return &RecipeService{recipeRepo: recipeRepo, categoryRepo: categoryRepo, favRepo: favRepo}
}

func (s *RecipeService) ListRecipes(keyword string, categoryID uint, taste, cookTime, difficulty, healthTags string, page, pageSize int) ([]model.Recipe, int64, error) {
	return s.recipeRepo.List(keyword, categoryID, taste, cookTime, difficulty, healthTags, page, pageSize)
}

func (s *RecipeService) GetRecipeDetail(id uint, userID uint) (*model.Recipe, error) {
	recipe, err := s.recipeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.recipeRepo.IncrementViewCount(id)

	if cat, err := s.categoryRepo.FindByID(recipe.CategoryID); err == nil {
		recipe.CategoryName = cat.Name
	}

	if userID > 0 {
		recipe.IsFavorited = s.favRepo.Exists(userID, id)
	}

	return recipe, nil
}

func (s *RecipeService) GetHotRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindHot(limit)
}

func (s *RecipeService) GetRandomRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindRandom(limit)
}

func (s *RecipeService) CreateRecipe(recipe *model.Recipe) error {
	return s.recipeRepo.Create(recipe)
}

func (s *RecipeService) UpdateRecipe(recipe *model.Recipe) error {
	return s.recipeRepo.Update(recipe)
}

func (s *RecipeService) DeleteRecipe(id uint) error {
	return s.recipeRepo.Delete(id)
}

func (s *RecipeService) ListAllRecipes(keyword string, categoryID uint, status *int8, page, pageSize int) ([]model.Recipe, int64, error) {
	return s.recipeRepo.List(keyword, categoryID, "", "", "", "", page, pageSize)
}
