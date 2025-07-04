package service

import (
	"CollabDoc-go/global"
	"github.com/Shopify/sarama"
)

type KafkaService struct{}

func (kafkaService *KafkaService) SendKafkaMessage(topic string, key string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	_, _, err := global.Kafka.SendMessage(msg)
	return err
}
