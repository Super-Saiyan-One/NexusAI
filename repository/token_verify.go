package repository

import (
	"fmt"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"time"
)

func TokenVerify(tokenKey string) (*dto.Token, error) {
	tokenRepo := NewTokenRepository(model.GetDB())
	token, err := tokenRepo.GetByKey(tokenKey)
	if err != nil {
		return nil, fmt.Errorf("token not found")
	}

	// 检查令牌是否已被删除
	if token.DeletedAt != nil {
		return nil, fmt.Errorf("token has been deleted")
	}

	// 检查令牌状态是否禁用 (假设状态0为禁用)
	if token.Status == 0 {
		return nil, fmt.Errorf("token is disabled")
	}

	// 检查令牌是否过期
	if !token.ExpireTime.IsZero() && time.Time(token.ExpireTime).Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return token, nil
}
