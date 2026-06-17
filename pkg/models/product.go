package models

type Product struct {
	Basic
	Name             string   `json:"name" binding:"required,min=3,max=100" gorm:"not null;size:100;"`
	WarningThreshold int      `json:"warning_threshold" binding:"required,min=1" gorm:"not null;default:10"`
	Unit             Unit     `json:"unit" gorm:"foreignKey:UnitID"`
	Category         Category `json:"category" gorm:"foreignKey:CategoryID"`
	CategoryID       uint     `json:"category_id" binding:"required" gorm:"not null;index:idx_category_id"`
	UnitID           uint     `json:"unit_id" binding:"required" gorm:"not null;index:idx_unit_id"`
}

type ProductUpdate struct {
	Name             string `json:"name" binding:"required,min=3,max=100"`
	WarningThreshold int    `json:"warning_threshold" binding:"required,min=1" `
	CategoryID       uint   `json:"category_id" binding:"required"`
	UnitID           uint   `json:"unit_id" binding:"required"`
}
