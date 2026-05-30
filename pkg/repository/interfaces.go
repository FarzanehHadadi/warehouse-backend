package repository

import (
	"warehouse/pkg/models"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByPhone(email string, password string) (*models.User, error)
	// Add more methods as needed
}
type CategoryRepository interface {
	Create(cat *models.Category) (*models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Delete(id uint) error
	Update(id uint, cat *models.Category) error
	GetList() ([]models.Category, error)
}
type UnitRepository interface {
	Create(unit *models.Unit) (*models.Unit, error)
	FindByID(unitId uint) (*models.Unit, error)
	Delete(unitId uint) error
	Update(unitId uint, unit *models.Unit) error
	GetList() ([]models.Unit, error)
}

// Repository holds all repositories
type Repository struct {
	User     UserRepository
	Category CategoryRepository
	Unit     UnitRepository
}
