package bot

import (
	"context"
	"sync"
	"time"
	feedclient "trading-bot/internal/client/feed_client"
	"trading-bot/internal/models"
	executionmodule "trading-bot/internal/modules/execution_module"
	ingressmodule "trading-bot/internal/modules/ingress_module"
	riskmodule "trading-bot/internal/modules/risk_module"
	tradestrategymodule "trading-bot/internal/modules/trade_strategy_module"

	"gorm.io/gorm"
)

type Bot struct {
	TotalCapital   float64
	RemCapital     float64
	Strategy       tradestrategymodule.ITradeStrategyModule
	Execution      executionmodule.IExecutionStrategy
	FeedClient     feedclient.IFeedClient
	IngressModule  ingressmodule.IIngressModule
	Bus            *models.Bus
	Risk           riskmodule.IRiskModule
	OpenPositions  map[string]*models.Position
	wg             *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	candleDuration time.Duration
}

func NewBot(strategy tradestrategymodule.TradeStrategyName, execution executionmodule.ExecutionStrategyName, bus *models.Bus, risk riskmodule.RiskStrategyName, db *gorm.DB, candleInterval time.Duration) *Bot {
	riskInstance := riskmodule.RiskFactory(risk)
	strategyInstance := tradestrategymodule.TradeStrategyFactory(strategy, riskInstance)
	executionInstance := executionmodule.ExecutionFactory(execution, db)
	feedclientInstance := feedclient.FeedClientFactory(feedclient.ZERODHA, bus)
	ingressmoduleInstance := ingressmodule.GetIngressModule()

	return &Bot{
		Strategy:       strategyInstance,
		Execution:      executionInstance,
		FeedClient:     feedclientInstance,
		IngressModule:  ingressmoduleInstance,
		Bus:            bus,
		Risk:           riskInstance,
		OpenPositions:  map[string]*models.Position{},
		candleDuration: candleInterval,
	}
}

func (b *Bot) Start() error {
	b.wg = &sync.WaitGroup{}
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.FeedClient.Start(b.wg)
	b.IngressModule.Runner(b.ctx, b.Bus, b.candleDuration, b.wg)
	b.Strategy.Runner(b.ctx, b.Bus, b.wg)
	b.Execution.Runner(b.ctx, b.Bus, b.wg)

	b.wg.Wait()
	return nil
}

func (b *Bot) Stop() error {
	b.FeedClient.Stop()
	b.cancel()
	return nil
}
