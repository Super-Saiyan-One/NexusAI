package utils

import (
	"fmt"
	"nexus-ai/constant"

	"github.com/gin-gonic/gin"
)

func SetupAPILog(server *gin.Engine) {
	colorBlue := "\033[36m"
	colorReset := "\033[0m"
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var requestID string
		if param.Keys != nil {
			requestID = param.Keys[string(constant.RequestIDKey)].(string)
		}
		return fmt.Sprintf("%s[API]%s     %s | %s | %d | %12v | %15s | %5s | %s%s%s\n",
			colorBlue,
			colorReset,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			requestID,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			colorBlue,
			param.Path,
			colorReset,
		)
	}))
}
