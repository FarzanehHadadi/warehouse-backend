package repository

import (
	"errors"
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
	return or.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			if isDuplicateKeyError(err) {
				return ErrDuplicateKey
			}
			return err
		}
		return or.adjustProductInventory(tx, order.ProductID, inventoryDelta(order))
	})
}

func (or *orderRepository) FindByID(orderId uint) (*models.Order, error) {
	return or.BaseRepository.FindByID(orderId, "Product", "Store", "Department")
}

func (or *orderRepository) GetList(req filter.Request) ([]*models.Order, *filter.CursorResponse, error) {
	return or.BaseRepository.GetList(req, filter.OrderFilterConfig, "Store", "Department", "Product")
}

func (or *orderRepository) Delete(orderId uint) error {
	return or.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Order
		if err := tx.First(&existing, orderId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := tx.Delete(&models.Order{}, orderId).Error; err != nil {
			return err
		}

		return or.adjustProductInventory(tx, existing.ProductID, -inventoryDelta(&existing))
	})
}

func (or *orderRepository) Update(orderId uint, order *models.OrderUpdate) error {
	return or.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Order
		if err := tx.First(&existing, orderId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		oldProductID, oldDelta := inventoryEffect(&existing)

		result := tx.Model(&models.Order{}).Where("id = ?", orderId).Updates(order)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}

		var updated models.Order
		if err := tx.First(&updated, orderId).Error; err != nil {
			return err
		}

		newProductID, newDelta := inventoryEffect(&updated)

		if oldProductID == newProductID {
			return or.adjustProductInventory(tx, newProductID, newDelta-oldDelta)
		}

		if err := or.adjustProductInventory(tx, oldProductID, -oldDelta); err != nil {
			return err
		}
		return or.adjustProductInventory(tx, newProductID, newDelta)
	})
}

func (or *orderRepository) GetListNoPagination(req filter.Request) ([]*models.Order, error) {
	return or.BaseRepository.GetListNoPagination(req, filter.OrderFilterConfig, "Store", "Department", "Product")
}

func inventoryDelta(order *models.Order) int {
	switch order.Type {
	case models.OrderTypeInbound:
		return order.Quantity
	case models.OrderTypeOutbound:
		return -order.Quantity
	default:
		return 0
	}
}

func inventoryEffect(order *models.Order) (productID uint, delta int) {
	return order.ProductID, inventoryDelta(order)
}

func (or *orderRepository) adjustProductInventory(tx *gorm.DB, productID uint, delta int) error {
	if delta == 0 {
		return nil
	}

	result := tx.Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("inventory_count", gorm.Expr("inventory_count + ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
