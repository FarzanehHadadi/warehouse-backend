package repository

import "gorm.io/gorm"

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
