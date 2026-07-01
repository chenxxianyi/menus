package service

import (
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

var ErrMenuRecipeIDsEmpty = errors.New("menu recipe ids empty")

type MenuService struct {
	menuRepo   *repository.MenuRepo
	recipeRepo *repository.RecipeRepo
}

type SaveMenuRequest struct {
	Name         string      `json:"name"`
	MenuType     string      `json:"menu_type"`
	MealType     string      `json:"meal_type"`
	PeopleCount  int         `json:"people_count"`
	Taste        string      `json:"taste"`
	HealthGoal   string      `json:"health_goal"`
	Dishes       interface{} `json:"dishes"`
	ShoppingList interface{} `json:"shopping_list"`
	Reason       string      `json:"reason"`
}

func NewMenuService(menuRepo *repository.MenuRepo, recipeRepo *repository.RecipeRepo) *MenuService {
	return &MenuService{menuRepo: menuRepo, recipeRepo: recipeRepo}
}

func (s *MenuService) Save(userID uint, req SaveMenuRequest) (*model.Menu, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = "我的菜单"
	}
	menuType := strings.TrimSpace(req.MenuType)
	if menuType == "" {
		menuType = "daily"
	}
	dishesJSON, err := json.Marshal(req.Dishes)
	if err != nil {
		return nil, err
	}
	shoppingJSON, err := json.Marshal(req.ShoppingList)
	if err != nil {
		return nil, err
	}
	menu := &model.Menu{
		UserID:           userID,
		MenuName:         name,
		MealType:         menuType,
		PeopleCount:      req.PeopleCount,
		Taste:            strings.TrimSpace(req.Taste),
		HealthGoal:       strings.TrimSpace(req.HealthGoal),
		DishesJSON:       model.JSON(dishesJSON),
		ShoppingListJSON: model.JSON(shoppingJSON),
		Reason:           strings.TrimSpace(req.Reason),
	}
	if strings.TrimSpace(req.MealType) != "" {
		menu.MealType = strings.TrimSpace(req.MealType)
	}
	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *MenuService) List(userID uint, page, pageSize int) ([]model.Menu, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.menuRepo.FindByUserID(userID, page, pageSize)
}

func (s *MenuService) Detail(userID, id uint) (*model.Menu, error) {
	if id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return s.menuRepo.FindByIDAndUserID(id, userID)
}

func (s *MenuService) Delete(userID, id uint) error {
	if id == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.menuRepo.DeleteByIDAndUserID(id, userID)
}

func (s *MenuService) Reuse(userID, id uint) (map[string]interface{}, error) {
	menu, err := s.Detail(userID, id)
	if err != nil {
		return nil, err
	}
	recipeIDs := extractRecipeIDsFromMenu(menu.DishesJSON)
	if len(recipeIDs) == 0 {
		return nil, ErrMenuRecipeIDsEmpty
	}
	recipes, err := s.recipeRepo.FindByIDs(recipeIDs)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"menu":       menu,
		"recipe_ids": recipeIDs,
		"recipes":    recipes,
	}, nil
}

func extractRecipeIDsFromMenu(raw model.JSON) []uint {
	var value interface{}
	if len(raw) == 0 || json.Unmarshal(raw, &value) != nil {
		return nil
	}
	ids := make([]uint, 0)
	seen := make(map[uint]bool)
	var walk func(interface{})
	walk = func(v interface{}) {
		switch item := v.(type) {
		case map[string]interface{}:
			for _, key := range []string{"recipe_id", "id"} {
				if id := numberToUint(item[key]); id > 0 && !seen[id] {
					seen[id] = true
					ids = append(ids, id)
				}
			}
			for _, child := range item {
				walk(child)
			}
		case []interface{}:
			for _, child := range item {
				walk(child)
			}
		}
	}
	walk(value)
	return ids
}

func numberToUint(value interface{}) uint {
	switch v := value.(type) {
	case float64:
		if v > 0 {
			return uint(v)
		}
	case int:
		if v > 0 {
			return uint(v)
		}
	case uint:
		return v
	}
	return 0
}
