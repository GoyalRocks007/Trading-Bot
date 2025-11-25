package utils

import (
	"context"
	kafkaclient "trading-bot/internal/client/kafka_client"
)

type NotificationEvent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEvent(event *NotificationEvent) {
	client := kafkaclient.GetKafkaClient([]string{"localhost:9092"})
	defer client.Close()

	client.Publish(context.Background(), "notifications", "", event)
}
