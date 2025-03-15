package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nexus-ai/api/luma/config"

	"nexus-ai/api/luma/client"
	"nexus-ai/api/luma/model"
)

type ImageService struct {
	client *client.HttpClient
	cfg    *config.Config
}

func NewImageService() (*ImageService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceInit, err)
	}
	return &ImageService{
		client: client.New(cfg.BaseURL, cfg.APIKey),
		cfg:    cfg,
	}, nil
}

func (s *ImageService) GenerateImage(ctx context.Context, req model.ImageGenerateRequest) (*model.GenerationResponse, error) {
	if err := validateImageRequest(req); err != nil {
		return nil, err
	}

	if req.Model == "" {
		req.Model = model.ModelPhoton1
	}
	if req.AspectRatio == "" {
		req.AspectRatio = "16:9"
	}

	resp, err := s.client.DoRequest(
		ctx,
		"POST",
		"/generations/image",
		req,
	)
	println(resp.Status)
	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer client.DrainBody(resp)

	return parseImageResponse(resp)
}

// GetStatus 复用视频服务的查询方法
func (s *ImageService) GetStatus(ctx context.Context, id string) (*model.GenerationResponse, error) {
	resp, err := s.client.DoRequest(ctx, "GET", "/generations/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer client.DrainBody(resp)

	return parseImageResponse(resp)
}

func validateImageRequest(req model.ImageGenerateRequest) error {
	// Validate model
	validModel := false
	for _, m := range model.ValidImageModels {
		if req.Model == m {
			validModel = true
			break
		}
	}
	if !validModel {
		return fmt.Errorf("%w: %s", ErrInvalidModel, req.Model)
	}

	// Validate aspect ratio
	validAspect := false
	for _, ar := range model.ValidAspectRatios {
		if req.AspectRatio == ar {
			validAspect = true
			break
		}
	}
	if !validAspect {
		return fmt.Errorf("%w: aspect ratio %s", ErrInvalidParam, req.AspectRatio)
	}

	// Validate reference weights
	for _, ref := range req.ImageRef {
		if ref.Weight < 0 || ref.Weight > 1 {
			return fmt.Errorf("%w: image_ref weight must be 0-1", ErrInvalidParam)
		}
	}

	return nil
}

func parseImageResponse(resp *http.Response) (*model.GenerationResponse, error) {
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
