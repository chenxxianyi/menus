package v1

import (
	"time"

	"github.com/gin-gonic/gin"

	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type FoodStatsHandler struct {
	service *service.FoodStatsService
}

func NewFoodStatsHandler(service *service.FoodStatsService) *FoodStatsHandler {
	return &FoodStatsHandler{service: service}
}

func (h *FoodStatsHandler) Get(c *gin.Context) {
	period := c.DefaultQuery("period", "week")
	if period != "week" && period != "month" {
		response.Error(c, errcode.ErrParam, "统计周期仅支持 week 或 month")
		return
	}
	result, err := h.service.Get(middleware.GetUserID(c), period, time.Now())
	if err != nil {
		response.Error(c, errcode.ErrServer, "统计数据读取失败")
		return
	}
	response.Success(c, result)
}
