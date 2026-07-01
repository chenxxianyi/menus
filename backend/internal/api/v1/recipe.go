package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type RecipeHandler struct {
	recipeService         *service.RecipeService
	recipeFeedbackService *service.UserRecipeFeedbackService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

func (h *RecipeHandler) SetRecipeFeedbackService(recipeFeedbackService *service.UserRecipeFeedbackService) {
	h.recipeFeedbackService = recipeFeedbackService
}

func (h *RecipeHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	taste := c.Query("taste")
	cookTime := c.Query("cook_time")
	difficulty := c.Query("difficulty")
	healthTags := c.Query("health_tags")
	sortBy := c.DefaultQuery("sort", "latest")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	if sortBy != "latest" && sortBy != "hot" {
		response.Error(c, errcode.ErrParam, "排序参数错误")
		return
	}

	recipes, total, err := h.recipeService.ListRecipes(keyword, uint(categoryID), taste, cookTime, difficulty, healthTags, sortBy, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}

	response.SuccessPage(c, recipes, total, page, pageSize)
}

func (h *RecipeHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	userID := middleware.GetUserID(c)
	recipe, err := h.recipeService.GetRecipeDetail(uint(id), userID)
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "菜谱不存在")
		return
	}

	response.Success(c, recipe)
}

func (h *RecipeHandler) FilterOptions(c *gin.Context) {
	options, err := h.recipeService.GetFilterOptions()
	if err != nil {
		response.Error(c, errcode.ErrServer, "筛选选项获取失败")
		return
	}
	response.Success(c, options)
}

type GenerateRecipeByAIRequest struct {
	DishName string `json:"dish_name" binding:"required"`
}

func (h *RecipeHandler) GenerateByAI(c *gin.Context) {
	var req GenerateRecipeByAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请输入想创建的菜品")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.recipeService.GenerateRecipeByAI(c.Request.Context(), userID, req.DishName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAIConfigMissing):
			response.Error(c, errcode.ErrParam, "AI 未配置，请先配置 AI_BASE_URL、AI_API_KEY 和 AI_MODEL")
		case errors.Is(err, service.ErrAIInvalidResponse):
			response.Error(c, errcode.ErrServer, "AI 返回的菜谱格式无效")
		case errors.Is(err, service.ErrAIUpstream):
			response.Error(c, errcode.ErrServer, "AI 服务暂时不可用或返回错误，请稍后重试")
		default:
			response.Error(c, errcode.ErrServer, "AI 生成菜谱失败")
		}
		return
	}
	response.Success(c, result)
}

type RecipeFeedbackRequest struct {
	Type   string `json:"type" binding:"required"`
	Source string `json:"source"`
}

func (h *RecipeHandler) SetFeedback(c *gin.Context) {
	if h.recipeFeedbackService == nil {
		response.Error(c, errcode.ErrServer, "菜谱反馈服务不可用")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, errcode.ErrParam, "无效的菜谱ID")
		return
	}

	var req RecipeFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "请选择反馈类型")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.recipeFeedbackService.Set(userID, uint(id), req.Type, req.Source); err != nil {
		if errors.Is(err, service.ErrInvalidRecipeFeedbackType) {
			response.Error(c, errcode.ErrParam, "反馈类型错误")
			return
		}
		response.Error(c, errcode.ErrServer, "提交反馈失败")
		return
	}

	status, _ := h.recipeFeedbackService.Status(userID, uint(id))
	response.Success(c, gin.H{"feedback": status})
}

func (h *RecipeHandler) DeleteFeedback(c *gin.Context) {
	if h.recipeFeedbackService == nil {
		response.Error(c, errcode.ErrServer, "菜谱反馈服务不可用")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, errcode.ErrParam, "无效的菜谱ID")
		return
	}

	feedbackType := c.Param("type")
	userID := middleware.GetUserID(c)
	if err := h.recipeFeedbackService.Delete(userID, uint(id), feedbackType); err != nil {
		if errors.Is(err, service.ErrInvalidRecipeFeedbackType) {
			response.Error(c, errcode.ErrParam, "反馈类型错误")
			return
		}
		response.Error(c, errcode.ErrServer, "取消反馈失败")
		return
	}

	status, _ := h.recipeFeedbackService.Status(userID, uint(id))
	response.Success(c, gin.H{"feedback": status})
}

func (h *RecipeHandler) GetUserRecipeFeedback(c *gin.Context) {
	if h.recipeFeedbackService == nil {
		response.Error(c, errcode.ErrServer, "菜谱反馈服务不可用")
		return
	}
	userID := middleware.GetUserID(c)
	items, err := h.recipeFeedbackService.List(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询反馈失败")
		return
	}
	response.Success(c, gin.H{"list": items})
}
