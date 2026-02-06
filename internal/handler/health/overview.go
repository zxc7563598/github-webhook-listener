package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/dto"
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
	result, err := h.svc.GetWebhookLogDetails("1")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Error(http.StatusInternalServerError, "数据获取失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, dto.Success(http.StatusOK, result))
}
