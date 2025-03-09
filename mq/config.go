package mq

import (
	"fmt"
	"nexus-ai/constant"
	"nexus-ai/utils"
	"strconv"
	"time"
)

// Config RabbitMQ配置结构体
type Config struct {
	// 基础配置组
	Basic struct {
		Host      string `json:"host" yaml:"host"`
		Port      int    `json:"port" yaml:"port"`
		User      string `json:"user" yaml:"user"`
		Password  string `json:"password" yaml:"password"`
		VHost     string `json:"vhost" yaml:"vhost"`
		PanelPort string `json:"panelPort" yaml:"panelPort"`
	}

	// 连接配置组
	Connection struct {
		MaxRetries      int           `json:"maxRetries" yaml:"maxRetries"`
		RetryInterval   time.Duration `json:"retryInterval" yaml:"retryInterval"`
		RequestTimeout  time.Duration `json:"requestTimeout" yaml:"requestTimeout"`
		HeartbeatDelay  time.Duration `json:"heartbeatDelay" yaml:"heartbeatDelay"`
		MaxChannels     int           `json:"maxChannels" yaml:"maxChannels"`
		ChannelPoolSize int           `json:"channelPoolSize" yaml:"channelPoolSize"`
		MaxFrameSize    int           `json:"maxFrameSize" yaml:"maxFrameSize"`
		ConsumerThreads int           `json:"consumerThreads" yaml:"consumerThreads"`
	}

	// 消息配置组
	Message struct {
		DefaultExchange    string        `json:"defaultExchange" yaml:"defaultExchange"`
		DelayedExchange    string        `json:"delayedExchange" yaml:"delayedExchange"`
		DefaultContentType string        `json:"defaultContentType" yaml:"defaultContentType"`
		DefaultExpiration  time.Duration `json:"defaultExpiration" yaml:"defaultExpiration"`
		MaxMessageSize     int           `json:"maxMessageSize" yaml:"maxMessageSize"`
		PrefetchCount      int           `json:"prefetchCount" yaml:"prefetchCount"`
		ReconnectInterval  time.Duration `json:"reconnectInterval" yaml:"reconnectInterval"`
		HealthCheckTimeout time.Duration `json:"healthCheckTimeout" yaml:"healthCheckTimeout"`
	}
}

// Validate 验证配置的合法性
func (c *Config) Validate() error {
	// 验证基础配置
	if c.Basic.Host == "" {
		return fmt.Errorf("rabbitmq host cannot be empty")
	}
	if c.Basic.Port < 1 || c.Basic.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Basic.Port)
	}
	if c.Basic.User == "" {
		return fmt.Errorf("rabbitmq user cannot be empty")
	}

	// 验证连接配置
	if c.Connection.MaxRetries <= 0 {
		return fmt.Errorf("max retries must be positive")
	}
	if c.Connection.RetryInterval <= 0 {
		return fmt.Errorf("retry interval must be positive")
	}
	if c.Connection.MaxChannels <= 0 {
		return fmt.Errorf("max channels must be positive")
	}
	if c.Connection.ConsumerThreads <= 0 {
		return fmt.Errorf("consumer threads must be positive")
	}

	// 验证消息配置
	if c.Message.DefaultExchange == "" {
		return fmt.Errorf("default exchange cannot be empty")
	}
	if c.Message.DelayedExchange == "" {
		return fmt.Errorf("delayed exchange cannot be empty")
	}
	if c.Message.MaxMessageSize <= 0 {
		return fmt.Errorf("max message size must be positive")
	}
	if c.Message.PrefetchCount <= 0 {
		return fmt.Errorf("prefetch count must be positive")
	}

	return nil
}

// DefaultConfig 返回默认的RabbitMQ配置
func DefaultConfig() *Config {
	cfg := &Config{}

	// 解析基础配置
	port, err := strconv.Atoi(utils.GetEnv("RABBITMQ_PORT", constant.RabbitMQDefaultPort))
	if err != nil {
		port, _ = strconv.Atoi(constant.RabbitMQDefaultPort)
	}

	cfg.Basic.Host = utils.GetEnv("RABBITMQ_HOST", constant.RabbitMQDefaultHost)
	cfg.Basic.Port = port
	cfg.Basic.User = utils.GetEnv("RABBITMQ_USER", constant.RabbitMQDefaultUser)
	cfg.Basic.Password = utils.GetEnv("RABBITMQ_PASSWORD", constant.RabbitMQDefaultPassword)
	cfg.Basic.VHost = utils.GetEnv("RABBITMQ_VHOST", constant.RabbitMQDefaultVHost)
	cfg.Basic.PanelPort = utils.GetEnv("RABBITMQ_PANEL_PORT", constant.RabbitMQDefaultPanelPort)

	utils.SysInfo(fmt.Sprintf("RabbitMQ | host: %s | port: %d | user: %s | vhost: %s",
		cfg.Basic.Host, cfg.Basic.Port, cfg.Basic.User, cfg.Basic.VHost))

	// 设置连接配置
	cfg.Connection.MaxRetries = constant.RabbitMQConnectionMaxRetries
	cfg.Connection.RetryInterval = constant.RabbitMQConnectionRetryInterval
	cfg.Connection.RequestTimeout = constant.RabbitMQConnectionRequestTimeout
	cfg.Connection.HeartbeatDelay = constant.RabbitMQConnectionHeartbeatDelay
	cfg.Connection.MaxChannels = constant.RabbitMQConnectionMaxChannels
	cfg.Connection.ChannelPoolSize = constant.RabbitMQConnectionChannelPoolSize
	cfg.Connection.MaxFrameSize = constant.RabbitMQConnectionMaxFrameSize
	cfg.Connection.ConsumerThreads = constant.RabbitMQConnectionConsumerThreads

	// 设置消息配置
	cfg.Message.DefaultExchange = constant.RabbitMQMessageDefaultExchange
	cfg.Message.DelayedExchange = constant.RabbitMQMessageDelayedExchange
	cfg.Message.DefaultContentType = constant.RabbitMQMessageDefaultContentType
	cfg.Message.DefaultExpiration = constant.RabbitMQMessageDefaultExpiration
	cfg.Message.MaxMessageSize = constant.RabbitMQMessageMaxMessageSize
	cfg.Message.PrefetchCount = constant.RabbitMQMessagePrefetchCount
	cfg.Message.ReconnectInterval = constant.RabbitMQMessageReconnectInterval
	cfg.Message.HealthCheckTimeout = constant.RabbitMQMessageHealthCheckTimeout

	return cfg
}
