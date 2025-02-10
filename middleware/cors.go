package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	// 1. 允许所有来源
	corsConfig.AllowAllOrigins = true
	// 2. 允许携带凭证
	corsConfig.AllowCredentials = true
	// 3. 允许所有方法
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// 4. 允许所有头
	corsConfig.AllowHeaders = []string{"*"}
	return cors.New(corsConfig)
}
