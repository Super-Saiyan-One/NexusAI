package mq

import (
	"fmt"
	"nexus-ai/utils"
	"sync"
	"testing"
	"time"
)

// BenchmarkPublishConsume 测试消息发布和消费的性能
func BenchmarkPublishConsume(b *testing.B) {
	// 创建队列
	queues := []string{"publish_current", "publish_concurrent", "publish_delay", "publish_condelay"}
	for _, queue := range queues {
		if queue == "publish_delay" || queue == "publish_condelay" {
			CreateQueue(queue, true)
		} else {
			CreateQueue(queue, false)
		}
		utils.SysInfo("RabbitMQ Benchmarking in queue: " + queue)
	}

	// 逐个发布即时消息到 publish_current 队列
	benchPublishCurrent(b)
	utils.SysInfo("Publish current benchmark completed")

	// 批量发布即时消息到 publish_concurrent 队列
	benchBulkPublishConcurrent(b)
	utils.SysInfo("Bulk publish concurrent benchmark completed")

	// 逐个发布延时消息到 publish_delay 队列
	benchDelayedMessage(b)
	utils.SysInfo("Delayed message benchmark completed")

	// 批量发布延时消息到 publish_condelay 队列
	benchBulkDelayedPublish(b)
	utils.SysInfo("Bulk delayed publish benchmark completed")

	// 为上述四个队列生成多个消费者，来并发消费消息
	benchParallelConsume(b, queues)
	utils.SysInfo("Parallel consume benchmark completed")
}

func benchPublishCurrent(b *testing.B) {
	queue := "publish_current"
	msg := Message{
		Type: "test",
		Data: map[string]interface{}{
			"value": "test message",
		},
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := PublishMessage("", queue, msg); err != nil {
			b.Error(err)
			return
		}
	}
}

func benchBulkPublishConcurrent(b *testing.B) {
	queue := "publish_concurrent"
	batchSize := 100
	b.ResetTimer()
	for i := 0; i < b.N; i += batchSize {
		for j := 0; j < batchSize && i+j < b.N; j++ {
			msg := Message{
				Type: "test",
				Data: map[string]interface{}{
					"value": fmt.Sprintf("bulk message %d", i+j),
				},
				CreatedAt: time.Now(),
			}

			if err := PublishMessage("", queue, msg); err != nil {
				b.Error(err)
				return
			}
		}
	}
}

func benchDelayedMessage(b *testing.B) {
	queue := "publish_delay"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := Message{
			Type: "test",
			Data: map[string]interface{}{
				"value": "delayed message",
			},
			DelayedTime: time.Now().Add(5 * time.Second),
			CreatedAt:   time.Now(),
		}

		if err := PublishMessage("", queue, msg); err != nil {
			b.Error(err)
			return
		}
	}
}

func benchBulkDelayedPublish(b *testing.B) {
	queue := "publish_condelay"
	batchSize := 100
	b.ResetTimer()
	for i := 0; i < b.N; i += batchSize {
		for j := 0; j < batchSize && i+j < b.N; j++ {
			msg := Message{
				Type: "test",
				Data: map[string]interface{}{
					"value": fmt.Sprintf("bulk delayed message %d", i+j),
				},
				DelayedTime: time.Now().Add(10 * time.Second),
				CreatedAt:   time.Now(),
			}

			if err := PublishMessage("", queue, msg); err != nil {
				b.Error(err)
				return
			}
		}
	}
}

func benchParallelConsume(b *testing.B, queues []string) {
	received := 0
	numConsumers := 4
	var mu sync.Mutex

	// 为每个队列创建多个消费者
	for _, queue := range queues {
		for i := 0; i < numConsumers; i++ {
			_, err := CreateConsumer(queue, func(msg Message) error {
				mu.Lock()
				defer mu.Unlock()
				received++
				return nil
			})
			if err != nil {
				b.Error(err)
				return
			}
		}
	}

	// 等待消息被消费
	timeout := time.After(10 * time.Second)
	for received < b.N {
		select {
		case <-timeout:
			b.Error("timeout waiting for messages to be consumed")
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
