package repository

import (
	"context"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (ur *unitRepository) Create(unit *models.Unit) (*models.Unit, error) {
	result := ur.db.Create(&unit)
	if result.Error != nil {
		switch {
		case isDuplicateKeyError(result.Error):
			return nil, result.Error
		default:
			return nil, result.Error
		}

	}
	return unit, nil
}

func (ur *unitRepository) FindByID(unitId uint) (*models.Unit, error) {
	var unit *models.Unit
	result := ur.db.Where(&models.Unit{Id: unitId}).First(&unit)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return unit, nil
}
func (ur *unitRepository) Delete(unitId uint) error {
	result := ur.db.Delete(&models.Unit{}, unitId)

	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error

}
func (ur *unitRepository) Update(unitId uint, unit *models.Unit) error {
	result := ur.db.Model(&models.Unit{}).Where(&models.Unit{Id: unitId}).Updates(unit)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (ur *unitRepository) GetList() ([]models.Unit, error) {
	var units []models.Unit
	ctx, cancelFunc := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancelFunc()
	result := ur.db.WithContext(ctx).Find(&units)
	if result.Error != nil {
		return nil, result.Error
	}
	if units == nil {
		return []models.Unit{}, nil
	}
	return units, nil

}
