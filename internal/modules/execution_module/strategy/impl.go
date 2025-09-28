package paper

import (
	"sync"
	"trading-bot/internal/models"
	"trading-bot/internal/modules/db_module/order"
)

type PaperStrategy struct {
	OpenPositions map[string]*models.Position
	OrderService  order.IOrderCore
}

func NewPaperStrategy(orderService order.IOrderCore) *PaperStrategy {
	return &PaperStrategy{
		OrderService: orderService,
	}
}

func (p *PaperStrategy) PlaceOrder(order *models.Order) error {
	err := p.OrderService.CreateOrder(order)
	if err != nil {
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
		return err
	}
	delete(p.OpenPositions, position.Symbol)
	return nil
}

func (p *PaperStrategy) SquareOff() error {
	for _, pos := range p.OpenPositions {
		err := p.ClosePosition(pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PaperStrategy) Runner(bus *models.Bus, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for o := range bus.Orders {
			if _, ok := p.OpenPositions[o.Symbol]; ok {
				o.Status = models.Cancelled
				p.PlaceOrder(o)
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
				p.OpenPositions[o.Symbol] = pos
				p.OrderService.CreatePosition(pos)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range bus.PositionTicks {
			if pos, ok := p.OpenPositions[t.Symbol]; ok {
				switch pos.Side {
				case models.Buy:
					if t.LTP >= pos.Target || t.LTP <= pos.Stop {
						pos.Exit = t.LTP
						p.ClosePosition(pos)
					}
				case models.Sell:
					if t.LTP <= pos.Target || t.LTP >= pos.Stop {
						pos.Exit = t.LTP
						p.ClosePosition(pos)
					}
				}
			}
		}
	}()

	return nil
}
