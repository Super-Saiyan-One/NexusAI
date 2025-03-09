package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nexus-ai/api/luma/client"
	"nexus-ai/api/luma/config"
	"nexus-ai/api/luma/model"
)

var (
	ErrServiceInit  = errors.New("service initialization failed")
	ErrInvalidModel = errors.New("invalid model")
	ErrInvalidParam = errors.New("invalid parameter")
	ErrAPIError     = errors.New("api error")
)

type VideoService struct {
	client *client.HttpClient
	cfg    *config.Config
}

// NewVideoService 自动初始化配置的服务构造函数
func NewVideoService() (*VideoService, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceInit, err)
	}

	// 初始化客户端
	httpClient := client.New(cfg.BaseURL, cfg.APIKey)

	return &VideoService{
		client: httpClient,
		cfg:    cfg,
	}, nil
}

func (s *VideoService) GenerateVideo(ctx context.Context, req model.VideoGenerateRequest) (*model.GenerationResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	// 设置固定参数
	req.GenerationType = "video"
	if req.AspectRatio == "" {
		req.AspectRatio = "16:9"
	}

	resp, err := s.client.DoRequest(ctx, "POST", "/generations", req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer client.DrainBody(resp)

	return parseResponse(resp)
}

func (s *VideoService) GetStatus(ctx context.Context, id string) (*model.GenerationResponse, error) {
	resp, err := s.client.DoRequest(ctx, "GET", "/generations/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer client.DrainBody(resp)

	return parseResponse(resp)
}

// 私有方法
func validateRequest(req model.VideoGenerateRequest) error {
	// 校验模型类型
	validModel := false
	for _, m := range model.ValidModels {
		if req.Model == m {
			validModel = true
			break
		}
	}
	if !validModel {
		return fmt.Errorf("%w: %s", ErrInvalidModel, req.Model)
	}

	// 校验分辨率
	if req.Resolution != nil {
		validRes := false
		for _, r := range model.ValidResolutions {
			if *req.Resolution == r {
				validRes = true
				break
			}
		}
		if !validRes {
			return fmt.Errorf("%w: %s", ErrInvalidParam, *req.Resolution)
		}
	}

	// 校验回调URL
	if req.CallbackURL != nil && *req.CallbackURL != "" {
		if _, err := url.ParseRequestURI(*req.CallbackURL); err != nil {
			return fmt.Errorf("%w: invalid callback URL", ErrInvalidParam)
		}
	}

	return nil
}

func parseResponse(resp *http.Response) (*model.GenerationResponse, error) {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrAPIError, string(body))
	}

	var response model.GenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &response, nil
}
