package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type CoupleHandler struct {
	coupleService *service.CoupleService
}

func NewCoupleHandler(coupleService *service.CoupleService) *CoupleHandler {
	return &CoupleHandler{coupleService: coupleService}
}

// GetInviteCode returns or generates an invite code
func (h *CoupleHandler) GetInviteCode(c *gin.Context) {
	userID := middleware.GetUserID(c)
	binding, err := h.coupleService.FindOrCreateInvite(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成邀请码失败")
		return
	}
	response.Success(c, gin.H{
		"invite_code": binding.InviteCode,
		"couple_id":   binding.ID,
	})
}

type BindRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

// Bind creates a couple binding
func (h *CoupleHandler) Bind(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req BindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入邀请码")
		return
	}

	// 前置检查：不能用自己生成的邀请码绑定
	info, _ := h.coupleService.GetInviteOwner(req.InviteCode)
	if info != nil && info.UserAID == userID {
		response.Error(c, errcode.ErrParam, "不能和自己绑定")
		return
	}

	binding, err := h.coupleService.BindByCode(userID, req.InviteCode)
	if err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, gin.H{
		"couple_id":   binding.ID,
		"couple_name": binding.CoupleName,
		"partner": gin.H{
			"id":       binding.UserA.ID,
			"nickname": binding.UserA.Nickname,
			"avatar":   binding.UserA.Avatar,
		},
	})
}

// GetInfo returns couple info
func (h *CoupleHandler) GetInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	info, err := h.coupleService.GetInfo(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	if info == nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, info)
}

// Unbind removes couple binding
func (h *CoupleHandler) Unbind(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.coupleService.Unbind(userID); err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}
	response.Success(c, nil)
}

type CreateOrderRequest struct {
	DishName string `json:"dish_name" binding:"required"`
	RecipeID *uint  `json:"recipe_id"`
	MealType string `json:"meal_type"`
	MealDate string `json:"meal_date"`
	Note     string `json:"note"`
}

// CreateOrder creates a couple order
func (h *CoupleHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入菜品名称")
		return
	}

	if req.MealDate == "" {
		req.MealDate = "today" // will be handled by frontend default
	}

	order, err := h.coupleService.CreateOrder(userID, req.DishName, req.RecipeID, req.MealType, req.MealDate, req.Note)
	if err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, order)
}

// GetOrders returns couple orders
func (h *CoupleHandler) GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	mealDate := c.Query("meal_date")

	orders, err := h.coupleService.GetOrders(userID, mealDate)
	if err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, orders)
}

type UpdateOrderStatusRequest struct {
	Status int8 `json:"status"`
}

// UpdateOrderStatus updates order status
func (h *CoupleHandler) UpdateOrderStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	order, err := h.coupleService.UpdateOrderStatus(uint(orderID), userID, req.Status)
	if err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, order)
}

// DeleteOrder deletes a couple order
func (h *CoupleHandler) DeleteOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	if err := h.coupleService.DeleteOrder(uint(orderID), userID); err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, nil)
}

type GenerateShoppingListRequest struct {
	MealDate string `json:"meal_date"`
	MealType string `json:"meal_type"`
}

// GenerateShoppingList generates merged shopping list
func (h *CoupleHandler) GenerateShoppingList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req GenerateShoppingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	result, err := h.coupleService.GenerateShoppingList(userID, req.MealDate, req.MealType)
	if err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, result)
}

type SetCoupleNameRequest struct {
	CoupleName string `json:"couple_name" binding:"required"`
}

// SetCoupleName sets the couple nickname
func (h *CoupleHandler) SetCoupleName(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req SetCoupleNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入昵称")
		return
	}

	if err := h.coupleService.SetCoupleName(userID, req.CoupleName); err != nil {
		response.Error(c, errcode.ErrParam, err.Error())
		return
	}

	response.Success(c, nil)
}
