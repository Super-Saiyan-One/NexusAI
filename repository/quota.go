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

// QuotaRepository 配额仓储接口
type QuotaRepository interface {
	Create(quota *dto.Quota) error
	Update(quota *dto.Quota) error
	Delete(quotaID string) error
	GetByID(quotaID string) (*dto.Quota, error)
	GetByPaymentID(paymentID string) (*dto.Quota, error)
	List(page, pageSize int) ([]*dto.Quota, int64, error)
	ListByUser(userID string, page, pageSize int) ([]*dto.Quota, int64, error)
	ListByStatus(status int8, page, pageSize int) ([]*dto.Quota, int64, error)
	ListByType(quotaType string, page, pageSize int) ([]*dto.Quota, int64, error)
	Benchmark(count int) error
}

type quotaRepository struct {
	db *gorm.DB
}

// NewQuotaRepository 创建配额仓储实例
func NewQuotaRepository(db *gorm.DB) QuotaRepository {
	return &quotaRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *quotaRepository) convertToDTO(model *model.Quota) *dto.Quota {
	if model == nil {
		return nil
	}

	var quotaOptions dto.QuotaOptions
	if err := model.QuotaOptions.ToStruct(&quotaOptions); err != nil {
		utils.SysError("解析配额选项失败:" + err.Error())
	}

	return &dto.Quota{
		QuotaID:         model.QuotaID,
		UserID:          model.UserID,
		QuotaType:       model.QuotaType,
		ValidPeriod:     model.ValidPeriod,
		Status:          model.Status,
		QuotaAmount:     model.QuotaAmount,
		RemainingAmount: model.RemainingAmount,
		FrozenAmount:    model.FrozenAmount,
		PaymentID:       model.PaymentID,
		StartTime:       model.StartTime,
		ExpireTime:      model.ExpireTime,
		QuotaOptions:    quotaOptions,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		DeletedAt:       utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *quotaRepository) convertToModel(dto *dto.Quota) (*model.Quota, error) {
	if dto == nil {
		return nil, nil
	}

	quotaOptionsJSON, err := common.FromStruct(dto.QuotaOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配额选项失败: %w", err)
	}

	return &model.Quota{
		QuotaID:         dto.QuotaID,
		UserID:          dto.UserID,
		QuotaType:       dto.QuotaType,
		ValidPeriod:     dto.ValidPeriod,
		Status:          dto.Status,
		QuotaAmount:     dto.QuotaAmount,
		RemainingAmount: dto.RemainingAmount,
		FrozenAmount:    dto.FrozenAmount,
		PaymentID:       dto.PaymentID,
		StartTime:       dto.StartTime,
		ExpireTime:      dto.ExpireTime,
		QuotaOptions:    quotaOptionsJSON,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
		DeletedAt:       utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建配额记录
func (r *quotaRepository) Create(quota *dto.Quota) error {
	model, err := r.convertToModel(quota)
	if err != nil {
		return err
	}
	return r.db.Create(model).Error
}

// Update 更新配额记录
func (r *quotaRepository) Update(quota *dto.Quota) error {
	modelData, err := r.convertToModel(quota)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Quota{}).Where("quota_id = ?", quota.QuotaID).Updates(modelData).Error
}

// Delete 删除配额记录
func (r *quotaRepository) Delete(quotaID string) error {
	return r.db.Delete(&model.Quota{}, "quota_id = ?", quotaID).Error
}

// GetByID 根据ID获取配额记录
func (r *quotaRepository) GetByID(quotaID string) (*dto.Quota, error) {
	var quota model.Quota
	if err := r.db.Where("quota_id = ?", quotaID).First(&quota).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&quota), nil
}

// GetByPaymentID 根据支付ID获取配额记录
func (r *quotaRepository) GetByPaymentID(paymentID string) (*dto.Quota, error) {
	var quota model.Quota
	if err := r.db.Where("payment_id = ?", paymentID).First(&quota).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&quota), nil
}

// List 获取配额记录列表
func (r *quotaRepository) List(page, pageSize int) ([]*dto.Quota, int64, error) {
	var total int64
	var quotas []model.Quota

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Quota{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Quota, len(quotas))
	for i, q := range quotas {
		dtoList[i] = r.convertToDTO(&q)
	}

	return dtoList, total, nil
}

// ListByUser 获取用户的配额记录列表
func (r *quotaRepository) ListByUser(userID string, page, pageSize int) ([]*dto.Quota, int64, error) {
	var total int64
	var quotas []model.Quota

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Quota{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Quota, len(quotas))
	for i, q := range quotas {
		dtoList[i] = r.convertToDTO(&q)
	}

	return dtoList, total, nil
}

// ListByStatus 获取指定状态的配额记录列表
func (r *quotaRepository) ListByStatus(status int8, page, pageSize int) ([]*dto.Quota, int64, error) {
	var total int64
	var quotas []model.Quota

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Quota{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Quota, len(quotas))
	for i, q := range quotas {
		dtoList[i] = r.convertToDTO(&q)
	}

	return dtoList, total, nil
}

// ListByType 获取指定类型的配额记录列表
func (r *quotaRepository) ListByType(quotaType string, page, pageSize int) ([]*dto.Quota, int64, error) {
	var total int64
	var quotas []model.Quota

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Quota{}).Where("quota_type = ?", quotaType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Quota, len(quotas))
	for i, q := range quotas {
		dtoList[i] = r.convertToDTO(&q)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *quotaRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行配额记录基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testQuota := &dto.Quota{
			QuotaID:     utils.GenerateRandomUUID(12),
			UserID:      utils.GenerateRandomUUID(12),
			QuotaType:   []string{"recharge", "gift", "reward"}[rand.Intn(3)],
			ValidPeriod: rand.Intn(365) + 1,
			Status:      int8(rand.Intn(3) + 1),
			QuotaAmount: float64(rand.Intn(10000)) / 100,
			PaymentID:   utils.GenerateRandomUUID(12),
			StartTime:   utils.MySQLTime(time.Now()),
			ExpireTime:  utils.MySQLTime(time.Now().AddDate(0, 0, rand.Intn(365)+1)),
			QuotaOptions: dto.QuotaOptions{
				QuotaRemark: fmt.Sprintf("测试配额 %d", i),
			},
		}

		testQuota.FrozenAmount = float64(rand.Intn(int(testQuota.QuotaAmount*100))) / 100
		testQuota.RemainingAmount = testQuota.QuotaAmount - testQuota.FrozenAmount

		// 创建
		if err := r.Create(testQuota); err != nil {
			utils.SysError("创建配额记录失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdQuota, err := r.GetByID(testQuota.QuotaID)
		if err != nil {
			utils.SysError("获取创建的配额记录失败: " + err.Error())
			return err
		}

		// 更新
		createdQuota.Status = int8(rand.Intn(3) + 1)
		if err := r.Update(createdQuota); err != nil {
			utils.SysError("更新配额记录失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdQuota.QuotaID); err != nil {
			utils.SysError("删除配额记录失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
