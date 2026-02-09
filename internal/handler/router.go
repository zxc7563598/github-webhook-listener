package handler

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/middleware"
	"github.com/zxc7563598/github-webhook-listener/internal/webui"

	healthHandler "github.com/zxc7563598/github-webhook-listener/internal/handler/health"
	webhookHandler "github.com/zxc7563598/github-webhook-listener/internal/handler/webhook"
)

func Register(isWeb *bool, webUser, webPass *string, r *gin.Engine, webhookHandler *webhookHandler.Handler, healthHandler *healthHandler.Handler) *gin.Engine {
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	// 中间件注册
	r.Use(gin.Logger(), gin.Recovery())
	// 路由
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "若已开启web模式, 请访问 /web 页面")
	})
	r.POST("/webhook", webhookHandler.MakeWebhookHandler)
	r.POST("/api/overview", healthHandler.GetOverview)
	r.POST("/api/webhook-log-details", healthHandler.GetWebhookLogDetails)

	if *isWeb {
		web := r.Group("/web")
		web.Use(middleware.WebBasicAuth(webUser, webPass))
		RegisterWeb(web)
	}
	return r
}

func RegisterWeb(web *gin.RouterGroup) {
	sub, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.StripPrefix(
		"/web/",
		http.FileServer(http.FS(sub)),
	)
	// /web → /web/
	web.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})
	// /web/* 交给 FileServer
	web.GET("/*filepath", func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
