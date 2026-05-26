package dto

import "warehouse/pkg/models"

type Category struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
type SuccessCategoryResponse struct {
	Data models.Category `json:"data"`
}
type SuccessCategoriesResponse struct {
	Data []models.Category `json:"data"`
}

// ErrorResponse wraps error responses
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
