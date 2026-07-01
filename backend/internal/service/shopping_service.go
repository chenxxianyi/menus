package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidShoppingItemIndices = errors.New("invalid shopping item indices")
var ErrRecipeIngredientsEmpty = errors.New("recipe ingredients empty")

type ShoppingService struct {
	repo       *repository.ShoppingRepo
	recipeRepo *repository.RecipeRepo
	aiClient   *AIClient
}

func NewShoppingService(repo *repository.ShoppingRepo, recipeRepo *repository.RecipeRepo) *ShoppingService {
	return &ShoppingService{repo: repo, recipeRepo: recipeRepo}
}

func (s *ShoppingService) SetAIClient(aiClient *AIClient) {
	s.aiClient = aiClient
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
}

type DishShoppingListResult struct {
	List   *model.ShoppingList `json:"list"`
	Recipe *model.Recipe       `json:"recipe"`
	Items  []DishShoppingItem  `json:"items"`
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

func (s *ShoppingService) GenerateFromDishByAI(ctx context.Context, userID uint, dishName string, preview bool) (*DishShoppingListResult, error) {
	name := strings.TrimSpace(dishName)
	if name == "" {
		return nil, ErrAIInvalidResponse
	}
	if s.aiClient == nil || !s.aiClient.IsConfigured() {
		return nil, ErrAIConfigMissing
	}

	items, err := s.aiClient.SuggestShoppingItems(ctx, name)
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
			Name:      name + " AI建议采购清单",
			ItemsJSON: model.JSON(itemsJSON),
		}
		if err := s.repo.Create(list); err != nil {
			return nil, err
		}
	}

	return &DishShoppingListResult{
		List:  list,
		Items: items,
	}, nil
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
	return strings.ReplaceAll(strings.TrimSpace(name), " ", "")
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
