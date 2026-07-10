package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{}
	ready := true
	if h.db == nil {
		checks["database"] = "not configured"
		ready = false
	} else if sqlDB, err := h.db.DB(); err != nil {
		checks["database"] = "unavailable"
		ready = false
	} else if err := sqlDB.PingContext(ctx); err != nil {
		checks["database"] = "unavailable"
		ready = false
	} else {
		checks["database"] = "ok"
	}

	if h.rdb == nil {
		checks["redis"] = "not configured"
		ready = false
	} else if err := h.rdb.Ping(ctx).Err(); err != nil {
		checks["redis"] = "unavailable"
		ready = false
	} else {
		checks["redis"] = "ok"
	}

	if !ready {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unready", "checks": checks})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready", "checks": checks})
}
