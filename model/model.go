package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统支持的AI模型信息，包括定价、配置和服务提供商等信息
type Model struct {
	ModelID          string `gorm:"column:model_id;type:char(36);primaryKey;default:(UUID())" json:"model_id"` // 模型ID
	ModelGroupID     string `gorm:"column:model_group_id;type:char(36);not null;index" json:"model_group_id"`  // 模型组ID
	ModelName        string `gorm:"column:model_name;size:100;uniqueIndex;not null" json:"model_name"`         // 模型名称
	ModelDescription string `gorm:"column:model_description;type:text" json:"model_description"`               // 模型描述
	ModelType        string `gorm:"column:model_type;size:50;index;not null" json:"model_type"`                // 模型类型 文本/语音/图像/视频/多模态 ...
	Provider         string `gorm:"column:provider;size:50;index;not null" json:"provider"`                    // 模型提供商

	PriceType    string      `gorm:"column:price_type;size:20;not null" json:"price_type"` // 计费类型
	ModelPrice   common.JSON `gorm:"column:model_price;type:json" json:"model_price"`      // 模型价格配置 请求 补全 响应
	Status       int8        `gorm:"column:status;index;not null;default:1" json:"status"` // 模型状态 1:正常 0:禁用
	ModelAlias   common.JSON `gorm:"column:model_alias;type:json" json:"model_alias"`      // 模型映射
	ModelOptions common.JSON `gorm:"column:model_options;type:json" json:"model_options"`  // 模型配置

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 表名
func (Model) TableName() string {
	return "models"
}

// BeforeCreate 在创建记录前自动设置时间
func (model *Model) BeforeCreate(tx *gorm.DB) error {
	model.CreatedAt = utils.MySQLTime(utils.GetTime())
	model.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (model *Model) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
