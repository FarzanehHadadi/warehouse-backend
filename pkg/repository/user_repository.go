package repository

import (
	"warehouse/pkg/models"
	"warehouse/pkg/utils"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	return nil, nil
}

func (r *userRepository) FindByPhone(phone string, password string) (*models.User, error) {
	var user *models.User
	result := r.db.Where("phone = ?", phone).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
