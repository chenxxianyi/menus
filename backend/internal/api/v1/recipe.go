package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type RecipeHandler struct {
	recipeService *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

func (h *RecipeHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	taste := c.Query("taste")
	cookTime := c.Query("cook_time")
	difficulty := c.Query("difficulty")
	healthTags := c.Query("health_tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	recipes, total, err := h.recipeService.ListRecipes(keyword, uint(categoryID), taste, cookTime, difficulty, healthTags, page, pageSize)
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
