package repository

import (
	"warehouse/pkg/api/filter"
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
	GetList(req filter.Request) ([]*models.Category, *filter.CursorResponse, error)
}
type UnitRepository interface {
	Create(unit *models.Unit) (*models.Unit, error)
	FindByID(unitId uint) (*models.Unit, error)
	Delete(unitId uint) error
	Update(unitId uint, unit *models.Unit) error
	GetList(req filter.Request) ([]*models.Unit, *filter.CursorResponse, error)
}
type DepartmentRepository interface {
	Create(department *models.Department) (*models.Department, error)
	FindByID(departmentId uint) (*models.Department, error)
	Delete(departmentId uint) error
	Update(departmentId uint, department *models.DepartmentUpdate) error
	GetList(req filter.Request) ([]*models.Department, *filter.CursorResponse, error)
}
type ManagerRepository interface {
	Create(Manager *models.Manager) (*models.Manager, error)
	FindByID(managerId uint) (*models.Manager, error)
	Delete(managerId uint) error
	Update(managerId uint, Manager *models.ManagerUpdate) error
	GetList(req filter.Request) ([]*models.Manager, *filter.CursorResponse, error)
}
type ProductRepository interface {
	Create(product *models.Product) (*models.Product, error)
	FindByID(productId uint) (*models.Product, error)
	Delete(productId uint) error
	Update(productId uint, product *models.ProductUpdate) error
	GetList(req filter.Request) ([]*models.Product, *filter.CursorResponse, error)
}
type StoreRepository interface {
	Create(store *models.Store) (*models.Store, error)
	FindByID(storeId uint) (*models.Store, error)
	Delete(storeId uint) error
	Update(storeId uint, store *models.StoreUpdate) error
	GetList(req filter.Request) ([]*models.Store, *filter.CursorResponse, error)
}
type OrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
	FindByID(orderId uint) (*models.Order, error)
	Delete(orderId uint) error
	Update(orderId uint, store *models.OrderUpdate) error
	GetList(req filter.Request) ([]*models.Order, *filter.CursorResponse, error)
}

// Repository holds all repositories
type Repository struct {
	User       UserRepository
	Category   CategoryRepository
	Unit       UnitRepository
	Department DepartmentRepository
	Manager    ManagerRepository
	Product    ProductRepository
	Store      StoreRepository
	Order      OrderRepository
}
