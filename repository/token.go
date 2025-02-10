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

// TokenRepository 令牌仓储接口
type TokenRepository interface {
	Create(token *dto.Token) (*dto.Token, error)
	Update(token *dto.Token) (*dto.Token, error)
	Delete(tokenID string) error
	GetByID(tokenID string) (*dto.Token, error)
	GetByKey(tokenKey string) (*dto.Token, error)
	GetByUserID(userID string) ([]*dto.Token, error)
	List(page, pageSize int) ([]*dto.Token, int64, error)
	ListByQuota(minQuota float64, page, pageSize int) ([]*dto.Token, int64, error)
	ListByOptions(options dto.TokenOptions, page, pageSize int) ([]*dto.Token, int64, error)
	Benchmark(count int) error
}

type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository 创建令牌仓储实例
func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *tokenRepository) convertToDTO(model *model.Token) *dto.Token {
	if model == nil {
		return nil
	}

	var channels dto.TokenChannels
	var models dto.TokenModels
	var options dto.TokenOptions

	if err := model.TokenChannels.ToStruct(&channels); err != nil {
		utils.SysError("解析渠道配置失败:" + err.Error())
	}

	if err := model.TokenModels.ToStruct(&models); err != nil {
		utils.SysError("解析模型配置失败:" + err.Error())
	}

	if err := model.TokenOptions.ToStruct(&options); err != nil {
		utils.SysError("解析配置选项失败:" + err.Error())
	}

	return &dto.Token{
		TokenID:          model.TokenID,
		UserID:           model.UserID,
		TokenName:        model.TokenName,
		TokenKey:         model.TokenKey,
		Status:           model.Status,
		TokenQuotaTotal:  model.TokenQuotaTotal,
		TokenQuotaUsed:   model.TokenQuotaUsed,
		TokenQuotaLeft:   model.TokenQuotaLeft,
		TokenQuotaFrozen: model.TokenQuotaFrozen,
		TokenChannels:    channels,
		TokenModels:      models,
		TokenOptions:     options,
		ExpireTime:       model.ExpireTime,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		DeletedAt:        utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *tokenRepository) convertToModel(dto *dto.Token) (*model.Token, error) {
	if dto == nil {
		return nil, nil
	}

	channelsJSON, err := common.FromStruct(dto.TokenChannels)
	if err != nil {
		return nil, fmt.Errorf("转换渠道配置失败: %w", err)
	}

	modelsJSON, err := common.FromStruct(dto.TokenModels)
	if err != nil {
		return nil, fmt.Errorf("转换模型配置失败: %w", err)
	}

	optionsJSON, err := common.FromStruct(dto.TokenOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配置选项失败: %w", err)
	}

	return &model.Token{
		TokenID:          dto.TokenID,
		UserID:           dto.UserID,
		TokenName:        dto.TokenName,
		TokenKey:         dto.TokenKey,
		Status:           dto.Status,
		TokenQuotaTotal:  dto.TokenQuotaTotal,
		TokenQuotaUsed:   dto.TokenQuotaUsed,
		TokenQuotaLeft:   dto.TokenQuotaLeft,
		TokenQuotaFrozen: dto.TokenQuotaFrozen,
		TokenChannels:    channelsJSON,
		TokenModels:      modelsJSON,
		TokenOptions:     optionsJSON,
		ExpireTime:       dto.ExpireTime,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
		DeletedAt:        utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建令牌
func (r *tokenRepository) Create(token *dto.Token) (*dto.Token, error) {
	model, err := r.convertToModel(token)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新令牌
func (r *tokenRepository) Update(token *dto.Token) (*dto.Token, error) {
	modelData, err := r.convertToModel(token)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Token{}).Where("token_id = ?", token.TokenID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(token.TokenID)
}

// Delete 删除令牌
func (r *tokenRepository) Delete(tokenID string) error {
	return r.db.Delete(&model.Token{}, "token_id = ?", tokenID).Error
}

// GetByID 根据ID获取令牌
func (r *tokenRepository) GetByID(tokenID string) (*dto.Token, error) {
	var token model.Token
	if err := r.db.Where("token_id = ?", tokenID).First(&token).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&token), nil
}

// GetByKey 根据密钥获取令牌
func (r *tokenRepository) GetByKey(tokenKey string) (*dto.Token, error) {
	var token model.Token
	if err := r.db.Where("token_key = ?", tokenKey).First(&token).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&token), nil
}

// GetByUserID 获取用户的所有令牌
func (r *tokenRepository) GetByUserID(userID string) ([]*dto.Token, error) {
	var tokens []model.Token
	if err := r.db.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	dtoList := make([]*dto.Token, len(tokens))
	for i, t := range tokens {
		dtoList[i] = r.convertToDTO(&t)
	}
	return dtoList, nil
}

// List 获取令牌列表
func (r *tokenRepository) List(page, pageSize int) ([]*dto.Token, int64, error) {
	var total int64
	var tokens []model.Token

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Token{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&tokens).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Token, len(tokens))
	for i, t := range tokens {
		dtoList[i] = r.convertToDTO(&t)
	}

	return dtoList, total, nil
}

// ListByQuota 根据配额筛选令牌
func (r *tokenRepository) ListByQuota(minQuota float64, page, pageSize int) ([]*dto.Token, int64, error) {
	var total int64
	var tokens []model.Token

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Token{}).Where("token_quota_left >= ?", minQuota)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&tokens).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Token, len(tokens))
	for i, t := range tokens {
		dtoList[i] = r.convertToDTO(&t)
	}

	return dtoList, total, nil
}

// ListByOptions 根据配置选项筛选令牌
func (r *tokenRepository) ListByOptions(options dto.TokenOptions, page, pageSize int) ([]*dto.Token, int64, error) {
	var total int64
	var tokens []model.Token

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Token{})

	optionsJSON, err := common.FromStruct(options)
	if err != nil {
		return nil, 0, fmt.Errorf("转换配置选项失败: %w", err)
	}

	query = query.Where("token_options @> ?", optionsJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&tokens).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Token, len(tokens))
	for i, t := range tokens {
		dtoList[i] = r.convertToDTO(&t)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *tokenRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行令牌基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testToken := &dto.Token{
			TokenID:         utils.GenerateRandomUUID(12),
			TokenName:       fmt.Sprintf("benchmark_token_%d", i),
			UserID:          fmt.Sprintf("test_user_%d", i),
			TokenKey:        utils.GenerateRandomString(32),
			Status:          1,
			TokenQuotaTotal: float64(rand.Intn(10000)),
			TokenQuotaUsed:  0,
			TokenQuotaLeft:  float64(rand.Intn(10000)),
			TokenChannels: dto.TokenChannels{
				ExtraAllowedChannels: []string{"openai", "anthropic"},
				PriorityChannels:     []string{"openai"},
				DefaultTestChannel:   "openai",
			},
			TokenModels: dto.TokenModels{
				AllowedModels:    []string{"gpt-3.5-turbo", "gpt-4"},
				DefaultTestModel: "gpt-3.5-turbo",
			},
			TokenOptions: dto.TokenOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				MaxRequestsPerMinute:  rand.Intn(100) + 1,
				MaxRequestsPerHour:    rand.Intn(1000) + 1,
				MaxRequestsPerDay:     rand.Intn(10000) + 1,
				RequireSignature:      rand.Intn(2) == 1,
				DisableRateLimit:      rand.Intn(2) == 1,
				AvailableLevels:       []int{1, 2, 3},
			},
		}

		// 创建
		_, err := r.Create(testToken)
		if err != nil {
			utils.SysError("创建令牌失败: " + err.Error())
			return err
		}

		// 获取创建后的记录以获得ID
		createdToken, err := r.GetByID(testToken.TokenID)
		if err != nil {
			utils.SysError("获取创建的令牌失败: " + err.Error())
			return err
		}

		// 更新
		createdToken.TokenOptions.MaxConcurrentRequests = rand.Intn(20) + 1
		_, err = r.Update(createdToken)
		if err != nil {
			utils.SysError("更新令牌失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdToken.TokenID); err != nil {
			utils.SysError("删除令牌失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
