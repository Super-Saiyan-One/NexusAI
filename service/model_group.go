package service

import (
	"errors"
	modelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type ModelGroupService interface {
	ModelGroupCreate(repo repository.ModelGroupRepository, modelGroup *dto.ModelGroup) (*dto.ModelGroup, error)
	ModelGroupUpdate(repo repository.ModelGroupRepository, modelGroup *dto.ModelGroup) (*dto.ModelGroup, error)
	ModelGroupSearch(repo repository.ModelGroupRepository, modelGroupSearch *modelGroupDto.ModelGroupSearchRequest) ([]*dto.ModelGroup, error)
	ModelGroupDelete(repo repository.ModelGroupRepository, modelGroupID string) error
}

type modelGroupService struct{}

func NewModelGroupService() ModelGroupService {
	return &modelGroupService{}
}

// ModelGroupCreate 创建模型组
func (mgs *modelGroupService) ModelGroupCreate(repo repository.ModelGroupRepository, modelGroup *dto.ModelGroup) (*dto.ModelGroup, error) {
	existingGroup, _ := repo.GetByName(modelGroup.ModelGroupName)
	if existingGroup != nil && existingGroup.DeletedAt == nil {
		return nil, errors.New("model group already exists")
	}
	return repo.Create(modelGroup)
}

// ModelGroupUpdate 更新模型组
func (mgs *modelGroupService) ModelGroupUpdate(repo repository.ModelGroupRepository, modelGroup *dto.ModelGroup) (*dto.ModelGroup, error) {
	existingGroup, err := repo.GetByID(modelGroup.ModelGroupID)
	if err != nil {
		return nil, errors.New("model group not found")
	}
	if existingGroup.DeletedAt != nil {
		return nil, errors.New("model group already deleted")
	}
	return repo.Update(modelGroup)
}

// ModelGroupSearch 搜索模型组
func (mgs *modelGroupService) ModelGroupSearch(repo repository.ModelGroupRepository, modelGroupSearch *modelGroupDto.ModelGroupSearchRequest) ([]*dto.ModelGroup, error) {
	groups, _, err := repo.Search(modelGroupSearch)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// ModelGroupDelete 删除模型组
func (mgs *modelGroupService) ModelGroupDelete(repo repository.ModelGroupRepository, modelGroupID string) error {
	existingGroup, err := repo.GetByID(modelGroupID)
	if err != nil {
		return errors.New("model group not found")
	}
	if existingGroup.DeletedAt != nil {
		return errors.New("model group already deleted")
	}
	return repo.Delete(modelGroupID)
}
