package handler

import (
	"fmt"
	"nexus-ai/constant"
	"nexus-ai/model"
	"nexus-ai/mq"
	"nexus-ai/repository"
	"nexus-ai/service"
	"time"
)

// quotaExpireMessageData 配额过期消息数据结构
type quotaExpireMessageData struct {
	QuotaID    string
	UserID     string
	ExpireTime time.Time
}

type QuotaExpire interface {
	CreateQueue() error
	CreateConsumers(consumerNum int) error
}

type quotaExpire struct{}

func NewQuotaExpire() QuotaExpire {
	return &quotaExpire{}
}

const (
	quotaExpireQueue = "quota_expire"
)

func (qe *quotaExpire) CreateQueue() error {
	err := mq.CreateQueue(quotaExpireQueue, true)
	if err != nil {
		return err
	}
	return nil
}

func (qe *quotaExpire) CreateConsumers(consumerNum int) error {
	for i := 0; i < consumerNum; i++ {
		_, err := mq.CreateConsumer(quotaExpireQueue, quotaExpireHandle)
		if err != nil {
			return err
		}
	}
	return nil
}

func PublishMessage(messageData quotaExpireMessageData) error {
	message := mq.Message{
		Data:        map[string]interface{}{"quota_id": messageData.QuotaID, "user_id": messageData.UserID},
		DelayedTime: messageData.ExpireTime,
	}
	return mq.PublishMessage(constant.RabbitMQMessageDelayedExchange, quotaExpireQueue, message)
}

func quotaExpireHandle(msg mq.Message) error {
	// 直接从map中获取数据，不需要json.Unmarshal
	quotaID, ok := msg.Data["quota_id"].(string)
	if !ok {
		return fmt.Errorf("invalid quota id")
	}

	userID, ok := msg.Data["user_id"].(string)
	if !ok {
		return fmt.Errorf("invalid user id")
	}

	quotaRepo := repository.NewQuotaRepository(model.GetDB())
	return service.QuotaExpireAsyncHandle(quotaRepo, quotaID, userID)
}
