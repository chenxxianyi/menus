package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/model"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type FeedbackHandler struct {
	feedbackService *service.FeedbackService
}

func NewFeedbackHandler(feedbackService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService}
}

type CreateFeedbackRequest struct {
	Content string `json:"content" binding:"required"`
	Contact string `json:"contact"`
}

func (h *FeedbackHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	fb := &model.Feedback{
		UserID:  &userID,
		Content: req.Content,
		Contact: req.Contact,
	}

	if err := h.feedbackService.Create(fb); err != nil {
		response.Error(c, errcode.ErrServer, "提交失败")
		return
	}
	response.Success(c, nil)
}
