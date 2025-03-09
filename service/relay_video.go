package service

import (
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type RelayVideoService interface {
	Text2Video(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
	Image2Video(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
	VideoExtend(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
	LipSync(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
}

type relayVideoService struct{}

func NewRelayVideoService() RelayVideoService {
	return &relayVideoService{}
}

func (rvs *relayVideoService) Text2Video(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	// TODO
	return nil, nil
}

func (rvs *relayVideoService) Image2Video(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	// TODO
	return nil, nil
}

func (rvs *relayVideoService) VideoExtend(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	// TODO
	return nil, nil
}

func (rvs *relayVideoService) LipSync(repo repository.TaskVideoRepository, taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	// TODO
	return nil, nil
}
