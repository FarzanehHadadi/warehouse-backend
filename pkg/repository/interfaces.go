package repository

import (
	"warehouse/pkg/models"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	// Add more methods as needed
}
type CategoryRepository interface {
	Create(cat *models.Category) (*models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Delete(id uint) error
	Update(id uint, cat *models.Category) error
	GetList() ([]models.Category, error)
}

// Repository holds all repositories
type Repository struct {
	User     UserRepository
	Category CategoryRepository
}
