// controller/video_controller.go
package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexus-ai/api/hailuo/model"
	"nexus-ai/api/hailuo/service"
	"time"
)

type VideoController struct {
	service *service.VideoService
}

func NewVideoController(s *service.VideoService) *VideoController {
	return &VideoController{service: s}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// CreateVideoTask godoc
// @Summary Create video generation task
// @Accept  json
// @Produce json
// @Param   request body model.VideoGenerationRequest true "Task parameters"
// @Success 201 {object} model.TaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [post]
func (c *VideoController) CreateVideoTask(w http.ResponseWriter, r *http.Request) {
	var req model.VideoGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := c.service.GenerateVideo(r.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidModelParams) {
			statusCode = http.StatusBadRequest
		}
		respondError(w, statusCode, "Failed to create task", err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// GetTaskStatus godoc
// @Summary Get task status
// @Produce json
// @Param   task_id path string true "Task ID"
// @Success 200 {object} model.TaskStatus
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tasks/{task_id} [get]
func (c *VideoController) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	if taskID == "" {
		respondError(w, http.StatusBadRequest, "Missing task ID", nil)
		return
	}

	status, err := c.service.GetTaskStatus(r.Context(), taskID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, context.DeadlineExceeded) {
			statusCode = http.StatusGatewayTimeout
		}
		respondError(w, statusCode, "Failed to get task status", err)
		return
	}

	respondJSON(w, http.StatusOK, status)
}

// HandleCallback 处理异步回调
func (c *VideoController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	var cbReq model.CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&cbReq); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid callback data", err)
		return
	}

	// 处理挑战请求
	if cbReq.Challenge != "" {
		respondJSON(w, http.StatusOK, model.CallbackChallenge{Challenge: cbReq.Challenge})
		return
	}

	// 处理状态更新
	switch cbReq.Status {
	case "processing":
		fmt.Printf("[%s] Task %s processing\n", time.Now().Format(time.RFC3339), cbReq.TaskID)
	case "success":
		fmt.Printf("[%s] Task %s completed\n", time.Now().Format(time.RFC3339), cbReq.TaskID)
	case "failed":
		fmt.Printf("[%s] Task %s failed: %s\n", time.Now().Format(time.RFC3339), cbReq.TaskID, cbReq.Status)
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string, err error) {
	resp := ErrorResponse{Error: message}
	if err != nil {
		resp.Details = err.Error()
	}
	respondJSON(w, status, resp)
}
