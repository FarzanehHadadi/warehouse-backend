package dto

type ProductSummary struct {
	ID               uint    `json:"id" example:"1"`
	Name             string  `json:"name" example:"product1"`
	WarningThreshold int     `json:"warning_threshold" example:"1" `
	CategoryID       uint    `json:"category_id" example:"2"`
	UnitID           uint    `json:"unit_id" example:"3"`
	Category         Summary `json:"category"`
	Unit             Summary `json:"unit"`
}

type productListResponse struct {
	Data       []ProductSummary `json:"items"`
	NextCursor string           `json:"next_cursor" example:""`
	HasMore    bool             `json:"has_more" example:"true"`
	Limit      int              `json:"limit" example:"20"`
}
