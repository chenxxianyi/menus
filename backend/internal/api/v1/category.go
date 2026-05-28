package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, categories)
}
