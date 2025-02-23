package service

import (
	"context"
	"errors"
	"fmt"
	"nexus-ai/api/hailuo/api"
	"nexus-ai/api/hailuo/model"
	"time"
)

var (
	ErrInvalidModelParams = errors.New("invalid model parameters")
	ErrTaskTimeout        = errors.New("task processing timeout")
)

type VideoService struct {
	client  *api.VideoClient
	timeout time.Duration
}

func NewVideoService(client *api.VideoClient) *VideoService {
	return &VideoService{
		client:  client,
		timeout: 30 * time.Minute,
	}
}

func (s *VideoService) GenerateVideo(ctx context.Context, req model.VideoGenerationRequest) (*model.TaskResponse, error) {
	if err := validateGenerationRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidModelParams, err)
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.CreateTask(ctx, req)
}

func (s *VideoService) GetTaskStatus(ctx context.Context, taskID string) (*model.TaskStatus, error) {
	if taskID == "" {
		return nil, errors.New("empty task ID")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.client.QueryTaskStatus(ctx, taskID)
}

func (s *VideoService) WaitForCompletion(ctx context.Context, taskID string) (*model.TaskStatus, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, ErrTaskTimeout
		case <-ticker.C:
			status, err := s.GetTaskStatus(ctx, taskID)
			if err != nil {
				return nil, err
			}

			switch status.Status {
			case "Success":
				return status, nil
			case "Fail":
				return nil, fmt.Errorf("task failed: %s", status.BaseResp.StatusMsg)
			}
		}
	}
}

func validateGenerationRequest(req model.VideoGenerationRequest) error {
	switch req.Model {
	case "I2V-01", "I2V-01-live":
		if req.FirstFrameImage == "" {
			return errors.New("first frame image required")
		}
	case "S2V-01":
		if len(req.SubjectReference) != 1 {
			return errors.New("exactly one subject reference required")
		}
	case "T2V-01-Director", "T2V-01":
		if req.Prompt == "" {
			return errors.New("prompt required")
		}
	default:
		return errors.New("unsupported model")
	}

	if len(req.Prompt) > 2000 {
		return errors.New("prompt too long")
	}

	return nil
}
