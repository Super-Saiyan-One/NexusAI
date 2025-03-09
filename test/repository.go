package test

import (
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/repository/log"
	"nexus-ai/repository/worker"
)

func TestRepository(num int) {
	TestUserRepository(num)
	TestUserGroupRepository(num)
	TestTokenRepository(num)
	TestModelRepository(num)
	TestChannelRepository(num)
	TestModelGroupRepository(num)
	TestChannelGroupRepository(num)
	TestUsageRepository(num)
	TestPaymentRepository(num)
	TestQuotaRepository(num)
	TestBillingRepository(num)
	TestMessageSaveRepository(num)
	TestTaskVideoRepository(num)
	TestGatewayLogRepository(num)
	TestMasterLogRepository(num)
	TestMasterMySQLLogRepository(num)
	TestRedisLogRepository(num)
	TestRedisPersistLogRepository(num)
	TestRelayLogRepository(num)
	TestRequestLogRepository(num)
	TestWorkerLogRepository(num)
	TestWorkerMySQLLogRepository(num)
	TestWorkerClusterRepository(num)
	TestWorkerGroupRepository(num)
	TestWorkerNodeRepository(num)
}

func TestUserRepository(num int) {
	userRepo := repository.NewUserRepository(model.GetDB())
	userRepo.Benchmark(num)
}

func TestUserGroupRepository(num int) {
	userGroupRepo := repository.NewUserGroupRepository(model.GetDB())
	userGroupRepo.Benchmark(num)
}

func TestTokenRepository(num int) {
	tokenRepo := repository.NewTokenRepository(model.GetDB())
	tokenRepo.Benchmark(num)
}

func TestModelRepository(num int) {
	modelRepo := repository.NewModelRepository(model.GetDB())
	modelRepo.Benchmark(num)
}

func TestChannelRepository(num int) {
	channelRepo := repository.NewChannelRepository(model.GetDB())
	channelRepo.Benchmark(num)
}

func TestModelGroupRepository(num int) {
	modelGroupRepo := repository.NewModelGroupRepository(model.GetDB())
	modelGroupRepo.Benchmark(num)
}

func TestChannelGroupRepository(num int) {
	channelGroupRepo := repository.NewChannelGroupRepository(model.GetDB())
	channelGroupRepo.Benchmark(num)
}

func TestUsageRepository(num int) {
	usageRepo := repository.NewUsageRepository(model.GetDB())
	usageRepo.Benchmark(num)
}

func TestPaymentRepository(num int) {
	paymentRepo := repository.NewPaymentRepository(model.GetDB())
	paymentRepo.Benchmark(num)
}

func TestQuotaRepository(num int) {
	quotaRepo := repository.NewQuotaRepository(model.GetDB())
	quotaRepo.Benchmark(num)
}

func TestBillingRepository(num int) {
	billingRepo := repository.NewBillingRepository(model.GetDB())
	billingRepo.Benchmark(num)
}

func TestMessageSaveRepository(num int) {
	messageSaveRepo := repository.NewMessageSaveRepository(model.GetDB())
	messageSaveRepo.Benchmark(num)
}

func TestTaskVideoRepository(num int) {
	taskVideoRepo := repository.NewTaskVideoRepository(model.GetDB())
	taskVideoRepo.Benchmark(num)
}

func TestGatewayLogRepository(num int) {
	gatewayLogRepo := log.NewGatewayLogRepository(model.GetDB())
	gatewayLogRepo.Benchmark(num)
}

func TestMasterLogRepository(num int) {
	masterLogRepo := log.NewMasterLogRepository(model.GetDB())
	masterLogRepo.Benchmark(num)
}

func TestMasterMySQLLogRepository(num int) {
	masterMySQLLogRepo := log.NewMasterMySQLLogRepository(model.GetDB())
	masterMySQLLogRepo.Benchmark(num)
}

func TestRedisLogRepository(num int) {
	redisLogRepo := log.NewRedisLogRepository(model.GetDB())
	redisLogRepo.Benchmark(num)
}

func TestRedisPersistLogRepository(num int) {
	redisPersistLogRepo := log.NewRedisPersistLogRepository(model.GetDB())
	redisPersistLogRepo.Benchmark(num)
}

func TestRelayLogRepository(num int) {
	relayLogRepo := log.NewRelayLogRepository(model.GetDB())
	relayLogRepo.Benchmark(num)
}

func TestRequestLogRepository(num int) {
	requestLogRepo := log.NewRequestLogRepository(model.GetDB())
	requestLogRepo.Benchmark(num)
}

func TestWorkerLogRepository(num int) {
	workerLogRepo := log.NewWorkerLogRepository(model.GetDB())
	workerLogRepo.Benchmark(num)
}

func TestWorkerMySQLLogRepository(num int) {
	workerMySQLLogRepo := log.NewWorkerMySQLLogRepository(model.GetDB())
	workerMySQLLogRepo.Benchmark(num)
}

func TestWorkerClusterRepository(num int) {
	workerClusterRepo := worker.NewWorkerClusterRepository(model.GetDB())
	workerClusterRepo.Benchmark(num)
}

func TestWorkerGroupRepository(num int) {
	workerGroupRepo := worker.NewWorkerGroupRepository(model.GetDB())
	workerGroupRepo.Benchmark(num)
}

func TestWorkerNodeRepository(num int) {
	workerNodeRepo := worker.NewWorkerNodeRepository(model.GetDB())
	workerNodeRepo.Benchmark(num)
}
