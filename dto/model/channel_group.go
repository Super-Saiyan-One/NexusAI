package model

import (
	"nexus-ai/utils"
)

// ChannelGroupPriceFactor 渠道组价格系数
type ChannelGroupPriceFactor struct {
	RequestPriceFactor    float64 `json:"request_price_factor"`    // 请求价格系数
	ResponsePriceFactor   float64 `json:"response_price_factor"`   // 响应价格系数
	CompletionPriceFactor float64 `json:"completion_price_factor"` // 补全价格系数
	CachePriceFactor      float64 `json:"cache_price_factor"`      // 缓存价格系数
}

// ChannelGroupOptions 渠道组配置选项
type ChannelGroupOptions struct {
	MaxConcurrentRequests int             `json:"max_concurrent_requests"` // 最大并发请求数
	DefaultLevel          int             `json:"default_level"`           // 默认等级
	Discount              float64         `json:"discount"`                // 折扣
	DiscountExpireAt      utils.MySQLTime `json:"discount_expire_at"`      // 折扣过期时间
}

// ChannelGroupChannels 渠道组渠道
type ChannelGroupChannels struct {
	Channels  []string            `json:"channels"`   // 渠道
	ModelsMap map[string][]string `json:"models_map"` // 模型可用渠道表
}

// ChannelGroup DTO结构
type ChannelGroup struct {
	ChannelGroupID          string                  `json:"channel_group_id"`           // 渠道组ID
	ChannelGroupName        string                  `json:"channel_group_name"`         // 渠道组名称
	ChannelGroupDescription string                  `json:"channel_group_description"`  // 渠道组描述
	ChannelGroupPriceFactor ChannelGroupPriceFactor `json:"channel_group_price_factor"` // 渠道组价格系数
	ChannelGroupOptions     ChannelGroupOptions     `json:"channel_group_options"`      // 渠道组配置
	ChannelGroupChannels    ChannelGroupChannels    `json:"channel_group_channels"`     // 渠道组下属渠道

	CreatedAt utils.MySQLTime  `json:"created_at"`           // 创建时间
	UpdatedAt utils.MySQLTime  `json:"updated_at"`           // 更新时间
	DeletedAt *utils.MySQLTime `json:"deleted_at,omitempty"` // 删除时间
}
