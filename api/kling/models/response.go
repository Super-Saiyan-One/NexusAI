package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// ------------------------- 时间处理 -------------------------
type UnixTime time.Time

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}
	*ut = UnixTime(time.UnixMilli(timestamp))
	return nil
}

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ut).UnixMilli())
}

func (ut UnixTime) Time() time.Time {
	return time.Time(ut)
}

// ------------------------- 基础结构 -------------------------
type BaseResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func (r *BaseResponse) Error() string {
	return fmt.Sprintf("API Error %d: %s (request_id: %s)", r.Code, r.Message, r.RequestID)
}

func (r *BaseResponse) IsSuccess() bool {
	return r.Code == ErrCodeSuccess
}

// ------------------------- 创建任务响应 -------------------------
type CreateTaskResponse struct {
	BaseResponse
	Data CreateTaskData `json:"data"`
}

type CreateTaskData struct {
	TaskID     string   `json:"task_id"`
	TaskInfo   TaskInfo `json:"task_info"`
	TaskStatus string   `json:"task_status"`
	CreatedAt  UnixTime `json:"created_at"`
	UpdatedAt  UnixTime `json:"updated_at"`
}

// ------------------------- 查询任务通用响应 -------------------------
type QueryTaskResponse struct {
	BaseResponse
	Data QueryTaskData `json:"data"`
}

type QueryTaskData struct {
	TaskID        string      `json:"task_id"`
	TaskStatus    string      `json:"task_status"`
	TaskStatusMsg string      `json:"task_status_msg,omitempty"`
	TaskInfo      TaskInfo    `json:"task_info"`
	TaskResult    *TaskResult `json:"task_result,omitempty"`
	CreatedAt     UnixTime    `json:"created_at"`
	UpdatedAt     UnixTime    `json:"updated_at"`
}

// ------------------------- 任务结果结构 -------------------------
type TaskInfo struct {
	ExternalTaskID string `json:"external_task_id"`
}

// 统一结果容器（视频和图片独立存在）
type TaskResult struct {
	// 视频生成结果
	Videos []VideoResult `json:"videos,omitempty"`
	// 图片生成结果
	Images []ImageResult `json:"images,omitempty"`
}

type VideoResult struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Duration string `json:"duration"` // 单位：秒
}

type ImageResult struct {
	Index int    `json:"index"` // 图片编号 0-9
	URL   string `json:"url"`
}

// ------------------------- 任务列表响应 -------------------------
type TaskListResponse struct {
	BaseResponse
	Data []TaskListResult `json:"data"` // 直接返回结果数组
}

// 单个任务结果（兼容图片/视频）
type TaskListResult struct {
	TaskID        string      `json:"task_id"`
	TaskStatus    string      `json:"task_status"`
	TaskStatusMsg string      `json:"task_status_msg,omitempty"`
	CreatedAt     UnixTime    `json:"created_at"`
	UpdatedAt     UnixTime    `json:"updated_at"`
	TaskResult    *TaskResult `json:"task_result,omitempty"` // 复用统一结果结构
}

// ------------------------- 常量定义 -------------------------
const (
	TaskStatusSubmitted  = "submitted"
	TaskStatusProcessing = "processing"
	TaskStatusSucceed    = "succeed"
	TaskStatusFailed     = "failed"
)

const (
	ErrCodeSuccess           = 0
	ErrCodeAuthFailed        = 1000
	ErrCodeInvalidParams     = 1200
	ErrCodeContentModeration = 1301
	ErrCodeRateLimit         = 1302
)

// ------------------------- 工具方法 -------------------------
// 获取首个视频URL
func (r *QueryTaskResponse) FirstVideoURL() (string, error) {
	if !r.IsSuccess() || r.Data.TaskResult == nil {
		return "", fmt.Errorf("no video result")
	}

	if len(r.Data.TaskResult.Videos) == 0 {
		return "", fmt.Errorf("empty videos")
	}

	return r.Data.TaskResult.Videos[0].URL, nil
}

// 按索引获取图片
func (r *QueryTaskResponse) GetImage(index int) (*ImageResult, error) {
	if !r.IsSuccess() || r.Data.TaskResult == nil {
		return nil, fmt.Errorf("no image result")
	}

	for _, img := range r.Data.TaskResult.Images {
		if img.Index == index {
			return &img, nil
		}
	}

	return nil, fmt.Errorf("image index %d not found", index)
}
