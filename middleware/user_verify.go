package middleware

import (
	"context"
	"errors"
	"net/http"
	"nexus-ai/constant"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(string(constant.AccessTokenKey))
		accessToken = strings.TrimPrefix(accessToken, "Bearer ") // 去掉Bearer
		refreshToken := c.GetHeader(string(constant.RefreshTokenKey))
		refreshToken = strings.TrimPrefix(refreshToken, "Bearer ") // 去掉Bearer
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

				// 使用新的access token继续验证
				_, err = service.ValidateToken(tokenPair.AccessToken, userID)
				if err != nil { // 如果新的access token验证失败
					utils.AbortWhenUserVerifyFailed(c, http.StatusUnauthorized, "login expired, please login again")
					return
				}
				// 更新access token和refresh token
				accessToken = tokenPair.AccessToken
				refreshToken = tokenPair.RefreshToken
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
		// 设置新的token到响应头
		c.Header(string(constant.AccessTokenKey), accessToken)
		c.Header(string(constant.RefreshTokenKey), refreshToken)
		c.Header(string(constant.UserIDKey), user.UserID)

		// 设置新的token到gin.Context和gin.Context.Request.Context中
		c.Set(string(constant.AccessTokenKey), accessToken)
		ctx := context.WithValue(c.Request.Context(), constant.AccessTokenKey, accessToken)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(constant.RefreshTokenKey), refreshToken)
		ctx = context.WithValue(c.Request.Context(), constant.RefreshTokenKey, refreshToken)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(constant.UserIDKey), user.UserID)
		ctx = context.WithValue(c.Request.Context(), constant.UserIDKey, user.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(constant.UserKey), user)
		ctx = context.WithValue(c.Request.Context(), constant.UserKey, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
