package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/repository"
)

type UserStatsHandler struct {
	favRepo      *repository.FavoriteRepo
	menuRepo     *repository.MenuRepo
	shoppingRepo *repository.ShoppingRepo
}

func NewUserStatsHandler(
	favRepo *repository.FavoriteRepo,
	menuRepo *repository.MenuRepo,
	shoppingRepo *repository.ShoppingRepo,
) *UserStatsHandler {
	return &UserStatsHandler{
		favRepo:      favRepo,
		menuRepo:     menuRepo,
		shoppingRepo: shoppingRepo,
	}
}

func (h *UserStatsHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)

	response.Success(c, gin.H{
		"favorite_count":      h.favRepo.CountByUserID(userID),
		"menu_count":          h.menuRepo.CountByUserID(userID),
		"shopping_list_count": h.shoppingRepo.CountByUserID(userID),
	})
}
