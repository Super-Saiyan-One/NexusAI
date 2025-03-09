package model

import (
	"nexus-ai/utils"
)

// QuotaOptions 配额特殊配置
type QuotaOptions struct {
	QuotaRemark string `json:"quota_remark"` // 配额备注
}

type QuotaType string

const (
	QuotaTypeTopup QuotaType = "topup"
	QuotaTypeGift  QuotaType = "gift"
)

// Quota DTO结构
type Quota struct {
	QuotaID         string           `json:"quota_id"`             // 配额记录唯一标识
	UserID          string           `json:"user_id"`              // 关联的用户ID
	QuotaType       QuotaType        `json:"quota_type"`           // 配额类型
	ValidPeriod     int              `json:"valid_period"`         // 有效期
	Status          int8             `json:"status"`               // 配额状态
	QuotaAmount     float64          `json:"quota_amount"`         // 配额总金额
	RemainingAmount float64          `json:"remaining_amount"`     // 剩余金额
	FrozenAmount    float64          `json:"frozen_amount"`        // 冻结金额
	PaymentID       string           `json:"payment_id"`           // 关联的支付ID
	StartTime       utils.MySQLTime  `json:"start_time"`           // 生效时间
	ExpireTime      utils.MySQLTime  `json:"expire_time"`          // 过期时间
	QuotaOptions    QuotaOptions     `json:"quota_options"`        // 配额特殊配置
	CreatedAt       utils.MySQLTime  `json:"created_at"`           // 创建时间
	UpdatedAt       utils.MySQLTime  `json:"updated_at"`           // 更新时间
	DeletedAt       *utils.MySQLTime `json:"deleted_at,omitempty"` // 删除时间
}
