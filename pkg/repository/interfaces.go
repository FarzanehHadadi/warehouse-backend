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
type DepartmentRepository interface {
	Create(department *models.Department) (*models.Department, error)
	FindByID(departmentId uint) (*models.Department, error)
	Delete(departmentId uint) error
	Update(departmentId uint, department *models.DepartmentUpdate) error
	GetList() ([]models.Department, error)
}
type ManagerRepository interface {
	Create(Manager *models.Manager) (*models.Manager, error)
	FindByID(managerId uint) (*models.Manager, error)
	Delete(managerId uint) error
	Update(managerId uint, Manager *models.ManagerUpdate) error
	GetList() ([]models.Manager, error)
}

// Repository holds all repositories
type Repository struct {
	User       UserRepository
	Category   CategoryRepository
	Unit       UnitRepository
	Department DepartmentRepository
	Manager    ManagerRepository
}
