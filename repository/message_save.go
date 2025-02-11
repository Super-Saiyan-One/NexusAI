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

// MessageSaveRepository 消息存储仓储接口
type MessageSaveRepository interface {
	Create(message *dto.MessageSave) error
	Update(message *dto.MessageSave) error
	Delete(messageID string) error
	GetByID(messageID string) (*dto.MessageSave, error)
	GetByRequestID(requestID string) (*dto.MessageSave, error)
	List(page, pageSize int) ([]*dto.MessageSave, int64, error)
	ListByUser(userID string, page, pageSize int) ([]*dto.MessageSave, int64, error)
	ListByToken(tokenID string, page, pageSize int) ([]*dto.MessageSave, int64, error)
	ListByModel(modelID string, page, pageSize int) ([]*dto.MessageSave, int64, error)
	ListByChannel(channelID string, page, pageSize int) ([]*dto.MessageSave, int64, error)
	ListByStatus(status int8, page, pageSize int) ([]*dto.MessageSave, int64, error)
	Benchmark(count int) error
}

type messageSaveRepository struct {
	db *gorm.DB
}

// NewMessageSaveRepository 创建消息存储仓储实例
func NewMessageSaveRepository(db *gorm.DB) MessageSaveRepository {
	return &messageSaveRepository{db: db}
}

// convertToDTO 将数据模型转换为DTO
func (r *messageSaveRepository) convertToDTO(model *model.MessageSave) *dto.MessageSave {
	if model == nil {
		return nil
	}

	var messageOptions dto.MessageOptions
	var streamOptions dto.StreamOptions
	var messageExtra dto.MessageExtra

	if err := model.MessageOptions.ToStruct(&messageOptions); err != nil {
		utils.SysError("解析消息配置失败:" + err.Error())
	}

	if err := model.StreamOptions.ToStruct(&streamOptions); err != nil {
		utils.SysError("解析流式输出配置失败:" + err.Error())
	}

	if err := model.MessageExtra.ToStruct(&messageExtra); err != nil {
		utils.SysError("解析消息额外信息失败:" + err.Error())
	}

	return &dto.MessageSave{
		MessageID:        model.MessageID,
		UserID:           model.UserID,
		TokenID:          model.TokenID,
		ModelID:          model.ModelID,
		ChannelID:        model.ChannelID,
		RequestID:        model.RequestID,
		RequestTitle:     model.RequestTitle,
		RequestContent:   model.RequestContent,
		ResponseContent:  model.ResponseContent,
		Input:            model.Input,
		Instruction:      model.Instruction,
		Prompt:           model.Prompt,
		Stream:           model.Stream,
		StreamOptions:    streamOptions,
		RequestTokens:    model.RequestTokens,
		PromptTokens:     model.PromptTokens,
		CompletionTokens: model.CompletionTokens,
		Latency:          model.Latency,
		MessageStatus:    model.MessageStatus,
		RetryCount:       model.RetryCount,
		ErrorType:        model.ErrorType,
		ErrorInfo:        model.ErrorInfo,
		MessageOptions:   messageOptions,
		MessageExtra:     messageExtra,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		DeletedAt:        utils.FromDeletedAt(model.DeletedAt),
	}
}

// convertToModel 将DTO转换为数据模型
func (r *messageSaveRepository) convertToModel(dto *dto.MessageSave) (*model.MessageSave, error) {
	if dto == nil {
		return nil, nil
	}

	messageOptionsJSON, err := common.FromStruct(dto.MessageOptions)
	if err != nil {
		return nil, fmt.Errorf("转换消息配置失败: %w", err)
	}

	streamOptionsJSON, err := common.FromStruct(dto.StreamOptions)
	if err != nil {
		return nil, fmt.Errorf("转换流式输出配置失败: %w", err)
	}

	messageExtraJSON, err := common.FromStruct(dto.MessageExtra)
	if err != nil {
		return nil, fmt.Errorf("转换消息额外信息失败: %w", err)
	}

	return &model.MessageSave{
		MessageID:        dto.MessageID,
		UserID:           dto.UserID,
		TokenID:          dto.TokenID,
		ModelID:          dto.ModelID,
		ChannelID:        dto.ChannelID,
		RequestID:        dto.RequestID,
		RequestTitle:     dto.RequestTitle,
		RequestContent:   dto.RequestContent,
		ResponseContent:  dto.ResponseContent,
		Input:            dto.Input,
		Instruction:      dto.Instruction,
		Prompt:           dto.Prompt,
		Stream:           dto.Stream,
		RequestTokens:    dto.RequestTokens,
		PromptTokens:     dto.PromptTokens,
		CompletionTokens: dto.CompletionTokens,
		Latency:          dto.Latency,
		MessageStatus:    dto.MessageStatus,
		RetryCount:       dto.RetryCount,
		ErrorType:        dto.ErrorType,
		ErrorInfo:        dto.ErrorInfo,
		MessageOptions:   messageOptionsJSON,
		StreamOptions:    streamOptionsJSON,
		MessageExtra:     messageExtraJSON,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
		DeletedAt:        utils.ToDeletedAt(dto.DeletedAt),
	}, nil
}

// Create 创建消息记录
func (r *messageSaveRepository) Create(message *dto.MessageSave) error {
	model, err := r.convertToModel(message)
	if err != nil {
		return err
	}
	return r.db.Create(model).Error
}

// Update 更新消息记录
func (r *messageSaveRepository) Update(message *dto.MessageSave) error {
	modelData, err := r.convertToModel(message)
	if err != nil {
		return err
	}
	return r.db.Model(&model.MessageSave{}).Where("message_id = ?", message.MessageID).Updates(modelData).Error
}

// Delete 删除消息记录
func (r *messageSaveRepository) Delete(messageID string) error {
	return r.db.Delete(&model.MessageSave{}, "message_id = ?", messageID).Error
}

// GetByID 根据ID获取消息记录
func (r *messageSaveRepository) GetByID(messageID string) (*dto.MessageSave, error) {
	var message model.MessageSave
	if err := r.db.Where("message_id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&message), nil
}

// GetByRequestID 根据请求ID获取消息记录
func (r *messageSaveRepository) GetByRequestID(requestID string) (*dto.MessageSave, error) {
	var message model.MessageSave
	if err := r.db.Where("request_id = ?", requestID).First(&message).Error; err != nil {
		return nil, err
	}
	return r.convertToDTO(&message), nil
}

// List 获取消息记录列表
func (r *messageSaveRepository) List(page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.MessageSave{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByUser 获取用户的消息记录列表
func (r *messageSaveRepository) ListByUser(userID string, page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.MessageSave{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByToken 获取令牌的消息记录列表
func (r *messageSaveRepository) ListByToken(tokenID string, page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.MessageSave{}).Where("token_id = ?", tokenID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByModel 获取模型的消息记录列表
func (r *messageSaveRepository) ListByModel(modelID string, page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.MessageSave{}).Where("model_id = ?", modelID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByChannel 获取渠道的消息记录列表
func (r *messageSaveRepository) ListByChannel(channelID string, page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.MessageSave{}).Where("channel_id = ?", channelID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// ListByStatus 获取指定状态的消息记录列表
func (r *messageSaveRepository) ListByStatus(status int8, page, pageSize int) ([]*dto.MessageSave, int64, error) {
	var total int64
	var messages []model.MessageSave

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.MessageSave{}).Where("message_status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	dtoList := make([]*dto.MessageSave, len(messages))
	for i, m := range messages {
		dtoList[i] = r.convertToDTO(&m)
	}

	return dtoList, total, nil
}

// Benchmark 执行基准测试
func (r *messageSaveRepository) Benchmark(count int) error {
	utils.SysInfo("开始执行消息记录基准测试...")
	startTime := time.Now()

	for i := 0; i < count; i++ {
		testMessage := &dto.MessageSave{
			MessageID:        utils.GenerateRandomUUID(12),
			UserID:           utils.GenerateRandomUUID(12),
			TokenID:          utils.GenerateRandomUUID(12),
			ModelID:          utils.GenerateRandomUUID(12),
			ChannelID:        utils.GenerateRandomUUID(12),
			RequestID:        fmt.Sprintf("REQ_%s", utils.GenerateRandomString(16)),
			RequestTitle:     fmt.Sprintf("测试消息 %d", i),
			RequestContent:   fmt.Sprintf("这是一条测试消息内容 %d", i),
			ResponseContent:  fmt.Sprintf("这是一条测试消息响应内容 %d", i),
			Input:            fmt.Sprintf("这是一条测试消息输入 %d", i),
			Instruction:      fmt.Sprintf("这是一条测试消息指令 %d", i),
			Prompt:           fmt.Sprintf("这是一条测试消息提示词 %d", i),
			Stream:           int8(rand.Intn(2)),
			RequestTokens:    rand.Intn(1000),
			PromptTokens:     rand.Intn(500),
			CompletionTokens: rand.Intn(500),
			Latency:          rand.Intn(2000),

			MessageOptions: dto.MessageOptions{
				MaxTokens:           rand.Intn(2000) + 1000,
				MaxCompletionTokens: rand.Intn(1000) + 500,
				Temperature:         rand.Intn(100),
				TopP:                rand.Intn(100),
				TopK:                rand.Intn(100),
				Stop:                []string{"test"},
				ResponseFormat:      "json",
				EncodingFormat:      "utf-8",
				N:                   1,
				Size:                1,
				Seed:                1,
				FunctionCall:        1,
				PresencePenalty:     1,
				FrequencyPenalty:    1,
				BestOf:              1,
				Dimensions:          1,
			},
			MessageStatus: int8(rand.Intn(3) + 1),
			RetryCount:    rand.Intn(3),
			ErrorType:     "test",
			ErrorInfo:     "test",
		}

		// 创建
		if err := r.Create(testMessage); err != nil {
			utils.SysError("创建消息记录失败: " + err.Error())
			return err
		}

		// 获取创建后的记录
		createdMessage, err := r.GetByID(testMessage.MessageID)
		if err != nil {
			utils.SysError("获取创建的消息记录失败: " + err.Error())
			return err
		}

		// 更新
		createdMessage.MessageStatus = int8(rand.Intn(3) + 1)
		if err := r.Update(createdMessage); err != nil {
			utils.SysError("更新消息记录失败: " + err.Error())
			return err
		}

		// 删除
		if err := r.Delete(createdMessage.MessageID); err != nil {
			utils.SysError("删除消息记录失败: " + err.Error())
			return err
		}
	}

	duration := time.Since(startTime)
	utils.SysInfo("基准测试完成，总耗时: " + duration.String() + ", 平均每组操作耗时: " + (duration / time.Duration(count)).String())
	return nil
}
