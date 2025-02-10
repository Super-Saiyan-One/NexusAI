package middleware

import (
	"context"
	"nexus-ai/constant"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func TestRootMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从数据库中获取 root 用户信息
		userRepo := repository.NewUserRepository(model.GetDB())
		user, err := userRepo.GetByUsername(constant.RootUserName)
		if err != nil {
			utils.SysError("获取root用户信息失败: " + err.Error())
			c.AbortWithStatusJSON(500, gin.H{"error": "系统错误"})
			return
		}

		// 1. 设置UserID到gin.Context中
		c.Set(string(constant.UserIDKey), user.UserID)
		// 2. 设置UserID到gin.Context的Request.Context中
		ctx := context.WithValue(c.Request.Context(), constant.UserIDKey, user.UserID)
		c.Request = c.Request.WithContext(ctx)
		// 3. 设置UserID到gin.Context的Header中
		c.Header(string(constant.UserIDKey), user.UserID)
		// 4. 继续处理请求
		c.Next()
	}
}
