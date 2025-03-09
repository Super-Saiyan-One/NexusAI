package repository

import (
	"fmt"
	"nexus-ai/common"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/utils"
	"time"

	"gorm.io/gorm"
)

// TaskVideoRepository 视频任务仓储接口
type TaskVideoRepository interface {
	Create(taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
	Update(taskVideo *dto.TaskVideo) (*dto.TaskVideo, error)
	Delete(taskVideoID string) error
	GetByID(taskVideoID string) (*dto.TaskVideo, error)
	List(page, pageSize int) ([]*dto.TaskVideo, int64, error)
	ListByUser(userID string, page, pageSize int) ([]*dto.TaskVideo, int64, error)
	Benchmark(count int) error
}

type taskVideoRepository struct {
	db *gorm.DB
}

// NewTaskVideoRepository 创建视频任务仓储实例
func NewTaskVideoRepository(db *gorm.DB) TaskVideoRepository {
	return &taskVideoRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *taskVideoRepository) convertToDTO(model *model.TaskVideo) *dto.TaskVideo {
	if model == nil {
		return nil
	}
	var requestParams dto.RequestParams
	var resultURLs dto.ResultURLs
	var errorLogs dto.ErrorLogs
	var nextRetryAt dto.NextRetryAt

	if err := model.RequestParams.ToStruct(&requestParams); err != nil {
		utils.SysError("解析请求参数失败:" + err.Error())
	}

	if err := model.ResultURLs.ToStruct(&resultURLs); err != nil {
		utils.SysError("解析结果URL失败:" + err.Error())
	}

	if err := model.ErrorLogs.ToStruct(&errorLogs); err != nil {
		utils.SysError("解析错误日志失败:" + err.Error())
	}

	if err := model.NextRetryAt.ToStruct(&nextRetryAt); err != nil {
		utils.SysError("解析下次重试时间失败:" + err.Error())
	}

	return &dto.TaskVideo{
		TaskVideoID:       model.TaskVideoID,
		RequestID:         model.RequestID,
		UserID:            model.UserID,
		TokenID:           model.TokenID,
		ModelID:           model.ModelID,
		ChannelID:         model.ChannelID,
		Status:            dto.TaskVideoStatus(model.Status),
		Provider:          dto.TaskVideoProvider(model.Provider),
		CallbackURL:       model.CallbackURL,
		RequestParams:     requestParams,
		ExternalRequestID: model.ExternalRequestID,
		RetryCount:        model.RetryCount,
		ErrorLogs:         errorLogs,
		NextRetryAt:       nextRetryAt,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		DeletedAt:         utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *taskVideoRepository) convertToModel(dto *dto.TaskVideo) (*model.TaskVideo, error) {
	if dto == nil {
		return nil, nil
	}
	requestParamsJSON, err := common.FromStruct(dto.RequestParams)
	if err != nil {
		return nil, fmt.Errorf("解析请求参数失败: %w", err)
	}

	resultURLsJSON, err := common.FromStruct(dto.ResultURLs)
	if err != nil {
		return nil, fmt.Errorf("解析结果URL失败: %w", err)
	}

	errorLogsJSON, err := common.FromStruct(dto.ErrorLogs)
	if err != nil {
		return nil, fmt.Errorf("解析错误日志失败: %w", err)
	}

	nextRetryAtJSON, err := common.FromStruct(dto.NextRetryAt)
	if err != nil {
		return nil, fmt.Errorf("解析下次重试时间失败: %w", err)
	}

	return &model.TaskVideo{
		TaskVideoID:       dto.TaskVideoID,
		RequestID:         dto.RequestID,
		UserID:            dto.UserID,
		TokenID:           dto.TokenID,
		ModelID:           dto.ModelID,
		ChannelID:         dto.ChannelID,
		Status:            string(dto.Status),
		Provider:          string(dto.Provider),
		CallbackURL:       dto.CallbackURL,
		RequestParams:     requestParamsJSON,
		ExternalRequestID: dto.ExternalRequestID,
		RetryCount:        dto.RetryCount,
		ResultURLs:        resultURLsJSON,
		ErrorLogs:         errorLogsJSON,
		NextRetryAt:       nextRetryAtJSON,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
		DeletedAt:         utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建视频任务记录
func (r *taskVideoRepository) Create(taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	model, err := r.convertToModel(taskVideo)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(model), nil
}

// Update 更新视频任务记录
func (r *taskVideoRepository) Update(taskVideo *dto.TaskVideo) (*dto.TaskVideo, error) {
	modelData, err := r.convertToModel(taskVideo)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.TaskVideo{}).Where("task_video_id = ?", taskVideo.TaskVideoID).Updates(modelData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(taskVideo.TaskVideoID)
}

// Delete 删除视频任务记录
func (r *taskVideoRepository) Delete(taskVideoID string) error {
	return r.db.Delete(&model.TaskVideo{}, "task_video_id = ?", taskVideoID).Error
}

// GetByID 根据ID获取视频任务记录
func (r *taskVideoRepository) GetByID(taskVideoID string) (*dto.TaskVideo, error) {
	var model model.TaskVideo
	if err := r.db.Where("task_video_id = ?", taskVideoID).First(&model).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&model), nil
}

// List 列出视频任务记录
func (r *taskVideoRepository) List(page, pageSize int) ([]*dto.TaskVideo, int64, error) {
	var total int64
	var models []model.TaskVideo

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.TaskVideo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.TaskVideo, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByUser 根据用户ID列出视频任务记录
func (r *taskVideoRepository) ListByUser(userID string, page, pageSize int) ([]*dto.TaskVideo, int64, error) {
	var total int64
	var models []model.TaskVideo

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.TaskVideo{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.TaskVideo, len(models))
	for i, m := range models {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *taskVideoRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行视频任务基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testTaskVideo := &dto.TaskVideo{
			TaskVideoID:       utils.GenerateRandomUUID(12),
			RequestID:         utils.GenerateRandomUUID(12),
			UserID:            utils.GenerateRandomUUID(12),
			TokenID:           utils.GenerateRandomUUID(12),
			ModelID:           utils.GenerateRandomUUID(12),
			ChannelID:         utils.GenerateRandomUUID(12),
			Status:            dto.TaskVideoStatusPending,
			Provider:          dto.TaskVideoProviderLuma,
			CallbackURL:       "http://example.com/callback",
			RequestParams:     dto.RequestParams{},
			ExternalRequestID: utils.GenerateRandomUUID(12),
			RetryCount:        0,
			ErrorLogs:         dto.ErrorLogs{},
			NextRetryAt:       dto.NextRetryAt{},
		}

		// 创建
		if _, err := r.Create(testTaskVideo); err != nil {
			utils.SysError("创建视频任务失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdTaskVideo, err := r.GetByID(testTaskVideo.TaskVideoID)
		if err != nil {
			utils.SysError("获取创建的视频任务失败: " + err.Error())
			return err
		}

		// 更新
		createdTaskVideo.Status = dto.TaskVideoStatusSuccess
		if _, err := r.Update(createdTaskVideo); err != nil {
			utils.SysError("更新视频任务失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdTaskVideo.TaskVideoID); err != nil {
			utils.SysError("删除视频任务失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
