package repository

import (
	"context"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (dr *departmentRepository) Create(department *models.Department) (*models.Department, error) {
	result := dr.db.Create(&department)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey

		} else {
			return nil, result.Error
		}
	}

	return department, nil

}

func (dr *departmentRepository) FindByID(departmentId uint) (*models.Department, error) {
	var department *models.Department
	result := dr.db.First(&department, departmentId)
	if result.Error != nil {
		return nil, result.Error
	}

	return department, nil
}
func (dr *departmentRepository) Delete(departmentId uint) error {
	result := dr.db.Delete(&models.Department{}, departmentId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (dr *departmentRepository) Update(departmentId uint, department *models.Department) error {
	result := dr.db.Model(&models.Department{}).Where("id = ?", departmentId).Updates(department)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (dr *departmentRepository) GetList() ([]models.Department, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancelFunc()

	var departments []models.Department
	result := dr.db.WithContext(ctx).Find(&departments)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return []models.Department{}, nil
	}
	return departments, nil
}
