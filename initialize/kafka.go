package initialize

import (
	"CollabDoc-go/global"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"os"
)

// InitKafka 初始化 Kafka 同步生产者客户端并返回
func InitKafka() sarama.SyncProducer {
	kafkaCfg := global.Config.Kafka
	if !kafkaCfg.Enabled {
		global.Log.Warn("Kafka 未启用，跳过初始化")
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	if kafkaCfg.Username != "" && kafkaCfg.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kafkaCfg.Username
		config.Net.SASL.Password = kafkaCfg.Password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, config)
	if err != nil {
		global.Log.Error("初始化 Kafka 生产者失败", zap.Error(err))
		os.Exit(1)
	}

	global.Log.Info("Kafka 生产者初始化成功", zap.Any("brokers", kafkaCfg.Brokers))
	return producer
}
