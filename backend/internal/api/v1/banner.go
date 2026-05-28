package v1

import (
	"github.com/gin-gonic/gin"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type BannerHandler struct {
	bannerService *service.BannerService
}

func NewBannerHandler(bannerService *service.BannerService) *BannerHandler {
	return &BannerHandler{bannerService: bannerService}
}

func (h *BannerHandler) List(c *gin.Context) {
	banners, err := h.bannerService.GetActive()
	if err != nil {
		response.Error(c, errcode.ErrServer, "查询失败")
		return
	}
	response.Success(c, banners)
}
