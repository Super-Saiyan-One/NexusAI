package utils

import (
	"fmt"
	"nexus-ai/constant"

	"github.com/gin-gonic/gin"
)

func CommonError(c *gin.Context, statusCode int, message string, typeString string) {
	requestID := c.GetString(string(constant.RequestIDKey))
	c.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
			"type":    typeString,
		},
	})
	LogError(c.Request.Context(), fmt.Sprintf("requestID: %s, type: %s, message: %s", requestID, typeString, message))
}

func CommonSuccess(c *gin.Context, statusCode int, message string, typeString string, data gin.H) {
	requestID := c.GetString(string(constant.RequestIDKey))
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
	LogInfo(c.Request.Context(), fmt.Sprintf("requestID: %s, type: %s, message: %s", requestID, typeString, message))
}

func CommonWarn(c *gin.Context, statusCode int, message string, typeString string) {
	requestID := c.GetString(string(constant.RequestIDKey))
	c.JSON(statusCode, gin.H{
		"success": true,
		"warn": gin.H{
			"message": message,
			"type":    typeString,
		},
	})
	LogWarn(c.Request.Context(), fmt.Sprintf("requestID: %s, type: %s, message: %s", requestID, typeString, message))
}
