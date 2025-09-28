package executionmodule

import (
	"sync"
	"trading-bot/internal/models"
)

type IExecutionStrategy interface {
	PlaceOrder(order *models.Order) error
	// OnOrderUpdate(order *models.Order) error
	CancelOrder(order *models.Order) error
	// OnPositionUpdate(position *models.Position) error
	ClosePosition(position *models.Position) error
	Runner(bus *models.Bus, wg *sync.WaitGroup) error
	SquareOff() error
}
