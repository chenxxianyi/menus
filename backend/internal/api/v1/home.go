package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/repository"
	"menu-recommend/internal/service"
)

type HomeHandler struct {
	recipeService   *service.RecipeService
	bannerService   *service.BannerService
	categoryService *service.CategoryService
	favoriteRepo    *repository.FavoriteRepo
}

func NewHomeHandler(rs *service.RecipeService, bs *service.BannerService, cs *service.CategoryService, favoriteRepo *repository.FavoriteRepo) *HomeHandler {
	return &HomeHandler{recipeService: rs, bannerService: bs, categoryService: cs, favoriteRepo: favoriteRepo}
}

func (h *HomeHandler) GetHome(c *gin.Context) {
	banners, _ := h.bannerService.GetActive()
	categories, _ := h.categoryService.GetAll()
	hotRecipes, _ := h.recipeService.GetHotRecipes(10)
	mealType := c.Query("meal_type")
	todayRecommend, _ := h.recipeService.GetMealRecommendations(mealType, 1)

	var recommend interface{}
	if len(todayRecommend) > 0 {
		recommend = todayRecommend[0]
	}
	userID := middleware.GetUserID(c)
	if userID > 0 && h.favoriteRepo != nil {
		for index := range hotRecipes {
			hotRecipes[index].IsFavorited = h.favoriteRepo.Exists(userID, hotRecipes[index].ID)
		}
		if len(todayRecommend) > 0 {
			todayRecommend[0].IsFavorited = h.favoriteRepo.Exists(userID, todayRecommend[0].ID)
			recommend = todayRecommend[0]
		}
	}

	response.Success(c, gin.H{
		"banners":         banners,
		"today_recommend": recommend,
		"categories":      categories,
		"hot_recipes":     hotRecipes,
		"meal_type":       mealType,
	})
}
