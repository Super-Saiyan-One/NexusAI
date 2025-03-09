package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	userDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *dto.User) (*dto.User, error)
	Update(user *dto.User) (*dto.User, error)
	Delete(userID string) error
	Search(userSearch *userDto.UserSearchRequest) ([]*dto.User, int64, error)
	UpdateLoginTime(userID string) error
	GetByID(userID string) (*dto.User, error)
	GetByEmail(email string) (*dto.User, error)
	GetByPhone(phone string) (*dto.User, error)
	GetByUsername(username string) (*dto.User, error)
	List(page, pageSize int) ([]*dto.User, int64, error)
	ListByQuota(quota dto.UserQuota, page, pageSize int) ([]*dto.User, int64, error)
	ListByOptions(options dto.UserOptions, page, pageSize int) ([]*dto.User, int64, error)
	Benchmark(count int) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *userRepository) convertToDTO(model *model.User) *dto.User {
	if model == nil {
		return nil
	}

	var oauthInfo dto.OAuthInfo
	var userQuota dto.UserQuota
	var userOptions dto.UserOptions

	if err := model.OAuthInfo.ToStruct(&oauthInfo); err != nil {
		utils.SysError("解析OAuth信息失败:" + err.Error())
	}

	if err := model.UserQuota.ToStruct(&userQuota); err != nil {
		utils.SysError("解析配额信息失败:" + err.Error())
	}

	if err := model.UserOptions.ToStruct(&userOptions); err != nil {
		utils.SysError("解析用户选项失败:" + err.Error())
	}

	return &dto.User{
		UserID:        model.UserID,
		UserGroupID:   model.UserGroupID,
		Username:      model.Username,
		Password:      model.Password,
		Email:         model.Email,
		Phone:         model.Phone,
		OAuthInfo:     oauthInfo,
		UserQuota:     userQuota,
		UserOptions:   userOptions,
		LastLoginTime: model.LastLoginTime,
		LastLoginIP:   model.LastLoginIP,
		Status:        model.Status,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		DeletedAt:     utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *userRepository) convertToModel(dto *dto.User) (*model.User, error) {
	if dto == nil {
		return nil, nil
	}

	oauthInfoJSON, err := common.FromStruct(dto.OAuthInfo)
	if err != nil {
		return nil, fmt.Errorf("转换OAuth信息失败: %w", err)
	}

	userQuotaJSON, err := common.FromStruct(dto.UserQuota)
	if err != nil {
		return nil, fmt.Errorf("转换配额信息失败: %w", err)
	}

	userOptionsJSON, err := common.FromStruct(dto.UserOptions)
	if err != nil {
		return nil, fmt.Errorf("转换用户选项失败: %w", err)
	}

	return &model.User{
		UserID:        dto.UserID,
		UserGroupID:   dto.UserGroupID,
		Username:      dto.Username,
		Password:      dto.Password,
		Email:         dto.Email,
		Phone:         dto.Phone,
		OAuthInfo:     oauthInfoJSON,
		UserQuota:     userQuotaJSON,
		UserOptions:   userOptionsJSON,
		LastLoginTime: dto.LastLoginTime,
		LastLoginIP:   dto.LastLoginIP,
		Status:        dto.Status,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
		DeletedAt:     utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建用户
func (r *userRepository) Create(user *dto.User) (*dto.User, error) {
	model, err := r.convertToModel(user)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新用户
func (r *userRepository) Update(user *dto.User) (*dto.User, error) {
	modelData, err := r.convertToModel(user)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.User{}).Where("user_id = ?", user.UserID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(user.UserID)
}

// Delete 删除用户
func (r *userRepository) Delete(userID string) error {
	return r.db.Delete(&model.User{}, "user_id = ?", userID).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(userID string) (*dto.User, error) {
	var user model.User
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&user), nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*dto.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&user), nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(phone string) (*dto.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&user), nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*dto.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&user), nil
}

// List 获取用户列表
func (r *userRepository) List(page, pageSize int) ([]*dto.User, int64, error) {
	var total int64
	var users []model.User

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.User, len(users))
	for i, u := range users {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByQuota 根据配额筛选用户
func (r *userRepository) ListByQuota(quota dto.UserQuota, page, pageSize int) ([]*dto.User, int64, error) {
	var total int64
	var users []model.User

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.User{})

	quotaJSON, err := common.FromStruct(quota)
	if err != nil {
		return nil, 0, fmt.Errorf("转换配额信息失败: %w", err)
	}

	query = query.Where("user_quota @> ?", quotaJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.User, len(users))
	for i, u := range users {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// ListByOptions 根据用户选项筛选用户
func (r *userRepository) ListByOptions(options dto.UserOptions, page, pageSize int) ([]*dto.User, int64, error) {
	var total int64
	var users []model.User

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.User{})

	optionsJSON, err := common.FromStruct(options)
	if err != nil {
		return nil, 0, fmt.Errorf("转换用户选项失败: %w", err)
	}

	query = query.Where("user_options @> ?", optionsJSON)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.User, len(users))
	for i, u := range users {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *userRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行用户基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testUser := &dto.User{
			UserID:   utils.GenerateRandomUUID(12),
			Username: fmt.Sprintf("benchmark_user_%s", utils.GenerateRandomUUID(12)),
			Email:    fmt.Sprintf("benchmark_%d%s@example.com", i, utils.GenerateRandomString(8)),
			Phone:    fmt.Sprintf("1%010d", rand.Intn(10000000000)),
			Password: utils.HashPassword(fmt.Sprintf("test%d", i)),
			UserQuota: dto.UserQuota{
				TotalQuota:  float64(rand.Intn(1000)),
				FrozenQuota: float64(rand.Intn(100)),
				GiftQuota:   float64(rand.Intn(100)),
			},
			UserOptions: dto.UserOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				DefaultLevel:          rand.Intn(3) + 1,
				APIDiscount:           float64(rand.Intn(50)+50) / 100,
			},
			Status: 1,
		}

		// 创建
		if _, err := r.Create(testUser); err != nil {
			utils.SysError("创建用户失败: " + err.Error())
			return err
		}

		// 在创建后获取ID
		createdUser, err := r.GetByID(testUser.UserID)
		if err != nil {
			utils.SysError("获取创建的用户失败: " + err.Error())
			return err
		}

		// 更新
		createdUser.UserOptions.DefaultLevel = rand.Intn(5) + 1
		if _, err := r.Update(createdUser); err != nil {
			utils.SysError("更新用户失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdUser.UserID); err != nil {
			utils.SysError("删除用户失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}

// Search 根据搜索条件筛选用户
func (r *userRepository) Search(req *userDto.UserSearchRequest) ([]*dto.User, int64, error) {
	var total int64
	var users []model.User

	query := r.db.Model(&model.User{})

	// 添加筛选条件
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.UserGroupID != "" {
		query = query.Where("user_group_id = ?", req.UserGroupID)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if len(req.Levels) > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_options, '$.default_level') AS SIGNED) = ANY(?)", req.Levels)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 修改 JSON 字段查询的语法
	if req.MinConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_options, '$.max_concurrent_requests') AS SIGNED) >= ?", req.MinConcurrentRequests)
	}
	if req.MaxConcurrentRequests > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_options, '$.max_concurrent_requests') AS SIGNED) <= ?", req.MaxConcurrentRequests)
	}

	// API折扣范围查询
	if req.MinAPIDiscount > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_options, '$.api_discount') AS DECIMAL(10,2)) >= ?", req.MinAPIDiscount)
	}
	if req.MaxAPIDiscount > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_options, '$.api_discount') AS DECIMAL(10,2)) <= ?", req.MaxAPIDiscount)
	}

	// 配额范围查询
	if req.MinQuota > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_quota, '$.total_quota') AS DECIMAL(10,2)) >= ?", req.MinQuota)
	}
	if req.MaxQuota > 0 {
		query = query.Where("CAST(JSON_EXTRACT(user_quota, '$.total_quota') AS DECIMAL(10,2)) <= ?", req.MaxQuota)
	}

	// 时间范围查询
	if !req.EarlyLastLoginTime.IsZero() {
		query = query.Where("last_login_time >= ?", req.EarlyLastLoginTime)
	}
	if !req.LateLastLoginTime.IsZero() {
		query = query.Where("last_login_time <= ?", req.LateLastLoginTime)
	}
	if !req.EarlyCreatedTime.IsZero() {
		query = query.Where("created_at >= ?", req.EarlyCreatedTime)
	}
	if !req.LateCreatedTime.IsZero() {
		query = query.Where("created_at <= ?", req.LateCreatedTime)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtoList := make([]*dto.User, len(users))
	for i, u := range users {
		dtoList[i] = r.convertToDTO(&u)
	}

	return dtoList, total, nil
}

// UpdateLoginTime 更新用户最后一次登录时间
func (r *userRepository) UpdateLoginTime(userID string) error {
	return r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("last_login_time", utils.MySQLTime(time.Now())).Error
}
