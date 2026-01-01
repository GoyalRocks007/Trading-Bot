package executionmodule

import (
	"trading-bot/internal/modules/db_module/order"
	paper "trading-bot/internal/modules/execution_module/strategy/paper"

	"gorm.io/gorm"
)

type ExecutionStrategyName string

const (
	PAPER ExecutionStrategyName = "PAPER"
)

func ExecutionFactory(executionStrategyName ExecutionStrategyName, db *gorm.DB) IExecutionStrategy {
	switch executionStrategyName {
	case PAPER:
		return paper.NewPaperStrategy(order.NewOrderCoreInstance(db))
	default:
		return paper.NewPaperStrategy(order.NewOrderCoreInstance(db))
	}
}
