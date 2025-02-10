package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupTokenRouter(server *gin.Engine) {
	// 初始化controller
	tokenController := controller.NewTokenController()

	// 设置路由
	tokenRouter := server.Group("/token")
	tokenRouter.Use(middleware.UserVerifyMiddleware()) // 用户验证中间件
	{
		tokenRouter.POST("/create", tokenController.TokenCreate)
		tokenRouter.POST("/update", tokenController.TokenUpdate)
		tokenRouter.GET("/search", tokenController.TokenSearch)
		tokenRouter.DELETE("/delete/:token_id", tokenController.TokenDelete)
	}
}
