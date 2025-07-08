package main

import (
	"CollabDoc-go/core"
	"CollabDoc-go/global"
	"CollabDoc-go/initialize"
	"context"
)

func main() {
	ctx := context.Background()
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.Minio = initialize.InitMinio()
	global.Kafka = initialize.InitKafka()
	global.Mongo = initialize.InitMongo()
	// Kafka 消费者（异步）
	initialize.StartDiffConsumer(ctx)

	defer global.Redis.Close()

	core.RunServer()

}
