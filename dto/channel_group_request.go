package dto

import "time"

type ChannelGroupSearchRequest struct {
	ChannelGroupID   string `json:"channel_group_id"`   // 渠道组ID
	ChannelGroupName string `json:"channel_group_name"` // 渠道组名称
	Levels           []int  `json:"levels"`             // 等级

	Channels []string `json:"channels"` // 渠道列表
	Models   []string `json:"models"`   // 模型列表

	MinAPIDiscount        float64 `json:"min_api_discount"`        // API折扣下限
	MaxAPIDiscount        float64 `json:"max_api_discount"`        // API折扣上限
	MinConcurrentRequests int     `json:"min_concurrent_requests"` // 最大并发请求数下限
	MaxConcurrentRequests int     `json:"max_concurrent_requests"` // 最大并发请求数上限

	MinRequestPriceFactor    float64 `json:"min_request_price_factor"`    // 请求价格系数下限
	MaxRequestPriceFactor    float64 `json:"max_request_price_factor"`    // 请求价格系数上限
	MinResponsePriceFactor   float64 `json:"min_response_price_factor"`   // 响应价格系数下限
	MaxResponsePriceFactor   float64 `json:"max_response_price_factor"`   // 响应价格系数上限
	MinCompletionPriceFactor float64 `json:"min_completion_price_factor"` // 补全价格系数下限
	MaxCompletionPriceFactor float64 `json:"max_completion_price_factor"` // 补全价格系数上限
	MinCachePriceFactor      float64 `json:"min_cache_price_factor"`      // 缓存价格系数下限
	MaxCachePriceFactor      float64 `json:"max_cache_price_factor"`      // 缓存价格系数上限

	EarlyCreatedTime time.Time `json:"early_created_time"` // 最早创建时间
	LateCreatedTime  time.Time `json:"late_created_time"`  // 最晚创建时间
}
