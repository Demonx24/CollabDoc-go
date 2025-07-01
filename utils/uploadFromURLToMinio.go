package utils

import (
	"CollabDoc-go/global"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
)

func UploadFromURLToMinio(fileURL, objectKey string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("URL请求失败，状态码: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = global.Minio.PutObject(
		context.Background(),
		global.Config.Minio.Bucket,
		objectKey,
		resp.Body,
		resp.ContentLength,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return fmt.Errorf("上传到MinIO失败: %v", err)
	}

	log.Printf("✅ 文件已从 URL 上传至 MinIO：%s", objectKey)
	return nil
}
