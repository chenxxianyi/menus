package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidShoppingItemIndices = errors.New("invalid shopping item indices")
var ErrRecipeIngredientsEmpty = errors.New("recipe ingredients empty")
var ErrShoppingRecipeIDsEmpty = errors.New("shopping recipe ids empty")

type ShoppingService struct {
	repo       *repository.ShoppingRepo
	recipeRepo *repository.RecipeRepo
	aiClient   *AIClient
	aiLogSvc   *AIGenerationLogService
	recipeSvc  *RecipeService
}

func NewShoppingService(repo *repository.ShoppingRepo, recipeRepo *repository.RecipeRepo) *ShoppingService {
	return &ShoppingService{repo: repo, recipeRepo: recipeRepo}
}

func (s *ShoppingService) SetAIClient(aiClient *AIClient) {
	s.aiClient = aiClient
}

func (s *ShoppingService) SetAIGenerationLogService(aiLogSvc *AIGenerationLogService) {
	s.aiLogSvc = aiLogSvc
}

func (s *ShoppingService) SetRecipeService(recipeSvc *RecipeService) {
	s.recipeSvc = recipeSvc
}

func (s *ShoppingService) GetLists(userID uint) ([]model.ShoppingList, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ShoppingService) CreateList(list *model.ShoppingList) error {
	return s.repo.Create(list)
}

func (s *ShoppingService) UpdateList(list *model.ShoppingList) error {
	return s.repo.Update(list)
}

func (s *ShoppingService) DeleteList(id uint) error {
	return s.repo.Delete(id)
}

func (s *ShoppingService) GetListByID(id uint) (*model.ShoppingList, error) {
	return s.repo.FindByID(id)
}

func (s *ShoppingService) DeleteItems(userID, listID uint, indices []int) (*model.ShoppingList, int, error) {
	list, err := s.repo.FindByIDAndUserID(listID, userID)
	if err != nil {
		return nil, 0, err
	}

	items, deletedCount, err := removeShoppingItems(list.ItemsJSON, indices)
	if err != nil {
		return nil, 0, err
	}

	if err := s.repo.UpdateItemsByUserID(listID, userID, items); err != nil {
		return nil, 0, err
	}

	list.ItemsJSON = items
	return list, deletedCount, nil
}

type DishShoppingItem struct {
	Name     string  `json:"name"`
	Amount   string  `json:"amount"`
	Emoji    string  `json:"emoji"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Checked  bool    `json:"checked"`
	Status   string  `json:"status,omitempty"`
}

type DishShoppingListResult struct {
	List   *model.ShoppingList `json:"list"`
	Recipe *model.Recipe       `json:"recipe"`
	Items  []DishShoppingItem  `json:"items"`
}

type RecipesShoppingListResult struct {
	List    *model.ShoppingList `json:"list"`
	Recipes []model.Recipe      `json:"recipes"`
	Items   []DishShoppingItem  `json:"items"`
}

func (s *ShoppingService) GenerateFromDish(userID uint, dishName string, preview bool) (*DishShoppingListResult, error) {
	name := strings.TrimSpace(dishName)
	if name == "" {
		return nil, gorm.ErrRecordNotFound
	}

	recipe, err := s.recipeRepo.FindBestMatch(name)
	if err != nil {
		return nil, err
	}

	items, err := shoppingItemsFromRecipe(recipe.Ingredients)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrRecipeIngredientsEmpty
	}

	var list *model.ShoppingList
	if !preview {
		itemsJSON, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("encode generated shopping items: %w", err)
		}

		list = &model.ShoppingList{
			UserID:    userID,
			Name:      recipe.Title + "采购清单",
			ItemsJSON: model.JSON(itemsJSON),
		}
		if err := s.repo.Create(list); err != nil {
			return nil, err
		}
	}

	return &DishShoppingListResult{
		List:   list,
		Recipe: recipe,
		Items:  items,
	}, nil
}

func (s *ShoppingService) GenerateFromRecipe(userID uint, recipeID uint, preview bool) (*DishShoppingListResult, error) {
	if recipeID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	recipe, err := s.recipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, err
	}

	items, err := shoppingItemsFromRecipe(recipe.Ingredients)
	if err != nil {
		return nil, err
	}

	var list *model.ShoppingList
	if !preview {
		list, err = s.createShoppingListFromItems(userID, recipe.Title+"采购清单", items)
		if err != nil {
			return nil, err
		}
	}

	return &DishShoppingListResult{
		List:   list,
		Recipe: recipe,
		Items:  items,
	}, nil
}

func (s *ShoppingService) GenerateFromRecipes(userID uint, recipeIDs []uint, name string, preview bool) (*RecipesShoppingListResult, error) {
	ids := normalizeRecipeIDs(recipeIDs)
	if len(ids) == 0 {
		return nil, ErrShoppingRecipeIDsEmpty
	}

	recipes, err := s.recipeRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	if len(recipes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	items := make([]DishShoppingItem, 0)
	for i := range recipes {
		recipeItems, err := shoppingItemsFromRecipe(recipes[i].Ingredients)
		if err != nil {
			continue
		}
		items = append(items, recipeItems...)
	}
	items = mergeShoppingItems(items)
	if len(items) == 0 {
		return nil, ErrRecipeIngredientsEmpty
	}

	listName := strings.TrimSpace(name)
	if listName == "" {
		listName = "推荐采购清单"
	}

	var list *model.ShoppingList
	if !preview {
		list, err = s.createShoppingListFromItems(userID, listName, items)
		if err != nil {
			return nil, err
		}
	}

	return &RecipesShoppingListResult{
		List:    list,
		Recipes: recipes,
		Items:   items,
	}, nil
}

func (s *ShoppingService) GenerateFromDishByAI(ctx context.Context, userID uint, dishName string, preview bool) (result *DishShoppingListResult, err error) {
	start := time.Now()
	name := strings.TrimSpace(dishName)
	defer func() {
		status := "success"
		if err != nil {
			status = "failed"
		}
		s.aiLogSvc.Record(AIGenerationLogPayload{
			UserID:         userID,
			GenerationType: "shopping",
			Model:          aiModelName(s.aiClient),
			Input:          map[string]interface{}{"dish_name": name, "preview": preview},
			Output:         summarizeShoppingAIResult(result),
			Status:         status,
			ErrorMessage:   aiErrorText(err),
			Duration:       time.Since(start),
			RecipeIDs:      recipeIDsFromShoppingAIResult(result),
		})
	}()
	if name == "" {
		return nil, ErrAIInvalidResponse
	}
	if s.recipeSvc == nil {
		return nil, ErrAIConfigMissing
	}
	if s.aiClient == nil || !s.aiClient.IsConfigured() {
		return nil, ErrAIConfigMissing
	}

	generatedRecipe, err := s.recipeSvc.generateRecipeByAI(ctx, userID, name, false)
	if err != nil {
		return nil, err
	}
	items, err := shoppingItemsFromRecipe(generatedRecipe.Recipe.Ingredients)
	if err != nil {
		return nil, err
	}

	var list *model.ShoppingList
	if !preview {
		itemsJSON, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("encode ai shopping items: %w", err)
		}

		list = &model.ShoppingList{
			UserID:    userID,
			Name:      generatedRecipe.Recipe.Title + "采购清单",
			ItemsJSON: model.JSON(itemsJSON),
		}
		if err := s.repo.Create(list); err != nil {
			return nil, err
		}
	}

	return &DishShoppingListResult{
		List:   list,
		Recipe: generatedRecipe.Recipe,
		Items:  items,
	}, nil
}

func (s *ShoppingService) createShoppingListFromItems(userID uint, name string, items []DishShoppingItem) (*model.ShoppingList, error) {
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("encode generated shopping items: %w", err)
	}

	list := &model.ShoppingList{
		UserID:    userID,
		Name:      name,
		ItemsJSON: model.JSON(itemsJSON),
	}
	if err := s.repo.Create(list); err != nil {
		return nil, err
	}
	return list, nil
}

func shoppingItemsFromRecipe(raw model.JSON) ([]DishShoppingItem, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, ErrRecipeIngredientsEmpty
	}

	var ingredients []struct {
		Name     string      `json:"name"`
		Amount   interface{} `json:"amount"`
		Unit     string      `json:"unit"`
		Emoji    string      `json:"emoji"`
		Category string      `json:"category"`
		Price    float64     `json:"price"`
	}
	if err := json.Unmarshal(raw, &ingredients); err != nil {
		return nil, fmt.Errorf("decode recipe ingredients: %w", err)
	}

	items := make([]DishShoppingItem, 0, len(ingredients))
	seen := make(map[string]int)
	for _, ingredient := range ingredients {
		name := strings.TrimSpace(ingredient.Name)
		if name == "" {
			continue
		}

		item := DishShoppingItem{
			Name:     name,
			Amount:   formatIngredientAmount(ingredient.Amount, ingredient.Unit),
			Emoji:    ingredient.Emoji,
			Category: inferShoppingCategory(name, ingredient.Category),
			Price:    ingredient.Price,
			Checked:  false,
			Status:   "pending",
		}

		key := normalizeShoppingItemName(name)
		if idx, ok := seen[key]; ok {
			items[idx].Amount = mergeAmountText(items[idx].Amount, item.Amount)
			continue
		}
		seen[key] = len(items)
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrRecipeIngredientsEmpty
	}

	return items, nil
}

func formatIngredientAmount(amount interface{}, unit string) string {
	var amountText string
	switch value := amount.(type) {
	case string:
		amountText = strings.TrimSpace(value)
	case float64:
		if value == float64(int64(value)) {
			amountText = fmt.Sprintf("%d", int64(value))
		} else {
			amountText = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
		}
	case nil:
		amountText = ""
	default:
		amountText = strings.TrimSpace(fmt.Sprint(value))
	}

	unit = strings.TrimSpace(unit)
	if amountText == "" && unit == "" {
		return "按菜谱适量"
	}
	if amountText == "" {
		return unit
	}
	if unit == "" || strings.Contains(amountText, unit) {
		return amountText
	}
	return amountText + unit
}

func normalizeShoppingItemName(name string) string {
	value := strings.ReplaceAll(strings.TrimSpace(name), " ", "")
	value = strings.ReplaceAll(value, "　", "")
	switch value {
	case "番茄":
		return "西红柿"
	case "蛋":
		return "鸡蛋"
	case "马铃薯", "洋芋":
		return "土豆"
	default:
		return value
	}
}

func mergeAmountText(left, right string) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || left == "按菜谱适量" {
		return right
	}
	if right == "" || right == "按菜谱适量" || right == left {
		return left
	}
	return left + "、" + right
}

var shoppingAmountPattern = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)(.*)$`)

func mergeShoppingItems(items []DishShoppingItem) []DishShoppingItem {
	merged := make([]DishShoppingItem, 0, len(items))
	seen := make(map[string]int)
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		item.Name = name
		item.Category = inferShoppingCategory(item.Name, item.Category)
		item.Status = normalizeShoppingItemStatus(item.Status, item.Checked)
		item.Checked = item.Status == "bought"

		key := normalizeShoppingItemName(item.Name)
		if idx, ok := seen[key]; ok {
			existing := &merged[idx]
			existing.Amount = mergeShoppingAmount(existing.Amount, item.Amount)
			existing.Status = mergeShoppingItemStatus(existing.Status, item.Status)
			existing.Checked = existing.Status == "bought"
			if existing.Emoji == "" {
				existing.Emoji = item.Emoji
			}
			if existing.Category == "" {
				existing.Category = item.Category
			}
			existing.Price += item.Price
			continue
		}
		seen[key] = len(merged)
		merged = append(merged, item)
	}
	return merged
}

func normalizeShoppingItemStatus(status string, checked bool) string {
	switch strings.TrimSpace(status) {
	case "pending", "bought", "owned":
		return strings.TrimSpace(status)
	default:
		if checked {
			return "bought"
		}
		return "pending"
	}
}

func mergeShoppingItemStatus(left, right string) string {
	left = normalizeShoppingItemStatus(left, false)
	right = normalizeShoppingItemStatus(right, false)
	if left == "pending" || right == "pending" {
		return "pending"
	}
	if left == "owned" || right == "owned" {
		return "owned"
	}
	return "bought"
}

func mergeShoppingAmount(left, right string) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || left == "按菜谱适量" || left == "适量" {
		if right == "" {
			return "适量"
		}
		return right
	}
	if right == "" || right == "按菜谱适量" || right == "适量" || right == left {
		return left
	}

	leftNum, leftUnit, leftOK := splitAmount(left)
	rightNum, rightUnit, rightOK := splitAmount(right)
	if leftOK && rightOK && leftUnit == rightUnit {
		total := leftNum + rightNum
		return formatNumber(total) + leftUnit
	}
	return left + "、" + right
}

func splitAmount(text string) (float64, string, bool) {
	matches := shoppingAmountPattern.FindStringSubmatch(strings.TrimSpace(text))
	if len(matches) != 3 {
		return 0, "", false
	}
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", false
	}
	return value, strings.TrimSpace(matches[2]), true
}

func formatNumber(value float64) string {
	if value == float64(int64(value)) {
		return fmt.Sprintf("%d", int64(value))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
}

func normalizeRecipeIDs(ids []uint) []uint {
	result := make([]uint, 0, len(ids))
	seen := make(map[uint]bool, len(ids))
	for _, id := range ids {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, id)
	}
	return result
}

func inferShoppingCategory(name, category string) string {
	if strings.TrimSpace(category) != "" {
		return category
	}
	if strings.ContainsAny(name, "鸡鸭鱼虾肉牛猪羊蛋排骨") {
		return "肉蛋水产"
	}
	if strings.ContainsAny(name, "盐糖酱醋油姜蒜葱胡椒淀粉料酒") {
		return "调味"
	}
	if strings.ContainsAny(name, "米面粥饭馒头粉") {
		return "主食"
	}
	return "蔬菜"
}

func removeShoppingItems(itemsJSON model.JSON, indices []int) (model.JSON, int, error) {
	if len(indices) == 0 {
		return nil, 0, fmt.Errorf("%w: indices cannot be empty", ErrInvalidShoppingItemIndices)
	}

	items := make([]json.RawMessage, 0)
	if len(itemsJSON) > 0 && string(itemsJSON) != "null" {
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			return nil, 0, fmt.Errorf("decode shopping items: %w", err)
		}
	}

	removeSet := make(map[int]struct{}, len(indices))
	for _, index := range indices {
		if index < 0 || index >= len(items) {
			return nil, 0, fmt.Errorf("%w: index %d out of range", ErrInvalidShoppingItemIndices, index)
		}
		if _, exists := removeSet[index]; exists {
			return nil, 0, fmt.Errorf("%w: duplicate index %d", ErrInvalidShoppingItemIndices, index)
		}
		removeSet[index] = struct{}{}
	}

	remaining := make([]json.RawMessage, 0, len(items)-len(removeSet))
	for index, item := range items {
		if _, remove := removeSet[index]; !remove {
			remaining = append(remaining, item)
		}
	}

	encoded, err := json.Marshal(remaining)
	if err != nil {
		return nil, 0, fmt.Errorf("encode shopping items: %w", err)
	}
	return model.JSON(encoded), len(removeSet), nil
}
