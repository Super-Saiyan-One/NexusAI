package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserGroupRouter(server *gin.Engine) {
	// 初始化controller
	userGroupController := controller.NewUserGroupController()

	userGroupRouter := server.Group("/user_group")
	// 需要认证且是root用户的路由
	userGroupRouter.Use(middleware.UserVerifyMiddleware(), middleware.RootVerifyMiddleware())
	{
		userGroupRouter.POST("/create", userGroupController.UserGroupCreate)
		userGroupRouter.POST("/update", userGroupController.UserGroupUpdate)
		userGroupRouter.GET("/search", userGroupController.UserGroupSearch)
		userGroupRouter.DELETE("/delete/:user_group_id", userGroupController.UserGroupDelete)
		userGroupRouter.GET("/available_models/:user_group_id", userGroupController.UserGroupAvailableModels)
		userGroupRouter.GET("/available_channels/:user_group_id", userGroupController.UserGroupAvailableChannels)
	}
}
