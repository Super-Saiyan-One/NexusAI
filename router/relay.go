package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRelayRouter(server *gin.Engine) {
	relayV1Controller := controller.NewRelayV1Controller()
	relayVideoController := controller.NewRelayVideoController()
	relayImageController := controller.NewRelayImageController()
	relayRouter := server.Group("/relay")
	// 验证令牌是否有效
	// 验证模型是否可用
	// 根据模型和请求内容验证配额是否足够
	relayRouter.Use(middleware.TokenVerifyMiddleware(), middleware.ModelVerifyMiddleware(), middleware.QuotaVerifyMiddleware())
	{
		relayV1Router := relayRouter.Group("/v1")
		{
			relayV1Router.POST("/chat/completions", relayV1Controller.RelayChatCompletions)
			relayVideoRouter := relayV1Router.Group("/videos")
			{
				relayVideoRouter.POST("/text2video", relayVideoController.Text2Video)
				relayVideoRouter.GET("/text2video", relayVideoController.Text2Video)
				relayVideoRouter.GET("/text2video/:task_id", relayVideoController.Text2Video)

				relayVideoRouter.POST("/image2video", relayVideoController.Image2Video)
				relayVideoRouter.GET("/image2video", relayVideoController.Image2Video)
				relayVideoRouter.GET("/image2video/:task_id", relayVideoController.Image2Video)

				relayVideoRouter.POST("/video-extend", relayVideoController.VideoExtend)
				relayVideoRouter.GET("/video-extend", relayVideoController.VideoExtend)
				relayVideoRouter.GET("/video-extend/:task_id", relayVideoController.VideoExtend)

				relayVideoRouter.POST("/lip-sync", relayVideoController.LipSync)
				relayVideoRouter.GET("/lip-sync", relayVideoController.LipSync)
				relayVideoRouter.GET("/lip-sync/:task_id", relayVideoController.LipSync)
			}
		}
		relayImageRouter := relayV1Router.Group("/images")
		{
			relayImageRouter.POST("/generations", relayImageController.RelayImage)
			relayImageRouter.GET("/generations", relayImageController.RelayImage)
			relayImageRouter.GET("/generations/:task_id", relayImageController.RelayImage)
		}
	}
}
