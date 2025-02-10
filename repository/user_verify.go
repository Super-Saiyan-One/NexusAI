package repository

import (
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
)

func UserVerify(accessToken string, refreshToken string, userID string) (*dto.User, error) {
	userRepo := NewUserRepository(model.GetDB())
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
