package tradestrategymodule

import (
	"time"
	riskmodule "trading-bot/internal/modules/risk_module"
	orbstrategy "trading-bot/internal/modules/trade_strategy_module/strategy/orb"
)

type TradeStrategyName string

const (
	ORB TradeStrategyName = "ORB"
)

func TradeStrategyFactory(tradeStrategyName TradeStrategyName, riskmodule riskmodule.IRiskModule) ITradeStrategyModule {
	switch tradeStrategyName {
	case ORB:
		return orbstrategy.NewORB(time.Now(), 15*time.Minute, riskmodule)
	default:
		return orbstrategy.NewORB(time.Now(), 15*time.Minute, riskmodule)
	}
}
