package service

import (
	"errors"
	"math"
	"nexus-ai/constant"
	dto "nexus-ai/dto/model"
	"nexus-ai/dto/relay"
	"nexus-ai/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func ParseRelayInfo(c *gin.Context, channel *dto.Channel) *relay.RelayInfo {
	relayInfo := &relay.RelayInfo{
		RelayMode:        constant.RelayModeChatCompletion,
		BaseURL:          channel.UpstreamOptions.Endpoint,
		AlternateBaseURL: channel.UpstreamOptions.ProxyURL,
		APIKey:           channel.AuthOptions.APIKey,
		APISecret:        channel.AuthOptions.APISecret,
		APIToken:         channel.AuthOptions.BearerToken,
		StartTime:        time.Now(),
	}
	return relayInfo
}

func ParseRelayTextRequest(c *gin.Context, relayInfo *relay.RelayInfo) (*relay.OpenAIRequest, error) {
	textRequest := &relay.OpenAIRequest{}
	err := utils.UnmarshalRequestBody(c, textRequest)
	if err != nil {
		return nil, err
	}
	if relayInfo.RelayMode == constant.RelayModeModeration && textRequest.Model == "" {
		textRequest.Model = "text-moderation-latest"
	}
	if relayInfo.RelayMode == constant.RelayModeEmbedding && textRequest.Model == "" {
		textRequest.Model = c.Param("model")
	}
	if textRequest.Model == "" {
		return nil, errors.New("model is required")
	}
	if textRequest.MaxTokens <= 0 || textRequest.MaxTokens > math.MaxInt32/2 {
		return nil, errors.New("max_tokens is invalid")
	}
	switch relayInfo.RelayMode {
	case constant.RelayModeCompletion:
		if textRequest.Prompt == "" {
			return nil, errors.New("field prompt is required")
		}
	case constant.RelayModeChatCompletion:
		if len(textRequest.Messages) == 0 {
			return nil, errors.New("field messages is required")
		}
	case constant.RelayModeEmbedding:
	case constant.RelayModeModeration:
		if textRequest.Input == "" {
			return nil, errors.New("field input is required")
		}
	case constant.RelayModeEdit:
		if textRequest.Instruction == "" {
			return nil, errors.New("field instruction is required")
		}
	}
	relayInfo.IsStream = textRequest.Stream
	return textRequest, nil
}

func checkTextRequestSensitive(textRequest *relay.OpenAIRequest, relayInfo *relay.RelayInfo) error {
	var err error
	switch relayInfo.RelayMode {
	case constant.RelayModeChatCompletion:
		err = CheckSensitiveMessages(textRequest.Messages)
	case constant.RelayModeCompletion:
		err = CheckSensitiveInput(textRequest.Prompt)
	case constant.RelayModeModeration:
		err = CheckSensitiveInput(textRequest.Input)
	case constant.RelayModeEmbedding:
		err = CheckSensitiveInput(textRequest.Input)
	}
	return err
}
