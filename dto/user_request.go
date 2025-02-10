package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱，必填，格式需正确
	Password string `json:"password" binding:"required"`    // 密码，必填
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱，必填
	Password string `json:"password" binding:"required"`    // 密码，必填
}

type PasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码，必填
	NewPassword string `json:"new_password" binding:"required"` // 新密码，必填
}

type SearchRequest struct {
	UserID                string    `json:"user_id"`                 // 用户ID
	UserGroupID           string    `json:"user_group_id"`           // 用户组ID
	Username              string    `json:"username"`                // 用户名
	Email                 string    `json:"email"`                   // 邮箱
	Phone                 string    `json:"phone"`                   // 手机号
	Levels                []int     `json:"levels"`                  // 等级
	Status                int       `json:"status"`                  // 状态
	MinConcurrentRequests int       `json:"min_concurrent_requests"` // 最大并发请求数下限
	MaxConcurrentRequests int       `json:"max_concurrent_requests"` // 最大并发请求数上限
	MinAPIDiscount        float64   `json:"min_api_discount"`        // API折扣下限
	MaxAPIDiscount        float64   `json:"max_api_discount"`        // API折扣上限
	MinQuota              int       `json:"min_quota"`               // 配额下限
	MaxQuota              int       `json:"max_quota"`               // 配额上限
	EarlyLastLoginTime    time.Time `json:"early_last_login_time"`   // 最早登录时间
	LateLastLoginTime     time.Time `json:"late_last_login_time"`    // 最晚登录时间
	EarlyCreatedTime      time.Time `json:"early_created_time"`      // 最早创建时间
	LateCreatedTime       time.Time `json:"late_created_time"`       // 最晚创建时间
}
