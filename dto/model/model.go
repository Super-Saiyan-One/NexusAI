package model

import (
	"nexus-ai/dto"
	"nexus-ai/utils"
)

// ModelPrice 模型价格配置
type ModelPrice struct {
	RequestPrice    float64        `json:"request_price"`    // 请求价格
	ResponsePrice   float64        `json:"response_price"`   // 响应价格
	CompletionPrice float64        `json:"completion_price"` // 补全价格
	CachePrice      float64        `json:"cache_price"`      // 缓存价格
	VideoPrice      dto.VideoPrice `json:"video_price"`      // 视频价格
	ImagePrice      dto.ImagePrice `json:"image_price"`      // 图片价格
}

// ModelAlias 模型映射配置
type ModelAlias struct {
	DisplayName string `json:"display_name"` // 显示名称
	RequestName string `json:"request_name"` // 请求名称
}

// ModelOptions 模型配置选项
type ModelOptions struct {
	APIDiscount         float64         `json:"api_discount"`           // API折扣
	APIDiscountExpireAt utils.MySQLTime `json:"api_discount_expire_at"` // API折扣过期时间
}

// Model DTO结构
type Model struct {
	ModelID          string           `json:"model_id"`             // 模型ID
	ModelGroupID     string           `json:"model_group_id"`       // 模型组ID
	ModelName        string           `json:"model_name"`           // 模型名称
	ModelDescription string           `json:"model_description"`    // 模型描述
	ModelType        string           `json:"model_type"`           // 模型类型
	Provider         string           `json:"provider"`             // 模型提供商
	PriceType        string           `json:"price_type"`           // 计费类型
	ModelPrice       ModelPrice       `json:"model_price"`          // 模型价格配置
	Status           int8             `json:"status"`               // 模型状态
	ModelAlias       ModelAlias       `json:"model_alias"`          // 模型映射
	ModelOptions     ModelOptions     `json:"model_options"`        // 模型配置
	CreatedAt        utils.MySQLTime  `json:"created_at"`           // 创建时间
	UpdatedAt        utils.MySQLTime  `json:"updated_at"`           // 更新时间
	DeletedAt        *utils.MySQLTime `json:"deleted_at,omitempty"` // 删除时间
}
