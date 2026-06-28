package main

import (
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"menu-recommend/config"
	"menu-recommend/internal/model"
	"menu-recommend/internal/router"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Logger
	var logger *zap.Logger
	if cfg.App.Debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	// Database
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("connect database failed", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	if cfg.Database.ConnMaxLifetimeMinutes > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)
	}

	// Auto migrate
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&model.User{},
			&model.UserPreference{},
			&model.Recipe{},
			&model.RecipeCategory{},
			&model.Ingredient{},
			&model.Favorite{},
			&model.Menu{},
			&model.ShoppingList{},
			&model.RecommendLog{},
			&model.Feedback{},
			&model.Banner{},
			&model.AdminUser{},
			&model.CoupleBinding{},
			&model.CoupleOrder{},
		); err != nil {
			logger.Fatal("auto migrate failed", zap.Error(err))
		}
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_ = rdb

	// Router
	r := router.Setup(cfg, db, logger)

	// Start
	addr := cfg.Server.Addr()
	logger.Info("server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
	fmt.Printf("Server running on %s\n", addr)
}
