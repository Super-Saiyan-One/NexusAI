package service

import (
	"fmt"
	"nexus-ai/constant"
	"nexus-ai/dto/relay"
	"nexus-ai/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func GeneralErrorHandle(c *gin.Context, err error, message string, statusCode int, errorType string) *relay.GeneralErrorWithStatusCode {
	errContent := strings.ToLower(err.Error())
	if strings.Contains(errContent, "post") || strings.Contains(errContent, "dial") || strings.Contains(errContent, "http") {
		utils.LogError(c, fmt.Sprintf("error: %s", errContent))
	}
	generalError := relay.GeneralErrorContent{
		Message: message,
		Type:    errorType,
		Code:    statusCode,
	}
	return &relay.GeneralErrorWithStatusCode{
		StatusCode: statusCode,
		Error:      generalError,
	}
}

func GeneralErrorLocalHandle(c *gin.Context, err error, message string, statusCode int) *relay.GeneralErrorWithStatusCode {
	generalError := GeneralErrorHandle(c, err, message, statusCode, constant.ErrorTypeRelayLocal)
	generalError.IsLocal = true
	return generalError
}

func GeneralErrorUpstreamHandle(c *gin.Context, err error, message string, statusCode int) *relay.GeneralErrorWithStatusCode {
	generalError := GeneralErrorHandle(c, err, message, statusCode, constant.ErrorTypeRelayUpstream)
	generalError.IsLocal = false
	return generalError
}
