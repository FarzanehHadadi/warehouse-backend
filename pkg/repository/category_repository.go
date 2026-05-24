package repository

import (
	"context"
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
func (cr *categoryRepository) GetList() ([]models.Category, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancelFunc()
	var categories []models.Category
	result := cr.db.WithContext(ctx).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	if categories == nil {
		return []models.Category{}, nil
	}
	return categories, nil

}
