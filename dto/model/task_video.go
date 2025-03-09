package model

import "nexus-ai/utils"

type RequestParams []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TaskVideo 视频任务
type ResultURLs struct {
	URL []string `json:"url"`
}

type ErrorLogs struct {
	ErrorLog []string `json:"error_log"`
}

type NextRetryAt struct {
	NextRetryAt []utils.MySQLTime `json:"next_retry_at"`
}

type TaskVideoStatus string

const (
	TaskVideoStatusPending    TaskVideoStatus = "pending"
	TaskVideoStatusProcessing TaskVideoStatus = "processing"
	TaskVideoStatusSuccess    TaskVideoStatus = "success"
	TaskVideoStatusFailed     TaskVideoStatus = "failed"
	TaskVideoStatusRetrying   TaskVideoStatus = "retrying"
)

type TaskVideoProvider string

const (
	TaskVideoProviderLuma      TaskVideoProvider = "luma"
	TaskVideoProviderStability TaskVideoProvider = "stability"
	TaskVideoProviderRunway    TaskVideoProvider = "runway"
)

// TaskVideo 视频任务
type TaskVideo struct {
	TaskVideoID   string            `json:"task_video_id"`
	RequestID     string            `json:"request_id"`
	UserID        string            `json:"user_id"`
	TokenID       string            `json:"token_id"`
	ModelID       string            `json:"model_id"`
	ChannelID     string            `json:"channel_id"`
	Status        TaskVideoStatus   `json:"status"`
	Provider      TaskVideoProvider `json:"provider"`
	CallbackURL   string            `json:"callback_url"`
	RequestParams RequestParams     `json:"request_params"`

	ExternalRequestID string     `json:"external_request_id"`
	ResultURLs        ResultURLs `json:"result_urls"`

	RetryCount  int         `json:"retry_count"`
	ErrorLogs   ErrorLogs   `json:"error_logs"`
	NextRetryAt NextRetryAt `json:"next_retry_at"`

	CreatedAt utils.MySQLTime  `json:"created_at"`
	UpdatedAt utils.MySQLTime  `json:"updated_at"`
	DeletedAt *utils.MySQLTime `json:"deleted_at"`
}
