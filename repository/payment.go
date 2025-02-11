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

// PaymentRepository 支付记录仓储接口
type PaymentRepository interface {
	Create(payment *dto.Payment) error
	Update(payment *dto.Payment) error
	Delete(paymentID string) error
	GetByID(paymentID string) (*dto.Payment, error)
	GetByOrderNo(orderNo string) (*dto.Payment, error)
	GetByTransactionID(transactionID string) (*dto.Payment, error)
	List(page, pageSize int) ([]*dto.Payment, int64, error)
	ListByUser(userID string, page, pageSize int) ([]*dto.Payment, int64, error)
	ListByStatus(status int8, page, pageSize int) ([]*dto.Payment, int64, error)
	ListByPlatform(platform string, page, pageSize int) ([]*dto.Payment, int64, error)
	Benchmark(count int) error
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建支付记录仓储实例
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *paymentRepository) convertToDTO(model *model.Payment) *dto.Payment {
	if model == nil {
		return nil
	}

	var paymentInfo dto.PaymentInfo
	var companyInfo dto.CompanyInfo
	var platformOptions dto.PlatformOptions
	var refundInfo dto.RefundInfo

	if err := model.PaymentInfo.ToStruct(&paymentInfo); err != nil {
		utils.SysError("解析支付信息失败:" + err.Error())
	}

	if err := model.CompanyInfo.ToStruct(&companyInfo); err != nil {
		utils.SysError("解析企业信息失败:" + err.Error())
	}

	if err := model.PlatformOptions.ToStruct(&platformOptions); err != nil {
		utils.SysError("解析平台配置失败:" + err.Error())
	}

	if err := model.RefundInfo.ToStruct(&refundInfo); err != nil {
		utils.SysError("解析退款信息失败:" + err.Error())
	}

	var callbackData interface{}
	if err := model.CallbackData.ToStruct(&callbackData); err != nil {
		utils.SysError("解析回调数据失败:" + err.Error())
	}

	return &dto.Payment{
		PaymentID:            model.PaymentID,
		UserID:               model.UserID,
		PaymentPlatform:      model.PaymentPlatform,
		PaymentScene:         model.PaymentScene,
		PaymentMethod:        model.PaymentMethod,
		PaymentCurrency:      model.PaymentCurrency,
		PaymentAmount:        model.PaymentAmount,
		PaymentStatus:        model.PaymentStatus,
		PaymentType:          model.PaymentType,
		PaymentOrderNo:       model.PaymentOrderNo,
		PaymentTransactionID: model.PaymentTransactionID,
		PaymentTitle:         model.PaymentTitle,
		PaymentDescription:   model.PaymentDescription,
		NotifyURL:            model.NotifyURL,
		PaymenterName:        model.PaymenterName,
		PaymenterEmail:       model.PaymenterEmail,
		PaymenterPhone:       model.PaymenterPhone,
		PaymentInfo:          paymentInfo,
		CompanyInfo:          companyInfo,
		CallbackData:         callbackData,
		PlatformOptions:      platformOptions,
		PaymentTime:          model.PaymentTime,
		ExpireTime:           model.ExpireTime,
		RefundInfo:           refundInfo,
		CreatedAt:            model.CreatedAt,
		UpdatedAt:            model.UpdatedAt,
		DeletedAt:            utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *paymentRepository) convertToModel(dto *dto.Payment) (*model.Payment, error) {
	if dto == nil {
		return nil, nil
	}

	paymentInfoJSON, err := common.FromStruct(dto.PaymentInfo)
	if err != nil {
		return nil, fmt.Errorf("转换支付信息失败: %w", err)
	}

	companyInfoJSON, err := common.FromStruct(dto.CompanyInfo)
	if err != nil {
		return nil, fmt.Errorf("转换企业信息失败: %w", err)
	}

	callbackDataJSON, err := common.FromStruct(dto.CallbackData)
	if err != nil {
		return nil, fmt.Errorf("转换回调数据失败: %w", err)
	}

	platformOptionsJSON, err := common.FromStruct(dto.PlatformOptions)
	if err != nil {
		return nil, fmt.Errorf("转换平台配置失败: %w", err)
	}

	refundInfoJSON, err := common.FromStruct(dto.RefundInfo)
	if err != nil {
		return nil, fmt.Errorf("转换退款信息失败: %w", err)
	}

	return &model.Payment{
		PaymentID:            dto.PaymentID,
		UserID:               dto.UserID,
		PaymentPlatform:      dto.PaymentPlatform,
		PaymentScene:         dto.PaymentScene,
		PaymentMethod:        dto.PaymentMethod,
		PaymentCurrency:      dto.PaymentCurrency,
		PaymentAmount:        dto.PaymentAmount,
		PaymentStatus:        dto.PaymentStatus,
		PaymentType:          dto.PaymentType,
		PaymentOrderNo:       dto.PaymentOrderNo,
		PaymentTransactionID: dto.PaymentTransactionID,
		PaymentTitle:         dto.PaymentTitle,
		PaymentDescription:   dto.PaymentDescription,
		NotifyURL:            dto.NotifyURL,
		PaymenterName:        dto.PaymenterName,
		PaymenterEmail:       dto.PaymenterEmail,
		PaymenterPhone:       dto.PaymenterPhone,
		PaymentInfo:          paymentInfoJSON,
		CompanyInfo:          companyInfoJSON,
		CallbackData:         callbackDataJSON,
		PlatformOptions:      platformOptionsJSON,
		PaymentTime:          dto.PaymentTime,
		ExpireTime:           dto.ExpireTime,
		RefundInfo:           refundInfoJSON,
		CreatedAt:            dto.CreatedAt,
		UpdatedAt:            dto.UpdatedAt,
		DeletedAt:            utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建支付记录
func (r *paymentRepository) Create(payment *dto.Payment) error {
	model, err := r.convertToModel(payment)
	if err != nil {
		return err
	}
	return r.db.Create(model).Error
}

// Update 更新支付记录
func (r *paymentRepository) Update(payment *dto.Payment) error {
	modelData, err := r.convertToModel(payment)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Payment{}).Where("payment_id = ?", payment.PaymentID).Updates(modelData).Error
}

// Delete 删除支付记录
func (r *paymentRepository) Delete(paymentID string) error {
	return r.db.Delete(&model.Payment{}, "payment_id = ?", paymentID).Error
}

// GetByID 根据ID获取支付记录
func (r *paymentRepository) GetByID(paymentID string) (*dto.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("payment_id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&payment), nil
}

// GetByOrderNo 根据订单号获取支付记录
func (r *paymentRepository) GetByOrderNo(orderNo string) (*dto.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("payment_order_no = ?", orderNo).First(&payment).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&payment), nil
}

// GetByTransactionID 根据交易ID获取支付记录
func (r *paymentRepository) GetByTransactionID(transactionID string) (*dto.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("payment_transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&payment), nil
}

// List 获取支付记录列表
func (r *paymentRepository) List(page, pageSize int) ([]*dto.Payment, int64, error) {
	var total int64
	var payments []model.Payment

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Payment, len(payments))
	for i, p := range payments {
		dtoList[i] = r.convertToDTO(&p)
	}

	return dtoList, total, nil
}

// ListByUser 获取用户的支付记录列表
func (r *paymentRepository) ListByUser(userID string, page, pageSize int) ([]*dto.Payment, int64, error) {
	var total int64
	var payments []model.Payment

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Payment{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Payment, len(payments))
	for i, p := range payments {
		dtoList[i] = r.convertToDTO(&p)
	}

	return dtoList, total, nil
}

// ListByStatus 获取指定状态的支付记录列表
func (r *paymentRepository) ListByStatus(status int8, page, pageSize int) ([]*dto.Payment, int64, error) {
	var total int64
	var payments []model.Payment

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Payment{}).Where("payment_status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Payment, len(payments))
	for i, p := range payments {
		dtoList[i] = r.convertToDTO(&p)
	}

	return dtoList, total, nil
}

// ListByPlatform 获取指定平台的支付记录列表
func (r *paymentRepository) ListByPlatform(platform string, page, pageSize int) ([]*dto.Payment, int64, error) {
	var total int64
	var payments []model.Payment

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Payment{}).Where("payment_platform = ?", platform)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.Payment, len(payments))
	for i, p := range payments {
		dtoList[i] = r.convertToDTO(&p)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *paymentRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行支付记录基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testPayment := &dto.Payment{
			PaymentID:       utils.GenerateRandomUUID(12),
			UserID:          utils.GenerateRandomUUID(12),
			PaymentPlatform: []string{"stripe", "paddle", "alipay", "wechat"}[rand.Intn(4)],
			PaymentScene:    []string{"recharge", "subscription"}[rand.Intn(2)],
			PaymentMethod:   []string{"credit_card", "bank_transfer", "alipay", "wechat"}[rand.Intn(4)],
			PaymentCurrency: []string{"USD", "CNY", "HKD", "EUR"}[rand.Intn(4)],
			PaymentAmount:   float64(rand.Intn(1000000)) / 100,
			PaymentStatus:   int8(rand.Intn(5) + 1),
			PaymentType:     []string{"personal", "enterprise"}[rand.Intn(2)],
			PaymentOrderNo:  fmt.Sprintf("ORDER_%s", utils.GenerateRandomString(16)),
			PaymentTitle:    fmt.Sprintf("测试支付 %d", i),
			PaymentInfo: dto.PaymentInfo{
				InvoiceTitle: fmt.Sprintf("发票抬头 %d", i),
				InvoiceType:  []string{"common", "special"}[rand.Intn(2)],
				TaxNumber:    utils.GenerateRandomString(15),
			},
			CompanyInfo: dto.CompanyInfo{
				CompanyName:     fmt.Sprintf("测试公司 %d", i),
				BusinessLicense: utils.GenerateRandomString(18),
				ContactPerson:   fmt.Sprintf("联系人 %d", i),
			},
			PlatformOptions: dto.PlatformOptions{
				AppID:       utils.GenerateRandomString(16),
				MerchantID:  utils.GenerateRandomString(12),
				Environment: "sandbox",
			},
		}

		// 创建
		if err := r.Create(testPayment); err != nil {
			utils.SysError("创建支付记录失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdPayment, err := r.GetByID(testPayment.PaymentID)
		if err != nil {
			utils.SysError("获取创建的支付记录失败: " + err.Error())
			return err
		}

		// 更新
		createdPayment.PaymentStatus = int8(rand.Intn(5) + 1)
		if err := r.Update(createdPayment); err != nil {
			utils.SysError("更新支付记录失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdPayment.PaymentID); err != nil {
			utils.SysError("删除支付记录失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
