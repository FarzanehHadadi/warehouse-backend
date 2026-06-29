package models

type ThresholdProximityReport struct {
	ProductID        uint   `json:"product_id" example:"1"`
	ProductName      string `json:"product_name" example:"Product A"`
	CategoryID       uint   `json:"category_id" example:"2"`
	CategoryName     string `json:"category_name" example:"Category A"`
	UnitID           uint   `json:"unit_id" example:"3"`
	UnitName         string `json:"unit_name" example:"kg"`
	WarningThreshold int    `json:"warning_threshold" example:"10"`
	CurrentQuantity  int    `json:"current_quantity" example:"5"`
	StoresCount      int    `json:"stores_count" example:"2"`
}

type StoreProductQuantityReport struct {
	ProductID     uint `json:"product_id" example:"1"`
	CategoryID    uint `json:"category_id" example:"2"`
	StoreID       uint `json:"store_id" example:"3"`
	TotalQuantity int  `json:"total_quantity" example:"50"`
}
