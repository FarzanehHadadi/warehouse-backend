package repository

import (
	"context"
	"warehouse/pkg/api/filter"
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
	if department.ManagerID != nil {
		var count int64
		dr.db.Model(&models.Manager{}).Where("id = ?", *department.ManagerID).Count(&count)
		if count == 0 {
			return nil, ErrNotFound // or custom ErrManagerNotFound
		}
	}
	result := dr.db.Create(&department)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey

		} else {
			return nil, result.Error
		}
	}
	dr.db.Preload("Manager").First(department, department.ID)
	return department, nil

}

func (dr *departmentRepository) FindByID(departmentId uint) (*models.Department, error) {
	var department *models.Department
	result := dr.db.Preload("Manager").First(&department, departmentId)
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
func (dr *departmentRepository) Update(departmentId uint, department *models.DepartmentUpdate) error {
	// Prevent changing ID
	result := dr.db.Model(&models.Department{}).Where("id = ?", departmentId).Updates(department)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (dr *departmentRepository) GetList(req filter.Request) ([]*models.Department, *filter.CursorResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := dr.db.WithContext(ctx).Model(&models.Department{})

	query, err := filter.Apply(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.SimpleFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query = query.Preload("Manager")
	var departments []*models.Department
	if err := query.Find(&departments).Error; err != nil {
		return nil, nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	hasMore := len(departments) > limit
	if hasMore {
		departments = departments[:limit]
	}

	var nextCursor string
	if len(departments) > 0 {
		last := departments[len(departments)-1]
		if hasMore {
			nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
		} else {
			nextCursor = ""
		}
	}

	return departments, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil

}
