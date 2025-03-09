package service

import (
	"fmt"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"
	"time"
)

type QuotaService interface {
	TopupQuota(repo repository.QuotaRepository, req *dto.Quota) error
	GiftQuota(repo repository.QuotaRepository, req *dto.Quota) error
	FreezeQuota(repo repository.QuotaRepository, quotaID string, amount float64) error
	ExpireQuota(repo repository.QuotaRepository, quotaID string, expireTime time.Time) error
	CostQuota(repo repository.QuotaRepository, userID string, amount float64) error
}

type quotaService struct{}

func NewQuotaService() QuotaService {
	return &quotaService{}
}

func (qs *quotaService) TopupQuota(repo repository.QuotaRepository, req *dto.Quota) error {
	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}
	if req.QuotaType != dto.QuotaTypeTopup {
		return fmt.Errorf("quota type must be topup")
	}
	if !req.ExpireTime.IsZero() && req.ExpireTime.Before(time.Now()) {
		return fmt.Errorf("expire time must be greater than current time")
	}
	if req.QuotaAmount <= 0 {
		return fmt.Errorf("quota amount must be greater than 0")
	}
	req.QuotaID = utils.GenerateRandomUUID(12)
	req.Status = 1
	req.RemainingAmount = req.QuotaAmount
	req.FrozenAmount = 0

	err := repo.Create(req)
	if err != nil {
		return fmt.Errorf("failed to topup quota record: %w", err)
	}
	topupUserQuota(req.UserID, req.QuotaAmount)
	return nil
}

func (qs *quotaService) GiftQuota(repo repository.QuotaRepository, req *dto.Quota) error {
	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}
	if req.QuotaType != dto.QuotaTypeGift {
		return fmt.Errorf("quota type must be gift")
	}
	if !req.ExpireTime.IsZero() && req.ExpireTime.Before(time.Now()) {
		return fmt.Errorf("expire time must be greater than current time")
	}
	if req.QuotaAmount <= 0 {
		return fmt.Errorf("quota amount must be greater than 0")
	}
	req.QuotaID = utils.GenerateRandomUUID(12)
	req.Status = 1
	req.RemainingAmount = req.QuotaAmount
	req.FrozenAmount = 0

	err := repo.Create(req)
	if err != nil {
		return fmt.Errorf("failed to gift quota record: %w", err)
	}
	giftUserQuota(req.UserID, req.QuotaAmount)
	return nil
}

func (qs *quotaService) FreezeQuota(repo repository.QuotaRepository, quotaID string, amount float64) error {
	req, err := repo.GetByID(quotaID)
	if err != nil {
		return fmt.Errorf("failed to get quota record: %w", err)
	}
	req.Status = 2
	req.FrozenAmount = amount

	err = repo.Update(req)
	if err != nil {
		return fmt.Errorf("failed to freeze quota record: %w", err)
	}
	freezeUserQuota(req.UserID, amount)
	return nil
}

func (qs *quotaService) ExpireQuota(repo repository.QuotaRepository, quotaID string, expireTime time.Time) error {
	req, err := repo.GetByID(quotaID)
	if err != nil {
		return fmt.Errorf("failed to get quota record: %w", err)
	}
	req.ExpireTime = utils.MySQLTime(expireTime)

	err = repo.Update(req)
	if err != nil {
		return fmt.Errorf("failed to expire quota record: %w", err)
	}
	return nil
}

func (qs *quotaService) CostQuota(repo repository.QuotaRepository, userID string, amount float64) error {
	return nil
}

func costUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user quota: %w", err)
	}
	user.UserQuota.LeftQuota -= amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user quota: %w", err)
	}
	return nil
}

func topupUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	// 增加用户配额 总配额、剩余配额
	user.UserQuota.TotalQuota += amount
	user.UserQuota.LeftQuota += amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user quota: %w", err)
	}
	return nil
}

func giftUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	// 增加用户配额 总配额、赠送配额、剩余配额
	user.UserQuota.TotalQuota += amount
	user.UserQuota.GiftQuota += amount
	user.UserQuota.LeftQuota += amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user quota: %w", err)
	}
	return nil
}

func freezeUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.UserQuota.FrozenQuota += amount
	user.UserQuota.LeftQuota -= amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user quota: %w", err)
	}
	return nil
}

func unfreezeUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.UserQuota.FrozenQuota -= amount
	user.UserQuota.LeftQuota += amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to unfreeze user quota: %w", err)
	}
	return nil
}

func expireUserQuota(userID string, amount float64) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.UserQuota.ExpireQuota += amount
	user.UserQuota.LeftQuota -= amount
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to expire user quota: %w", err)
	}
	return nil
}

func LimitUserQuota(userID string, amount float64, expireTime time.Time) error {
	user, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.UserQuota.LimitQuota = amount
	user.UserQuota.LimitExpire = utils.MySQLTime(expireTime)
	err = updateUser(user)
	if err != nil {
		return fmt.Errorf("failed to set expire user quota: %w", err)
	}
	return nil
}

func QuotaExpireAsyncHandle(repo repository.QuotaRepository, quotaID string, userID string) error {
	req, err := repo.GetByID(quotaID)
	if err != nil {
		return fmt.Errorf("failed to get quota record: %w", err)
	}
	req.Status = 3
	err = repo.Update(req)
	if err != nil {
		return fmt.Errorf("failed to expire quota record: %w", err)
	}
	expireUserQuota(userID, req.RemainingAmount) // 过期用户剩余配额
	return nil
}
