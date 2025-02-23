package main

import (
	"context"
	"fmt"
	"log"
	"nexus-ai/api/luma/service"
	"time"

	"nexus-ai/api/luma/model"
)

func main() {
	// 初始化服务（自动加载配置）
	svc, err := service.NewVideoService()
	if err != nil {
		log.Fatalf("Service initialization failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 构建请求参数
	resolution := model.Res4k
	duration := model.Duration5s
	req := model.VideoGenerateRequest{
		Prompt:     "In a fresh and beautiful style, it shows a beautiful natural landscape in the mountains. At the beginning of the video, the camera slowly zooms in from the misty peaks in the distance, and the sun shines through the clouds and sprinkles on the green mountains, with light and shadow interlaced. Then the screen switches to a clear stream, the stream is gurgling, the pebbles at the bottom of the water are clearly visible, and occasionally a few small fish swim by. Along the stream, shoot the unknown wild flowers blooming by the stream, colorful, and bees busy collecting honey in the flowers.\nThen the camera moves up to show the dense woods, the tall trees cover the sky, the sun shines through the gaps in the leaves, and the birds in the forest chirp, and occasionally you can see squirrels jumping between the branches. Continue to walk up the mountain and shoot a spectacular scene of a waterfall pouring down from a high place, with water splashing and mist filling the air, bringing a hint of coolness.\nThen the camera switches to the top of the mountain to capture the vast view. The mountains in the distance are rolling and connected to the sky. There are a few white clouds floating in the sky, just like a beautiful painting. Finally, the picture slowly zooms out to show the whole natural landscape of the mountains. With the soft and soothing background music, it gives people a sense of tranquility and relaxation, as if they are in this beautiful nature.",
		Model:      model.ModelRay2,
		Resolution: &resolution,
		Duration:   &duration,
	}

	// 提交生成请求
	genResp, err := svc.GenerateVideo(ctx, req)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	fmt.Printf("🚀 Generation started! ID: %s\n", genResp.ID)
	fmt.Printf("⏳ Initial state: %s\n", genResp.State)

	// 状态轮询
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Fatal("⏰ Operation timed out")
		case <-ticker.C:
			status, err := svc.GetStatus(ctx, genResp.ID)
			if err != nil {
				log.Fatalf("🔧 Status check failed: %v", err)
			}

			fmt.Printf("\n📊 Current state: %s", status.State)
			switch status.State {
			case "completed":
				if status.Assets != nil && status.Assets.Video != nil {
					fmt.Printf("\n🎉 Video ready: %s\n", *status.Assets.Video)
				}
				return
			case "failed":
				if status.FailureReason != nil {
					log.Fatalf("\n❌ Failed: %s", *status.FailureReason)
				}
				return
			}
		}
	}
}
