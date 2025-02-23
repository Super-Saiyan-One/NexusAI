package router

import (
	"net/http"
	controller2 "nexus-ai/api/kling/controller"
)

func NewRouter(videoCtrl *controller2.VideoController, imageCtrl *controller2.ImageController) http.Handler {
	mux := http.NewServeMux()

	// 视频相关路由
	mux.HandleFunc("POST /v1/videos", videoCtrl.CreateVideoTask)
	mux.HandleFunc("GET /v1/videos/{taskId}", videoCtrl.GetVideoTask)

	// 图像相关路由
	mux.HandleFunc("POST /v1/images", imageCtrl.CreateImageTask)
	mux.HandleFunc("GET /v1/images/{taskId}", imageCtrl.GetImageTask)
	mux.HandleFunc("GET /v1/images", imageCtrl.ListImageTasks)

	return mux
}
