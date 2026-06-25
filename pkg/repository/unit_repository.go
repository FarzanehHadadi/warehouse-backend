package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type unitRepository struct {
	*BaseRepository[models.Unit]
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{
		BaseRepository: NewBaseRepository[models.Unit](db),
	}
}

func (ur *unitRepository) Create(unit *models.Unit) error {
	return ur.BaseRepository.Create(unit)

}

func (ur *unitRepository) FindByID(unitId uint) (*models.Unit, error) {
	return ur.BaseRepository.FindByID(unitId)

}
func (ur *unitRepository) Delete(unitId uint) error {
	return ur.BaseRepository.Delete(unitId)

}
func (ur *unitRepository) Update(unitId uint, unit *models.Unit) error {
	return ur.BaseRepository.Update(unitId, unit)

}
func (ur *unitRepository) GetList(req filter.Request) ([]*models.Unit, *filter.CursorResponse, error) {
	return ur.BaseRepository.GetList(req, filter.SimpleFilterConfig)
}
