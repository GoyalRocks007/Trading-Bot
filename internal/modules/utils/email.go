package utils

import (
	"context"
	"trading-bot/internal/client/publisher"
)

type NotificationEvent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEvent(event *NotificationEvent) {
	client := publisher.GetPublisher([]string{"localhost:9092"})
	defer client.Close()

	client.Publish(context.Background(), "notifications", "", event)
}
