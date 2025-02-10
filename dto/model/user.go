package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"
)

type OAuthProvider struct {
	ProviderID           string          `json:"provider_id"`            // 第三方登录提供商ID
	ProviderType         string          `json:"provider_type"`          // 第三方登录提供商类型
	ProviderName         string          `json:"provider_name"`          // 第三方登录提供商名称
	ProviderIcon         string          `json:"provider_icon"`          // 第三方登录提供商图标
	ProviderURL          string          `json:"provider_url"`           // 第三方登录提供商URL
	ProviderAccessToken  string          `json:"provider_access_token"`  // 第三方登录提供商的访问令牌
	ProviderRefreshToken string          `json:"provider_refresh_token"` // 第三方登录提供商的刷新令牌
	ProviderExpireTime   utils.MySQLTime `json:"provider_expire_time"`   // 第三方登录提供商的令牌过期时间
	ProviderScope        string          `json:"provider_scope"`         // 第三方登录提供商的权限范围
	ProviderUserID       string          `json:"provider_user_id"`       // 第三方登录提供商的用户ID
	ProviderUserName     string          `json:"provider_user_name"`     // 第三方登录提供商的用户名
	ProviderUserEmail    string          `json:"provider_user_email"`    // 第三方登录提供商的用户邮箱
	ProviderUserAvatar   string          `json:"provider_user_avatar"`   // 第三方登录提供商的用户头像
	ProviderUserPhone    string          `json:"provider_user_phone"`    // 第三方登录提供商的用户手机号
	ProviderUserGender   string          `json:"provider_user_gender"`   // 第三方登录提供商的用户性别
	ProviderUserBirthday utils.MySQLTime `json:"provider_user_birthday"` // 第三方登录提供商的用户生日
	ProviderUserLocation string          `json:"provider_user_location"` // 第三方登录提供商的用户位置
	ProviderUserLanguage string          `json:"provider_user_language"` // 第三方登录提供商的用户语言
	ProviderUserStatus   int8            `json:"provider_user_status"`   // 第三方登录提供商的用户状态
	ProviderOtherInfo    common.JSON     `json:"provider_other_info"`    // 第三方登录提供商的其他信息
}

type OAuthInfo struct {
	Google    OAuthProvider `json:"google"`    // Google登录信息
	GitHub    OAuthProvider `json:"github"`    // GitHub登录信息
	Discord   OAuthProvider `json:"discord"`   // Discord登录信息
	Apple     OAuthProvider `json:"apple"`     // Apple登录信息
	Telegram  OAuthProvider `json:"telegram"`  // Telegram登录信息
	Microsoft OAuthProvider `json:"microsoft"` // Microsoft登录信息
	Facebook  OAuthProvider `json:"facebook"`  // Facebook登录信息
	QQ        OAuthProvider `json:"qq"`        // QQ登录信息
	WX        OAuthProvider `json:"wx"`        // 微信登录信息
	X         OAuthProvider `json:"x"`         // X登录信息
}

type UserQuota struct {
	TotalQuota      float64         `json:"total_quota"`      // 总配额
	FrozenQuota     float64         `json:"frozen_quota"`     // 冻结配额
	GiftQuota       float64         `json:"gift_quota"`       // 赠送配额
	LeftQuota       float64         `json:"left_quota"`       // 剩余配额
	ExceedQuota     float64         `json:"exceed_quota"`     // 超额配额
	NotifyThreshold float64         `json:"notify_threshold"` // 通知阈值
	LimitQuota      float64         `json:"limit_quota"`      // 限制配额
	LimitExpire     utils.MySQLTime `json:"limit_expire"`     // 限制过期时间
}

type UserOptions struct {
	MaxConcurrentRequests    int      `json:"max_concurrent_requests"`     // 最大并发请求数
	DefaultLevel             int      `json:"default_level"`               // 默认等级
	ExtraAllowedModels       []string `json:"extra_allowed_models"`        // 额外允许使用的模型列表
	ExtraAllowedChannels     []string `json:"extra_allowed_channels"`      // 额外允许使用的渠道列表
	APIDiscount              float64  `json:"api_discount"`                // API折扣
	UserAvatar               string   `json:"user_avatar"`                 // 用户头像
	UserBackground           string   `json:"user_background"`             // 用户背景
	UserLanguage             string   `json:"user_language"`               // 用户语言
	QuotaNotifyEmail         string   `json:"quota_notify_email"`          // 配额通知邮件
	ReceiveQuotaNotifyMail   int8     `json:"receive_quota_notify_mail"`   // 是否接收配额邮件通知
	BillingNotifyEmail       string   `json:"billing_notify_email"`        // 账单通知邮件
	ReceiveBillingNotifyMail int8     `json:"receive_billing_notify_mail"` // 是否接收账单邮件通知
	PaymentNotifyEmail       string   `json:"payment_notify_email"`        // 支付通知邮件
	ReceivePaymentNotifyMail int8     `json:"receive_payment_notify_mail"` // 是否接收支付邮件通知
}

type User struct {
	UserID        string           `json:"user_id"`              // 用户唯一标识
	UserGroupID   string           `json:"user_group_id"`        // 用户组ID
	Username      string           `json:"username"`             // 用户名
	Password      string           `json:"password"`             // 密码
	Email         string           `json:"email"`                // 邮箱
	Phone         string           `json:"phone"`                // 手机号
	OAuthInfo     OAuthInfo        `json:"oauth_info"`           // 第三方登录相关的认证信息(如GitHub、微信、Google等)
	UserQuota     UserQuota        `json:"user_quota"`           // 各种资源使用配额信息
	UserOptions   UserOptions      `json:"user_options"`         // 个性化设置和偏好选项
	LastLoginTime utils.MySQLTime  `json:"last_login_time"`      // 最后一次成功登录的时间
	LastLoginIP   string           `json:"last_login_ip"`        // 最后一次登录的IP地址
	Status        int8             `json:"status"`               // 用户状态，1:正常 0:禁用
	CreatedAt     utils.MySQLTime  `json:"created_at"`           // 用户账号的创建时间
	UpdatedAt     utils.MySQLTime  `json:"updated_at"`           // 用户信息的最后更新时间
	DeletedAt     *utils.MySQLTime `json:"deleted_at,omitempty"` // 删除时间
}
