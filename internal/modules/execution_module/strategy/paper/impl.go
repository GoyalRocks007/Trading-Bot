package paper

import (
	"fmt"
	"sync"
	"trading-bot/internal/models"
	"trading-bot/internal/modules/db_module/order"
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
	p.DeletePosSafe(position)
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

func (p *PaperStrategy) Runner(bus *models.Bus, wg *sync.WaitGroup) error {
	fmt.Println("execution strategy is running")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for o := range bus.Orders {
			if _, ok := p.ReadPosSafe(o.Symbol); ok {
				o.Status = models.Cancelled
				p.PlaceOrder(o)
				fmt.Printf("%s, order couldn't be placed due to already open order\n", o.Symbol)
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
				fmt.Printf("%s, order placed and position created\n", o.Symbol)
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
						bus.UpdateRemCap((float64(pos.Qty) * pos.Exit))
						remCap := bus.GetRemCap()
						fmt.Printf("%s position closed, rem cap is %f\n", t.Symbol, remCap)
					}

				case models.Sell:
					if t.LTP <= pos.Target || t.LTP >= pos.Stop {
						pos.Exit = t.LTP
						p.ClosePosition(pos)
						bus.UpdateRemCap((float64(pos.Qty) * (pos.Entry - pos.Exit)))
						remCap := bus.GetRemCap()
						fmt.Printf("%s position closed, rem cap is %f\n", t.Symbol, remCap)
					}
				}
			}
		}
	}()

	return nil
}
