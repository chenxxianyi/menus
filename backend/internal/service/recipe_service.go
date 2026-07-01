package service

import (
	"context"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type RecipeService struct {
	recipeRepo   *repository.RecipeRepo
	categoryRepo *repository.CategoryRepo
	favRepo      *repository.FavoriteRepo
	historyRepo  *repository.BrowseHistoryRepo
	aiClient     *AIClient
}

func NewRecipeService(recipeRepo *repository.RecipeRepo, categoryRepo *repository.CategoryRepo, favRepo *repository.FavoriteRepo, historyRepo *repository.BrowseHistoryRepo) *RecipeService {
	return &RecipeService{recipeRepo: recipeRepo, categoryRepo: categoryRepo, favRepo: favRepo, historyRepo: historyRepo}
}

func (s *RecipeService) SetAIClient(aiClient *AIClient) {
	s.aiClient = aiClient
}

func (s *RecipeService) ListRecipes(keyword string, categoryID uint, taste, cookTime, difficulty, healthTags, sortBy string, page, pageSize int) ([]model.Recipe, int64, error) {
	if sortBy != "hot" {
		sortBy = "latest"
	}
	return s.recipeRepo.List(keyword, categoryID, taste, cookTime, difficulty, healthTags, sortBy, page, pageSize)
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
		_ = s.historyRepo.Record(userID, id)
	}

	return recipe, nil
}

func (s *RecipeService) GetHotRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindHot(limit)
}

func (s *RecipeService) GetRandomRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindRandom(limit)
}

func (s *RecipeService) GetFilterOptions() (map[string]interface{}, error) {
	tastes, err := s.recipeRepo.DistinctTastes()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"tastes": tastes,
	}, nil
}

func (s *RecipeService) CreateRecipe(recipe *model.Recipe) error {
	return s.recipeRepo.Create(recipe)
}

type AIRecipeGenerateResult struct {
	Recipe  *model.Recipe `json:"recipe"`
	Created bool          `json:"created"`
}

func (s *RecipeService) GenerateRecipeByAI(ctx context.Context, dishName string) (*AIRecipeGenerateResult, error) {
	name := strings.TrimSpace(dishName)
	if name == "" {
		return nil, ErrAIInvalidResponse
	}
	if existing, err := s.recipeRepo.FindByTitle(name); err == nil {
		return &AIRecipeGenerateResult{Recipe: existing, Created: false}, nil
	}
	if s.aiClient == nil || !s.aiClient.IsConfigured() {
		return nil, ErrAIConfigMissing
	}

	draft, err := s.aiClient.GenerateRecipeDraft(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing, err := s.recipeRepo.FindByTitle(draft.Title); err == nil {
		return &AIRecipeGenerateResult{Recipe: existing, Created: false}, nil
	}

	recipe, err := draft.ToRecipe(s.defaultRecipeCategoryID())
	if err != nil {
		return nil, err
	}
	if err := s.recipeRepo.Create(recipe); err != nil {
		return nil, err
	}
	return &AIRecipeGenerateResult{Recipe: recipe, Created: true}, nil
}

func (s *RecipeService) CreateOrReuseAIRecipeDraft(draft *AIRecipeDraft) (*AIRecipeGenerateResult, error) {
	if draft == nil {
		return nil, ErrAIInvalidResponse
	}
	name := strings.TrimSpace(draft.Title)
	if name == "" {
		return nil, ErrAIInvalidResponse
	}
	if existing, err := s.recipeRepo.FindByTitle(name); err == nil {
		return &AIRecipeGenerateResult{Recipe: existing, Created: false}, nil
	}

	recipe, err := draft.ToRecipe(s.defaultRecipeCategoryID())
	if err != nil {
		return nil, err
	}
	if err := s.recipeRepo.Create(recipe); err != nil {
		return nil, err
	}
	return &AIRecipeGenerateResult{Recipe: recipe, Created: true}, nil
}

func (s *RecipeService) defaultRecipeCategoryID() uint {
	categories, err := s.categoryRepo.FindAll()
	if err != nil || len(categories) == 0 {
		return 0
	}
	return categories[0].ID
}

func (s *RecipeService) UpdateRecipe(recipe *model.Recipe) error {
	return s.recipeRepo.Update(recipe)
}

func (s *RecipeService) DeleteRecipe(id uint) error {
	return s.recipeRepo.Delete(id)
}

func (s *RecipeService) ListAllRecipes(keyword string, categoryID uint, status *int8, page, pageSize int) ([]model.Recipe, int64, error) {
	return s.recipeRepo.List(keyword, categoryID, "", "", "", "", "latest", page, pageSize)
}
