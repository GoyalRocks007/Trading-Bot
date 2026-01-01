package simplestrategy

import (
	"sync"
	"trading-bot/internal/models"
)

type SimpleStrategy struct {
	RiskPerShare float64
	RiskPerTrade float64
	RRR          float64
}

var (
	simpleStrategyInstance *SimpleStrategy
	once                   sync.Once
)

func NewSimpleStrategy(riskPerShare float64, riskPerTrade float64, rrr float64) *SimpleStrategy {
	once.Do(func() {
		simpleStrategyInstance = &SimpleStrategy{
			RiskPerShare: riskPerShare,
			RiskPerTrade: riskPerTrade,
			RRR:          rrr,
		}
	})
	return simpleStrategyInstance
}

func (s *SimpleStrategy) CalculateRiskPerShare(price float64) float64 {
	return s.RiskPerShare
}

func (s *SimpleStrategy) CalculatePositionSize(price float64, capital float64) int {
	riskPerTrade := capital * (s.RiskPerTrade / 100) * 0.1 // Only 10% of capital
	riskPerShare := price * (s.RiskPerShare / 100)

	qty := int(riskPerTrade / riskPerShare)

	if qty < 1 {
		qty = 1
	}
	return qty
}

func (s *SimpleStrategy) CalculateStopLossAndTarget(price float64, sig models.Side) (float64, float64) {
	switch sig {
	case models.Buy:
		risk := price * (s.RiskPerShare / 100)
		stop := price - risk
		target := price + risk*s.RRR
		return stop, target
	case models.Sell:
		risk := price * (s.RiskPerShare / 100)
		stop := price + risk
		target := price - risk*s.RRR
		return stop, target
	default:
		return -1, -1
	}
}
