package ingressmodule

import (
	"context"
	"sync"
	"time"
	"trading-bot/internal/models"
	"trading-bot/logger"

	"go.uber.org/zap"
)

var (
	GetIngressModule = func() IIngressModule {
		return &IngressModule{}
	}
)

type IngressModule struct{}

func (im *IngressModule) Runner(ctx context.Context, bus *models.Bus, interval time.Duration, wg *sync.WaitGroup) error {

	out := bus.Candles
	tickCh := bus.Ticks

	logger.Log.Info("Starting Ingress Module....")

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		// ONE candle per symbol
		current := make(map[string]*models.Candle)

		for {
			select {
			case tick, ok := <-tickCh:
				if !ok {
					logger.Log.Info("tick channel closed")
					// Flush remaining candles on shutdown
					for _, c := range current {
						out <- *c
					}
					return
				}
				windowStart := tick.Time.Truncate(interval)
				windowEnd := windowStart.Add(interval)

				c, exists := current[tick.Symbol]

				// First tick for symbol OR new window
				if !exists || !tick.Time.Before(c.End) {
					// flush old candle if exists
					if exists {
						logger.Log.Info("emitting candle", zap.String("symbol", c.Symbol))
						out <- *c
					}

					// start new candle
					current[tick.Symbol] = &models.Candle{
						Symbol: tick.Symbol,
						Open:   tick.LTP,
						High:   tick.LTP,
						Low:    tick.LTP,
						Close:  tick.LTP,
						Volume: int64(tick.LastTradedQuantity),
						Start:  windowStart,
						End:    windowEnd,
					}
					continue
				}

				// Update existing candle
				if tick.LTP > c.High {
					c.High = tick.LTP
				}
				if tick.LTP < c.Low {
					c.Low = tick.LTP
				}
				c.Close = tick.LTP
				c.Volume += int64(tick.LastTradedQuantity)

			case <-ctx.Done():
				return
			}

		}
	}()

	return nil
}
