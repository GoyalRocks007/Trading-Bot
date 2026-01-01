package models

import "sync"

type Bus struct {
	Ticks         chan Tick
	PositionTicks chan Tick
	Candles       chan Candle
	Orders        chan *Order // risk → execution
	Fills         chan Fill   // execution → store/metrics
	OrderUpdates  chan *OrderEvent
	Equity        float64
	remCap        float64
	mu            sync.RWMutex
}

func NewBus(equity float64) *Bus {
	return &Bus{
		Ticks:         make(chan Tick, 15000),
		PositionTicks: make(chan Tick, 15000),
		Orders:        make(chan *Order, 256),
		Fills:         make(chan Fill, 256),
		Candles:       make(chan Candle, 5000),
		Equity:        equity,
		remCap:        equity,
	}
}

func (b *Bus) GetRemCap() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.remCap
}

func (b *Bus) UpdateRemCap(diff float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	cur := b.remCap
	b.remCap = cur + diff
}
