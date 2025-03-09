package main

import (
	"flag"
	"fmt"
	"net/http"
	"nexus-ai/constant"
	"nexus-ai/middleware"
	"nexus-ai/model"
	"nexus-ai/mq"
	"nexus-ai/mysql"
	"nexus-ai/redis"
	"nexus-ai/repository"
	"nexus-ai/test"
	"testing"

	"nexus-ai/router"
	"nexus-ai/utils"

	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envPath := utils.GetEnv("ENV_PATH", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		utils.FatalLog("Failed to load .env file: " + err.Error())
	}

	// 解析命令行参数
	redisTest := flag.Int("redis", 0, "是否执行Redis基准测试")
	mysqlTest := flag.Int("mysql", 0, "是否执行MySQL基准测试")
	rabbitmqTest := flag.Int("rabbitmq", 0, "是否执行RabbitMQ基准测试")
	flag.Parse()

	// 初始化服务日志
	utils.SetupLog()

	// 初始化MySQL
	if err := mysql.Setup(); err != nil {
		utils.FatalLog("MySQL | " + err.Error())
	}
	utils.SysInfo("MySQL setup completed")
	defer func() {
		if err := mysql.Shutdown(); err != nil {
			utils.SysError("MySQL | " + err.Error())
		}
	}()

	// 初始化Gorm
	if err := model.InitGorm(); err != nil {
		utils.FatalLog("Gorm | " + err.Error())
	}
	utils.SysInfo("Gorm setup completed")

	// 初始化Redis
	if err := redis.Setup(); err != nil {
		utils.FatalLog("Redis | " + err.Error())
	}
	utils.SysInfo("Redis setup completed")
	defer func() {
		if err := redis.Shutdown(); err != nil {
			utils.SysError("Redis | " + err.Error())
		}
	}()

	// 初始化RabbitMQ
	for i := 1; i > 0; i-- {
		utils.SysInfo(fmt.Sprintf("等待 %d 秒以初始化 RabbitMQ...", i))
		time.Sleep(1 * time.Second) // 每秒打印一次
	}
	if err := mq.Setup(); err != nil {
		utils.FatalLog("RabbitMQ | " + err.Error())
	}
	utils.SysInfo("RabbitMQ setup completed")
	defer func() {
		if err := mq.Shutdown(); err != nil {
			utils.SysError("RabbitMQ | " + err.Error())
		}
	}()

	// 执行MySQL基准测试
	if *mysqlTest > 0 {
		utils.SysInfo("执行MySQL基准测试")
		test.TestRepository(*mysqlTest)
		utils.SysInfo("MySQL基准测试完成")
	}

	// 执行root用户生成
	repository.RootGenerate(model.GetDB())
	utils.SysInfo("RootGenerate completed")

	// 执行Redis基准测试
	if *redisTest > 0 {
		utils.SysInfo("执行Redis基准测试")
		if err := redis.RunBenchmarks(); err != nil {
			utils.SysError("Redis benchmarks failed: " + err.Error())
		} else {
			utils.SysInfo("Redis benchmarks completed successfully")
		}
	}

	// 执行RabbitMQ基准测试
	if *rabbitmqTest > 0 {
		utils.SysInfo("执行RabbitMQ基准测试")
		b := &testing.B{N: *rabbitmqTest}
		mq.BenchmarkPublishConsume(b)
		utils.SysInfo("RabbitMQ基准测试完成")
	}

	// 设置gin模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.New()
	server.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		utils.SysError(fmt.Sprintf("panic detected: %v", err)) // 出现panic时，输出错误日志
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": fmt.Sprintf("panic detected: %v. Please submit an issue to %s", err, constant.GitRepoURL), // 输出错误信息，并提示用户提交issue
				"type":    constant.ErrorTypeInternalServerError,                                                     // nexus-ai panic error
			},
		})
	}))
	utils.SysInfo("GIN Config setup completed")

	server.Use(middleware.RequestIDGenerateMiddleware()) // 添加requestID中间件 生成requestID
	utils.SysInfo("Middleware setup completed")

	utils.SetupAPILog(server) // 为gin Engine添加日志服务 记录请求日志
	utils.SysInfo("API Log setup completed")

	router.SetupRouter(server)
	utils.SysInfo("Router setup completed")

	backendPort, _ := strconv.Atoi(utils.GetEnv("BACKEND_PORT", constant.BackendPort))
	utils.SysInfo("Server starting on port " + strconv.Itoa(backendPort))
	err = server.Run(":" + strconv.Itoa(backendPort))
	if err != nil {
		utils.FatalLog("Failed to start HTTP server: " + err.Error())
	}
}
