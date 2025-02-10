package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 渠道组信息，包括渠道组ID、名称、描述、价格系数和配置等
type ChannelGroup struct {
	ChannelGroupID          string `gorm:"column:channel_group_id;type:char(36);primaryKey;default:(UUID())" json:"channel_group_id"` // 渠道组ID
	ChannelGroupName        string `gorm:"column:channel_group_name;size:100;not null" json:"channel_group_name"`                     // 渠道组名称
	ChannelGroupDescription string `gorm:"column:channel_group_description;type:text" json:"channel_group_description"`               // 渠道组描述

	ChannelGroupPriceFactor common.JSON `gorm:"column:channel_group_price_factor;type:json" json:"channel_group_price_factor"` // 渠道组价格系数
	ChannelGroupOptions     common.JSON `gorm:"column:channel_group_options;type:json" json:"channel_group_options"`           // 渠道组配置
	ChannelGroupChannels    common.JSON `gorm:"column:channel_group_channels;type:json" json:"channel_group_channels"`         // 渠道组渠道

	CreatedAt utils.MySQLTime `gorm:"column:created_at;index;not null" json:"created_at"` // 创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"`       // 更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ChannelGroup) TableName() string {
	return "channel_groups"
}

// BeforeCreate 在创建记录前自动设置时间
func (channelGroup *ChannelGroup) BeforeCreate(tx *gorm.DB) error {
	channelGroup.CreatedAt = utils.MySQLTime(utils.GetTime())
	channelGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (channelGroup *ChannelGroup) BeforeUpdate(tx *gorm.DB) error {
	channelGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
