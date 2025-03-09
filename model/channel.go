package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统中所有AI服务渠道的配置信息，包括上游服务配置、认证信息、重试策略等
type Channel struct {
	ChannelID          string `gorm:"column:channel_id;type:char(36);primaryKey;default:(UUID())" json:"channel_id"` // 渠道唯一标识
	ChannelGroupID     string `gorm:"column:channel_group_id;type:char(36);not null;index" json:"channel_group_id"`  // 渠道组ID
	ChannelName        string `gorm:"column:channel_name;size:100;not null" json:"channel_name"`                     // 渠道名称
	ChannelDescription string `gorm:"column:channel_description;type:text" json:"channel_description"`               // 渠道描述信息
	Status             int8   `gorm:"column:status;index;not null;default:1" json:"status"`                          // 渠道状态，1:正常 0:禁用

	ChannelModels      common.JSON `gorm:"column:channel_models;type:json" json:"channel_models"`             // 支持的模型列表配置
	ChannelPriceFactor common.JSON `gorm:"column:channel_price_factor;type:json" json:"channel_price_factor"` // 渠道价格系数
	UpstreamOptions    common.JSON `gorm:"column:upstream_options;type:json" json:"upstream_options"`         // 上游服务配置信息(endpoint/timeout等)
	AuthOptions        common.JSON `gorm:"column:auth_options;type:json" json:"auth_options"`                 // 认证配置信息(key/secret等)
	RetryOptions       common.JSON `gorm:"column:retry_options;type:json" json:"retry_options"`               // 重试策略配置
	RateLimit          common.JSON `gorm:"column:rate_limit;type:json" json:"rate_limit"`                     // 速率限制配置
	ModelMapping       common.JSON `gorm:"column:model_mapping;type:json" json:"model_mapping"`               // 模型映射配置，用于模型名称转换
	TestModels         common.JSON `gorm:"column:test_models;type:json" json:"test_models"`                   // 测试模型配置

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 渠道创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 渠道信息最后更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 渠道删除时间
}

// TableName 表名
func (Channel) TableName() string {
	return "channels"
}

// BeforeCreate 在创建记录前自动设置时间
func (channel *Channel) BeforeCreate(tx *gorm.DB) error {
	channel.CreatedAt = utils.MySQLTime(utils.GetTime())
	channel.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (channel *Channel) BeforeUpdate(tx *gorm.DB) error {
	channel.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
