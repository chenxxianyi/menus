package admin

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/response"

	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) Get(c *gin.Context) {
	var userCount, recipeCount, favoriteCount, recommendCount int64

	h.db.Model(&model.User{}).Count(&userCount)
	h.db.Model(&model.Recipe{}).Count(&recipeCount)
	h.db.Model(&model.Favorite{}).Count(&favoriteCount)
	h.db.Model(&model.RecommendLog{}).Count(&recommendCount)

	response.Success(c, gin.H{
		"user_count":      userCount,
		"recipe_count":    recipeCount,
		"favorite_count":  favoriteCount,
		"recommend_count": recommendCount,
	})
}
