package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统中所有AI对话消息记录
type MessageSave struct {
	MessageID string `gorm:"column:message_id;type:char(36);primaryKey;default:(UUID())" json:"message_id"`                  // 消息唯一标识
	RequestID string `gorm:"column:request_id;size:64;index" json:"request_id"`                                              // 关联的请求ID
	UserID    string `gorm:"column:user_id;type:char(36);index;not null;foreignKey:User(UserID)" json:"user_id"`             // 关联的用户ID
	TokenID   string `gorm:"column:token_id;type:char(36);index;not null;foreignKey:Token(TokenID)" json:"token_id"`         // 关联的令牌ID
	ModelID   string `gorm:"column:model_id;type:char(36);index;not null;foreignKey:Model(ModelID)" json:"model_id"`         // 关联的模型ID
	ChannelID string `gorm:"column:channel_id;type:char(36);index;not null;foreignKey:Channel(ChannelID)" json:"channel_id"` // 关联的渠道ID

	RequestTitle     string      `gorm:"column:request_title;size:255" json:"request_title"`                   // 请求标题
	RequestContent   string      `gorm:"column:request_content;type:text;not null" json:"request_content"`     // 请求内容
	ResponseContent  string      `gorm:"column:response_content;type:text" json:"response_content"`            // 响应内容
	Input            string      `gorm:"column:input;type:text" json:"input"`                                  // 输入
	Instruction      string      `gorm:"column:instruction;type:text" json:"instruction"`                      // 指令
	Prompt           string      `gorm:"column:prompt;type:text" json:"prompt"`                                // 提示词
	Stream           int8        `gorm:"column:stream;not null;default:0" json:"stream"`                       // 是否流式输出
	StreamOptions    common.JSON `gorm:"column:stream_options;type:json" json:"stream_options"`                // 流式输出配置
	RequestTokens    int         `gorm:"column:request_tokens;not null;default:0" json:"request_tokens"`       // 请求token数
	PromptTokens     int         `gorm:"column:prompt_tokens;not null;default:0" json:"prompt_tokens"`         // 提示词token数
	CompletionTokens int         `gorm:"column:completion_tokens;not null;default:0" json:"completion_tokens"` // 补全token数
	Latency          int         `gorm:"column:latency;not null;default:0" json:"latency"`                     // 响应延迟(ms)
	MessageOptions   common.JSON `gorm:"column:message_options;type:json" json:"message_options"`              // 消息配置(温度/top_p等)
	MessageExtra     common.JSON `gorm:"column:message_extra;type:json" json:"message_extra"`                  // 消息额外信息

	MessageStatus int8   `gorm:"column:message_status;index;not null;default:1" json:"message_status"` // 消息状态(1:正常 2:中断 3:异常)
	RetryCount    int    `gorm:"column:retry_count;not null;default:0" json:"retry_count"`             // 重试次数
	ErrorType     string `gorm:"column:error_type;size:50" json:"error_type"`                          // 错误类型
	ErrorInfo     string `gorm:"column:error_info;type:text" json:"error_info"`                        // 错误信息

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 表名
func (MessageSave) TableName() string {
	return "message_saves"
}

// BeforeCreate 在创建记录前自动设置时间
func (messageSave *MessageSave) BeforeCreate(tx *gorm.DB) error {
	messageSave.CreatedAt = utils.MySQLTime(utils.GetTime())
	messageSave.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (messageSave *MessageSave) BeforeUpdate(tx *gorm.DB) error {
	messageSave.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
