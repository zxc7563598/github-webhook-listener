package router

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/handler"
	"github.com/zxc7563598/github-webhook-listener/internal/middleware"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
	"github.com/zxc7563598/github-webhook-listener/internal/service"
	"github.com/zxc7563598/github-webhook-listener/internal/webui"
)

func SetupRouter(isWeb bool, cfg config.Config, scheduler *queue.ShellScheduler, webUser, webPass string) *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	// 中间件注册
	r.Use(gin.Logger(), gin.Recovery())
	// 依赖注入
	shellQueueService := service.NewContainer(scheduler)
	container := handler.NewContainer(isWeb, cfg, shellQueueService)
	// 路由
	r.GET("/", container.Lead)
	r.POST("/webhook", container.MakeWebhookHandler)
	if isWeb {
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
