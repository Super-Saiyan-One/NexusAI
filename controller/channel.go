package controller

import (
	"net/http"
	"nexus-ai/constant"
	channelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/redis"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type ChannelController interface {
	GetChannelRepo() repository.ChannelRepository
	ChannelCreate(c *gin.Context)
	ChannelUpdate(c *gin.Context)
	ChannelSearch(c *gin.Context)
	ChannelDelete(c *gin.Context)
}

type channelController struct {
	service service.ChannelService
}

func NewChannelController() ChannelController {
	channelService := service.NewChannelService()
	return &channelController{service: channelService}
}

func (cc *channelController) GetChannelRepo() repository.ChannelRepository {
	return repository.NewChannelRepository(model.GetDB())
}

// ChannelCreate 创建渠道
func (cc *channelController) ChannelCreate(c *gin.Context) {
	var channel dto.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelPrefix+"_create")
		return
	}

	channel.ChannelID = utils.GenerateRandomUUID(12)
	if channel.ChannelName == "" {
		channel.ChannelName = utils.GenerateRandomString(12)
	}
	if channel.ChannelGroupID == "" {
		channel.ChannelGroupID, _ = redis.Get(c, "ordinary_channel_group_id")
	}

	// 设置默认值
	if channel.UpstreamOptions.Timeout <= 0 {
		channel.UpstreamOptions.Timeout = constant.DefaultChannelUpstreamTimeout
	}
	if channel.UpstreamOptions.MaxRetries <= 0 {
		channel.UpstreamOptions.MaxRetries = constant.DefaultChannelUpstreamMaxRetries
	}
	if channel.UpstreamOptions.DialTimeout <= 0 {
		channel.UpstreamOptions.DialTimeout = constant.DefaultChannelUpstreamDialTimeout
	}

	if channel.RetryOptions.MaxRetries <= 0 {
		channel.RetryOptions.MaxRetries = constant.DefaultChannelRetryMaxRetries
	}
	if channel.RetryOptions.RetryInterval <= 0 {
		channel.RetryOptions.RetryInterval = constant.DefaultChannelRetryInterval
	}
	if channel.RetryOptions.MaxRetryBackoff <= 0 {
		channel.RetryOptions.MaxRetryBackoff = constant.DefaultChannelRetryMaxRetryBackoff
	}
	if len(channel.RetryOptions.RetryStatuses) == 0 {
		channel.RetryOptions.RetryStatuses = constant.DefaultChannelRetryRetryStatuses
	}

	if channel.RateLimit.RequestsPerSecond <= 0 {
		channel.RateLimit.RequestsPerSecond = constant.DefaultChannelRateLimitRequestsPerSecond
	}
	if channel.RateLimit.RequestsPerMinute <= 0 {
		channel.RateLimit.RequestsPerMinute = constant.DefaultChannelRateLimitRequestsPerMinute
	}
	if channel.RateLimit.RequestsPerHour <= 0 {
		channel.RateLimit.RequestsPerHour = constant.DefaultChannelRateLimitRequestsPerHour
	}
	if channel.RateLimit.RequestsPerDay <= 0 {
		channel.RateLimit.RequestsPerDay = constant.DefaultChannelRateLimitRequestsPerDay
	}

	if channel.ChannelPriceFactor.RequestPriceFactor <= 0 {
		channel.ChannelPriceFactor.RequestPriceFactor = constant.DefaultChannelRequestPriceFactor
	}
	if channel.ChannelPriceFactor.ResponsePriceFactor <= 0 {
		channel.ChannelPriceFactor.ResponsePriceFactor = constant.DefaultChannelResponsePriceFactor
	}
	if channel.ChannelPriceFactor.CompletionPriceFactor <= 0 {
		channel.ChannelPriceFactor.CompletionPriceFactor = constant.DefaultChannelCompletionPriceFactor
	}
	if channel.ChannelPriceFactor.CachePriceFactor <= 0 {
		channel.ChannelPriceFactor.CachePriceFactor = constant.DefaultChannelCachePriceFactor
	}

	channelRepo := cc.GetChannelRepo()
	createdChannel, err := cc.service.ChannelCreate(channelRepo, &channel)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to create channel: "+err.Error(), constant.ErrorTypeChannelPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel created successfully", constant.SuccessTypeChannelPrefix+"_create", gin.H{"channel": createdChannel})
}

// ChannelUpdate 更新渠道
func (cc *channelController) ChannelUpdate(c *gin.Context) {
	var channel dto.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelPrefix+"_update")
		return
	}

	channelRepo := cc.GetChannelRepo()
	_, err := channelRepo.GetByID(channel.ChannelID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Channel not found", constant.ErrorTypeChannelPrefix+"_update")
		return
	}

	updatedChannel, err := cc.service.ChannelUpdate(channelRepo, &channel)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update channel: "+err.Error(), constant.ErrorTypeChannelPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel updated successfully", constant.SuccessTypeChannelPrefix+"_update", gin.H{"channel": updatedChannel})
}

// ChannelSearch 搜索渠道
func (cc *channelController) ChannelSearch(c *gin.Context) {
	var channelSearch channelDto.ChannelSearchRequest
	if err := c.ShouldBindJSON(&channelSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelPrefix+"_search")
		return
	}

	channelRepo := cc.GetChannelRepo()
	channels, err := cc.service.ChannelSearch(channelRepo, &channelSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search channels: "+err.Error(), constant.ErrorTypeChannelPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channels searched successfully", constant.SuccessTypeChannelPrefix+"_search", gin.H{"channels": channels})
}

// ChannelDelete 删除渠道
func (cc *channelController) ChannelDelete(c *gin.Context) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		utils.CommonError(c, http.StatusBadRequest, "channel_id is required", constant.ErrorTypeChannelPrefix+"_delete")
		return
	}

	channelRepo := cc.GetChannelRepo()
	err := cc.service.ChannelDelete(channelRepo, channelID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete channel: "+err.Error(), constant.ErrorTypeChannelPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel deleted successfully", constant.SuccessTypeChannelPrefix+"_delete", gin.H{"message": "Channel deleted successfully"})
}
