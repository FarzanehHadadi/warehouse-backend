package repository

import (
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type storeRepository struct {
	*BaseRepository[models.Store]
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{
		BaseRepository: NewBaseRepository[models.Store](db),
	}
}

func (sr *storeRepository) Create(store *models.Store) error {
	return sr.BaseRepository.Create(store)

}
func (sr *storeRepository) FindByID(storeId uint) (*models.Store, error) {
	return sr.BaseRepository.FindByID(storeId, "Manager")

}
func (sr *storeRepository) Delete(storeId uint) error {
	return sr.BaseRepository.Delete(storeId)

}
func (sr *storeRepository) Update(storeId uint, store *models.StoreUpdate) error {
	return sr.BaseRepository.Update(storeId, store)

}
func (sr *storeRepository) GetList(req filter.Request) ([]*models.Store, *filter.CursorResponse, error) {
	return sr.BaseRepository.GetList(req, filter.StoreFilterConfig, "Manager")
}
