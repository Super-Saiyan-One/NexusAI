package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupModelGroupRouter(server *gin.Engine) {
	// 初始化controller
	modelGroupController := controller.NewModelGroupController()

	modelGroupRouter := server.Group("/model_group")
	// 需要认证且是root用户的路由
	modelGroupRouter.Use(middleware.UserVerifyMiddleware(), middleware.RootVerifyMiddleware())
	{
		modelGroupRouter.POST("/create", modelGroupController.ModelGroupCreate)
		modelGroupRouter.POST("/update", modelGroupController.ModelGroupUpdate)
		modelGroupRouter.GET("/search", modelGroupController.ModelGroupSearch)
		modelGroupRouter.DELETE("/delete/:model_group_id", modelGroupController.ModelGroupDelete)
	}
}
