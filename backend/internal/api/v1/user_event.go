package v1

import (
	"github.com/gin-gonic/gin"

	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type UserEventHandler struct {
	eventService *service.UserEventService
}

func NewUserEventHandler(eventService *service.UserEventService) *UserEventHandler {
	return &UserEventHandler{eventService: eventService}
}

type TrackUserEventRequest struct {
	EventName  string                 `json:"event_name" binding:"required"`
	EntityType string                 `json:"entity_type"`
	EntityID   uint                   `json:"entity_id"`
	Payload    map[string]interface{} `json:"payload"`
}

func (h *UserEventHandler) Track(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req TrackUserEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "埋点参数错误")
		return
	}

	if err := h.eventService.Track(userID, req.EventName, req.EntityType, req.EntityID, req.Payload); err != nil {
		response.Error(c, errcode.ErrServer, "记录行为失败")
		return
	}
	response.Success(c, nil)
}
