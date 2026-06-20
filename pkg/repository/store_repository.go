package repository

import (
	"errors"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/logger"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (sr *storeRepository) Create(store *models.Store) (*models.Store, error) {
	result := sr.db.Create(store)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey
		}
		logger.Log.Error(result.Error.Error())

		return nil, result.Error
	}
	sr.db.Preload("Manager").First(store, store.ID)
	return store, nil
}
func (sr *storeRepository) FindByID(storeId uint) (*models.Store, error) {
	var store *models.Store
	err := sr.db.Preload("Manager").First(&store, storeId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return store, nil

}
func (sr *storeRepository) Delete(storeId uint) error {
	result := sr.db.Delete(&models.Store{}, storeId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (sr *storeRepository) Update(storeId uint, store *models.StoreUpdate) error {
	result := sr.db.Model(&models.Store{}).Where("id = ?", storeId).Updates(store)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (sr *storeRepository) GetList(req filter.Request) ([]*models.Store, *filter.CursorResponse, error) {
	query := sr.db.Model(&models.Store{})
	query, err := filter.Apply(query, req, filter.StoreFilterConfig)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.StoreFilterConfig)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err
	}
	query = query.Preload("Manager")
	var storesList []*models.Store
	if err := query.Find(&storesList).Error; err != nil {
		logger.Log.Error(err.Error())
		return nil, nil, err

	}
	limit := req.Limit
	if limit < 0 {
		limit = 20
	}
	HasMore := len(storesList) > limit
	if HasMore {
		storesList = storesList[:limit]
	}
	var nextCursor string
	if len(storesList) > 0 {
		last := storesList[len(storesList)-1]
		nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)

	} else {
		nextCursor = ""
	}

	return storesList, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    HasMore,
	}, nil
}
