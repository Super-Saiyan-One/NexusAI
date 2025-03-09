package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChannelGroupRouter(server *gin.Engine) {
	// 初始化controller
	channelGroupController := controller.NewChannelGroupController()

	channelGroupRouter := server.Group("/channel_group")
	// 需要认证且是root用户的路由
	channelGroupRouter.Use(middleware.UserVerifyMiddleware(), middleware.RootVerifyMiddleware())
	{
		channelGroupRouter.POST("/create", channelGroupController.ChannelGroupCreate)
		channelGroupRouter.POST("/update", channelGroupController.ChannelGroupUpdate)
		channelGroupRouter.GET("/search", channelGroupController.ChannelGroupSearch)
		channelGroupRouter.DELETE("/delete/:channel_group_id", channelGroupController.ChannelGroupDelete)
		channelGroupRouter.GET("/available_models/:channel_group_id", channelGroupController.ChannelGroupAvailableModels)
	}
}
