package controller

import (
	"net/http"
	"nexus-ai/constant"
	modelGroupDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type ModelGroupController interface {
	GetModelGroupRepo() repository.ModelGroupRepository
	ModelGroupCreate(c *gin.Context)
	ModelGroupUpdate(c *gin.Context)
	ModelGroupSearch(c *gin.Context)
	ModelGroupDelete(c *gin.Context)
}

type modelGroupController struct {
	service service.ModelGroupService
}

func NewModelGroupController() ModelGroupController {
	modelGroupService := service.NewModelGroupService()
	return &modelGroupController{service: modelGroupService}
}

func (mgc *modelGroupController) GetModelGroupRepo() repository.ModelGroupRepository {
	return repository.NewModelGroupRepository(model.GetDB())
}

// ModelGroupCreate 创建模型组
func (mgc *modelGroupController) ModelGroupCreate(c *gin.Context) {
	var modelGroup dto.ModelGroup
	if err := c.ShouldBindJSON(&modelGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_create")
		return
	}

	modelGroup.ModelGroupID = utils.GenerateRandomUUID(12)
	if modelGroup.ModelGroupName == "" {
		modelGroup.ModelGroupName = utils.GenerateRandomString(12)
	}

	// 设置默认值
	if modelGroup.ModelGroupOptions.MaxConcurrentRequests <= 0 {
		modelGroup.ModelGroupOptions.MaxConcurrentRequests = constant.DefaultModelGroupMaxConcurrentRequests
	}
	if modelGroup.ModelGroupOptions.DefaultLevel <= 0 {
		modelGroup.ModelGroupOptions.DefaultLevel = constant.DefaultModelGroupDefaultLevel
	}
	if modelGroup.ModelGroupOptions.APIDiscount <= 0 {
		modelGroup.ModelGroupOptions.APIDiscount = constant.DefaultModelGroupAPIDiscount
	}
	if modelGroup.ModelGroupPriceFactor.RequestPriceFactor <= 0 {
		modelGroup.ModelGroupPriceFactor.RequestPriceFactor = constant.DefaultModelGroupRequestPriceFactor
	}
	if modelGroup.ModelGroupPriceFactor.ResponsePriceFactor <= 0 {
		modelGroup.ModelGroupPriceFactor.ResponsePriceFactor = constant.DefaultModelGroupResponsePriceFactor
	}
	if modelGroup.ModelGroupPriceFactor.CompletionPriceFactor <= 0 {
		modelGroup.ModelGroupPriceFactor.CompletionPriceFactor = constant.DefaultModelGroupCompletionPriceFactor
	}
	if modelGroup.ModelGroupPriceFactor.CachePriceFactor <= 0 {
		modelGroup.ModelGroupPriceFactor.CachePriceFactor = constant.DefaultModelGroupCachePriceFactor
	}

	modelGroupRepo := mgc.GetModelGroupRepo()
	createdModelGroup, err := mgc.service.ModelGroupCreate(modelGroupRepo, &modelGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to create model group: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model group created successfully", constant.SuccessTypeModelGroupPrefix+"_create", gin.H{"model_group": createdModelGroup})
}

// ModelGroupUpdate 更新模型组
func (mgc *modelGroupController) ModelGroupUpdate(c *gin.Context) {
	var modelGroup dto.ModelGroup
	if err := c.ShouldBindJSON(&modelGroup); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_update")
		return
	}

	modelGroupRepo := mgc.GetModelGroupRepo()
	_, err := modelGroupRepo.GetByID(modelGroup.ModelGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Model group not found", constant.ErrorTypeModelGroupPrefix+"_update")
		return
	}

	updatedModelGroup, err := mgc.service.ModelGroupUpdate(modelGroupRepo, &modelGroup)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update model group: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model group updated successfully", constant.SuccessTypeModelGroupPrefix+"_update", gin.H{"model_group": updatedModelGroup})
}

// ModelGroupSearch 搜索模型组
func (mgc *modelGroupController) ModelGroupSearch(c *gin.Context) {
	var modelGroupSearch modelGroupDto.ModelGroupSearchRequest
	if err := c.ShouldBindJSON(&modelGroupSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_search")
		return
	}

	modelGroupRepo := mgc.GetModelGroupRepo()
	modelGroups, err := mgc.service.ModelGroupSearch(modelGroupRepo, &modelGroupSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search model groups: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model groups searched successfully", constant.SuccessTypeModelGroupPrefix+"_search", gin.H{"model_groups": modelGroups})
}

// ModelGroupDelete 删除模型组
func (mgc *modelGroupController) ModelGroupDelete(c *gin.Context) {
	modelGroupID := c.Param("model_group_id")
	if modelGroupID == "" {
		utils.CommonError(c, http.StatusBadRequest, "model_group_id is required", constant.ErrorTypeModelGroupPrefix+"_delete")
		return
	}

	modelGroupRepo := mgc.GetModelGroupRepo()
	err := mgc.service.ModelGroupDelete(modelGroupRepo, modelGroupID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete model group: "+err.Error(), constant.ErrorTypeModelGroupPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model group deleted successfully", constant.SuccessTypeModelGroupPrefix+"_delete", gin.H{"message": "Model group deleted successfully"})
}
