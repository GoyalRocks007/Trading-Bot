package paper

import (
	"context"
	"sync"
	"trading-bot/internal/models"
	"trading-bot/internal/modules/db_module/order"
	"trading-bot/internal/modules/utils"
	"trading-bot/logger"

	"go.uber.org/zap"
)

type PaperStrategy struct {
	OpenPositions map[string]*models.Position
	OrderService  order.IOrderCore
	mu            sync.RWMutex
}

func NewPaperStrategy(orderService order.IOrderCore) *PaperStrategy {
	return &PaperStrategy{
		OrderService:  orderService,
		OpenPositions: make(map[string]*models.Position),
	}
}

func (p *PaperStrategy) PlaceOrder(order *models.Order) error {
	err := p.OrderService.CreateOrder(order)
	if err != nil {
		logger.Log.Error("order couldn't be placed", zap.Error(err))
		utils.SendEvent(&utils.NotificationEvent{
			Subject: "Order couldn't be placed",
			Body:    err.Error(),
			To:      "uddyan.goyal@gmail.com",
		})
		return err
	}
	return nil
}

func (p *PaperStrategy) CancelOrder(order *models.Order) error {
	return nil
}

func (p *PaperStrategy) UpdateOrder(order *models.Order) error {
	return nil
}

func (p *PaperStrategy) ClosePosition(position *models.Position) error {
	position.Status = models.Closed
	err := p.OrderService.UpdatePosition(position)
	if err != nil {
		logger.Log.Error("position couldn't be closed", zap.Error(err))
		utils.SendEvent(&utils.NotificationEvent{
			Subject: "Position couldn't be closed",
			Body:    err.Error(),
			To:      "uddyan.goyal@gmail.com",
		})
		return err
	}
	p.DeletePosSafe(position)
	return nil
}

func (p *PaperStrategy) SquareOff() error {
	p.mu.RLock()
	positions := make([]*models.Position, 0, len(p.OpenPositions))
	for _, pos := range p.OpenPositions {
		positions = append(positions, pos)
	}
	p.mu.RUnlock()
	for _, pos := range positions {
		err := p.ClosePosition(pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PaperStrategy) AddPosSafe(pos *models.Position) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.OpenPositions[pos.Symbol] = pos
}

func (p *PaperStrategy) DeletePosSafe(pos *models.Position) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.OpenPositions, pos.Symbol)
}

func (p *PaperStrategy) ReadPosSafe(symbol string) (*models.Position, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pos, ok := p.OpenPositions[symbol]
	return pos, ok
}

func (p *PaperStrategy) Runner(ctx context.Context, bus *models.Bus, wg *sync.WaitGroup) error {
	logger.Log.Info("execution strategy is running")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case o, ok := <-bus.Orders:
				if !ok {
					logger.Log.Info("order channel closed")
					return
				}

				if _, ok := p.ReadPosSafe(o.Symbol); ok {
					o.Status = models.Cancelled
					p.PlaceOrder(o)
					logger.Log.Warn("order couldn't be placed due to already open order", zap.String("symbol", o.Symbol))
				} else {
					o.Status = models.Executed
					p.PlaceOrder(o)
					pos := &models.Position{
						Symbol: o.Symbol,
						Side:   o.Side,
						Entry:  o.Entry,
						Stop:   o.Stop,
						Target: o.Target,
						Qty:    o.Qty,
						Status: models.Open,
					}
					p.AddPosSafe(pos)
					p.OrderService.CreatePosition(pos)
					if o.Side == models.Buy {
						bus.UpdateRemCap(-(float64(o.Qty) * o.Entry))
					}
					logger.Log.Info("Order placed and position created", zap.String("symbol", o.Symbol))
				}
			case <-ctx.Done():
				p.SquareOff()
				return
			}

		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case t, ok := <-bus.PositionTicks:
				if !ok {
					logger.Log.Info("position tick channel closed")
					return
				}
				if pos, ok := p.OpenPositions[t.Symbol]; ok {
					switch pos.Side {
					case models.Buy:
						if t.LTP >= pos.Target || t.LTP <= pos.Stop {
							pos.Exit = t.LTP
							p.ClosePosition(pos)
							bus.UpdateRemCap((float64(pos.Qty) * pos.Exit))
							remCap := bus.GetRemCap()
							logger.Log.Info("position closed side buy", zap.String("symbol", t.Symbol), zap.Float64("remCap", remCap))
						}

					case models.Sell:
						if t.LTP <= pos.Target || t.LTP >= pos.Stop {
							pos.Exit = t.LTP
							p.ClosePosition(pos)
							bus.UpdateRemCap((float64(pos.Qty) * (pos.Entry - pos.Exit)))
							remCap := bus.GetRemCap()
							logger.Log.Info("position closed side sell", zap.String("symbol", t.Symbol), zap.Float64("remCap", remCap))
						}
					}
				}

			}

		}
	}()

	return nil
}
