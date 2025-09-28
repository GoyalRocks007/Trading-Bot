package riskmodule

import "trading-bot/internal/models"

type IRiskModule interface {
	CalculateRiskPerShare(price float64) float64
	CalculatePositionSize(price float64, capital float64) int
	CalculateStopLossAndTarget(price float64, sig models.Side) (float64, float64)
}
