package tradestrategymodule

import (
	"context"
	"sync"
	"trading-bot/internal/models"
)

type ITradeStrategyModule interface {
	Core(c models.Candle) *models.Order
	Runner(ctx context.Context, bus *models.Bus, wg *sync.WaitGroup)
}
