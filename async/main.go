package async

import (
	"nexus-ai/async/handler"
	"nexus-ai/constant"
)

func SetupAsync() error {
	err := setupQuotaExpire() // 设置配额过期异步处理队列
	if err != nil {
		return err
	}
	return nil
}

func setupQuotaExpire() error {
	quotaExpire := handler.NewQuotaExpire()
	err := quotaExpire.CreateQueue()
	if err != nil {
		return err
	}
	err = quotaExpire.CreateConsumers(constant.RabbitMQDefaultCustomerNum)
	if err != nil {
		return err
	}
	return nil
}
