package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupModelRouter(server *gin.Engine) {
	// 初始化controller
	modelController := controller.NewModelController()

	modelRouter := server.Group("/model")
	// 需要认证且是root用户的路由
	modelRouter.Use(middleware.UserVerifyMiddleware(), middleware.RootVerifyMiddleware())
	{
		modelRouter.POST("/create", modelController.ModelCreate)
		modelRouter.POST("/update", modelController.ModelUpdate)
		modelRouter.GET("/search", modelController.ModelSearch)
		modelRouter.DELETE("/delete/:model_id", modelController.ModelDelete)
	}
}
