package utils

import (
	"CollabDoc-go/global"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"os"
	"path/filepath"
)

func DownloadAndSaveFromMinio(objectKey, savePath string) error {
	ctx := context.Background()
	bucket := global.Config.Minio.Bucket

	// 从 MinIO 获取对象
	obj, err := global.Minio.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("从 MinIO 获取对象失败: %v", err)
	}
	defer obj.Close()

	// 创建保存目录
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建本地文件
	f, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer f.Close()

	// 拷贝数据流到本地文件
	_, err = io.Copy(f, obj)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
