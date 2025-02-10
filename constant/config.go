package constant

import "time"

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
	RabbitMQDefaultHost     = "localhost" // 默认rabbitmq地址
	RabbitMQDefaultPort     = "11003"     // 默认rabbitmq端口
	RabbitMQDefaultUser     = "nexus"     // 默认rabbitmq用户
	RabbitMQDefaultPassword = "nexus123"  // 默认rabbitmq密码
)

const (
	RootUserName = "root"
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
	DefaultMaxConcurrentRequests = 100
	DefaultUserQuota             = 1
	DefaultUserLevel             = 1
	DefaultAPIDiscount           = 1
)
