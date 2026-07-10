package v1

import (
	"github.com/gin-gonic/gin"
	"strings"

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

var allowedUserEvents = map[string]bool{
	"home_primary_action_clicked": true,
	"recommend_requested":         true,
	"recommend_result_viewed":     true,
	"recommend_result_rejected":   true,
	"recipe_favorited":            true,
	"recipe_added_to_menu":        true,
	"shopping_list_generated":     true,
	"shopping_item_completed":     true,
	"recipe_marked_cooked":        true,
	"recipe_feedback_submitted":   true,
	"ai_generation_started":       true,
	"ai_generation_completed":     true,
	"ai_generation_failed":        true,
	// Compatibility names emitted by versions before the v2 event dictionary.
	"home_action_click":      true,
	"recommend_start":        true,
	"recommend_result_click": true,
	"add_shopping_list":      true,
	"save_menu":              true,
	"couple_order_create":    true,
}

func (h *UserEventHandler) Track(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req TrackUserEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "埋点参数错误")
		return
	}
	req.EventName = strings.TrimSpace(req.EventName)
	if !allowedUserEvents[req.EventName] {
		response.Error(c, errcode.ErrParam, "不支持的埋点事件")
		return
	}
	if len(req.Payload) > 24 {
		response.Error(c, errcode.ErrParam, "埋点属性过多")
		return
	}

	if err := h.eventService.Track(userID, req.EventName, req.EntityType, req.EntityID, req.Payload); err != nil {
		response.Error(c, errcode.ErrServer, "记录行为失败")
		return
	}
	response.Success(c, nil)
}
