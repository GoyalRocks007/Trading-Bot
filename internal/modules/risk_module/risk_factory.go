package riskmodule

import simplestrategy "trading-bot/internal/modules/risk_module/strategy/simple_strategy"

type RiskStrategyName string

const (
	SIMPLE RiskStrategyName = "SIMPLE"
)

func RiskFactory(riskStrategyName RiskStrategyName) IRiskModule {
	switch riskStrategyName {
	case SIMPLE:
		return simplestrategy.NewSimpleStrategy(0.5, 0.5, 1)
	default:
		return simplestrategy.NewSimpleStrategy(0.5, 0.5, 1)
	}
}
