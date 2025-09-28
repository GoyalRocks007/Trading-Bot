package models

import (
	"context"
	"sync"
	"trading-bot/internal/models"
	executionmodule "trading-bot/internal/modules/execution_module"
	riskmodule "trading-bot/internal/modules/risk_module"
	tradestrategymodule "trading-bot/internal/modules/trade_strategy_module"

	"gorm.io/gorm"
)

type Bot struct {
	TotalCapital  float64
	RemCapital    float64
	Strategy      tradestrategymodule.ITradeStrategyModule
	Execution     executionmodule.IExecutionStrategy
	Bus           *models.Bus
	Risk          riskmodule.IRiskModule
	OpenPositions map[string]*models.Position
	wg            *sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewBot(totalCapital float64, strategy tradestrategymodule.TradeStrategyName, execution executionmodule.ExecutionStrategyName, bus *models.Bus, risk riskmodule.RiskStrategyName, db *gorm.DB) *Bot {
	riskInstance := riskmodule.RiskFactory(risk)
	strategyInstance := tradestrategymodule.TradeStrategyFactory(strategy, riskInstance)
	executionInstance := executionmodule.ExecutionFactory(execution, db)

	return &Bot{
		TotalCapital:  totalCapital,
		Strategy:      strategyInstance,
		Execution:     executionInstance,
		Bus:           bus,
		Risk:          riskInstance,
		OpenPositions: map[string]*models.Position{},
	}
}

func (b *Bot) Start() error {
	b.wg = &sync.WaitGroup{}
	b.ctx, b.cancel = context.WithCancel(context.Background())

	b.Strategy.Runner(b.Bus, b.wg)
	b.Execution.Runner(b.Bus, b.wg)

	b.wg.Wait()
	return nil
}
