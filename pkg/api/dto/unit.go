package dto

import "warehouse/pkg/models"

type Unit struct {
	Name string `json:"name" binding:"required, min=1, max=100"`
}

type SuccessUnitResponse struct {
	Data models.Unit `json:"data"`
}
type SuccessUnitListResponse struct {
	Items      []models.Unit `json:"items"`
	NextCursor string        `json:"next_cursor" example:""`
	HasMore    bool          `json:"has_more" example:"true"`
	Limit      int           `json:"limit" example:"20"`
}
type Summary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
