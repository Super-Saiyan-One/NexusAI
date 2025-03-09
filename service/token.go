package service

import (
	"fmt"
	tokenDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
)

type TokenService interface {
	TokenCreate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error)
	TokenUpdate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error)
	TokenSearch(repo repository.TokenRepository, tokenSearch *tokenDto.TokenSearchRequest) ([]*dto.Token, error)
	TokenDelete(repo repository.TokenRepository, tokenID string) error
	TokenAvailableModels(repo repository.TokenRepository, tokenID string) ([]*dto.Model, error)
	TokenAvailableChannels(repo repository.TokenRepository, tokenID string) ([]*dto.Channel, error)
}

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (ts *tokenService) TokenCreate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error) {
	if token.TokenQuotaTotal < 0 {
		return nil, fmt.Errorf("token quota cannot be negative")
	}

	createdToken, err := repo.Create(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	return createdToken, nil
}

func (ts *tokenService) TokenUpdate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error) {
	existingToken, err := repo.GetByID(token.TokenID)
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}

	existingToken.TokenName = token.TokenName
	existingToken.Status = token.Status
	existingToken.TokenOptions = token.TokenOptions
	existingToken.TokenChannels = token.TokenChannels
	existingToken.TokenModels = token.TokenModels

	if token.TokenQuotaTotal >= 0 {
		quotaDiff := token.TokenQuotaTotal - existingToken.TokenQuotaTotal
		existingToken.TokenQuotaTotal = token.TokenQuotaTotal
		existingToken.TokenQuotaLeft += quotaDiff
	}

	updatedToken, err := repo.Update(existingToken)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}
	return updatedToken, nil
}

func (ts *tokenService) TokenSearch(repo repository.TokenRepository, tokenSearch *tokenDto.TokenSearchRequest) ([]*dto.Token, error) {
	tokens, _, err := repo.Search(tokenSearch)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (ts *tokenService) TokenDelete(repo repository.TokenRepository, tokenID string) error {
	err := repo.Delete(tokenID)
	if err != nil {
		return err
	}
	return nil
}

// TokenAvailableModels 获取token可用模型
func (ts *tokenService) TokenAvailableModels(repo repository.TokenRepository, tokenID string) ([]*dto.Model, error) {
	// TODO: 获取token可用模型
	return []*dto.Model{}, nil
}

// TokenAvailableChannels 获取token可用渠道
func (ts *tokenService) TokenAvailableChannels(repo repository.TokenRepository, tokenID string) ([]*dto.Channel, error) {
	// TODO: 获取token可用渠道
	return []*dto.Channel{}, nil
}
