package dto

import "time"

type TokenSearchRequest struct {
	TokenID   string `json:"token_id"`   // 令牌ID
	UserID    string `json:"user_id"`    // 用户ID
	TokenName string `json:"token_name"` // 令牌名称
	TokenKey  string `json:"token_key"`  // 令牌密钥
	Status    int8   `json:"status"`     // 令牌状态

	MinTokenQuotaTotal  float64 `json:"min_token_quota_total"`  // 令牌配额下限
	MaxTokenQuotaTotal  float64 `json:"max_token_quota_total"`  // 令牌配额上限
	MinTokenQuotaUsed   float64 `json:"min_token_quota_used"`   // 已使用配额下限
	MaxTokenQuotaUsed   float64 `json:"max_token_quota_used"`   // 已使用配额上限
	MinTokenQuotaLeft   float64 `json:"min_token_quota_left"`   // 剩余配额下限
	MaxTokenQuotaLeft   float64 `json:"max_token_quota_left"`   // 剩余配额上限
	MinTokenQuotaFrozen float64 `json:"min_token_quota_frozen"` // 冻结配额下限
	MaxTokenQuotaFrozen float64 `json:"max_token_quota_frozen"` // 冻结配额上限

	ExtraAllowedChannels []string `json:"extra_allowed_channels"` // 额外允许使用的渠道ID列表
	PriorityChannels     []string `json:"priority_channels"`      // 优先使用的渠道ID列表
	AllowedModels        []string `json:"allowed_models"`         // 允许使用的模型ID列表

	MinConcurrentRequests         int `json:"min_concurrent_requests"`            // 最大并发请求数下限
	MaxConcurrentRequests         int `json:"max_concurrent_requests"`            // 最大并发请求数上限
	MinRateLimitRequestsPerMinute int `json:"min_rate_limit_requests_per_minute"` // 每分钟请求数限制下限
	MaxRateLimitRequestsPerMinute int `json:"max_rate_limit_requests_per_minute"` // 每分钟请求数限制上限
	MinRateLimitRequestsPerHour   int `json:"min_rate_limit_requests_per_hour"`   // 每小时请求数限制下限
	MaxRateLimitRequestsPerHour   int `json:"max_rate_limit_requests_per_hour"`   // 每小时请求数限制上限
	MinRateLimitRequestsPerDay    int `json:"min_rate_limit_requests_per_day"`    // 每天请求数限制下限
	MaxRateLimitRequestsPerDay    int `json:"max_rate_limit_requests_per_day"`    // 每天请求数限制上限

	AllowedIPs       []string `json:"allowed_ips"`        // IP白名单
	DisallowedIPs    []string `json:"disallowed_ips"`     // IP黑名单
	RequireSignature bool     `json:"require_signature"`  // 是否要求签名
	DisableRateLimit bool     `json:"disable_rate_limit"` // 是否禁用频率限制
	AvailableLevels  []int    `json:"available_levels"`   // 可用等级

	EarlyCreatedTime time.Time `json:"early_created_time"` // 最早创建时间
	LateCreatedTime  time.Time `json:"late_created_time"`  // 最晚创建时间
	EarlyUpdatedTime time.Time `json:"early_updated_time"` // 最早更新时间
	LateUpdatedTime  time.Time `json:"late_updated_time"`  // 最晚更新时间
	EarlyExpireTime  time.Time `json:"early_expire_time"`  // 最早过期时间
	LateExpireTime   time.Time `json:"late_expire_time"`   // 最晚过期时间
}
