package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统中所有用户的基本信息、认证信息、配额信息等
type User struct {
	UserID        string          `gorm:"column:user_id;type:char(36);primaryKey;default:(UUID())" json:"user_id"` // 用户唯一标识
	UserGroupID   string          `gorm:"column:user_group_id;type:char(36);not null" json:"user_group_id"`        // 用户组ID
	Username      string          `gorm:"column:username;size:50;uniqueIndex;not null" json:"username"`            // 用户名
	Password      string          `gorm:"column:password;size:255;not null" json:"password"`                       // 密码哈希
	Email         string          `gorm:"column:email;size:100;uniqueIndex;not null" json:"email"`                 // 邮箱地址
	Phone         string          `gorm:"column:phone;size:20;uniqueIndex" json:"phone"`                           // 手机号码
	OAuthInfo     common.JSON     `gorm:"column:oauth_info;type:json" json:"oauth_info"`                           // 第三方登录相关的认证信息(如GitHub、微信、Google等)
	UserQuota     common.JSON     `gorm:"column:user_quota;type:json" json:"user_quota"`                           // 各种资源使用配额信息
	UserOptions   common.JSON     `gorm:"column:user_options;type:json" json:"user_options"`                       // 个性化设置和偏好选项
	LastLoginTime utils.MySQLTime `gorm:"column:last_login_time" json:"last_login_time"`                           // 最后一次成功登录的时间
	LastLoginIP   string          `gorm:"column:last_login_ip;size:100" json:"last_login_ip"`                      // 最后一次登录的IP地址
	Status        int8            `gorm:"column:status;index;not null;default:1" json:"status"`                    // 用户状态，1:正常 0:禁用

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null;" json:"created_at"` // 用户账号的创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null;" json:"updated_at"` // 用户信息的最后更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`           // 软删除时间
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 在创建记录前自动设置时间
func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.CreatedAt = utils.MySQLTime(utils.GetTime())
	user.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (user *User) BeforeUpdate(tx *gorm.DB) error {
	user.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
