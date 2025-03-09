package model

import (
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统中所有API调用的用量记录，包括token使用量、计费信息等
type Usage struct {
	UsageID   string `gorm:"column:usage_id;type:char(36);primaryKey;default:(UUID())" json:"usage_id"`                      // 用量记录唯一标识
	TokenID   string `gorm:"column:token_id;type:char(36);index;not null;foreignKey:Token(TokenID)" json:"token_id"`         // 关联的令牌ID
	UserID    string `gorm:"column:user_id;type:char(36);index;not null;foreignKey:User(UserID)" json:"user_id"`             // 关联的用户ID
	ChannelID string `gorm:"column:channel_id;type:char(36);index;not null;foreignKey:Channel(ChannelID)" json:"channel_id"` // 关联的渠道ID
	ModelID   string `gorm:"column:model_id;type:char(36);index;not null;foreignKey:Model(ModelID)" json:"model_id"`         // 关联的模型ID

	RequestID        string  `gorm:"column:request_id;size:64;index" json:"request_id"`                                            // 关联的请求ID，用于日志追踪
	UsageType        string  `gorm:"column:usage_type;size:20;index;not null" json:"usage_type"`                                   // 计费类型(usage:按量,times:按次)
	UnitPrice        float64 `gorm:"column:unit_price;type:decimal(10,6);not null;default:0.000000" json:"unit_price"`             // 单价，支持6位小数
	TimesCount       int     `gorm:"column:times_count;not null;default:0" json:"times_count"`                                     // 请求次数(按次收费时使用)
	TokensCount      int     `gorm:"column:tokens_count;not null;default:0" json:"tokens_count"`                                   // Token数量(按量收费时使用)
	PriceTotalFactor float64 `gorm:"column:price_total_factor;type:decimal(10,2);not null;default:1.00" json:"price_total_factor"` // 价格总倍率
	TotalAmount      float64 `gorm:"column:total_amount;type:decimal(10,6);not null;default:0.000000" json:"total_amount"`         // 总金额，支持6位小数

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 记录创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 记录更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 软删除时间
}

// TableName 表名
func (Usage) TableName() string {
	return "usages"
}

// BeforeCreate 在创建记录前自动设置时间
func (usage *Usage) BeforeCreate(tx *gorm.DB) error {
	usage.CreatedAt = utils.MySQLTime(utils.GetTime())
	usage.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (usage *Usage) BeforeUpdate(tx *gorm.DB) error {
	usage.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
