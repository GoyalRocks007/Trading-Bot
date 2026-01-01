package orbstrategy

import (
	"fmt"
	"sync"
	"time"
	"trading-bot/internal/models"
	riskmodule "trading-bot/internal/modules/risk_module"
)

type ORB struct {
	start time.Time
	end   time.Time
	// per-symbol opening range
	orbHigh      map[string]float64
	orbLow       map[string]float64
	locked       map[string]bool // true once window ends
	RiskStrategy riskmodule.IRiskModule
}

func NewORB(dayStart time.Time, orbWindow time.Duration, riskStrategy riskmodule.IRiskModule) *ORB {
	return &ORB{
		start:        dayStart,
		end:          dayStart.Add(orbWindow),
		orbHigh:      map[string]float64{},
		orbLow:       map[string]float64{},
		locked:       map[string]bool{},
		RiskStrategy: riskStrategy,
	}
}

func (o *ORB) Core(c models.Candle) *models.Order {
	// During OR window: update range
	if !o.locked[c.Symbol] {
		if c.Start.Before(o.end) { // still inside window
			if c.High > o.orbHigh[c.Symbol] {
				o.orbHigh[c.Symbol] = c.High
			}
			if o.orbLow[c.Symbol] == 0 || c.Low < o.orbLow[c.Symbol] {
				o.orbLow[c.Symbol] = c.Low
			}
			return nil
		} else {
			o.locked[c.Symbol] = true
			fmt.Printf("Breakout range for %s: %f - %f\n", c.Symbol, o.orbLow[c.Symbol], o.orbHigh[c.Symbol])
		}
	}

	// After window: check breakouts (close outside range)
	if o.locked[c.Symbol] {
		// bullish breakout
		if c.Close > o.orbHigh[c.Symbol] && o.orbHigh[c.Symbol] > 0 {
			return &models.Order{Symbol: c.Symbol, Side: models.Buy, Entry: c.Close, Reason: "ORB_BREAKOUT_UP"}
		}
		// bearish breakout
		if c.Close < o.orbLow[c.Symbol] && o.orbLow[c.Symbol] > 0 {
			return &models.Order{Symbol: c.Symbol, Side: models.Sell, Entry: c.Close, Reason: "ORB_BREAKOUT_DOWN"}
		}
	}
	return nil
}

func (o *ORB) Runner(bus *models.Bus, wg *sync.WaitGroup) {
	fmt.Println("trade strategy is running")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for c := range bus.Candles {
			if ord := o.Core(c); ord != nil {
				stop, target := o.RiskStrategy.CalculateStopLossAndTarget(ord.Entry, ord.Side)
				qty := o.RiskStrategy.CalculatePositionSize(ord.Entry, bus.Equity)
				remCap := bus.GetRemCap()
				switch ord.Side {
				case models.Buy:
					if float64(qty)*ord.Entry > remCap {
						qty = int(remCap / ord.Entry)
					}
				}
				if qty > 0 && stop > 0 && target > 0 {
					ord.Stop = stop
					ord.Target = target
					ord.Qty = qty
					bus.Orders <- ord
				}
			}
		}
	}()
}
