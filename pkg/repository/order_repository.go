package repository

import (
	"context"
	"errors"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/logger"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (or *orderRepository) Create(order *models.Order) (*models.Order, error) {
	result := or.db.Create(order)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey
		}
		logger.Log.Error(result.Error.Error())
		return nil, result.Error
	}
	or.db.Preload("Product").First(order, order.ID)
	or.db.Preload("Store").First(order, order.ID)
	or.db.Preload("Department").First(order, order.ID)
	return order, nil
}

func (or *orderRepository) FindByID(orderId uint) (*models.Order, error) {
	var order *models.Order
	err := or.db.Preload("Product").Preload("Store").Preload("Department").First(&order, orderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		logger.Log.Error(err.Error())
		return nil, err
	}
	return order, nil
}
func (or *orderRepository) GetList(req filter.Request) ([]*models.Order, *filter.CursorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := or.db.WithContext(ctx).Model(&models.Order{})
	query, err := filter.Apply(query, req, filter.OrderFilterConfig)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.OrderFilterConfig)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}

	query = query.Preload("Store")
	query = query.Preload("Department")
	query = query.Preload("Product")
	var orders []*models.Order
	if err := query.Find(&orders).Error; err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	hasMore := len(orders) > limit
	if hasMore {
		orders = orders[:limit]
	}
	var nextCursor string
	if len(orders) > 0 {
		last := orders[len(orders)-1]
		nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
	} else {
		nextCursor = ""
	}
	return orders, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
func (or *orderRepository) Delete(orderId uint) error {
	result := or.db.Delete(&models.Order{}, orderId)
	if result.Error != nil {
		logger.Log.Error(result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.Log.Error("not found")
		return ErrNotFound
	}
	return nil
}
func (or *orderRepository) Update(orderId uint, order *models.OrderUpdate) error {
	result := or.db.Model(&models.Order{}).Where("id = ?", orderId).Updates(order)
	if result.Error != nil {
		logger.Log.Error(result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
