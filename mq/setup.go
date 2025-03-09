package mq

import (
	"context"
	"fmt"
	"nexus-ai/utils"
	"time"
)

var (
	stopChan               chan struct{}
	consecutiveFailures    int
	maxConsecutiveFailures = 3
)

func Setup() error {
	cfg := DefaultConfig()
	if err := InitMQ(cfg); err != nil {
		return fmt.Errorf("failed to initialize rabbitmq: %v", err)
	}
	stopChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := HealthCheck(ctx); err != nil {
					consecutiveFailures++
					utils.SysError(fmt.Sprintf("RabbitMQ | health check failed (%d/%d): %v",
						consecutiveFailures, maxConsecutiveFailures, err))

					if consecutiveFailures >= maxConsecutiveFailures {
						utils.SysError("RabbitMQ | attempting to reconnect...")
						if err := InitMQ(cfg); err != nil {
							utils.SysError(fmt.Sprintf("RabbitMQ | reconnection failed: %v", err))
						} else {
							utils.SysInfo("RabbitMQ | reconnection successful")
							consecutiveFailures = 0
						}
					}
				} else {
					if consecutiveFailures > 0 {
						utils.SysInfo("RabbitMQ | health check recovered")
					} else {
						utils.SysInfo("RabbitMQ | health check passed")
					}
					consecutiveFailures = 0
				}
				cancel()
			}
		}
	}()

	return nil
}

// Shutdown 关闭RabbitMQ连接
func Shutdown() error {
	// 停止健康检查
	if stopChan != nil {
		close(stopChan)
	}

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return GracefulClose(ctx)
}
