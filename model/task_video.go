package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 视频生成记录表
type TaskVideo struct {
	TaskVideoID   string      `gorm:"column:task_video_id;type:char(36);primaryKey;default:(UUID())" json:"task_video_id"`                                        // 任务ID
	RequestID     string      `gorm:"column:request_id;size:64;not null;unique" json:"request_id"`                                                                // 内部请求ID
	UserID        string      `gorm:"column:user_id;type:char(36);index;not null;foreignKey:User(UserID)" json:"user_id"`                                         // 发起用户ID
	TokenID       string      `gorm:"column:token_id;type:char(36);index;not null;foreignKey:Token(TokenID)" json:"token_id"`                                     // 关联的令牌ID
	ModelID       string      `gorm:"column:model_id;type:char(36);index;not null;foreignKey:Model(ModelID)" json:"model_id"`                                     // 关联的模型ID
	ChannelID     string      `gorm:"column:channel_id;type:char(36);index;not null;foreignKey:Channel(ChannelID)" json:"channel_id"`                             // 关联的渠道ID
	Status        string      `gorm:"column:status;type:enum('pending', 'processing', 'success', 'failed', 'retrying');not null;default:'pending'" json:"status"` // 任务状态
	Provider      string      `gorm:"column:provider;type:enum('luma', 'stability', 'runway');not null" json:"provider"`                                          // 使用的平台
	CallbackURL   string      `gorm:"column:callback_url" json:"callback_url"`                                                                                    // 结果回调地址
	RequestParams common.JSON `gorm:"column:request_params;type:json" json:"request_params"`                                                                      // 请求参数

	ExternalRequestID string      `gorm:"column:external_request_id" json:"external_request_id"` // 外部请求ID
	ResultURLs        common.JSON `gorm:"column:result_urls;type:json" json:"result_urls"`       // 生成的视频URL集合
	RetryCount        int         `gorm:"column:retry_count;default:0" json:"retry_count"`       // 重试次数
	ErrorLogs         common.JSON `gorm:"column:error_logs;type:json" json:"error_logs"`         // 错误日志
	NextRetryAt       common.JSON `gorm:"column:next_retry_at;type:json" json:"next_retry_at"`   // 每次重试时间

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 删除时间
}

// TableName 表名
func (TaskVideo) TableName() string {
	return "task_videos"
}

// BeforeCreate 在创建记录前自动设置时间
func (taskVideo *TaskVideo) BeforeCreate(tx *gorm.DB) error {
	taskVideo.CreatedAt = utils.MySQLTime(utils.GetTime())
	taskVideo.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (taskVideo *TaskVideo) BeforeUpdate(tx *gorm.DB) error {
	taskVideo.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
