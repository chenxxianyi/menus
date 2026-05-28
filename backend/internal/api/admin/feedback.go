package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *FeedbackHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		s8 := int8(v)
		status = &s8
	}

	items, total, err := h.feedbackService.List(status, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

type UpdateFeedbackRequest struct {
	Status int8   `json:"status"`
	Reply  string `json:"reply"`
}

func (h *FeedbackHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req UpdateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if err := h.feedbackService.UpdateStatus(uint(id), req.Status, req.Reply); err != nil {
		response.Error(c, errcode.ErrServer, "更新失败")
		return
	}
	response.Success(c, nil)
}
