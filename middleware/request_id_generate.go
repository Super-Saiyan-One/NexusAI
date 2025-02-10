package middleware

import (
	"context"
	"nexus-ai/constant"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func RequestIDGenerateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetTimeString() + "-" + utils.GenerateRandomString(12)
		// 1. 设置requestID到gin.Context中
		c.Set(string(constant.RequestIDKey), requestID)
		// 2. 设置requestID到gin.Context的Request.Context中
		ctx := context.WithValue(c.Request.Context(), constant.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		// 3. 设置requestID到gin.Context的Header中
		c.Header(string(constant.RequestIDKey), requestID)
		// 4. 继续处理请求
		c.Next()
	}
}
