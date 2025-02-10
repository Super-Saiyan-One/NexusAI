package middleware

import (
	"errors"
	"net/http"
	"nexus-ai/constant"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func UserVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(string(constant.AccessTokenKey))
		refreshToken := c.GetHeader(string(constant.RefreshTokenKey))
		userID := c.GetHeader(string(constant.UserIDKey))

		// 首先验证access token
		_, err := service.ValidateToken(accessToken, userID)
		if err != nil {
			if errors.Is(err, service.ErrExpiredToken) && refreshToken != "" { // 如果access token过期，尝试使用refresh token
				tokenPair, err := service.RefreshAccessToken(refreshToken, userID)
				if err != nil {
					utils.AbortWhenUserVerifyFailed(c, http.StatusUnauthorized, "login expired, please login again")
					return
				}

				// 设置新的token到响应头
				c.Header(string(constant.AccessTokenKey), tokenPair.AccessToken)
				c.Header(string(constant.RefreshTokenKey), tokenPair.RefreshToken)

				// 使用新的access token继续验证
				_, err = service.ValidateToken(tokenPair.AccessToken, userID)
				if err != nil { // 如果新的access token验证失败
					utils.AbortWhenUserVerifyFailed(c, http.StatusUnauthorized, "login expired, please login again")
					return
				}
			} else { // 如果access token过期，且没有refresh token
				utils.AbortWhenUserVerifyFailed(c, http.StatusUnauthorized, "invalid login, please login again")
				return
			}
		}

		// 验证用户信息
		user, err := repository.UserVerify(accessToken, refreshToken, userID)
		if err != nil {
			utils.AbortWhenUserVerifyFailed(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(string(constant.UserIDKey), user.UserID)
		c.Set(string(constant.UserKey), user)
		c.Next()
	}
}
