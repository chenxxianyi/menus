package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type ShoppingHandler struct {
	shoppingService *service.ShoppingService
}

func NewShoppingHandler(shoppingService *service.ShoppingService) *ShoppingHandler {
	return &ShoppingHandler{shoppingService: shoppingService}
}

func (h *ShoppingHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	lists, err := h.shoppingService.GetLists(userID)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, lists)
}

type CreateShoppingRequest struct {
	Name      string    `json:"name"`
	ItemsJSON model.JSON `json:"items_json"`
}

func (h *ShoppingHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateShoppingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	list := &model.ShoppingList{
		UserID:    userID,
		Name:      req.Name,
		ItemsJSON: req.ItemsJSON,
	}

	if err := h.shoppingService.CreateList(list); err != nil {
		response.Error(c, errcode.ErrServer, "创建失败")
		return
	}
	response.Success(c, list)
}

type UpdateShoppingRequest struct {
	Name      string    `json:"name"`
	ItemsJSON model.JSON `json:"items_json"`
}

func (h *ShoppingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	existing, err := h.shoppingService.GetListByID(uint(id))
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "清单不存在")
		return
	}

	var req UpdateShoppingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if len(req.ItemsJSON) > 0 {
		existing.ItemsJSON = req.ItemsJSON
	}

	if err := h.shoppingService.UpdateList(existing); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, existing)
}

func (h *ShoppingHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	if err := h.shoppingService.DeleteList(uint(id)); err != nil {
		response.Error(c, errcode.ErrServer, "删除失败")
		return
	}
	response.Success(c, nil)
}
