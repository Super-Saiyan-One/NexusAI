package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	modelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ModelGroupRepository 模型组仓储接口
type ModelGroupRepository interface {
	Create(modelGroup *dto.ModelGroup) (*dto.ModelGroup, error)
	Update(modelGroup *dto.ModelGroup) (*dto.ModelGroup, error)
	Delete(modelGroupID string) error
	GetByID(modelGroupID string) (*dto.ModelGroup, error)
	GetByName(name string) (*dto.ModelGroup, error)
	List(page, pageSize int) ([]*dto.ModelGroup, int64, error)
	Search(modelGroupSearch *modelGroupDto.ModelGroupSearchRequest) ([]*dto.ModelGroup, int64, error)
	Benchmark(count int) error
}

type modelGroupRepository struct {
	db *gorm.DB
}

// NewModelGroupRepository 创建模型组仓储实例
func NewModelGroupRepository(db *gorm.DB) ModelGroupRepository {
	return &modelGroupRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *modelGroupRepository) convertToDTO(model *model.ModelGroup) *dto.ModelGroup {
	if model == nil {
		return nil
	}

	var priceFactor dto.ModelGroupPriceFactor
	var options dto.ModelGroupOptions

	if err := model.ModelGroupPriceFactor.ToStruct(&priceFactor); err != nil {
		utils.SysError("解析价格系数失败:" + err.Error())
	}

	if err := model.ModelGroupOptions.ToStruct(&options); err != nil {
		utils.SysError("解析配置选项失败:" + err.Error())
	}

	return &dto.ModelGroup{
		ModelGroupID:          model.ModelGroupID,
		ModelGroupName:        model.ModelGroupName,
		ModelGroupDescription: model.ModelGroupDescription,
		ModelGroupPriceFactor: priceFactor,
		ModelGroupOptions:     options,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
		DeletedAt:             utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *modelGroupRepository) convertToModel(dto *dto.ModelGroup) (*model.ModelGroup, error) {
	if dto == nil {
		return nil, nil
	}

	priceFactorJSON, err := common.FromStruct(dto.ModelGroupPriceFactor)
	if err != nil {
		return nil, fmt.Errorf("转换价格系数失败: %w", err)
	}

	optionsJSON, err := common.FromStruct(dto.ModelGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配置选项失败: %w", err)
	}

	return &model.ModelGroup{
		ModelGroupID:          dto.ModelGroupID,
		ModelGroupName:        dto.ModelGroupName,
		ModelGroupDescription: dto.ModelGroupDescription,
		ModelGroupPriceFactor: priceFactorJSON,
		ModelGroupOptions:     optionsJSON,
		CreatedAt:             dto.CreatedAt,
		UpdatedAt:             dto.UpdatedAt,
		DeletedAt:             utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建模型组
func (r *modelGroupRepository) Create(modelGroup *dto.ModelGroup) (*dto.ModelGroup, error) {
	model, err := r.convertToModel(modelGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新模型组
func (r *modelGroupRepository) Update(modelGroup *dto.ModelGroup) (*dto.ModelGroup, error) {
	modelData, err := r.convertToModel(modelGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.ModelGroup{}).Where("model_group_id = ?", modelGroup.ModelGroupID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(modelGroup.ModelGroupID)
}

// Delete 删除模型组
func (r *modelGroupRepository) Delete(modelGroupID string) error {
	return r.db.Delete(&model.ModelGroup{}, "model_group_id = ?", modelGroupID).Error
}

// GetByID 根据ID获取模型组
func (r *modelGroupRepository) GetByID(modelGroupID string) (*dto.ModelGroup, error) {
	var modelGroup model.ModelGroup
	if err := r.db.Where("model_group_id = ?", modelGroupID).First(&modelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&modelGroup), nil
}

// GetByName 根据名称获取模型组
func (r *modelGroupRepository) GetByName(name string) (*dto.ModelGroup, error) {
	var modelGroup model.ModelGroup
	if err := r.db.Where("model_group_name = ?", name).First(&modelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&modelGroup), nil
}

// List 获取模型组列表
func (r *modelGroupRepository) List(page, pageSize int) ([]*dto.ModelGroup, int64, error) {
	var total int64
	var modelGroups []model.ModelGroup

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.ModelGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&modelGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.ModelGroup, len(modelGroups))
	for i, mg := range modelGroups {
		dtoList[i] = r.convertToDTO(&mg)
	}

	return dtoList, total, nil
}

// Search 根据搜索条件筛选模型组
func (r *modelGroupRepository) Search(req *modelGroupDto.ModelGroupSearchRequest) ([]*dto.ModelGroup, int64, error) {
	var total int64
	var modelGroups []model.ModelGroup

	query := r.db.Model(&model.ModelGroup{})

	// 基本信息筛选
	if req.ModelGroupID != "" {
		query = query.Where("model_group_id = ?", req.ModelGroupID)
	}
	if req.ModelGroupName != "" {
		query = query.Where("model_group_name LIKE ?", "%"+req.ModelGroupName+"%")
	}
	if len(req.Levels) > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_options, '$.default_level') AS SIGNED) = ANY(?)", req.Levels)
	}

	// API折扣范围查询，考虑折扣过期时间
	if req.MinAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(model_group_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(model_group_options, '$.api_discount') AS DECIMAL(10,2)) 
			END >= ?`, req.MinAPIDiscount)
	}
	if req.MaxAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(model_group_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(model_group_options, '$.api_discount') AS DECIMAL(10,2)) 
			END <= ?`, req.MaxAPIDiscount)
	}

	// 并发请求数范围查询
	if req.MinConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_options, '$.max_concurrent_requests') AS SIGNED) >= ?", req.MinConcurrentRequests)
	}
	if req.MaxConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_options, '$.max_concurrent_requests') AS SIGNED) <= ?", req.MaxConcurrentRequests)
	}

	// 价格系数范围查询
	if req.MinRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) >= ?", req.MinRequestPriceFactor)
	}
	if req.MaxRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxRequestPriceFactor)
	}

	if req.MinResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) >= ?", req.MinResponsePriceFactor)
	}
	if req.MaxResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxResponsePriceFactor)
	}

	if req.MinCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCompletionPriceFactor)
	}
	if req.MaxCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCompletionPriceFactor)
	}

	if req.MinCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCachePriceFactor)
	}
	if req.MaxCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(model_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCachePriceFactor)
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
	if err := query.Find(&modelGroups).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtoList := make([]*dto.ModelGroup, len(modelGroups))
	for i, mg := range modelGroups {
		dtoList[i] = r.convertToDTO(&mg)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *modelGroupRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行模型组基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testModelGroup := &dto.ModelGroup{
			ModelGroupID:          utils.GenerateRandomUUID(12),
			ModelGroupName:        fmt.Sprintf("benchmark_model_group_%d", i),
			ModelGroupDescription: fmt.Sprintf("这是第%d个基准测试模型组", i),
			ModelGroupPriceFactor: dto.ModelGroupPriceFactor{
				RequestPriceFactor:    float64(rand.Intn(50)+50) / 100,
				ResponsePriceFactor:   float64(rand.Intn(50)+50) / 100,
				CompletionPriceFactor: float64(rand.Intn(50)+50) / 100,
				CachePriceFactor:      float64(rand.Intn(50)+50) / 100,
			},
			ModelGroupOptions: dto.ModelGroupOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				DefaultLevel:          rand.Intn(3) + 1,
				APIDiscount:           float64(rand.Intn(50)+50) / 100,
				APIDiscountExpireAt:   utils.MySQLTime(time.Now().Add(time.Duration(rand.Intn(30)) * time.Hour)),
			},
		}

		// 创建
		if _, err := r.Create(testModelGroup); err != nil {
			utils.SysError("创建模型组失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdModelGroup, err := r.GetByID(testModelGroup.ModelGroupID)
		if err != nil {
			utils.SysError("获取创建的模型组失败: " + err.Error())
			return err
		}

		// 更新
		createdModelGroup.ModelGroupOptions.DefaultLevel = rand.Intn(5) + 1
		if _, err := r.Update(createdModelGroup); err != nil {
			utils.SysError("更新模型组失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdModelGroup.ModelGroupID); err != nil {
			utils.SysError("删除模型组失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
