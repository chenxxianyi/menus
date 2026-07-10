package service

import (
	"context"
	"strings"
	"time"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type RecipeService struct {
	recipeRepo   *repository.RecipeRepo
	categoryRepo *repository.CategoryRepo
	favRepo      *repository.FavoriteRepo
	historyRepo  *repository.BrowseHistoryRepo
	feedbackRepo *repository.UserRecipeFeedbackRepo
	aiLogService *AIGenerationLogService
	aiClient     *AIClient
}

func NewRecipeService(recipeRepo *repository.RecipeRepo, categoryRepo *repository.CategoryRepo, favRepo *repository.FavoriteRepo, historyRepo *repository.BrowseHistoryRepo) *RecipeService {
	return &RecipeService{recipeRepo: recipeRepo, categoryRepo: categoryRepo, favRepo: favRepo, historyRepo: historyRepo}
}

func (s *RecipeService) SetFeedbackRepo(feedbackRepo *repository.UserRecipeFeedbackRepo) {
	s.feedbackRepo = feedbackRepo
}

func (s *RecipeService) SetAIClient(aiClient *AIClient) {
	s.aiClient = aiClient
}

func (s *RecipeService) SetAIGenerationLogService(aiLogService *AIGenerationLogService) {
	s.aiLogService = aiLogService
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
		if s.feedbackRepo != nil {
			if items, err := s.feedbackRepo.FindByUserAndRecipe(userID, id); err == nil {
				recipe.Feedback = FeedbackStatusMap(items)
			}
		}
	}

	return recipe, nil
}

func (s *RecipeService) GetHotRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindHot(limit)
}

func (s *RecipeService) GetRandomRecipes(limit int) ([]model.Recipe, error) {
	return s.recipeRepo.FindRandom(limit)
}

// GetMealRecommendations keeps the homepage meal selector meaningful even
// before dedicated breakfast/lunch/dinner tags are maintained in the catalog.
// It applies conservative cooking-time preferences and still falls back to a
// random published recipe when the catalog is small.
func (s *RecipeService) GetMealRecommendations(mealType string, limit int) ([]model.Recipe, error) {
	maxCookTime := 0
	switch strings.TrimSpace(mealType) {
	case "早餐", "breakfast":
		maxCookTime = 30
	case "夜宵", "late_night":
		maxCookTime = 20
	case "午餐", "lunch", "晚餐", "dinner":
		maxCookTime = 60
	}
	if maxCookTime > 0 {
		if recipes, err := s.recipeRepo.FindRandomByMaxCookTime(limit, maxCookTime); err == nil && len(recipes) > 0 {
			return recipes, nil
		}
	}
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

func (s *RecipeService) GenerateRecipeByAI(ctx context.Context, userID uint, dishName string) (*AIRecipeGenerateResult, error) {
	return s.generateRecipeByAI(ctx, userID, dishName, true)
}

func (s *RecipeService) generateRecipeByAI(ctx context.Context, userID uint, dishName string, logEnabled bool) (result *AIRecipeGenerateResult, err error) {
	start := time.Now()
	name := strings.TrimSpace(dishName)
	if logEnabled {
		defer func() {
			status := "success"
			if err != nil {
				status = "failed"
			}
			s.aiLogService.Record(AIGenerationLogPayload{
				UserID:         userID,
				GenerationType: "recipe",
				Model:          aiModelName(s.aiClient),
				Input:          map[string]interface{}{"dish_name": name},
				Output:         summarizeAIRecipeGenerateResult(result),
				Status:         status,
				ErrorMessage:   aiErrorText(err),
				Duration:       time.Since(start),
				RecipeIDs:      recipeIDsFromAIRecipeGenerateResult(result),
			})
		}()
	}
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
