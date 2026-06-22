package repository

import (
	"context"
	"warehouse/pkg/api/filter"
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
	result := ur.db.Where("id = ?", unitId).First(&unit)
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
	result := ur.db.Model(&models.Unit{}).Where("id = ?", unitId).Updates(unit)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (ur *unitRepository) GetList(req filter.Request) ([]*models.Unit, *filter.CursorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := ur.db.WithContext(ctx).Model(&models.Unit{})

	query, err := filter.Apply(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}

	var units []*models.Unit
	if err := query.Find(&units).Error; err != nil {
		return nil, nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	hasMore := len(units) > limit
	if hasMore {
		units = units[:limit]
	}

	var nextCursor string
	if len(units) > 0 {
		last := units[len(units)-1]
		if hasMore {
			nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
		}
	}

	return units, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
