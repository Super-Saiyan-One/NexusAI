package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	modelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ModelRepository 模型仓储接口
type ModelRepository interface {
	Create(model *dto.Model) (*dto.Model, error)
	Update(model *dto.Model) (*dto.Model, error)
	Delete(modelID string) error
	GetByID(modelID string) (*dto.Model, error)
	GetByName(name string) (*dto.Model, error)
	GetByAlias(aliasName string) (*dto.Model, error)
	List(page, pageSize int) ([]*dto.Model, int64, error)
	ListByType(modelType string, page, pageSize int) ([]*dto.Model, int64, error)
	ListByProvider(provider string, page, pageSize int) ([]*dto.Model, int64, error)
	ListByOptions(options dto.ModelOptions, page, pageSize int) ([]*dto.Model, int64, error)
	Search(req *modelDto.ModelSearchRequest) ([]*dto.Model, int64, error)
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
func (r *modelRepository) Create(model *dto.Model) (*dto.Model, error) {
	modelData, err := r.convertToModel(model)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(modelData).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(modelData), nil
}

// Update 更新模型
func (r *modelRepository) Update(dtoModel *dto.Model) (*dto.Model, error) {
	modelData, err := r.convertToModel(dtoModel)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Model{}).Where("model_id = ?", dtoModel.ModelID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(dtoModel.ModelID)
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

// Search 根据搜索条件筛选模型
func (r *modelRepository) Search(req *modelDto.ModelSearchRequest) ([]*dto.Model, int64, error) {
	var total int64
	var models []model.Model

	query := r.db.Model(&model.Model{})

	// 基本信息筛选
	if req.ModelID != "" {
		query = query.Where("model_id = ?", req.ModelID)
	}
	if req.ModelGroupID != "" {
		query = query.Where("model_group_id = ?", req.ModelGroupID)
	}
	if req.ModelName != "" {
		query = query.Where("model_name LIKE ?", "%"+req.ModelName+"%")
	}
	if req.ModelType != "" {
		query = query.Where("model_type = ?", req.ModelType)
	}
	if req.Provider != "" {
		query = query.Where("provider = ?", req.Provider)
	}
	if req.PriceType != "" {
		query = query.Where("price_type = ?", req.PriceType)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 模型别名查询
	if req.ModelDisplayName != "" {
		query = query.Where("CAST(JSON_EXTRACT(model_alias, '$.display_name') AS CHAR(255)) = ?", req.ModelDisplayName)
	}
	if req.ModelRequestName != "" {
		query = query.Where("CAST(JSON_EXTRACT(model_alias, '$.request_name') AS CHAR(255)) = ?", req.ModelRequestName)
	}

	// API折扣范围查询，考虑折扣过期时间
	if req.MinAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(model_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(model_options, '$.api_discount') AS DECIMAL(10,2)) 
			END >= ?`, req.MinAPIDiscount)
	}
	if req.MaxAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(model_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(model_options, '$.api_discount') AS DECIMAL(10,2)) 
			END <= ?`, req.MaxAPIDiscount)
	}

	// 价格范围查询
	if req.MinRequestPrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.request_price') AS DECIMAL(10,2)) >= ?", req.MinRequestPrice)
	}
	if req.MaxRequestPrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.request_price') AS DECIMAL(10,2)) <= ?", req.MaxRequestPrice)
	}

	if req.MinResponsePrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.response_price') AS DECIMAL(10,2)) >= ?", req.MinResponsePrice)
	}
	if req.MaxResponsePrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.response_price') AS DECIMAL(10,2)) <= ?", req.MaxResponsePrice)
	}

	if req.MinCompletionPrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.completion_price') AS DECIMAL(10,2)) >= ?", req.MinCompletionPrice)
	}
	if req.MaxCompletionPrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.completion_price') AS DECIMAL(10,2)) <= ?", req.MaxCompletionPrice)
	}

	if req.MinCachePrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.cache_price') AS DECIMAL(10,2)) >= ?", req.MinCachePrice)
	}
	if req.MaxCachePrice > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_price, '$.cache_price') AS DECIMAL(10,2)) <= ?", req.MaxCachePrice)
	}

	// 时间范围查询
	if !req.EarlyCreatedTime.IsZero() {
		query = query.Where("created_at >= ?", req.EarlyCreatedTime)
	}
	if !req.LateCreatedTime.IsZero() {
		query = query.Where("created_at <= ?", req.LateCreatedTime)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
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
		createdModel, err := r.Create(testModel)
		if err != nil {
			utils.SysError("创建模型失败: " + err.Error())
			return err
		}

		// 更新
		createdModel.ModelOptions.APIDiscount = float64(rand.Intn(100)) / 100
		if _, err := r.Update(createdModel); err != nil {
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
