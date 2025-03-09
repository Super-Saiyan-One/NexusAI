package controller // TODO

import (
	"fmt"
	"net/http"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/dto/relay"
	"nexus-ai/service"
	"nexus-ai/utils"

	"github.com/gin-gonic/gin"
)

type RelayV1Controller interface {
	RelayChatCompletions(c *gin.Context)
}

type relayV1Controller struct {
}

func NewRelayV1Controller() RelayV1Controller {
	return &relayV1Controller{}
}

func (r *relayV1Controller) RelayChatCompletions(c *gin.Context) {
	var openAIError *relay.GeneralErrorWithStatusCode
	channel := c.MustGet(string(constant.ChannelKey)).(*dto.Channel)
	for i := 0; i <= channel.RetryOptions.MaxRetries; i++ {
		openAIError = relayRequest(c, channel, constant.RelayModeChatCompletion)
		if openAIError == nil {
			return // 如果请求成功，直接返回
		}
		// go TODO
		if !shouldRetry(c, channel, openAIError) {
			break
		}
	}
}

func relayRequest(c *gin.Context, channel *dto.Channel, relayMode int) *relay.GeneralErrorWithStatusCode {
	var err *relay.GeneralErrorWithStatusCode
	switch relayMode {
	case constant.RelayModeChatCompletion:
		err = relayText(c, channel)
	}
	return err
}

func shouldRetry(c *gin.Context, channel *dto.Channel, openAIError *relay.GeneralErrorWithStatusCode) bool {
	for _, statusCode := range channel.RetryOptions.RetryStatuses {
		if openAIError.StatusCode == statusCode {
			// 如果请求状态码在重试状态码列表中，则重试
			return true
		}
	}
	if openAIError == nil || openAIError.IsLocal || openAIError.StatusCode == http.StatusBadRequest {
		// 如果请求成功或者请求是本地错误或者请求状态码为400，则不重试
		return false
	}
	return true
}

func relayText(c *gin.Context, channel *dto.Channel) *relay.GeneralErrorWithStatusCode {
	relayInfo := service.ParseRelayInfo(c, channel)
	textRequest, err := service.ParseRelayTextRequest(c, relayInfo)
	if err != nil {
		utils.LogError(c, fmt.Sprintf("ParseRelayTextRequest failed: %v", err.Error()))
		return service.GeneralErrorLocalHandle(c, err, "Invalid text request", http.StatusBadRequest)
	}
	model := c.MustGet(string(constant.ModelKey)).(*dto.Model)
	relayInfo.RequestModelName = model.ModelName
	if constant.CheckSensitiveText { // 如果开启敏感词检测，则检测敏感词
		err = service.CheckSensitiveMessages(textRequest.Messages)
		if err != nil {
			utils.LogError(c, fmt.Sprintf("CheckSensitiveMessages failed: %v", err.Error()))
			return service.GeneralErrorLocalHandle(c, err, "Sensitive words detected", http.StatusBadRequest)
		}
	}
	
	return nil
}
