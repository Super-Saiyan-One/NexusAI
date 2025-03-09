package controller

import (
	"net/http"
	"nexus-ai/constant"
	tokenDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type TokenController interface {
	GetTokenRepo() repository.TokenRepository
	TokenCreate(c *gin.Context)
	TokenUpdate(c *gin.Context)
	TokenSearch(c *gin.Context)
	TokenDelete(c *gin.Context)
	TokenAvailableModels(c *gin.Context)
	TokenAvailableChannels(c *gin.Context)
}

type tokenController struct {
	tokenService service.TokenService
}

func NewTokenController() TokenController {
	tokenService := service.NewTokenService()
	return &tokenController{tokenService: tokenService}
}

func (tc *tokenController) GetTokenRepo() repository.TokenRepository {
	return repository.NewTokenRepository(model.GetDB())
}

func (tc *tokenController) TokenCreate(c *gin.Context) {
	var token dto.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeTokenPrefix+"_create")
		return
	}
	token.UserID = c.GetString(string(constant.UserIDKey))
	token.TokenID = utils.GenerateRandomUUID(12)
	token.TokenKey = utils.GenerateRandomString(32)
	token.Status = 1
	if token.TokenName == "" {
		token.TokenName = utils.GenerateRandomString(12)
	}

	// 设置默认值
	if token.TokenOptions.MaxConcurrentRequests <= 0 {
		token.TokenOptions.MaxConcurrentRequests = constant.DefaultTokenMaxConcurrentRequests
	}
	if token.TokenOptions.MaxRequestsPerMinute <= 0 {
		token.TokenOptions.MaxRequestsPerMinute = constant.DefaultTokenMaxRequestsPerMinute
	}
	if token.TokenOptions.MaxRequestsPerHour <= 0 {
		token.TokenOptions.MaxRequestsPerHour = constant.DefaultTokenMaxRequestsPerHour
	}
	if token.TokenOptions.MaxRequestsPerDay <= 0 {
		token.TokenOptions.MaxRequestsPerDay = constant.DefaultTokenMaxRequestsPerDay
	}

	if token.TokenOptions.RequireSignature == nil {
		token.TokenOptions.RequireSignature = constant.DefaultTokenRequireSignature
	}
	if token.TokenOptions.DisableRateLimit == nil {
		token.TokenOptions.DisableRateLimit = constant.DefaultTokenDisableRateLimit
	}

	if token.TokenQuotaTotal <= 0 {
		token.TokenQuotaTotal = constant.DefaultTokenQuotaTotal
		token.TokenQuotaLeft = constant.DefaultTokenQuotaLeft
		token.TokenQuotaUsed = constant.DefaultTokenQuotaUsed
		token.TokenQuotaFrozen = constant.DefaultTokenQuotaFrozen
	} else {
		token.TokenQuotaLeft = token.TokenQuotaTotal
		token.TokenQuotaUsed = 0
		token.TokenQuotaFrozen = 0
	}

	tokenRepo := tc.GetTokenRepo()
	createdToken, err := tc.tokenService.TokenCreate(tokenRepo, &token)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to create token: "+err.Error(), constant.ErrorTypeTokenPrefix+"_create")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Token created successfully", constant.SuccessTypeTokenPrefix+"_create", gin.H{"token": createdToken})
}

func (tc *tokenController) TokenUpdate(c *gin.Context) {
	var token dto.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeTokenPrefix+"_update")
		return
	}
	userID := c.GetString(string(constant.UserIDKey))
	tokenRepo := tc.GetTokenRepo()
	existingToken, err := tokenRepo.GetByID(token.TokenID)
	if err != nil || existingToken.UserID != userID {
		utils.CommonError(c, http.StatusForbidden, "Unauthorized operation", constant.ErrorTypeTokenPrefix+"_update")
		return
	}

	updatedToken, err := tc.tokenService.TokenUpdate(tokenRepo, &token)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to update token: "+err.Error(), constant.ErrorTypeTokenPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Token updated successfully", constant.SuccessTypeTokenPrefix+"_update", gin.H{"token": updatedToken})
}

func (tc *tokenController) TokenSearch(c *gin.Context) {
	var tokenSearch tokenDto.TokenSearchRequest
	if err := c.ShouldBindJSON(&tokenSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeTokenPrefix+"_search")
		return
	}

	tokenRepo := tc.GetTokenRepo()
	tokens, err := tc.tokenService.TokenSearch(tokenRepo, &tokenSearch)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to search tokens: "+err.Error(), constant.ErrorTypeTokenPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Tokens searched successfully", constant.SuccessTypeTokenPrefix+"_search", gin.H{"tokens": tokens})
}

func (tc *tokenController) TokenDelete(c *gin.Context) {
	tokenID := c.Param("token_id")
	if tokenID == "" {
		utils.CommonError(c, http.StatusBadRequest, "token_id is required", constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	userID := c.GetString(string(constant.UserIDKey))
	userRepo := repository.NewUserRepository(model.GetDB())
	userGroupRepo := repository.NewUserGroupRepository(model.GetDB())
	tokenRepo := tc.GetTokenRepo()

	user, err := userRepo.GetByID(userID)
	if err != nil { // 如果获取用户失败，则认为没有权限
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	userGroup, err := userGroupRepo.GetByID(user.UserGroupID)
	if err != nil { // 如果获取用户组失败，则认为没有权限
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	targetToken, err := tokenRepo.GetByID(tokenID)
	if err != nil { // 如果获取目标token失败，则认为没有权限
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	if userGroup.UserGroupOptions.DefaultLevel != 99 && user.UserOptions.DefaultLevel != 99 && targetToken.UserID != userID { // 如果用户组和用户默认等级都不是99，且目标token的userID不等于当前用户ID，则认为没有权限
		utils.CommonError(c, http.StatusForbidden, "You are not authorized to delete this token", constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	if err := tc.tokenService.TokenDelete(tokenRepo, tokenID); err != nil {
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	// 成功
	utils.CommonSuccess(c, http.StatusOK, "Token deleted successfully", constant.SuccessTypeTokenPrefix+"_delete", gin.H{"message": "Token deleted successfully"})
}

// TokenAvailableModels 获取token可用模型
func (tc *tokenController) TokenAvailableModels(c *gin.Context) {
	tokenID := c.Param("token_id")
	if tokenID == "" {
		utils.CommonError(c, http.StatusBadRequest, "token_id is required", constant.ErrorTypeTokenPrefix+"_available_models")
		return
	}

	tokenRepo := tc.GetTokenRepo()
	models, err := tc.tokenService.TokenAvailableModels(tokenRepo, tokenID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to get token available models: "+err.Error(), constant.ErrorTypeTokenPrefix+"_available_models")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Token available models fetched successfully", constant.SuccessTypeTokenPrefix+"_available_models", gin.H{"models": models})
}

// TokenAvailableChannels 获取token可用渠道
func (tc *tokenController) TokenAvailableChannels(c *gin.Context) {
	tokenID := c.Param("token_id")
	if tokenID == "" {
		utils.CommonError(c, http.StatusBadRequest, "token_id is required", constant.ErrorTypeTokenPrefix+"_available_channels")
		return
	}

	tokenRepo := tc.GetTokenRepo()
	channels, err := tc.tokenService.TokenAvailableChannels(tokenRepo, tokenID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to get token available channels: "+err.Error(), constant.ErrorTypeTokenPrefix+"_available_channels")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Token available channels fetched successfully", constant.SuccessTypeTokenPrefix+"_available_channels", gin.H{"channels": channels})
}
