package order

import (
	"trading-bot/internal/models"

	"gorm.io/gorm"
)

var (
	orderCoreInstance    *OrderCore
	NewOrderCoreInstance = func(db *gorm.DB) *OrderCore {
		if orderCoreInstance == nil {
			orderCoreInstance = &OrderCore{
				repo: NewOrderRepoInstance(db),
			}
		}
		return orderCoreInstance
	}
)

type IOrderCore interface {
	CreateOrder(order *models.Order) error
	UpdateOrder(order *models.Order) error
	CreatePosition(position *models.Position) error
	UpdatePosition(position *models.Position) error
}

type OrderCore struct {
	repo IOrderRepo
}

func (o *OrderCore) CreateOrder(order *models.Order) error {
	return o.repo.CreateOrder(order)
}

func (o *OrderCore) UpdateOrder(order *models.Order) error {
	return o.repo.UpdateOrderStatus(order)
}

func (o *OrderCore) CreatePosition(position *models.Position) error {
	return o.repo.CreatePosition(position)
}

func (o *OrderCore) UpdatePosition(position *models.Position) error {
	return o.repo.UpdatePosition(position)
}
