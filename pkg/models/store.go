package models

type Store struct {
	Basic
	Name      string   `json:"name" gorm:"required;min=3;max=100;uniqueIndex"`
	ManagerID uint     `json:"manager_id" gorm:"required;index"`
	Manager   *Manager `json:"manager" gorm:"foreignKey:ManagerID"`
	IsActive  bool     `json:"is_active" gorm:"default:true"`
}
type StoreUpdate struct {
	Name      string `json:"name" gorm:"required;min=3;max=100;uniqueIndex" binding:"omitempty,min=3,max=100"`
	ManagerID uint   `json:"manager_id" gorm:"required" binding:"omitempty"`
	IsActive  *bool  `json:"is_active" binding:"omitempty"`
}
