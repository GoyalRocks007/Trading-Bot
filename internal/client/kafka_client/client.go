package kafkaclient

import "context"

type IKafkaClient interface {
	Publish(ctx context.Context, topic, key string, message any) error
	Close() error
}
