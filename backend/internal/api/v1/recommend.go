package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"menu-recommend/internal/middleware"
	"menu-recommend/internal/pkg/errcode"
	"menu-recommend/internal/pkg/response"
	"menu-recommend/internal/service"
)

type RecommendHandler struct {
	recommendService *service.RecommendService
}

func NewRecommendHandler(recommendService *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{recommendService: recommendService}
}

func (h *RecommendHandler) Menu(c *gin.Context) {
	var params service.RecommendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	params.UserID = middleware.GetUserID(c)

	result, err := h.recommendService.RecommendMenu(&params)
	if err != nil {
		response.Error(c, errcode.ErrServer, "推荐失败")
		return
	}
	response.Success(c, result)
}

func (h *RecommendHandler) MenuAI(c *gin.Context) {
	var params service.RecommendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	params.UserID = userID
	result, err := h.recommendService.RecommendSceneByAI(c.Request.Context(), userID, &params)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAIConfigMissing):
			response.Error(c, errcode.ErrServer, "AI 未配置，请先配置 AI 服务后再生成")
		case errors.Is(err, service.ErrAIInvalidResponse):
			response.Error(c, errcode.ErrServer, "AI 返回的场景推荐格式无效，请稍后重试")
		case errors.Is(err, service.ErrAIUpstream):
			response.Error(c, errcode.ErrServer, "AI 服务暂时不可用或返回错误，请稍后重试")
		default:
			response.Error(c, errcode.ErrServer, "AI 场景推荐失败")
		}
		return
	}
	response.Success(c, result)
}

type ByIngredientsRequest struct {
	Ingredients []string `json:"ingredients" binding:"required"`
	Mode        string   `json:"mode"`
	Limit       int      `json:"limit"`
}

func (h *RecommendHandler) ByIngredients(c *gin.Context) {
	var req ByIngredientsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}

	if len(req.Ingredients) > 20 {
		response.Error(c, errcode.ErrParam, "食材数量不能超过20个")
		return
	}
	if req.Mode != "" && req.Mode != "ingredients" && req.Mode != "fridge" {
		response.Error(c, errcode.ErrParam, "推荐模式错误")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.recommendService.RecommendByIngredients(userID, req.Ingredients, req.Mode, req.Limit)
	if err != nil {
		response.Error(c, errcode.ErrServer, "推荐失败")
		return
	}
	response.Success(c, result)
}

func (h *RecommendHandler) WeekMenu(c *gin.Context) {
	var params service.RecommendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, errcode.ErrParam, "参数错误")
		return
	}
	params.UserID = middleware.GetUserID(c)

	result, err := h.recommendService.GenerateWeekMenu(&params)
	if err != nil {
		response.Error(c, errcode.ErrServer, "生成失败")
		return
	}
	response.Success(c, result)
}
