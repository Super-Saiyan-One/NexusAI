package repository

import (
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
)

type AvailableChannel struct {
	ChannelID         string
	ChannelGroupLevel int
}

func GetAvailableChannels(requestModel *dto.Model, token *dto.Token) (map[string]int, []AvailableChannel) {
	availableLevels := token.TokenOptions.AvailableLevels // 令牌可用等级
	channelGroupRepo := NewChannelGroupRepository(model.GetDB())
	if len(availableLevels) > 0 {
		availableChannels := make(map[string]int)        // 初始化map
		for _, availableLevel := range availableLevels { // 遍历token可用等级
			channelGroups, err := channelGroupRepo.GetListByDefaultLevel(availableLevel) // 获取对应等级的渠道组列表
			if err != nil {
				continue
			}
			for _, channelGroup := range channelGroups { // 遍历当前等级的渠道组
				modelAvailableChannels := channelGroup.ChannelGroupChannels.ModelsMap[requestModel.ModelID] // 当前渠道组请求模型可用的渠道列表
				// 将当前渠道组请求模型可用的渠道列表逐个添加到可用渠道列表
				for _, channelID := range modelAvailableChannels {
					availableChannels[channelID] = availableLevel
				}
			}
		}
		// 将map转换为结构体切片
		result := make([]AvailableChannel, 0, len(availableChannels))
		for channelID, level := range availableChannels {
			result = append(result, AvailableChannel{
				ChannelID:         channelID,
				ChannelGroupLevel: level,
			})
		}
		return availableChannels, result
	}
	return nil, nil // 令牌可用等级为空，返回空列表
}

func GetDefaultLevelByChannelID(channelID string) (int, error) {
	channelRepo := NewChannelRepository(model.GetDB())
	channel, err := channelRepo.GetByID(channelID)
	if err != nil {
		return -1, err
	}
	channelGroupRepo := NewChannelGroupRepository(model.GetDB())
	channelGroup, err := channelGroupRepo.GetByID(channel.ChannelGroupID)
	if err != nil {
		return -1, err
	}
	return channelGroup.ChannelGroupOptions.DefaultLevel, nil
}

func SelectChannel(availableChannels []AvailableChannel) (*dto.Channel, error) {
	channelRepo := NewChannelRepository(model.GetDB())
	channel, err := channelRepo.GetByID(availableChannels[0].ChannelID)
	if err != nil {
		return nil, err
	}
	return channel, nil
}
