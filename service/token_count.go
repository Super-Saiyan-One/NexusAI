package service

import (
	"fmt"
	"nexus-ai/utils"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

var defaultTokenEncoder *tiktoken.Tiktoken
var gpt4TokenEncoder *tiktoken.Tiktoken
var gpt4oTokenEncoder *tiktoken.Tiktoken

func InitTokenEncoder() {
	utils.SysInfo("Initalizing token encoders")
	gpt35TurboTokenEncoder, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		utils.FatalLog(fmt.Sprintf("Failed to initalize gpt-3.5-turbo token encoder: %v", err.Error()))
	}
	defaultTokenEncoder = gpt35TurboTokenEncoder // 默认使用gpt-3.5-turbo的编码器
	gpt4TokenEncoder, err = tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		utils.FatalLog(fmt.Sprintf("Failed to initalize gpt-4 token encoder: %v", err.Error()))
	}
	gpt4oTokenEncoder, err = tiktoken.EncodingForModel("gpt-4o")
	if err != nil {
		utils.FatalLog(fmt.Sprintf("Failed to initalize gpt-4o token encoder: %v", err.Error()))
	}
	utils.SysInfo("Token encoders initalized")
}

func getDefaultTokenEncoderByModel(model string) *tiktoken.Tiktoken {
	if strings.HasPrefix(model, "gpt-4o") || strings.HasPrefix(model, "chatgpt-4o") {
		return gpt4oTokenEncoder
	} else if strings.HasPrefix(model, "gpt-4") || strings.HasPrefix(model, "chatgpt-4") {
		return gpt4TokenEncoder // 如果模型名包含gpt-4或chatgpt-4，则使用gpt-4的编码器
	}
	// TODO
	return defaultTokenEncoder
}
