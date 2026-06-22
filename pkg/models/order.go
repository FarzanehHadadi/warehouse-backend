package models

import "time"

type OrderType string

const (
	OrderTypeInbound  OrderType = "inbound"
	OrderTypeOutbound OrderType = "outbound"
)

type ProductStatus string

const (
	ProductStatusGood      ProductStatus = "good"
	ProductStatusDefective ProductStatus = "defective"
	ProductStatusUnknown   ProductStatus = "unknown"
)

type Order struct {
	Basic
	ProductID     uint
	Product       *Product
	StoreID       uint
	Store         *Store
	Type          OrderType //enum :inbound,outbound
	Quantity      int
	Price         int
	ExpireDate    time.Time
	ProductStatus ProductStatus //good,defective,unknown
	DepartmentId  uint
	Department    *Department
	Description   string
}
type OrderUpdate struct{}
