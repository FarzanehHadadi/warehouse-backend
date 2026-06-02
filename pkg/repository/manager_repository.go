package repository

import (
	"context"
	"errors"
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

func (mr *managerRepository) GetList() ([]*models.Manager, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	var managers []*models.Manager
	result := mr.db.WithContext(ctx).Preload("Departments").Find(&managers)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	total := result.RowsAffected

	return managers, total, nil
}
