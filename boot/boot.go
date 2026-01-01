package boot

import (
	"trading-bot/internal/client"
	feedclient "trading-bot/internal/client/feed_client"
	"trading-bot/internal/models"
)

func InitClientRegistery(feedClientName feedclient.FeedClientName, bus *models.Bus, brokers []string) *client.ClientRegistry {
	return client.GetRegistry().
		WithFeedClient(feedClientName, bus).
		WithPublisher(brokers)
}
