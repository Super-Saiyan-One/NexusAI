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
