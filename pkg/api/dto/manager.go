package dto

import "warehouse/pkg/models"

type ManagerListResponse struct {
	Data []models.Manager `json:"data"`
}
type ManagerResponse struct {
	Data models.Manager `json:"data"`
}
type Manager struct {
	Name  string `json:"name" binding:"required,min=3, max=100" gorm:"not null;min:3;max:100;"`
	Phone string `json:"phone" gorm:"min:7;max:11" binding:"omitempty,min=7,max=11"`
}
