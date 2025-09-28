package models

type Bus struct {
	Ticks         chan Tick
	PositionTicks chan Tick
	Candles       chan Candle
	Orders        chan *Order // risk → execution
	Fills         chan Fill   // execution → store/metrics
	OrderUpdates  chan *OrderEvent
	Equity        float64
	remCap        float64
}

func NewBus(ticks chan Tick) *Bus {
	return &Bus{
		Ticks:  ticks,
		Orders: make(chan *Order, 256),
		Fills:  make(chan Fill, 256),
	}
}

func (b *Bus) GetRemCap() float64 {
	return b.remCap
}
