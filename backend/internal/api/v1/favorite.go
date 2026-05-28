package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type FavoriteHandler struct {
	favService *service.FavoriteService
}

func NewFavoriteHandler(favService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favService: favService}
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	recipeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.favService.AddFavorite(userID, uint(recipeID)); err != nil {
		response.Error(c, errcode.ErrServer, "收藏失败")
		return
	}
	response.Success(c, nil)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	recipeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam, "无效的ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.favService.RemoveFavorite(userID, uint(recipeID)); err != nil {
		response.Error(c, errcode.ErrServer, "取消收藏失败")
		return
	}
	response.Success(c, nil)
}
