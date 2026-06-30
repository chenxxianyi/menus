package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type BrowseHistoryHandler struct {
	historyService *service.BrowseHistoryService
}

func NewBrowseHistoryHandler(historyService *service.BrowseHistoryService) *BrowseHistoryHandler {
	return &BrowseHistoryHandler{historyService: historyService}
}

func (h *BrowseHistoryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	histories, total, err := h.historyService.List(userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}

	response.SuccessPage(c, histories, total, page, pageSize)
}

func (h *BrowseHistoryHandler) Clear(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.historyService.Clear(userID); err != nil {
		response.Error(c, errcode.ErrServer, "清空失败")
		return
	}

	response.Success(c, nil)
}
