package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexus-ai/api/kling/models"
	"nexus-ai/api/kling/service"
	"strings"
)

type VideoController struct {
	service service.VideoService
}

func NewVideoController(s service.VideoService) *VideoController {
	return &VideoController{service: s}
}

// CreateVideoTask 处理创建视频任务请求
func (c *VideoController) CreateVideoTask(w http.ResponseWriter, r *http.Request) {
	var req models.TextToVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := validateCreateRequest(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := c.service.CreateTextVideoTask(r.Context(), req)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// GetVideoTask 处理查询任务 GET /v1/videos/{taskId} 或 GET /v1/videos?external_task_id=
func (c *VideoController) GetVideoTask(w http.ResponseWriter, r *http.Request) {
	// 从URL路径获取taskId
	pathSegments := strings.Split(r.URL.Path, "/")
	taskID := ""
	if len(pathSegments) >= 4 {
		taskID = pathSegments[3]
	}

	// 从查询参数获取external_task_id
	externalTaskID := r.URL.Query().Get("external_task_id")

	// 参数校验
	if taskID == "" && externalTaskID == "" {
		respondError(w, http.StatusBadRequest, "必须提供 task_id 或 external_task_id 参数")
		return
	}

	// 调用服务层
	resp, err := c.service.GetVideoTask(r.Context(), taskID, externalTaskID, false)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// 私有错误处理方法
func (c *VideoController) handleServiceError(w http.ResponseWriter, err error) {
	var serviceErr *service.ServiceError

	// 使用 errors.As 处理错误链
	if errors.As(err, &serviceErr) {
		statusCode := serviceErr.HTTPStatus()
		respondJSON(w, statusCode, map[string]interface{}{
			"code":    serviceErr.Code,
			"message": serviceErr.Message,
		})
		return
	}

	// 未知错误处理
	respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"code":    5000,
		"message": "Internal server error",
	})
}

// controller/video_controller.go
func (c *VideoController) CreateImageToVideoTask(w http.ResponseWriter, r *http.Request) {
	var req models.ImageToVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// 自定义验证逻辑
	if err := validateImageToVideoRequest(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := c.service.CreateImageToVideoTask(r.Context(), req)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// validateImageToVideoRequest 自定义验证逻辑
func validateImageToVideoRequest(req models.ImageToVideoRequest) error {
	// 检查图像参数互斥
	if req.Image != "" && req.ImageTail != "" {
		return errors.New("image and image_tail cannot be used together")
	}

	// 检查三组参数互斥
	groupCount := 0
	if req.Image != "" || req.ImageTail != "" {
		groupCount++
	}
	if req.StaticMask != "" || len(req.DynamicMasks) > 0 {
		groupCount++
	}
	if req.CameraControl != nil {
		groupCount++
	}
	if groupCount > 1 {
		return errors.New("image/image_tail, masks and camera_control are mutually exclusive")
	}

	// 检查动态笔刷轨迹点坐标
	for _, dm := range req.DynamicMasks {
		for _, t := range dm.Trajectories {
			if t.X < 0 || t.Y < 0 {
				return errors.New("trajectory coordinates must be positive")
			}
		}
	}

	return nil
}

// 参数验证函数
func validateCreateRequest(req models.TextToVideoRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if len(req.Prompt) > 2500 {
		return fmt.Errorf("prompt exceeds maximum length of 2500 characters")
	}
	return nil
}

// 统一响应方法
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
