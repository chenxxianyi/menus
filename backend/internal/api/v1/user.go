package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "用户不存在")
		return
	}
	response.Success(c, user)
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if err := h.userService.UpdateProfile(userID, req.Nickname, req.Avatar, req.Gender); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) GetPreferences(c *gin.Context) {
	userID := middleware.GetUserID(c)
	pref, err := h.userService.GetPreferences(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "获取失败")
		return
	}
	response.Success(c, pref)
}

type UpdatePreferencesRequest struct {
	TastePreference     model.JSON `json:"taste_preference"`
	HealthGoal          string     `json:"health_goal"`
	AvoidIngredients    model.JSON `json:"avoid_ingredients"`
	FavoriteIngredients model.JSON `json:"favorite_ingredients"`
	CookTimePreference  string     `json:"cook_time_preference"`
	PeopleCount         int        `json:"people_count"`
}

func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	pref := &model.UserPreference{
		UserID:              userID,
		TastePreference:     req.TastePreference,
		HealthGoal:          req.HealthGoal,
		AvoidIngredients:    req.AvoidIngredients,
		FavoriteIngredients: req.FavoriteIngredients,
		CookTimePreference:  req.CookTimePreference,
		PeopleCount:         req.PeopleCount,
	}

	if err := h.userService.UpdatePreferences(pref); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, nil)
}
