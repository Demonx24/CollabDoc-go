package global

import (
	"CollabDoc-go/config"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      redis.Client
	BlackCache local_cache.Cache
	Mongo      *mongo.Client
	Minio      *minio.Client
	Kafka      sarama.SyncProducer
)
