// main.go
package main

import (
	"log"
	"net/http"
	"nexus-ai/api/hailuo/api"
	"nexus-ai/api/hailuo/auth"
	"nexus-ai/api/hailuo/controller"
	"nexus-ai/api/hailuo/service"
)

func main() {

	// 初始化依赖
	auth := auth.NewAuthenticator()
	client := api.NewVideoClient(auth)
	videoService := service.NewVideoService(client)
	videoController := controller.NewVideoController(videoService)

	// 配置路由
	http.Handle("POST /tasks", http.HandlerFunc(videoController.CreateVideoTask))
	http.Handle("GET /tasks/{task_id}", http.HandlerFunc(videoController.GetTaskStatus))
	http.Handle("POST /callbacks", http.HandlerFunc(videoController.HandleCallback))

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
