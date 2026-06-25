package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type departmentRepository struct {
	*BaseRepository[models.Department]
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		BaseRepository: NewBaseRepository[models.Department](db),
	}
}

func (dr *departmentRepository) Create(department *models.Department) error {
	return dr.BaseRepository.Create(department)

}

func (dr *departmentRepository) FindByID(departmentId uint) (*models.Department, error) {
	return dr.BaseRepository.FindByID(departmentId, "Manager")
}
func (dr *departmentRepository) Delete(departmentId uint) error {
	return dr.BaseRepository.Delete(departmentId)

}
func (dr *departmentRepository) Update(departmentId uint, department *models.DepartmentUpdate) error {
	return dr.BaseRepository.Update(departmentId, department)

}
func (dr *departmentRepository) GetList(req filter.Request) ([]*models.Department, *filter.CursorResponse, error) {

	return dr.BaseRepository.GetList(req, filter.SimpleFilterConfig, "Manager")
}
