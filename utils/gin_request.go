package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"nexus-ai/constant"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRequestBody(c *gin.Context) ([]byte, error) {
	requestBody, _ := c.Get(constant.KeyRequestBody) // 获取请求体
	if requestBody != nil {                          // 如果请求体已存在，直接返回
		return requestBody.([]byte), nil
	}

	requestBody, err := io.ReadAll(c.Request.Body) // 读取请求体
	if err != nil {
		return nil, err
	}
	_ = c.Request.Body.Close()                  // 关闭请求体
	c.Set(constant.KeyRequestBody, requestBody) // 将请求体存储在上下文中
	return requestBody.([]byte), nil            // 返回请求体
}

func UnmarshalRequestBody(c *gin.Context, v any) error { // 解析请求体为指定类型v
	requestBody, err := GetRequestBody(c)
	if err != nil {
		return err
	}
	contentType := c.Request.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		err = json.Unmarshal(requestBody, &v)
	} else {
		// skip for now
	}
	if err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	return nil
}
