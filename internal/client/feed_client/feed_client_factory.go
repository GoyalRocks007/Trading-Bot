package feedclient

import (
	zerodhafeedclient "trading-bot/internal/client/feed_client/zerodha"
	"trading-bot/internal/models"
)

type FeedClientName string

const (
	ZERODHA FeedClientName = "ZERODHA"
)

func FeedClientFactory(feedClientName FeedClientName, bus *models.Bus) IFeedClient {
	switch feedClientName {
	case ZERODHA:
		return zerodhafeedclient.GetZerodhaFeedClient(bus)
	default:
		return zerodhafeedclient.GetZerodhaFeedClient(bus)
	}
}
