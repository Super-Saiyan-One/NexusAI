package router

import (
	"net/http"
	"nexus-ai/middleware"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware()) // 跨域中间件
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is health check!"})
	})
	router.GET("/info", func(c *gin.Context) {
		ctx := c.Request.Context() // 获取请求上下文
		utils.LogInfo(ctx, "This is an info log.")
		utils.LogWarn(ctx, "This is a warning log.")
		utils.LogError(ctx, "This is an error log.")
		utils.SysInfo("This is a system info log.")
		utils.SysError("This is a system error log.")
		utils.SysWarn("This is a system warning log.")
		// utils.FatalLog("This is a fatal log.")
		c.JSON(http.StatusOK, gin.H{"message": "Hello, This is test router!"})
	})
	SetupBillingRouter(router)
	SetupChannelGroupRouter(router)
	SetupChannelRouter(router)
	SetupConfigRouter(router)
	SetupLogRouter(router)
	SetupMessageRouter(router)
	SetupModelGroupRouter(router)
	SetupModelRouter(router)
	SetupMonitorRouter(router)
	SetupNotifyRouter(router)
	SetupPaymentRouter(router)
	SetupQuotaRouter(router)
	SetupRelayRouter(router)
	SetupTokenRouter(router)
	SetupUsageRouter(router)
	SetupUserGroupRouter(router)
	SetupUserRouter(router)
}
