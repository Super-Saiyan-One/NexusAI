package model

import (
	"nexus-ai/dto"
	"nexus-ai/utils"
)

// ChannelModels 渠道支持的模型配置
type ChannelModels struct {
	AllowedModels    []string `json:"allowed_models"`     // 允许使用的模型ID列表
	DefaultTestModel string   `json:"default_test_model"` // 默认测试模型ID
}

// ChannelPriceFactor 渠道价格系数
type ChannelPriceFactor struct {
	RequestPriceFactor    float64        `json:"request_price_factor"`    // 请求价格系数
	ResponsePriceFactor   float64        `json:"response_price_factor"`   // 响应价格系数
	CompletionPriceFactor float64        `json:"completion_price_factor"` // 补全价格系数
	CachePriceFactor      float64        `json:"cache_price_factor"`      // 缓存价格系数
	VideoPriceFactor      dto.VideoPrice `json:"video_price_factor"`      // 视频价格系数
	ImagePriceFactor      dto.ImagePrice `json:"image_price_factor"`      // 图片价格系数
}

// UpstreamOptions 上游服务配置
type UpstreamOptions struct {
	Endpoint    string `json:"endpoint"`     // 服务端点
	Timeout     int    `json:"timeout"`      // 超时时间(秒)
	MaxRetries  int    `json:"max_retries"`  // 最大重试次数
	ProxyURL    string `json:"proxy_url"`    // 代理URL
	DialTimeout int    `json:"dial_timeout"` // 连接超时时间(秒)
}

// AuthOptions 认证配置
type AuthOptions struct {
	APIKey        string            `json:"api_key"`        // API密钥
	APISecret     string            `json:"api_secret"`     // API密钥对应的secret
	BearerToken   string            `json:"bearer_token"`   // Bearer令牌
	Headers       map[string]string `json:"headers"`        // 自定义认证头
	RequestParams map[string]string `json:"request_params"` // 请求参数
}

// RetryOptions 重试策略配置
type RetryOptions struct {
	MaxRetries      int   `json:"max_retries"`       // 最大重试次数
	RetryInterval   int   `json:"retry_interval"`    // 重试间隔(毫秒)
	MaxRetryBackoff int   `json:"max_retry_backoff"` // 最大重试退避时间(毫秒)
	RetryStatuses   []int `json:"retry_statuses"`    // 需要重试的HTTP状态码
}

// RateLimit 速率限制配置
type RateLimit struct {
	RequestsPerSecond int `json:"requests_per_second"` // 每秒请求数限制
	RequestsPerMinute int `json:"requests_per_minute"` // 每分钟请求数限制
	RequestsPerHour   int `json:"requests_per_hour"`   // 每小时请求数限制
	RequestsPerDay    int `json:"requests_per_day"`    // 每天请求数限制
}

// ModelMapping 模型映射配置
type ModelMapping struct {
	MappingModels map[string][]string `json:"mapping_models"` // 模型映射列表
}

// TestModels 测试模型配置
type TestModels struct {
	TestModels       []string `json:"test_models"`        // 测试模型列表
	DefaultTestModel string   `json:"default_test_model"` // 默认测试模型
}

// Channel DTO结构
type Channel struct {
	ChannelID          string             `json:"channel_id"`           // 渠道唯一标识
	ChannelGroupID     string             `json:"channel_group_id"`     // 渠道组ID
	ChannelName        string             `json:"channel_name"`         // 渠道名称
	ChannelDescription string             `json:"channel_description"`  // 渠道描述信息
	Status             int8               `json:"status"`               // 渠道状态
	ChannelModels      ChannelModels      `json:"channel_models"`       // 支持的模型列表配置
	ChannelPriceFactor ChannelPriceFactor `json:"channel_price_factor"` // 渠道价格系数
	UpstreamOptions    UpstreamOptions    `json:"upstream_options"`     // 上游服务配置信息
	AuthOptions        AuthOptions        `json:"auth_options"`         // 认证配置信息
	RetryOptions       RetryOptions       `json:"retry_options"`        // 重试策略配置
	RateLimit          RateLimit          `json:"rate_limit"`           // 速率限制配置
	ModelMapping       ModelMapping       `json:"model_mapping"`        // 模型映射配置
	TestModels         TestModels         `json:"test_models"`          // 测试模型配置
	CreatedAt          utils.MySQLTime    `json:"created_at"`           // 创建时间
	UpdatedAt          utils.MySQLTime    `json:"updated_at"`           // 更新时间
	DeletedAt          *utils.MySQLTime   `json:"deleted_at,omitempty"` // 删除时间
}
