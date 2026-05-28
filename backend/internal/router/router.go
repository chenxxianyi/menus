package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"menu-recommend/config"
	v1 "menu-recommend/internal/api/v1"
	"menu-recommend/internal/api/admin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/repository"
	"menu-recommend/internal/service"
)

func Setup(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware(cfg.CORS.Origins))

	// Repositories
	userRepo := repository.NewUserRepo(db)
	prefRepo := repository.NewUserPrefRepo(db)
	recipeRepo := repository.NewRecipeRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	favRepo := repository.NewFavoriteRepo(db)
	shoppingRepo := repository.NewShoppingRepo(db)
	bannerRepo := repository.NewBannerRepo(db)
	feedbackRepo := repository.NewFeedbackRepo(db)
	ingredientRepo := repository.NewIngredientRepo(db)
	_ = repository.NewMenuRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT.SecretKey, cfg.JWT.ExpireDuration())
	userService := service.NewUserService(userRepo, prefRepo)
	recipeService := service.NewRecipeService(recipeRepo, categoryRepo, favRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	bannerService := service.NewBannerService(bannerRepo)
	favService := service.NewFavoriteService(favRepo, recipeRepo)
	shoppingService := service.NewShoppingService(shoppingRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)
	recommendService := service.NewRecommendService(recipeRepo, ingredientRepo)

	// Handlers - v1
	authHandler := v1.NewAuthHandler(authService)
	userHandler := v1.NewUserHandler(userService)
	homeHandler := v1.NewHomeHandler(recipeService, bannerService, categoryService)
	recipeHandler := v1.NewRecipeHandler(recipeService)
	categoryHandler := v1.NewCategoryHandler(categoryService)
	bannerHandler := v1.NewBannerHandler(bannerService)
	favoriteHandler := v1.NewFavoriteHandler(favService)
	shoppingHandler := v1.NewShoppingHandler(shoppingService)
	feedbackHandler := v1.NewFeedbackHandler(feedbackService)
	recommendHandler := v1.NewRecommendHandler(recommendService)

	// Handlers - admin
	adminAuthHandler := admin.NewAuthHandler(authService, db)
	adminDashboardHandler := admin.NewDashboardHandler(db)
	adminRecipeHandler := admin.NewRecipeHandler(recipeService)
	adminCategoryHandler := admin.NewCategoryHandler(categoryService)
	adminIngredientHandler := admin.NewIngredientHandler(db)
	adminBannerHandler := admin.NewBannerHandler(bannerService)
	adminUserHandler := admin.NewUserHandler(userService)
	adminFeedbackHandler := admin.NewFeedbackHandler(feedbackService)

	// User API routes
	api := r.Group("/api")
	{
		// Public
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/home", homeHandler.GetHome)
		api.GET("/banners", bannerHandler.List)
		api.GET("/recipes", recipeHandler.List)
		api.GET("/recipes/:id", recipeHandler.Detail)
		api.GET("/categories", categoryHandler.List)

		// Authenticated
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(authService))
		{
			auth.GET("/user/info", userHandler.GetInfo)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.GET("/user/preferences", userHandler.GetPreferences)
			auth.PUT("/user/preferences", userHandler.UpdatePreferences)
			auth.POST("/recipes/:id/favorite", favoriteHandler.Add)
			auth.DELETE("/recipes/:id/favorite", favoriteHandler.Remove)
			auth.GET("/shopping-list", shoppingHandler.List)
			auth.POST("/shopping-list", shoppingHandler.Create)
			auth.PUT("/shopping-list/:id", shoppingHandler.Update)
			auth.DELETE("/shopping-list/:id", shoppingHandler.Delete)
			auth.POST("/recommend/menu", recommendHandler.Menu)
			auth.POST("/recommend/by-ingredients", recommendHandler.ByIngredients)
			auth.POST("/recommend/week-menu", recommendHandler.WeekMenu)
			auth.POST("/feedback", feedbackHandler.Create)
		}
	}

	// Admin routes
	adminGroup := r.Group("/api/admin")
	{
		adminGroup.POST("/login", adminAuthHandler.Login)

		adminAuth := adminGroup.Group("")
		adminAuth.Use(middleware.AdminMiddleware(authService))
		{
			adminAuth.GET("/dashboard", adminDashboardHandler.Get)
			adminAuth.GET("/recipes", adminRecipeHandler.List)
			adminAuth.POST("/recipes", adminRecipeHandler.Create)
			adminAuth.PUT("/recipes/:id", adminRecipeHandler.Update)
			adminAuth.DELETE("/recipes/:id", adminRecipeHandler.Delete)
			adminAuth.GET("/categories", adminCategoryHandler.List)
			adminAuth.POST("/categories", adminCategoryHandler.Create)
			adminAuth.PUT("/categories/:id", adminCategoryHandler.Update)
			adminAuth.DELETE("/categories/:id", adminCategoryHandler.Delete)
			adminAuth.GET("/ingredients", adminIngredientHandler.List)
			adminAuth.POST("/ingredients", adminIngredientHandler.Create)
			adminAuth.PUT("/ingredients/:id", adminIngredientHandler.Update)
			adminAuth.DELETE("/ingredients/:id", adminIngredientHandler.Delete)
			adminAuth.GET("/banners", adminBannerHandler.List)
			adminAuth.POST("/banners", adminBannerHandler.Create)
			adminAuth.PUT("/banners/:id", adminBannerHandler.Update)
			adminAuth.DELETE("/banners/:id", adminBannerHandler.Delete)
			adminAuth.GET("/users", adminUserHandler.List)
			adminAuth.GET("/users/:id", adminUserHandler.Detail)
			adminAuth.GET("/feedback", adminFeedbackHandler.List)
			adminAuth.PUT("/feedback/:id", adminFeedbackHandler.Update)
		}
	}

	return r
}
