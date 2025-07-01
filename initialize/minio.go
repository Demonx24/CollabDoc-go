package initialize

import (
	"CollabDoc-go/global"
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// InitMinio 初始化并返回 MinIO 客户端
func InitMinio() *minio.Client {
	minioCfg := global.Config.Minio

	// 创建 MinIO 客户端
	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKeyID, minioCfg.SecretAccessKey, ""),
		Secure: minioCfg.UseSSL,
	})
	if err != nil {
		global.Log.Error("初始化 MinIO 客户端失败", zap.Error(err))
		os.Exit(1)
	}

	// 测试连接（比如列桶或者Ping）
	ctx := context.Background()
	_, err = client.ListBuckets(ctx)
	if err != nil {
		global.Log.Error("连接 MinIO 失败", zap.Error(err))
		os.Exit(1)
	}

	global.Log.Info("MinIO 客户端初始化成功", zap.String("endpoint", minioCfg.Endpoint))
	return client
}
