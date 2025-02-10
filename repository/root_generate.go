package repository

import (
	dto "nexus-ai/dto/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

func RootGenerate(db *gorm.DB) {

	utils.SysInfo("RootGenerate start")
	userGroupRepo := NewUserGroupRepository(db)
	userRepo := NewUserRepository(db)

	// 检查administrator用户组是否存在
	adminGroup, err := userGroupRepo.GetByName("administrator")
	if err == gorm.ErrRecordNotFound {
		// 创建administrator用户组
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

		if err := userGroupRepo.Create(adminGroup); err != nil {
			utils.SysError("创建administrator用户组失败: " + err.Error())
			return
		}
	} else if err != nil {
		utils.SysError("查询administrator用户组失败: " + err.Error())
		return
	}

	// 检查root用户是否存在
	_, err = userRepo.GetByUsername("root")
	if err == gorm.ErrRecordNotFound {
		// 创建root用户
		rootUser := &dto.User{
			UserID:      utils.GenerateRandomUUID(12),
			UserGroupID: adminGroup.UserGroupID,
			Username:    "root",
			Password:    utils.HashPassword("root@2024"),
			Email:       "root@system.local",
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

		if _, err := userRepo.Create(rootUser); err != nil {
			utils.SysError("创建root用户失败: " + err.Error())
			return
		}
	} else if err != nil {
		utils.SysError("查询root用户失败: " + err.Error())
		return
	}

}
