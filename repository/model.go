package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ModelRepository 模型仓储接口
type ModelRepository interface {
	Create(model *dto.Model) error
	Update(model *dto.Model) error
	Delete(modelID string) error
	GetByID(modelID string) (*dto.Model, error)
	GetByName(name string) (*dto.Model, error)
	GetByAlias(aliasName string) (*dto.Model, error)
	List(page, pageSize int) ([]*dto.Model, int64, error)
	ListByType(modelType string, page, pageSize int) ([]*dto.Model, int64, error)
	ListByProvider(provider string, page, pageSize int) ([]*dto.Model, int64, error)
	ListByOptions(options dto.ModelOptions, page, pageSize int) ([]*dto.Model, int64, error)
	Benchmark(count int) error
}

type modelRepository struct {
	db *gorm.DB
}

// NewModelRepository 创建模型仓储实例
func NewModelRepository(db *gorm.DB) ModelRepository {
	return &modelRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *modelRepository) convertToDTO(model *model.Model) *dto.Model {
	if model == nil {
		return nil
	}

	var modelPrice dto.ModelPrice
	var modelAlias dto.ModelAlias
	var modelOptions dto.ModelOptions

	if err := model.ModelPrice.ToStruct(&modelPrice); err != nil {
		utils.SysError("解析价格配置失败:" + err.Error())
	}

	if err := model.ModelAlias.ToStruct(&modelAlias); err != nil {
		utils.SysError("解析模型映射失败:" + err.Error())
	}

	if err := model.ModelOptions.ToStruct(&modelOptions); err != nil {
		utils.SysError("解析模型配置失败:" + err.Error())
	}

	return &dto.Model{
		ModelID:          model.ModelID,
		ModelGroupID:     model.ModelGroupID,
		ModelName:        model.ModelName,
		ModelDescription: model.ModelDescription,
		ModelType:        model.ModelType,
		Provider:         model.Provider,
		PriceType:        model.PriceType,
		ModelPrice:       modelPrice,
		Status:           model.Status,
		ModelAlias:       modelAlias,
		ModelOptions:     modelOptions,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		DeletedAt:        utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *modelRepository) convertToModel(dto *dto.Model) (*model.Model, error) {
	if dto == nil {
		return nil, nil
	}

	modelPriceJSON, err := common.FromStruct(dto.ModelPrice)
	if err != nil {
		return nil, fmt.Errorf("转换价格配置失败: %w", err)
	}

	modelAliasJSON, err := common.FromStruct(dto.ModelAlias)
	if err != nil {
		return nil, fmt.Errorf("转换模型映射失败: %w", err)
	}

	modelOptionsJSON, err := common.FromStruct(dto.ModelOptions)
	if err != nil {
		return nil, fmt.Errorf("转换模型配置失败: %w", err)
	}

	return &model.Model{
		ModelID:          dto.ModelID,
		ModelGroupID:     dto.ModelGroupID,
		ModelName:        dto.ModelName,
		ModelDescription: dto.ModelDescription,
		ModelType:        dto.ModelType,
		Provider:         dto.Provider,
		PriceType:        dto.PriceType,
		ModelPrice:       modelPriceJSON,
		Status:           dto.Status,
		ModelAlias:       modelAliasJSON,
		ModelOptions:     modelOptionsJSON,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
		DeletedAt:        utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建模型
func (r *modelRepository) Create(model *dto.Model) error {
	modelData, err := r.convertToModel(model)
	if err != nil {
		return err
	}
	return r.db.Create(modelData).Error
}

// Update 更新模型
func (r *modelRepository) Update(dtoModel *dto.Model) error {
	modelData, err := r.convertToModel(dtoModel)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Model{}).Where("model_id = ?", dtoModel.ModelID).Updates(modelData).Error
}

// Delete 删除模型
func (r *modelRepository) Delete(modelID string) error {
	return r.db.Delete(&model.Model{}, "model_id = ?", modelID).Error
}

// GetByID 根据ID获取模型
func (r *modelRepository) GetByID(modelID string) (*dto.Model, error) {
	var model model.Model
	if err := r.db.Where("model_id = ?", modelID).First(&model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&model), nil
}

// GetByName 根据名称获取模型
func (r *modelRepository) GetByName(name string) (*dto.Model, error) {
	var model model.Model
	if err := r.db.Where("model_name = ? AND deleted_at IS NULL", name).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&model), nil
}

// GetByAlias 根据别名获取模型
func (r *modelRepository) GetByAlias(aliasName string) (*dto.Model, error) {
	var model model.Model
	if err := r.db.Where("(CAST(JSON_EXTRACT(model_alias, '$.display_name') AS CHAR(255)) = ? OR CAST(JSON_EXTRACT(model_alias, '$.request_name') AS CHAR(255)) = ?) AND deleted_at IS NULL",
		aliasName, aliasName).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&model), nil
}

// List 获取模型列表
func (r *modelRepository) List(page, pageSize int) ([]*dto.Model, int64, error) {
	var total int64
	var models []model.Model

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Model{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Model, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByType 根据类型获取模型列表
func (r *modelRepository) ListByType(modelType string, page, pageSize int) ([]*dto.Model, int64, error) {
	var total int64
	var models []model.Model

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Model{}).Where("model_type = ?", modelType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Model, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByProvider 根据提供商获取模型列表
func (r *modelRepository) ListByProvider(provider string, page, pageSize int) ([]*dto.Model, int64, error) {
	var total int64
	var models []model.Model

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Model{}).Where("provider = ?", provider)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Model, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByOptions 根据配置选项筛选模型
func (r *modelRepository) ListByOptions(options dto.ModelOptions, page, pageSize int) ([]*dto.Model, int64, error) {
	var total int64
	var models []model.Model

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Model{})

	optionsJSON, err := common.FromStruct(options)
	if err != nil {
		return nil, 0, fmt.Errorf("转换配置选项失败: %w", err)
	}

	query = query.Where("model_options @> ?", optionsJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Model, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *modelRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行模型基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testModel := &dto.Model{
			ModelID:          utils.GenerateRandomUUID(12),
			ModelGroupID:     utils.GenerateRandomUUID(12),
			ModelName:        fmt.Sprintf("benchmark_model_%d%s", i, utils.GenerateRandomUUID(4)),
			ModelDescription: fmt.Sprintf("这是第%d个基准测试模型", i),
			ModelType:        []string{"text", "image", "audio", "video"}[rand.Intn(4)],
			Provider:         []string{"openai", "anthropic", "google", "microsoft"}[rand.Intn(4)],
			PriceType:        "token",
			Status:           1,
			ModelPrice: dto.ModelPrice{
				RequestPrice:    float64(rand.Intn(100)) / 100,
				ResponsePrice:   float64(rand.Intn(100)) / 100,
				CompletionPrice: float64(rand.Intn(100)) / 100,
				CachePrice:      float64(rand.Intn(100)) / 100,
			},
			ModelAlias: dto.ModelAlias{
				DisplayName: fmt.Sprintf("Test Model %d", i),
				RequestName: fmt.Sprintf("TM%d", i),
			},
			ModelOptions: dto.ModelOptions{
				APIDiscount:         float64(rand.Intn(100)) / 100,
				APIDiscountExpireAt: utils.MySQLTime(time.Now().Add(time.Duration(rand.Intn(30)) * time.Hour)),
			},
		}

		// 创建
		if err := r.Create(testModel); err != nil {
			utils.SysError("创建模型失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdModel, err := r.GetByID(testModel.ModelID)
		if err != nil {
			utils.SysError("获取创建的模型失败: " + err.Error())
			return err
		}

		// 更新
		createdModel.ModelOptions.APIDiscount = float64(rand.Intn(100)) / 100
		if err := r.Update(createdModel); err != nil {
			utils.SysError("更新模型失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdModel.ModelID); err != nil {
			utils.SysError("删除模型失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
