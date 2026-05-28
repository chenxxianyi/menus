package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	recipes, total, err := h.recipeService.ListAllRecipes(keyword, uint(categoryID), nil, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.SuccessPage(c, recipes, total, page, pageSize)
}

func (h *RecipeHandler) Create(c *gin.Context) {
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if err := h.recipeService.CreateRecipe(&recipe); err != nil {
		response.Error(c, errcode.ErrServer, "创建失败")
		return
	}
	response.Success(c, recipe)
}

func (h *RecipeHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	recipe.ID = uint(id)

	if err := h.recipeService.UpdateRecipe(&recipe); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, recipe)
}

func (h *RecipeHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.recipeService.DeleteRecipe(uint(id)); err != nil {
		response.Error(c, errcode.ErrServer, "删除失败")
		return
	}
	response.Success(c, nil)
}
