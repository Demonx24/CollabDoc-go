package initialize

import (
	"CollabDoc-go/global"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"time"
)

func InitMongo() *mongo.Client {
	cfg := global.Config.Mongo
	if !cfg.Enabled {
		global.Log.Warn("MongoDB 未启用，跳过初始化")
		return nil
	}

	// 构造 Mongo URI
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	if cfg.Username != "" && cfg.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}
	clientOpts := options.Client().ApplyURI(uri)

	// 认证源、副本集、SSL
	if cfg.AuthSource != "" {
		clientOpts.SetAuth(options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.AuthSource,
		})
	}
	if cfg.ReplicaSet != "" {
		clientOpts.SetReplicaSet(cfg.ReplicaSet)
	}
	clientOpts.SetTLSConfig(nil) // 如果需要 SSL，可以设置 TLSConfig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		global.Log.Error("初始化 MongoDB 客户端失败", zap.Error(err))
		os.Exit(1)
	}

	if err := client.Ping(ctx, nil); err != nil {
		global.Log.Error("连接 MongoDB 失败", zap.Error(err))
		os.Exit(1)
	}

	global.Log.Info("MongoDB 客户端初始化成功", zap.String("uri", uri))
	return client
}
