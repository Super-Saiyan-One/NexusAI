package model

import "nexus-ai/utils"

// TokenChannels 令牌渠道配置
type TokenChannels struct {
	ExtraAllowedChannels []string `json:"extra_allowed_channels"` // 额外允许使用的渠道ID列表
	PriorityChannels     []string `json:"priority_channels"`      // 优先使用的渠道ID列表
	DefaultTestChannel   string   `json:"default_test_channel"`   // 默认测试渠道ID
}

// TokenModels 令牌模型配置
type TokenModels struct {
	AllowedModels    []string `json:"allowed_models"`     // 允许使用的模型ID列表
	DefaultTestModel string   `json:"default_test_model"` // 默认测试模型ID
}

// TokenOptions 令牌配置选项
type TokenOptions struct {
	MaxConcurrentRequests int      `json:"max_concurrent_requests"` // 最大并发请求数
	MaxRequestsPerMinute  int      `json:"max_requests_per_minute"` // 每分钟最大请求数
	MaxRequestsPerHour    int      `json:"max_requests_per_hour"`   // 每小时最大请求数
	MaxRequestsPerDay     int      `json:"max_requests_per_day"`    // 每天最大请求数
	AllowedIPs            []string `json:"allowed_ips"`             // IP白名单
	DisallowedIPs         []string `json:"disallowed_ips"`          // IP黑名单
	RequireSignature      bool     `json:"require_signature"`       // 是否要求签名
	DisableRateLimit      bool     `json:"disable_rate_limit"`      // 是否禁用频率限制
	AvailableLevels       []int    `json:"available_levels"`        // 可用等级
}

// Token DTO结构
type Token struct {
	TokenID   string `json:"token_id"`   // 令牌唯一标识
	UserID    string `json:"user_id"`    // 关联的用户ID
	TokenName string `json:"token_name"` // 令牌名称
	TokenKey  string `json:"token_key"`  // 令牌密钥
	Status    int8   `json:"status"`     // 令牌状态

	TokenQuotaTotal  float64 `json:"token_quota_total"`  // 令牌配额
	TokenQuotaUsed   float64 `json:"token_quota_used"`   // 已使用配额
	TokenQuotaLeft   float64 `json:"token_quota_left"`   // 剩余配额
	TokenQuotaFrozen float64 `json:"token_quota_frozen"` // 冻结配额

	TokenChannels TokenChannels    `json:"token_channels"`       // 令牌的渠道配置
	TokenModels   TokenModels      `json:"token_models"`         // 支持的模型列表配置
	TokenOptions  TokenOptions     `json:"token_options"`        // 令牌的特殊配置选项
	ExpireTime    utils.MySQLTime  `json:"expire_time"`          // 令牌过期时间
	CreatedAt     utils.MySQLTime  `json:"created_at"`           // 令牌创建时间
	UpdatedAt     utils.MySQLTime  `json:"updated_at"`           // 令牌信息最后更新时间
	DeletedAt     *utils.MySQLTime `json:"deleted_at,omitempty"` // 软删除时间
}
