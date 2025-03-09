package main

import (
	"fmt"
	"log"
	"net/http"
	"nexus-ai/api/kling/api"
	"nexus-ai/api/kling/config"
	controller2 "nexus-ai/api/kling/controller"
	"nexus-ai/api/kling/router"
	service2 "nexus-ai/api/kling/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("kling/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化依赖
	apiClient := api.NewClient(&cfg.KlingAI)
	videoService := service2.NewVideoService(apiClient)
	controllerService := service2.NewImageService(apiClient)

	// 注入控制器（自动类型匹配）
	videoCtrl := controller2.NewVideoController(videoService)
	imageCtrl := controller2.NewImageController(controllerService)

	// 创建路由
	router := router.NewRouter(videoCtrl, imageCtrl)

	// 启动服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	log.Printf("Server starting on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
