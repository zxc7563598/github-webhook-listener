package webhook

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/dto"
	webhookDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/webhook"
	"github.com/zxc7563598/github-webhook-listener/pkg/utils"
)

func (h *Handler) MakeWebhookHandler(ctx *gin.Context) {
	// 获取原始数据
	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Error(http.StatusBadRequest, "请求数据异常", nil))
		return
	}
	// 解析数据
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var req webhookDTO.GitHubWebhook
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Error(http.StatusBadRequest, "请求参数异常", nil))
		return
	}
	// 获取关键数据
	repoName := req.Repository.FullName
	if repoName == "" {
		log.Printf("[webhook] 无法从请求中提取到仓库名")
		ctx.JSON(http.StatusOK, dto.Error(http.StatusUnprocessableEntity, "无法获取仓库名", nil))
		return
	}
	ref := req.Ref
	if ref == "" {
		log.Printf("[webhook] 无法从 %s 仓库的请求中提取到分支信息", repoName)
		ctx.JSON(http.StatusOK, dto.Error(http.StatusUnprocessableEntity, "无法获取分支信息", nil))
		return
	}
	parts := strings.Split(ref, "/")
	branch := parts[len(parts)-1]
	if branch == "" {
		log.Printf("[webhook] 无法从 %s 仓库的分支信息携带的 %s 中提取到分支信息", repoName, ref)
		ctx.JSON(http.StatusOK, dto.Error(http.StatusUnprocessableEntity, "无法获取分支信息", nil))
		return
	}
	// 查找配置
	repoCfg, ok := h.cfg.Repos[repoName]
	if !ok {
		log.Printf("[webhook] 在配置中找不到存储库: %s", repoName)
		ctx.JSON(http.StatusOK, dto.Error(http.StatusNotFound, "未能找到存储库", nil))
		return
	}
	signature := ctx.GetHeader("X-Hub-Signature-256")
	if !utils.ValidateGitHubSignature(repoCfg.Secret, body, signature) {
		log.Printf("[webhook] 仓库 %s 的GitHub签名验证失败", repoName)
		ctx.JSON(http.StatusOK, dto.Error(http.StatusServiceUnavailable, "签名验证失败", nil))
		return
	}
	event := ctx.GetHeader("X-GitHub-Event")
	if event == "" {
		log.Printf("[webhook] 缺少 X-GitHub-Event 头")
		ctx.JSON(http.StatusOK, dto.Error(http.StatusNotFound, "缺少 X-GitHub-Event", nil))
		return
	}
	// 匹配规则
	if err := h.svc.MakeWebhookService(repoCfg, branch, event, repoName); err != nil {
		log.Printf("[webhook] 匹配规则执行shell失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Error(http.StatusInternalServerError, "匹配规则失败", nil))
		return
	}
	// 返回成功
	ctx.JSON(http.StatusOK, dto.Success(http.StatusOK, nil))
}
