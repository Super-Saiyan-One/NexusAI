package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	channelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ChannelRepository 渠道仓储接口
type ChannelRepository interface {
	Create(channel *dto.Channel) (*dto.Channel, error)
	Update(channel *dto.Channel) (*dto.Channel, error)
	Delete(channelID string) error
	GetByID(channelID string) (*dto.Channel, error)
	GetByName(name string) (*dto.Channel, error)
	List(page, pageSize int) ([]*dto.Channel, int64, error)
	ListByStatus(status int8, page, pageSize int) ([]*dto.Channel, int64, error)
	ListByGroup(groupID string, page, pageSize int) ([]*dto.Channel, int64, error)
	Search(req *channelDto.ChannelSearchRequest) ([]*dto.Channel, int64, error)
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
func (r *channelRepository) Create(channel *dto.Channel) (*dto.Channel, error) {
	model, err := r.convertToModel(channel)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新渠道
func (r *channelRepository) Update(channel *dto.Channel) (*dto.Channel, error) {
	modelData, err := r.convertToModel(channel)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Channel{}).Where("channel_id = ?", channel.ChannelID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(channel.ChannelID)
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

// Search 根据搜索条件筛选渠道
func (r *channelRepository) Search(req *channelDto.ChannelSearchRequest) ([]*dto.Channel, int64, error) {
	var total int64
	var channels []model.Channel

	query := r.db.Model(&model.Channel{})

	// 基本信息筛选
	if req.ChannelID != "" {
		query = query.Where("channel_id = ?", req.ChannelID)
	}
	if req.ChannelGroupID != "" {
		query = query.Where("channel_group_id = ?", req.ChannelGroupID)
	}
	if req.ChannelName != "" {
		query = query.Where("channel_name LIKE ?", "%"+req.ChannelName+"%")
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 允许的模型列表筛选
	if len(req.AllowedModels) > 0 {
		for _, modelID := range req.AllowedModels {
			query = query.Where("JSON_CONTAINS(CAST(channel_models->>'$.allowed_models' AS JSON), JSON_ARRAY(?))", modelID)
		}
	}

	// 上游服务配置筛选
	if req.UpstreamEndpoint != "" {
		query = query.Where("JSON_EXTRACT(upstream_options, '$.endpoint') = ?", req.UpstreamEndpoint)
	}
	if req.UpstreamProxyURL != "" {
		query = query.Where("JSON_EXTRACT(upstream_options, '$.proxy_url') = ?", req.UpstreamProxyURL)
	}
	if req.MinUpstreamTimeout > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.timeout') AS SIGNED) >= ?", req.MinUpstreamTimeout)
	}
	if req.MaxUpstreamTimeout > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.timeout') AS SIGNED) <= ?", req.MaxUpstreamTimeout)
	}
	if req.MinUpstreamMaxRetries > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.max_retries') AS SIGNED) >= ?", req.MinUpstreamMaxRetries)
	}
	if req.MaxUpstreamMaxRetries > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.max_retries') AS SIGNED) <= ?", req.MaxUpstreamMaxRetries)
	}
	if req.MinUpstreamDialTimeout > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.dial_timeout') AS SIGNED) >= ?", req.MinUpstreamDialTimeout)
	}
	if req.MaxUpstreamDialTimeout > 0 {
		query = query.Where("CAST(JSON_EXTRACT(upstream_options, '$.dial_timeout') AS SIGNED) <= ?", req.MaxUpstreamDialTimeout)
	}

	// 认证配置筛选
	if req.AuthAPIKey != "" {
		query = query.Where("JSON_EXTRACT(auth_options, '$.api_key') = ?", req.AuthAPIKey)
	}
	if req.AuthAPISecret != "" {
		query = query.Where("JSON_EXTRACT(auth_options, '$.api_secret') = ?", req.AuthAPISecret)
	}
	if req.AuthBearerToken != "" {
		query = query.Where("JSON_EXTRACT(auth_options, '$.bearer_token') = ?", req.AuthBearerToken)
	}

	// 重试配置筛选
	if req.MinRetryMaxRetries > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.max_retries') AS SIGNED) >= ?", req.MinRetryMaxRetries)
	}
	if req.MaxRetryMaxRetries > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.max_retries') AS SIGNED) <= ?", req.MaxRetryMaxRetries)
	}
	if req.MinRetryInterval > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.retry_interval') AS SIGNED) >= ?", req.MinRetryInterval)
	}
	if req.MaxRetryInterval > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.retry_interval') AS SIGNED) <= ?", req.MaxRetryInterval)
	}
	if req.MinRetryMaxRetryBackoff > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.max_retry_backoff') AS SIGNED) >= ?", req.MinRetryMaxRetryBackoff)
	}
	if req.MaxRetryMaxRetryBackoff > 0 {
		query = query.Where("CAST(JSON_EXTRACT(retry_options, '$.max_retry_backoff') AS SIGNED) <= ?", req.MaxRetryMaxRetryBackoff)
	}
	if len(req.RetryStatuses) > 0 {
		for _, status := range req.RetryStatuses {
			query = query.Where("JSON_CONTAINS(CAST(retry_options->>'$.retry_statuses' AS JSON), CAST(? AS JSON))", status)
		}
	}

	// 速率限制筛选
	if req.MinRateLimitRequestsPerSecond > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_second') AS SIGNED) >= ?", req.MinRateLimitRequestsPerSecond)
	}
	if req.MaxRateLimitRequestsPerSecond > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_second') AS SIGNED) <= ?", req.MaxRateLimitRequestsPerSecond)
	}
	if req.MinRateLimitRequestsPerMinute > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_minute') AS SIGNED) >= ?", req.MinRateLimitRequestsPerMinute)
	}
	if req.MaxRateLimitRequestsPerMinute > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_minute') AS SIGNED) <= ?", req.MaxRateLimitRequestsPerMinute)
	}
	if req.MinRateLimitRequestsPerHour > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_hour') AS SIGNED) >= ?", req.MinRateLimitRequestsPerHour)
	}
	if req.MaxRateLimitRequestsPerHour > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_hour') AS SIGNED) <= ?", req.MaxRateLimitRequestsPerHour)
	}
	if req.MinRateLimitRequestsPerDay > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_day') AS SIGNED) >= ?", req.MinRateLimitRequestsPerDay)
	}
	if req.MaxRateLimitRequestsPerDay > 0 {
		query = query.Where("CAST(JSON_EXTRACT(rate_limit, '$.requests_per_day') AS SIGNED) <= ?", req.MaxRateLimitRequestsPerDay)
	}

	// 价格系数筛选
	if req.MinRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) >= ?", req.MinRequestPriceFactor)
	}
	if req.MaxRequestPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.request_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxRequestPriceFactor)
	}
	if req.MinResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) >= ?", req.MinResponsePriceFactor)
	}
	if req.MaxResponsePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.response_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxResponsePriceFactor)
	}
	if req.MinCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCompletionPriceFactor)
	}
	if req.MaxCompletionPriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.completion_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCompletionPriceFactor)
	}
	if req.MinCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) >= ?", req.MinCachePriceFactor)
	}
	if req.MaxCachePriceFactor > 0 {
		query = query.Where("CAST(JSON_EXTRACT(channel_price_factor, '$.cache_price_factor') AS DECIMAL(10,2)) <= ?", req.MaxCachePriceFactor)
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
	if err := query.Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
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
				MappingModels: map[string][]string{"local-1": {"upstream-1"}},
			},
			TestModels: dto.TestModels{
				TestModels:       []string{"test-1", "test-2"},
				DefaultTestModel: "test-1",
			},
		}

		// 创建
		createdChannel, err := r.Create(testChannel)
		if err != nil {
			utils.SysError("创建渠道失败: " + err.Error())
			return err
		}

		// 更新
		createdChannel.Status = 2
		if _, err := r.Update(createdChannel); err != nil {
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
