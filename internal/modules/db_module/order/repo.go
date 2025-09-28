package order

import (
	"trading-bot/internal/models"

	"gorm.io/gorm"
)

var (
	OrderRepoInstance    *OrderRepo
	NewOrderRepoInstance = func(db *gorm.DB) *OrderRepo {
		if OrderRepoInstance == nil {
			OrderRepoInstance = &OrderRepo{
				BaseRepo: models.BaseRepo{
					Db: db,
				},
			}
		}
		return OrderRepoInstance
	}
)

type IOrderRepo interface {
	CreateOrder(order *models.Order) error
	UpdateOrderStatus(order *models.Order) error
	CreatePosition(position *models.Position) error
	UpdatePosition(position *models.Position) error
}

type OrderRepo struct {
	models.BaseRepo
}

func (o *OrderRepo) CreateOrder(order *models.Order) error {
	return o.Db.Create(order).Error
}

func (o *OrderRepo) UpdateOrderStatus(order *models.Order) error {
	return o.Db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(order).Error
}

func (o *OrderRepo) CreatePosition(position *models.Position) error {
	return o.Db.Create(position).Error
}

func (o *OrderRepo) UpdatePosition(position *models.Position) error {
	return o.Db.Model(&models.Position{}).Where("id = ?", position.Id).Updates(position).Error
}
