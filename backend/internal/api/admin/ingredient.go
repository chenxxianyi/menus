package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"

	"gorm.io/gorm"
)

type IngredientHandler struct {
	db *gorm.DB
}

func NewIngredientHandler(db *gorm.DB) *IngredientHandler {
	return &IngredientHandler{db: db}
}

func (h *IngredientHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var items []model.Ingredient
	var total int64

	query := h.db.Model(&model.Ingredient{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)

	response.SuccessPage(c, items, total, page, pageSize)
}

func (h *IngredientHandler) Create(c *gin.Context) {
	var item model.Ingredient
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *IngredientHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var item model.Ingredient
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	item.ID = uint(id)
	h.db.Save(&item)
	response.Success(c, item)
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Delete(&model.Ingredient{}, id)
	response.Success(c, nil)
}
