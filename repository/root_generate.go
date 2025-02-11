package repository

import (
	"context"
	"fmt"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/redis"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

func RootGenerate(db *gorm.DB) {
	utils.SysInfo("RootGenerate start")
	userGroupRepo := NewUserGroupRepository(db)
	userRepo := NewUserRepository(db)
	modelGroupRepo := NewModelGroupRepository(db)

	// 创建管理员用户组
	adminGroup, err := createAdministratorGroup(userGroupRepo)
	if err != nil {
		utils.SysError("创建administrator用户组失败: " + err.Error())
		return
	}

	// 创建普通用户组
	if err := createOrdinaryGroup(userGroupRepo); err != nil {
		utils.SysError("创建ordinary用户组失败: " + err.Error())
		return
	}

	// 创建root用户
	if err := createRootUser(userRepo, adminGroup.UserGroupID); err != nil {
		utils.SysError("创建root用户失败: " + err.Error())
		return
	}

	// 创建模型组
	if err := createModelGroups(modelGroupRepo); err != nil {
		utils.SysError("创建模型组失败: " + err.Error())
		return
	}
}

// createAdministratorGroup 创建管理员用户组
func createAdministratorGroup(userGroupRepo UserGroupRepository) (*dto.UserGroup, error) {
	adminGroup, err := userGroupRepo.GetByName("administrator")
	if err == nil {
		return adminGroup, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	adminGroup = &dto.UserGroup{
		UserGroupID:          utils.GenerateRandomUUID(12),
		UserGroupName:        "administrator",
		UserGroupDescription: "系统管理员组",
		UserGroupPriceFactor: dto.UserGroupPriceFactor{
			RequestPriceFactor:    0,
			ResponsePriceFactor:   0,
			CompletionPriceFactor: 0,
			CachePriceFactor:      0,
		},
		UserGroupOptions: dto.UserGroupOptions{
			MaxConcurrentRequests: 10000,
			DefaultLevel:          99,
			ExtraAllowedModels:    []string{"*"},
			ExtraAllowedChannels:  []string{"*"},
			APIDiscount:           0,
		},
		CreatedAt: utils.MySQLTime(time.Now()),
		UpdatedAt: utils.MySQLTime(time.Now()),
	}

	return userGroupRepo.Create(adminGroup)
}

// createOrdinaryGroup 创建普通用户组
func createOrdinaryGroup(userGroupRepo UserGroupRepository) error {
	ordinaryGroup, err := userGroupRepo.GetByName("ordinary")
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		ordinaryGroup = &dto.UserGroup{
			UserGroupID:          utils.GenerateRandomUUID(12),
			UserGroupName:        "ordinary",
			UserGroupDescription: "普通用户组",
			UserGroupPriceFactor: dto.UserGroupPriceFactor{
				RequestPriceFactor:    constant.DefaultUserGroupRequestPriceFactor,
				ResponsePriceFactor:   constant.DefaultUserGroupResponsePriceFactor,
				CompletionPriceFactor: constant.DefaultUserGroupCompletionPriceFactor,
				CachePriceFactor:      constant.DefaultUserGroupCachePriceFactor,
			},
			UserGroupOptions: dto.UserGroupOptions{
				MaxConcurrentRequests: constant.DefaultUserGroupMaxConcurrentRequests,
				DefaultLevel:          constant.DefaultUserGroupDefaultLevel,
				ExtraAllowedModels:    []string{},
				ExtraAllowedChannels:  []string{},
				APIDiscount:           constant.DefaultUserGroupAPIDiscount,
			},
			CreatedAt: utils.MySQLTime(time.Now()),
			UpdatedAt: utils.MySQLTime(time.Now()),
		}

		createdGroup, err := userGroupRepo.Create(ordinaryGroup)
		if err != nil {
			return err
		}
		ordinaryGroup = createdGroup
	}

	// 更新Redis中的ordinary用户组ID
	ctx := context.Background()
	if err := redis.Set(ctx, "ordinary_user_group_id", ordinaryGroup.UserGroupID, 0); err != nil {
		utils.SysError("更新ordinary用户组ID到Redis失败: " + err.Error())
	} else {
		utils.SysInfo("ordinary用户组ID已更新到Redis: " + ordinaryGroup.UserGroupID)
	}

	return nil
}

// createRootUser 创建root用户
func createRootUser(userRepo UserRepository, adminGroupID string) error {
	_, err := userRepo.GetByUsername("root")
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	rootUser := &dto.User{
		UserID:      utils.GenerateRandomUUID(12),
		UserGroupID: adminGroupID,
		Username:    constant.RootUserName,
		Password:    utils.HashPassword(constant.RootUserPassword),
		Email:       constant.RootUserEmail,
		Phone:       "",
		OAuthInfo:   dto.OAuthInfo{},
		UserQuota: dto.UserQuota{
			TotalQuota:  1000000,
			FrozenQuota: 0,
			GiftQuota:   1000000,
		},
		UserOptions: dto.UserOptions{
			MaxConcurrentRequests: 10000,
			DefaultLevel:          99,
			APIDiscount:           0,
		},
		Status:    1,
		CreatedAt: utils.MySQLTime(time.Now()),
		UpdatedAt: utils.MySQLTime(time.Now()),
	}

	_, err = userRepo.Create(rootUser)
	return err
}

// createModelGroups 创建模型组
func createModelGroups(modelGroupRepo ModelGroupRepository) error {
	// 创建 administrator 等级模型组
	_, err := modelGroupRepo.GetByName("administrator")
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		adminModelGroup := &dto.ModelGroup{
			ModelGroupID:          utils.GenerateRandomUUID(12),
			ModelGroupName:        "administrator",
			ModelGroupDescription: "管理员模型组",
			ModelGroupPriceFactor: dto.ModelGroupPriceFactor{
				RequestPriceFactor:    0,
				ResponsePriceFactor:   0,
				CompletionPriceFactor: 0,
				CachePriceFactor:      0,
			},
			ModelGroupOptions: dto.ModelGroupOptions{
				MaxConcurrentRequests: 10000,
				DefaultLevel:          99,
				APIDiscount:           0,
				APIDiscountExpireAt:   utils.MySQLTime(time.Now().Add(36500 * 24 * time.Hour)), // 100年后过期
			},
			CreatedAt: utils.MySQLTime(time.Now()),
			UpdatedAt: utils.MySQLTime(time.Now()),
		}

		if _, err := modelGroupRepo.Create(adminModelGroup); err != nil {
			return fmt.Errorf("创建administrator模型组失败: %w", err)
		}
	}

	// 创建 ordinary 等级模型组
	_, err = modelGroupRepo.GetByName("ordinary")
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		ordinaryModelGroup := &dto.ModelGroup{
			ModelGroupID:          utils.GenerateRandomUUID(12),
			ModelGroupName:        "ordinary",
			ModelGroupDescription: "普通用户模型组",
			ModelGroupPriceFactor: dto.ModelGroupPriceFactor{
				RequestPriceFactor:    constant.DefaultModelGroupRequestPriceFactor,
				ResponsePriceFactor:   constant.DefaultModelGroupResponsePriceFactor,
				CompletionPriceFactor: constant.DefaultModelGroupCompletionPriceFactor,
				CachePriceFactor:      constant.DefaultModelGroupCachePriceFactor,
			},
			ModelGroupOptions: dto.ModelGroupOptions{
				MaxConcurrentRequests: constant.DefaultModelGroupMaxConcurrentRequests,
				DefaultLevel:          constant.DefaultModelGroupDefaultLevel,
				APIDiscount:           constant.DefaultModelGroupAPIDiscount,
				APIDiscountExpireAt:   utils.MySQLTime(time.Now().Add(36500 * 24 * time.Hour)), // 100年后过期
			},
			CreatedAt: utils.MySQLTime(time.Now()),
			UpdatedAt: utils.MySQLTime(time.Now()),
		}

		if _, err := modelGroupRepo.Create(ordinaryModelGroup); err != nil {
			return fmt.Errorf("创建ordinary模型组失败: %w", err)
		}
	}

	return nil
}
