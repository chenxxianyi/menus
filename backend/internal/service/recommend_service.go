package service

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type RecommendService struct {
	recipeRepo     *repository.RecipeRepo
	ingredientRepo *repository.IngredientRepo
	prefRepo       *repository.UserPrefRepo
	recipeService  *RecipeService
	aiClient       *AIClient
}

func NewRecommendService(recipeRepo *repository.RecipeRepo, ingredientRepo *repository.IngredientRepo, prefRepo *repository.UserPrefRepo, recipeService *RecipeService, aiClient *AIClient) *RecommendService {
	return &RecommendService{
		recipeRepo:     recipeRepo,
		ingredientRepo: ingredientRepo,
		prefRepo:       prefRepo,
		recipeService:  recipeService,
		aiClient:       aiClient,
	}
}

type RecommendParams struct {
	Scene               string   `json:"scene"`
	PeopleCount         int      `json:"people_count"`
	MealType            string   `json:"meal_type"`
	TastePreference     []string `json:"taste_preference"`
	HealthGoal          string   `json:"health_goal"`
	AvoidIngredients    []string `json:"avoid_ingredients"`
	ExistingIngredients []string `json:"existing_ingredients"`
	CookTimePreference  string   `json:"cook_time_preference"`
	ExcludeRecipeIDs    []uint   `json:"-"`
}

type DishResult struct {
	RecipeID     uint     `json:"recipe_id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	CookTime     int      `json:"cook_time"`
	Difficulty   string   `json:"difficulty"`
	Ingredients  []string `json:"ingredients"`
	StepsSummary string   `json:"steps_summary"`
}

type RecommendResult struct {
	MenuName     string       `json:"menu_name"`
	Reason       string       `json:"reason"`
	Dishes       []DishResult `json:"dishes"`
	ShoppingList []string     `json:"shopping_list"`
	Source       string       `json:"source,omitempty"`
}

type scoredRecipe struct {
	recipe *model.Recipe
	score  float64
}

type IngredientRecommendResult struct {
	Recipe             model.Recipe `json:"recipe"`
	MatchRate          float64      `json:"match_rate"`
	MatchedIngredients []string     `json:"matched_ingredients"`
	MissingIngredients []string     `json:"missing_ingredients"`
	Reason             string       `json:"reason"`
}

type IngredientRecommendPayload struct {
	List  []IngredientRecommendResult `json:"list"`
	Total int                         `json:"total"`
}

var ingredientAliases = map[string]string{
	"番茄":  "西红柿",
	"西红柿": "西红柿",
	"土豆":  "土豆",
	"马铃薯": "土豆",
	"洋芋":  "土豆",
	"鸡蛋":  "鸡蛋",
	"蛋":   "鸡蛋",
}

func normalizeIngredientName(name string) string {
	value := strings.TrimSpace(name)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "　", "")
	if alias, ok := ingredientAliases[value]; ok {
		return alias
	}
	return value
}

func normalizeIngredientList(items []string, limit int) []string {
	if limit <= 0 {
		limit = 20
	}
	set := make(map[string]bool)
	result := make([]string, 0, len(items))
	for _, item := range items {
		name := normalizeIngredientName(item)
		if name == "" || set[name] {
			continue
		}
		set[name] = true
		result = append(result, name)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func recipeIngredientNames(raw model.JSON) []string {
	if len(raw) == 0 {
		return nil
	}

	var items []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil
	}

	names := make([]string, 0, len(items))
	for _, item := range items {
		if item.Name != "" {
			names = append(names, item.Name)
		}
	}
	return names
}

func containsRecipeID(ids []uint, id uint) bool {
	for _, item := range ids {
		if item == id {
			return true
		}
	}
	return false
}

func recipeHasAvoidIngredient(recipeIngredients []string, avoidIngredients []string) bool {
	if len(avoidIngredients) == 0 {
		return false
	}
	avoidSet := make(map[string]bool)
	for _, item := range normalizeIngredientList(avoidIngredients, 100) {
		avoidSet[item] = true
	}
	for _, item := range normalizeIngredientList(recipeIngredients, 100) {
		if avoidSet[item] {
			return true
		}
	}
	return false
}

func jsonStringList(raw model.JSON) []string {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil
	}
	return normalizeStringList(values, 20)
}

func sceneLabel(scene string) string {
	switch scene {
	case "quick_meal":
		return "快手一餐"
	case "family_dinner":
		return "家庭聚餐"
	case "fat_loss":
		return "减脂轻食"
	case "treat_guest":
		return "宴客招待"
	case "late_night":
		return "夜宵"
	default:
		return "日常用餐"
	}
}

func mergeStringList(primary []string, secondary []string, limit int) []string {
	merged := make([]string, 0, len(primary)+len(secondary))
	merged = append(merged, primary...)
	merged = append(merged, secondary...)
	return normalizeStringList(merged, limit)
}

func ensureSceneDefaults(params *RecommendParams, pref *model.UserPreference) AISceneRecommendContext {
	ctx := AISceneRecommendContext{
		Scene:              strings.TrimSpace(params.Scene),
		SceneLabel:         sceneLabel(params.Scene),
		MealType:           strings.TrimSpace(params.MealType),
		PeopleCount:        params.PeopleCount,
		TastePreference:    normalizeStringList(params.TastePreference, 12),
		HealthGoal:         strings.TrimSpace(params.HealthGoal),
		AvoidIngredients:   normalizeStringList(params.AvoidIngredients, 20),
		CookTimePreference: strings.TrimSpace(params.CookTimePreference),
	}
	if ctx.MealType == "" {
		ctx.MealType = "dinner"
	}
	if ctx.PeopleCount <= 0 {
		ctx.PeopleCount = 2
	}
	if pref != nil {
		ctx.TastePreference = mergeStringList(ctx.TastePreference, jsonStringList(pref.TastePreference), 12)
		ctx.AvoidIngredients = mergeStringList(ctx.AvoidIngredients, jsonStringList(pref.AvoidIngredients), 20)
		ctx.FavoriteIngredients = jsonStringList(pref.FavoriteIngredients)
		if ctx.HealthGoal == "" {
			ctx.HealthGoal = strings.TrimSpace(pref.HealthGoal)
		}
		if ctx.CookTimePreference == "" {
			ctx.CookTimePreference = strings.TrimSpace(pref.CookTimePreference)
		}
		if params.PeopleCount <= 0 && pref.PeopleCount > 0 {
			ctx.PeopleCount = pref.PeopleCount
		}
	}
	return ctx
}

func dishFromRecipe(recipe *model.Recipe, dishType string, reason string) DishResult {
	ingredients := recipeIngredientNames(recipe.Ingredients)
	return DishResult{
		RecipeID:     recipe.ID,
		Name:         recipe.Title,
		Type:         normalizeSceneDishType(dishType),
		CookTime:     recipe.CookTime,
		Difficulty:   recipe.Difficulty,
		Ingredients:  ingredients,
		StepsSummary: strings.TrimSpace(reason),
	}
}

func (s *RecommendService) RecommendSceneByAI(ctx context.Context, userID uint, params *RecommendParams) (*RecommendResult, error) {
	if s.aiClient == nil || !s.aiClient.IsConfigured() || s.recipeService == nil {
		return nil, ErrAIConfigMissing
	}

	var pref *model.UserPreference
	var err error
	if s.prefRepo != nil && userID > 0 {
		pref, err = s.prefRepo.FindByUserID(userID)
		if err != nil {
			return nil, err
		}
	}

	sceneCtx := ensureSceneDefaults(params, pref)
	draft, err := s.aiClient.GenerateSceneRecipeDrafts(ctx, sceneCtx)
	if err != nil {
		return nil, err
	}

	result := &RecommendResult{
		MenuName: draft.MenuName,
		Reason:   draft.Reason + " 已同步到菜谱库，可点击查看完整做法。",
		Source:   "ai",
	}
	shoppingSet := make(map[string]bool)
	for _, item := range draft.Dishes {
		if recipeHasAvoidIngredient(recipeIngredientNamesFromDraft(item.Recipe.Ingredients), sceneCtx.AvoidIngredients) {
			continue
		}
		created, err := s.recipeService.CreateOrReuseAIRecipeDraft(&item.Recipe)
		if err != nil {
			continue
		}
		dish := dishFromRecipe(created.Recipe, item.Type, item.Reason)
		for _, ingredient := range dish.Ingredients {
			shoppingSet[ingredient] = true
		}
		result.Dishes = append(result.Dishes, dish)
	}
	if len(result.Dishes) == 0 {
		return nil, ErrAIInvalidResponse
	}
	for item := range shoppingSet {
		result.ShoppingList = append(result.ShoppingList, item)
	}
	sort.Strings(result.ShoppingList)
	return result, nil
}

func recipeIngredientNamesFromDraft(items []AIRecipeIngredient) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Name) != "" {
			names = append(names, item.Name)
		}
	}
	return names
}

func (s *RecommendService) RecommendMenu(params *RecommendParams) (*RecommendResult, error) {
	recipes, err := s.recipeRepo.FindHot(100)
	if err != nil {
		return nil, err
	}

	var scored []scoredRecipe
	for i := range recipes {
		r := &recipes[i]
		if containsRecipeID(params.ExcludeRecipeIDs, r.ID) {
			continue
		}
		if recipeHasAvoidIngredient(recipeIngredientNames(r.Ingredients), params.AvoidIngredients) {
			continue
		}
		score := s.calcScore(*r, params)
		scored = append(scored, scoredRecipe{recipe: r, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := &RecommendResult{
		MenuName: params.MealType + "推荐菜单",
		Reason:   "根据您的口味偏好和饮食目标智能推荐",
	}

	dishTypes := []string{"主菜", "配菜", "汤"}
	usedIDs := make(map[uint]bool)
	shoppingSet := make(map[string]bool)

	for i, dishType := range dishTypes {
		if i >= len(scored) {
			break
		}
		for _, sr := range scored {
			if usedIDs[sr.recipe.ID] {
				continue
			}
			usedIDs[sr.recipe.ID] = true

			ingredients := recipeIngredientNames(sr.recipe.Ingredients)
			for _, item := range ingredients {
				shoppingSet[item] = true
			}

			result.Dishes = append(result.Dishes, DishResult{
				RecipeID:    sr.recipe.ID,
				Name:        sr.recipe.Title,
				Type:        dishType,
				CookTime:    sr.recipe.CookTime,
				Difficulty:  sr.recipe.Difficulty,
				Ingredients: ingredients,
			})
			break
		}
	}

	for item := range shoppingSet {
		result.ShoppingList = append(result.ShoppingList, item)
	}

	return result, nil
}

func (s *RecommendService) calcScore(r model.Recipe, params *RecommendParams) float64 {
	var score float64

	if len(params.TastePreference) > 0 {
		for _, t := range params.TastePreference {
			if strings.TrimSpace(r.Taste) == strings.TrimSpace(t) {
				score += 20
				break
			}
		}
	}

	if params.Scene == "quick_meal" && r.CookTime > 0 && r.CookTime <= 30 {
		score += 25
	}
	if params.Scene == "family_dinner" && r.PeopleCount >= 4 {
		score += 18
	}
	if params.Scene == "late_night" && r.CookTime > 0 && r.CookTime <= 20 {
		score += 16
	}
	if strings.Contains(params.CookTimePreference, "30") && r.CookTime > 0 && r.CookTime <= 30 {
		score += 10
	}
	if strings.Contains(params.CookTimePreference, "20") && r.CookTime > 0 && r.CookTime <= 20 {
		score += 10
	}
	if params.HealthGoal != "" && len(r.HealthTags) > 0 && strings.Contains(string(r.HealthTags), params.HealthGoal) {
		score += 12
	}

	if r.CookTime <= 15 {
		score += 15
	} else if r.CookTime <= 30 {
		score += 10
	}

	score += float64(r.FavoriteCount) * 0.01
	score += float64(r.ViewCount) * 0.001

	return score
}

func (s *RecommendService) RecommendByIngredients(ingredients []string, mode string, limit int) (*IngredientRecommendPayload, error) {
	recipes, err := s.recipeRepo.FindHot(50)
	if err != nil {
		return nil, err
	}

	normalized := normalizeIngredientList(ingredients, 20)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	var results []IngredientRecommendResult
	ingredientSet := make(map[string]bool)
	for _, ing := range normalized {
		ingredientSet[ing] = true
	}

	for _, r := range recipes {
		recipeIngredients := normalizeIngredientList(recipeIngredientNames(r.Ingredients), 100)
		if len(recipeIngredients) == 0 {
			continue
		}

		matched := make([]string, 0)
		missing := make([]string, 0)
		for _, ri := range recipeIngredients {
			if ingredientSet[ri] {
				matched = append(matched, ri)
			} else {
				missing = append(missing, ri)
			}
		}

		if len(matched) > 0 {
			matchRate := float64(len(matched)) / float64(len(recipeIngredients))
			reason := "匹配到您的现有食材"
			if len(missing) == 0 {
				reason = "现有食材已足够制作"
			} else if mode == "fridge" {
				reason = "优先消耗冰箱现有食材"
			}
			results = append(results, IngredientRecommendResult{
				Recipe:             r,
				MatchRate:          matchRate,
				MatchedIngredients: matched,
				MissingIngredients: missing,
				Reason:             reason,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].MatchRate == results[j].MatchRate {
			iMissing := len(results[i].MissingIngredients)
			jMissing := len(results[j].MissingIngredients)
			if mode == "fridge" && iMissing != jMissing {
				return iMissing < jMissing
			}
			if results[i].Recipe.FavoriteCount == results[j].Recipe.FavoriteCount {
				return results[i].Recipe.ViewCount > results[j].Recipe.ViewCount
			}
			return results[i].Recipe.FavoriteCount > results[j].Recipe.FavoriteCount
		}
		return results[i].MatchRate > results[j].MatchRate
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return &IngredientRecommendPayload{List: results, Total: len(results)}, nil
}

func (s *RecommendService) GenerateWeekMenu(params *RecommendParams) ([]map[string]interface{}, error) {
	var weekMenu []map[string]interface{}
	mealTypes := []string{"breakfast", "lunch", "dinner"}
	usedIDs := make(map[uint]bool)

	for day := 1; day <= 7; day++ {
		dayMenu := map[string]interface{}{
			"day":   day,
			"meals": []interface{}{},
		}

		for _, mealType := range mealTypes {
			p := *params
			p.MealType = mealType
			p.ExcludeRecipeIDs = make([]uint, 0, len(usedIDs))
			for id := range usedIDs {
				p.ExcludeRecipeIDs = append(p.ExcludeRecipeIDs, id)
			}
			result, err := s.RecommendMenu(&p)
			if err != nil {
				continue
			}
			for _, dish := range result.Dishes {
				usedIDs[dish.RecipeID] = true
			}
			dayMenu["meals"] = append(dayMenu["meals"].([]interface{}), result)
		}

		weekMenu = append(weekMenu, dayMenu)
	}

	return weekMenu, nil
}
