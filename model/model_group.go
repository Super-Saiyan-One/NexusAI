package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 模型组信息，包括模型组ID、名称、描述、配置和价格系数等
type ModelGroup struct {
	ModelGroupID          string `gorm:"column:model_group_id;type:char(36);primaryKey;default:(UUID())" json:"model_group_id"` // 模型组唯一标识
	ModelGroupName        string `gorm:"column:model_group_name;size:100;not null" json:"model_group_name"`                     // 模型组名称
	ModelGroupDescription string `gorm:"column:model_group_description;type:text" json:"model_group_description"`               // 模型组描述

	ModelGroupPriceFactor common.JSON `gorm:"column:model_group_price_factor;type:json" json:"model_group_price_factor"` // 模型组价格系数
	ModelGroupOptions     common.JSON `gorm:"column:model_group_options;type:json" json:"model_group_options"`           // 模型组配置
	ModelGroupModels      common.JSON `gorm:"column:model_group_models;type:json" json:"model_group_models"`             // 模型组模型

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 模型组创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 模型组信息最后更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 模型组删除时间
}

func (ModelGroup) TableName() string {
	return "model_groups"
}

// BeforeCreate 在创建记录前自动设置时间
func (modelGroup *ModelGroup) BeforeCreate(tx *gorm.DB) error {
	modelGroup.CreatedAt = utils.MySQLTime(utils.GetTime())
	modelGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (modelGroup *ModelGroup) BeforeUpdate(tx *gorm.DB) error {
	modelGroup.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
