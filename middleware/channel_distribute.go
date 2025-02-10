package middleware

import (
	"net/http"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func ChannelDistributeMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet(string(constant.TokenKey)).(*dto.Token)
		requestModel := c.MustGet(string(constant.ModelKey)).(*dto.Model)
		tokenExtraAllowedChannels := token.TokenChannels.ExtraAllowedChannels // 令牌额外允许的渠道列表
		tokenPriorityChannels := token.TokenChannels.PriorityChannels         // 令牌优先使用的渠道列表
		availableChannelsMap, availableChannels := repository.GetAvailableChannels(requestModel, token)
		if len(tokenExtraAllowedChannels) > 0 { // 如果令牌额外允许的渠道列表不为空 与可用渠道列表合并
			for _, extraAllowedChannelID := range tokenExtraAllowedChannels {
				extraAllowedChannelLevel, err := repository.GetDefaultLevelByChannelID(extraAllowedChannelID)
				if err != nil {
					continue
				}
				if _, ok := availableChannelsMap[extraAllowedChannelID]; !ok { // 如果当前渠道不在可用渠道列表中
					availableChannelsMap[extraAllowedChannelID] = extraAllowedChannelLevel
					availableChannels = append(availableChannels, repository.AvailableChannel{
						ChannelID:         extraAllowedChannelID,
						ChannelGroupLevel: extraAllowedChannelLevel,
					})
				}
			}
		}
		if len(tokenPriorityChannels) > 0 { // 如果令牌优先使用的渠道列表不为空 从可用渠道中获取优先使用的渠道
			priorityChannels := make([]repository.AvailableChannel, 0)
			for _, priorityChannelID := range tokenPriorityChannels {
				if _, ok := availableChannelsMap[priorityChannelID]; ok {
					priorityChannels = append(priorityChannels, repository.AvailableChannel{
						ChannelID:         priorityChannelID,
						ChannelGroupLevel: availableChannelsMap[priorityChannelID],
					})
				}
			}
			availableChannels = priorityChannels
		}
		selectedChannel, err := repository.SelectChannel(availableChannels)
		if err != nil {
			utils.AbortWhenChannelDistributeFailed(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.Set(string(constant.ChannelKey), selectedChannel) // 设置选中的渠道
		c.Next()
	}
}
