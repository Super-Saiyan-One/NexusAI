package middleware

import (
	"net/http"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

func ModelVerifyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelRequest, err := utils.ParseModelRequest(c)
		if err != nil {
			utils.AbortWhenModelVerifyFailed(c, http.StatusBadRequest, "Invalid request, "+err.Error())
			return
		}
		token := c.MustGet(string(constant.TokenKey)).(*dto.Token)
		allowedModels := token.TokenModels.AllowedModels

		// 如果允许的模型列表不为空，则进行验证
		if len(allowedModels) > 0 {
			allowedModelsMap := utils.SliceToMap(allowedModels)
			if _, exists := allowedModelsMap[modelRequest.Model]; exists { // 如果允许的模型列表中包含请求的模型
				model, err := repository.ModelVerify(modelRequest.Model)
				if err != nil {
					utils.AbortWhenModelVerifyFailed(c, http.StatusBadRequest, err.Error())
					return
				}
				c.Set(string(constant.ModelKey), model)
				c.Next()
			} else {
				utils.AbortWhenModelVerifyFailed(c, http.StatusBadRequest, "Using token not allowed to use requested model")
				return
			}
		} else {
			utils.AbortWhenModelVerifyFailed(c, http.StatusBadRequest, "Using token not allowed to use any models")
			return
		}
	}
}
