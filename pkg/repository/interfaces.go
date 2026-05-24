package repository

import "warehouse/pkg/models"

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	// Add more methods as needed
}

// Repository holds all repositories
type Repository struct {
	User UserRepository
}
