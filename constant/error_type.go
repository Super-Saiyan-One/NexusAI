package constant

const (
	ErrorTypeInternalServerError     = "nexus_ai_panic_error"               // 内部服务器错误
	ErrorTypeIPVerifyFailed          = "nexus_ai_ip_verify_failed"          // IP验证失败
	ErrorTypeTokenVerifyFailed       = "nexus_ai_token_verify_failed"       // 令牌验证失败
	ErrorTypeUserVerifyFailed        = "nexus_ai_user_verify_failed"        // 用户验证失败
	ErrorTypeQuotaVerifyFailed       = "nexus_ai_quota_verify_failed"       // 配额验证失败
	ErrorTypeModelVerifyFailed       = "nexus_ai_model_verify_failed"       // 模型验证失败
	ErrorTypeChannelDistributeFailed = "nexus_ai_channel_distribute_failed" // 渠道分配失败

	ErrorTypeBillingPrefix      = "nexus_ai_router_error_billing_"       // 计费通用错误前缀
	ErrorTypeChannelGroupPrefix = "nexus_ai_router_error_channel_group_" // 渠道组通用错误前缀
	ErrorTypeChannelPrefix      = "nexus_ai_router_error_channel_"       // 渠道通用错误前缀
	ErrorTypeConfigPrefix       = "nexus_ai_router_error_config_"        // 配置通用错误前缀
	ErrorTypeLogPrefix          = "nexus_ai_router_error_log_"           // 日志通用错误前缀
	ErrorTypeMessagePrefix      = "nexus_ai_router_error_message_"       // 消息通用错误前缀
	ErrorTypeModelGroupPrefix   = "nexus_ai_router_error_model_group_"   // 模型组通用错误前缀
	ErrorTypeModelPrefix        = "nexus_ai_router_error_model_"         // 模型通用错误前缀
	ErrorTypeMonitorPrefix      = "nexus_ai_router_error_monitor_"       // 监控通用错误前缀
	ErrorTypeNotifyPrefix       = "nexus_ai_router_error_notify_"        // 通知通用错误前缀
	ErrorTypePaymentPrefix      = "nexus_ai_router_error_payment_"       // 支付通用错误前缀
	ErrorTypeQuotaPrefix        = "nexus_ai_router_error_quota_"         // 配额通用错误前缀
	ErrorTypeRelayPrefix        = "nexus_ai_router_error_relay_"         // 转发通用错误前缀
	ErrorTypeTokenPrefix        = "nexus_ai_router_error_token_"         // 令牌通用错误前缀
	ErrorTypeUsagePrefix        = "nexus_ai_router_error_usage_"         // 使用量通用错误前缀
	ErrorTypeUserGroupPrefix    = "nexus_ai_router_error_user_group_"    // 用户组通用错误前缀
	ErrorTypeUserPrefix         = "nexus_ai_router_error_user_"          // 用户通用错误前缀
)
