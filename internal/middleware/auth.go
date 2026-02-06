package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebBasicAuth
// 当 username 和 password 都非空时，启用 Basic Auth
// 否则，直接放行（相当于 no-op middleware）
func WebBasicAuth(username, password *string) gin.HandlerFunc {
	// 如果未配置账号或密码，直接返回一个空中间件
	if username == nil || password == nil || *username == "" || *password == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		u, p, ok := c.Request.BasicAuth()
		if !ok || u != *username || p != *password {
			c.Header("WWW-Authenticate", `Basic realm="Web UI"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
