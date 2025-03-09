package service

import (
	"errors"
	channelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type ChannelGroupService interface {
	ChannelGroupCreate(repo repository.ChannelGroupRepository, channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error)
	ChannelGroupUpdate(repo repository.ChannelGroupRepository, channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error)
	ChannelGroupSearch(repo repository.ChannelGroupRepository, channelGroupSearch *channelGroupDto.ChannelGroupSearchRequest) ([]*dto.ChannelGroup, error)
	ChannelGroupDelete(repo repository.ChannelGroupRepository, channelGroupID string) error
	ChannelGroupAvailableModels(repo repository.ChannelGroupRepository, channelGroupID string) ([]*dto.Model, error)
}

type channelGroupService struct{}

func NewChannelGroupService() ChannelGroupService {
	return &channelGroupService{}
}

// ChannelGroupCreate 创建渠道组
func (cgs *channelGroupService) ChannelGroupCreate(repo repository.ChannelGroupRepository, channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error) {
	existingGroup, _ := repo.GetByName(channelGroup.ChannelGroupName)
	if existingGroup != nil && existingGroup.DeletedAt == nil {
		return nil, errors.New("channel group already exists")
	}
	return repo.Create(channelGroup)
}

// ChannelGroupUpdate 更新渠道组
func (cgs *channelGroupService) ChannelGroupUpdate(repo repository.ChannelGroupRepository, channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error) {
	existingGroup, err := repo.GetByID(channelGroup.ChannelGroupID)
	if err != nil {
		return nil, errors.New("channel group not found")
	}
	if existingGroup.DeletedAt != nil {
		return nil, errors.New("channel group already deleted")
	}
	return repo.Update(channelGroup)
}

// ChannelGroupSearch 搜索渠道组
func (cgs *channelGroupService) ChannelGroupSearch(repo repository.ChannelGroupRepository, channelGroupSearch *channelGroupDto.ChannelGroupSearchRequest) ([]*dto.ChannelGroup, error) {
	groups, _, err := repo.Search(channelGroupSearch)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// ChannelGroupDelete 删除渠道组
func (cgs *channelGroupService) ChannelGroupDelete(repo repository.ChannelGroupRepository, channelGroupID string) error {
	existingGroup, err := repo.GetByID(channelGroupID)
	if err != nil {
		return errors.New("channel group not found")
	}
	if existingGroup.DeletedAt != nil {
		return errors.New("channel group already deleted")
	}
	return repo.Delete(channelGroupID)
}

// ChannelGroupAvailableModels 获取渠道组可用模型
func (cgs *channelGroupService) ChannelGroupAvailableModels(repo repository.ChannelGroupRepository, channelGroupID string) ([]*dto.Model, error) {
	// TODO: 获取渠道组可用模型
	return []*dto.Model{}, nil
}
