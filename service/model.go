package service

import (
	"errors"
	modelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type ModelService interface {
	ModelCreate(repo repository.ModelRepository, model *dto.Model) (*dto.Model, error)
	ModelUpdate(repo repository.ModelRepository, model *dto.Model) (*dto.Model, error)
	ModelSearch(repo repository.ModelRepository, modelSearch *modelDto.ModelSearchRequest) ([]*dto.Model, error)
	ModelDelete(repo repository.ModelRepository, modelID string) error
	ModelAvailableChannels(repo repository.ModelRepository, modelID string) ([]*dto.Channel, error)
}

type modelService struct{}

func NewModelService() ModelService {
	return &modelService{}
}

// ModelCreate 创建模型
func (ms *modelService) ModelCreate(repo repository.ModelRepository, model *dto.Model) (*dto.Model, error) {
	existingModel, _ := repo.GetByName(model.ModelName)
	if existingModel != nil && existingModel.DeletedAt == nil {
		return nil, errors.New("model already exists")
	}
	return repo.Create(model)
}

// ModelUpdate 更新模型
func (ms *modelService) ModelUpdate(repo repository.ModelRepository, model *dto.Model) (*dto.Model, error) {
	existingModel, err := repo.GetByID(model.ModelID)
	if err != nil {
		return nil, errors.New("model not found")
	}
	if existingModel.DeletedAt != nil {
		return nil, errors.New("model already deleted")
	}
	return repo.Update(model)
}

// ModelSearch 搜索模型
func (ms *modelService) ModelSearch(repo repository.ModelRepository, modelSearch *modelDto.ModelSearchRequest) ([]*dto.Model, error) {
	models, _, err := repo.Search(modelSearch)
	if err != nil {
		return nil, err
	}
	return models, nil
}

// ModelDelete 删除模型
func (ms *modelService) ModelDelete(repo repository.ModelRepository, modelID string) error {
	existingModel, err := repo.GetByID(modelID)
	if err != nil {
		return errors.New("model not found")
	}
	if existingModel.DeletedAt != nil {
		return errors.New("model already deleted")
	}
	return repo.Delete(modelID)
}

// ModelAvailableChannels 获取模型可用渠道
func (ms *modelService) ModelAvailableChannels(repo repository.ModelRepository, modelID string) ([]*dto.Channel, error) {
	// TODO: 获取模型可用渠道
	return []*dto.Channel{}, nil
}
