package repository

import (
	"fmt"
	"math/rand"
	"nexus-ai/common"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// ChannelGroupRepository 渠道组仓储接口
type ChannelGroupRepository interface {
	Create(channelGroup *dto.ChannelGroup) error
	Update(channelGroup *dto.ChannelGroup) error
	Delete(channelGroupID string) error
	GetByID(channelGroupID string) (*dto.ChannelGroup, error)
	GetByName(name string) (*dto.ChannelGroup, error)
	List(page, pageSize int) ([]*dto.ChannelGroup, int64, error)
	GetListByDefaultLevel(defaultLevel int) ([]*dto.ChannelGroup, error)
	Benchmark(count int) error
}

type channelGroupRepository struct {
	db *gorm.DB
}

// NewChannelGroupRepository 创建渠道组仓储实例
func NewChannelGroupRepository(db *gorm.DB) ChannelGroupRepository {
	return &channelGroupRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *channelGroupRepository) convertToDTO(model *model.ChannelGroup) *dto.ChannelGroup {
	if model == nil {
		return nil
	}

	var priceFactor dto.ChannelGroupPriceFactor
	var options dto.ChannelGroupOptions
	var channels dto.ChannelGroupChannels

	if err := model.ChannelGroupPriceFactor.ToStruct(&priceFactor); err != nil {
		utils.SysError("解析价格系数失败:" + err.Error())
	}

	if err := model.ChannelGroupOptions.ToStruct(&options); err != nil {
		utils.SysError("解析配置选项失败:" + err.Error())
	}

	if err := model.ChannelGroupChannels.ToStruct(&channels); err != nil {
		utils.SysError("解析渠道组渠道失败:" + err.Error())
	}

	var deletedAt *utils.MySQLTime
	if model.DeletedAt.Valid {
		t := utils.MySQLTime(model.DeletedAt.Time)
		deletedAt = &t
	}

	return &dto.ChannelGroup{
		ChannelGroupID:          model.ChannelGroupID,
		ChannelGroupName:        model.ChannelGroupName,
		ChannelGroupDescription: model.ChannelGroupDescription,
		ChannelGroupPriceFactor: priceFactor,
		ChannelGroupOptions:     options,
		ChannelGroupChannels:    channels,
		CreatedAt:               model.CreatedAt,
		UpdatedAt:               model.UpdatedAt,
		DeletedAt:               deletedAt,
	}
}

// convertToModel 将DTO转换为数据模型
func (r *channelGroupRepository) convertToModel(dto *dto.ChannelGroup) (*model.ChannelGroup, error) {
	if dto == nil {
		return nil, nil
	}

	priceFactorJSON, err := common.FromStruct(dto.ChannelGroupPriceFactor)
	if err != nil {
		return nil, fmt.Errorf("转换价格系数失败: %w", err)
	}

	optionsJSON, err := common.FromStruct(dto.ChannelGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("转换配置选项失败: %w", err)
	}

	channelsJSON, err := common.FromStruct(dto.ChannelGroupChannels)
	if err != nil {
		return nil, fmt.Errorf("转换渠道组渠道失败: %w", err)
	}

	var deletedAt gorm.DeletedAt
	if dto.DeletedAt != nil {
		deletedAt.Time = time.Time(*dto.DeletedAt)
		deletedAt.Valid = true
	}

	return &model.ChannelGroup{
		ChannelGroupID:          dto.ChannelGroupID,
		ChannelGroupName:        dto.ChannelGroupName,
		ChannelGroupDescription: dto.ChannelGroupDescription,
		ChannelGroupPriceFactor: priceFactorJSON,
		ChannelGroupOptions:     optionsJSON,
		ChannelGroupChannels:    channelsJSON,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,
		DeletedAt:               deletedAt,
	}, nil
}

// Create 创建渠道组
func (r *channelGroupRepository) Create(channelGroup *dto.ChannelGroup) error {
	model, err := r.convertToModel(channelGroup)
	if err != nil {
		return err
	}
	return r.db.Create(model).Error
}

// Update 更新渠道组
func (r *channelGroupRepository) Update(channelGroup *dto.ChannelGroup) error {
	modelData, err := r.convertToModel(channelGroup)
	if err != nil {
		return err
	}
	return r.db.Model(&model.ChannelGroup{}).Where("channel_group_id = ?", channelGroup.ChannelGroupID).Updates(modelData).Error
}

// Delete 删除渠道组
func (r *channelGroupRepository) Delete(channelGroupID string) error {
	return r.db.Delete(&model.ChannelGroup{}, "channel_group_id = ?", channelGroupID).Error
}

// GetByID 根据ID获取渠道组
func (r *channelGroupRepository) GetByID(channelGroupID string) (*dto.ChannelGroup, error) {
	var channelGroup model.ChannelGroup
	if err := r.db.Where("channel_group_id = ?", channelGroupID).First(&channelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channelGroup), nil
}

// GetByName 根据名称获取渠道组
func (r *channelGroupRepository) GetByName(name string) (*dto.ChannelGroup, error) {
	var channelGroup model.ChannelGroup
	if err := r.db.Where("channel_group_name = ?", name).First(&channelGroup).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&channelGroup), nil
}

// List 获取渠道组列表
func (r *channelGroupRepository) List(page, pageSize int) ([]*dto.ChannelGroup, int64, error) {
	var total int64
	var channelGroups []model.ChannelGroup

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.ChannelGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&channelGroups).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.ChannelGroup, len(channelGroups))
	for i, cg := range channelGroups {
		dtoList[i] = r.convertToDTO(&cg)
	}

	return dtoList, total, nil
}

// GetListByDefaultLevel 根据默认等级获取渠道组列表
func (r *channelGroupRepository) GetListByDefaultLevel(defaultLevel int) ([]*dto.ChannelGroup, error) {
	var channelGroups []model.ChannelGroup
	if err := r.db.Where("CAST(JSON_EXTRACT(channel_group_options, '$.default_level') AS SIGNED) = ?", defaultLevel).Find(&channelGroups).Error; err != nil {
		return nil, err
	}
	dtoList := make([]*dto.ChannelGroup, len(channelGroups))
	for i, cg := range channelGroups {
		dtoList[i] = r.convertToDTO(&cg)
	}
	return dtoList, nil
}

// Benchmark 执行基准测试
func (r *channelGroupRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行渠道组基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testChannelGroup := &dto.ChannelGroup{
			ChannelGroupID:          utils.GenerateRandomUUID(12),
			ChannelGroupName:        fmt.Sprintf("benchmark_channel_group_%d", i),
			ChannelGroupDescription: fmt.Sprintf("这是第%d个基准测试渠道组", i),
			ChannelGroupPriceFactor: dto.ChannelGroupPriceFactor{
				RequestPriceFactor:    float64(rand.Intn(50)+50) / 100,
				ResponsePriceFactor:   float64(rand.Intn(50)+50) / 100,
				CompletionPriceFactor: float64(rand.Intn(50)+50) / 100,
				CachePriceFactor:      float64(rand.Intn(50)+50) / 100,
			},
			ChannelGroupOptions: dto.ChannelGroupOptions{
				MaxConcurrentRequests: rand.Intn(10) + 1,
				DefaultLevel:          rand.Intn(3) + 1,
				Discount:              float64(rand.Intn(50)+50) / 100,
				DiscountExpireAt:      utils.MySQLTime(time.Now().Add(time.Duration(rand.Intn(30)) * time.Hour)),
			},
			ChannelGroupChannels: dto.ChannelGroupChannels{
				Channels:  []string{},
				ModelsMap: map[string][]string{},
			},
		}

		// 创建
		if err := r.Create(testChannelGroup); err != nil {
			utils.SysError("创建渠道组失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdChannelGroup, err := r.GetByID(testChannelGroup.ChannelGroupID)
		if err != nil {
			utils.SysError("获取创建的渠道组失败: " + err.Error())
			return err
		}

		// 更新
		createdChannelGroup.ChannelGroupOptions.DefaultLevel = rand.Intn(5) + 1
		if err := r.Update(createdChannelGroup); err != nil {
			utils.SysError("更新渠道组失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdChannelGroup.ChannelGroupID); err != nil {
			utils.SysError("删除渠道组失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
