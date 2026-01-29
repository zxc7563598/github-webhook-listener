package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c Container) Lead(ctx *gin.Context) {
	if c.webEnabled {
		ctx.String(http.StatusOK, "web的参数内容为true")
	} else {
		ctx.String(http.StatusOK, "web的参数内容为false")
	}
}
