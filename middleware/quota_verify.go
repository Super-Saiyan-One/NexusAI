package middleware

import (
	"net/http"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func QuotaVerifyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet(string(constant.TokenKey)).(*dto.Token)      // 从Gin上下文中获取令牌信息
		tokenQuotaLeft := token.TokenQuotaLeft                          // 获取令牌剩余配额
		userID := token.UserID                                          // 获取令牌所属用户ID
		user, ok, err := repository.QuotaVerify(userID, tokenQuotaLeft) // 验证配额
		if !ok {
			utils.AbortWhenQuotaVerifyFailed(c, http.StatusUnauthorized, err.Error())
			return
		}
		c.Set(string(constant.UserKey), user)
		c.Next()
	}
}
