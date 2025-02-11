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

// ChannelRepository 渠道仓储接口
type ChannelRepository interface {
	Create(channel *dto.Channel) error
	Update(channel *dto.Channel) error
	Delete(channelID string) error
	GetByID(channelID string) (*dto.Channel, error)
	GetByName(name string) (*dto.Channel, error)
	List(page, pageSize int) ([]*dto.Channel, int64, error)
	ListByStatus(status int8, page, pageSize int) ([]*dto.Channel, int64, error)
	ListByGroup(groupID string, page, pageSize int) ([]*dto.Channel, int64, error)
	Benchmark(count int) error
}

type channelRepository struct {
	db *gorm.DB
}

// NewChannelRepository 创建渠道仓储实例
func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *channelRepository) convertToDTO(model *model.Channel) *dto.Channel {
	if model == nil {
		return nil
	}

	var channelModels dto.ChannelModels
	var channelPriceFactor dto.ChannelPriceFactor
	var upstreamOptions dto.UpstreamOptions
	var authOptions dto.AuthOptions
	var retryOptions dto.RetryOptions
	var rateLimit dto.RateLimit
	var modelMapping dto.ModelMapping
	var testModels dto.TestModels

	if err := model.ChannelModels.ToStruct(&channelModels); err != nil {
		utils.SysError("解析渠道模型配置失败:" + err.Error())
	}

	if err := model.ChannelPriceFactor.ToStruct(&channelPriceFactor); err != nil {
		utils.SysError("解析价格系数失败:" + err.Error())
	}

	if err := model.UpstreamOptions.ToStruct(&upstreamOptions); err != nil {
		utils.SysError("解析上游配置失败:" + err.Error())
	}

	if err := model.AuthOptions.ToStruct(&authOptions); err != nil {
		utils.SysError("解析认证配置失败:" + err.Error())
	}

	if err := model.RetryOptions.ToStruct(&retryOptions); err != nil {
		utils.SysError("解析重试配置失败:" + err.Error())
	}

	if err := model.RateLimit.ToStruct(&rateLimit); err != nil {
		utils.SysError("解析速率限制失败:" + err.Error())
	}

	if err := model.ModelMapping.ToStruct(&modelMapping); err != nil {
		utils.SysError("解析模型映射失败:" + err.Error())
	}

	if err := model.TestModels.ToStruct(&testModels); err != nil {
		utils.SysError("解析测试模型配置失败:" + err.Error())
	}

	return &dto.Channel{
		ChannelID:          model.ChannelID,
		ChannelGroupID:     model.ChannelGroupID,
		ChannelName:        model.ChannelName,
		ChannelDescription: model.ChannelDescription,
		Status:             model.Status,
		ChannelModels:      channelModels,
		ChannelPriceFactor: channelPriceFactor,
		UpstreamOptions:    upstreamOptions,
		AuthOptions:        authOptions,
		RetryOptions:       retryOptions,
		RateLimit:          rateLimit,
		ModelMapping:       modelMapping,
		TestModels:         testModels,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		DeletedAt:          utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *channelRepository) convertToModel(dto *dto.Channel) (*model.Channel, error) {
	if dto == nil {
		return nil, nil
	}

	channelModelsJSON, err := common.FromStruct(dto.ChannelModels)
	if err != nil {
		return nil, fmt.Errorf("转换渠道模型配置失败: %w", err)
	}

	channelPriceFactorJSON, err := common.FromStruct(dto.ChannelPriceFactor)
	if err != nil {
		return nil, fmt.Errorf("转换价格系数失败: %w", err)
	}

	upstreamOptionsJSON, err := common.FromStruct(dto.UpstreamOptions)
	if err != nil {
		return nil, fmt.Errorf("转换上游配置失败: %w", err)
	}

	authOptionsJSON, err := common.FromStruct(dto.AuthOptions)
	if err != nil {
		return nil, fmt.Errorf("转换认证配置失败: %w", err)
	}

	retryOptionsJSON, err := common.FromStruct(dto.RetryOptions)
	if err != nil {
		return nil, fmt.Errorf("转换重试配置失败: %w", err)
	}

	rateLimitJSON, err := common.FromStruct(dto.RateLimit)
	if err != nil {
		return nil, fmt.Errorf("转换速率限制失败: %w", err)
	}

	modelMappingJSON, err := common.FromStruct(dto.ModelMapping)
	if err != nil {
		return nil, fmt.Errorf("转换模型映射失败: %w", err)
	}

	testModelsJSON, err := common.FromStruct(dto.TestModels)
	if err != nil {
		return nil, fmt.Errorf("转换测试模型配置失败: %w", err)
	}

	return &model.Channel{
		ChannelID:          dto.ChannelID,
		ChannelGroupID:     dto.ChannelGroupID,
		ChannelName:        dto.ChannelName,
		ChannelDescription: dto.ChannelDescription,
		Status:             dto.Status,
		ChannelModels:      channelModelsJSON,
		ChannelPriceFactor: channelPriceFactorJSON,
		UpstreamOptions:    upstreamOptionsJSON,
		AuthOptions:        authOptionsJSON,
		RetryOptions:       retryOptionsJSON,
		RateLimit:          rateLimitJSON,
		ModelMapping:       modelMappingJSON,
		TestModels:         testModelsJSON,
		CreatedAt:          dto.CreatedAt,
		UpdatedAt:          dto.UpdatedAt,
		DeletedAt:          utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建渠道
func (r *channelRepository) Create(channel *dto.Channel) error {
	model, err := r.convertToModel(channel)
	if err != nil {
		return err
	}
	return r.db.Create(model).Error
}

// Update 更新渠道
func (r *channelRepository) Update(channel *dto.Channel) error {
	modelData, err := r.convertToModel(channel)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Channel{}).Where("channel_id = ?", channel.ChannelID).Updates(modelData).Error
}

// Delete 删除渠道
func (r *channelRepository) Delete(channelID string) error {
	return r.db.Delete(&model.Channel{}, "channel_id = ?", channelID).Error
}

// GetByID 根据ID获取渠道
func (r *channelRepository) GetByID(channelID string) (*dto.Channel, error) {
	var channel model.Channel
	if err := r.db.Where("channel_id = ?", channelID).First(&channel).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channel), nil
}

// GetByName 根据名称获取渠道
func (r *channelRepository) GetByName(name string) (*dto.Channel, error) {
	var channel model.Channel
	if err := r.db.Where("channel_name = ?", name).First(&channel).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channel), nil
}

// List 获取渠道列表
func (r *channelRepository) List(page, pageSize int) ([]*dto.Channel, int64, error) {
	var total int64
	var channels []model.Channel

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Channel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Channel, len(channels))
	for i, c := range channels {
		dtoList[i] = r.convertToDTO(&c)
	}

	return dtoList, total, nil
}

// ListByStatus 根据状态获取渠道列表
func (r *channelRepository) ListByStatus(status int8, page, pageSize int) ([]*dto.Channel, int64, error) {
	var total int64
	var channels []model.Channel

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Channel{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Channel, len(channels))
	for i, c := range channels {
		dtoList[i] = r.convertToDTO(&c)
	}

	return dtoList, total, nil
}

// ListByGroup 根据渠道组获取渠道列表
func (r *channelRepository) ListByGroup(groupID string, page, pageSize int) ([]*dto.Channel, int64, error) {
	var total int64
	var channels []model.Channel

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Channel{}).Where("channel_group_id = ?", groupID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Channel, len(channels))
	for i, c := range channels {
		dtoList[i] = r.convertToDTO(&c)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *channelRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行渠道基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testChannel := &dto.Channel{
			ChannelID:          utils.GenerateRandomUUID(12),
			ChannelGroupID:     utils.GenerateRandomUUID(12),
			ChannelName:        fmt.Sprintf("benchmark_channel_%d", i),
			ChannelDescription: fmt.Sprintf("这是第%d个基准测试渠道", i),
			Status:             1,
			ChannelModels: dto.ChannelModels{
				AllowedModels:    []string{"gpt-3.5-turbo", "gpt-4"},
				DefaultTestModel: "gpt-3.5-turbo",
			},
			ChannelPriceFactor: dto.ChannelPriceFactor{
				RequestPriceFactor:    float64(rand.Intn(50)+50) / 100,
				ResponsePriceFactor:   float64(rand.Intn(50)+50) / 100,
				CompletionPriceFactor: float64(rand.Intn(50)+50) / 100,
				CachePriceFactor:      float64(rand.Intn(50)+50) / 100,
			},
			UpstreamOptions: dto.UpstreamOptions{
				Endpoint:    "https://api.example.com",
				Timeout:     30,
				MaxRetries:  3,
				ProxyURL:    "http://proxy.example.com",
				DialTimeout: 5,
			},
			AuthOptions: dto.AuthOptions{
				APIKey:        utils.GenerateRandomString(32),
				APISecret:     utils.GenerateRandomString(64),
				BearerToken:   utils.GenerateRandomString(48),
				Headers:       map[string]string{"X-Custom": "value"},
				RequestParams: map[string]string{"version": "v1"},
			},
			RetryOptions: dto.RetryOptions{
				MaxRetries:      3,
				RetryInterval:   1000,
				MaxRetryBackoff: 5000,
				RetryStatuses:   []int{429, 500, 502, 503, 504},
			},
			RateLimit: dto.RateLimit{
				RequestsPerSecond: rand.Intn(50) + 10,
				RequestsPerMinute: rand.Intn(1000) + 100,
				RequestsPerHour:   rand.Intn(10000) + 1000,
				RequestsPerDay:    rand.Intn(100000) + 10000,
			},
			ModelMapping: dto.ModelMapping{
				MappingModels: []string{"local-1", "upstream-1"},
			},
			TestModels: dto.TestModels{
				TestModels:       []string{"test-1", "test-2"},
				DefaultTestModel: "test-1",
			},
		}

		// 创建
		if err := r.Create(testChannel); err != nil {
			utils.SysError("创建渠道失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdChannel, err := r.GetByID(testChannel.ChannelID)
		if err != nil {
			utils.SysError("获取创建的渠道失败: " + err.Error())
			return err
		}

		// 更新
		createdChannel.Status = 2
		if err := r.Update(createdChannel); err != nil {
			utils.SysError("更新渠道失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdChannel.ChannelID); err != nil {
			utils.SysError("删除渠道失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
