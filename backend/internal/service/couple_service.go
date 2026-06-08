package service

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"gorm.io/gorm"

	"menu-recommend/internal/model"
	"menu-recommend/internal/repository"
)

type CoupleService struct {
	coupleRepo  *repository.CoupleRepo
	recipeRepo  *repository.RecipeRepo
	userRepo    *repository.UserRepo
}

func NewCoupleService(coupleRepo *repository.CoupleRepo, recipeRepo *repository.RecipeRepo, userRepo *repository.UserRepo) *CoupleService {
	return &CoupleService{coupleRepo: coupleRepo, recipeRepo: recipeRepo, userRepo: userRepo}
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
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Category string `json:"category"`
}

// GenerateShoppingList merges ingredients from confirmed orders
func (s *CoupleService) GenerateShoppingList(userID uint, mealDate, mealType string) (map[string]interface{}, error) {
	binding, err := s.coupleRepo.FindActiveBindingByUser(userID)
	if err != nil {
		return nil, errors.New("请先绑定情侣关系")
	}

	orders, err := s.coupleRepo.FindOrdersByCoupleIDAndType(binding.ID, mealDate, mealType)
	if err != nil {
		return nil, err
	}

	// Collect recipe IDs from orders that have them
	var recipeIDs []uint
	for _, o := range orders {
		if o.RecipeID != nil {
			recipeIDs = append(recipeIDs, *o.RecipeID)
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
				if existing, ok := ingredientMap[ing.Name]; ok {
					if existing.Amount != ing.Amount {
						existing.Amount = existing.Amount + " + " + ing.Amount
					}
				} else {
					ingredientMap[ing.Name] = &ShoppingListItem{
						Name:   ing.Name,
						Amount: ing.Amount,
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
				if _, ok := ingredientMap[s.Name]; !ok {
					ingredientMap[s.Name] = &ShoppingListItem{
						Name:   s.Name,
						Amount: s.Amount,
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

	return map[string]interface{}{
		"orders":        orders,
		"shopping_list": shoppingList,
		"total_items":   len(shoppingList),
	}, nil
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
