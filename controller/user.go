package controller

import (
	"net/http"
	"nexus-ai/constant"
	userDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUserRepo() repository.UserRepository
	UserRegister(c *gin.Context)
	UserLogin(c *gin.Context)
	UserSearch(c *gin.Context)
	UserDelete(c *gin.Context)
	UserLogout(c *gin.Context)
	UserUpdate(c *gin.Context)
	UserPassword(c *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController() UserController {
	userService := service.NewUserService()
	return &userController{service: userService}
}

func (uc *userController) GetUserRepo() repository.UserRepository {
	return repository.NewUserRepository(model.GetDB())
}

// UserRegister 用户注册
func (uc *userController) UserRegister(c *gin.Context) {
	var user dto.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserPrefix+"_register")
		return
	}
	user.UserID = utils.GenerateRandomUUID(12)
	if user.Username == "" {
		user.Username = utils.GenerateRandomString(8)
	}
	user.Password = utils.HashPassword(user.Password)
	user.Status = 1
	user.OAuthInfo = dto.OAuthInfo{}
	user.UserQuota = dto.UserQuota{
		TotalQuota:  constant.DefaultUserQuota,
		FrozenQuota: 0,
		GiftQuota:   constant.DefaultUserQuota,
	}
	user.UserOptions = dto.UserOptions{
		MaxConcurrentRequests:    constant.DefaultMaxConcurrentRequests,
		DefaultLevel:             constant.DefaultUserLevel,
		APIDiscount:              constant.DefaultAPIDiscount,
		QuotaNotifyEmail:         user.Email,
		ReceiveQuotaNotifyMail:   1,
		BillingNotifyEmail:       user.Email,
		ReceiveBillingNotifyMail: 1,
		PaymentNotifyEmail:       user.Email,
		ReceivePaymentNotifyMail: 1,
	}
	userRepo := uc.GetUserRepo()
	registeredUser, err := uc.service.UserRegister(userRepo, &user)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to register user: "+err.Error(), constant.ErrorTypeUserPrefix+"_register")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User registered successfully", constant.SuccessTypeUserPrefix+"_register", gin.H{"user": registeredUser})
}

// UserLogin 用户登录
func (uc *userController) UserLogin(c *gin.Context) {
	var loginRequest userDto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserPrefix+"_login")
		return
	}

	userRepo := uc.GetUserRepo()
	user, tokens, err := uc.service.UserLogin(userRepo, &loginRequest)
	if err != nil {
		utils.CommonError(c, http.StatusUnauthorized, "Login failed: "+err.Error(), constant.ErrorTypeUserPrefix+"_login")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Login successful", constant.SuccessTypeUserPrefix+"_login", gin.H{"user": user, "tokens": tokens})
}

// UserSearch 用户搜索
func (uc *userController) UserSearch(c *gin.Context) {
	var userSearch userDto.SearchRequest
	if err := c.ShouldBindJSON(&userSearch); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserPrefix+"_search")
		return
	}

	userRepo := uc.GetUserRepo()
	users, err := uc.service.UserSearch(userRepo, &userSearch)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to search users: "+err.Error(), constant.ErrorTypeUserPrefix+"_search")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Users searched successfully", constant.SuccessTypeUserPrefix+"_search", gin.H{"users": users})
}

// UserDelete 删除用户
func (uc *userController) UserDelete(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		utils.CommonError(c, http.StatusBadRequest, "user_id is required", constant.ErrorTypeUserPrefix+"_delete")
		return
	}

	userRepo := uc.GetUserRepo()
	err := uc.service.UserDelete(userRepo, userID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to delete user: "+err.Error(), constant.ErrorTypeUserPrefix+"_delete")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User deleted successfully", constant.SuccessTypeUserPrefix+"_delete", gin.H{"message": "User deleted successfully"})
}

// UserLogout 用户登出
func (uc *userController) UserLogout(c *gin.Context) {
	userID := c.GetString(string(constant.UserIDKey))
	if userID == "" {
		utils.CommonError(c, http.StatusBadRequest, "user_id is required", constant.ErrorTypeUserPrefix+"_logout")
		return
	}

	userRepo := uc.GetUserRepo()
	err := uc.service.UserLogout(userRepo, userID)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to logout user: "+err.Error(), constant.ErrorTypeUserPrefix+"_logout")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User logged out successfully", constant.SuccessTypeUserPrefix+"_logout", gin.H{"message": "User logged out successfully"})
}

// UserUpdate 更新用户信息
func (uc *userController) UserUpdate(c *gin.Context) {
	var user dto.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserPrefix+"_update")
		return
	}

	userRepo := uc.GetUserRepo()
	existingUser, err := userRepo.GetByID(user.UserID)
	if err != nil || existingUser.UserID != user.UserID {
		utils.CommonError(c, http.StatusForbidden, "Unauthorized operation", constant.ErrorTypeUserPrefix+"_update")
		return
	}

	updatedUser, err := uc.service.UserUpdate(userRepo, &user)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update user: "+err.Error(), constant.ErrorTypeUserPrefix+"_update")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "User updated successfully", constant.SuccessTypeUserPrefix+"_update", gin.H{"user": updatedUser})
}

// UserPassword 更新密码
func (uc *userController) UserPassword(c *gin.Context) {
	var passwordRequest userDto.PasswordRequest

	if err := c.ShouldBindJSON(&passwordRequest); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypeUserPrefix+"_update_password")
		return
	}

	userID := c.GetString(string(constant.UserIDKey))
	userRepo := uc.GetUserRepo()
	err := uc.service.UserPassword(userRepo, userID, passwordRequest.OldPassword, passwordRequest.NewPassword)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to update password: "+err.Error(), constant.ErrorTypeUserPrefix+"_update_password")
		return
	}

	utils.CommonSuccess(c, http.StatusOK, "Password updated successfully", constant.SuccessTypeUserPrefix+"_update_password", gin.H{})
}
