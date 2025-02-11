package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	channelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ChannelGroupRepository 渠道组仓储接口
type ChannelGroupRepository interface {
	Create(channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error)
	Update(channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error)
	Delete(channelGroupID string) error
	GetByID(channelGroupID string) (*dto.ChannelGroup, error)
	GetByName(name string) (*dto.ChannelGroup, error)
	List(page, pageSize int) ([]*dto.ChannelGroup, int64, error)
	GetListByDefaultLevel(defaultLevel int) ([]*dto.ChannelGroup, error)
	Search(req *channelGroupDto.ChannelGroupSearchRequest) ([]*dto.ChannelGroup, int64, error)
	Benchmark(count int) error
}

type channelGroupRepository struct {
	db *gorm.DB
}

// NewChannelGroupRepository 创建渠道组仓储实例
func NewChannelGroupRepository(db *gorm.DB) ChannelGroupRepository {
	return &channelGroupRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *channelGroupRepository) convertToDTO(model *model.ChannelGroup) *dto.ChannelGroup {
	if model == nil {
		return nil
	}

	var priceFactor dto.ChannelGroupPriceFactor
	var options dto.ChannelGroupOptions
	var channels dto.ChannelGroupChannels

	if err := model.ChannelGroupPriceFactor.ToStruct(&priceFactor); err != nil {
		utils.SysError("解析价格系数失败:" + err.Error())
	}

	if err := model.ChannelGroupOptions.ToStruct(&options); err != nil {
		utils.SysError("解析配置选项失败:" + err.Error())
	}

	if err := model.ChannelGroupChannels.ToStruct(&channels); err != nil {
		utils.SysError("解析渠道组渠道失败:" + err.Error())
	}

	return &dto.ChannelGroup{
		ChannelGroupID:          model.ChannelGroupID,
		ChannelGroupName:        model.ChannelGroupName,
		ChannelGroupDescription: model.ChannelGroupDescription,
		ChannelGroupPriceFactor: priceFactor,
		ChannelGroupOptions:     options,
		ChannelGroupChannels:    channels,
		CreatedAt:               model.CreatedAt,
		UpdatedAt:               model.UpdatedAt,
		DeletedAt:               utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *channelGroupRepository) convertToModel(dto *dto.ChannelGroup) (*model.ChannelGroup, error) {
	if dto == nil {
		return nil, nil
	}

	priceFactorJSON, err := common.FromStruct(dto.ChannelGroupPriceFactor)
	if err != nil {
		return nil, fmt.Errorf("转换价格系数失败: %w", err)
	}

	optionsJSON, err := common.FromStruct(dto.ChannelGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配置选项失败: %w", err)
	}

	channelsJSON, err := common.FromStruct(dto.ChannelGroupChannels)
	if err != nil {
		return nil, fmt.Errorf("转换渠道组渠道失败: %w", err)
	}

	return &model.ChannelGroup{
		ChannelGroupID:          dto.ChannelGroupID,
		ChannelGroupName:        dto.ChannelGroupName,
		ChannelGroupDescription: dto.ChannelGroupDescription,
		ChannelGroupPriceFactor: priceFactorJSON,
		ChannelGroupOptions:     optionsJSON,
		ChannelGroupChannels:    channelsJSON,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,
		DeletedAt:               utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建渠道组
func (r *channelGroupRepository) Create(channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error) {
	model, err := r.convertToModel(channelGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新渠道组
func (r *channelGroupRepository) Update(channelGroup *dto.ChannelGroup) (*dto.ChannelGroup, error) {
	modelData, err := r.convertToModel(channelGroup)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.ChannelGroup{}).Where("channel_group_id = ?", channelGroup.ChannelGroupID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(channelGroup.ChannelGroupID)
}

// Delete 删除渠道组
func (r *channelGroupRepository) Delete(channelGroupID string) error {
	return r.db.Delete(&model.ChannelGroup{}, "channel_group_id = ?", channelGroupID).Error
}

// GetByID 根据ID获取渠道组
func (r *channelGroupRepository) GetByID(channelGroupID string) (*dto.ChannelGroup, error) {
	var channelGroup model.ChannelGroup
	if err := r.db.Where("channel_group_id = ?", channelGroupID).First(&channelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channelGroup), nil
}

// GetByName 根据名称获取渠道组
func (r *channelGroupRepository) GetByName(name string) (*dto.ChannelGroup, error) {
	var channelGroup model.ChannelGroup
	if err := r.db.Where("channel_group_name = ?", name).First(&channelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channelGroup), nil
}

// List 获取渠道组列表
func (r *channelGroupRepository) List(page, pageSize int) ([]*dto.ChannelGroup, int64, error) {
	var total int64
	var channelGroups []model.ChannelGroup

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.ChannelGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&channelGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.ChannelGroup, len(channelGroups))
	for i, cg := range channelGroups {
		dtoList[i] = r.convertToDTO(&cg)
	}

	return dtoList, total, nil
}

// GetListByDefaultLevel 根据默认等级获取渠道组列表
func (r *channelGroupRepository) GetListByDefaultLevel(defaultLevel int) ([]*dto.ChannelGroup, error) {
	var channelGroups []model.ChannelGroup
	if err := r.db.Where("CAST(JSON_EXTRACT(channel_group_options, '$.default_level') AS SIGNED) = ?", defaultLevel).Find(&channelGroups).Error; err != nil {
		return nil, err
	}
	dtoList := make([]*dto.ChannelGroup, len(channelGroups))
	for i, cg := range channelGroups {
		dtoList[i] = r.convertToDTO(&cg)
	}
	return dtoList, nil
}

// Search 根据搜索条件筛选渠道组
func (r *channelGroupRepository) Search(req *channelGroupDto.ChannelGroupSearchRequest) ([]*dto.ChannelGroup, int64, error) {
	var total int64
	var channelGroups []model.ChannelGroup

	query := r.db.Model(&model.ChannelGroup{})

	// 基本信息筛选
	if req.ChannelGroupID != "" {
		query = query.Where("channel_group_id = ?", req.ChannelGroupID)
	}
	if req.ChannelGroupName != "" {
		query = query.Where("channel_group_name LIKE ?", "%"+req.ChannelGroupName+"%")
	}

	// 等级筛选
	if len(req.Levels) > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_options, '$.default_level') AS SIGNED) IN (?)", req.Levels)
	}

	// 渠道和模型筛选
	if len(req.Channels) > 0 {
		for _, channelID := range req.Channels {
			query = query.Where("JSON_CONTAINS(CAST(channel_group_channels->>'$.channels' AS JSON), JSON_ARRAY(?))", channelID)
		}
	}
	if len(req.Models) > 0 {
		for _, modelID := range req.Models {
			query = query.Where("JSON_CONTAINS(JSON_KEYS(CAST(channel_group_channels->>'$.models_map' AS JSON)), JSON_ARRAY(?))", modelID)
		}
	}

	// API折扣范围查询，考虑折扣过期时间
	if req.MinAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(channel_group_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(channel_group_options, '$.api_discount') AS DECIMAL(10,2)) 
			END >= ?`, req.MinAPIDiscount)
	}
	if req.MaxAPIDiscount > 0 {
		query = query.Where(`
			CASE 
				WHEN CAST(JSON_EXTRACT(channel_group_options, '$.api_discount_expire_at') AS DATETIME) < NOW() 
				THEN 1 
				ELSE CAST(JSON_EXTRACT(channel_group_options, '$.api_discount') AS DECIMAL(10,2)) 
			END <= ?`, req.MaxAPIDiscount)
	}

	// 并发请求数范围查询
	if req.MinConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_options, '$.max_concurrent_requests') AS SIGNED) >= ?", req.MinConcurrentRequests)
	}
	if req.MaxConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_options, '$.max_concurrent_requests') AS SIGNED) <= ?", req.MaxConcurrentRequests)
	}

	// 价格系数范围查询
	if req.MinRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) >= ?", req.MinRequestPriceFactor)
	}
	if req.MaxRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxRequestPriceFactor)
	}
	if req.MinResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) >= ?", req.MinResponsePriceFactor)
	}
	if req.MaxResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxResponsePriceFactor)
	}
	if req.MinCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCompletionPriceFactor)
	}
	if req.MaxCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCompletionPriceFactor)
	}
	if req.MinCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCachePriceFactor)
	}
	if req.MaxCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_group_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCachePriceFactor)
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
	if err := query.Find(&channelGroups).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtoList := make([]*dto.ChannelGroup, len(channelGroups))
	for i, cg := range channelGroups {
		dtoList[i] = r.convertToDTO(&cg)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *channelGroupRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行渠道组基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testChannelGroup := &dto.ChannelGroup{
			ChannelGroupID:          utils.GenerateRandomUUID(12),
			ChannelGroupName:        fmt.Sprintf("benchmark_channel_group_%d", i),
			ChannelGroupDescription: fmt.Sprintf("这是第%d个基准测试渠道组", i),
			ChannelGroupPriceFactor: dto.ChannelGroupPriceFactor{
				RequestPriceFactor:    float64(rand.Intn(50)+50) / 100,
				ResponsePriceFactor:   float64(rand.Intn(50)+50) / 100,
				CompletionPriceFactor: float64(rand.Intn(50)+50) / 100,
				CachePriceFactor:      float64(rand.Intn(50)+50) / 100,
			},
			ChannelGroupOptions: dto.ChannelGroupOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				DefaultLevel:          rand.Intn(3) + 1,
				APIDiscount:           float64(rand.Intn(50)+50) / 100,
				APIDiscountExpireAt:   utils.MySQLTime(time.Now().Add(time.Duration(rand.Intn(30)) * time.Hour)),
			},
			ChannelGroupChannels: dto.ChannelGroupChannels{
				Channels:  []string{},
				ModelsMap: map[string][]string{},
			},
		}

		// 创建
		createdChannelGroup, err := r.Create(testChannelGroup)
		if err != nil {
			utils.SysError("创建渠道组失败: " + err.Error())
			return err
		}

		// 更新
		createdChannelGroup.ChannelGroupOptions.DefaultLevel = rand.Intn(5) + 1
		if _, err := r.Update(createdChannelGroup); err != nil {
			utils.SysError("更新渠道组失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdChannelGroup.ChannelGroupID); err != nil {
			utils.SysError("删除渠道组失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
