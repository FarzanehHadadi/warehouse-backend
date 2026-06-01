package dto

import "warehouse/pkg/models"

type Department struct {
	Name      string `json:"name" binding:"min=1,max=100, required"`
	ManagerId uint   `json:"manager_id" `
}

type DepartmentResponse struct {
	Data models.Department `json:"data"`
}
type DepartmentListResponse struct {
	Data []models.Department `json:"data"`
}
