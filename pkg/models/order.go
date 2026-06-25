package models

import (
	"encoding/json"
	"fmt"
)

type OrderType string

const (
	OrderTypeInbound  OrderType = "inbound"
	OrderTypeOutbound OrderType = "outbound"
)

func (t OrderType) IsValid() bool {
	switch t {
	case OrderTypeInbound, OrderTypeOutbound:
		return true
	default:
		return false
	}
}

func (t *OrderType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*t = OrderType(s)
	if !t.IsValid() {
		return fmt.Errorf("invalid order type %q, must be one of: inbound, outbound", s)
	}

	return nil
}

func (t OrderType) MarshalJSON() ([]byte, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid order type %q", t)
	}
	return json.Marshal(string(t))
}

type ProductStatus string

const (
	ProductStatusGood      ProductStatus = "good"
	ProductStatusDefective ProductStatus = "defective"
	ProductStatusUnknown   ProductStatus = "unknown"
)

func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusGood, ProductStatusDefective, ProductStatusUnknown:
		return true
	default:
		return false
	}
}

func (s *ProductStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = ProductStatus(v)
	if !s.IsValid() {
		return fmt.Errorf("invalid product status %q, must be one of: good, defective, unknown", v)
	}

	return nil
}

func (s ProductStatus) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid product status %q", s)
	}
	return json.Marshal(string(s))
}

type Order struct {
	Basic
	ProductID     uint          `json:"product_id" binding:"required" gorm:"not null;index"`
	Product       *Product      `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	StoreID       uint          `json:"store_id" binding:"required" gorm:"not null;index"`
	Store         *Store        `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Type          OrderType     `json:"type" binding:"required,oneof=inbound outbound" gorm:"type:varchar(20);not null"`
	Quantity      int           `json:"quantity" binding:"required,min=1" gorm:"not null"`
	Price         int           `json:"price" binding:"required,min=0" gorm:"not null"`
	ExpireDate    Date          `json:"expire_date" binding:"required" gorm:"type:date;not null"`
	ProductStatus ProductStatus `json:"product_status" binding:"required,oneof=good defective unknown" gorm:"type:varchar(20);not null"`
	DepartmentID  uint          `json:"department_id" binding:"required" gorm:"not null;index"`
	Department    *Department   `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Description   *string       `json:"description" binding:"omitempty"`
}

type OrderUpdate struct {
	ProductID     *uint          `json:"product_id" binding:"omitempty"`
	StoreID       *uint          `json:"store_id" binding:"omitempty"`
	DepartmentID  *uint          `json:"department_id" binding:"omitempty"`
	Type          *OrderType     `json:"type" binding:"omitempty,oneof=inbound outbound"`
	Quantity      *int           `json:"quantity" binding:"omitempty,min=1"`
	Price         *int           `json:"price" binding:"omitempty,min=0"`
	ExpireDate    *Date          `json:"expire_date" binding:"omitempty" gorm:"type:date"`
	ProductStatus *ProductStatus `json:"product_status" binding:"omitempty,oneof=good defective unknown"`
	Description   *string        `json:"description" binding:"omitempty"`
}
