package dto

import "warehouse/pkg/models"

type ThresholdProximityListResponse struct {
	Items      []models.ThresholdProximityReport `json:"items"`
	NextCursor string                            `json:"next_cursor" example:""`
	HasMore    bool                              `json:"has_more" example:"true"`
	Limit      int                               `json:"limit" example:"20"`
}

type StoreProductQuantityListResponse struct {
	Items      []models.StoreProductQuantityReport `json:"items"`
	NextCursor string                              `json:"next_cursor" example:""`
	HasMore    bool                                `json:"has_more" example:"true"`
	Limit      int                                 `json:"limit" example:"20"`
}
