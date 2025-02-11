package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChannelRouter(server *gin.Engine) {
	// 初始化controller
	channelController := controller.NewChannelController()

	channelRouter := server.Group("/channel")
	// 需要认证且是root用户的路由
	channelRouter.Use(middleware.UserVerifyMiddleware(), middleware.RootVerifyMiddleware())
	{
		channelRouter.POST("/create", channelController.ChannelCreate)
		channelRouter.POST("/update", channelController.ChannelUpdate)
		channelRouter.GET("/search", channelController.ChannelSearch)
		channelRouter.DELETE("/delete/:channel_id", channelController.ChannelDelete)
	}
}
