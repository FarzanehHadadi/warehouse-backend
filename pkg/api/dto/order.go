package dto

import (
	"warehouse/pkg/models"
)

type SimpleSummary struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Product 1"`
}
type OrderSummary struct {
	ID            uint                 `json:"id" example:"1"`
	Product       *SimpleSummary       `json:"product,omitempty"`
	Store         *SimpleSummary       `json:"store,omitempty"`
	Department    *SimpleSummary       `json:"department,omitempty"`
	Description   *string              `json:"description" example:"Stock intake"`
	Quantity      int                  `json:"quantity" example:"10"`
	Price         int                  `json:"price" example:"1000"`
	ExpireDate    models.Date          `json:"expire_date" swaggertype:"string" format:"date" example:"2026-12-31"`
	ProductStatus models.ProductStatus `json:"product_status" example:"good"`
	Type          models.OrderType     `json:"type" example:"inbound"`
	ProductID     uint                 `json:"product_id" example:"1"`
	StoreID       uint                 `json:"store_id" example:"1"`
	DepartmentID  uint                 `json:"department_id" example:"1"`
}

type orderListResponse struct {
	Items      []OrderSummary `json:"items"`
	NextCursor string         `json:"next_cursor" example:""`
	HasMore    bool           `json:"has_more" example:"true"`
	Limit      int            `json:"limit" example:"20"`
}

type CreateOrderRequest struct {
	ProductID     uint                 `json:"product_id" binding:"required" example:"1"`
	StoreID       uint                 `json:"store_id" binding:"required" example:"1"`
	DepartmentID  uint                 `json:"department_id" binding:"required" example:"1"`
	Type          models.OrderType     `json:"type" binding:"required,oneof=inbound outbound" enums:"inbound,outbound" example:"inbound"`
	Quantity      int                  `json:"quantity" binding:"required,min=1" example:"10"`
	Price         int                  `json:"price" binding:"required,min=0" example:"1000"`
	ExpireDate    models.Date          `json:"expire_date" binding:"required" swaggertype:"string" format:"date" example:"2026-12-31"`
	ProductStatus models.ProductStatus `json:"product_status" binding:"required,oneof=good defective unknown" enums:"good,defective,unknown" example:"good"`
	Description   *string              `json:"description" binding:"omitempty" example:"Stock intake"`
}

type UpdateOrderRequest struct {
	ProductID     *uint                 `json:"product_id" binding:"omitempty" example:"1"`
	StoreID       *uint                 `json:"store_id" binding:"omitempty" example:"1"`
	DepartmentID  *uint                 `json:"department_id" binding:"omitempty" example:"1"`
	Type          *models.OrderType     `json:"type" binding:"omitempty,oneof=inbound outbound" enums:"inbound,outbound" example:"outbound"`
	Quantity      *int                  `json:"quantity" binding:"omitempty,min=1" example:"5"`
	Price         *int                  `json:"price" binding:"omitempty,min=0" example:"500"`
	ExpireDate    *models.Date          `json:"expire_date" binding:"omitempty" swaggertype:"string" format:"date" example:"2026-12-31"`
	ProductStatus *models.ProductStatus `json:"product_status" binding:"omitempty,oneof=good defective unknown" enums:"good,defective,unknown" example:"defective"`
	Description   *string               `json:"description" binding:"omitempty" example:"Updated note"`
}
