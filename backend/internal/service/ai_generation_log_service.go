package service

import (
	"encoding/json"
	"strings"
	"time"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type AIGenerationLogService struct {
	repo *repository.AIGenerationLogRepo
}

type AIGenerationLogPayload struct {
	UserID         uint
	GenerationType string
	Scene          string
	Model          string
	Input          interface{}
	Output         interface{}
	Status         string
	ErrorMessage   string
	Duration       time.Duration
	RecipeIDs      []uint
}

func NewAIGenerationLogService(repo *repository.AIGenerationLogRepo) *AIGenerationLogService {
	return &AIGenerationLogService{repo: repo}
}

func (s *AIGenerationLogService) Record(payload AIGenerationLogPayload) {
	if s == nil || s.repo == nil {
		return
	}

	status := strings.TrimSpace(payload.Status)
	if status == "" {
		status = "failed"
	}

	log := &model.AIGenerationLog{
		UserID:         payload.UserID,
		GenerationType: strings.TrimSpace(payload.GenerationType),
		Scene:          strings.TrimSpace(payload.Scene),
		Model:          strings.TrimSpace(payload.Model),
		InputJSON:      marshalLogJSON(payload.Input),
		OutputJSON:     marshalLogJSON(payload.Output),
		Status:         status,
		ErrorMessage:   strings.TrimSpace(payload.ErrorMessage),
		DurationMS:     payload.Duration.Milliseconds(),
		RecipeIDsJSON:  marshalLogJSON(payload.RecipeIDs),
	}
	_ = s.repo.Create(log)
}

func marshalLogJSON(value interface{}) model.JSON {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return model.JSON(raw)
}

func aiModelName(client *AIClient) string {
	if client == nil {
		return ""
	}
	return client.Model()
}

func aiErrorText(err error) string {
	if err == nil {
		return ""
	}
	text := strings.TrimSpace(err.Error())
	runes := []rune(text)
	if len(runes) > 1000 {
		return string(runes[:1000])
	}
	return text
}

func summarizeAIRecipeGenerateResult(result *AIRecipeGenerateResult) map[string]interface{} {
	if result == nil || result.Recipe == nil {
		return nil
	}
	return map[string]interface{}{
		"recipe_id": result.Recipe.ID,
		"title":     result.Recipe.Title,
		"created":   result.Created,
	}
}

func recipeIDsFromAIRecipeGenerateResult(result *AIRecipeGenerateResult) []uint {
	if result == nil || result.Recipe == nil || result.Recipe.ID == 0 {
		return nil
	}
	return []uint{result.Recipe.ID}
}

func summarizeRecommendResult(result *RecommendResult) map[string]interface{} {
	if result == nil {
		return nil
	}
	dishes := make([]map[string]interface{}, 0, len(result.Dishes))
	for _, dish := range result.Dishes {
		dishes = append(dishes, map[string]interface{}{
			"recipe_id": dish.RecipeID,
			"name":      dish.Name,
			"type":      dish.Type,
		})
	}
	return map[string]interface{}{
		"menu_name": result.MenuName,
		"source":    result.Source,
		"dishes":    dishes,
	}
}

func recipeIDsFromDishes(dishes []DishResult) []uint {
	ids := make([]uint, 0, len(dishes))
	seen := make(map[uint]bool, len(dishes))
	for _, dish := range dishes {
		if dish.RecipeID == 0 || seen[dish.RecipeID] {
			continue
		}
		seen[dish.RecipeID] = true
		ids = append(ids, dish.RecipeID)
	}
	return ids
}

func summarizeShoppingAIResult(result *DishShoppingListResult) map[string]interface{} {
	if result == nil {
		return nil
	}
	items := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		if strings.TrimSpace(item.Name) != "" {
			items = append(items, item.Name)
		}
	}
	output := map[string]interface{}{
		"items": items,
	}
	if result.Recipe != nil {
		output["recipe_id"] = result.Recipe.ID
		output["recipe_title"] = result.Recipe.Title
	}
	if result.List != nil {
		output["shopping_list_id"] = result.List.ID
		output["shopping_list_name"] = result.List.Name
	}
	return output
}

func recipeIDsFromShoppingAIResult(result *DishShoppingListResult) []uint {
	if result == nil || result.Recipe == nil || result.Recipe.ID == 0 {
		return nil
	}
	return []uint{result.Recipe.ID}
}
