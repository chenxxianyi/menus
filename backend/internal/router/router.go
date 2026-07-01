package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"menu-recommend/config"
	"menu-recommend/internal/api/admin"
	v1 "menu-recommend/internal/api/v1"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/repository"
	"menu-recommend/internal/service"
)

func Setup(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware(cfg.CORS.Origins))
	r.Static("/uploads", cfg.Upload.Dir)

	// Repositories
	userRepo := repository.NewUserRepo(db)
	prefRepo := repository.NewUserPrefRepo(db)
	recipeRepo := repository.NewRecipeRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	favRepo := repository.NewFavoriteRepo(db)
	shoppingRepo := repository.NewShoppingRepo(db)
	bannerRepo := repository.NewBannerRepo(db)
	appConfigRepo := repository.NewAppConfigRepo(db)
	feedbackRepo := repository.NewFeedbackRepo(db)
	ingredientRepo := repository.NewIngredientRepo(db)
	menuRepo := repository.NewMenuRepo(db)
	coupleRepo := repository.NewCoupleRepo(db)
	historyRepo := repository.NewBrowseHistoryRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT.SecretKey, cfg.JWT.ExpireDuration())
	userService := service.NewUserService(userRepo, prefRepo)
	recipeService := service.NewRecipeService(recipeRepo, categoryRepo, favRepo, historyRepo)
	aiClient := service.NewAIClient(cfg.AI)
	recipeService.SetAIClient(aiClient)
	categoryService := service.NewCategoryService(categoryRepo)
	bannerService := service.NewBannerService(bannerRepo)
	appConfigService := service.NewAppConfigService(appConfigRepo)
	favService := service.NewFavoriteService(favRepo, recipeRepo)
	shoppingService := service.NewShoppingService(shoppingRepo, recipeRepo)
	shoppingService.SetAIClient(aiClient)
	shoppingService.SetRecipeService(recipeService)
	feedbackService := service.NewFeedbackService(feedbackRepo)
	recommendService := service.NewRecommendService(recipeRepo, ingredientRepo, prefRepo, recipeService, aiClient)
	coupleService := service.NewCoupleService(coupleRepo, recipeRepo, userRepo)
	historyService := service.NewBrowseHistoryService(historyRepo)

	// Handlers - v1
	authHandler := v1.NewAuthHandler(authService)
	userHandler := v1.NewUserHandler(userService)
	userStatsHandler := v1.NewUserStatsHandler(favRepo, menuRepo, shoppingRepo)
	homeHandler := v1.NewHomeHandler(recipeService, bannerService, categoryService)
	recipeHandler := v1.NewRecipeHandler(recipeService)
	historyHandler := v1.NewBrowseHistoryHandler(historyService)
	categoryHandler := v1.NewCategoryHandler(categoryService)
	bannerHandler := v1.NewBannerHandler(bannerService)
	appConfigHandler := v1.NewAppConfigHandler(appConfigService)
	favoriteHandler := v1.NewFavoriteHandler(favService)
	shoppingHandler := v1.NewShoppingHandler(shoppingService)
	feedbackHandler := v1.NewFeedbackHandler(feedbackService)
	recommendHandler := v1.NewRecommendHandler(recommendService)
	ingredientHandler := v1.NewIngredientHandler(ingredientRepo)
	coupleHandler := v1.NewCoupleHandler(coupleService)
	uploadHandler := v1.NewUploadHandler(&cfg.Upload)

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
		api.GET("/about", appConfigHandler.About)
		api.GET("/banners", bannerHandler.List)
		api.GET("/recipes/filter-options", recipeHandler.FilterOptions)
		api.GET("/recipes", recipeHandler.List)
		api.GET("/recipes/:id", middleware.OptionalAuth(authService), recipeHandler.Detail)
		api.GET("/categories", categoryHandler.List)
		api.GET("/ingredients/options", ingredientHandler.Options)

		// Authenticated
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(authService))
		{
			auth.GET("/user/info", userHandler.GetInfo)
			auth.GET("/user/stats", userStatsHandler.Get)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.POST("/upload/avatar", uploadHandler.Avatar)
			auth.GET("/user/preferences", userHandler.GetPreferences)
			auth.PUT("/user/preferences", userHandler.UpdatePreferences)
			auth.POST("/recipes/:id/favorite", favoriteHandler.Add)
			auth.DELETE("/recipes/:id/favorite", favoriteHandler.Remove)
			auth.POST("/recipes/generate-by-ai", recipeHandler.GenerateByAI)
			auth.GET("/user/favorites", favoriteHandler.List)
			auth.GET("/user/favorites/count", favoriteHandler.Count)
			auth.GET("/user/history", historyHandler.List)
			auth.DELETE("/user/history", historyHandler.Clear)
			auth.GET("/shopping-list", shoppingHandler.List)
			auth.POST("/shopping-list", shoppingHandler.Create)
			auth.POST("/shopping-list/generate-by-dish", shoppingHandler.GenerateByDish)
			auth.POST("/shopping-list/generate-by-ai", shoppingHandler.GenerateByAI)
			auth.PUT("/shopping-list/:id", shoppingHandler.Update)
			auth.DELETE("/shopping-list/:id/items", shoppingHandler.DeleteItems)
			auth.DELETE("/shopping-list/:id", shoppingHandler.Delete)
			auth.POST("/recommend/menu", recommendHandler.Menu)
			auth.POST("/recommend/menu-ai", recommendHandler.MenuAI)
			auth.POST("/recommend/by-ingredients", recommendHandler.ByIngredients)
			auth.POST("/recommend/week-menu", recommendHandler.WeekMenu)
			auth.POST("/feedback", feedbackHandler.Create)

			// Couple
			auth.GET("/couple/invite-code", coupleHandler.GetInviteCode)
			auth.POST("/couple/bind", coupleHandler.Bind)
			auth.GET("/couple/info", coupleHandler.GetInfo)
			auth.POST("/couple/unbind", coupleHandler.Unbind)
			auth.PUT("/couple/name", coupleHandler.SetCoupleName)
			auth.POST("/couple/orders", coupleHandler.CreateOrder)
			auth.GET("/couple/orders", coupleHandler.GetOrders)
			auth.PUT("/couple/orders/:id", coupleHandler.UpdateOrderStatus)
			auth.DELETE("/couple/orders/:id", coupleHandler.DeleteOrder)
			auth.POST("/couple/orders/generate-shopping-list", coupleHandler.GenerateShoppingList)
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
