package repository

import (
	"fmt"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
)

func QuotaVerify(userID string, tokenQuotaLeft float64) (*dto.User, bool, error) {
	userRepo := NewUserRepository(model.GetDB())
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return nil, false, fmt.Errorf("user not found")
	}
	if tokenQuotaLeft < constant.MinimumQuota {
		return nil, false, fmt.Errorf("token quota left is less than minimum needed quota")
	}
	userQuotaLeft := user.UserQuota.LeftQuota
	if userQuotaLeft < constant.MinimumQuota {
		return nil, false, fmt.Errorf("token owner's quota left is less than minimum needed quota")
	}
	return user, true, nil
}
