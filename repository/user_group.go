package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	userGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// UserGroupRepository 用户组仓储接口
type UserGroupRepository interface {
	Create(userGroup *dto.UserGroup) (*dto.UserGroup, error)
	Update(userGroup *dto.UserGroup) (*dto.UserGroup, error)
	Delete(userGroupID string) error
	GetByID(userGroupID string) (*dto.UserGroup, error)
	GetByName(name string) (*dto.UserGroup, error)
	List(page, pageSize int) ([]*dto.UserGroup, int64, error)
	ListByPriceFactor(factor dto.UserGroupPriceFactor, page, pageSize int) ([]*dto.UserGroup, int64, error)
	ListByOptions(options dto.UserGroupOptions, page, pageSize int) ([]*dto.UserGroup, int64, error)
	Search(userGroupSearch *userGroupDto.UserGroupSearchRequest) ([]*dto.UserGroup, int64, error)
	Benchmark(count int) error
}

type userGroupRepository struct {
	db *gorm.DB
}

// NewUserGroupRepository 创建用户组仓储实例
func NewUserGroupRepository(db *gorm.DB) UserGroupRepository {
	return &userGroupRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *userGroupRepository) convertToDTO(model *model.UserGroup) *dto.UserGroup {
	if model == nil {
		return nil
	}

	var priceFactor dto.UserGroupPriceFactor
	var options dto.UserGroupOptions

	if err := model.UserGroupPriceFactor.ToStruct(&priceFactor); err != nil {
		utils.SysError("解析价格系数失败:" + err.Error())
	}

	if err := model.UserGroupOptions.ToStruct(&options); err != nil {
		utils.SysError("解析配置选项失败:" + err.Error())
	}

	return &dto.UserGroup{
		UserGroupID:          model.UserGroupID,
		UserGroupName:        model.UserGroupName,
		UserGroupDescription: model.UserGroupDescription,
		UserGroupPriceFactor: priceFactor,
		UserGroupOptions:     options,
		CreatedAt:            model.CreatedAt,
		UpdatedAt:            model.UpdatedAt,
		DeletedAt:            utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *userGroupRepository) convertToModel(dto *dto.UserGroup) (*model.UserGroup, error) {
	if dto == nil {
		return nil, nil
	}

	priceFactorJSON, err := common.FromStruct(dto.UserGroupPriceFactor)
	if err != nil {
		return nil, fmt.Errorf("转换价格系数失败: %w", err)
	}

	optionsJSON, err := common.FromStruct(dto.UserGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配置选项失败: %w", err)
	}

	return &model.UserGroup{
		UserGroupID:          dto.UserGroupID,
		UserGroupName:        dto.UserGroupName,
		UserGroupDescription: dto.UserGroupDescription,
		UserGroupPriceFactor: priceFactorJSON,
		UserGroupOptions:     optionsJSON,
		CreatedAt:            dto.CreatedAt,
		UpdatedAt:            dto.UpdatedAt,
		DeletedAt:            utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建用户组
func (r *userGroupRepository) Create(userGroup *dto.UserGroup) (*dto.UserGroup, error) {
	model, err := r.convertToModel(userGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新用户组
func (r *userGroupRepository) Update(userGroup *dto.UserGroup) (*dto.UserGroup, error) {
	modelData, err := r.convertToModel(userGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.UserGroup{}).Where("user_group_id = ?", userGroup.UserGroupID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(userGroup.UserGroupID)
}

// Delete 删除用户组
func (r *userGroupRepository) Delete(userGroupID string) error {
	return r.db.Delete(&model.UserGroup{}, "user_group_id = ?", userGroupID).Error
}

// GetByID 根据ID获取用户组
func (r *userGroupRepository) GetByID(userGroupID string) (*dto.UserGroup, error) {
	var userGroup model.UserGroup
	if err := r.db.Where("user_group_id = ?", userGroupID).First(&userGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&userGroup), nil
}

// GetByName 根据名称获取用户组
func (r *userGroupRepository) GetByName(name string) (*dto.UserGroup, error) {
	var userGroup model.UserGroup
	if err := r.db.Where("user_group_name = ?", name).First(&userGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&userGroup), nil
}

// List 获取用户组列表
func (r *userGroupRepository) List(page, pageSize int) ([]*dto.UserGroup, int64, error) {
	var total int64
	var userGroups []model.UserGroup

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.UserGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&userGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.UserGroup, len(userGroups))
	for i, ug := range userGroups {
		dtoList[i] = r.convertToDTO(&ug)
	}

	return dtoList, total, nil
}

// ListByPriceFactor 根据价格系数筛选用户组
func (r *userGroupRepository) ListByPriceFactor(factor dto.UserGroupPriceFactor, page, pageSize int) ([]*dto.UserGroup, int64, error) {
	var total int64
	var userGroups []model.UserGroup

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.UserGroup{})

	// 构建查询条件
	factorJSON, err := common.FromStruct(factor)
	if err != nil {
		return nil, 0, fmt.Errorf("转换价格系数失败: %w", err)
	}

	query = query.Where("user_group_price_factor @> ?", factorJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&userGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.UserGroup, len(userGroups))
	for i, ug := range userGroups {
		dtoList[i] = r.convertToDTO(&ug)
	}

	return dtoList, total, nil
}

// ListByOptions 根据配置选项筛选用户组
func (r *userGroupRepository) ListByOptions(options dto.UserGroupOptions, page, pageSize int) ([]*dto.UserGroup, int64, error) {
	var total int64
	var userGroups []model.UserGroup

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.UserGroup{})

	// 构建查询条件
	optionsJSON, err := common.FromStruct(options)
	if err != nil {
		return nil, 0, fmt.Errorf("转换配置选项失败: %w", err)
	}

	query = query.Where("user_group_options @> ?", optionsJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&userGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.UserGroup, len(userGroups))
	for i, ug := range userGroups {
		dtoList[i] = r.convertToDTO(&ug)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *userGroupRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行用户组基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		// 创建测试用户组
		testGroup := &dto.UserGroup{
			UserGroupID:          utils.GenerateRandomUUID(12),
			UserGroupName:        fmt.Sprintf("benchmark_group_%d", i),
			UserGroupDescription: "基准测试用户组",
			UserGroupPriceFactor: dto.UserGroupPriceFactor{
				RequestPriceFactor:    float64(rand.Intn(50)+50) / 100,
				ResponsePriceFactor:   float64(rand.Intn(50)+50) / 100,
				CompletionPriceFactor: float64(rand.Intn(50)+50) / 100,
				CachePriceFactor:      float64(rand.Intn(50)+50) / 100,
			},
			UserGroupOptions: dto.UserGroupOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				DefaultLevel:          rand.Intn(3) + 1,
				ExtraAllowedModels:    []string{"gpt-3.5-turbo", "gpt-4"},
				ExtraAllowedChannels:  []string{"openai", "anthropic"},
				APIDiscount:           float64(rand.Intn(50)+50) / 100,
			},
		}

		// 创建
		_, err := r.Create(testGroup)
		if err != nil {
			utils.SysError("创建用户组失败: " + err.Error())
			return err
		}

		// 在创建后获取ID
		createdGroup, err := r.GetByID(testGroup.UserGroupID)
		if err != nil {
			utils.SysError("获取创建的用户组失败: " + err.Error())
			return err
		}

		// 更新时使用获取到的完整记录
		createdGroup.UserGroupDescription = "已更新的基准测试用户组"
		if _, err := r.Update(createdGroup); err != nil {
			utils.SysError("更新用户组失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdGroup.UserGroupID); err != nil {
			utils.SysError("删除用户组失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}

// Search 根据搜索条件筛选用户组
func (r *userGroupRepository) Search(req *userGroupDto.UserGroupSearchRequest) ([]*dto.UserGroup, int64, error) {
	var total int64
	var userGroups []model.UserGroup

	query := r.db.Model(&model.UserGroup{})

	// 基本信息筛选
	if req.UserGroupID != "" {
		query = query.Where("user_group_id = ?", req.UserGroupID)
	}
	if req.UserGroupName != "" {
		query = query.Where("user_group_name LIKE ?", "%"+req.UserGroupName+"%")
	}
	if len(req.Levels) > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_options, '$.default_level') AS SIGNED) = ANY(?)", req.Levels)
	}

	// 额外允许的模型和渠道筛选
	if len(req.ExtraAllowedModels) > 0 {
		for _, model := range req.ExtraAllowedModels {
			query = query.Where("JSON_CONTAINS(user_group_options, JSON_ARRAY(?), '$.extra_allowed_models')", model)
		}
	}
	if len(req.ExtraAllowedChannels) > 0 {
		for _, channel := range req.ExtraAllowedChannels {
			query = query.Where("JSON_CONTAINS(user_group_options, JSON_ARRAY(?), '$.extra_allowed_channels')", channel)
		}
	}

	// API折扣范围查询
	if req.MinAPIDiscount > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_options, '$.api_discount') AS DECIMAL(10,2)) >= ?", req.MinAPIDiscount)
	}
	if req.MaxAPIDiscount > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_options, '$.api_discount') AS DECIMAL(10,2)) <= ?", req.MaxAPIDiscount)
	}

	// 并发请求数范围查询
	if req.MinConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_options, '$.max_concurrent_requests') AS SIGNED) >= ?", req.MinConcurrentRequests)
	}
	if req.MaxConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_options, '$.max_concurrent_requests') AS SIGNED) <= ?", req.MaxConcurrentRequests)
	}

	// 价格系数范围查询
	if req.MinRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) >= ?", req.MinRequestPriceFactor)
	}
	if req.MaxRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxRequestPriceFactor)
	}

	if req.MinResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) >= ?", req.MinResponsePriceFactor)
	}
	if req.MaxResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxResponsePriceFactor)
	}

	if req.MinCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCompletionPriceFactor)
	}
	if req.MaxCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCompletionPriceFactor)
	}

	if req.MinCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCachePriceFactor)
	}
	if req.MaxCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCachePriceFactor)
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
	if err := query.Find(&userGroups).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtoList := make([]*dto.UserGroup, len(userGroups))
	for i, ug := range userGroups {
		dtoList[i] = r.convertToDTO(&ug)
	}

	return dtoList, total, nil
}
