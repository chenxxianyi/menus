package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
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

func (h *CategoryHandler) Create(c *gin.Context) {
	var cat model.RecipeCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	if err := h.categoryService.Create(&cat); err != nil {
		response.Error(c, errcode.ErrServer, "创建失败")
		return
	}
	response.Success(c, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var cat model.RecipeCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	cat.ID = uint(id)
	if err := h.categoryService.Update(&cat); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.categoryService.Delete(uint(id)); err != nil {
		response.Error(c, errcode.ErrServer, "删除失败")
		return
	}
	response.Success(c, nil)
}
