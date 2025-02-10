package model

import (
	"nexus-ai/utils"
)

// ModelGroupPriceFactor 模型组价格系数
type ModelGroupPriceFactor struct {
	RequestPriceFactor    float64 `json:"request_price_factor"`    // 请求价格系数
	ResponsePriceFactor   float64 `json:"response_price_factor"`   // 响应价格系数
	CompletionPriceFactor float64 `json:"completion_price_factor"` // 补全价格系数
	CachePriceFactor      float64 `json:"cache_price_factor"`      // 缓存价格系数
}

// ModelGroupOptions 模型组配置选项
type ModelGroupOptions struct {
	MaxConcurrentRequests int             `json:"max_concurrent_requests"` // 最大并发请求数
	DefaultLevel          int             `json:"default_level"`           // 默认等级
	Discount              float64         `json:"discount"`                // 折扣
	DiscountExpireAt      utils.MySQLTime `json:"discount_expire_at"`      // 折扣过期时间
}

// ModelGroup DTO结构
type ModelGroup struct {
	ModelGroupID          string                `json:"model_group_id"`           // 模型组唯一标识
	ModelGroupName        string                `json:"model_group_name"`         // 模型组名称
	ModelGroupDescription string                `json:"model_group_description"`  // 模型组描述
	ModelGroupPriceFactor ModelGroupPriceFactor `json:"model_group_price_factor"` // 模型组价格系数
	ModelGroupOptions     ModelGroupOptions     `json:"model_group_options"`      // 模型组配置
	CreatedAt             utils.MySQLTime       `json:"created_at"`               // 创建时间
	UpdatedAt             utils.MySQLTime       `json:"updated_at"`               // 更新时间
	DeletedAt             *utils.MySQLTime      `json:"deleted_at,omitempty"`     // 删除时间
}
