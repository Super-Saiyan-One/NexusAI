package controller

import (
	"net/http"
	"nexus-ai/constant"
	channelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type ChannelGroupController interface {
	GetChannelGroupRepo() repository.ChannelGroupRepository
	ChannelGroupCreate(c *gin.Context)
	ChannelGroupUpdate(c *gin.Context)
	ChannelGroupSearch(c *gin.Context)
	ChannelGroupDelete(c *gin.Context)
}

type channelGroupController struct {
	service service.ChannelGroupService
}

func NewChannelGroupController() ChannelGroupController {
	channelGroupService := service.NewChannelGroupService()
	return &channelGroupController{service: channelGroupService}
}

func (cgc *channelGroupController) GetChannelGroupRepo() repository.ChannelGroupRepository {
	return repository.NewChannelGroupRepository(model.GetDB())
}

// ChannelGroupCreate 创建渠道组
func (cgc *channelGroupController) ChannelGroupCreate(c *gin.Context) {
	var channelGroup dto.ChannelGroup
	if err := c.ShouldBindJSON(&channelGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_create")
		return
	}

	channelGroup.ChannelGroupID = utils.GenerateRandomUUID(12)
	if channelGroup.ChannelGroupName == "" {
		channelGroup.ChannelGroupName = utils.GenerateRandomString(12)
	}

	// 设置默认值
	if channelGroup.ChannelGroupOptions.MaxConcurrentRequests <= 0 {
		channelGroup.ChannelGroupOptions.MaxConcurrentRequests = constant.DefaultChannelGroupMaxConcurrentRequests
	}
	if channelGroup.ChannelGroupOptions.DefaultLevel <= 0 {
		channelGroup.ChannelGroupOptions.DefaultLevel = constant.DefaultChannelGroupDefaultLevel
	}
	if channelGroup.ChannelGroupOptions.APIDiscount <= 0 {
		channelGroup.ChannelGroupOptions.APIDiscount = constant.DefaultChannelGroupAPIDiscount
	}
	if channelGroup.ChannelGroupPriceFactor.RequestPriceFactor <= 0 {
		channelGroup.ChannelGroupPriceFactor.RequestPriceFactor = constant.DefaultChannelGroupRequestPriceFactor
	}
	if channelGroup.ChannelGroupPriceFactor.ResponsePriceFactor <= 0 {
		channelGroup.ChannelGroupPriceFactor.ResponsePriceFactor = constant.DefaultChannelGroupResponsePriceFactor
	}
	if channelGroup.ChannelGroupPriceFactor.CompletionPriceFactor <= 0 {
		channelGroup.ChannelGroupPriceFactor.CompletionPriceFactor = constant.DefaultChannelGroupCompletionPriceFactor
	}
	if channelGroup.ChannelGroupPriceFactor.CachePriceFactor <= 0 {
		channelGroup.ChannelGroupPriceFactor.CachePriceFactor = constant.DefaultChannelGroupCachePriceFactor
	}

	channelGroupRepo := cgc.GetChannelGroupRepo()
	createdChannelGroup, err := cgc.service.ChannelGroupCreate(channelGroupRepo, &channelGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to create channel group: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel group created successfully", constant.SuccessTypeChannelGroupPrefix+"_create", gin.H{"channel_group": createdChannelGroup})
}

// ChannelGroupUpdate 更新渠道组
func (cgc *channelGroupController) ChannelGroupUpdate(c *gin.Context) {
	var channelGroup dto.ChannelGroup
	if err := c.ShouldBindJSON(&channelGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_update")
		return
	}

	channelGroupRepo := cgc.GetChannelGroupRepo()
	_, err := channelGroupRepo.GetByID(channelGroup.ChannelGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Channel group not found", constant.ErrorTypeChannelGroupPrefix+"_update")
		return
	}

	updatedChannelGroup, err := cgc.service.ChannelGroupUpdate(channelGroupRepo, &channelGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update channel group: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel group updated successfully", constant.SuccessTypeChannelGroupPrefix+"_update", gin.H{"channel_group": updatedChannelGroup})
}

// ChannelGroupSearch 搜索渠道组
func (cgc *channelGroupController) ChannelGroupSearch(c *gin.Context) {
	var channelGroupSearch channelGroupDto.ChannelGroupSearchRequest
	if err := c.ShouldBindJSON(&channelGroupSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_search")
		return
	}

	channelGroupRepo := cgc.GetChannelGroupRepo()
	channelGroups, err := cgc.service.ChannelGroupSearch(channelGroupRepo, &channelGroupSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search channel groups: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel groups searched successfully", constant.SuccessTypeChannelGroupPrefix+"_search", gin.H{"channel_groups": channelGroups})
}

// ChannelGroupDelete 删除渠道组
func (cgc *channelGroupController) ChannelGroupDelete(c *gin.Context) {
	channelGroupID := c.Param("channel_group_id")
	if channelGroupID == "" {
		utils.CommonError(c, http.StatusBadRequest, "channel_group_id is required", constant.ErrorTypeChannelGroupPrefix+"_delete")
		return
	}

	channelGroupRepo := cgc.GetChannelGroupRepo()
	err := cgc.service.ChannelGroupDelete(channelGroupRepo, channelGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete channel group: "+err.Error(), constant.ErrorTypeChannelGroupPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Channel group deleted successfully", constant.SuccessTypeChannelGroupPrefix+"_delete", gin.H{"message": "Channel group deleted successfully"})
}
