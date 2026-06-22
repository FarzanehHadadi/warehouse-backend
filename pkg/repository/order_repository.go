package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (dr *orderRepository) Create(order *models.Order) (*models.Order, error) {

	return order, nil

}

func (dr *orderRepository) FindByID(orderId uint) (*models.Order, error) {

	return nil, nil
}
func (dr *orderRepository) Delete(orderId uint) error {

	return nil
}
func (dr *orderRepository) Update(orderId uint, order *models.OrderUpdate) error {

	return nil
}
func (dr *orderRepository) GetList(req filter.Request) ([]*models.Order, *filter.CursorResponse, error) {

	return nil, nil, nil

}
