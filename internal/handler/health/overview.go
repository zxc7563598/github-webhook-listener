package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/dto"
	healthDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/health"
)

func (h *Handler) GetOverview(ctx *gin.Context) {
	result, err := h.svc.GetOverview(h.cfg.Repos)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Error(http.StatusInternalServerError, "数据获取失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, dto.Success(http.StatusOK, result))
}

func (h *Handler) GetWebhookLogDetails(ctx *gin.Context) {
	var req healthDTO.GetWebhookLogDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Error(http.StatusBadRequest, "请求参数异常", nil))
		return
	}
	result, err := h.svc.GetWebhookLogDetails(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Error(http.StatusInternalServerError, "数据获取失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, dto.Success(http.StatusOK, result))
}
