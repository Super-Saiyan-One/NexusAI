package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL             string `yaml:"base_url"`
	APIKey              string `yaml:"api_key"`
	TimeoutSeconds      int    `yaml:"timeout_seconds"`
	MaxRetries          int    `yaml:"max_retries"`
	PollIntervalSeconds int    `yaml:"poll_interval_seconds"`
}

func (c *Config) Timeout() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

func (c *Config) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSeconds) * time.Second
}

func LoadConfig() (*Config, error) {
	// 获取配置文件路径（假设配置文件位于项目根目录/config目录下）
	configPath := "api/luma/config/config.yaml"

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证必要配置项
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required in config")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required in config")
	}

	return &cfg, nil
}
