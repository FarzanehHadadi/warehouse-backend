package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type managerRepository struct {
	*BaseRepository[models.Manager]
}

func NewManagerRepository(db *gorm.DB) ManagerRepository {
	return &managerRepository{
		BaseRepository: NewBaseRepository[models.Manager](db),
	}
}

func (mr *managerRepository) Create(manager *models.Manager) error {

	return mr.BaseRepository.Create(manager)

}
func (mr *managerRepository) FindByID(managerId uint) (*models.Manager, error) {
	return mr.BaseRepository.FindByID(managerId, "Departments")

}
func (mr *managerRepository) Delete(managerId uint) error {
	return mr.BaseRepository.Delete(managerId)

}
func (mr *managerRepository) Update(managerId uint, manager *models.ManagerUpdate) error {
	return mr.BaseRepository.Update(managerId, manager)

}

func (mr *managerRepository) GetList(req filter.Request) ([]*models.Manager, *filter.CursorResponse, error) {

	return mr.BaseRepository.GetList(req, filter.SimpleFilterConfig, "Departments")

}
