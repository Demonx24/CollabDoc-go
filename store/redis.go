package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis init failed: %v", err)
	}
}
