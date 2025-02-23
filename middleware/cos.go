package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
	"nexus-ai/constant" // 替换成你的项目实际路径
)

// UploadFileToCOS 函数用于将指定 URL 的文件上传到腾讯云 COS 并返回文件 URL
func UploadFileToCOS(key, fileURL string) (string, error) {
	// 创建 COS 客户端
	baseURLStr := fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", constant.Bucket, constant.AppId, constant.Region)
	baseURL, err := url.Parse(baseURLStr)
	client := cos.NewClient(&cos.BaseURL{BucketURL: baseURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  constant.SecretID,
			SecretKey: constant.SecretKey,
		},
	})

	// 从 URL 获取文件内容
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to get file from URL: %v", err)
	}
	defer resp.Body.Close()

	// 上传文件到 COS
	_, err = client.Object.Put(context.Background(), key, resp.Body, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to COS: %v", err)
	}

	// 生成文件的访问 URL
	// 对于公有读的存储桶，可以直接使用以下方式生成 URL
	publicURL := fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com/%s", constant.Bucket, constant.AppId, constant.Region, key)
	return publicURL, nil
}
