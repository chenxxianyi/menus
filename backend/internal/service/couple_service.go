package service

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type CoupleService struct {
	coupleRepo  *repository.CoupleRepo
	recipeRepo  *repository.RecipeRepo
	userRepo    *repository.UserRepo
	shoppingRepo *repository.ShoppingRepo
	prefRepo    *repository.UserPrefRepo
}

func NewCoupleService(coupleRepo *repository.CoupleRepo, recipeRepo *repository.RecipeRepo, userRepo *repository.UserRepo) *CoupleService {
	return &CoupleService{coupleRepo: coupleRepo, recipeRepo: recipeRepo, userRepo: userRepo}
}

func (s *CoupleService) SetShoppingRepo(shoppingRepo *repository.ShoppingRepo) {
	s.shoppingRepo = shoppingRepo
}

func (s *CoupleService) SetUserPrefRepo(prefRepo *repository.UserPrefRepo) {
	s.prefRepo = prefRepo
}

func generateInviteCode() (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		code[i] = chars[n.Int64()]
	}
	return string(code), nil
}

// GetInviteCode returns existing invite code or generates a new one
func (s *CoupleService) GetInviteCode(userID uint) (string, error) {
	// Check if already bound
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err == nil && binding.ID > 0 {
		return binding.InviteCode, nil
	}

	code, err := generateInviteCode()
	if err != nil {
		return "", err
	}
	return code, nil
}

// Bind creates a couple binding using an invite code
func (s *CoupleService) Bind(userID uint, inviteCode string) (*model.CoupleBinding, error) {
	// Check if user is already bound
	existing, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err == nil && existing.ID > 0 && existing.UserBID != nil {
		return nil, errors.New("你已经有情侣关系了，请先解除绑定")
	}

	// Find the binding by invite code
	binding, err := s.coupleRepo.FindBindingByInviteCode(inviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邀请码无效或已过期")
		}
		return nil, err
	}

	// Can't bind with yourself
	if binding.UserAID == userID || (binding.UserBID != nil && *binding.UserBID == userID) {
		return nil, errors.New("不能和自己绑定")
	}

	// Check if the other user is already bound (complete binding only)
	otherUser, _ := s.coupleRepo.FindActiveBindingByUser(binding.UserAID)
	if otherUser != nil && otherUser.ID > 0 && otherUser.UserBID != nil && otherUser.ID != binding.ID {
		return nil, errors.New("对方已经有情侣关系了")
	}

	// Update binding with the second user
	now := time.Now()
	binding.UserBID = &userID
	binding.Status = 1
	binding.BoundAt = now

	if err := s.coupleRepo.UpdateBinding(binding); err != nil {
		return nil, err
	}

	// Load user info
	binding.UserA, _ = s.userRepo.FindByID(binding.UserAID)
	binding.UserB, _ = s.userRepo.FindByID(userID)

	return binding, nil
}

// CreateInvite creates a new invite (pending binding)
func (s *CoupleService) CreateInvite(userID uint) (*model.CoupleBinding, error) {
	// Check if already bound (must be a complete binding with both users)
	existing, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err == nil && existing.ID > 0 {
		if existing.UserBID != nil {
			return existing, nil
		}
		// Stale binding without UserB — clean it up
		existing.Status = 2
		s.coupleRepo.UpdateBinding(existing)
	}

	// 清理旧的未完成邀请
	s.coupleRepo.CleanupStaleInvites(userID)

	code, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	binding := &model.CoupleBinding{
		UserAID:    userID,
		InviteCode: code,
		Status:     0, // pending until someone binds
		BoundAt:    time.Now(),
	}

	if err := s.coupleRepo.CreateBinding(binding); err != nil {
		return nil, err
	}

	return binding, nil
}

// GetInfo returns couple info for the given user
func (s *CoupleService) GetInfo(userID uint) (map[string]interface{}, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not bound
		}
		return nil, err
	}

	// 只有 status=1 且双方都绑定完成的才算真正的情侣关系
	if binding.Status != 1 || binding.UserBID == nil {
		return nil, nil
	}

	var partner *model.User
	if binding.UserAID == userID {
		partner = binding.UserB
	} else {
		partner = binding.UserA
	}

	if partner == nil {
		return nil, nil
	}

	partnerInfo := map[string]interface{}{
		"id":       partner.ID,
		"nickname": partner.Nickname,
		"avatar":   partner.Avatar,
	}

	return map[string]interface{}{
		"couple_id":   binding.ID,
		"couple_name": binding.CoupleName,
		"invite_code": binding.InviteCode,
		"partner":     partnerInfo,
		"bound_at":    binding.BoundAt,
	}, nil
}

// Unbind removes a couple binding
func (s *CoupleService) Unbind(userID uint) error {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return errors.New("未找到情侣关系")
	}

	binding.Status = 0
	return s.coupleRepo.UpdateBinding(binding)
}

// CreateOrder creates a new couple order
func (s *CoupleService) CreateOrder(userID uint, dishName string, recipeID *uint, mealType, mealDate, note string) (*model.CoupleOrder, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return nil, errors.New("请先绑定情侣关系")
	}

	if recipeID == nil {
		if recipe, err := s.recipeRepo.FindBestMatch(dishName); err == nil && recipe.ID > 0 {
			matchedID := recipe.ID
			recipeID = &matchedID
		}
	}

	order := &model.CoupleOrder{
		CoupleID: binding.ID,
		UserID:   userID,
		RecipeID: recipeID,
		DishName: dishName,
		MealType: mealType,
		MealDate: mealDate,
		Note:     note,
		Status:   0,
	}

	if err := s.coupleRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	// Load user info
	order.User, _ = s.userRepo.FindByID(userID)
	if order.RecipeID != nil {
		order.Recipe, _ = s.recipeRepo.FindByID(*order.RecipeID)
	}

	return order, nil
}

// GetOrders returns couple orders, optionally filtered by date
func (s *CoupleService) GetOrders(userID uint, mealDate string) ([]model.CoupleOrder, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return nil, errors.New("请先绑定情侣关系")
	}

	return s.coupleRepo.FindOrdersByCoupleID(binding.ID, mealDate)
}

// UpdateOrderStatus updates order status (confirm/cancel)
func (s *CoupleService) UpdateOrderStatus(orderID, userID uint, status int8) (*model.CoupleOrder, error) {
	order, err := s.coupleRepo.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("点餐记录不存在")
	}

	// Verify the order belongs to the user's couple
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil || binding.ID != order.CoupleID {
		return nil, errors.New("无权操作")
	}

	order.Status = status
	if err := s.coupleRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

// DeleteOrder deletes a couple order
func (s *CoupleService) DeleteOrder(orderID, userID uint) error {
	order, err := s.coupleRepo.FindOrderByID(orderID)
	if err != nil {
		return errors.New("点餐记录不存在")
	}

	// Only the creator can delete
	if order.UserID != userID {
		return errors.New("只能删除自己的点餐")
	}

	return s.coupleRepo.DeleteOrder(orderID)
}

// ShoppingListItem represents a merged shopping item
type ShoppingListItem struct {
	Name     string  `json:"name"`
	Amount   string  `json:"amount"`
	Category string  `json:"category"`
	Emoji    string  `json:"emoji,omitempty"`
	Price    float64 `json:"price"`
	Checked  bool    `json:"checked"`
	Status   string  `json:"status,omitempty"`
}

type CoupleMenuDish struct {
	RecipeID   uint   `json:"recipe_id"`
	Name       string `json:"name"`
	Reason     string `json:"reason"`
	Source     string `json:"source"`
	CookTime   int    `json:"cook_time"`
	Difficulty string `json:"difficulty"`
}

// GenerateShoppingList merges ingredients from couple orders and can create a copy for both users.
func (s *CoupleService) GenerateShoppingList(userID uint, mealDate, mealType string, saveShared bool) (map[string]interface{}, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return nil, errors.New("请先绑定情侣关系")
	}

	orders, err := s.coupleRepo.FindOrdersByCoupleIDAndType(binding.ID, mealDate, mealType)
	if err != nil {
		return nil, err
	}

	agreedDishes, compromiseDishes := s.buildCoupleMenuDishes(binding, orders)

	// Collect recipe IDs from agreed/confirmed/compromise dishes that have them.
	recipeIDSet := make(map[uint]bool)
	var recipeIDs []uint
	addRecipeID := func(id uint) {
		if id == 0 || recipeIDSet[id] {
			return
		}
		recipeIDSet[id] = true
		recipeIDs = append(recipeIDs, id)
	}
	for _, dish := range agreedDishes {
		addRecipeID(dish.RecipeID)
	}
	if len(agreedDishes) == 0 {
		for _, dish := range compromiseDishes {
			addRecipeID(dish.RecipeID)
		}
	}
	for _, o := range orders {
		if o.RecipeID != nil {
			addRecipeID(*o.RecipeID)
		}
	}

	// Fetch recipes to get ingredients
	ingredientMap := make(map[string]*ShoppingListItem)

	if len(recipeIDs) > 0 {
		recipes, _ := s.recipeRepo.FindByIDs(recipeIDs)
		for _, recipe := range recipes {
			var ingredients []struct {
				Name   string `json:"name"`
				Amount string `json:"amount"`
			}
			if len(recipe.Ingredients) > 0 {
				json.Unmarshal(recipe.Ingredients, &ingredients)
			}
			for _, ing := range ingredients {
				name := strings.TrimSpace(ing.Name)
				if name == "" {
					continue
				}
				if existing, ok := ingredientMap[ing.Name]; ok {
					if existing.Amount != ing.Amount {
						existing.Amount = existing.Amount + " + " + ing.Amount
					}
				} else {
					ingredientMap[ing.Name] = &ShoppingListItem{
						Name:     ing.Name,
						Amount:   firstNonEmpty(ing.Amount, "适量"),
						Category: inferShoppingCategory(ing.Name, ""),
						Status:   "pending",
					}
				}
			}
			// Also add seasonings
			var seasonings []struct {
				Name   string `json:"name"`
				Amount string `json:"amount"`
			}
			if len(recipe.Seasonings) > 0 {
				json.Unmarshal(recipe.Seasonings, &seasonings)
			}
			for _, s := range seasonings {
				name := strings.TrimSpace(s.Name)
				if name == "" {
					continue
				}
				if _, ok := ingredientMap[name]; !ok {
					ingredientMap[name] = &ShoppingListItem{
						Name:     name,
						Amount:   firstNonEmpty(s.Amount, "适量"),
						Category: "调味",
						Status:   "pending",
					}
				}
			}
		}
	}

	// Convert map to slice
	var shoppingList []ShoppingListItem
	for _, item := range ingredientMap {
		shoppingList = append(shoppingList, *item)
	}
	sort.Slice(shoppingList, func(i, j int) bool {
		if shoppingList[i].Category == shoppingList[j].Category {
			return shoppingList[i].Name < shoppingList[j].Name
		}
		return shoppingList[i].Category < shoppingList[j].Category
	})

	var sharedLists []model.ShoppingList
	if saveShared && len(shoppingList) > 0 {
		sharedLists, err = s.createSharedShoppingLists(binding, mealDate, shoppingList)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"orders":            orders,
		"agreed_dishes":     agreedDishes,
		"compromise_dishes": compromiseDishes,
		"shopping_list":     shoppingList,
		"shared_lists":      sharedLists,
		"saved_shared":      saveShared && len(sharedLists) > 0,
		"total_items":       len(shoppingList),
	}, nil
}

func (s *CoupleService) buildCoupleMenuDishes(binding *model.CoupleBinding, orders []model.CoupleOrder) ([]CoupleMenuDish, []CoupleMenuDish) {
	if binding == nil {
		return nil, nil
	}
	byKey := make(map[string][]model.CoupleOrder)
	for _, order := range orders {
		key := coupleDishKey(order)
		if key == "" {
			continue
		}
		byKey[key] = append(byKey[key], order)
	}

	agreed := make([]CoupleMenuDish, 0)
	for _, group := range byKey {
		userSet := make(map[uint]bool)
		var recipe *model.Recipe
		dishName := ""
		for _, order := range group {
			userSet[order.UserID] = true
			if dishName == "" {
				dishName = order.DishName
			}
			if recipe == nil && order.Recipe != nil {
				recipe = order.Recipe
			}
		}
		if len(userSet) < 2 {
			continue
		}
		agreed = append(agreed, coupleMenuDishFromRecipe(recipe, dishName, "双方都想吃，优先安排。", "overlap"))
	}
	sort.Slice(agreed, func(i, j int) bool { return agreed[i].Name < agreed[j].Name })
	if len(agreed) > 0 {
		return agreed, nil
	}

	compromise := s.buildCompromiseDishes(binding, orders)
	return nil, compromise
}

func coupleDishKey(order model.CoupleOrder) string {
	if order.RecipeID != nil && *order.RecipeID > 0 {
		return fmt.Sprintf("recipe:%d", *order.RecipeID)
	}
	return "dish:" + strings.TrimSpace(strings.ToLower(order.DishName))
}

func coupleMenuDishFromRecipe(recipe *model.Recipe, fallbackName, reason, source string) CoupleMenuDish {
	dish := CoupleMenuDish{
		Name:   strings.TrimSpace(fallbackName),
		Reason: reason,
		Source: source,
	}
	if recipe != nil {
		dish.RecipeID = recipe.ID
		dish.Name = recipe.Title
		dish.CookTime = recipe.CookTime
		dish.Difficulty = recipe.Difficulty
	}
	if dish.Name == "" {
		dish.Name = "合意推荐"
	}
	return dish
}

func (s *CoupleService) buildCompromiseDishes(binding *model.CoupleBinding, orders []model.CoupleOrder) []CoupleMenuDish {
	recipeMap := make(map[uint]*model.Recipe)
	for _, order := range orders {
		if order.Recipe != nil && order.Recipe.ID > 0 {
			recipeMap[order.Recipe.ID] = order.Recipe
		}
	}
	if len(recipeMap) == 0 {
		recipes, err := s.recipeRepo.FindHot(6)
		if err != nil {
			return nil
		}
		for i := range recipes {
			recipe := recipes[i]
			recipeMap[recipe.ID] = &recipe
		}
	}

	prefs := s.couplePreferenceContext(binding)
	type scored struct {
		recipe *model.Recipe
		score  int
		reason string
	}
	scoredRecipes := make([]scored, 0, len(recipeMap))
	for _, recipe := range recipeMap {
		if recipe == nil {
			continue
		}
		score, reason := scoreCompromiseRecipe(recipe, prefs)
		scoredRecipes = append(scoredRecipes, scored{recipe: recipe, score: score, reason: reason})
	}
	sort.Slice(scoredRecipes, func(i, j int) bool {
		if scoredRecipes[i].score == scoredRecipes[j].score {
			return scoredRecipes[i].recipe.FavoriteCount+scoredRecipes[i].recipe.ViewCount > scoredRecipes[j].recipe.FavoriteCount+scoredRecipes[j].recipe.ViewCount
		}
		return scoredRecipes[i].score > scoredRecipes[j].score
	})

	limit := 3
	if len(scoredRecipes) < limit {
		limit = len(scoredRecipes)
	}
	result := make([]CoupleMenuDish, 0, limit)
	for i := 0; i < limit; i++ {
		item := scoredRecipes[i]
		result = append(result, coupleMenuDishFromRecipe(item.recipe, "", item.reason, "compromise"))
	}
	return result
}

type couplePreferenceContext struct {
	tastes    map[string]int
	health    map[string]int
	cookQuick int
}

func (s *CoupleService) couplePreferenceContext(binding *model.CoupleBinding) couplePreferenceContext {
	ctx := couplePreferenceContext{
		tastes: make(map[string]int),
		health: make(map[string]int),
	}
	if s.prefRepo == nil || binding == nil {
		return ctx
	}
	userIDs := []uint{binding.UserAID}
	if binding.UserBID != nil {
		userIDs = append(userIDs, *binding.UserBID)
	}
	for _, userID := range userIDs {
		pref, err := s.prefRepo.FindByUserID(userID)
		if err != nil || pref == nil {
			continue
		}
		for _, taste := range jsonStringList(pref.TastePreference) {
			ctx.tastes[taste]++
		}
		if strings.TrimSpace(pref.HealthGoal) != "" {
			ctx.health[strings.TrimSpace(pref.HealthGoal)]++
		}
		if strings.Contains(pref.CookTimePreference, "20") || strings.Contains(pref.CookTimePreference, "30") {
			ctx.cookQuick++
		}
	}
	return ctx
}

func scoreCompromiseRecipe(recipe *model.Recipe, prefs couplePreferenceContext) (int, string) {
	score := recipe.FavoriteCount/10 + recipe.ViewCount/100
	reasons := make([]string, 0, 3)
	if count := prefs.tastes[strings.TrimSpace(recipe.Taste)]; count > 0 {
		score += 30 * count
		reasons = append(reasons, "兼顾双方口味")
	}
	for _, tag := range jsonStringList(recipe.HealthTags) {
		if count := prefs.health[tag]; count > 0 {
			score += 22 * count
			reasons = append(reasons, "照顾健康目标")
			break
		}
	}
	if prefs.cookQuick > 0 && recipe.CookTime > 0 && recipe.CookTime <= 30 {
		score += 16 * prefs.cookQuick
		reasons = append(reasons, "做起来不费时间")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "从两人的候选菜中选择热度更高、接受度更稳的一道")
	}
	return score, strings.Join(reasons, "，") + "。"
}

func (s *CoupleService) createSharedShoppingLists(binding *model.CoupleBinding, mealDate string, items []ShoppingListItem) ([]model.ShoppingList, error) {
	if s.shoppingRepo == nil || binding == nil || binding.UserBID == nil {
		return nil, errors.New("共享购物清单服务不可用")
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(mealDate)
	if name == "" {
		name = time.Now().Format("2006-01-02")
	}
	name += " 情侣点餐采购清单"
	userIDs := []uint{binding.UserAID, *binding.UserBID}
	lists := make([]model.ShoppingList, 0, len(userIDs))
	for _, uid := range userIDs {
		list := model.ShoppingList{
			UserID:    uid,
			Name:      name,
			ItemsJSON: model.JSON(raw),
		}
		if err := s.shoppingRepo.Create(&list); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

// Helper to get user binding info
func (s *CoupleService) IsUserBound(userID uint) bool {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	return err == nil && binding.ID > 0
}

func (s *CoupleService) GetPartnerID(userID uint) (uint, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return 0, err
	}
	if binding.UserAID == userID && binding.UserBID != nil {
		return *binding.UserBID, nil
	}
	return binding.UserAID, nil
}

func (s *CoupleService) GetBindingByUser(userID uint) (*model.CoupleBinding, error) {
	return s.coupleRepo.FindActiveBindingByUser(userID)
}

func (s *CoupleService) SetCoupleName(userID uint, name string) error {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return errors.New("未找到情侣关系")
	}
	binding.CoupleName = name
	return s.coupleRepo.UpdateBinding(binding)
}

func (s *CoupleService) FindOrCreateInvite(userID uint) (*model.CoupleBinding, error) {
	// Check if already bound (must be a complete binding with both users)
	existing, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err == nil && existing.ID > 0 && existing.UserBID != nil {
		return existing, nil
	}
	return s.CreateInvite(userID)
}

func (s *CoupleService) BindByCode(userID uint, code string) (*model.CoupleBinding, error) {
	return s.Bind(userID, code)
}

func (s *CoupleService) GetInviteOwner(inviteCode string) (*model.CoupleBinding, error) {
	return s.coupleRepo.FindBindingByInviteCode(inviteCode)
}

func (s *CoupleService) GetOrdersGroupedByDate(userID uint) (map[string][]model.CoupleOrder, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return nil, errors.New("请先绑定情侣关系")
	}

	orders, err := s.coupleRepo.FindOrdersByCoupleID(binding.ID, "")
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]model.CoupleOrder)
	for _, o := range orders {
		grouped[o.MealDate] = append(grouped[o.MealDate], o)
	}

	return grouped, nil
}
