package kafkaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"trading-bot/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	writer *kafka.Writer
}

func GetKafkaClient(brokers []string) *KafkaClient {
	return &KafkaClient{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{}, // distributes load
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (k *KafkaClient) Publish(ctx context.Context, topic, key string, message any) error {
	valueBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: valueBytes,
		Time:  time.Now(),
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		logger.Log.Error(fmt.Sprintf("❌ Failed to publish to topic %s: %v", topic, err))
		return err
	}

	logger.Log.Info(fmt.Sprintf("📤 Published to topic=%s key=%s value=%s", topic, key, string(valueBytes)))
	return nil
}

func (k *KafkaClient) Close() error {
	return k.writer.Close()
}
