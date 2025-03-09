package mq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nexus-ai/utils"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn             *amqp.Connection
	channelPool      *ChannelPool
	mu               sync.RWMutex
	conf             *Config
	consumers        = make(map[string]*Consumer)
	consumersMu      sync.RWMutex
	reconnectBackoff = 1 * time.Second
)

// Message 消息结构体
type Message struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Data        map[string]interface{} `json:"data"`
	DelayedTime time.Time              `json:"delayed_time,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Consumer 消费者结构体
type Consumer struct {
	Queue    string
	Handler  func(Message) error
	Channel  *amqp.Channel
	Done     chan struct{}
	IsClosed bool
}

// GetChannel 获取RabbitMQ通道
func GetChannel() (*amqp.Channel, error) {
	if channelPool == nil {
		return nil, fmt.Errorf("channel pool is nil")
	}
	return channelPool.Get()
}

// HealthCheck 健康检查
func HealthCheck(ctx context.Context) error {
	mu.RLock()
	defer mu.RUnlock()

	if conn == nil || conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}

	if channelPool == nil || channelPool.closed {
		return fmt.Errorf("rabbitmq channel pool is closed")
	}

	return nil
}

// CreateVHost 创建虚拟主机
func CreateVHost(user, host, password, vhost, panelPort string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%s/api/vhosts/%s", user, password, host, panelPort, vhost)
	// 创建虚拟主机的请求体
	vhostData := map[string]interface{}{
		"name": vhost,
	}

	data, err := json.Marshal(vhostData)
	if err != nil {
		return fmt.Errorf("failed to marshal vhost data: %v", err)
	}

	// 发送HTTP PUT请求以创建虚拟主机
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create vhost, status code: %d", resp.StatusCode)
	}

	return nil
}

// SetPermissions 设置用户权限
func SetPermissions(user, host, password, vhost, panelPort string) error {
	url := fmt.Sprintf("http://%s:%s@%s:%s/api/permissions/%s/%s", user, password, host, panelPort, user, vhost)

	// 创建权限的请求体
	permissionsData := map[string]interface{}{
		"configure": ".*",
		"write":     ".*",
		"read":      ".*",
	}

	data, err := json.Marshal(permissionsData)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions data: %v", err)
	}

	// 发送HTTP PUT请求以设置权限
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to set permissions, status code: %d", resp.StatusCode)
	}

	return nil
}

// InitMQ 初始化RabbitMQ连接
func InitMQ(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	conf = cfg
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid rabbitmq config: %v", err)
	}

	// 创建虚拟主机
	if err := CreateVHost(cfg.Basic.User, cfg.Basic.Host, cfg.Basic.Password, cfg.Basic.VHost, cfg.Basic.PanelPort); err != nil {
		utils.SysError(fmt.Sprintf("Failed to create vhost: %v", err))
	}

	// 设置权限
	if err := SetPermissions(cfg.Basic.User, cfg.Basic.Host, cfg.Basic.Password, cfg.Basic.VHost, cfg.Basic.PanelPort); err != nil {
		utils.SysError(fmt.Sprintf("Failed to set permissions: %v", err))
	}

	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Basic.User,
		cfg.Basic.Password,
		cfg.Basic.Host,
		cfg.Basic.Port,
		cfg.Basic.VHost,
	)

	var lastErr error
	// 使用指数退避重试连接
	for i := 0; i < cfg.Connection.MaxRetries; i++ {
		if err := connect(url); err != nil {
			lastErr = err
			sleepTime := time.Duration(1<<uint(i)) * cfg.Connection.RetryInterval
			if sleepTime > 30*time.Second {
				sleepTime = 30 * time.Second
			}
			time.Sleep(sleepTime)
			continue
		}

		// 初始化Channel池
		pool, err := NewChannelPool(cfg.Connection.MaxChannels)
		if err != nil {
			return fmt.Errorf("failed to create channel pool: %v", err)
		}
		channelPool = pool

		return setupExchanges()
	}

	return fmt.Errorf("failed to connect to rabbitmq after %d attempts: %v",
		cfg.Connection.MaxRetries, lastErr)
}

// connect 建立RabbitMQ连接
func connect(url string) error {
	connection, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	mu.Lock()
	conn = connection
	mu.Unlock()

	// 监听连接关闭
	go func() {
		<-connection.NotifyClose(make(chan *amqp.Error))
		mu.Lock()
		conn = nil
		mu.Unlock()

		// 重新连接
		time.Sleep(reconnectBackoff)
		if err := InitMQ(conf); err != nil {
			utils.SysError(fmt.Sprintf("Failed to reconnect to RabbitMQ: %v", err))
		}

		// 重新启动所有消费者
		restartConsumers()
	}()

	return nil
}

// setupExchange 设置单个交换机
func setupExchange(exchange string, isDelayed bool) error {
	if exchange == "" {
		return fmt.Errorf("empty exchange name")
	}
	ch, err := GetChannel()
	if err != nil {
		return fmt.Errorf("channel is nil")
	}
	defer channelPool.Put(ch)

	// 声明交换机
	if isDelayed {
		args := make(amqp.Table)
		args["x-delayed-type"] = "direct"
		err = ch.ExchangeDeclare(
			exchange,
			"x-delayed-message",
			true,
			false,
			false,
			false,
			args,
		)
	} else {
		err = ch.ExchangeDeclare(
			exchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
	}
	return err
}

// setupExchanges 设置交换机
func setupExchanges() error {
	ch, err := GetChannel()
	if err != nil {
		return fmt.Errorf("channel is nil")
	}
	defer channelPool.Put(ch)

	// 声明默认交换机
	err = ch.ExchangeDeclare(
		conf.Message.DefaultExchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// 声明延迟交换机
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err = ch.ExchangeDeclare(
		conf.Message.DelayedExchange,
		"x-delayed-message",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		args,  // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

// PublishMessage 发布消息
func PublishMessage(exchange, routingKey string, msg Message) error {
	ch, err := GetChannel()
	if err != nil {
		return fmt.Errorf("channel is nil")
	}
	defer channelPool.Put(ch)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	headers := make(amqp.Table)
	if !msg.DelayedTime.IsZero() { // 延迟消息
		delay := time.Until(msg.DelayedTime)
		if delay > 0 {
			headers["x-delay"] = int64(delay.Milliseconds())
		}
		if exchange != conf.Message.DelayedExchange { // 如果交换机不是默认延迟交换机，则创建新的延迟交换机
			err = setupExchange(exchange, true)
			if err != nil {
				utils.SysError(fmt.Sprintf("Failed to setup new delayed exchange %s: %v", exchange, err))
				exchange = conf.Message.DelayedExchange // 使用默认的延迟交换机
			}
		}
	} else { // 非延迟消息
		if exchange != conf.Message.DefaultExchange { // 如果交换机不是默认交换机，则创建新的默认交换机
			err = setupExchange(exchange, false)
			if err != nil {
				utils.SysError(fmt.Sprintf("Failed to setup new exchange %s: %v", exchange, err))
				exchange = conf.Message.DefaultExchange // 使用默认交换机
			}
		}
	}

	err = ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  conf.Message.DefaultContentType,
			Body:         body,
			Headers:      headers,
			Expiration:   fmt.Sprintf("%d", int64(conf.Message.DefaultExpiration.Seconds()*1000)),
			Timestamp:    time.Now(),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		utils.SysError(fmt.Sprintf("Failed to publish message: %v", err))
	}

	return err
}

// CreateQueue 创建队列
func CreateQueue(queue string, isDelayed bool) error {
	ch, err := GetChannel()
	if err != nil {
		return fmt.Errorf("channel is nil")
	}
	defer channelPool.Put(ch)

	// 声明队列
	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	var exchange string
	if isDelayed {
		exchange = conf.Message.DelayedExchange
	} else {
		exchange = conf.Message.DefaultExchange
	}
	err = ch.QueueBind(
		queue,
		queue, // 使用队列名称作为路由键
		exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s: %v", queue, exchange, err)
	}

	return nil
}

// CreateConsumer 创建消费者
func CreateConsumer(queue string, handler func(Message) error) (*Consumer, error) {
	ch, err := GetChannel()
	if err != nil {
		return nil, fmt.Errorf("channel is nil")
	}

	consumer := &Consumer{
		Queue:    queue,
		Handler:  handler,
		Channel:  ch,
		Done:     make(chan struct{}),
		IsClosed: false,
	}

	// 设置QoS
	err = ch.Qos(
		conf.Message.PrefetchCount, // prefetch count
		0,                          // prefetch size
		false,                      // global
	)
	if err != nil {
		return nil, err
	}

	// 启动消费者
	deliveries, err := ch.Consume(
		queue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	// 注册消费者
	consumersMu.Lock()
	consumers[queue] = consumer
	consumersMu.Unlock()

	// 处理消息
	go func() {
		for {
			select {
			case <-consumer.Done:
				return
			case d := <-deliveries:
				var msg Message
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					d.Reject(false)
					continue
				}

				if err := handler(msg); err != nil {
					d.Reject(true) // 重新入队
					continue
				}

				d.Ack(false)
			}
		}
	}()

	return consumer, nil
}

// restartConsumers 重启所有消费者
func restartConsumers() {
	consumersMu.Lock()
	defer consumersMu.Unlock()

	for queue, consumer := range consumers {
		if !consumer.IsClosed {
			close(consumer.Done)
			if newConsumer, err := CreateConsumer(queue, consumer.Handler); err != nil {
				utils.SysError(fmt.Sprintf("Failed to restart consumer for queue %s: %v", queue, err))
			} else {
				consumers[queue] = newConsumer
			}
		}
	}
}

// GracefulClose 优雅关闭RabbitMQ连接
func GracefulClose(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	// 关闭所有消费者
	consumersMu.Lock()
	defer consumersMu.Unlock() // 确保在函数结束时解锁

	for _, consumer := range consumers {
		if !consumer.IsClosed {
			close(consumer.Done) // 关闭通道
			consumer.IsClosed = true
		}
	}

	if channelPool != nil {
		if err := channelPool.Close(); err != nil {
			utils.SysError(fmt.Sprintf("Error closing channel pool: %v", err)) // 记录错误
			return fmt.Errorf("error closing channel pool: %v", err)
		}
		channelPool = nil
	}

	if conn != nil {
		if err := conn.Close(); err != nil {
			utils.SysError(fmt.Sprintf("Error closing connection: %v", err)) // 记录错误
			return fmt.Errorf("error closing connection: %v", err)
		}
		conn = nil
	}

	return nil
}
