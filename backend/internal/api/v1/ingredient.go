package v1

import (
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/repository"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	ingredientRepo *repository.IngredientRepo
}

func NewIngredientHandler(ingredientRepo *repository.IngredientRepo) *IngredientHandler {
	return &IngredientHandler{ingredientRepo: ingredientRepo}
}

func (h *IngredientHandler) Options(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")

	items, _, err := h.ingredientRepo.List(keyword, category, 1, 100)
	if err != nil {
		response.Error(c, errcode.ErrServer, "食材选项获取失败")
		return
	}

	response.Success(c, gin.H{"list": items})
}
