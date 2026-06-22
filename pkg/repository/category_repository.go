package repository

import (
	"context"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {

	return &categoryRepository{db: db}
}
func (cr *categoryRepository) Create(cat *models.Category) (*models.Category, error) {
	result := cr.db.Create(&cat)
	if result.Error != nil {
		switch {
		case isDuplicateKeyError(result.Error):
			return nil, ErrDuplicateKey
		default:
			return nil, result.Error
		}
	}
	return cat, nil
}
func (cr *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var cat *models.Category
	result := cr.db.First(&cat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return cat, nil
}
func (cr *categoryRepository) Delete(id uint) error {

	result := cr.db.Delete(&models.Category{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
func (cr *categoryRepository) Update(id uint, cat *models.Category) error {

	result := cr.db.Model(&models.Category{}).Where("id = ?", id).Updates(cat)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (cr *categoryRepository) GetList(req filter.Request) ([]*models.Category, *filter.CursorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := cr.db.WithContext(ctx).Model(&models.Category{})

	query, err := filter.Apply(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}

	var categories []*models.Category
	if err := query.Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	hasMore := len(categories) > limit
	if hasMore {
		categories = categories[:limit]
	}

	var nextCursor string
	if len(categories) > 0 {
		last := categories[len(categories)-1]
		if hasMore {
			nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
		}
	}

	return categories, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
