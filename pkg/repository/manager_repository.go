package repository

import (
	"context"
	"errors"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type managerRepository struct {
	db *gorm.DB
}

func NewManagerRepository(db *gorm.DB) ManagerRepository {
	return &managerRepository{db: db}
}

func (mr *managerRepository) Create(manager *models.Manager) (*models.Manager, error) {

	result := mr.db.Create(manager)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey
		}
		return nil, result.Error
	}

	// Preload departments after creation (optional but useful)
	mr.db.Preload("Departments").First(manager, manager.ID)

	return manager, nil
}
func (mr *managerRepository) FindByID(managerId uint) (*models.Manager, error) {
	var manager *models.Manager
	err := mr.db.Preload("Departments").First(&manager, managerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return manager, nil
}
func (mr *managerRepository) Delete(managerId uint) error {
	result := mr.db.Delete(&models.Manager{}, managerId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (mr *managerRepository) Update(managerId uint, manager *models.ManagerUpdate) error {
	result := mr.db.Model(&models.Manager{}).Where("id = ?", managerId).Updates(manager)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (mr *managerRepository) GetList(req filter.Request) ([]*models.Manager, *filter.CursorResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := mr.db.WithContext(ctx).Model(&models.Manager{})

	query, err := filter.Apply(query, req, filter.ManagerFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.ManagerFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query = query.Preload("Departments")
	var managers []*models.Manager
	if err := query.Find(&managers).Error; err != nil {
		return nil, nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	hasMore := len(managers) > limit
	if hasMore {
		managers = managers[:limit]
	}

	var nextCursor string
	if len(managers) > 0 {
		last := managers[len(managers)-1]
		if hasMore {
			nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
		} else {
			nextCursor = ""
		}
	}

	return managers, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
