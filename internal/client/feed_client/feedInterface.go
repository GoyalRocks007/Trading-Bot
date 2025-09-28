package feedclient

import (
	"context"
	"trading-bot/internal/models"
)

type Feed interface {
	Start(ctx context.Context) (<-chan models.Tick, error)
	Stop()
}
