package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type categoryRepository struct {
	*BaseRepository[models.Category]
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {

	return &categoryRepository{
		BaseRepository: NewBaseRepository[models.Category](db),
	}
}
func (cr *categoryRepository) Create(cat *models.Category) error {
	return cr.BaseRepository.Create(cat)

}
func (cr *categoryRepository) FindByID(id uint) (*models.Category, error) {
	return cr.BaseRepository.FindByID(id)
}
func (cr *categoryRepository) Delete(id uint) error {

	return cr.BaseRepository.Delete(id)
}
func (cr *categoryRepository) Update(id uint, cat *models.Category) error {
	return cr.BaseRepository.Update(id, cat)
}
func (cr *categoryRepository) GetList(req filter.Request) ([]*models.Category, *filter.CursorResponse, error) {
	return cr.BaseRepository.GetList(req, filter.SimpleFilterConfig)

}
