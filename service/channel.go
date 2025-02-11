package service

import (
	"errors"
	channelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type ChannelService interface {
	ChannelCreate(repo repository.ChannelRepository, channel *dto.Channel) (*dto.Channel, error)
	ChannelUpdate(repo repository.ChannelRepository, channel *dto.Channel) (*dto.Channel, error)
	ChannelSearch(repo repository.ChannelRepository, channelSearch *channelDto.ChannelSearchRequest) ([]*dto.Channel, error)
	ChannelDelete(repo repository.ChannelRepository, channelID string) error
}

type channelService struct{}

func NewChannelService() ChannelService {
	return &channelService{}
}

// ChannelCreate 创建渠道
func (cs *channelService) ChannelCreate(repo repository.ChannelRepository, channel *dto.Channel) (*dto.Channel, error) {
	existingChannel, _ := repo.GetByName(channel.ChannelName)
	if existingChannel != nil && existingChannel.DeletedAt == nil {
		return nil, errors.New("channel already exists")
	}
	return repo.Create(channel)
}

// ChannelUpdate 更新渠道
func (cs *channelService) ChannelUpdate(repo repository.ChannelRepository, channel *dto.Channel) (*dto.Channel, error) {
	existingChannel, err := repo.GetByID(channel.ChannelID)
	if err != nil {
		return nil, errors.New("channel not found")
	}
	if existingChannel.DeletedAt != nil {
		return nil, errors.New("channel already deleted")
	}
	return repo.Update(channel)
}

// ChannelSearch 搜索渠道
func (cs *channelService) ChannelSearch(repo repository.ChannelRepository, channelSearch *channelDto.ChannelSearchRequest) ([]*dto.Channel, error) {
	channels, _, err := repo.Search(channelSearch)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// ChannelDelete 删除渠道
func (cs *channelService) ChannelDelete(repo repository.ChannelRepository, channelID string) error {
	existingChannel, err := repo.GetByID(channelID)
	if err != nil {
		return errors.New("channel not found")
	}
	if existingChannel.DeletedAt != nil {
		return errors.New("channel already deleted")
	}
	return repo.Delete(channelID)
}
