package service

import (
	"fmt"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"
)

type TokenService interface {
	TokenCreate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error)
	TokenUpdate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error)
	TokenSearch(repo repository.TokenRepository, tokenID string, tokenKey string, userID string) ([]*dto.Token, error)
	TokenDelete(repo repository.TokenRepository, tokenID string) error
}

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (ts *tokenService) TokenCreate(repo repository.TokenRepository, token *dto.Token) (*dto.Token, error) {
	token.TokenID = utils.GenerateRandomUUID(12)
	token.TokenKey = utils.GenerateRandomString(32)
	token.Status = 1

	if token.TokenQuotaTotal < 0 {
		return nil, fmt.Errorf("token quota cannot be negative")
	}
	token.TokenQuotaLeft = token.TokenQuotaTotal
	token.TokenQuotaUsed = 0
	token.TokenQuotaFrozen = 0

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

func (ts *tokenService) TokenSearch(repo repository.TokenRepository, tokenID string, tokenKey string, userID string) ([]*dto.Token, error) {
	if tokenID != "" {
		token, err := repo.GetByID(tokenID)
		if err != nil {
			return nil, err
		}
		return []*dto.Token{token}, nil
	}
	if tokenKey != "" {
		token, err := repo.GetByKey(tokenKey)
		if err != nil {
			return nil, err
		}
		return []*dto.Token{token}, nil
	}
	if userID != "" {
		tokens, err := repo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		return tokens, nil
	}
	return nil, nil
}

func (ts *tokenService) TokenDelete(repo repository.TokenRepository, tokenID string) error {
	err := repo.Delete(tokenID)
	if err != nil {
		return err
	}
	return nil
}
