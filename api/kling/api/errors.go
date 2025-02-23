package api

import "fmt"

// APIError 实现 GetCode(), GetMessage(), GetRequestId() 接口
type APIError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// 实现 error 接口
func (e *APIError) Error() string {
	return fmt.Sprintf("code=%d, message=%s, request_id=%s", e.Code, e.Message, e.RequestID)
}

// 新增 Get 方法
func (e *APIError) GetCode() int         { return e.Code }
func (e *APIError) GetMessage() string   { return e.Message }
func (e *APIError) GetRequestID() string { return e.RequestID } // 注意方法名和字段名是否一致！
