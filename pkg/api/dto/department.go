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
	Items      []models.Department `json:"items"`
	NextCursor string              `json:"next_cursor" example:""`
	HasMore    bool                `json:"has_more" example:"true"`
	Limit      int                 `json:"limit" example:"20"`
}
