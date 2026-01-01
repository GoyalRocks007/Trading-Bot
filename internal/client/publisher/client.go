package publisher

import (
	"context"
	kafkaclient "trading-bot/internal/client/publisher/kafka_publisher"
)

func GetPublisher(brokers []string) IPublisher {
	return kafkaclient.GetKafkaClient(brokers)
}

type IPublisher interface {
	Publish(ctx context.Context, topic, key string, message any) error
	Close() error
}
