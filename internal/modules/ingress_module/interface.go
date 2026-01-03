package ingressmodule

import (
	"context"
	"sync"
	"time"
	"trading-bot/internal/models"
)

type IIngressModule interface {
	Runner(ctx context.Context, bus *models.Bus, interval time.Duration, wg *sync.WaitGroup) error
}
