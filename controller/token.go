package controller

import (
	"net/http"
	"nexus-ai/constant"
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

	tokenRepo := tc.GetTokenRepo()
	existingToken, err := tokenRepo.GetByID(token.TokenID)
	if err != nil || existingToken.UserID != token.UserID {
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
	tokenID := c.Query("token_id")
	tokenKey := c.Query("token_key")
	userID := c.Query("user_id")

	tokenRepo := tc.GetTokenRepo()
	tokens, err := tc.tokenService.TokenSearch(tokenRepo, tokenID, tokenKey, userID)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_search")
		return
	}
	// 成功 返回查询的令牌
	utils.CommonSuccess(c, http.StatusOK, "Token searched successfully", constant.SuccessTypeTokenPrefix+"_search", gin.H{"tokens": tokens})
}

func (tc *tokenController) TokenDelete(c *gin.Context) {
	tokenID := c.Param("token_id")
	if tokenID == "" {
		utils.CommonError(c, http.StatusBadRequest, "token_id is required", constant.ErrorTypeTokenPrefix+"_delete")
		return
	}

	tokenRepo := tc.GetTokenRepo()
	if err := tc.tokenService.TokenDelete(tokenRepo, tokenID); err != nil {
		utils.CommonError(c, http.StatusInternalServerError, err.Error(), constant.ErrorTypeTokenPrefix+"_delete")
		return
	}
	// 成功
	utils.CommonSuccess(c, http.StatusOK, "Token deleted successfully", constant.SuccessTypeTokenPrefix+"_delete", gin.H{"message": "Token deleted successfully"})
}
