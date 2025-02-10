package middleware

import (
	"net/http"
	"nexus-ai/constant"
	"nexus-ai/repository"
	"nexus-ai/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenVerifyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		var tokenParts []string
		tokenKey := c.Request.Header.Get("Authorization")  // 获取Authorization头
		tokenKey = strings.TrimPrefix(tokenKey, "Bearer ") // 去掉Bearer
		if tokenKey == "midjourney-proxy" {
			tokenKey = c.Request.Header.Get("mj-api-secret")
			tokenKey = strings.TrimPrefix(tokenKey, "Bearer ")
			tokenKey = strings.TrimPrefix(tokenKey, "sk-")
			tokenParts = strings.Split(tokenKey, "-")
			tokenKey = tokenParts[0]
		} else {
			tokenKey = strings.TrimPrefix(tokenKey, "sk-")
			tokenParts = strings.Split(tokenKey, "-")
			tokenKey = tokenParts[0]
		}
		token, err := repository.TokenVerify(tokenKey) // 验证token是否有效 返回dto.Token
		if err != nil {
			utils.AbortWhenTokenVerifyFailed(c, http.StatusUnauthorized, err.Error())
			return
		}
		clientIP := c.ClientIP()
		allowedIPs := token.TokenOptions.AllowedIPs
		disallowedIPs := token.TokenOptions.DisallowedIPs
		if !utils.IsIPAllowed(clientIP, allowedIPs, disallowedIPs) {
			utils.AbortWhenIPVerifyFailed(c, http.StatusUnauthorized, "IP not allowed")
			return
		}

		c.Set(string(constant.TokenKey), token) // 将令牌信息存储在Gin上下文中
		c.Next()
	}
}
