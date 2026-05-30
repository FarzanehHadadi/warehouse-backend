package dto

import "warehouse/pkg/models"

type Unit struct {
	Name string `json:"name" binding:"required, min=1, max=100"`
}

type SuccessUnitResponse struct {
	Data models.Unit `json:"data"`
}
type SuccessUnitListResponse struct {
	Data []models.Unit `json:"data"`
}
