package dto

import "time"

type ModelRequest struct {
	Model string `json:"model"`
}

type ModelSearchRequest struct {
	ModelID      string `json:"model_id"`       // 模型ID
	ModelGroupID string `json:"model_group_id"` // 模型组ID
	ModelName    string `json:"model_name"`     // 模型名称
	ModelType    string `json:"model_type"`     // 模型类型
	Provider     string `json:"provider"`       // 模型提供商
	PriceType    string `json:"price_type"`     // 计费类型
	Status       int8   `json:"status"`         // 模型状态

	ModelDisplayName string `json:"model_display_name"` // 模型显示名称
	ModelRequestName string `json:"model_request_name"` // 模型请求名称

	MinAPIDiscount float64 `json:"min_api_discount"` // API折扣下限
	MaxAPIDiscount float64 `json:"max_api_discount"` // API折扣上限

	MinRequestPrice    float64 `json:"min_request_price"`    // 请求价格下限
	MaxRequestPrice    float64 `json:"max_request_price"`    // 请求价格上限
	MinResponsePrice   float64 `json:"min_response_price"`   // 响应价格下限
	MaxResponsePrice   float64 `json:"max_response_price"`   // 响应价格上限
	MinCompletionPrice float64 `json:"min_completion_price"` // 补全价格下限
	MaxCompletionPrice float64 `json:"max_completion_price"` // 补全价格上限
	MinCachePrice      float64 `json:"min_cache_price"`      // 缓存价格下限
	MaxCachePrice      float64 `json:"max_cache_price"`      // 缓存价格上限

	EarlyCreatedTime time.Time `json:"early_created_time"` // 最早创建时间
	LateCreatedTime  time.Time `json:"late_created_time"`  // 最晚创建时间
}
