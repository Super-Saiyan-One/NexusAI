package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

type UserGroup struct {
	UserGroupID          string `gorm:"column:user_group_id;type:char(36);primaryKey;default:(UUID())" json:"user_group_id"` // 用户组唯一标识
	UserGroupName        string `gorm:"column:user_group_name;size:100;uniqueIndex;not null" json:"user_group_name"`         // 用户组名称
	UserGroupDescription string `gorm:"column:user_group_description;type:text" json:"user_group_description"`               // 用户组描述

	UserGroupPriceFactor common.JSON `gorm:"column:user_group_price_factor;type:json" json:"user_group_price_factor"` // 用户组价格系数
	UserGroupOptions     common.JSON `gorm:"column:user_group_options;type:json" json:"user_group_options"`           // 用户组配置
	UserGroupUsers       common.JSON `gorm:"column:user_group_users;type:json" json:"user_group_users"`               // 用户组用户
	UserGroupModels      common.JSON `gorm:"column:user_group_models;type:json" json:"user_group_models"`             // 用户组模型

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 删除时间
}

func (UserGroup) TableName() string {
	return "user_groups"
}

// BeforeCreate 在创建记录前自动设置时间
func (userGroup *UserGroup) BeforeCreate(tx *gorm.DB) error {
	userGroup.CreatedAt = utils.MySQLTime(utils.GetTime())
	userGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (userGroup *UserGroup) BeforeUpdate(tx *gorm.DB) error {
	userGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
