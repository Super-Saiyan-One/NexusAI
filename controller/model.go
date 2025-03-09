package controller

import (
	"net/http"
	"nexus-ai/constant"
	modelDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/redis"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type ModelController interface {
	GetModelRepo() repository.ModelRepository
	ModelCreate(c *gin.Context)
	ModelUpdate(c *gin.Context)
	ModelSearch(c *gin.Context)
	ModelDelete(c *gin.Context)
	ModelAvailableChannels(c *gin.Context)
}

type modelController struct {
	service service.ModelService
}

func NewModelController() ModelController {
	modelService := service.NewModelService()
	return &modelController{service: modelService}
}

func (mc *modelController) GetModelRepo() repository.ModelRepository {
	return repository.NewModelRepository(model.GetDB())
}

// ModelCreate 创建模型
func (mc *modelController) ModelCreate(c *gin.Context) {
	var model dto.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelPrefix+"_create")
		return
	}

	model.ModelID = utils.GenerateRandomUUID(12)
	if model.ModelGroupID == "" {
		model.ModelGroupID, _ = redis.Get(c, "ordinary_model_group_id")
	}
	if model.ModelName == "" {
		model.ModelName = utils.GenerateRandomString(12)
	}

	// 设置默认值
	if model.ModelOptions.APIDiscount <= 0 {
		model.ModelOptions.APIDiscount = constant.DefaultModelAPIDiscount
	}

	if model.ModelPrice.RequestPrice <= 0 {
		model.ModelPrice.RequestPrice = constant.DefaultModelRequestPrice
	}
	if model.ModelPrice.ResponsePrice <= 0 {
		model.ModelPrice.ResponsePrice = constant.DefaultModelResponsePrice
	}
	if model.ModelPrice.CompletionPrice <= 0 {
		model.ModelPrice.CompletionPrice = constant.DefaultModelCompletionPrice
	}
	if model.ModelPrice.CachePrice <= 0 {
		model.ModelPrice.CachePrice = constant.DefaultModelCachePrice
	}

	modelRepo := mc.GetModelRepo()
	createdModel, err := mc.service.ModelCreate(modelRepo, &model)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to create model: "+err.Error(), constant.ErrorTypeModelPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model created successfully", constant.SuccessTypeModelPrefix+"_create", gin.H{"model": createdModel})
}

// ModelUpdate 更新模型
func (mc *modelController) ModelUpdate(c *gin.Context) {
	var model dto.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelPrefix+"_update")
		return
	}

	modelRepo := mc.GetModelRepo()
	_, err := modelRepo.GetByID(model.ModelID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Model not found", constant.ErrorTypeModelPrefix+"_update")
		return
	}

	updatedModel, err := mc.service.ModelUpdate(modelRepo, &model)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update model: "+err.Error(), constant.ErrorTypeModelPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model updated successfully", constant.SuccessTypeModelPrefix+"_update", gin.H{"model": updatedModel})
}

// ModelSearch 搜索模型
func (mc *modelController) ModelSearch(c *gin.Context) {
	var modelSearch modelDto.ModelSearchRequest
	if err := c.ShouldBindJSON(&modelSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeModelPrefix+"_search")
		return
	}

	modelRepo := mc.GetModelRepo()
	models, err := mc.service.ModelSearch(modelRepo, &modelSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search models: "+err.Error(), constant.ErrorTypeModelPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Models searched successfully", constant.SuccessTypeModelPrefix+"_search", gin.H{"models": models})
}

// ModelDelete 删除模型
func (mc *modelController) ModelDelete(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		utils.CommonError(c, http.StatusBadRequest, "model_id is required", constant.ErrorTypeModelPrefix+"_delete")
		return
	}

	modelRepo := mc.GetModelRepo()
	err := mc.service.ModelDelete(modelRepo, modelID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete model: "+err.Error(), constant.ErrorTypeModelPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model deleted successfully", constant.SuccessTypeModelPrefix+"_delete", gin.H{"message": "Model deleted successfully"})
}

// ModelAvailableChannels 获取模型可用渠道
func (mc *modelController) ModelAvailableChannels(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		utils.CommonError(c, http.StatusBadRequest, "model_id is required", constant.ErrorTypeModelPrefix+"_available_channels")
		return
	}

	modelRepo := mc.GetModelRepo()
	channels, err := mc.service.ModelAvailableChannels(modelRepo, modelID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to get model available channels: "+err.Error(), constant.ErrorTypeModelPrefix+"_available_channels")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Model available channels fetched successfully", constant.SuccessTypeModelPrefix+"_available_channels", gin.H{"channels": channels})
}
