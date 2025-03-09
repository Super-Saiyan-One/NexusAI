package service

import (
	"errors"
	userGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type UserGroupService interface {
	UserGroupCreate(repo repository.UserGroupRepository, userGroup *dto.UserGroup) (*dto.UserGroup, error)
	UserGroupUpdate(repo repository.UserGroupRepository, userGroup *dto.UserGroup) (*dto.UserGroup, error)
	UserGroupSearch(repo repository.UserGroupRepository, userGroupSearch *userGroupDto.UserGroupSearchRequest) ([]*dto.UserGroup, error)
	UserGroupDelete(repo repository.UserGroupRepository, userGroupID string) error
	UserGroupAvailableModels(repo repository.UserGroupRepository, userGroupID string) ([]*dto.Model, error)
	UserGroupAvailableChannels(repo repository.UserGroupRepository, userGroupID string) ([]*dto.Channel, error)
}

type userGroupService struct{}

func NewUserGroupService() UserGroupService {
	return &userGroupService{}
}

// UserGroupCreate 创建用户组
func (ugs *userGroupService) UserGroupCreate(repo repository.UserGroupRepository, userGroup *dto.UserGroup) (*dto.UserGroup, error) {
	return repo.Create(userGroup)
}

// UserGroupUpdate 更新用户组
func (ugs *userGroupService) UserGroupUpdate(repo repository.UserGroupRepository, userGroup *dto.UserGroup) (*dto.UserGroup, error) {
	existingGroup, err := repo.GetByID(userGroup.UserGroupID)
	if err != nil {
		return nil, errors.New("user group not found")
	}
	if existingGroup.DeletedAt != nil {
		return nil, errors.New("user group already deleted")
	}

	return repo.Update(userGroup)
}

// UserGroupSearch 搜索用户组
func (ugs *userGroupService) UserGroupSearch(repo repository.UserGroupRepository, userGroupSearch *userGroupDto.UserGroupSearchRequest) ([]*dto.UserGroup, error) {
	groups, _, err := repo.Search(userGroupSearch)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// UserGroupDelete 删除用户组
func (ugs *userGroupService) UserGroupDelete(repo repository.UserGroupRepository, userGroupID string) error {
	existingGroup, err := repo.GetByID(userGroupID)
	if err != nil {
		return errors.New("user group not found")
	}
	if existingGroup.DeletedAt != nil {
		return errors.New("user group already deleted")
	}
	return repo.Delete(userGroupID)
}

// UserGroupAvailableModels 获取用户组模型
// user_group.extra_allowed_models 用户组额外允许的模型列表
// user_group.priority_models 用户组优先使用的模型列表
func (ugs *userGroupService) UserGroupAvailableModels(repo repository.UserGroupRepository, userGroupID string) ([]*dto.Model, error) {
	// TODO: 获取用户组模型
	return []*dto.Model{}, nil
}

// UserGroupAvailableChannels 获取用户组渠道
func (ugs *userGroupService) UserGroupAvailableChannels(repo repository.UserGroupRepository, userGroupID string) ([]*dto.Channel, error) {
	// TODO: 获取用户组渠道
	return []*dto.Channel{}, nil
}
