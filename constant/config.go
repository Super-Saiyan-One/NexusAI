package constant

import (
	"nexus-ai/dto"
	"os"
	"time"
)

type UserIDKeyType string       // 基于string的UserIDKeyType，可以避免与其他库的冲突
type RequestIDKeyType string    // 基于string的RequestIDKeyType，可以避免与其他库的冲突
type TokenKeyType string        // 基于string的TokenKeyType，可以避免与其他库的冲突
type ModelKeyType string        // 基于string的ModelKeyType，可以避免与其他库的冲突
type UserKeyType string         // 基于string的UserKeyType，可以避免与其他库的冲突
type ChannelKeyType string      // 基于string的ChannelKeyType，可以避免与其他库的冲突
type AccessTokenKeyType string  // 基于string的AccessTokenKeyType，可以避免与其他库的冲突
type RefreshTokenKeyType string // 基于string的RefreshTokenKeyType，可以避免与其他库的冲突

const (
	FrontendPort    = "11000"
	BackendPort     = "10000"
	LogMaxCount     = 100000000
	LogDir          = "./logs"
	GitRepoURL      = "https://github.com/pandalla/NexusAI.git"
	UserIDKey       = UserIDKeyType("X-Nexus-AI-User-ID")
	RequestIDKey    = RequestIDKeyType("X-Nexus-AI-Request-ID")
	TokenKey        = TokenKeyType("X-Nexus-AI-Token")
	ModelKey        = ModelKeyType("X-Nexus-AI-Model")
	UserKey         = UserKeyType("X-Nexus-AI-User")
	ChannelKey      = ChannelKeyType("X-Nexus-AI-Channel")
	AccessTokenKey  = AccessTokenKeyType("X-Nexus-AI-Access-Token")
	RefreshTokenKey = RefreshTokenKeyType("X-Nexus-AI-Refresh-Token")
	MinimumQuota    = 0.05 // 单词请求最小配额
	KeyCharset      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberCharset   = "0123456789"
)

const ( // 默认mysql配置
	MySQLDefaultHost     = "localhost" // 默认mysql地址
	MySQLDefaultPort     = "11001"     // 默认mysql端口
	MySQLDefaultUser     = "nexus"     // 默认mysql用户
	MySQLDefaultPassword = "nexus123"  // 默认mysql密码
	MySQLDefaultDatabase = "nexus"     // 默认mysql数据库
)

const ( // 默认redis配置
	RedisDefaultHost         = "localhost" // 默认redis地址 localhost
	RedisDefaultPort         = "11002"     // 默认redis端口 11002
	RedisDefaultPassword     = "nexus123"  // 默认redis密码 nexus123
	RedisDefaultDB           = "0"         // 默认redis数据库 0
	RedisDefaultMaxPoolSize  = "10000"     // 默认redis连接池最大连接数 1e4
	RedisDefaultMinIdleConns = "100"       // 默认redis连接池最小空闲连接数 1e2
)

const ( // 默认rabbitmq配置
	RabbitMQDefaultHost      = "localhost" // 默认rabbitmq地址
	RabbitMQDefaultPort      = "11003"     // 默认rabbitmq端口
	RabbitMQDefaultUser      = "nexus"     // 默认rabbitmq用户
	RabbitMQDefaultPassword  = "nexus123"  // 默认rabbitmq密码
	RabbitMQDefaultVHost     = "nexus"     // 默认rabbitmq虚拟主机
	RabbitMQDefaultPanelPort = "11004"     // 默认rabbitmq面板端口

	RabbitMQConnectionMaxRetries      = 5
	RabbitMQConnectionRetryInterval   = 5 * time.Second
	RabbitMQConnectionRequestTimeout  = 30 * time.Second
	RabbitMQConnectionHeartbeatDelay  = 10 * time.Second
	RabbitMQConnectionMaxChannels     = 1024
	RabbitMQConnectionChannelPoolSize = 100
	RabbitMQConnectionMaxFrameSize    = 131072 // 128KB
	RabbitMQConnectionConsumerThreads = 4

	RabbitMQMessageDefaultExchange    = "nexus.default"
	RabbitMQMessageDelayedExchange    = "nexus.delayed"
	RabbitMQMessageDefaultContentType = "application/json"
	RabbitMQMessageDefaultExpiration  = 24 * time.Hour
	RabbitMQMessageMaxMessageSize     = 1048576 // 1MB
	RabbitMQMessagePrefetchCount      = 1
	RabbitMQMessageReconnectInterval  = 5 * time.Second
	RabbitMQMessageHealthCheckTimeout = 5 * time.Second

	RabbitMQDefaultCustomerNum = 3
)

const (
	RootUserName     = "root"
	RootUserEmail    = "root@nexus.ai"
	RootUserPassword = "nexusai@2025"
)

const (
	JwtSecret             = "98yq2*5&!EQxZ4haVn33^D48B&jk@F##i87PC69J$&t5JpU2t^Yxo25R#V@e&eH#"
	JwtAccessTokenExpiry  = time.Hour * 2
	JwtRefreshTokenExpiry = time.Hour * 24 * 7
)

const (
	KeyRequestBody = "key_request_body"
)

const (
	DefaultUserMaxConcurrentRequests = 100
	DefaultUserQuota                 = 1
	DefaultUserLevel                 = 1
	DefaultUserAPIDiscount           = 1
)

const (
	DefaultUserGroupMaxConcurrentRequests = 100
	DefaultUserGroupDefaultLevel          = 1
	DefaultUserGroupAPIDiscount           = 1
	DefaultUserGroupRequestPriceFactor    = 1
	DefaultUserGroupResponsePriceFactor   = 1
	DefaultUserGroupCompletionPriceFactor = 1
	DefaultUserGroupCachePriceFactor      = 1
)

const (
	DefaultModelAPIDiscount     = 1
	DefaultModelRequestPrice    = 1
	DefaultModelResponsePrice   = 1
	DefaultModelCompletionPrice = 1
	DefaultModelCachePrice      = 1
)

var (
	DefaultModelVideoPrice = dto.VideoPrice{
		"std": {"5": 1},
		"pro": {"5": 3.5},
		"diy": {"1920x1080 5": 10},
	}
	DefaultModelImagePrice = dto.ImagePrice{
		"std":       1,
		"pro":       3.5,
		"1024x1024": 5,
		"1920x1080": 10,
	}
)

const (
	DefaultModelGroupMaxConcurrentRequests = 100
	DefaultModelGroupDefaultLevel          = 1
	DefaultModelGroupAPIDiscount           = 1
	DefaultModelGroupRequestPriceFactor    = 1
	DefaultModelGroupResponsePriceFactor   = 1
	DefaultModelGroupCompletionPriceFactor = 1
	DefaultModelGroupCachePriceFactor      = 1
)

const (
	DefaultChannelUpstreamTimeout            = 120
	DefaultChannelUpstreamMaxRetries         = 5
	DefaultChannelUpstreamDialTimeout        = 60
	DefaultChannelRetryMaxRetries            = 5
	DefaultChannelRetryInterval              = 1000
	DefaultChannelRetryMaxRetryBackoff       = 10000
	DefaultChannelRateLimitRequestsPerSecond = 100
	DefaultChannelRateLimitRequestsPerMinute = 3000
	DefaultChannelRateLimitRequestsPerHour   = 100000
	DefaultChannelRateLimitRequestsPerDay    = 1000000
	DefaultChannelRequestPriceFactor         = 1
	DefaultChannelResponsePriceFactor        = 1
	DefaultChannelCompletionPriceFactor      = 1
	DefaultChannelCachePriceFactor           = 1
)

var DefaultChannelRetryRetryStatuses = []int{500, 502, 503, 504}

const (
	DefaultChannelGroupMaxConcurrentRequests = 1000
	DefaultChannelGroupDefaultLevel          = 1
	DefaultChannelGroupAPIDiscount           = 1
	DefaultChannelGroupRequestPriceFactor    = 1
	DefaultChannelGroupResponsePriceFactor   = 1
	DefaultChannelGroupCompletionPriceFactor = 1
	DefaultChannelGroupCachePriceFactor      = 1
)

const (
	DefaultTokenMaxConcurrentRequests = 100
	DefaultTokenMaxRequestsPerMinute  = 3000
	DefaultTokenMaxRequestsPerHour    = 100000
	DefaultTokenMaxRequestsPerDay     = 1000000
	DefaultTokenQuotaTotal            = 0
	DefaultTokenQuotaUsed             = 0
	DefaultTokenQuotaLeft             = 0
	DefaultTokenQuotaFrozen           = 0
)

var (
	DefaultTokenRequireSignature = new(bool) // false by default
	DefaultTokenDisableRateLimit = new(bool) // false by default
)

var (
	SecretID  string
	SecretKey string
	Bucket    string
	AppId     string
	Region    string
)

func LoadCosConfig() {
	SecretID = getString(os.Getenv("COS_SECRET_ID"), "")
	SecretKey = getString(os.Getenv("COS_SECRET_KEY"), "")
	Bucket = getString(os.Getenv("COS_BUCKET"), "")
	AppId = getString(os.Getenv("COS_APP_ID"), "")
	Region = getString(os.Getenv("COS_REGION"), "")
}

var (
	StripeKey           string
	StripeWebhookSecret string
	AlipayPriceID       string
	CardPriceID         string
	ServerAddress       string

	StripeSuccessURL string
	StripeCancelURL  string
)

func getString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadStripeConfig() {
	StripeKey = getString(os.Getenv("STRIPE_KEY"), "")
	StripeWebhookSecret = getString(os.Getenv("STRIPE_WEBHOOK_SECRET"), "")
	AlipayPriceID = getString(os.Getenv("ALIPAY_PRICE_ID"), "")
	CardPriceID = getString(os.Getenv("CARD_PRICE_ID"), "")
	ServerAddress = getString(os.Getenv("SERVER_ADDRESS"), "")
	StripeSuccessURL = getString(os.Getenv("STRIPE_SUCCESS_URL"), "")
	StripeCancelURL = getString(os.Getenv("STRIPE_CANCEL_URL"), "")
}

var (
	CheckSensitiveText bool
)

func LoadCheckSensitive() {
	CheckSensitiveText = getString(os.Getenv("CHECK_SENSITIVE_TEXT"), "false") == "true"
}
