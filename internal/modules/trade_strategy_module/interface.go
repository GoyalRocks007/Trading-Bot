package tradestrategymodule

import (
	"sync"
	"trading-bot/internal/models"
)

type ITradeStrategyModule interface {
	Core(c models.Candle) *models.Order
	Runner(bus *models.Bus, wg *sync.WaitGroup)
}
