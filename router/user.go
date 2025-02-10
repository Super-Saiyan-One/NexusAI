package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(server *gin.Engine) {
	// 初始化controller
	userController := controller.NewUserController()

	// 设置路由
	userRouter := server.Group("/user")
	{
		// 无需认证的路由
		userRouter.POST("/login", userController.UserLogin)
		userRouter.POST("/register", userController.UserRegister)
		userRouter.GET("/search", userController.UserSearch)
		userRouter.DELETE("/delete/:user_id", userController.UserDelete)

		// 需要认证的路由
		authRouter := userRouter.Group("/:user_id")
		authRouter.Use(middleware.UserVerifyMiddleware())
		{
			authRouter.POST("/logout", userController.UserLogout)
			authRouter.POST("/update", userController.UserUpdate)
			authRouter.POST("/password", userController.UserPassword)
		}
	}
}
