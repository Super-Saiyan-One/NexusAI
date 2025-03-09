package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 用户的配额信息，包括额度、使用情况和有效期等
type Quota struct {
	QuotaID         string          `gorm:"column:quota_id;type:char(36);primaryKey;default:(UUID())" json:"quota_id"`                   // 配额记录唯一标识 uuid
	UserID          string          `gorm:"column:user_id;type:char(36);index;not null;foreignKey:User(UserID)" json:"user_id"`          // 关联的用户ID
	QuotaType       string          `gorm:"column:quota_type;type:enum('topup', 'gift');index;not null" json:"quota_type"`              // 配额类型(充值/赠送/奖励等)
	ValidPeriod     int             `gorm:"column:valid_period;not null;default:0" json:"valid_period"`                                  // 有效期(天)
	Status          int8            `gorm:"column:status;index;not null;default:1" json:"status"`                                        // 配额状态(1:正常 2:冻结 3:过期)
	QuotaAmount     float64         `gorm:"column:quota_amount;type:decimal(10,2);not null;default:0.00" json:"quota_amount"`            // 配额总金额
	RemainingAmount float64         `gorm:"column:remaining_amount;type:decimal(10,2);not null;default:0.00" json:"remaining_amount"`    // 剩余金额
	FrozenAmount    float64         `gorm:"column:frozen_amount;type:decimal(10,2);not null;default:0.00" json:"frozen_amount"`          // 冻结金额(处理中的请求)
	PaymentID       string          `gorm:"column:payment_id;type:char(36);uniqueIndex;foreignKey:Payment(PaymentID)" json:"payment_id"` // 关联的支付ID 充值
	StartTime       utils.MySQLTime `gorm:"column:start_time;not null" json:"start_time"`                                                // 生效时间
	ExpireTime      utils.MySQLTime `gorm:"column:expire_time;not null" json:"expire_time"`                                              // 过期时间
	QuotaOptions    common.JSON     `gorm:"column:quota_options;type:json" json:"quota_options"`                                         // 配额特殊配置

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 表名
func (Quota) TableName() string {
	return "quotas"
}

// BeforeCreate 在创建记录前自动设置时间
func (quota *Quota) BeforeCreate(tx *gorm.DB) error {
	quota.CreatedAt = utils.MySQLTime(utils.GetTime())
	quota.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (quota *Quota) BeforeUpdate(tx *gorm.DB) error {
	quota.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
