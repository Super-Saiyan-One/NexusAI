package middleware

import (
	"net/http"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func RootVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.MustGet(string(constant.UserKey)).(*dto.User)
		if !ok {
			utils.AbortWhenRootVerifyFailed(c, http.StatusUnauthorized, "Current user has no root permission")
			return
		}
		userGroupID := user.UserGroupID
		if userGroupID == "" {
			utils.AbortWhenRootVerifyFailed(c, http.StatusUnauthorized, "Current user has no root permission")
			return
		}
		userGroupRepo := repository.NewUserGroupRepository(model.GetDB())
		userGroup, err := userGroupRepo.GetByID(userGroupID)
		if err != nil { // 如果获取用户组失败，则认为没有root权限
			utils.AbortWhenRootVerifyFailed(c, http.StatusUnauthorized, "Current user has no root permission")
			return
		}
		if userGroup.UserGroupOptions.DefaultLevel != 99 && user.UserOptions.DefaultLevel != 99 { // 如果用户组和用户默认等级都不是99，则认为没有root权限
			utils.AbortWhenRootVerifyFailed(c, http.StatusUnauthorized, "Current user has no root permission")
			return
		}
		c.Next()
	}
}
