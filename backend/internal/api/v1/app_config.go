package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type AppConfigHandler struct {
	appConfigService *service.AppConfigService
}

func NewAppConfigHandler(appConfigService *service.AppConfigService) *AppConfigHandler {
	return &AppConfigHandler{appConfigService: appConfigService}
}

func (h *AppConfigHandler) About(c *gin.Context) {
	info, err := h.appConfigService.GetAboutInfo()
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, info)
}
