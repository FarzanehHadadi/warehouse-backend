package dto

import "warehouse/pkg/models"

type Category struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
type SuccessCategoryResponse struct {
	Data models.Category `json:"data"`
}
type SuccessCategoriesResponse struct {
	Items      []models.Category `json:"items"`
	NextCursor string            `json:"next_cursor" example:""`
	HasMore    bool              `json:"has_more" example:"true"`
	Limit      int               `json:"limit" example:"20"`
}

// ErrorResponse wraps error responses
