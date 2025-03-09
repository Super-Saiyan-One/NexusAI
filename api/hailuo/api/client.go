package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nexus-ai/api/hailuo/auth"
	"nexus-ai/api/hailuo/config"
	"nexus-ai/api/hailuo/model"
	"time"
)

const (
	baseURL          = "https://api.minimax.chat/v1"
	createTaskPath   = "/video_generation"
	queryTaskPath    = "/query/video_generation"
	retrieveFilePath = "/file"
)

type VideoClient struct {
	httpClient *http.Client
}

func NewVideoClient(auth *auth.Authenticator) *VideoClient {
	return &VideoClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *VideoClient) CreateTask(ctx context.Context, req model.VideoGenerationRequest) (*model.TaskResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+createTaskPath, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", config.API_KEY)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taskResp model.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &taskResp, nil
}

func (c *VideoClient) QueryTaskStatus(ctx context.Context, taskID string) (*model.TaskStatus, error) {
	url := fmt.Sprintf("%s%s?task_id=%s", baseURL, queryTaskPath, taskID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", config.API_KEY)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var status model.TaskStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &status, nil
}

func (c *VideoClient) RetrieveFile(ctx context.Context, fileID string) (*model.FileResponse, error) {
	url := fmt.Sprintf("%s%s/%s", baseURL, retrieveFilePath, fileID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", config.API_KEY)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var fileResp model.FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &fileResp, nil
}
