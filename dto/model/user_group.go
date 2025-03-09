package model

import (
	"nexus-ai/dto"
	"nexus-ai/utils"
)

// 用户组价格系数
type UserGroupPriceFactor struct {
	RequestPriceFactor    float64        `json:"request_price_factor"`    // 请求价格系数
	ResponsePriceFactor   float64        `json:"response_price_factor"`   // 响应价格系数
	CompletionPriceFactor float64        `json:"completion_price_factor"` // 补全价格系数
	CachePriceFactor      float64        `json:"cache_price_factor"`      // 缓存价格系数
	VideoPriceFactor      dto.VideoPrice `json:"video_price_factor"`      // 视频价格系数
	ImagePriceFactor      dto.ImagePrice `json:"image_price_factor"`      // 图片价格系数
}

// 用户组配置选项
type UserGroupOptions struct {
	MaxConcurrentRequests int      `json:"max_concurrent_requests"` // 最大并发请求数
	DefaultLevel          int      `json:"default_level"`           // 默认等级
	ExtraAllowedModels    []string `json:"extra_allowed_models"`    // 额外允许使用的模型列表
	ExtraAllowedChannels  []string `json:"extra_allowed_channels"`  // 额外允许使用的渠道列表
	APIDiscount           float64  `json:"api_discount"`            // API折扣
}

// 用户组用户
type UserGroupUsers struct {
	Users []string `json:"users"` // 用户组用户
}

// 用户组模型
type UserGroupModels struct {
	Models []string `json:"models"` // 用户组模型
}

// 用户组 DTO结构
type UserGroup struct {
	UserGroupID          string               `json:"user_group_id"`           // 用户组唯一标识
	UserGroupName        string               `json:"user_group_name"`         // 用户组名称
	UserGroupDescription string               `json:"user_group_description"`  // 用户组描述
	UserGroupPriceFactor UserGroupPriceFactor `json:"user_group_price_factor"` // 用户组价格系数
	UserGroupOptions     UserGroupOptions     `json:"user_group_options"`      // 用户组配置
	UserGroupUsers       UserGroupUsers       `json:"user_group_users"`        // 用户组用户
	UserGroupModels      UserGroupModels      `json:"user_group_models"`       // 用户组模型
	CreatedAt            utils.MySQLTime      `json:"created_at"`              // 创建时间
	UpdatedAt            utils.MySQLTime      `json:"updated_at"`              // 更新时间
	DeletedAt            *utils.MySQLTime     `json:"deleted_at,omitempty"`    // 删除时间
}
