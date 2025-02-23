package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"nexus-ai/api/kling/models"
	service2 "nexus-ai/api/kling/service"
	"strconv"
	"strings"
)

type ImageController struct {
	service service2.ImageService
}

func NewImageController(s service2.ImageService) *ImageController {
	return &ImageController{service: s}
}

// CreateImageTask 处理创建图像任务请求 POST /v1/images/generations
func (c *ImageController) CreateImageTask(w http.ResponseWriter, r *http.Request) {
	var req models.TextToImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := validateImageRequest(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := c.service.CreateImageTask(r.Context(), req)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// GetImageTask 处理查询单个任务请求 GET /v1/images/generations/{taskId}
func (c *ImageController) GetImageTask(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	taskID := ""
	if len(pathSegments) >= 4 {
		taskID = pathSegments[3]
	}

	if taskID == "" {
		respondError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	resp, err := c.service.GetImageTask(r.Context(), taskID)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ListImageTasks 处理任务列表查询 GET /v1/images/generations
func (c *ImageController) ListImageTasks(w http.ResponseWriter, r *http.Request) {
	pageNum, _ := strconv.Atoi(r.URL.Query().Get("pageNum"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	// 参数安全处理
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 30
	}

	resp, err := c.service.ListImageTasks(r.Context(), pageNum, pageSize)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// 私有方法
func (c *ImageController) handleServiceError(w http.ResponseWriter, err error) {
	var serviceErr *service2.ServiceError

	if errors.As(err, &serviceErr) {
		respondJSON(w, serviceErr.HTTPStatus(), map[string]interface{}{
			"code":    serviceErr.Code,
			"message": serviceErr.Message,
		})
		return
	}

	respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"code":    5000,
		"message": "Internal server error",
	})
}

func validateImageRequest(req models.TextToImageRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	if len(req.Prompt) > 500 {
		return fmt.Errorf("prompt exceeds maximum length of 500 characters")
	}

	if req.NegativePrompt != "" && len(req.NegativePrompt) > 200 {
		return fmt.Errorf("negative prompt exceeds maximum length of 200 characters")
	}

	if req.Image != "" {
		if _, err := url.ParseRequestURI(req.Image); err != nil {
			if !isValidBase64(req.Image) {
				return fmt.Errorf("invalid image format, must be URL or base64")
			}
		}
	}

	if req.ImageFidelity < 0 || req.ImageFidelity > 1 {
		return fmt.Errorf("image fidelity must be between 0 and 1")
	}

	if req.N < 1 || req.N > 9 {
		return fmt.Errorf("n must be between 1 and 9")
	}

	validAspectRatios := map[string]bool{
		"16:9": true, "9:16": true, "1:1": true,
		"4:3": true, "3:4": true, "3:2": true, "2:3": true,
	}
	if req.AspectRatio != "" && !validAspectRatios[req.AspectRatio] {
		return fmt.Errorf("invalid aspect ratio")
	}

	if req.CallbackURL != "" {
		if _, err := url.ParseRequestURI(req.CallbackURL); err != nil {
			return fmt.Errorf("invalid callback URL format")
		}
	}

	return nil
}

func isValidBase64(s string) bool {
	// 简化版base64验证
	return !strings.ContainsAny(s, " \t\r\n") &&
		len(s)%4 == 0 &&
		!strings.Contains(s, "data:image")
}
