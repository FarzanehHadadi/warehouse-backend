package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type orderRepository struct {
	*BaseRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		BaseRepository: NewBaseRepository[models.Order](db),
	}
}

func (or *orderRepository) Create(order *models.Order) error {
	return or.BaseRepository.Create(order)

}

func (or *orderRepository) FindByID(orderId uint) (*models.Order, error) {
	return or.BaseRepository.FindByID(orderId, "Product", "Store", "Department")

}
func (or *orderRepository) GetList(req filter.Request) ([]*models.Order, *filter.CursorResponse, error) {
	return or.BaseRepository.GetList(req, filter.OrderFilterConfig, "Store", "Department", "Product")

}
func (or *orderRepository) Delete(orderId uint) error {
	return or.BaseRepository.Delete(orderId)

}
func (or *orderRepository) Update(orderId uint, order *models.OrderUpdate) error {
	return or.BaseRepository.Update(orderId, order)

}

func (or *orderRepository) GetListNoPagination(req filter.Request) ([]*models.Order, error) {
	return or.BaseRepository.GetListNoPagination(req, filter.OrderFilterConfig, "Store", "Department", "Product")

}
