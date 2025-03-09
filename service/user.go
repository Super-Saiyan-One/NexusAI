package service

import (
	"errors"
	"fmt"
	userDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/utils"
)

type UserService interface {
	UserRegister(repo repository.UserRepository, user *dto.User) (*dto.User, error)
	UserLogin(repo repository.UserRepository, user *userDto.LoginRequest) (*dto.User, []string, error)
	UserSearch(repo repository.UserRepository, userSearch *userDto.UserSearchRequest) ([]*dto.User, error)
	UserDelete(repo repository.UserRepository, userID string) error
	UserLogout(repo repository.UserRepository, userID string, accessToken, refreshToken string) error
	UserUpdate(repo repository.UserRepository, user *dto.User) (*dto.User, error)
	UserPassword(repo repository.UserRepository, userID string, oldPassword, newPassword string) error
	UserAvailableModels(repo repository.UserRepository, userID string) ([]*dto.Model, error)
	UserAvailableChannels(repo repository.UserRepository, userID string) ([]*dto.Channel, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (us *userService) UserRegister(repo repository.UserRepository, user *dto.User) (*dto.User, error) {
	existingUser, _ := repo.GetByEmail(user.Email)
	if existingUser != nil && existingUser.DeletedAt == nil { // 用户已存在且未被删除
		return existingUser, errors.New("user already exists")
	}
	return repo.Create(user)
}

func (us *userService) UserLogin(repo repository.UserRepository, user *userDto.LoginRequest) (*dto.User, []string, error) {
	targetUser, err := repo.GetByEmail(user.Email)
	if err != nil { // 用户不存在
		return nil, nil, err
	}

	if targetUser.Status == 0 { // 用户被禁用
		return nil, nil, errors.New("user is disabled")
	}

	if targetUser.DeletedAt != nil { // 用户被删除
		return nil, nil, errors.New("user is deleted")
	}

	if !utils.CheckPasswordHash(user.Password, targetUser.Password) { // 密码错误
		return nil, nil, errors.New("password is incorrect")
	}

	accessToken, err := GenerateToken(targetUser.UserID, AccessToken) // 生成access token
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := GenerateToken(targetUser.UserID, RefreshToken) // 生成refresh token
	if err != nil {
		return nil, nil, err
	}
	repo.UpdateLoginTime(targetUser.UserID)                     // 更新用户最后一次登录时间
	return targetUser, []string{accessToken, refreshToken}, nil // 返回用户信息和token
}

func (us *userService) UserSearch(repo repository.UserRepository, userSearch *userDto.UserSearchRequest) ([]*dto.User, error) {
	users, _, err := repo.Search(userSearch)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) UserDelete(repo repository.UserRepository, userID string) error {
	existingUser, _ := repo.GetByID(userID)
	if existingUser == nil { // 用户不存在
		return errors.New("user not found")
	}
	if existingUser.DeletedAt != nil { // 用户已删除
		return errors.New("user already deleted")
	}
	err := repo.Delete(userID)
	if err != nil { // 删除用户失败
		return err
	}
	InvalidateAllUserTokens(userID) // 使指定用户的所有token失效
	return nil
}

func (us *userService) UserLogout(repo repository.UserRepository, userID string, accessToken, refreshToken string) error {
	existingUser, _ := repo.GetByID(userID)
	if existingUser == nil { // 用户不存在
		return errors.New("user not found")
	}
	if existingUser.DeletedAt != nil { // 用户已删除
		return errors.New("user already deleted")
	}
	if existingUser.Status == 0 { // 用户被禁用
		return errors.New("user is disabled")
	}
	InvalidateToken(accessToken)
	InvalidateToken(refreshToken)
	return nil
}

func (us *userService) UserUpdate(repo repository.UserRepository, user *dto.User) (*dto.User, error) {
	existingUser, _ := repo.GetByID(user.UserID)
	if existingUser == nil { // 用户不存在
		return nil, errors.New("user not found")
	}
	if existingUser.DeletedAt != nil { // 用户已删除
		return nil, errors.New("user already deleted")
	}
	return repo.Update(user)
}

func (us *userService) UserPassword(repo repository.UserRepository, userID string, oldPassword, newPassword string) error {
	user, err := repo.GetByID(userID)
	if err != nil { // 用户不存在
		return err
	}

	if user.DeletedAt != nil { // 用户已删除
		return errors.New("user is deleted")
	}

	if user.Status == 0 { // 用户被禁用
		return errors.New("user is disabled")
	}

	if !utils.CheckPasswordHash(oldPassword, user.Password) { // 旧密码错误
		return errors.New("old password is incorrect")
	}

	user.Password = utils.HashPassword(newPassword)
	_, err = repo.Update(user)
	if err != nil { // 更新用户密码失败
		return err
	}

	return nil
}

// UserAvailableModels 获取用户可用模型
func (us *userService) UserAvailableModels(repo repository.UserRepository, userID string) ([]*dto.Model, error) {
	// TODO: 获取用户可用模型
	return []*dto.Model{}, nil
}

// UserAvailableChannels 获取用户可用渠道
func (us *userService) UserAvailableChannels(repo repository.UserRepository, userID string) ([]*dto.Channel, error) {
	// TODO: 获取用户可用渠道
	return []*dto.Channel{}, nil
}

func getUser(userID string) (*dto.User, error) {
	userRepo := repository.NewUserRepository(model.GetDB())
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func updateUser(user *dto.User) error {
	userRepo := repository.NewUserRepository(model.GetDB())
	_, err := userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
