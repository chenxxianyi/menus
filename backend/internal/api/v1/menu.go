package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	menus, total, err := h.menuService.List(userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询菜单失败")
		return
	}
	response.SuccessPage(c, menus, total, page, pageSize)
}

func (h *MenuHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req service.SaveMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "菜单参数错误")
		return
	}
	menu, err := h.menuService.Save(userID, req)
	if err != nil {
		response.Error(c, errcode.ErrServer, "保存菜单失败")
		return
	}
	response.Success(c, menu)
}

func (h *MenuHandler) Detail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	menu, err := h.menuService.Detail(userID, uint(id))
	if err != nil {
		response.Error(c, errcode.ErrNotFound, "菜单不存在")
		return
	}
	response.Success(c, menu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.menuService.Delete(userID, uint(id)); err != nil {
		response.Error(c, errcode.ErrNotFound, "菜单不存在")
		return
	}
	response.Success(c, nil)
}

func (h *MenuHandler) Reuse(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.menuService.Reuse(userID, uint(id))
	if err != nil {
		if errors.Is(err, service.ErrMenuRecipeIDsEmpty) {
			response.Error(c, errcode.ErrParam, "该菜单没有可复用的菜谱")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrNotFound, "菜单不存在")
			return
		}
		response.Error(c, errcode.ErrServer, "复用菜单失败")
		return
	}
	response.Success(c, result)
}
