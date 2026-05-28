package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type HomeHandler struct {
	recipeService   *service.RecipeService
	bannerService   *service.BannerService
	categoryService *service.CategoryService
}

func NewHomeHandler(rs *service.RecipeService, bs *service.BannerService, cs *service.CategoryService) *HomeHandler {
	return &HomeHandler{recipeService: rs, bannerService: bs, categoryService: cs}
}

func (h *HomeHandler) GetHome(c *gin.Context) {
	banners, _ := h.bannerService.GetActive()
	categories, _ := h.categoryService.GetAll()
	hotRecipes, _ := h.recipeService.GetHotRecipes(10)
	todayRecommend, _ := h.recipeService.GetRandomRecipes(1)

	var recommend interface{}
	if len(todayRecommend) > 0 {
		recommend = todayRecommend[0]
	}

	response.Success(c, gin.H{
		"banners":          banners,
		"today_recommend":  recommend,
		"categories":       categories,
		"hot_recipes":      hotRecipes,
	})
}
