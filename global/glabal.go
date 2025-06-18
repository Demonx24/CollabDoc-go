package global

import (
	"CollabDoc-go/config"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-redis/redis"
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
)
