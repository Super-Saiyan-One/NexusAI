package repository

import (
	"fmt"
	"math/rand"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// UsageRepository 用量记录仓储接口
type UsageRepository interface {
	Create(usage *dto.Usage) error
	Update(usage *dto.Usage) error
	Delete(usageID string) error
	GetByID(usageID string) (*dto.Usage, error)
	GetByRequestID(requestID string) (*dto.Usage, error)
	List(page, pageSize int) ([]*dto.Usage, int64, error)
	ListByUser(userID string, page, pageSize int) ([]*dto.Usage, int64, error)
	ListByToken(tokenID string, page, pageSize int) ([]*dto.Usage, int64, error)
	ListByModel(modelID string, page, pageSize int) ([]*dto.Usage, int64, error)
	ListByChannel(channelID string, page, pageSize int) ([]*dto.Usage, int64, error)
	Benchmark(count int) error
}

type usageRepository struct {
	db *gorm.DB
}

// NewUsageRepository 创建用量记录仓储实例
func NewUsageRepository(db *gorm.DB) UsageRepository {
	return &usageRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *usageRepository) convertToDTO(model *model.Usage) *dto.Usage {
	if model == nil {
		return nil
	}
	return &dto.Usage{
		UsageID:          model.UsageID,
		TokenID:          model.TokenID,
		UserID:           model.UserID,
		ChannelID:        model.ChannelID,
		ModelID:          model.ModelID,
		RequestID:        model.RequestID,
		UsageType:        model.UsageType,
		UnitPrice:        model.UnitPrice,
		TimesCount:       model.TimesCount,
		TokensCount:      model.TokensCount,
		PriceTotalFactor: model.PriceTotalFactor,
		TotalAmount:      model.TotalAmount,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		DeletedAt:        utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *usageRepository) convertToModel(dto *dto.Usage) *model.Usage {
	if dto == nil {
		return nil
	}

	return &model.Usage{
		UsageID:          dto.UsageID,
		TokenID:          dto.TokenID,
		UserID:           dto.UserID,
		ChannelID:        dto.ChannelID,
		ModelID:          dto.ModelID,
		RequestID:        dto.RequestID,
		UsageType:        dto.UsageType,
		UnitPrice:        dto.UnitPrice,
		TimesCount:       dto.TimesCount,
		TokensCount:      dto.TokensCount,
		PriceTotalFactor: dto.PriceTotalFactor,
		TotalAmount:      dto.TotalAmount,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
		DeletedAt:        utils.ToDeletedAt(dto.DeletedAt),
	}
}

// Create 创建用量记录
func (r *usageRepository) Create(usage *dto.Usage) error {
	return r.db.Create(r.convertToModel(usage)).Error
}

// Update 更新用量记录
func (r *usageRepository) Update(usage *dto.Usage) error {
	return r.db.Model(&model.Usage{}).Where("usage_id = ?", usage.UsageID).Updates(r.convertToModel(usage)).Error
}

// Delete 删除用量记录
func (r *usageRepository) Delete(usageID string) error {
	return r.db.Delete(&model.Usage{}, "usage_id = ?", usageID).Error
}

// GetByID 根据ID获取用量记录
func (r *usageRepository) GetByID(usageID string) (*dto.Usage, error) {
	var usage model.Usage
	if err := r.db.Where("usage_id = ?", usageID).First(&usage).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&usage), nil
}

// GetByRequestID 根据请求ID获取用量记录
func (r *usageRepository) GetByRequestID(requestID string) (*dto.Usage, error) {
	var usage model.Usage
	if err := r.db.Where("request_id = ?", requestID).First(&usage).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&usage), nil
}

// List 获取用量记录列表
func (r *usageRepository) List(page, pageSize int) ([]*dto.Usage, int64, error) {
	var total int64
	var usages []model.Usage

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Usage{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&usages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Usage, len(usages))
	for i, u := range usages {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByUser 获取用户的用量记录列表
func (r *usageRepository) ListByUser(userID string, page, pageSize int) ([]*dto.Usage, int64, error) {
	var total int64
	var usages []model.Usage

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Usage{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&usages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Usage, len(usages))
	for i, u := range usages {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByToken 获取令牌的用量记录列表
func (r *usageRepository) ListByToken(tokenID string, page, pageSize int) ([]*dto.Usage, int64, error) {
	var total int64
	var usages []model.Usage

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Usage{}).Where("token_id = ?", tokenID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&usages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Usage, len(usages))
	for i, u := range usages {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByModel 获取模型的用量记录列表
func (r *usageRepository) ListByModel(modelID string, page, pageSize int) ([]*dto.Usage, int64, error) {
	var total int64
	var usages []model.Usage

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Usage{}).Where("model_id = ?", modelID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&usages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Usage, len(usages))
	for i, u := range usages {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByChannel 获取渠道的用量记录列表
func (r *usageRepository) ListByChannel(channelID string, page, pageSize int) ([]*dto.Usage, int64, error) {
	var total int64
	var usages []model.Usage

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Usage{}).Where("channel_id = ?", channelID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&usages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Usage, len(usages))
	for i, u := range usages {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *usageRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行用量记录基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testUsage := &dto.Usage{
			UsageID:          utils.GenerateRandomUUID(12),
			TokenID:          utils.GenerateRandomUUID(12),
			UserID:           utils.GenerateRandomUUID(12),
			ChannelID:        utils.GenerateRandomUUID(12),
			ModelID:          utils.GenerateRandomUUID(12),
			RequestID:        fmt.Sprintf("req_%s", utils.GenerateRandomString(16)),
			UsageType:        []string{"usage", "times"}[rand.Intn(2)],
			UnitPrice:        float64(rand.Intn(100)) / 100000,
			TimesCount:       rand.Intn(100),
			TokensCount:      rand.Intn(1000),
			PriceTotalFactor: float64(rand.Intn(50)+50) / 100,
			TotalAmount:      float64(rand.Intn(10000)) / 100000,
		}

		// 创建
		if err := r.Create(testUsage); err != nil {
			utils.SysError("创建用量记录失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdUsage, err := r.GetByID(testUsage.UsageID)
		if err != nil {
			utils.SysError("获取创建的用量记录失败: " + err.Error())
			return err
		}

		// 更新
		createdUsage.TokensCount = rand.Intn(1000)
		createdUsage.TotalAmount = float64(createdUsage.TokensCount) * createdUsage.UnitPrice * createdUsage.PriceTotalFactor
		if err := r.Update(createdUsage); err != nil {
			utils.SysError("更新用量记录失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdUsage.UsageID); err != nil {
			utils.SysError("删除用量记录失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
