package client

import (
	feedclient "trading-bot/internal/client/feed_client"
	"trading-bot/internal/client/publisher"
	"trading-bot/internal/models"
)

var (
	registry *ClientRegistry
)

type ClientRegistry struct {
	FeedClient feedclient.IFeedClient
	Publisher  publisher.IPublisher
}

func GetRegistry() *ClientRegistry {
	if registry == nil {
		registry = &ClientRegistry{}
	}
	return registry
}

func (cr *ClientRegistry) WithPublisher(brokers []string) *ClientRegistry {
	cr.Publisher = publisher.GetPublisher(brokers)
	return cr
}

func (cr *ClientRegistry) WithFeedClient(feedclientName feedclient.FeedClientName, bus *models.Bus) *ClientRegistry {
	cr.FeedClient = feedclient.FeedClientFactory(feedclientName, bus)
	return cr
}
