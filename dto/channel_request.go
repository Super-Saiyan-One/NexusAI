package dto

import "time"

type ChannelSearchRequest struct {
	ChannelID      string `json:"channel_id"`       // 渠道ID
	ChannelGroupID string `json:"channel_group_id"` // 渠道组ID
	ChannelName    string `json:"channel_name"`     // 渠道名称
	Status         int8   `json:"status"`           // 渠道状态

	AllowedModels []string `json:"allowed_models"` // 允许使用的模型ID列表

	UpstreamEndpoint       string `json:"upstream_endpoint"`         // 上游服务端点
	UpstreamProxyURL       string `json:"upstream_proxy_url"`        // 上游服务代理URL
	MinUpstreamTimeout     int    `json:"min_upstream_timeout"`      // 上游服务超时时间(秒)下限
	MaxUpstreamTimeout     int    `json:"max_upstream_timeout"`      // 上游服务超时时间(秒)上限
	MinUpstreamMaxRetries  int    `json:"min_upstream_max_retries"`  // 上游服务最大重试次数下限
	MaxUpstreamMaxRetries  int    `json:"max_upstream_max_retries"`  // 上游服务最大重试次数上限
	MinUpstreamDialTimeout int    `json:"min_upstream_dial_timeout"` // 上游服务连接超时时间(秒)下限
	MaxUpstreamDialTimeout int    `json:"max_upstream_dial_timeout"` // 上游服务连接超时时间(秒)上限

	AuthAPIKey      string `json:"auth_api_key"`      // API密钥
	AuthAPISecret   string `json:"auth_api_secret"`   // API密钥对应的secret
	AuthBearerToken string `json:"auth_bearer_token"` // Bearer令牌

	MinRetryMaxRetries      int   `json:"min_retry_max_retries"`       // 重试最大次数下限
	MaxRetryMaxRetries      int   `json:"max_retry_max_retries"`       // 重试最大次数上限
	MinRetryInterval        int   `json:"min_retry_interval"`          // 重试间隔(毫秒)下限
	MaxRetryInterval        int   `json:"max_retry_interval"`          // 重试间隔(毫秒)上限
	MinRetryMaxRetryBackoff int   `json:"min_retry_max_retry_backoff"` // 重试最大退避时间(毫秒)下限
	MaxRetryMaxRetryBackoff int   `json:"max_retry_max_retry_backoff"` // 重试最大退避时间(毫秒)上限
	RetryStatuses           []int `json:"retry_statuses"`              // 需要重试的HTTP状态码

	MinRateLimitRequestsPerSecond int `json:"min_rate_limit_requests_per_second"` // 每秒请求数限制下限
	MaxRateLimitRequestsPerSecond int `json:"max_rate_limit_requests_per_second"` // 每秒请求数限制上限
	MinRateLimitRequestsPerMinute int `json:"min_rate_limit_requests_per_minute"` // 每分钟请求数限制下限
	MaxRateLimitRequestsPerMinute int `json:"max_rate_limit_requests_per_minute"` // 每分钟请求数限制上限
	MinRateLimitRequestsPerHour   int `json:"min_rate_limit_requests_per_hour"`   // 每小时请求数限制下限
	MaxRateLimitRequestsPerHour   int `json:"max_rate_limit_requests_per_hour"`   // 每小时请求数限制上限
	MinRateLimitRequestsPerDay    int `json:"min_rate_limit_requests_per_day"`    // 每天请求数限制下限
	MaxRateLimitRequestsPerDay    int `json:"max_rate_limit_requests_per_day"`    // 每天请求数限制上限

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
