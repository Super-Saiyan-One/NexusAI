package controller

import (
	"net/http"
	"nexus-ai/constant"
	userGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type UserGroupController interface {
	GetUserGroupRepo() repository.UserGroupRepository
	UserGroupCreate(c *gin.Context)
	UserGroupUpdate(c *gin.Context)
	UserGroupSearch(c *gin.Context)
	UserGroupDelete(c *gin.Context)
}

type userGroupController struct {
	service service.UserGroupService
}

func NewUserGroupController() UserGroupController {
	userGroupService := service.NewUserGroupService()
	return &userGroupController{service: userGroupService}
}

func (uc *userGroupController) GetUserGroupRepo() repository.UserGroupRepository {
	return repository.NewUserGroupRepository(model.GetDB())
}

// UserGroupCreate 创建用户组
func (ugc *userGroupController) UserGroupCreate(c *gin.Context) {
	var userGroup dto.UserGroup
	if err := c.ShouldBindJSON(&userGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_create")
		return
	}

	userGroup.UserGroupID = utils.GenerateRandomUUID(12)
	if userGroup.UserGroupName == "" {
		userGroup.UserGroupName = utils.GenerateRandomString(12)
	}

	// 设置默认最大并发请求数、默认等级、API折扣、价格系数
	if userGroup.UserGroupOptions.MaxConcurrentRequests <= 0 {
		userGroup.UserGroupOptions.MaxConcurrentRequests = constant.DefaultUserGroupMaxConcurrentRequests
	}
	if userGroup.UserGroupOptions.DefaultLevel <= 0 {
		userGroup.UserGroupOptions.DefaultLevel = constant.DefaultUserGroupDefaultLevel
	}
	if userGroup.UserGroupOptions.APIDiscount <= 0 {
		userGroup.UserGroupOptions.APIDiscount = constant.DefaultUserGroupAPIDiscount
	}
	if userGroup.UserGroupPriceFactor.RequestPriceFactor <= 0 {
		userGroup.UserGroupPriceFactor.RequestPriceFactor = constant.DefaultUserGroupRequestPriceFactor
	}
	if userGroup.UserGroupPriceFactor.ResponsePriceFactor <= 0 {
		userGroup.UserGroupPriceFactor.ResponsePriceFactor = constant.DefaultUserGroupResponsePriceFactor
	}
	if userGroup.UserGroupPriceFactor.CompletionPriceFactor <= 0 {
		userGroup.UserGroupPriceFactor.CompletionPriceFactor = constant.DefaultUserGroupCompletionPriceFactor
	}
	if userGroup.UserGroupPriceFactor.CachePriceFactor <= 0 {
		userGroup.UserGroupPriceFactor.CachePriceFactor = constant.DefaultUserGroupCachePriceFactor
	}

	userGroupRepo := ugc.GetUserGroupRepo()
	createdUserGroup, err := ugc.service.UserGroupCreate(userGroupRepo, &userGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to create user group: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User group created successfully", constant.SuccessTypeUserGroupPrefix+"_create", gin.H{"user_group": createdUserGroup})
}

// UserGroupUpdate 更新用户组
func (ugc *userGroupController) UserGroupUpdate(c *gin.Context) {
	var userGroup dto.UserGroup
	if err := c.ShouldBindJSON(&userGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_update")
		return
	}

	userGroupRepo := ugc.GetUserGroupRepo()
	_, err := userGroupRepo.GetByID(userGroup.UserGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "User group not found", constant.ErrorTypeUserGroupPrefix+"_update")
		return
	}

	updatedUserGroup, err := ugc.service.UserGroupUpdate(userGroupRepo, &userGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update user group: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User group updated successfully", constant.SuccessTypeUserGroupPrefix+"_update", gin.H{"user_group": updatedUserGroup})
}

// UserGroupSearch 搜索用户组
func (ugc *userGroupController) UserGroupSearch(c *gin.Context) {
	var userGroupSearch userGroupDto.UserGroupSearchRequest
	if err := c.ShouldBindJSON(&userGroupSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_search")
		return
	}

	userGroupRepo := ugc.GetUserGroupRepo()
	userGroups, err := ugc.service.UserGroupSearch(userGroupRepo, &userGroupSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search user groups: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User groups searched successfully", constant.SuccessTypeUserGroupPrefix+"_search", gin.H{"user_groups": userGroups})
}

// UserGroupDelete 删除用户组
func (ugc *userGroupController) UserGroupDelete(c *gin.Context) {
	userGroupID := c.Param("user_group_id")
	if userGroupID == "" {
		utils.CommonError(c, http.StatusBadRequest, "user_group_id is required", constant.ErrorTypeUserGroupPrefix+"_delete")
		return
	}

	userGroupRepo := ugc.GetUserGroupRepo()
	err := ugc.service.UserGroupDelete(userGroupRepo, userGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete user group: "+err.Error(), constant.ErrorTypeUserGroupPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User group deleted successfully", constant.SuccessTypeUserGroupPrefix+"_delete", gin.H{"message": "User group deleted successfully"})
}
