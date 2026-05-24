package repository

import (
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
	return nil
}
func (cr *categoryRepository) Update(cat *models.Category) (*models.Category, error) {
	return nil, nil
}
