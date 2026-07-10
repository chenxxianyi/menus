package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
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
			&model.BrowseHistory{},
			&model.Feedback{},
			&model.UserRecipeFeedback{},
			&model.UserEvent{},
			&model.AIGenerationLog{},
			&model.Banner{},
			&model.AppConfig{},
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
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		if cfg.App.IsProduction() {
			logger.Fatal("connect redis failed", zap.Error(err))
		}
		logger.Warn("redis is unavailable; readiness checks will report it", zap.Error(err))
	}

	// Router
	r := router.Setup(cfg, db, logger, rdb)

	// Start and gracefully stop accepting new traffic on deployment signals.
	addr := cfg.Server.Addr()
	logger.Info("server starting", zap.String("addr", addr))
	server := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 10 * time.Second}
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	case <-signalContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			logger.Error("graceful shutdown failed", zap.Error(err))
		}
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	_ = rdb.Close()
	fmt.Printf("Server stopped on %s\n", addr)
}
